package handler

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"

	"tides-server/pkg/config"
	"tides-server/pkg/logger"
	"tides-server/pkg/models"
	"tides-server/pkg/restapi/operations/resource"
)

func ValidateResourceHandler(params resource.ValidateResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewValidateResourceUnauthorized()
	}

	body := params.ReqBody

	u, err := soap.ParseURL(body.Host)
	if err != nil {
		return resource.NewValidateResourceNotFound()
	}
	u.User = url.UserPassword(body.Username, body.Password)
	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/validate/: [404] Connection failure")
		return resource.NewValidateResourceNotFound()
	}

	var resBody resource.ValidateResourceOKBody

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/validate/: [404] Container view failure")
		return resource.NewValidateResourceNotFound()
	}

	defer v.Destroy(ctx)

	var dss []mo.Datacenter
	err = v.Retrieve(ctx, []string{"Datacenter"}, []string{}, &dss)
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/validate/: [404] Datacenter not found")
		return resource.NewValidateResourceNotFound()
	}
	for _, dc := range dss {
		fmt.Println(dc.ManagedEntity.Name)
		resBody.Results = append(resBody.Results, dc.ManagedEntity.Name)
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/validate/: [200] Resource validation success")

	resBody.Message = "Success"
	return resource.NewValidateResourceOK().WithPayload(&resBody)
}

