package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/logger"
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
		logger.SetLogLevel("ERROR")
		logger.Error("/policy/add/: [400] Policy creation failure")
		return policy.NewAddPolicyBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/policy/add/: [200] Policy creation success")
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
	if db.Where("id = ?", body.ID).First(&pol).Error != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/policy/update/: [404] Policy not found")
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
		logger.SetLogLevel("ERROR")
		logger.Error("/policy/update/: [400] Policy update failure")
		return policy.NewUpdatePolicyBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/policy/update/: [200] Policy update success")
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

	db.Where("user_id = ?", uid).Find(&policies)

	results := []*policy.ResultsItems0{}
	for _, pol := range policies {
		resources := []*models.Resource{}
		var pro models.Project
		db.Where("policy_id = ?", pol.Model.ID).Find(&resources)
		db.Where("id = ?", pol.ProjectID).First(&pro)

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

	logger.SetLogLevel("INFO")
	logger.Info("/policy/list/: [200] Policy retrival success")
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
		logger.SetLogLevel("ERROR")
		logger.Error("/policy/remove/: [404] Policy not found")
		return policy.NewRemovePolicyNotFound()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/policy/remove/: [200] Policy deletion success")
	return policy.NewRemovePolicyOK().WithPayload(&policy.RemovePolicyOKBody{
		Message: "success",
	})
}
