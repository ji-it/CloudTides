package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/policy"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

func AddPolicyHandler(params policy.AddPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewAddPolicyUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	body := params.ReqBody

	newPolicy := models.Policy{
		AccountType:     body.AccountType,
		BoincPassword:   body.BoincPassword,
		BoincUsername:   body.BoincUsername,
		DeployType:      body.DeployType,
		IdlePolicy:      body.Idle,
		IsDestroy:       body.IsDestroy,
		Name:            body.Name,
		ProjectID:       uint(body.ProjectID),
		TemplateID:      uint(body.TemplateID),
		ThresholdPolicy: body.Threshold,
		UserID:          uid,
	}

	db := config.GetDB()
	err := db.Create(&newPolicy).Error

	if err != nil {
		return policy.NewAddPolicyBadRequest()
	}

	return policy.NewAddPolicyOK().WithPayload(&policy.AddPolicyOKBody{
		Message: "success",
		ID:      int64(newPolicy.Model.ID),
	})
}

func UpdatePolicyHandler(params policy.UpdatePolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewAddPolicyUnauthorized()
	}

	body := params.ReqBody
	var pol models.Policy
	db := config.GetDB()
	if db.Where("id = ?", params.ID).First(&pol).Error != nil {
		return policy.NewUpdatePolicyNotFound()
	}

	pol.AccountType = body.AccountType
	pol.BoincPassword = body.BoincPassword
	pol.BoincUsername = body.BoincUsername
	pol.DeployType = body.DeployType
	pol.IdlePolicy = body.Idle
	pol.IsDestroy = body.IsDestroy
	pol.Name = body.Name
	pol.ProjectID = uint(body.ProjectID)
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

func ListPolicyHandler(params policy.ListPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewListPolicyUnauthorized()
	}

	policies := []*models.Policy{}
	db := config.GetDB()

	db.Find(&policies)

	results := []*policy.ListPolicyOKBodyResultsItems0{}
	for _, pol := range policies {
		var pro models.Project
		db.Where("id = ?", pol.ProjectID).First(&pro)
		newResult := policy.ListPolicyOKBodyResultsItems0{
			DeployType:      pol.DeployType,
			ID:              int64(pol.Model.ID),
			IdlePolicy:      pol.IdlePolicy,
			IsDestroy:       pol.IsDestroy,
			Name:            pol.Name,
			ProjectName:     pro.ProjectName,
			ThresholdPolicy: pol.ThresholdPolicy,
		}

		results = append(results, &newResult)
	}

	return policy.NewListPolicyOK().WithPayload(&policy.ListPolicyOKBody{
		Message: "success",
		Results: results,
	})
}

func RemovePolicyHandler(params policy.RemovePolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewRemovePolicyUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	var pol models.Policy
	db := config.GetDB()

	err := db.Unscoped().Where("id = ? AND user_id = ?", params.ID, uid).Delete(&pol).Error
	if err != nil {
		return policy.NewRemovePolicyNotFound()
	}

	return policy.NewRemovePolicyOK().WithPayload(&policy.RemovePolicyOKBody{
		Message: "success",
	})
}

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
	var pro models.Project
	db.Where("id = ?", pol.ProjectID).First(&pro)

	response := policy.GetPolicyOKBody{
		DeployType:      pol.DeployType,
		HostsAssigned:   int64(len(resources)),
		IdlePolicy:      pol.IdlePolicy,
		IsDestroy:       pol.IsDestroy,
		Name:            pol.Name,
		ProjectName:     pro.ProjectName,
		ThresholdPolicy: pol.ThresholdPolicy,
		User:            user.Username,
	}

	return policy.NewGetPolicyOK().WithPayload(&response)
}