// Register new clusters or resource pools
func AddResourceHandler(params resource.AddResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewAddResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	body := params.ReqBody

	u, err := soap.ParseURL(body.HostAddress)
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/add/: [404] Connection failure")
		return resource.NewAddResourceNotFound()
	}
	u.User = url.UserPassword(body.Username, body.Password)
	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/add: [404] Connection failure")
		return resource.NewAddResourceNotFound()
	}

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/add/: [404] Container view failure")
		return resource.NewAddResourceNotFound()
	}

	defer v.Destroy(ctx)

	var dss []mo.Datacenter
	err = v.Retrieve(ctx, []string{"Datacenter"}, []string{}, &dss)
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/add/: [404] Datacenter not found")
		return resource.NewAddResourceNotFound()
	}

	var datacenter mo.Datacenter
	found := false
	for _, dc := range dss {
		if dc.ManagedEntity.Name == body.Datacenters {
			datacenter = dc
			found = true
			break
		}
	}

	if !found {
		return resource.NewAddResourceNotFound()
	}

	db := config.GetDB()

	res := []*models.ResourceAddItem{}

	if !body.IsResourcePool {
		var clu models.Resource
		if db.Where("host_address = ? AND cluster = ?", body.HostAddress, body.Cluster).First(&clu).RowsAffected > 0 {
			if clu.IsResourcePool {
				logger.SetLogLevel("ERROR")
				logger.Error("/resource/add/: [404] Child resource pool already registered")
				return resource.NewAddResourceNotFound().WithPayload(&resource.AddResourceNotFoundBody{
					Message: "Child resource pool already registered!",
				})
			} else {
				logger.SetLogLevel("ERROR")
				logger.Error("/resource/add/: [404] Cluster already registered")
				return resource.NewAddResourceNotFound().WithPayload(&resource.AddResourceNotFoundBody{
					Message: "Cluster already registered!",
				})
			}
		}

		v, err = m.CreateContainerView(ctx, datacenter.HostFolder, []string{"ClusterComputeResource"}, true)
		var hss []mo.ClusterComputeResource
		err = v.Retrieve(ctx, []string{"ClusterComputeResource"}, []string{}, &hss)
		for _, newclus := range body.Resources {
			for _, hs := range hss {
				if newclus == hs.ManagedEntity.Name {
					newres := models.Resource{
						Cluster:        newclus,
						Datacenter:     body.Datacenters,
						HostAddress:    body.HostAddress,
						IsActive:       true,
						IsResourcePool: false,
						JobCompleted:   0,
						Monitored:      false,
						Name:           newclus,
						Password:       body.Password,
						PlatformType:   body.Vmtype,
						Status:         "unknown",
						TotalJobs:      0,
						UserID:         uid,
						Username:       body.Username,
					}

					db.Create(&newres)

					Summary := hs.Summary.GetComputeResourceSummary()

					CurrentCPU := float64(Summary.TotalCpu-Summary.EffectiveCpu) / 1000.0
					CurrentRAM := float64(Summary.TotalMemory)/(1024.0*1024.0*1024.0) - float64(Summary.EffectiveMemory)/1024.0
					TotalCPU := float64(Summary.TotalCpu) / 1000.0
					TotalRAM := float64(Summary.TotalMemory) / (1024.0 * 1024.0 * 1024.0)

					newResourceUsage := models.ResourceUsage{
						CurrentCPU:  CurrentCPU,
						CurrentRAM:  CurrentRAM,
						HostAddress: body.HostAddress,
						Name:        newclus,
						PercentCPU:  CurrentCPU / TotalCPU,
						PercentRAM:  CurrentRAM / TotalRAM,
						TotalCPU:    TotalCPU,
						TotalRAM:    TotalRAM,
						ResourceID:  newres.Model.ID,
					}

					db.Create(&newResourceUsage)

					re := models.ResourceAddItem{
						CPUPercent:     CurrentCPU / TotalCPU,
						RAMPercent:     CurrentRAM / TotalRAM,
						Cluster:        body.Cluster,
						CurrentCPU:     CurrentCPU,
						CurrentRAM:     CurrentRAM,
						Datacenter:     body.Datacenters,
						ID:             int64(newres.Model.ID),
						IsActive:       true,
						IsResourcePool: false,
						JobCompleted:   0,
						Monitored:      false,
						Name:           newclus,
						PlatformType:   body.Vmtype,
						Status:         "unknown",
						TotalCPU:       TotalCPU,
						TotalJobs:      0,
						TotalRAM:       TotalRAM,
					}

					res = append(res, &re)
				}
			}
		}
	} else {
		var clu models.Resource
		if db.Where("host_address = ? AND cluster = ?", body.HostAddress, body.Cluster).First(&clu).Error != nil {
			logger.SetLogLevel("ERROR")
			logger.Error("/resource/add/: [404] Parent cluster already registered")
			return resource.NewAddResourceNotFound().WithPayload(&resource.AddResourceNotFoundBody{
				Message: "Parent cluster already registered!",
			})
		}

		v, err = m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"ResourcePool"}, true)
		var pool []mo.ResourcePool
		err = v.Retrieve(ctx, []string{"ResourcePool"}, []string{}, &pool)

		for _, newpo := range body.Resources {
			for _, po := range pool {
				if newpo == po.ManagedEntity.Name {
					newres := models.Resource{
						Cluster:        body.Cluster,
						Datacenter:     body.Datacenters,
						HostAddress:    body.HostAddress,
						IsActive:       true,
						IsResourcePool: true,
						JobCompleted:   0,
						Monitored:      false,
						Name:           newpo,
						Password:       body.Password,
						PlatformType:   body.Vmtype,
						Status:         "unknown",
						TotalJobs:      0,
						UserID:         uid,
						Username:       body.Username,
					}

					db.Create(&newres)

					Summary := po.Summary.GetResourcePoolSummary().QuickStats
					RuntimeInfo := po.Runtime

					CurrentCPU := float64(RuntimeInfo.Cpu.OverallUsage)
					TotalCPU := float64(RuntimeInfo.Cpu.MaxUsage)
					CurrentRAM := float64(RuntimeInfo.Memory.OverallUsage) / (1024.0 * 1024.0)
					TotalRAM := float64(RuntimeInfo.Memory.MaxUsage) / (1024.0 * 1024.0)

					fmt.Println(CurrentCPU, TotalCPU, CurrentRAM, TotalRAM, Summary.GuestMemoryUsage, Summary.HostMemoryUsage, Summary.DistributedCpuEntitlement,
						Summary.DistributedMemoryEntitlement, Summary.PrivateMemory, Summary.SharedMemory,
						Summary.SwappedMemory, Summary.BalloonedMemory, Summary.OverheadMemory, Summary.ConsumedOverheadMemory,
						Summary.CompressedMemory, RuntimeInfo.Cpu.OverallUsage, RuntimeInfo.Cpu.MaxUsage,
						RuntimeInfo.Memory.OverallUsage, RuntimeInfo.Memory.MaxUsage)

					newResourceUsage := models.ResourceUsage{
						CurrentCPU:  CurrentCPU,
						CurrentRAM:  CurrentRAM,
						HostAddress: body.HostAddress,
						Name:        newpo,
						PercentCPU:  CurrentCPU / TotalCPU,
						PercentRAM:  CurrentRAM / TotalRAM,
						TotalCPU:    TotalCPU,
						TotalRAM:    TotalRAM,
						ResourceID:  newres.Model.ID,
					}

					db.Create(&newResourceUsage)

					re := models.ResourceAddItem{
						CPUPercent:     CurrentCPU / TotalCPU,
						RAMPercent:     CurrentRAM / TotalRAM,
						Cluster:        body.Cluster,
						CurrentCPU:     CurrentCPU,
						CurrentRAM:     CurrentRAM,
						Datacenter:     body.Datacenters,
						ID:             int64(newres.Model.ID),
						IsActive:       true,
						IsResourcePool: true,
						JobCompleted:   0,
						Monitored:      false,
						Name:           newpo,
						PlatformType:   body.Vmtype,
						Status:         "unknown",
						TotalCPU:       TotalCPU,
						TotalJobs:      0,
						TotalRAM:       TotalRAM,
					}

					res = append(res, &re)
				}
			}
		}
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/add/: [200] Resource registration success")

	return resource.NewAddResourceOK().WithPayload(&resource.AddResourceOKBody{
		Message: "success",
		Results: res,
	})
}

