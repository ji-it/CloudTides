package handler

import (
	"time"

	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/logger"
	"tides-server/pkg/restapi/operations/usage"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

func AddResourceUsageHandler(params usage.AddResourceUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var res models.Resource

	if db.Where("host_address = ? AND name = ?", body.HostAddress, body.Name).First(&res).Error != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/usage/add_resource/: [404] Resource not found")
		return usage.NewAddResourceUsageNotFound()
	}

	PercentCPU := body.CurrentCPU / body.TotalCPU
	PercentRAM := body.CurrentRAM / body.TotalRAM

	newResourceUsage := models.ResourceUsage{
		CurrentCPU:  body.CurrentCPU,
		CurrentRAM:  body.CurrentRAM,
		HostAddress: body.HostAddress,
		Name:        body.Name,
		PercentCPU:  PercentCPU,
		PercentRAM:  PercentRAM,
		TotalCPU:    body.TotalCPU,
		TotalRAM:    body.TotalRAM,
		ResourceID:  res.Model.ID,
	}

	err := db.Create(&newResourceUsage).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/usage/add_resource/: [400] Resource usage add failure")
		return usage.NewAddResourceUsageBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/usage/add_resource/: [200] Resource usage add success")
	return usage.NewAddResourceUsageOK().WithPayload(&usage.AddResourceUsageOKBody{
		Message: "success",
	})
}

func UpdateResourceUsageHandler(params usage.UpdateResourceUsageParams) middleware.Responder {
	body := params.ReqBody
	db := config.GetDB()

	var hu models.ResourceUsage
	if db.Where("host_address = ? AND name = ?", body.HostAddress, body.Name).First(&hu).Error != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/usage/update_resource/: [404] Resource usage not found")
		return usage.NewUpdateResourceUsageNotFound()
	}

	hu.CurrentCPU = body.CurrentCPU
	hu.CurrentRAM = body.CurrentRAM
	hu.TotalCPU = body.TotalCPU
	hu.TotalRAM = body.TotalRAM
	hu.PercentCPU = body.CurrentCPU / body.TotalCPU
	hu.PercentRAM = body.CurrentRAM / body.TotalRAM

	db.Save(&hu)

	logger.SetLogLevel("INFO")
	logger.Info("/usage/update_resource/: [200] Resource usage update success")
	return usage.NewUpdateResourceUsageOK().WithPayload(&usage.UpdateResourceUsageOKBody{
		Message: "resource usage recorded",
	})
}

func DeleteResourceUsageHandler(params usage.DeleteResourceUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var hu models.ResourceUsage

	err := db.Unscoped().Where("host_address = ? AND datacenter = ?", body.HostAddress, body.Datacenter).Delete(&hu).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/usage/delete_resource/: [400] Resource usage deletion failure")
		return usage.NewDeleteResourceUsageBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/usage/delete_resource/: [200] Resource usage deletion success")
	return usage.NewDeleteResourceUsageOK().WithPayload(&usage.DeleteResourceUsageOKBody{
		Message: "success",
	})
}

func AddVMUsageHandler(params usage.AddVMUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()

	for ip, val := range body.VMs {
		var vm models.VM
		if db.Where("ip_address = ?", ip).First(&vm).Error != nil {
			var res models.Resource
			err := db.Where("name = ?", body.Name).First(&res).Error
			if err != nil {
				logger.SetLogLevel("ERROR")
				logger.Error("/usage/addVM/: [400] Resource not found")
				return usage.NewAddVMUsageBadRequest()
			}

			boincTime, _ := time.Parse("2020-06-27 19:51:26", val.BoincStartTime)

			newvm := models.VM{
				BoincTime:   boincTime,
				GuestOS:     val.GuestOS,
				IPAddress:   ip,
				IsDestroyed: false,
				Name:        val.Name,
				NumCPU:      val.NumCPU,
				PoweredOn:   val.PoweredOn,
				ResourceID:  res.Model.ID,
			}

			db.Create(&newvm)

			newvmUsage := models.VMUsage{
				CurrentCPU: val.CurrentCPU,
				CurrentRAM: val.CurrentRAM,
				TotalCPU:   val.TotalCPU,
				TotalRAM:   val.TotalRAM,
				VmID:       newvm.Model.ID,
			}

			db.Create(&newvmUsage)
		} else {
			var vmu models.VMUsage
			db.Where("vm_id = ?", vm.Model.ID).First(&vmu)

			vmu.CurrentCPU = val.CurrentCPU
			vmu.CurrentRAM = val.CurrentRAM

			db.Save(&vmu)
		}
	}

	logger.SetLogLevel("INFO")
	logger.Info("/usage/addVM/: [200] VM usage add success")
	return usage.NewAddVMUsageOK().WithPayload(&usage.AddVMUsageOKBody{
		Message: "success",
	})
}
