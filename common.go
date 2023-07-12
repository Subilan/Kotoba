package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wissance/stringFormatter"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func f(str string, args ...any) string {
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
		"ok":        code == 200,
		"code":      code,
		"message":   msg,
		"data":      data,
		"timestamp": time.Now().UnixMilli(),
	})
}

func bindJSON(c *gin.Context, objRef any) bool {
	body, ioErr := io.ReadAll(c.Request.Body)

	if ioErr != nil {
		return false
	}

	json.Unmarshal([]byte(body), objRef)
	return true
}

func structHasEmptyValue(obj any) bool {
	ref := reflect.ValueOf(obj)
	val := make([]interface{}, ref.NumField())
	for i := 0; i < ref.NumField(); i++ {
		if val[i] == "" || val[i] == nil {
			return true
		}
	}
	return false
}

func mapHasEmptyValue(obj primitive.M) bool {
	for _, v := range obj {
		if v == "" || v == nil {
			return true
		}
	}
	return false
}

func parseInt(tg string) (int64, error) {
	return strconv.ParseInt(tg, 0, 64)
}