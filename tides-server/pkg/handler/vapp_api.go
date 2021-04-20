package handler

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"math/rand"
	"tides-server/pkg/config"
	"tides-server/pkg/models"
	"tides-server/pkg/restapi/operations/vapp"
	"time"
)

func randSeqT(n int) string {
	b := make([]rune, n)
	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetVdc(resourceID uint) *govcd.Vdc{
	var res models.Resource
	var vendor models.Vendor
	var vcd models.Vcd
	db := config.GetDB()
	if db.Where("id = ?", resourceID).First(&res).RowsAffected == 0 {
		return nil
	}

	if db.Where("url = ?", res.HostAddress).First(&vendor).RowsAffected == 0 {
		return nil
	}

	if db.Where("resource_id = ?", res.ID).First(&vcd).RowsAffected == 0{
		return nil
	}

	conf := VcdConfig{
		Href: vendor.URL,
		Password: res.Password,
		User: res.Username,
		Org: vcd.Organization,
		VDC: res.Datacenter,
	}

	client, err := conf.Client() // We now have a client
	if err != nil {
		fmt.Println(err)
		return nil
	}

	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	vdc, err := org.GetVDCByName(conf.VDC, false)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return vdc
}

// New version of deploy VAPP
func DeployVAPP(org *govcd.Org, vdc *govcd.Vdc, temName string, VMs []models.VMTemp, cataName string, vAppName string, netName string, vAppID uint){
	var vappDB models.Vapp
	db := config.GetDB()

	fmt.Println(vAppID)

	if db.Preload("VMs").Where("id = ?", vAppID).First(&vappDB).RowsAffected == 0 {
		fmt.Printf("id is %d", vAppID)
	}

	catalog, _ := org.GetCatalogByName(cataName, true)
	cataItem, _ := catalog.GetCatalogItemByName(temName, true)
	vappTem, _ := cataItem.GetVAppTemplate()
	net, err := vdc.GetOrgVdcNetworkByName(netName, true)
	networks := []*types.OrgVDCNetwork{}

	networks = append(networks, net.OrgVDCNetwork)

	storageProf := vdc.Vdc.VdcStorageProfiles.VdcStorageProfile[0]

	task, err := vdc.ComposeVApp(networks, vappTem, *storageProf, vAppName, "test purpose", true)
	task.WaitTaskCompletion()

	if err != nil {
		fmt.Println(err)
	}

	vapp, err := vdc.GetVAppByName(vAppName, true)
	task, err = vapp.PowerOn()
	task.WaitTaskCompletion()

	for _, VM := range VMs {
		vm, err := vapp.GetVMByName(VM.VMName, true)

		task, err = vm.Undeploy()
		task.WaitTaskCompletion()

		task, err = vm.ChangeMemorySize(VM.VMem * 1024)
		task.WaitTaskCompletion()

		cus, _ := vm.GetGuestCustomizationSection()
		cus.Enabled = new(bool)
		*cus.Enabled = true
		cus.AdminPasswordAuto = new(bool)
		*cus.AdminPasswordAuto = false
		cus.AdminPasswordEnabled = new(bool)
		*cus.AdminPasswordEnabled = true
		passWord := randSeqT(10)
		cus.AdminPassword = passWord
		// cus.ComputerName = "tides-" + randSeq(5)
		vm.SetGuestCustomizationSection(cus)
		err = vm.PowerOnAndForceCustomization()
		if err != nil {
			fmt.Println(err)
		}
	}

	if vapp != nil{
		vappDB.IPAddress = vapp.VApp.HREF
		vappDB.IsDestroyed = false
		vappDB.Name = vapp.VApp.Name
		vappDB.PoweredOn = true
		vappDB.Status = "Running"
		db.Save(&vappDB)
		for _, VM := range vappDB.VMs {
			vm, err := vapp.GetVMByName(VM.Name, true)
			if err != nil {
				VM.Status = "Error"
				return
			}
			VM.Status = "Running"
			if len(vm.VM.NetworkConnectionSection.NetworkConnection) > 0 {
				VM.ExternalIPAddress = vm.VM.NetworkConnectionSection.NetworkConnection[0].ExternalIPAddress
				VM.IPAddress = vm.VM.NetworkConnectionSection.NetworkConnection[0].IPAddress
				VM.UserName = "root"
				VM.PassWord = vm.VM.GuestCustomizationSection.AdminPassword
			}
			db.Save(&VM)
		}
	}
}

// Deploy VAPP
func DeployVapp(org *govcd.Org, vdc *govcd.Vdc, temName string, VMName string, cataName string, vAppName string, netName string) *govcd.VApp {
	catalog, _ := org.GetCatalogByName(cataName, true)
	cataItem, _ := catalog.GetCatalogItemByName(temName, true)
	vappTem, _ := cataItem.GetVAppTemplate()
	net, err := vdc.GetOrgVdcNetworkByName(netName, true)
	networks := []*types.OrgVDCNetwork{}

	networks = append(networks, net.OrgVDCNetwork)

	storageProf := vdc.Vdc.VdcStorageProfiles.VdcStorageProfile[0]

	task, err := vdc.ComposeVApp(networks, vappTem, *storageProf, vAppName, "test purpose", true)
	task.WaitTaskCompletion()

	if err != nil {
		fmt.Println(err)
		return nil
	}

	vapp, err := vdc.GetVAppByName(vAppName, true)
	task, err = vapp.PowerOn()
	task.WaitTaskCompletion()

	vm, err := vapp.GetVMByName(VMName, true)

	task, err = vm.Undeploy()
	task.WaitTaskCompletion()

	task, err = vm.ChangeMemorySize(4096)
	task.WaitTaskCompletion()

	cus, _ := vm.GetGuestCustomizationSection()
	cus.Enabled = new(bool)
	*cus.Enabled = true
	// cus.ComputerName = "tides-" + randSeq(5)
	vm.SetGuestCustomizationSection(cus)
	err = vm.PowerOnAndForceCustomization()
	if err != nil {
		fmt.Println(err)
	}

	return vapp
}

// Power on suspended VAPP
func PowerOnVapp(vdc *govcd.Vdc, vAppName string) error {
	vapp, err := vdc.GetVAppByName(vAppName, true)
	if vapp == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return err
	}
	task, err := vapp.PowerOn()
	task.WaitTaskCompletion()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Undeploy VAPP
func UndeployVapp(vdc *govcd.Vdc, vAppName string) error {
	vapp, err := vdc.GetVAppByName(vAppName, true)
	if vapp == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return err
	}
	task, err := vapp.Undeploy()
	task.WaitTaskCompletion()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Destroy VAPP
func DestroyVAPP(vdc *govcd.Vdc, vAppName string, vApp *models.Vapp) error {
	vapp, err := vdc.GetVAppByName(vAppName, true)
	if vapp == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return nil
	}
	if vapp.VApp.Deployed {
		task, err := vapp.Undeploy()
		task.WaitTaskCompletion()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	task, err := vapp.Delete()
	task.WaitTaskCompletion()
	if err != nil {
		fmt.Println(err)
		return err
	}
	db := config.GetDB()
	vappQuery, _ := vdc.GetVAppByName(vAppName, true)
	if err == nil && vappQuery == nil {
		for _, VM := range vApp.VMs {
			db.Unscoped().Delete(&VM)
		}
		db.Unscoped().Delete(&vApp)
	}

	return nil
}

// Destroy VAPP
func DestroyVapp(vdc *govcd.Vdc, vAppName string) error {
	vapp, err := vdc.GetVAppByName(vAppName, true)
	if vapp == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return nil
	}
	if vapp.VApp.Deployed {
		task, err := vapp.Undeploy()
		task.WaitTaskCompletion()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	task, err := vapp.Delete()
	task.WaitTaskCompletion()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// AddVAPPHandler is the API handler for /vapp POST
func AddVAPPHandler(params vapp.AddVappParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vapp.NewAddVappUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)

	body := params.ReqBody

	db := config.GetDB()
	var vendor models.Vendor
	var res models.Resource
	var vcd models.Vcd
	var tem models.Template
	if db.Where("name = ?", body.Vendor).First(&vendor).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Vendor not found",
		})
	}
	if db.Where("host_address = ? AND datacenter = ?", vendor.URL, body.Datacenter).First(&res).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Datacenter not found",
		})
	}
	if db.Where("resource_id = ?", res.ID).First(&vcd).RowsAffected == 0{
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Vcd not found",
		})
	}
	if db.Preload("VMTemps").Where("ID = ?", body.Template).First(&tem).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Template not found",
		})
	}

	if res.Type != "Fixed" {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Resource is not fixed, cannot create Vapp manually",
		})
	}

	conf := VcdConfig{
		Href: vendor.URL,
		Password: res.Password,
		User: res.Username,
		Org: vcd.Organization,
		VDC: res.Datacenter,
	}
	client, err := conf.Client() // We now have a client
	if err != nil {
		fmt.Println(err)
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create client failed",
		})
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		fmt.Println(err)
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create org failed",
		})
	}
	vdc, err := org.GetVDCByName(conf.VDC, false)
	if err != nil {
		fmt.Println(err)
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create vdc failed",
		})
	}

	newVapp := models.Vapp{
		UserId: uid,
		IsDestroyed: false,
		PoweredOn:   false,
		ResourceID:  res.ID,
		Template: tem.Name,
		Status: "Creating",
		Name: body.Name,
	}

	if db.Create(&newVapp).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create vApp failed",
		})
	}

	for _, VM := range tem.VMTemps {
		newVM := models.VMachine{
			Name: VM.VMName,
			VMem: VM.VMem,
			VCPU: VM.VCPU,
			VappID: newVapp.ID,
			Disk: VM.Disk,
			UsedMoney: 0,
			IPAddress: "TBD",
			ExternalIPAddress: "TBD",
			UserName: "TBD",
			PassWord: "TBD",
			Status: "Creating",
		}
		db.Create(&newVM)
	}

	go DeployVAPP(org, vdc, tem.Name, tem.VMTemps, res.Catalog, body.Name, res.Network, newVapp.ID)
	/*Vapp := DeployVAPP(org, vdc, tem.Name, tem.VMTemps, res.Catalog, body.Name, res.Network)
	if Vapp != nil{
		newVapp := models.Vapp{
			UserId: uid,
			IPAddress:   Vapp.VApp.HREF,
			IsDestroyed: false,
			Name:        Vapp.VApp.Name,
			PoweredOn:   true,
			ResourceID:  res.ID,
			Template: tem.Name,
		}
		db.Create(&newVapp)
		for _, VM := range tem.VMTemps {
			newVM := models.VMachine{
				Name: VM.VMName,
				VMem: VM.VMem,
				VCPU: VM.VCPU,
				VappID: newVapp.ID,
				Disk: VM.Disk,
				UsedMoney: 1.1,
				IpAddress: "TBD",
				UserName: "TBD",
				PassWord: "TBD",
				Status: "creating",
			}
			db.Create(&newVM)
		}
	}*/
	return vapp.NewAddVappOK().WithPayload(&vapp.AddVappOKBody{
		ID: 1,
		Message: "Create VApp success",
	})
}

