package kotoba

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	if !writeJSON(c, &obj) {
		respond(c, 500, "Invalid to decode request body as JSON.", nil)
		return
	}

	if hasEmptyValue(obj) {
		respond(c, 400, "Not enough argument.", nil)
		return
	}

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

	tg_username := res["username"].(string)
	tg_hash := res["hash"].(primitive.Binary).Data
	tg_avatar := res["avatar"].(string)
	tg_website := res["website"].(string)

	hashErr := bcrypt.CompareHashAndPassword(tg_hash, []byte(obj.Password))

	if hashErr != nil {
		respond(c, 403, "Invalid credentials", nil)
		return
	}

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

	if !writeJSON(c, &obj) {
		respond(c, 500, "Invalid to decode request body as JSON.", nil)
		return
	}

	if hasEmptyValue(obj) {
		respond(c, 400, "Not enough argument", nil)
		return
	}

	dupCount, countErr := mongoCount("accounts", bson.M{"username": obj.Username})

	if countErr != nil {
		respond(c, 500, countErr.Error(), nil)
		return
	}

	if dupCount != 0 {
		respond(c, 409, "Username is already in use.", nil)
		return
	}

	hashed, hashErr := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)

	if hashErr != nil {
		respond(c, 500, hashErr.Error(), nil)
		return
	}

	doc := bson.M{
		"username": obj.Username,
		"hash":     hashed,
		"avatar":   obj.Avatar,
		"website":  obj.Website,
	}

	insErr := mongoInsertOne("accounts", doc)

	if insErr != nil {
		respond(c, 500, insErr.Error(), nil)
		return
	}

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
