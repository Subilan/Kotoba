package main

import "github.com/gin-gonic/gin"

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