func ListResourceHandler(params resource.ListResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewListResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	resources := []*models.Resource{}
	db := config.GetDB()

	db.Where("user_id = ?", uid).Find(&resources)

	var response []*models.ResourceListItem
	for _, res := range resources {
		var resUsage models.ResourceUsage
		db.Where("resource_id = ?", res.Model.ID).First(&resUsage)
		var pol models.Policy
		db.Where("id = ?", res.Policy.Model.ID).First(&pol)

		newResultItem := models.ResourceListItem{
			CPUPercent:   resUsage.PercentCPU,
			RAMPercent:   resUsage.PercentRAM,
			Cluster:      res.Cluster,
			CurrentCPU:   resUsage.CurrentCPU,
			CurrentRAM:   resUsage.CurrentRAM,
			DateAdded:    time.Time.String(res.Model.CreatedAt),
			Datacenter:   res.Datacenter,
			HostAddress:  res.HostAddress,
			ID:           int64(res.Model.ID),
			IsActive:     res.IsActive,
			JobCompleted: res.JobCompleted,
			Monitored:    res.Monitored,
			Name:         res.Name,
			PlatformType: res.PlatformType,
			PolicyName:   pol.Name,
			Status:       res.Status,
			TotalCPU:     resUsage.TotalCPU,
			TotalRAM:     resUsage.TotalRAM,
			TotalJobs:    res.TotalJobs,
		}

		response = append(response, &newResultItem)
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/list/: [200] Resource retrival success")

	return resource.NewListResourceOK().WithPayload(&resource.ListResourceOKBody{
		Message: "success",
		Results: response,
	})
}

func DeleteResourceHandler(params resource.DeleteResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewListResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	rid := params.ReqBody.ID
	var res models.Resource

	db := config.GetDB()
	// db.Where("id = ? AND user_id = ?", rid, uid).Delete(&res)
	err := db.Unscoped().Where("id = ? AND user_id = ?", rid, uid).Delete(&res).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/delete/: [404] Resource not found")
		return resource.NewDeleteResourceNotFound()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/delete/: [200] Resource deletion success")
	return resource.NewDeleteResourceOK()
}

/*
func UpdateHostHandler(params resource.UpdateHostParams) middleware.Responder {
	body := params.ReqBody
	db := config.GetDB()

	var res models.Resource
	db.Where("host_address = ? AND host_name = ?", body.HostAddress, body.HostName).First(&res)
	if res.HostAddress == "" {
		return resource.NewUpdateHostNotFound()
	}

	res.CurrentCPU = body.CurrentCPU
	res.CurrentRAM = body.CurrentRAM
	db.Save(&res)

	return resource.NewUpdateHostOK().WithPayload(&resource.UpdateHostOKBody{
		Message: "success",
	})
}
*/

func ToggleActiveHandler(params resource.ToggleActiveParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewToggleActiveUnauthorized()
	}

	rid := params.ReqBody.ID
	var res models.Resource
	db := config.GetDB()
	if db.Where("id = ?", rid).First(&res).Error != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/toggle_active/: [404] Resource not found")
		return resource.NewToggleActiveNotFound()
	}

	res.IsActive = !res.IsActive
	db.Save(&res)

	logger.SetLogLevel("INFO")
	logger.Info("/resource/toggle_active/: [200] Toggle active success")
	return resource.NewToggleActiveOK().WithPayload(&resource.ToggleActiveOKBody{
		Message: "success",
	})
}

