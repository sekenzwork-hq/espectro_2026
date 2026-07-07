package handlers

import "github.com/gin-gonic/gin"

func RegisterUser(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"message": "connected",
	})
}
