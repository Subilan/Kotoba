package kotoba

import (
	"github.com/gin-gonic/gin"
)

func main() {
	global := gin.New()
	global.Use(middleFormatter())

	public := global.Group("/public")
	public.POST("/auth/login", login)
	public.POST("/auth/register", register)
	public.POST("/check-token", checkToken)

	private := global.Group("/private")
	private.Use(middleTokenChecker())
	private.DELETE("/user/delete/:username", deleteAccount)

	global.Run("localhost:9080")
}