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

func (i InvestorHandlers) UpdateInvestor(ctx *gin.Context) {

	investorId, idExists := ctx.GetPostForm("investor_id")
	nameForm, nameExists := ctx.GetPostForm("name")
	websiteUrlForm, websiteUrlExists := ctx.GetPostForm("website_url")
	phoneNumberForm, phoneNumberExists := ctx.GetPostForm("phone_number")
	emailForm, emailExists := ctx.GetPostForm("email")
	logo, _ := ctx.FormFile("logo")

	if !idExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide investor id"})
		return
	}

	var name *string
	var websiteUrl *string
	var phoneNumber *string
	var email *string

	if nameExists {
		name = &nameForm
	}
	if websiteUrlExists {
		websiteUrl = &websiteUrlForm
	}
	if phoneNumberExists {
		phoneNumber = &phoneNumberForm
	}
	if emailExists {
		email = &emailForm
	}

	newInvestor, err := i.usecases.UpdateInvestor(investorId, name, phoneNumber, email, websiteUrl, logo)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Investor has been updated", "investor": newInvestor})
	}
}

func (i InvestorHandlers) DeleteInvestor(ctx *gin.Context) {

	investorId := ctx.Query("investor_id")

	err := i.usecases.DeleteInvestor(investorId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Investor has been deleted"})
	}
}
