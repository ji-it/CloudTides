package handler

import (
	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/project"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

// AddProjectHandler is API handler for /project POST
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
		return project.NewAddProjectBadRequest()
	}

	return project.NewAddProjectOK().WithPayload(&project.AddProjectOKBody{
		Message: "success",
		ID:      int64(newProject.Model.ID),
	})
}

// ListProjectHandler is API handler for /project GET
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

	return project.NewListProjectOK().WithPayload(res)
}

// UpdateProjectHandler is API handler for /project/{id} PUT
func UpdateProjectHandler(params project.UpdateProjectParams) middleware.Responder {
	body := params.ReqBody

	var pro models.Project
	db := config.GetDB()
	if db.Where("id = ?", params.ID).First(&pro).Error != nil {
		return project.NewUpdateProjectNotFound()
	}

	pro.ProjectName = body.ProjectName
	pro.HasAccountManager = body.HasAccountManager
	pro.URL = body.URL

	err := db.Save(&pro).Error
	if err != nil {
		return project.NewUpdateProjectBadRequest()
	}

	return project.NewUpdateProjectOK().WithPayload(&project.UpdateProjectOKBody{
		Message: "success",
	})
}

// DeleteProjectHandler is API handler for /project/{id} DELETE
func DeleteProjectHandler(params project.DeleteProjectParams) middleware.Responder {

	db := config.GetDB()
	var pro models.Project
	err := db.Unscoped().Where("id = ?", params.ID).Delete(&pro).Error
	if err != nil {
		return project.NewDeleteProjectNotFound()
	}

	return project.NewDeleteProjectOK().WithPayload(&project.DeleteProjectOKBody{
		Message: "success",
	})
}
