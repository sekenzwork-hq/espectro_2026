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

func (i InvestorHandlers) CreateInvestor(ctx *gin.Context) {

	nameForm, nameExists := ctx.GetPostForm("name")
	websiteUrlForm, websiteUrlExists := ctx.GetPostForm("website_url")
	phoneNumberForm, phoneNumExists := ctx.GetPostForm("phone_number")
	emailForm, emailExists := ctx.GetPostForm("email")
	logo, _ := ctx.FormFile("logo")
	if !nameExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide name"})
		return
	}
	if !phoneNumExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide phone number"})
		return
	}

	if !emailExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide email"})
		return
	}

	var websiteUrl *string

	if websiteUrlExists {
		websiteUrl = &websiteUrlForm
	}

	investor, err := i.usecases.CreateInvestor(nameForm, websiteUrl, phoneNumberForm, emailForm, logo)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Investor has been created", "investor": investor})
}

func (i InvestorHandlers) UpdateInvestor(ctx *gin.Context) {

	investorId, _ := ctx.GetPostForm("investor_id")
	nameForm, nameExists := ctx.GetPostForm("name")
	websiteUrlForm, websiteUrlExists := ctx.GetPostForm("website_url")
	phoneNumberForm, phoneNumberExists := ctx.GetPostForm("phone_number")
	emailForm, emailExists := ctx.GetPostForm("email")
	logo, _ := ctx.FormFile("logo")

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

	investorId, exists := ctx.GetQuery("investor_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide investor id"})
		return
	}

	err := i.usecases.DeleteInvestor(investorId)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Investor has been deleted"})
	}
}

func (i InvestorHandlers) RetrieveInvestors(ctx *gin.Context) {

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))
	if limitErr != nil {
		code := pkg.GetStatusCodeForError(limitErr)
		ctx.JSON(code, gin.H{"status": code, "message": limitErr.Error()})
		return
	}

	investors, retrievalErr := i.usecases.RetrieveInvestors(limit, page)
	if retrievalErr != nil {
		code := pkg.GetStatusCodeForError(retrievalErr)
		ctx.JSON(code, gin.H{"status": code, "message": retrievalErr.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "investors": investors})
	}
}