// AddVappHandler is the API handler for /vapp POST
func AddVappHandler(params vapp.AddVappParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vapp.NewAddVappUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)

	body := params.ReqBody

	db := config.GetDB()
	var vendor models.Vendor
	var res models.Resource
	var vcd models.Vcd
	var tem models.Template
	if db.Where("name = ?", body.Vendor).First(&vendor).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Vendor not found",
		})
	}
	if db.Where("host_address = ? AND datacenter = ?", vendor.URL, body.Datacenter).First(&res).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Datacenter not found",
		})
	}
	if db.Where("resource_id = ?", res.ID).First(&vcd).RowsAffected == 0{
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Vcd not found",
		})
	}
	if db.Where("ID = ?", body.Template).First(&tem).RowsAffected == 0 {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Template not found",
		})
	}

	if res.Type != "Fixed" {
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Resource is not fixed, cannot create Vapp manually",
		})
	}

	conf := VcdConfig{
		Href: vendor.URL,
		Password: res.Password,
		User: res.Username,
		Org: vcd.Organization,
		VDC: res.Datacenter,
	}
	client, err := conf.Client() // We now have a client
	if err != nil {
		fmt.Println(err)
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create client failed",
		})
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		fmt.Println(err)
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create org failed",
		})
	}
	vdc, err := org.GetVDCByName(conf.VDC, false)
	if err != nil {
		fmt.Println(err)
		return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
			Message: "Create vdc failed",
		})
	}

	Vapp := DeployVapp(org, vdc, tem.Name, tem.VMName, res.Catalog, body.Name, res.Network)
	if Vapp != nil{
		newVapp := models.Vapp{
			UserId: uid,
			IPAddress:   Vapp.VApp.HREF,
			IsDestroyed: false,
			Name:        Vapp.VApp.Name,
			PoweredOn:   true,
			ResourceID:  res.ID,
			Template: tem.Name,
		}
		db.Create(&newVapp)
	}
	return vapp.NewAddVappOK().WithPayload(&vapp.AddVappOKBody{
		ID: 1,
		Message: "Create VApp success",
	})
}

