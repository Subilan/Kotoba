package main

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func checkToken(c *gin.Context) {
	token := c.Request.Header.Get("Token")

	err := checkJWT(token)

	if err != nil {
		respond(c, 403, err.Error(), false)
		return
	}

	respond(c, 200, "OK", true)
}

func login(c *gin.Context) {
	var obj reqCommonLogin

	if err := c.BindJSON(&obj); err != nil {
		switch err.(type) {
		default:
			respond(c, 500, err.Error(), nil)
		case *json.UnmarshalTypeError:
			respond(c, 500, "Invalid argument type.", nil)
		}
		return
	}

	if len(obj.Password) == 0 || len(obj.Username) == 0 {
		respond(c, 400, "Not enough argument.", nil)
		return
	}

	// Check: if there is the user.
	res, dbErr := mongoGetOne("accounts", bson.M{
		"username": obj.Username,
	})

	if dbErr != nil {
		if dbErr == mongo.ErrNoDocuments {
			respond(c, 404, f("No such user named {0}.", obj.Username), nil)
			return
		}
		respond(c, 500, dbErr.Error(), nil)
		return
	}

	// Get data from the result.
	tg_username := res["username"].(string)
	tg_hash := res["hash"].(primitive.Binary).Data
	tg_avatar := res["avatar"].(string)
	tg_website := res["website"].(string)

	// Compare hash and password.
	hashErr := bcrypt.CompareHashAndPassword(tg_hash, []byte(obj.Password))

	if hashErr != nil {
		respond(c, 403, "Invalid credentials", nil)
		return
	}

	// Generate login token.
	token, tokenErr := genJWT(map[string]string{
		"username": tg_username,
		"website":  tg_website,
		"avatar":   tg_avatar,
	})

	if tokenErr != nil {
		respond(c, 500, tokenErr.Error(), nil)
		return
	}

	respond(c, 200, "", token)
}

func register(c *gin.Context) {
	var obj reqCommonRegister

	if err := c.BindJSON(&obj); err != nil {
		respond(c, 500, err.Error(), nil)
		return
	}

	// Check: if there is already a user named that.
	dupCount, countErr := mongoCount("accounts", bson.M{"username": obj.Username})

	if countErr != nil {
		respond(c, 500, countErr.Error(), nil)
		return
	}

	if dupCount != 0 {
		respond(c, 409, "Username is already in use.", nil)
		return
	}

	// Generate password hash.
	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		respond(c, 500, hashErr.Error(), nil)
		return
	}

	// Insert user into db.
	insErr := mongoInsertOne("accounts", bson.M{
		"username": obj.Username,
		"hash":     hashed,
		"avatar":   obj.Avatar,
		"website":  obj.Website,
		"created_at": time.Now().UnixMilli(),
	})

	if insErr != nil {
		respond(c, 500, insErr.Error(), nil)
		return
	}

	// Generate login token. This should be the same as that in the login process.
	token, tokenErr := genJWT(map[string]string{
		"username": obj.Username,
		"website":  obj.Website,
		"avatar":   obj.Avatar,
	})

	if tokenErr != nil {
		respond(c, 500, tokenErr.Error(), nil)
		return
	}

	respond(c, 200, f("Successfully created account {0}", obj.Username), token)
}

func deleteComment(c *gin.Context) {
	q_commentId := c.Param("commentId")

	if q_commentId == "" {
		respond(c, 500, "Not enough argument.", nil)
	}

	decoded, err := extractJWT(c.Request.Header.Get("Token"))

	if err != nil {
		respond(c, 403, err.Error(), nil)
		return
	}

	tg_username := decoded["username"].(string)
	tg_commentId := q_commentId

	res, getErr := mongoGetOne("comments", bson.M{
		"comment_id": tg_commentId,
	})

	if getErr != nil {
		respond(c, 500, getErr.Error(), nil)
		return
	}

	res_username := res["username"].(string)

	if res_username != tg_username {
		respond(c, 403, "You cannot delete comments that is not your own.", nil)
		return
	}

	_, delErr := mongoDeleteMany("comments", bson.M{
		"comment_id": tg_commentId,
	})

	if delErr != nil {
		respond(c, 500, delErr.Error(), nil)
		return
	}

	respond(c, 200, "", nil)
}

func createComment(c *gin.Context) {
	var obj reqCreateComment

	// Get username from JWT
	decoded, err := extractJWT(c.Request.Header.Get("Token"))

	if err != nil {
		respond(c, 403, err.Error(), nil)
		return
	}

	tg_username := decoded["username"].(string)

	if err := c.BindJSON(&obj); err != nil {
		respond(c, 500, err.Error(), nil)
		return
	}

	tg_text := obj.Text

	// Insert the comment into db
	insErr := mongoInsertOne("comments", bson.M{
		"username":   tg_username,
		"text":       tg_text,
		"comment_id": uuid.New().String(),
		"created_at": time.Now().UnixMilli(),
		"updated_at": time.Now().UnixMilli(),
	})

	if insErr != nil {
		respond(c, 500, insErr.Error(), nil)
		return
	}

	respond(c, 200, "Your comment has been sent!", nil)
}

