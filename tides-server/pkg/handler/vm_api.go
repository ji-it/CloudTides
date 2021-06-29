package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"tides-server/pkg/config"
	"tides-server/pkg/models"
	"tides-server/pkg/restapi/operations/vm"
)

// ListVMHandler is API handler for /vapp/vm/{id} GET
func ListVMHandler(params vm.ListVMParams) middleware.Responder {
	/*if !VerifyUser(params.HTTPRequest) {
		return vapp.NewDeleteVappUnauthorized()
	}*/

	vappID := params.ID
	db := config.GetDB()
	var VMs []*models.VMachine
	db.Where("vapp_id = ?", vappID).Find(&VMs)
	var result []*vm.ListVMOKBodyItems0

	for _, VM := range VMs {
		newItem := vm.ListVMOKBodyItems0{
			Disk: int64(VM.Disk),
			Vmem: int64(VM.VMem),
			Vcpu: int64(VM.VCPU),
			ID: int64(VM.ID),
			IPAddress: VM.IPAddress,
			ExternalIPAddress: VM.ExternalIPAddress,
			Name: VM.Name,
			Username: VM.UserName,
			Password: VM.PassWord,
			Status: VM.Status,
			UsedMoney: VM.UsedMoney,
		}
		result = append(result, &newItem)
	}

	return vm.NewListVMOK().WithPayload(result)
}