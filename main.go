package main

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
	public.GET("/comments", getComments)

	private := global.Group("/private")
	private.Use(middleTokenChecker())
	private.POST("/comment/delete", deleteComment)
	private.POST("/comment/create", createComment)
	private.POST("/comment/toggle-reaction", toggleReaction)

	global.Run("localhost:9080")
}