func toggleReaction(c *gin.Context) {
	var obj reqReaction

	if err := c.BindJSON(&obj); err != nil {
		respond(c, 500, err.Error(), nil)
		return
	}

	decoded, err := extractJWT(c.Request.Header.Get("Token"))

	if err != nil {
		respond(c, 403, err.Error(), nil)
		return
	}

	tg_username := decoded["username"].(string)
	tg_commentId := obj.CommentId
	tg_emoji := obj.Emoji

	doc := bson.M{
		"emoji":      tg_emoji,
		"username":   tg_username,
		"comment_id": tg_commentId,
	}

	count, countErr := mongoCount("comment_reactions", doc)

	if countErr != nil {
		respond(c, 500, countErr.Error(), nil)
		return
	}

	if count > 0 {
		// The response is present, just delete it.
		_, delErr := mongoDeleteMany("comment_reactions", doc)
		if delErr != nil {
			respond(c, 500, delErr.Error(), nil)
			return
		}
		respond(c, 200, "Deleted target reaction.", nil)
	} else {
		doc["created_at"] = time.Now().UnixMilli()

		// The response is not present, create it.
		insErr := mongoInsertOne("comment_reactions", doc)

		if insErr != nil {
			respond(c, 500, insErr.Error(), nil)
			return
		}

		respond(c, 200, "Created target reaction.", nil)
	}

}

func getComments(c *gin.Context) {
	q_borderTimestamp := c.Query("border_timestamp")
	q_limit := c.Query("limit")
	tg_order := c.Query("order")

	if tg_order == "" {
		respond(c, 500, "Not enough argument.", nil)
		return
	}

	if tg_order != "desc" && tg_order != "asc" {
		respond(c, 500, "The `order` field in the request must be either `desc` or `asc`.", nil)
		return
	}

	var tg_borderTimestamp int64
	var tg_limit int64

	if q_borderTimestamp != "" {
		parseInt(q_borderTimestamp, &tg_borderTimestamp)
	}

	if q_limit != "" {
		parseInt(q_limit, &tg_limit)
	}

	var tg_order_num int
	filter := bson.M{}
	// desc: the newest at the top; asc: the oldest at the top
	if tg_order == "desc" {
		if tg_borderTimestamp != 0 {
			filter = bson.M{
				"created_at": bson.M{
					// when using desc, the comments we'll get are bound to have dates that are smaller than what is already present.
					"$lt": tg_borderTimestamp,
				},
			}
		}
		tg_order_num = -1
	} else {
		if tg_borderTimestamp != 0 {
			filter = bson.M{
				"created_at": bson.M{
					// when using asc, the comments we'll get are bound to have dates that are bigger than what is already present.
					"$gt": tg_borderTimestamp,
				},
			}
		}
		tg_order_num = 1
	}

	findOptions := options.Find()
	if tg_order_num != 0 {
		findOptions.SetSort(bson.M{"created_at": tg_order_num})
	}

	if tg_limit != 0 {
		findOptions.SetLimit(tg_limit)
	}

	res, err := mongoGetMany("comments", filter, findOptions)

	if err != nil {
		respond(c, 500, err.Error(), nil)
		return
	}

	for _, r := range res {
		user, err := mongoGetOne("accounts", bson.M{"username": r["username"]})

		if err != nil {
			respond(c, 500, err.Error(), nil)
			return
		}

		r["user_avatar"] = user["avatar"]
		r["user_website"] = user["website"]
	}

	respond(c, 200, f("Found {0} comments.", len(res)), res)
}

func updateComment(c *gin.Context) {
	var obj reqAlterComment

	if err := c.BindJSON(&obj); err != nil {
		respond(c, 500, err.Error(), nil)
		return
	}

	tg_commentId := obj.CommentId
	tg_text := obj.Text

	err := mongoUpdateOne("comments", bson.M{
		"comment_id": tg_commentId,
	}, bson.M{
		"$set": bson.M{
			"text": tg_text,
			"updated_at": time.Now().UnixMilli(),
		},
	})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			respond(c, 500, f("No such comment with id {0}", tg_commentId), nil)
		}
		respond(c, 500, err.Error(), nil)
		return
	}

	respond(c, 200, "Successfully altered the comment.", nil)
}

func deleteAccount(c *gin.Context) {
	decoded, decErr := extractJWT(c.Request.Header.Get("Token"))

	if decErr != nil {
		respond(c, 403, decErr.Error(), nil)
		return
	}

	_, delErr := mongoDeleteMany("accounts", bson.M{
		"username": decoded["username"],
	})

	if delErr != nil {
		respond(c, 500, delErr.Error(), nil)
		return
	}

	respond(c, 200, "Successfully deleted your account.", nil)
}

func updateAccount(c *gin.Context) {
	var obj reqAlterAccount

	decoded, err := extractJWT(c.Request.Header.Get("Token"))

	if err != nil {
		respond(c, 403, err.Error(), nil)
		return
	}

	tg_username := decoded["username"]
	tg_avatar := obj.Avatar
	tg_website := obj.Website
	tg_hash, hashErr := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		respond(c, 500, hashErr.Error(), nil)
		return
	}

	updErr := mongoUpdateOne("accounts", bson.M{
		"username": tg_username,
	}, bson.M{
		"$set": bson.M{
			"avatar":  tg_avatar,
			"website": tg_website,
			"hash":    tg_hash,
		},
	})

	if updErr != nil {
		respond(c, 500, updErr.Error(), nil)
		return
	}

	respond(c, 200, "Successfully altered user information.", nil)
}

func getReactions(c *gin.Context) {
	tg_commentid := c.Query("uid")

	if tg_commentid == "" {
		respond(c, 500, "Not enough argument.", nil)
		return
	}

	res, getErr := mongoGetMany("comment_reactions", bson.M{
		"comment_id": tg_commentid,
	}, options.Find())

	if getErr != nil {
		respond(c, 500, getErr.Error(), nil)
		return
	}

	respond(c, 200, f("Found {0} reactions for {1}", len(res), tg_commentid), res)
}