package kotoba

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
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

func writeJSON[refT any](c *gin.Context, objRef refT) bool {
	body, ioErr := io.ReadAll(c.Request.Body)

	if ioErr != nil {
		return false
	}

	json.Unmarshal([]byte(body), objRef)
	return true
}

func hasEmptyValue[structT any](obj structT) bool {
	ref := reflect.ValueOf(obj)
	val := make([]interface{}, ref.NumField())
	for i := 0; i < ref.NumField(); i++ {
		if val[i] == "" {
			return true
		}
	}
	return false
}
