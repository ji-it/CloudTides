package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"tides-server/pkg/config"
	"tides-server/pkg/models"
	"tides-server/pkg/restapi/operations/port"
)

// ListPortsHandler is API handler for /vm/ports/{id} GET
func ListPortsHandler(params port.ListPortsParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return port.NewListPortsUnauthorized()
	}

	db := config.GetDB()
	var ports []*models.Port
	db.Where("vmachine_id = ?", params.ID).Find(&ports)
	var result []*port.ListPortsOKBodyItems0

	for _, P := range ports {
		newItem := port.ListPortsOKBodyItems0{
			Port: int64(P.Port),
			URL: P.URL,
		}
		result = append(result, &newItem)
	}

	return port.NewListPortsOK().WithPayload(result)
}