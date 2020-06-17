package handler

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/usage"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

func AddHostUsageHandler(params usage.AddHostUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var res models.Resource

	if db.Where("host_address = ? AND host_name = ?", body.HostAddress, body.HostName).First(&res).RecordNotFound() {
		return usage.NewAddHostUsageNotFound()
	}

	PercentCPU := body.CurrentCPU / body.TotalCPU
	PercentRAM := body.CurrentRAM / body.TotalRAM

	newHostUsage := models.HostUsage{
		CurrentCPU:  body.CurrentCPU,
		CurrentRAM:  body.CurrentRAM,
		HostAddress: body.HostAddress,
		HostName:    body.HostName,
		PercentCPU:  PercentCPU,
		PercentRAM:  PercentRAM,
		TotalCPU:    body.TotalCPU,
		TotalRAM:    body.TotalRAM,
		ResourceRef: res.Model.ID,
	}

	err := db.Create(&newHostUsage).Error
	if err != nil {
		return usage.NewAddHostUsageBadRequest()
	}

	return usage.NewAddHostUsageOK().WithPayload(&usage.AddHostUsageOKBody{
		Message: "success",
	})
}

func UpdateHostUsageHandler(params usage.UpdateHostUsageParams) middleware.Responder {
	body := params.ReqBody
	fmt.Println("Request received")
	db := config.GetDB()

	var hu models.HostUsage
	if db.Where("host_address = ? AND host_name = ?", body.HostAddress, body.HostName).First(&hu).RecordNotFound() {
		return usage.NewUpdateHostUsageNotFound()
	}

	hu.CurrentCPU = body.CurrentCPU
	hu.CurrentRAM = body.CurrentRAM
	hu.PercentCPU = body.CurrentCPU / hu.TotalCPU
	hu.PercentRAM = body.CurrentRAM / hu.TotalRAM

	db.Save(&hu)

	return usage.NewUpdateHostUsageOK().WithPayload(&usage.UpdateHostUsageOKBody{
		Message: "host usage recorded",
	})
}

func DeleteHostUsageHandler(params usage.DeleteHostUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var hu models.HostUsage

	err := db.Unscoped().Where("host_address = ? AND host_name = ?", body.HostAddress, body.HostName).Delete(&hu).Error
	if err != nil {
		return usage.NewDeleteHostUsageBadRequest()
	}

	return usage.NewDeleteHostUsageOK().WithPayload(&usage.DeleteHostUsageOKBody{
		Message: "success",
	})
}

func AddVMUsageHandler(params usage.AddVMUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()

	for ip, val := range body {
		var vm models.VM
		if db.Where("ip_address = ?", ip).First(&vm).RecordNotFound() {
			var host models.Resource
			err := db.Where("host_name = ?", val.HostName).First(&host).Error
			if err != nil {
				return usage.NewAddVMUsageBadRequest()
			}

			boincTime, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", val.BoincStartTime)

			newvm := models.VM{
				BoincTime:   boincTime,
				GuestOS:     val.GuestOS,
				IPAddress:   ip,
				IsDestroyed: false,
				Name:        val.Name,
				NumCPU:      val.NumCPU,
				PoweredOn:   val.PoweredOn,
				ResourceRef: host.Model.ID,
			}

			db.Create(&newvm)

			newvmUsage := models.VMUsage{
				CurrentCPU: val.CurrentCPU,
				CurrentRAM: val.CurrentRAM,
				TotalCPU:   val.TotalCPU,
				TotalRAM:   val.TotalRAM,
				VmRef:      newvm.Model.ID,
			}

			db.Create(&newvmUsage)
		} else {
			var vmu models.VMUsage
			db.Where("vm_ref = ?", vm.Model.ID).First(&vmu)

			vmu.CurrentCPU = val.CurrentCPU
			vmu.CurrentRAM = val.CurrentRAM

			db.Save(&vm)
		}
	}

	return usage.NewAddVMUsageOK().WithPayload(&usage.AddVMUsageOKBody{
		Message: "success",
	})
}
