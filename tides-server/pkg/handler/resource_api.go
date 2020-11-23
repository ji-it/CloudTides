package handler

import (
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
	"tides-server/pkg/restapi/operations/resource"
)

func ValidateVsphereResourceHandler(params resource.ValidateVsphereResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewValidateVsphereResourceUnauthorized()
	}

	body := params.ReqBody

	u, err := soap.ParseURL(body.Host)
	if err != nil {
		return resource.NewValidateVsphereResourceNotFound()
	}
	u.User = url.UserPassword(body.Username, body.Password)
	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return resource.NewValidateVsphereResourceNotFound()
	}

	var resBody resource.ValidateVsphereResourceOKBody

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		return resource.NewValidateVsphereResourceNotFound()
	}

	defer v.Destroy(ctx)

	var dss []mo.Datacenter
	err = v.Retrieve(ctx, []string{"Datacenter"}, []string{}, &dss)
	if err != nil {
		return resource.NewValidateVsphereResourceNotFound()
	}
	for _, dc := range dss {
		fmt.Println(dc.ManagedEntity.Name)
		resBody.Results = append(resBody.Results, dc.ManagedEntity.Name)
	}

	resBody.Message = "Success"
	return resource.NewValidateVsphereResourceOK().WithPayload(&resBody)
}

