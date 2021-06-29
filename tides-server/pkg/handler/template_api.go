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
		Name:        body.Name,
		Tag:         body.Tag,
		Description: body.Description,
	}

	/*db := config.GetDB()
	db.Create(&newTem)
	var vendor models.Vendor
	var res models.Resource
	var vcd models.Vcd
	if db.Where("name = ?", "ThinkCloud").First(&vendor).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Vendor not found",
		})
	}
	if db.Where("host_address = ? AND datacenter = ?", vendor.URL, "ElasticCloudDatacenter").First(&res).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Datacenter not found",
		})
	}
	if db.Where("resource_id = ?", res.ID).First(&vcd).RowsAffected == 0{
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Vcd not found",
		})
	}
	if res.Type != "Fixed" {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Resource is not fixed, cannot create Vapp manually",
		})
	}

	conf := config.VcdConfig{
		Href: vendor.URL,
		Password: res.Password,
		User: res.Username,
		Org: vcd.Organization,
		VDC: res.Datacenter,
	}
	client, err := conf.Client() // We now have a client
	if err != nil {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create client failed",
		})
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create org failed",
		})
	}
	catalog, err := org.GetCatalogByName(res.Catalog, true)
	if err != nil {
		fmt.Println(err)
	}
	cataItem, err := catalog.GetCatalogItemByName(temName, true)
	if err != nil {
		fmt.Println(err)
	}
	vappTem, err := cataItem.GetVAppTemplate()
	if err != nil {
		fmt.Println(err)
	}
	for _, vm := range vappTem.VAppTemplate.Children.VM {
		var VMDB models.VMTemp
		VMDB.VMName = vm.Name
		VMDB.TemplateID = newTem.ID
		VMDB.Disk = 64
		VMDB.VMem = 20
		VMDB.VCPU = 8
		db.Create(&VMDB)
	}*/
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
			Tag:              tem.Tag,
			Description:      tem.Description,
			ResourceID:       int64(tem.ResourceID),
			Vcpu:             tem.VCPUSize,
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
	_, err := CheckPorts(body.Ports)
	if (err != nil) {
		return vmtemp.NewAddVMTempBadRequest()
	}
	newVMTemp := models.VMTemp{
		VMName: body.Name,
		VCPU: int(body.Vcpu),
		VMem: int(body.Vmem),
		Disk: int(body.Disk),
		Ports: body.Ports,
		TemplateID: uint(body.TemplateID),
	}

	db := config.GetDB()
	err = db.Create(&newVMTemp).Error
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
	/*if !VerifyUser(params.HTTPRequest) {
		return vmtemp.NewListVMTempUnauthorized()
	}*/
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
			Ports: vmt.Ports,
		}

		result = append(result, &newItem)
	}

	return vmtemp.NewListVMTempOK().WithPayload(result)
}

// UpdateVMTemplateHandler is API handler for /vmtemp PUT
func UpdateVMTemplateHandler(params vmtemp.UpdateVMTempParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vmtemp.NewUpdateVMTempUnauthorized()
	}

	db := config.GetDB()
	var VMTemp models.VMTemp
	body := params.ReqBody
	if db.Where("id = ?", body.ID).First(&VMTemp).RowsAffected == 0 {
		return vmtemp.NewUpdateVMTempNotFound()
	}

	VMTemp.VCPU = int(body.Vcpu)
	VMTemp.Disk = int(body.Disk)
	VMTemp.VMem = int(body.Vmem)
	VMTemp.Ports = body.Ports

	if db.Save(&VMTemp).RowsAffected == 0 {
		return vmtemp.NewUpdateVMTempBadRequest()
	}

	return vmtemp.NewUpdateVMTempOK().WithPayload(&vmtemp.UpdateVMTempOKBody{ID: int64(VMTemp.ID), Message: "success"})
}

// DeleteVMTemplateHandler is API handler for /vmtemp/{id} DELETE
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
