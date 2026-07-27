package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrganizationHandlers struct {
	usecases usecases.OrganizationUsecases
}

func NewOrganizationHandlers(usecases usecases.OrganizationUsecases) OrganizationHandlers {
	return OrganizationHandlers{usecases: usecases}
}

func (o OrganizationHandlers) CreateOrganization(ctx *gin.Context) {
	email := ctx.PostForm("email")
	logoForm, _ := ctx.MultipartForm()
	fullname := ctx.PostForm("fullname")
	phoneNumber := ctx.PostForm("phone_number")
	websiteUrlForm, websiteUrlExists := ctx.GetPostForm("website_url")
	industry := ctx.PostForm("industry")
	headquarters := ctx.PostForm("headquarters")

	var websiteUrl *string
	var logo *multipart.FileHeader

	if websiteUrlExists {
		websiteUrl = &websiteUrlForm
	}
	if logoForm != nil && len(logoForm.File) != 0 && len(logoForm.File["logo"]) != 0 {
		logo = logoForm.File["logo"][0]
	}

	newOrg, err := o.usecases.CreateOrganization(entity.OrganizationCreateEntity{
		Fullname:     fullname,
		Email:        email,
		Logo:         logo,
		PhoneNumber:  phoneNumber,
		WebsiteUrl:   websiteUrl,
		Industry:     industry,
		HeadQuarters: headquarters,
	})

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Organization has been created", "organization": newOrg})
}