// Register new clusters or resource pools
func AddVsphereResourceHandler(params resource.AddVsphereResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewAddVsphereResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	body := params.ReqBody

	u, err := soap.ParseURL(body.HostAddress)
	if err != nil {
		return resource.NewAddVsphereResourceNotFound()
	}
	u.User = url.UserPassword(body.Username, body.Password)
	ctx := context.Background()
	c, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		return resource.NewAddVsphereResourceNotFound()
	}

	m := view.NewManager(c.Client)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"Datacenter"}, true)
	if err != nil {
		return resource.NewAddVsphereResourceNotFound()
	}

	defer v.Destroy(ctx)

	var dss []mo.Datacenter
	err = v.Retrieve(ctx, []string{"Datacenter"}, []string{}, &dss)
	if err != nil {
		return resource.NewAddVsphereResourceNotFound()
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
		return resource.NewAddVsphereResourceNotFound()
	}

	db := config.GetDB()

	res := []*models.ResourceAddItem{}

	if !body.IsResourcePool {
		var vs models.Vsphere

		if db.Where("cluster = ?", body.Cluster).First(&vs).RowsAffected > 0 {
			if vs.IsResourcePool {
				return resource.NewAddVsphereResourceNotFound().WithPayload(&resource.AddVsphereResourceNotFoundBody{
					Message: "Child resource pool already registered!",
				})
			} else {
				return resource.NewAddVsphereResourceNotFound().WithPayload(&resource.AddVsphereResourceNotFoundBody{
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
						Datacenter:   body.Datacenters,
						HostAddress:  body.HostAddress,
						IsActive:     true,
						JobCompleted: 0,
						Monitored:    false,
						Name:         newclus,
						Password:     body.Password,
						Status:       "unknown",
						TotalJobs:    0,
						UserID:       uid,
						Username:     body.Username,
					}

					db.Create(&newres)

					newvs := models.Vsphere{
						Cluster:        body.Cluster,
						IsResourcePool: body.IsResourcePool,
						ResourceID:     newres.Model.ID,
					}

					db.Create(&newvs)

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
			return resource.NewAddVsphereResourceNotFound().WithPayload(&resource.AddVsphereResourceNotFoundBody{
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
						Datacenter:   body.Datacenters,
						HostAddress:  body.HostAddress,
						IsActive:     true,
						JobCompleted: 0,
						Monitored:    false,
						Name:         newpo,
						Password:     body.Password,
						Status:       "unknown",
						TotalJobs:    0,
						UserID:       uid,
						Username:     body.Username,
					}

					db.Create(&newres)

					newvs := models.Vsphere{
						Cluster:        body.Cluster,
						IsResourcePool: body.IsResourcePool,
						ResourceID:     newres.Model.ID,
					}

					db.Create(&newvs)

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

	return resource.NewAddVsphereResourceOK().WithPayload(&resource.AddVsphereResourceOKBody{
		Message: "success",
		Results: res,
	})
}

func ListVsphereResourceHandler(params resource.ListVsphereResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewListVsphereResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	resources := []*models.Resource{}
	db := config.GetDB()

	if VerifyAdmin(params.HTTPRequest) {
		db.Find(&resources)
	} else {
		db.Where("user_id = ?", uid).Find(&resources)
	}

	var response []*models.ResourceListItem
	for _, res := range resources {
		var resUsage models.ResourceUsage
		db.Where("resource_id = ?", res.Model.ID).First(&resUsage)
		var pol models.Policy
		db.Where("id = ?", res.Policy.Model.ID).First(&pol)

		newResultItem := models.ResourceListItem{
			CPUPercent:   resUsage.PercentCPU,
			RAMPercent:   resUsage.PercentRAM,
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
			PolicyName:   pol.Name,
			Status:       res.Status,
			TotalCPU:     resUsage.TotalCPU,
			TotalRAM:     resUsage.TotalRAM,
			TotalJobs:    res.TotalJobs,
		}

		response = append(response, &newResultItem)
	}

	return resource.NewListVsphereResourceOK().WithPayload(&resource.ListVsphereResourceOKBody{
		Message: "success",
		Results: response,
	})
}

func ValidateVcdResourceHandler(params resource.ValidateVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewValidateVcdResourceUnauthorized()
	}

	body := params.ReqBody
	conf := VcdConfig{
		User:     body.Username,
		Password: body.Password,
		Org:      body.Org,
		Href:     body.Href,
		VDC:      body.Vdc,
	}
	client, err := conf.Client() // We now have a client
	if err != nil {
		return resource.NewValidateVcdResourceNotFound().WithPayload(&resource.ValidateVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		return resource.NewValidateVcdResourceNotFound().WithPayload(&resource.ValidateVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}
	_, err = org.GetVDCByName(conf.VDC, false)
	if err != nil {
		return resource.NewValidateVcdResourceNotFound().WithPayload(&resource.ValidateVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}

	return resource.NewValidateVcdResourceOK().WithPayload(&resource.ValidateVcdResourceOKBody{
		Message: "success",
	})
}

func AddVcdResourceHandler(params resource.AddVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewAddVcdResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	body := params.ReqBody
	db := config.GetDB()
	var res models.Resource
	if db.Where("datacenter = ? AND host_address = ?", body.Datacenter, body.Href).First(&res).RowsAffected > 0 {
		return resource.NewAddVcdResourceNotFound().WithPayload(&resource.AddVcdResourceNotFoundBody{
			Message: "VCD already registered",
		})
	}
	conf := VcdConfig{
		User:     body.Username,
		Password: body.Password,
		Org:      body.Org,
		Href:     body.Href,
		VDC:      body.Datacenter,
	}
	client, err := conf.Client() // We now have a client
	if err != nil {
		return resource.NewAddVcdResourceNotFound().WithPayload(&resource.AddVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		return resource.NewAddVcdResourceNotFound().WithPayload(&resource.AddVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}
	vdc, err := org.GetVDCByName(conf.VDC, false)
	if err != nil {
		return resource.NewAddVcdResourceNotFound().WithPayload(&resource.AddVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}
	adminOrg, err := client.GetAdminOrgByName(conf.Org)
	rand.Seed(time.Now().UnixNano())
	username := "cloudtides-" + randSeq(10)
	password := randSeq(20)
	// fmt.Println(username, password)
	var userDefinition = govcd.OrgUserConfiguration{
		Name:         username,
		Password:     password,
		RoleName:     govcd.OrgUserRoleOrganizationAdministrator,
		ProviderType: govcd.OrgUserProviderIntegrated,
		IsEnabled:    true,
	}
	user, err := adminOrg.CreateUserSimple(userDefinition)
	if err != nil {
		return resource.NewAddVcdResourceNotFound().WithPayload(&resource.AddVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}

	newres := models.Resource{
		Datacenter:   body.Datacenter,
		HostAddress:  body.Href,
		IsActive:     true,
		JobCompleted: 0,
		Monitored:    false,
		Name:         body.Datacenter,
		Password:     password,
		PlatformType: models.ResourcePlatformTypeVcd,
		Status:       models.ResourceStatusUnknown,
		TotalJobs:    0,
		UserID:       uid,
		Username:     user.User.Name,
	}
	if body.Policy > 0 {
		newres.PolicyID = new(uint)
		*newres.PolicyID = uint(body.Policy)
	} else {
		newres.PolicyID = nil
	}

	err = db.Create(&newres).Error
	if err != nil {
		return resource.NewAddVcdResourceNotFound().WithPayload(&resource.AddVcdResourceNotFoundBody{
			Message: err.Error(),
		})
	}
	newvcd := models.Vcd{
		AllocationModel: vdc.Vdc.AllocationModel,
		Organization:    body.Org,
		ResourceID:      newres.Model.ID,
	}
	db.Create(&newvcd)

	CurrentCPU := float64(vdc.Vdc.ComputeCapacity[0].CPU.Used)
	TotalCPU := float64(vdc.Vdc.ComputeCapacity[0].CPU.Limit)
	CurrentRAM := float64(vdc.Vdc.ComputeCapacity[0].Memory.Used)
	TotalRAM := float64(vdc.Vdc.ComputeCapacity[0].Memory.Limit)
	newVcdUsage := models.ResourceUsage{
		CurrentCPU:  CurrentCPU,
		CurrentRAM:  CurrentRAM,
		HostAddress: body.Href,
		Name:        body.Datacenter,
		PercentCPU:  CurrentCPU / TotalCPU,
		PercentRAM:  CurrentRAM / TotalRAM,
		TotalCPU:    TotalCPU,
		TotalRAM:    TotalRAM,
		ResourceID:  newres.Model.ID,
	}
	db.Create(&newVcdUsage)

	newVcdPastUsage := models.ResourcePastUsage{
		CurrentCPU: CurrentCPU,
		CurrentRAM: CurrentRAM,
		PercentCPU: CurrentCPU / TotalCPU,
		PercentRAM: CurrentRAM / TotalRAM,
		TotalCPU:   TotalCPU,
		TotalRAM:   TotalRAM,
		ResourceID: newres.Model.ID,
	}
	db.Create(&newVcdPastUsage)

	return resource.NewAddVcdResourceOK().WithPayload(&resource.AddVcdResourceOKBody{
		Message: "success",
		Results: &resource.AddVcdResourceOKBodyResults{
			ResourceID: int64(newres.Model.ID),
			Username:   user.User.Name,
			VcdID:      int64(newvcd.Model.ID),
		},
	})
}

func ListVcdResourceHandler(params resource.ListVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewListVcdResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	resources := []*models.Resource{}
	db := config.GetDB()

	if VerifyAdmin(params.HTTPRequest) {
		db.Find(&resources)
	} else {
		db.Where("user_id = ?", uid).Find(&resources)
	}

	var responses []*resource.ListVcdResourceOKBodyItems0
	for _, res := range resources {
		var vcd models.Vcd
		if db.Where("resource_id = ?", res.Model.ID).First(&vcd).RowsAffected == 0 {
			continue
		}

		newres := resource.ListVcdResourceOKBodyItems0{
			AllocationModel: vcd.AllocationModel,
			Datacenter:      res.Datacenter,
			Href:            res.HostAddress,
			IsActive:        res.IsActive,
			Monitored:       res.Monitored,
			Organization:    vcd.Organization,
			Status:          res.Status,
			ID:              int64(res.ID),
		}
		responses = append(responses, &newres)
	}

	return resource.NewListVcdResourceOK().WithPayload(responses)
}

func GetVcdResourceHandler(params resource.GetVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewGetVcdResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	vcdId := params.ID
	db := config.GetDB()
	var vcd models.Vcd
	if db.Where("id = ?", vcdId).First(&vcd).RowsAffected == 0 {
		return resource.NewGetVcdResourceNotFound()
	}

	var vcdUsage models.ResourceUsage
	db.Where("resource_id = ?", vcd.ResourceID).First(&vcdUsage)

	var res models.Resource
	if !VerifyAdmin(params.HTTPRequest) {
		if db.Where("id = ? AND user_id = ?", vcd.ResourceID, uid).First(&res).RowsAffected == 0 {
			return resource.NewGetVcdResourceForbidden()
		}
	} else {
		db.Where("id = ?", vcd.ResourceID).First(&res)
	}

	policy := 0
	if res.PolicyID != nil {
		policy = int(*res.PolicyID)
	}
	response := resource.GetVcdResourceOKBody{
		AllocationModel: vcd.AllocationModel,
		CurrentCPU:      vcdUsage.CurrentCPU,
		CurrentRAM:      vcdUsage.CurrentRAM,
		Datacenter:      res.Datacenter,
		Href:            res.HostAddress,
		IsActive:        res.IsActive,
		JobCompleted:    res.JobCompleted,
		Monitored:       res.Monitored,
		Organization:    vcd.Organization,
		Policy:          int64(policy),
		Status:          res.Status,
		TotalCPU:        vcdUsage.TotalCPU,
		TotalJobs:       res.TotalJobs,
		TotalRAM:        vcdUsage.TotalRAM,
		TotalVMs:        res.TotalVMs,
	}

	return resource.NewGetVcdResourceOK().WithPayload(&response)
}

func DeleteVcdResourceHandler(params resource.DeleteVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewDeleteVcdResourceUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	vcdId := params.ID
	db := config.GetDB()
	var vcd models.Vcd
	if db.Where("id = ?", vcdId).First(&vcd).RowsAffected == 0 {
		return resource.NewDeleteVcdResourceNotFound()
	}

	if db.Unscoped().Where("id = ? AND user_id = ?", vcd.ResourceID, uid).Delete(&models.Resource{}).RowsAffected == 0 {
		return resource.NewDeleteVcdResourceForbidden()
	}

	return resource.NewDeleteVcdResourceOK().WithPayload(&resource.DeleteVcdResourceOKBody{
		Message: "success",
	})
}

func UpdateResourceHandler(params resource.UpdateResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewUpdateResourceUnauthorized()
	}

	rid := params.ID
	db := config.GetDB()
	var res models.Resource
	if db.Where("id = ?", rid).First(&res).RowsAffected == 0 {
		return resource.NewUpdateResourceNotFound()
	}

	if params.ReqBody.Active == true || params.ReqBody.Active == false {
		res.IsActive = params.ReqBody.Active
		if params.ReqBody.Active == false {
			res.Monitored = false
		}
	}
	if params.ReqBody.Policy > 0 {
		if res.PolicyID != nil {
			*res.PolicyID = uint(params.ReqBody.Policy)
		} else {
			res.PolicyID = new(uint)
			*res.PolicyID = uint(params.ReqBody.Policy)
		}
	}
	err := db.Save(&res).Error
	if err != nil {
		return resource.NewUpdateResourceNotFound().WithPayload(&resource.UpdateResourceNotFoundBody{
			Message: err.Error(),
		})
	}

	return resource.NewUpdateResourceOK().WithPayload(&resource.UpdateResourceOKBody{
		Message: "success",
	})
}
