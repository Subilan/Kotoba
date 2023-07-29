package main

import (
	"github.com/gin-gonic/gin"
)

func middleFormatter() gin.HandlerFunc {
	return gin.LoggerWithFormatter(loggerFormatter)
}

func middleTokenChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		checkErr := checkJWT(c.Request.Header.Get("Token"))

		if checkErr != nil {
			respond(c, 403, "Invalid credentials", nil)
			return
		}
	}
}

func middleCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

        c.Next()
	}
}