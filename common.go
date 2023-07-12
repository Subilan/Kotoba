package main

import (
	"fmt"
	"net/http"
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

func mapHasEmptyValue(obj primitive.M) bool {
	for _, v := range obj {
		if v == "" || v == nil {
			return true
		}
	}
	return false
}

func parseInt(tg string, ref *int64) {
	res, err := strconv.ParseInt(tg, 0, 64)

	if err != nil {
		*ref = 0
	} else {
		*ref = res
	}
}