func UpdateStatusHandler(params resource.UpdateStatusParams) middleware.Responder {
	body := params.ReqBody
	var res models.Resource
	db := config.GetDB()
	if db.Where("id = ?", body.ResourceID).First(&res).Error != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/update_status/: [404] Resource not found")
		return resource.NewUpdateStatusNotFound()
	}

	if body.Status != "" {
		res.Status = body.Status
	}
	res.Monitored = body.Monitored
	err := db.Save(&res).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/update_status/: [400] Update status failed")
		return resource.NewUpdateStatusBadRequest()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/update_status/: [200] Update status success")
	return resource.NewUpdateStatusOK()
}

func AssignPolicyHandler(params resource.AssignPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewAssignPolicyUnauthorized()
	}

	body := params.ReqBody
	rid := body.ResourceID
	pid := body.PolicyID

	var res models.Resource
	var pol models.Policy
	db := config.GetDB()
	db.Where("id = ?", rid).First(&res)
	db.Where("id = ?", pid).First(&pol)

	if res.Model.ID == 0 || pol.Model.ID == 0 {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/assign_policy/: [404] Resource not found")
		return resource.NewAssignPolicyNotFound()
	}

	res.Policy = pol
	db.Save(&res)

	logger.SetLogLevel("INFO")
	logger.Info("/resource/assign_policy/: [200] Assign policy success")
	return resource.NewAssignPolicyOK().WithPayload(&resource.AssignPolicyOKBody{
		Message: "success",
	})
}

