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
		return resource.NewValidateResourceNotFound()
	}

	var resBody resource.ValidateResourceOKBody

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		return resource.NewValidateResourceNotFound()
	}

	defer v.Destroy(ctx)

	var dss []mo.Datacenter
	err = v.Retrieve(ctx, []string{"Datacenter"}, []string{}, &dss)
	if err != nil {
		return resource.NewValidateResourceNotFound()
	}
	for _, dc := range dss {
		fmt.Println(dc.ManagedEntity.Name)
		resBody.Results = append(resBody.Results, dc.ManagedEntity.Name)
	}

	resBody.Message = "Success"
	return resource.NewValidateResourceOK().WithPayload(&resBody)
}

func AddResourceHandler(params resource.AddResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewAddResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	body := params.ReqBody

	u, err := soap.ParseURL(body.Host)
	if err != nil {
		return resource.NewAddResourceNotFound()
	}
	u.User = url.UserPassword(body.Username, body.Password)
	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return resource.NewAddResourceNotFound()
	}

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		return resource.NewAddResourceNotFound()
	}

	defer v.Destroy(ctx)

	var dss []mo.Datacenter
	err = v.Retrieve(ctx, []string{"Datacenter"}, []string{}, &dss)
	if err != nil {
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

	v, err = m.CreateContainerView(ctx, datacenter.HostFolder, []string{"HostSystem"}, true)
	var hss []mo.HostSystem
	err = v.Retrieve(ctx, []string{"HostSystem"}, []string{}, &hss)
	db := config.GetDB()
	var res []*models.ResourceAddItem

	for _, hs := range hss {
		hostName := hs.Summary.Config.Name
		totalCPU := float64(hs.Summary.Hardware.CpuMhz) * float64(hs.Summary.Hardware.NumCpuCores) / float64(1000)
		currentCPU := float64(hs.Summary.QuickStats.OverallCpuUsage) / float64(1000)
		totalMem := float64(hs.Summary.Hardware.MemorySize) / float64(1024*1024*1024)
		currentMem := float64(hs.Summary.QuickStats.OverallMemoryUsage) / float64(1024)

		newResource := models.Resource{
			Datacenter:   body.Datacenters,
			HostAddress:  body.Host,
			HostName:     hostName,
			IsActive:     true,
			JobCompleted: 0,
			Monitored:    false,
			Password:     body.Password,
			PlatformType: body.Vmtype,
			PolicyRef:    nil,
			Status:       "unknown",
			UserRef:      uid,
			Username:     body.Username,
		}
		db.Create(&newResource)

		newHostUsage := models.HostUsage{
			CurrentCPU:  currentCPU,
			CurrentRAM:  currentMem,
			HostAddress: body.Host,
			HostName:    hostName,
			PercentCPU:  currentCPU / totalCPU,
			PercentRAM:  currentMem / totalMem,
			TotalCPU:    totalCPU,
			TotalRAM:    totalMem,
			ResourceRef: newResource.Model.ID,
		}
		db.Create(&newHostUsage)

		newResultItem := models.ResourceAddItem{
			CPUPercent:   currentCPU / totalCPU,
			RAMPercent:   currentMem / totalMem,
			CurrentCPU:   currentCPU,
			CurrentRAM:   currentMem,
			Datacenter:   body.Datacenters,
			HostName:     hostName,
			IsActive:     true,
			JobCompleted: 0,
			Monitored:    false,
			PlatformType: body.Vmtype,
			TotalCPU:     totalCPU,
			TotalRAM:     totalMem,
		}

		res = append(res, &newResultItem)

	}

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

	db.Where("user_ref = ?", uid).Find(&resources)
	if len(resources) == 0 {
		return resource.NewListResourceNotFound()
	}

	var response []*models.ResourceListItem
	for _, res := range resources {
		var resUsage models.HostUsage
		db.Where("resource_ref = ?", res.Model.ID).First(&resUsage)
		var pol models.Policy
		db.Where("id = ?", res.PolicyRef).First(&pol)

		newResultItem := models.ResourceListItem{
			CPUPercent:   resUsage.PercentCPU,
			RAMPercent:   resUsage.PercentRAM,
			CurrentCPU:   resUsage.CurrentCPU,
			CurrentRAM:   resUsage.CurrentRAM,
			DateAdded:    time.Time.String(res.Model.CreatedAt),
			Datacenter:   res.Datacenter,
			HostName:     res.HostName,
			ID:           int64(res.Model.ID),
			IsActive:     res.IsActive,
			JobCompleted: res.JobCompleted,
			Monitored:    res.Monitored,
			PlatformType: res.PlatformType,
			PolicyName:   pol.Name,
			Status:       res.Status,
			TotalCPU:     resUsage.TotalCPU,
			TotalRAM:     resUsage.TotalRAM,
			TotalJobs:    res.TotalJobs,
		}

		response = append(response, &newResultItem)
	}

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
	// db.Where("id = ? AND user_ref = ?", rid, uid).Delete(&res)
	err := db.Unscoped().Where("id = ? AND user_ref = ?", rid, uid).Delete(&res).Error
	if err != nil {
		return resource.NewDeleteResourceNotFound()
	}

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
	db.Where("id = ?", rid).First(&res)

	if res.HostName == "" {
		return resource.NewToggleActiveNotFound()
	}

	res.IsActive = !res.IsActive
	db.Save(&res)

	return resource.NewToggleActiveOK().WithPayload(&resource.ToggleActiveOKBody{
		Message: "success",
	})
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
		return resource.NewAssignPolicyNotFound()
	}

	res.PolicyRef = new(uint)
	*res.PolicyRef = uint(pid)
	db.Save(&res)

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
		return resource.NewDestroyVMNotFound()
	}

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
	db.Where("user_ref = ?", uid).Find(&resources)

	results := []*models.ResourceInfoItem{}
	for _, res := range resources {
		vms := []*models.VM{}
		db.Where("resource_ref = ?", res.Model.ID).Find(&vms)
		totalVMs := len(vms)
		var pol models.Policy
		db.Where("id = ?", res.PolicyRef).First(&pol)
		var hu models.HostUsage
		db.Where("resource_ref = ?", res.Model.ID).First(&hu)

		result := models.ResourceInfoItem{
			CPUPercent:   hu.PercentCPU,
			RAMPercent:   hu.PercentRAM,
			CurrentCPU:   hu.CurrentCPU,
			CurrentRAM:   hu.CurrentRAM,
			Datacenter:   res.Datacenter,
			DateAdded:    time.Time.String(res.Model.CreatedAt),
			HostName:     res.HostName,
			ID:           int64(res.Model.ID),
			IsActive:     res.IsActive,
			JobCompleted: res.JobCompleted,
			Monitored:    res.Monitored,
			PlatformType: res.PlatformType,
			PolicyName:   pol.Name,
			Status:       res.Status,
			TotalCPU:     hu.TotalCPU,
			TotalRAM:     hu.TotalRAM,
			TotalJobs:    res.TotalJobs,
			TotalVMs:     int64(totalVMs),
		}

		if totalVMs > 0 {
			var activeVMs int
			db.Where("resource_ref = ? AND is_destroyed = ? AND powered_on = ?",
				res.Model.ID, false, true).Find(&vms).Count(&activeVMs)
			var vm models.VM
			db.Where("resource_ref = ?", res.Model.ID).Order("created_at").First(&vm)
			result.ActiveVMs = int64(activeVMs)
			result.LastDeployed = time.Time.String(vm.Model.CreatedAt)
		}

		results = append(results, &result)
	}

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
	db.Where("user_ref = ?", uid).Find(&resources)

	results := [][]*models.ResourceVMInfoItem{}
	for _, res := range resources {
		vms := []*models.VM{}
		curvms := []*models.ResourceVMInfoItem{}
		db.Where("resource_ref = ?", res.Model.ID).Find(&vms)

		for _, vm := range vms {
			var vmu models.VMUsage
			var newvm models.ResourceVMInfoItem
			if db.Where("vm_ref = ?", vm.Model.ID).First(&vmu).RecordNotFound() {
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
