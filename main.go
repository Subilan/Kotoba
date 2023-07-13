package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	router.Use(middleFormatter())

	auth := router.Group("/auth")

	auth.POST("/auth/login", login)
	auth.POST("/auth/register", register)
	auth.POST("/auth/check-token", checkToken)

	account := router.Group("/account")
	account.Use(middleTokenChecker())
	account.POST("/account/update", updateAccount)
	account.DELETE("/account/delete", deleteAccount)

	comment := router.Group("/comment")
	comment.Use(middleTokenChecker())
	comment.DELETE("/comment/delete/:commentId", deleteComment)
	comment.POST("/comment/create", createComment)
	comment.POST("/comment/update", updateComment)
	comment.POST("/comment/toggle-reaction", toggleReaction)

	publicGet := router.Group("/public/get")
	publicGet.GET("/comment", getComments)

	router.Run("localhost:9080")
}
