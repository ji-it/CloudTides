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
		ProjectRef:      uint(body.ProjectID),
		TemplateRef:     uint(body.TemplateID),
		ThresholdPolicy: body.Threshold,
		UserRef:         uid,
	}

	db := config.GetDB()
	err := db.Create(&newPolicy).Error

	if err != nil {
		return policy.NewAddPolicyBadRequest()
	}

	return policy.NewAddPolicyOK().WithPayload(&policy.AddPolicyOKBody{
		Message: "success",
	})
}

func UpdatePolicyHandler(params policy.UpdatePolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewAddPolicyUnauthorized()
	}

	body := params.ReqBody
	var pol models.Policy
	db := config.GetDB()
	if db.Where("id = ?", body.ID).First(&pol).RecordNotFound() {
		return policy.NewUpdatePolicyNotFound()
	}

	pol.AccountType = body.AccountType
	pol.BoincPassword = body.BoincPassword
	pol.BoincUsername = body.BoincUsername
	pol.DeployType = body.DeployType
	pol.IdlePolicy = body.Idle
	pol.IsDestroy = body.IsDestroy
	pol.Name = body.Name
	pol.ProjectRef = uint(body.ProjectID)
	pol.TemplateRef = uint(body.TemplateID)
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

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)

	policies := []*models.Policy{}
	db := config.GetDB()

	db.Where("user_ref = ?", uid).Find(&policies)
	if len(policies) == 0 {
		return policy.NewListPolicyNotFound()
	}

	results := []*policy.ResultsItems0{}
	for _, pol := range policies {
		resources := []*models.Resource{}
		var pro models.Project
		db.Where("policy_ref = ?", pol.Model.ID).Find(&resources)
		db.Where("id = ?", pol.ProjectRef).First(&pro)

		newResult := policy.ResultsItems0{
			DeployType:      pol.DeployType,
			HostsAssigned:   int64(len(resources)),
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

	body := params.ReqBody
	var pol models.Policy
	db := config.GetDB()

	err := db.Unscoped().Where("id = ?", body.ID).Delete(&pol).Error
	if err != nil {
		return policy.NewRemovePolicyNotFound()
	}

	return policy.NewRemovePolicyOK().WithPayload(&policy.RemovePolicyOKBody{
		Message: "success",
	})
}
