package main

import (
	"strconv"
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
		respond(c, 500, err.Error(), nil)
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
	var obj reqDeleteComment

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

	tg_commentId := obj.CommentId

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
		"created_at": time.Now().Format(time.DateTime),
	})

	if insErr != nil {
		respond(c, 500, insErr.Error(), nil)
		return
	}

	respond(c, 200, "Your comment has been sent!", nil)
}

func toggleReaction(c *gin.Context) {
	var obj reqReaction

	decoded, err := extractJWT(c.Request.Header.Get("Token"))

	if err != nil {
		respond(c, 403, err.Error(), nil)
		return
	}

	tg_username := decoded["username"].(string)
	tg_commentId := obj.CommentId
	tg_emoji := obj.Emoji

	doc := bson.M{
		"emoji":    tg_emoji,
		"username": tg_username,
		"comment_id":   tg_commentId,
	}

	count, countErr := mongoCount("comments", doc)

	if countErr != nil {
		respond(c, 500, countErr.Error(), nil)
		return
	}

	if count > 0 {
		// The response is present, just delete it.
		_, delErr := mongoDeleteMany("comments", doc)
		if delErr != nil {
			respond(c, 500, delErr.Error(), nil)
			return
		}
		respond(c, 200, "", nil)
	} else {
		doc["created_at"] = time.Now().Format(time.DateTime)

		// The response is not present, create it.
		insErr := mongoInsertOne("comments", doc)

		if insErr != nil {
			respond(c, 500, insErr.Error(), nil)
			return
		}

		respond(c, 200, "", nil)
	}

}

func getComments(c *gin.Context) {
	tg_borderTimestamp, parseErr1 := parseInt(c.Query("border_timestamp"))

	if parseErr1 != nil {
		if parseErr1 == strconv.ErrSyntax {
			respond(c, 500, "Invalid number for field `limit`.", nil)
		} else {
			respond(c, 500, parseErr1.Error(), nil)
		}
		return
	}

	tg_limit, parseErr2 := parseInt(c.Query("limit"))

	if parseErr2 != nil {
		if parseErr2 == strconv.ErrSyntax {
			respond(c, 500, "Invalid number for field `limit`.", nil)
		} else {
			respond(c, 500, parseErr2.Error(), nil)
		}
		return
	}

	tg_order := c.Query("order")

	if tg_order != "desc" && tg_order != "asc" {
		respond(c, 500, "The `order` field in the request must be either `desc` or `asc`.", nil)
		return
	}

	tg_borderTime := time.Unix(int64(tg_borderTimestamp), 0)
	var tg_operator string
	var tg_order_num int
	// desc: the newest at the top; asc: the oldest at the top
	if tg_order == "desc" {
		// when using desc, the comments we'll get are bound to have dates that are smaller than what is already present.
		tg_operator = "$lte"
		tg_order_num = -1
	} else {
		// when using asc, the comments we'll get are bound to have dates that are bigger than what is already present.
		tg_operator = "$gte"
		tg_order_num = 1
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"created_at": tg_order_num})
	findOptions.SetLimit(tg_limit)

	res, err := mongoGetMany("comments", bson.M{
		"created_at": bson.M{
			tg_operator: primitive.NewDateTimeFromTime(tg_borderTime),
		},
	}, findOptions)

	if err != nil {
		respond(c, 500, err.Error(), nil)
		return
	}

	respond(c, 200, f("Found {0} comments.", len(res)), res)
}
