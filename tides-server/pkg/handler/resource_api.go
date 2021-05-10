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
	"tides-server/pkg/controller"
	"tides-server/pkg/models"
	"tides-server/pkg/restapi/operations/resource"
)

// ValidateVsphereResourceHandler is API handler for /resource/vsphere/validate GET, deprecated
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

// AddVsphereResourceHandler is API handler for /resource/vsphere POST, deprecated
func AddVsphereResourceHandler(params resource.AddVsphereResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewAddVsphereResourceUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
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
			}
			return resource.NewAddVsphereResourceNotFound().WithPayload(&resource.AddVsphereResourceNotFoundBody{
				Message: "Cluster already registered!",
			})
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

// ListVsphereResourceHandler is API handler for /resource/vsphere GET, deprecated
func ListVsphereResourceHandler(params resource.ListVsphereResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewListVsphereResourceUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
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

// ValidateVcdResourceHandler is API handler for /resource/vcd/validate GET
func ValidateVcdResourceHandler(params resource.ValidateVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewValidateVcdResourceUnauthorized()
	}

	body := params.ReqBody
	conf := config.VcdConfig{
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

// AddVcdResourceHandler is API handler for /resource/vcd POST
func AddVcdResourceHandler(params resource.AddVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewAddVcdResourceUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	body := params.ReqBody
	db := config.GetDB()
	var res models.Resource
	if db.Where("datacenter = ? AND host_address = ?", body.Datacenter, body.Href).First(&res).RowsAffected > 0 {
		return resource.NewAddVcdResourceNotFound().WithPayload(&resource.AddVcdResourceNotFoundBody{
			Message: "VCD already registered",
		})
	}
	conf := config.VcdConfig{
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
		Activated:    false,
		Datacenter:   body.Datacenter,
		HostAddress:  body.Href,
		IsActive:     true,
		JobCompleted: 0,
		Monitored:    false,
		Name:         body.Datacenter,
		Password:     password,
		PlatformType: models.ResourcePlatformTypeVcd,
		SetupStatus:  "Initializing",
		Status:       models.ResourceStatusUnknown,
		TotalJobs:    0,
		UserID:       uid,
		Username:     user.User.Name,
		Type:         body.ResType,
		Catalog:      body.Catalog,
		Network:      body.Network,
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
	storageRef := vdc.Vdc.VdcStorageProfiles.VdcStorageProfile[0].HREF
	storage, err := govcd.GetStorageProfileByHref(client, storageRef)
	CurrentDisk := float64(storage.StorageUsedMB)
	TotalDisk := float64(storage.Limit)
	newVcdUsage := models.ResourceUsage{
		CurrentCPU:  CurrentCPU,
		CurrentDisk: CurrentDisk,
		CurrentRAM:  CurrentRAM,
		HostAddress: body.Href,
		Name:        body.Datacenter,
		PercentCPU:  CurrentCPU / TotalCPU,
		PercentDisk: CurrentDisk / TotalDisk,
		PercentRAM:  CurrentRAM / TotalRAM,
		TotalCPU:    TotalCPU,
		TotalDisk:   TotalDisk,
		TotalRAM:    TotalRAM,
		ResourceID:  newres.Model.ID,
	}
	db.Create(&newVcdUsage)

	newVcdPastUsage := models.ResourcePastUsage{
		CurrentCPU:  CurrentCPU,
		CurrentDisk: CurrentDisk,
		CurrentRAM:  CurrentRAM,
		PercentCPU:  CurrentCPU / TotalCPU,
		PercentDisk: CurrentDisk / TotalDisk,
		PercentRAM:  CurrentRAM / TotalRAM,
		TotalCPU:    TotalCPU,
		TotalDisk:   TotalDisk,
		TotalRAM:    TotalRAM,
		ResourceID:  newres.Model.ID,
	}
	db.Create(&newVcdPastUsage)

	confi := config.VcdConfig{
		User:     username,
		Password: password,
		Org:      body.Org,
		Href:     body.Href,
		VDC:      body.Datacenter,
	}

	network = body.Network
	catalog = body.Catalog
	go initValidation(&confi, body.Catalog, body.Network, &newres)

	return resource.NewAddVcdResourceOK().WithPayload(&resource.AddVcdResourceOKBody{
		Message: "success",
		Results: &resource.AddVcdResourceOKBodyResults{
			ResourceID: int64(newres.Model.ID),
			Username:   user.User.Name,
			VcdID:      int64(newvcd.Model.ID),
		},
	})
}

// ListVcdResourceHandler is API handler for /resource/vcd GET
func ListVcdResourceHandler(params resource.ListVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewListVcdResourceUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
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

		var vendor models.Vendor
		db.Where("url = ?", res.HostAddress).First(&vendor)

		newres := resource.ListVcdResourceOKBodyItems0{
			AllocationModel: vcd.AllocationModel,
			Datacenter:      res.Datacenter,
			Href:            res.HostAddress,
			IsActive:        res.IsActive,
			Monitored:       res.Monitored,
			Organization:    vcd.Organization,
			Status:          res.Status,
			ID:              int64(res.ID),
			VcdID:           int64(vcd.ID),
			Vendor:			 vendor.Name,
			ResType:         res.Type,
		}
		responses = append(responses, &newres)
	}

	return resource.NewListVcdResourceOK().WithPayload(responses)
}

// GetVcdResourceHandler is API handler for /resource/vcd/{id} GET
func GetVcdResourceHandler(params resource.GetVcdResourceParams) middleware.Responder {
	/*if !VerifyUser(params.HTTPRequest) {
		return resource.NewGetVcdResourceUnauthorized()
	}*/

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	vcdID := params.ID
	db := config.GetDB()
	var vcd models.Vcd
	if db.Where("id = ?", vcdID).First(&vcd).RowsAffected == 0 {
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

	var vendor models.Vendor
	db.Where("url = ?", res.HostAddress).First(&vendor)

	response := resource.GetVcdResourceOKBody{
		AllocationModel: vcd.AllocationModel,
		CurrentCPU:      vcdUsage.CurrentCPU,
		CurrentDisk:     vcdUsage.CurrentDisk,
		CurrentRAM:      vcdUsage.CurrentRAM,
		Datacenter:      res.Datacenter,
		Href:            res.HostAddress,
		IsActive:        res.IsActive,
		JobCompleted:    res.JobCompleted,
		Monitored:       res.Monitored,
		Organization:    vcd.Organization,
		Policy:          int64(policy),
		SetupStatus:     res.SetupStatus,
		Status:          res.Status,
		TotalCPU:        vcdUsage.TotalCPU,
		TotalDisk:       vcdUsage.TotalDisk,
		TotalJobs:       res.TotalJobs,
		TotalRAM:        vcdUsage.TotalRAM,
		TotalVMs:        res.TotalVMs,
		Vendor:          vendor.Name,
	}

	return resource.NewGetVcdResourceOK().WithPayload(&response)
}

// DeleteVcdResourceHandler is API handler for /resource/vcd/{id} DELETE
func DeleteVcdResourceHandler(params resource.DeleteVcdResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewDeleteVcdResourceUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	vcdID := params.ID
	db := config.GetDB()
	var vcd models.Vcd
	if db.Where("id = ?", vcdID).First(&vcd).RowsAffected == 0 {
		return resource.NewDeleteVcdResourceNotFound()
	}

	if db.Unscoped().Where("id = ? AND user_id = ?", vcd.ResourceID, uid).Delete(&models.Resource{}).RowsAffected == 0 {
		return resource.NewDeleteVcdResourceForbidden()
	}
	controller.RemoveJob(vcd.ResourceID)

	return resource.NewDeleteVcdResourceOK().WithPayload(&resource.DeleteVcdResourceOKBody{
		Message: "success",
	})
}

// AssignPolicyHandler is API handler for /resource/policy/{id} PUT
func AssignPolicyHandler(params resource.AssignPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewUpdateResourceUnauthorized()
	}

	rid := params.ID
	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	db := config.GetDB()
	var res models.Resource
	if db.Where("id = ? AND user_id = ?", rid, uid).First(&res).RowsAffected == 0 {
		return resource.NewUpdateResourceForbidden()
	}

	if params.ReqBody.Policy > 0 {
		if res.PolicyID != nil {
			*res.PolicyID = uint(params.ReqBody.Policy)
		} else {
			res.PolicyID = new(uint)
			*res.PolicyID = uint(params.ReqBody.Policy)
		}
	} else {
		res.PolicyID = nil
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

// ActivateResourceHandler is API handler for /resource/activate/{id} PUT
func ActivateResourceHandler(params resource.ActivateResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewActivateResourceUnauthorized()
	}
	if !VerifyAdmin(params.HTTPRequest) {
		return resource.NewActivateResourceForbidden()
	}

	var res models.Resource
	db := config.GetDB()
	if db.Where("id = ?", params.ID).First(&res).RowsAffected == 0 {
		return resource.NewActivateResourceNotFound()
	}

	res.Activated = !res.Activated
	res.SetupStatus = "Validated"
	if !res.Activated {
		res.Monitored = false
		controller.RemoveJob(res.ID)
	}
	db.Save(&res)

	if res.PlatformType == models.ResourcePlatformTypeVcd {
		var vcd models.Vcd
		db.Where("resource_id = ?", res.ID).First(&vcd)
		conf := config.VcdConfig{
			User:     res.Username,
			Password: res.Password,
			Org:      vcd.Organization,
			Href:     res.HostAddress,
			VDC:      res.Datacenter,
		}
		go initDestruction(&conf)
		go PolicySetup(res.ID, res.UserID, network, catalog)
	}

	return resource.NewActivateResourceOK().WithPayload(&resource.ActivateResourceOKBody{
		Message:   "success",
		Activated: res.Activated,
	})
}

// ContributeResourceHandler is API handler for /resource/contribute/{id} PUT
func ContributeResourceHandler(params resource.ContributeResourceParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return resource.NewContributeResourceUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	db := config.GetDB()
	var res models.Resource
	if db.Where("id = ? AND user_id = ?", params.ID, uid).First(&res).RowsAffected == 0 {
		return resource.NewContributeResourceForbidden()
	}

	res.IsActive = !res.IsActive
	if !res.IsActive {
		res.Monitored = false
		controller.RemoveJob(res.ID)
	}
	db.Save(&res)

	return resource.NewContributeResourceOK().WithPayload(&resource.ContributeResourceOKBody{
		Message:     "success",
		Contributed: res.IsActive,
	})
}
