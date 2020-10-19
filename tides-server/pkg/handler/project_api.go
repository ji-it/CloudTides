package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/logger"
	"tides-server/pkg/restapi/operations/project"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

func AddProjectHandler(params project.AddProjectParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	newProject := models.Project{
		HasAccountManager: body.HasAccountManager,
		ProjectName:       body.ProjectName,
		URL:               body.URL,
	}

	err := db.Create(&newProject).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/project/add/: [400] Project creation failure")
		return project.NewAddProjectBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/policy/add/: [200] Project creation success")
	return project.NewAddProjectOK().WithPayload(&project.AddProjectOKBody{
		Message: "success",
		ID:      int64(newProject.Model.ID),
	})
}

func ListProjectHandler(params project.ListProjectParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return project.NewAddProjectUnauthorized()
	}

	projects := []*models.Project{}
	db := config.GetDB()

	db.Find(&projects)

	res := []*project.ListProjectOKBodyItems0{}
	for _, pro := range projects {
		respro := project.ListProjectOKBodyItems0{
			ID:                int64(pro.Model.ID),
			HasAccountManager: pro.HasAccountManager,
			ProjectName:       pro.ProjectName,
			URL:               pro.URL,
		}

		res = append(res, &respro)
	}

	logger.SetLogLevel("INFO")
	logger.Info("/policy/list/: [200] Project retrieval success")
	return project.NewListProjectOK().WithPayload(res)
}

func UpdateProjectHandler(params project.UpdateProjectParams) middleware.Responder {
	body := params.ReqBody

	var pro models.Project
	db := config.GetDB()
	if db.Where("id = ?", body.ID).First(&pro).Error != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/project/update/: [404] Project not found")
		return project.NewUpdateProjectNotFound()
	}

	pro.ProjectName = body.ProjectName
	pro.HasAccountManager = body.HasAccountManager
	pro.URL = body.URL

	err := db.Save(&pro).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/project/update/: [400] Project update failure")
		return project.NewUpdateProjectBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Error("/project/add/: [200] Project update success")
	return project.NewUpdateProjectOK().WithPayload(&project.UpdateProjectOKBody{
		Message: "success",
	})
}

func DeleteProjectHandler(params project.DeleteProjectParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var pro models.Project
	err := db.Unscoped().Where("id = ?", body.ID).Delete(&pro).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/project/delete/: [404] Project deletion failure")
		return project.NewDeleteProjectNotFound()
	}

	logger.SetLogLevel("INFO")
	logger.Error("/project/delete/: [200] Project deletion success")

	return project.NewDeleteProjectOK().WithPayload(&project.DeleteProjectOKBody{
		Message: "success",
	})
}
