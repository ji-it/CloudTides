package handler

import (
	"time"

	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/usage"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

func AddResourceUsageHandler(params usage.AddResourceUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var res models.Resource

	if db.Where("host_address = ? AND name = ?", body.HostAddress, body.Name).First(&res).Error != nil {
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
		return usage.NewAddResourceUsageBadRequest()
	}

	return usage.NewAddResourceUsageOK().WithPayload(&usage.AddResourceUsageOKBody{
		Message: "success",
	})
}

func GetResourceUsageHandler(params usage.GetResourceUsageParams) middleware.Responder {
	db := config.GetDB()

	var us models.ResourceUsage
	if db.Where("resource_id = ?", params.ID).First(&us).RowsAffected == 0 {
		return usage.NewGetResourceUsageNotFound()
	}

	res := usage.GetResourceUsageOKBody{
		CurrentCPU: us.CurrentCPU,
		CurrentRAM: us.CurrentRAM,
		Name:       us.Name,
		PercentCPU: us.PercentCPU,
		PercentRAM: us.PercentRAM,
		TotalCPU:   us.TotalCPU,
		TotalRAM:   us.TotalRAM,
	}

	return usage.NewGetResourceUsageOK().WithPayload(&res)
}

func UpdateResourceUsageHandler(params usage.UpdateResourceUsageParams) middleware.Responder {
	body := params.ReqBody
	db := config.GetDB()

	var hu models.ResourceUsage
	if db.Where("resource_id = ?", params.ID).First(&hu).Error != nil {
		return usage.NewUpdateResourceUsageNotFound()
	}

	hu.CurrentCPU = body.CurrentCPU
	hu.CurrentRAM = body.CurrentRAM
	hu.TotalCPU = body.TotalCPU
	hu.TotalRAM = body.TotalRAM
	hu.PercentCPU = body.CurrentCPU / body.TotalCPU
	hu.PercentRAM = body.CurrentRAM / body.TotalRAM

	db.Save(&hu)

	return usage.NewUpdateResourceUsageOK().WithPayload(&usage.UpdateResourceUsageOKBody{
		Message: "resource usage recorded",
	})
}

func DeleteResourceUsageHandler(params usage.DeleteResourceUsageParams) middleware.Responder {

	db := config.GetDB()
	var hu models.ResourceUsage

	err := db.Unscoped().Where("resource_id", params.ID).Delete(&hu).Error
	if err != nil {
		return usage.NewDeleteResourceUsageBadRequest()
	}

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

	return usage.NewAddVMUsageOK().WithPayload(&usage.AddVMUsageOKBody{
		Message: "success",
	})
}