// ListVappHandler is API handler for /vapp GET
func ListVappHandler(params vapp.ListVappsParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vapp.NewListVappsUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	vapps := []*models.Vapp{}
	db := config.GetDB()

	if VerifyAdmin(params.HTTPRequest) {
		db.Find(&vapps)
	} else {
		db.Where("user_id = ?", uid).Find(&vapps)
	}

	var response []*vapp.ListVappsOKBodyItems0

	for _, vap := range vapps {
		var vendor = models.Vendor{}
		var res = models.Resource{}
		db.Where("id = ?", vap.ResourceID).First(&res)
		db.Where("url = ?", res.HostAddress).First(&vendor)
		newvapp := vapp.ListVappsOKBodyItems0{
			Datacenter: res.Datacenter,
			ID: int64(vap.ID),
			Name: vap.Name,
			Template: vap.Template,
			Vendor: vendor.Name,
			Ipaddress: vap.IPAddress,
			Status: vap.Status,
			PoweredOn: vap.PoweredOn,
		}

		response = append(response, &newvapp)
	}

	return vapp.NewListVappsOK().WithPayload(response)
}

// new version of delete vapp
func DeleteVAPPHandler(params vapp.DeleteVappParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vapp.NewDeleteVappUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	var vApp models.Vapp
	db := config.GetDB()
	if db.Preload("VMs").Where("id = ?", params.ID).First(&vApp).RowsAffected == 0 {
		return vapp.NewDeleteVappNotFound()
	}
	if VerifyAdmin(params.HTTPRequest) {
	} else {
		if(vApp.UserId != uid){
			return vapp.NewDeleteVappUnauthorized()
		}
	}

	if(vApp.Status == "Creating" || vApp.Status == "Deleting") {
		return vapp.NewDeleteVappNotFound()
	}

	if vApp.IsDestroyed {
		return vapp.NewDeleteVappNotFound()
	}

	vdc := GetVdc(vApp.ResourceID)
	if vdc == nil {
		return vapp.NewDeleteVappNotFound().WithPayload(&vapp.DeleteVappNotFoundBody{
			Message: "Vdc not found",
		})
	}

	vApp.Status = "Deleting"
	db.Save(&vApp)
	for _, VM := range vApp.VMs {
		VM.Status = "Deleting"
		db.Save(&VM)
	}

	/*for i := 0; i < 3; i++ {
		err := DestroyVapp(vdc, vApp.Name)
		vappQuery, _ := vdc.GetVAppByName(vApp.Name, true)
		if err == nil && vappQuery == nil {
			for _, VM := range vApp.VMs {
				db.Unscoped().Delete(&VM)
			}
			db.Unscoped().Delete(&vApp)
			return vapp.NewDeleteVappOK().WithPayload(&vapp.DeleteVappOKBody{
				Message: "success",
			})
		}
	}

	return vapp.NewDeleteVappForbidden()*/

	go DestroyVAPP(vdc, vApp.Name, &vApp)

	return vapp.NewDeleteVappOK().WithPayload(&vapp.DeleteVappOKBody{
		Message: "success",
	})
}

