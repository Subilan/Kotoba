package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wissance/stringFormatter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func t(str string, args ...any) string {
	return stringFormatter.Format(str, args)
}

func loggerFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.DateTime),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}

func respond(c *gin.Context, code int, msg string, data any) {
	c.JSON(http.StatusOK, primitive.M{
		"ok": code == 200,
		"code": code,
		"message": msg,
		"data": data,
		"timestamp": time.Now().UnixMilli(),
	})
}

func writeJSON[refT any](c *gin.Context, objRef *refT) bool {
	body, ioErr := io.ReadAll(c.Request.Body)

	if ioErr != nil {
		return false
	}

	json.Unmarshal([]byte(body), *objRef);
	return true
}

func hasEmptyValue[structT any](obj structT) bool {
	ref := reflect.ValueOf(obj)
	val := make([]interface{}, ref.NumField())
	for i := 0; i < ref.NumField(); i++ {
		if (val[i] == "") {
			return true
		}
	}
	return false
}

func serve() {
	router := gin.New()
	router.Use(gin.LoggerWithFormatter(loggerFormatter))
	router.POST("/auth/login", login)
	router.POST("/auth/register", register)
	router.Run("localhost:9080")
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

	res, dbErr := get_one("accounts", bson.M{
		"username": obj.Username,
	})

	if dbErr != nil {
		if (dbErr == mongo.ErrNoDocuments) {
			respond(c, 404, t("No such user named {0}.", obj.Username), nil)
			return
		}
		respond(c, 500, dbErr.Error(), nil)
		return
	}

	tg_username := res["username"]
	tg_hash := res["hash"]
	tg_avatar := res["avatar"]
	tg_website := res["website"]

	hashErr := bcrypt.CompareHashAndPassword(tg_hash.([]byte), []byte(obj.Password))

	if hashErr != nil {
		respond(c, 403, "Incorrect credentials", nil)
	} else {
		respond(c, 200, "", resLogin{Username: tg_username.(string), Avatar: tg_avatar.(string), Website: tg_website.(string)})
	}
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

	dupCount, countErr := count("accounts", bson.M{"username": obj.Username})

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
		"hash": hashed,
		"avatar":   obj.Avatar,
		"website":  obj.Website,
	}

	
	err := insert_one("accounts", doc)

	if err != nil {
		respond(c, 500, err.Error(), nil)
	} else {
		respond(c, 200, t("Successfully created account {0}", obj.Username), nil)
	}
}
