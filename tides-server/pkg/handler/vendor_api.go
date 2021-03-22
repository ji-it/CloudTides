package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"tides-server/pkg/config"
	"tides-server/pkg/models"
	"tides-server/pkg/restapi/operations/vendor_swagger"
)

// ListVendorsHandler is API handler for /vendor GET
func ListVendorsHandler(params vendor_swagger.ListVendorParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vendor_swagger.NewListVendorUnauthorized()
	}

	vendors := []*models.Vendor{}
	db := config.GetDB()

	db.Find(&vendors)

	var responses []*vendor_swagger.ListVendorOKBodyItems0
	for _, vendor := range vendors{
		newven := vendor_swagger.ListVendorOKBodyItems0{
			ID: int64(vendor.Id),
			Name: vendor.Name,
			URL: vendor.URL,
			VendorType: vendor.Type,
			Version: vendor.Version,
		}
		responses = append(responses, &newven)
	}

	return vendor_swagger.NewListVendorOK().WithPayload(responses)
}

// AddVendorHandler is API handler for /vendor POST
func AddVendorHandler(params vendor_swagger.AddVendorParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vendor_swagger.NewListVendorUnauthorized()
	}

	vendor := models.Vendor{}
	db := config.GetDB()

	body := params.ReqBody
	if db.Where("name = ? AND url = ?", body.Name, body.URL).First(&vendor).RowsAffected > 0 {
		return vendor_swagger.NewAddVendorNotFound().WithPayload(&vendor_swagger.AddVendorNotFoundBody{
			Message: "This vendor already exist",
		})
	}

	vendors := []*models.Vendor{}

	id := int(db.Find(&vendors).RowsAffected)

	newvendor := models.Vendor{
		Name: body.Name,
		Id: id + 1,
		URL: body.URL,
		Type: body.VendorType,
		Version: body.Version,
	}

	db.Create(&newvendor)

	return vendor_swagger.NewAddVendorOK().WithPayload(&vendor_swagger.AddVendorOKBody{
		Message: "Add Vendor Success",
		ID: int64(id + 1),
	})
}

// DeleteVendorHandler is API handler for /vendors/{id} DELETE

func DeleteVendorHandler(params vendor_swagger.DeleteVendorParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vendor_swagger.NewDeleteVendorUnauthorized()
	}

	vendorId := params.ID
	db := config.GetDB()
	var vendor models.Vendor

	if db.Where("id = ?", vendorId).First(&vendor).RowsAffected == 0 {
		return vendor_swagger.NewDeleteVendorNotFound()
	}

	if db.Unscoped().Where("id = ?", vendorId).Delete(&models.Vendor{}).RowsAffected == 0 {
		return vendor_swagger.NewDeleteVendorForbidden()
	}

	return vendor_swagger.NewDeleteVendorOK().WithPayload(&vendor_swagger.DeleteVendorOKBody{
		Message: "success",
	})
}