// DeleteVappHandler is API handler for /vapp/{id} DELETE
func DeleteVappHandler(params vapp.DeleteVappParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return vapp.NewDeleteVappUnauthorized()
	}

	uid, _ := ParseUserIDFromToken(params.HTTPRequest)
	var vApp models.Vapp
	db := config.GetDB()
	if db.Where("id = ?", params.ID).First(&vApp).RowsAffected == 0 {
		return vapp.NewDeleteVappNotFound()
	}
	if VerifyAdmin(params.HTTPRequest) {
	} else {
		if(vApp.UserId != uid){
			return vapp.NewDeleteVappUnauthorized()
		}
	}

	if vApp.IsDestroyed {
		return vapp.NewDeleteVappNotFound()
	}

	vdc := GetVdc(vApp.ResourceID)
	if vdc == nil {
		return vapp.NewDeleteVappNotFound().WithPayload(&vapp.DeleteVappNotFoundBody{
			Message: "Vdc not found",
		})
	}

	for i := 0; i < 3; i++ {
		err := DestroyVapp(vdc, vApp.Name)
		vappQuery, _ := vdc.GetVAppByName(vApp.Name, true)
		if err == nil && vappQuery == nil {
			db.Unscoped().Delete(&vApp)
			return vapp.NewDeleteVappOK().WithPayload(&vapp.DeleteVappOKBody{
				Message: "success",
			})
		}
	}

	return vapp.NewDeleteVappForbidden()
}
