package handler

import (
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations/usage"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

// AddResourceUsageHandler is API handler for /usage POST
func AddResourceUsageHandler(params usage.AddResourceUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var res models.Resource

	if db.Where("host_address = ? AND name = ?", body.HostAddress, body.Name).First(&res).Error != nil {
		return usage.NewAddResourceUsageNotFound()
	}

	PercentCPU := body.CurrentCPU / body.TotalCPU
	PercentRAM := body.CurrentRAM / body.TotalRAM
	PercentDisk := body.CurrentDisk / body.TotalDisk

	newResourceUsage := models.ResourceUsage{
		CurrentCPU:  body.CurrentCPU,
		CurrentDisk: body.CurrentDisk,
		CurrentRAM:  body.CurrentRAM,
		HostAddress: body.HostAddress,
		Name:        body.Name,
		PercentCPU:  PercentCPU,
		PercentDisk: PercentDisk,
		PercentRAM:  PercentRAM,
		TotalCPU:    body.TotalCPU,
		TotalDisk:   body.TotalDisk,
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

// GetResourceUsageHandler is API handler for /usage/{id} GET
func GetResourceUsageHandler(params usage.GetResourceUsageParams) middleware.Responder {
	db := config.GetDB()

	var us models.ResourceUsage
	if db.Where("resource_id = ?", params.ID).First(&us).RowsAffected == 0 {
		return usage.NewGetResourceUsageNotFound()
	}

	res := usage.GetResourceUsageOKBody{
		CurrentCPU:  us.CurrentCPU,
		CurrentDisk: us.CurrentDisk,
		CurrentRAM:  us.CurrentRAM,
		Name:        us.Name,
		PercentCPU:  us.PercentCPU,
		PercentDisk: us.PercentDisk,
		PercentRAM:  us.PercentRAM,
		TotalCPU:    us.TotalCPU,
		TotalDisk:   us.TotalDisk,
		TotalRAM:    us.TotalRAM,
	}

	return usage.NewGetResourceUsageOK().WithPayload(&res)
}

// UpdateResourceUsageHandler is API handler for /usage/{id} PUT
func UpdateResourceUsageHandler(params usage.UpdateResourceUsageParams) middleware.Responder {
	body := params.ReqBody
	db := config.GetDB()

	var hu models.ResourceUsage
	if db.Where("resource_id = ?", params.ID).First(&hu).Error != nil {
		return usage.NewUpdateResourceUsageNotFound()
	}

	hu.CurrentCPU = body.CurrentCPU
	hu.CurrentDisk = body.CurrentDisk
	hu.CurrentRAM = body.CurrentRAM
	hu.TotalCPU = body.TotalCPU
	hu.TotalDisk = body.TotalDisk
	hu.TotalRAM = body.TotalRAM
	hu.PercentCPU = body.CurrentCPU / body.TotalCPU
	hu.PercentDisk = body.CurrentDisk / body.TotalDisk
	hu.PercentRAM = body.CurrentRAM / body.TotalRAM

	db.Save(&hu)

	return usage.NewUpdateResourceUsageOK().WithPayload(&usage.UpdateResourceUsageOKBody{
		Message: "resource usage recorded",
	})
}

// DeleteResourceUsageHandler is API handler for /resource/{id} DELETE
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

// AddVMUsageHandler is API handler for /usage/vm POST, deprecated
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
				VMID:       newvm.Model.ID,
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

// GetPastUsageHandler is API handler for /usage/past/{id} GET
func GetPastUsageHandler(params usage.GetPastUsageParams) middleware.Responder {

	db := config.GetDB()
	pastTime := time.Now().Local().Add(-time.Hour * time.Duration(params.ReqBody.TimeLength))
	var pastUsage []*models.ResourcePastUsage

	if db.Where("resource_id = ? AND created_at > ?", params.ID, pastTime).Find(&pastUsage).RowsAffected == 0 {
		return usage.NewGetPastUsageNotFound()
	}
	var responses []*usage.GetPastUsageOKBodyItems0

	for _, us := range pastUsage {
		response := usage.GetPastUsageOKBodyItems0{
			CurrentCPU:  us.CurrentCPU,
			CurrentDisk: us.CurrentDisk,
			CurrentRAM:  us.CurrentRAM,
			PercentCPU:  us.PercentCPU,
			PercentDisk: us.PercentDisk,
			PercentRAM:  us.PercentRAM,
			Time:        strfmt.DateTime(us.CreatedAt),
		}

		responses = append(responses, &response)
	}

	return usage.NewGetPastUsageOK().WithPayload(responses)
}
