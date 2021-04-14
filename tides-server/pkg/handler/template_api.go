package handler

import (
	"tides-server/pkg/restapi/operations/vmtemp"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/template"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

// AddTemplateHandler is API handler for /template POST
func AddTemplateHandler(params template.AddTemplateParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return template.NewAddTemplateUnauthorized()
	}

	body := params.ReqBody
	newTem := models.Template{
		Compatibility:    body.Compat,
		GuestOS:          body.Os,
		MemorySize:       body.Memsize,
		VCPUSize:         body.Vcpu,
		Name:             body.Name,
		ProvisionedSpace: body.Space,
		TemplateType:     body.Source,
		VMName:           body.VMName,
		ResourceID:       uint(body.ResourceID),
	}

	db := config.GetDB()
	err := db.Create(&newTem).Error
	if err != nil {
		return template.NewAddTemplateBadRequest()
	}

	return template.NewAddTemplateOK().WithPayload(&template.AddTemplateOKBody{
		Message: "success",
		ID:      int64(newTem.Model.ID),
	})
}

// ListTemplateHandler is API handler for /template GET
func ListTemplateHandler(params template.ListTemplateParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return template.NewListTemplateUnauthorized()
	}

	db := config.GetDB()
	var templates []*models.Template
	db.Find(&templates)
	var result []*template.ListTemplateOKBodyItems0

	for _, tem := range templates {
		newItem := template.ListTemplateOKBodyItems0{
			ID:               int64(tem.ID),
			Compatibility:    tem.Compatibility,
			DateAdded:        time.Time.String(tem.Model.CreatedAt),
			GuestOS:          tem.GuestOS,
			MemorySize:       tem.MemorySize,
			Name:             tem.Name,
			ProvisionedSpace: tem.ProvisionedSpace,
			TemplateType:     tem.TemplateType,
			ResourceID:       int64(tem.ResourceID),
		}
		result = append(result, &newItem)
	}

	return template.NewListTemplateOK().WithPayload(result)
}

// DeleteTemplateHandler is API handler for /template/{id} DELETE
func DeleteTemplateHandler(params template.DeleteTemplateParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return template.NewDeleteTemplateUnauthorized()
	}

	//uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	db := config.GetDB()

	var vmtemps []*models.VMTemp
	db.Where("template_id = ?", params.ID).Find(&vmtemps)
	for _, vmt := range vmtemps {
		if db.Unscoped().Where("id = ?", vmt.ID).Delete(&models.VMTemp{}).RowsAffected == 0 {
			return template.NewDeleteTemplateForbidden()
		}
	}

	if db.Unscoped().Where("id = ?", params.ID).Delete(&models.Template{}).RowsAffected == 0 {
		return template.NewDeleteTemplateForbidden()
	}

	return template.NewDeleteTemplateOK().WithPayload(&template.DeleteTemplateOKBody{
		Message: "success",
	})
}

// AddVMTemplateHandler is API handler for /vmtemp POST
func AddVMTemplateHandler(params vmtemp.AddVMTempParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vmtemp.NewAddVMTempUnauthorized()
	}

	body := params.ReqBody
	newVMTemp := models.VMTemp{
		VMName: body.Name,
		VCPU: int(body.Vcpu),
		VMem: int(body.Vmem),
		Disk: int(body.Disk),
		TemplateID: uint(body.TemplateID),
	}

	db := config.GetDB()
	err := db.Create(&newVMTemp).Error
	if err != nil {
		return vmtemp.NewAddVMTempBadRequest()
	}

	return vmtemp.NewAddVMTempOK().WithPayload(&vmtemp.AddVMTempOKBody{
		Message: "success",
		ID:      int64(newVMTemp.Model.ID),
	})
}

// ListVMTemplateHandler is API handler for /template/vmtemp/{id} GET
func ListVMTemplateHandler(params vmtemp.ListVMTempParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vmtemp.NewListVMTempUnauthorized()
	}
	templateID := params.ID

	db := config.GetDB()
	var vmtemps []*models.VMTemp
	db.Where("template_id = ?", templateID).Find(&vmtemps)
	var result []*vmtemp.ListVMTempOKBodyItems0

	for _, vmt := range vmtemps {
		newItem := vmtemp.ListVMTempOKBodyItems0{
			ID: int64(vmt.ID),
			Name: vmt.VMName,
			Vmem: int64(vmt.VMem),
			Vcpu: int64(vmt.VCPU),
			Disk: int64(vmt.Disk),
		}

		result = append(result, &newItem)
	}

	return vmtemp.NewListVMTempOK().WithPayload(result)
}

//
func DeleteVMTemplateHandler(params vmtemp.DeleteVMTempParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vmtemp.NewDeleteVMTempUnauthorized()
	}

	//uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	db := config.GetDB()

	if db.Unscoped().Where("id = ?", params.ID).Delete(&models.VMTemp{}).RowsAffected == 0 {
		return vmtemp.NewDeleteVMTempForbidden()
	}

	return vmtemp.NewDeleteVMTempOK().WithPayload(&vmtemp.DeleteVMTempOKBody{
		Message: "success",
	})
}
