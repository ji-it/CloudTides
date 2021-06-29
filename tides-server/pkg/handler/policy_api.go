package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/policy"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

// AddPolicyHandler is API handler for /policy POST
func AddPolicyHandler(params policy.AddPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewAddPolicyUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	body := params.ReqBody

	newPolicy := models.Policy{
		DeployType:      body.DeployType,
		IdlePolicy:      body.Idle,
		IsDestroy:       body.IsDestroy,
		Name:            body.Name,
		PlatformType:    body.PlatformType,
		TemplateID:      uint(body.TemplateID),
		ThresholdPolicy: body.Threshold,
		UserID:          uid,
	}

	db := config.GetDB()
	err := db.Create(&newPolicy).Error

	if err != nil {
		return policy.NewAddPolicyBadRequest().WithPayload(&policy.AddPolicyBadRequestBody{
			Message: err.Error(),
		})
	}

	if body.PlatformType == models.ResourcePlatformTypeVcd {
		newVcdPolicy := models.VcdPolicy{
			Catalog:  body.Catalog,
			Network:  body.Network,
			PolicyID: newPolicy.ID,
		}

		err = db.Create(&newVcdPolicy).Error

		if err != nil {
			return policy.NewAddPolicyBadRequest().WithPayload(&policy.AddPolicyBadRequestBody{
				Message: err.Error(),
			})
		}
	}

	return policy.NewAddPolicyOK().WithPayload(&policy.AddPolicyOKBody{
		Message: "success",
		ID:      int64(newPolicy.Model.ID),
	})
}

// UpdatePolicyHandler is API handler for /policy/{id} PUT
func UpdatePolicyHandler(params policy.UpdatePolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewAddPolicyUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	body := params.ReqBody
	var pol models.Policy
	db := config.GetDB()
	if db.Where("id = ? AND user_id = ?", params.ID, uid).First(&pol).RowsAffected == 0 {
		return policy.NewUpdatePolicyForbidden()
	}

	pol.DeployType = body.DeployType
	pol.IdlePolicy = body.Idle
	pol.IsDestroy = body.IsDestroy
	pol.Name = body.Name
	pol.TemplateID = uint(body.TemplateID)
	pol.ThresholdPolicy = body.Threshold

	err := db.Save(&pol).Error
	if err != nil {
		return policy.NewUpdatePolicyBadRequest()
	}

	return policy.NewUpdatePolicyOK().WithPayload(&policy.UpdatePolicyOKBody{
		Message: "success",
	})
}

// ListPolicyHandler is API handler for /policy GET
func ListPolicyHandler(params policy.ListPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewListPolicyUnauthorized()
	}

	policies := []*models.Policy{}
	db := config.GetDB()

	db.Find(&policies)

	results := []*policy.ListPolicyOKBodyResultsItems0{}
	for _, pol := range policies {
		newResult := policy.ListPolicyOKBodyResultsItems0{
			DeployType:      pol.DeployType,
			ID:              int64(pol.Model.ID),
			IdlePolicy:      pol.IdlePolicy,
			IsDestroy:       pol.IsDestroy,
			Name:            pol.Name,
			PlatformType:    pol.PlatformType,
			ThresholdPolicy: pol.ThresholdPolicy,
		}

		results = append(results, &newResult)
	}

	return policy.NewListPolicyOK().WithPayload(&policy.ListPolicyOKBody{
		Message: "success",
		Results: results,
	})
}

// RemovePolicyHandler is API handler for /policy/{id} DELETE
func RemovePolicyHandler(params policy.RemovePolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewRemovePolicyUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	var pol models.Policy
	db := config.GetDB()

	if db.Unscoped().Where("id = ? AND user_id = ?", params.ID, uid).Delete(&pol).RowsAffected == 0 {
		return policy.NewRemovePolicyForbidden()
	}

	return policy.NewRemovePolicyOK().WithPayload(&policy.RemovePolicyOKBody{
		Message: "success",
	})
}

// GetPolicyHandler is API handler for /policy/{id} GET
func GetPolicyHandler(params policy.GetPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewRemovePolicyUnauthorized()
	}

	db := config.GetDB()
	var pol models.Policy
	if db.Where("id = ?", params.ID).First(&pol).RowsAffected == 0 {
		return policy.NewGetPolicyNotFound()
	}

	resources := []*models.Resource{}
	db.Where("policy_id = ?", pol.Model.ID).Find(&resources)
	var user models.User
	db.Where("id = ?", pol.UserID).First(&user)

	response := policy.GetPolicyOKBody{
		DeployType:      pol.DeployType,
		HostsAssigned:   int64(len(resources)),
		IdlePolicy:      pol.IdlePolicy,
		IsDestroy:       pol.IsDestroy,
		Name:            pol.Name,
		PlatformType:    pol.PlatformType,
		ThresholdPolicy: pol.ThresholdPolicy,
		User:            user.Username,
	}

	return policy.NewGetPolicyOK().WithPayload(&response)
}
