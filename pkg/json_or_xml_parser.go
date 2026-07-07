package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseJsonOrXML(ctx *gin.Context, obj any) bool {

	parseErr := ctx.ShouldBind(obj)

	if parseErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid Json"})
		return false
	}

	return true
}
