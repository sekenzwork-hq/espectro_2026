package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseJson(ctx *gin.Context, obj any) bool {

	parseErr := ctx.ShouldBindBodyWithJSON(obj)

	if parseErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid Json"})
		return false
	}

	return true
}