func DestroyVMHandler(params resource.DestroyVMParams) middleware.Responder {
	ip := params.ReqBody.IPAddress
	var vm models.VM

	db := config.GetDB()
	err := db.Unscoped().Where("ip_address = ?", ip).Delete(&vm).Error
	if err != nil {
		logger.SetLogLevel("ERROR")
		logger.Error("/resource/destroy_vm/: [404] VM not found")
		return resource.NewDestroyVMNotFound()
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/destroy_vm/: [200] VM destroyed")
	return resource.NewDestroyVMOK().WithPayload(&resource.DestroyVMOKBody{
		Message: "success",
	})
}

func ResourceInfoHandler(params resource.ResourceInfoParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewResourceInfoUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	resources := []*models.Resource{}
	db := config.GetDB()
	db.Where("user_id = ?", uid).Find(&resources)

	results := []*models.ResourceInfoItem{}
	for _, res := range resources {
		vms := []*models.VM{}
		db.Where("resource_id = ?", res.Model.ID).Find(&vms)
		totalVMs := len(vms)
		var pol models.Policy
		db.Where("id = ?", res.Policy.Model.ID).First(&pol)
		var hu models.ResourceUsage
		db.Where("resource_id = ?", res.Model.ID).First(&hu)

		result := models.ResourceInfoItem{
			CPUPercent:   hu.PercentCPU,
			RAMPercent:   hu.PercentRAM,
			Cluster:      res.Cluster,
			CurrentCPU:   hu.CurrentCPU,
			CurrentRAM:   hu.CurrentRAM,
			Datacenter:   res.Datacenter,
			DateAdded:    time.Time.String(res.Model.CreatedAt),
			HostAddress:  res.HostAddress,
			ID:           int64(res.Model.ID),
			IsActive:     res.IsActive,
			JobCompleted: res.JobCompleted,
			Monitored:    res.Monitored,
			Name:         res.Name,
			PlatformType: res.PlatformType,
			PolicyName:   pol.Name,
			Status:       res.Status,
			TotalCPU:     hu.TotalCPU,
			TotalRAM:     hu.TotalRAM,
			TotalJobs:    res.TotalJobs,
			TotalVMs:     int64(totalVMs),
		}

		if totalVMs > 0 {
			var activeVMs int64
			db.Where("resource_id = ? AND is_destroyed = ? AND powered_on = ?",
				res.Model.ID, false, true).Find(&vms).Count(&activeVMs)
			var vm models.VM
			db.Where("resource_id = ?", res.Model.ID).Order("created_at").First(&vm)
			result.ActiveVMs = int64(activeVMs)
			result.LastDeployed = time.Time.String(vm.Model.CreatedAt)
		}

		results = append(results, &result)
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/get_details/: [200] Resource info retrival success")
	return resource.NewResourceInfoOK().WithPayload(&resource.ResourceInfoOKBody{
		Message: "success",
		Results: results,
	})
}

func ResourceVMsInfoHandler(params resource.ResourceVMsInfoParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewResourceVMsInfoUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	resources := []*models.Resource{}
	db := config.GetDB()
	db.Where("user_id = ?", uid).Find(&resources)

	results := [][]*models.ResourceVMInfoItem{}
	for _, res := range resources {
		vms := []*models.VM{}
		curvms := []*models.ResourceVMInfoItem{}
		db.Where("resource_id = ?", res.Model.ID).Find(&vms)

		for _, vm := range vms {
			var vmu models.VMUsage
			var newvm models.ResourceVMInfoItem
			if db.Where("vm_id = ?", vm.Model.ID).First(&vmu).Error != nil {
				newvm = models.ResourceVMInfoItem{
					BoincTime:   time.Time.String(vm.BoincTime),
					DateCreated: time.Time.String(vm.Model.CreatedAt),
					GuestOS:     vm.GuestOS,
					ID:          int64(vm.Model.ID),
					IPAddress:   vm.IPAddress,
					Name:        vm.Name,
					PoweredOn:   vm.PoweredOn,
				}
			} else {
				newvm = models.ResourceVMInfoItem{
					CPUPercent:  vmu.CurrentCPU / vmu.TotalCPU,
					RAMPercent:  vmu.CurrentRAM / vmu.TotalRAM,
					BoincTime:   time.Time.String(vm.BoincTime),
					CurrentCPU:  vmu.CurrentCPU,
					CurrentRAM:  vmu.CurrentRAM,
					DateCreated: time.Time.String(vm.Model.CreatedAt),
					GuestOS:     vm.GuestOS,
					ID:          int64(vm.Model.ID),
					IPAddress:   vm.IPAddress,
					Name:        vm.Name,
					PoweredOn:   vm.PoweredOn,
					TotalCPU:    vmu.TotalCPU,
					TotalRAM:    vmu.TotalRAM,
				}
			}

			curvms = append(curvms, &newvm)
		}

		results = append(results, curvms)
	}

	logger.SetLogLevel("INFO")
	logger.Info("/resource/get_vm_details/: [200] Resource VM info retrival success")
	return resource.NewResourceVMsInfoOK().WithPayload(&resource.ResourceVMsInfoOKBody{
		Message: "success",
		Results: results,
	})
}

/*
func OverviewStatsHandler(params resource.OverviewStatsParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewOverviewStatsUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	db := config.GetDB()
	resources := []*models.Resource{}
	db.Where("user_ref = ?", uid).Find(&resources)

	totalVMs := 0
	power := 0
	cost := 0.0
	totalActiveVM := 0
	totalDestroyedVM := 0
	costCPU := 0.0
	costRAM := 0.0
	for _, res := range resources {
		var activeVMs int
		var inactiveVMs int
		vms := []*models.VM{}
		totalVMs += int(res.TotalVMs)
		db.Where("resource_ref = ? AND is_destroyed = ? AND powered_on = ?",
			res.Model.ID, false, true).Find(&vms).Count(&activeVMs)
		totalActiveVM += activeVMs

		for _, vm := range vms {
			var vmusage models.VMUsage
			db.Where("vm_ref = ?", vm.Model.ID).First(&vmusage)
			costCPU += vmusage.CurrentCPU / vmusage.TotalCPU
			costRAM += vmusage.CurrentRAM / vmusage.TotalRAM
		}
		db.Where("resource_ref = ? AND is_destroyed = ?", res.Model.ID, true).
			Find(&vms).Count(&inactiveVMs)
		totalDestroyedVM += inactiveVMs
	}
	cost = 600*costCPU + 200*costRAM

	resourcesUsed := len(resources)
}
*/
