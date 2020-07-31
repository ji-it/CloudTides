package handler

import (
	"time"

	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/template"

	"tides-server/pkg/config"
	"tides-server/pkg/logger"
	"tides-server/pkg/models"
)

func AddTemplateHandler(params template.AddTemplateParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return template.NewAddTemplateUnauthorized()
	}

	body := params.ReqBody
	newTem := models.Template{
		Compatibility:    body.Compat,
		GuestOS:          body.Os,
		MemorySize:       body.Memsize,
		Name:             body.Name,
		ProvisionedSpace: body.Space,
		TemplateType:     body.Source,
	}

	db := config.GetDB()
	err := db.Create(&newTem).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/template/add/: [400] Template add failure")
		return template.NewAddTemplateBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/template/add/: [200] Template add success")
	return template.NewAddTemplateOK()
}

func ListTemplateHandler(params template.ListTemplateParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return template.NewListTemplateUnauthorized()
	}

	db := config.GetDB()
	var templates []*models.Template
	db.Find(&templates)
	var result []*template.ResultsItems0

	for _, tem := range templates {
		newItem := template.ResultsItems0{
			Compatibility:    tem.Compatibility,
			DateAdded:        time.Time.String(tem.Model.CreatedAt),
			GuestOS:          tem.GuestOS,
			MemorySize:       tem.MemorySize,
			Name:             tem.Name,
			ProvisionedSpace: tem.ProvisionedSpace,
			TemplateType:     tem.TemplateType,
		}
		result = append(result, &newItem)
	}

	logger.SetLogLevel("INFO")
	logger.Info("/template/list/: [200] Template retrival success")
	return template.NewListTemplateOK().WithPayload(&template.ListTemplateOKBody{
		Message: "success",
		Results: result,
	})
}

func DeleteTemplateHandler(params template.DeleteTemplateParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return template.NewDeleteTemplateUnauthorized()
	}

	body := params.ReqBody
	db := config.GetDB()

	err := db.Unscoped().Where("name = ?", body.Name).Delete(&models.Template{}).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/template/delete/: [404] Template not found")
		return template.NewDeleteTemplateNotFound()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/template/delete/: [200] Template deletion success")
	return template.NewDeleteTemplateOK().WithPayload(&template.DeleteTemplateOKBody{
		Message: "success",
	})
}
