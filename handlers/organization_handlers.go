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
	organizationName := ctx.PostForm("organization_name")

	var websiteUrl *string
	var logo *multipart.FileHeader

	if websiteUrlExists {
		websiteUrl = &websiteUrlForm
	}
	if logoForm != nil && len(logoForm.File) != 0 && len(logoForm.File["logo"]) != 0 {
		logo = logoForm.File["logo"][0]
	}

	newOrg, err := o.usecases.CreateOrganization(entity.OrganizationCreateEntity{
		Fullname:         fullname,
		Email:            email,
		Logo:             logo,
		PhoneNumber:      phoneNumber,
		WebsiteUrl:       websiteUrl,
		Industry:         industry,
		HeadQuarters:     headquarters,
		OrganizationName: organizationName,
	})

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Organization has been created", "organization": newOrg})
}

func (o OrganizationHandlers) UpdateOrganization(ctx *gin.Context) {

	id, idExists := ctx.GetPostForm("organization_id")
	if !idExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide organization id"})
		return
	}
	emailForm, emailExists := ctx.GetPostForm("email")
	logoForm, _ := ctx.MultipartForm()
	fullnameForm, fullnameExists := ctx.GetPostForm("fullname")
	phoneNumberForm, phoneNumberExists := ctx.GetPostForm("phone_number")
	websiteUrlForm, websiteUrlExists := ctx.GetPostForm("website_url")
	industryForm, industryExists := ctx.GetPostForm("industry")
	headquartersForm, headquartersExists := ctx.GetPostForm("headquarters")
	organizationNameForm, organizationNameExists := ctx.GetPostForm("organization_name")

	var email *string
	var fullname *string
	var phoneNumber *string
	var websiteUrl *string
	var logo *multipart.FileHeader
	var headquarters *string
	var indusctry *string
	var organizationName *string

	if emailExists {
		email = &emailForm
	}
	if fullnameExists {
		fullname = &fullnameForm
	}
	if phoneNumberExists {
		phoneNumber = &phoneNumberForm
	}
	if logoForm != nil && len(logoForm.File) != 0 && len(logoForm.File["logo"]) != 0 {
		logo = logoForm.File["logo"][0]
	}
	if websiteUrlExists {
		websiteUrl = &websiteUrlForm
	}

	if industryExists {
		indusctry = &industryForm
	}
	if headquartersExists {
		headquarters = &headquartersForm
	}
	if organizationNameExists {
		organizationName = &organizationNameForm
	}

	newOrg, err := o.usecases.UpdateOrganization(
		id,
		email,
		logo,
		fullname,
		phoneNumber,
		websiteUrl,
		indusctry,
		headquarters,
		organizationName,
	)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Organization has been updated", "organization": newOrg})
}

func (o OrganizationHandlers) DeleteOrganization(ctx *gin.Context) {

	id, exists := ctx.GetQuery("organization_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide organization id"})
		return
	}

	err := o.usecases.DeleteOrganization(id)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Organization has been deleted"})
}

func (o OrganizationHandlers) RetrieveOrganizations(ctx *gin.Context) {
	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))
	if limitErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": limitErr.Error()})
		return
	}

	organizations, err := o.usecases.RetrieveOrganizations(limit, page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "organizations": organizations})
}
