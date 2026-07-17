package handlers

import (
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvestorHandlers struct {
	usecases usecases.InvestorUsecases
}

func NewInvestorHandlers(usecases usecases.InvestorUsecases) InvestorHandlers {
	return InvestorHandlers{usecases: usecases}
}

func (i InvestorHandlers) AddInvestor(ctx *gin.Context) {

	nameForm := ctx.PostForm("name")
	websiteUrlForm := ctx.PostForm("website_url")
	phoneNumberForm := ctx.PostForm("phone_number")
	emailForm := ctx.PostForm("email")
	logo, _ := ctx.FormFile("logo")

	var websiteUrl *string

	if len(websiteUrlForm) != 0 {
		websiteUrl = &websiteUrlForm
	}

	investor, err := i.usecases.AddInvestor(nameForm, websiteUrl, phoneNumberForm, emailForm, logo)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Investor has been added", "investor": investor})
}
