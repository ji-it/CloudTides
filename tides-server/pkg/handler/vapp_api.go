package handler

import (
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"tides-server/pkg/config"
	"tides-server/pkg/controller"
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

	conf := config.VcdConfig{
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
func DeployVAPP(client *govcd.VCDClient, org *govcd.Org, vdc *govcd.Vdc, temName string, VMs []models.VMTemp,
	cataName string, vAppName string, netName string, vAppID uint) (err error){
	defer func() {
		if err != nil {
			var vappDB models.Vapp
			db := config.GetDB()

			if db.Preload("VMs").Where("id = ?", vAppID).First(&vappDB).RowsAffected == 0 {
				fmt.Printf("id is %d", vAppID)
			}

			DestroyVAPP(vdc, vAppName, &vappDB)
		}
	}()
	var vappDB models.Vapp
	db := config.GetDB()

	if db.Preload("VMs").Where("id = ?", vAppID).First(&vappDB).RowsAffected == 0 {
		fmt.Printf("id is %d", vAppID)
	}

	catalog, err := org.GetCatalogByName(cataName, true)
	if err != nil {
		fmt.Println(err)
	}
	cataItem, err := catalog.GetCatalogItemByName(temName, true)
	if err != nil {
		fmt.Println(err)
	}
	vappTem, err := cataItem.GetVAppTemplate()
	if err != nil {
		fmt.Println(err)
	}
	net, err := vdc.GetOrgVdcNetworkByName(netName, true)
	if err != nil {
		fmt.Println(err)
	}
	networks := []*types.OrgVDCNetwork{}

	networks = append(networks, net.OrgVDCNetwork)

	storageProf := vdc.Vdc.VdcStorageProfiles.VdcStorageProfile[1]

	task, err := vdc.ComposeVApp(networks, vappTem, *storageProf, vAppName, "", true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = task.WaitTaskCompletion()
	if err != nil {
		fmt.Println(err)
		return err
	}

	vApp, err := vdc.GetVAppByName(vAppName, true)
	if err != nil {
		return err
	}
	task, err = vApp.PowerOn()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = task.WaitTaskCompletion()
	if err != nil {
		fmt.Println(err)
		return err
	}
	/* Since ComposeVApp can only create one VM in the VApp template, so if there exist more than one VM in the template, we need to
	   Add these VMs one by one.
	*/
	var NetSec *types.NetworkConnectionSection
	for index, VM := range vappTem.VAppTemplate.Children.VM {
		if index == 0 {
			vm, err := vApp.GetVMByName(VM.Name, true)
			if err != nil {
				fmt.Println(err)
				return err
			}
			NetSec, err = vm.GetNetworkConnectionSection()
			if err != nil {
				fmt.Println(err)
				return err
			}
			continue
		}
		VMTemp := *govcd.NewVAppTemplate(&client.Client)
		VMTemp.VAppTemplate = VM
		task, err := vApp.AddNewVMWithStorageProfile(VM.Name, VMTemp, NetSec, storageProf, true)
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = task.WaitTaskCompletion()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	task, err = vApp.PowerOn()
	if err != nil {
		return err
	}
	err = task.WaitTaskCompletion()
	if err != nil {
		return err
	}

	for _, VM := range VMs {
		if VM.VMName == "Deploy" {
			continue
		} else {
			err := CusVM(vApp, &VM, "")
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	for _, VM := range VMs {
		if VM.VMName == "Deploy" {
			script := fmt.Sprintf("cd /root && ./client/client %s %s %s '%s' %s %d",
				os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USER"),
				os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), vappDB.ID)
			err := CusVM(vApp, &VM, script)
			if err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	if vApp != nil{
		vappDB.IPAddress = vApp.VApp.HREF
		vappDB.IsDestroyed = false
		vappDB.Name = vApp.VApp.Name
		vappDB.PoweredOn = true
		vappDB.Status = "Running"
		for _, VM := range vappDB.VMs {
			vm, err := vApp.GetVMByName(VM.Name, true)
			if err != nil {
				return err
			}
			vmMonitor, _ := controller.VMMonitors.Load(VM.ID)
			vmMonitor.Task.Start()
			if len(vm.VM.NetworkConnectionSection.NetworkConnection) > 0 {
				VM.ExternalIPAddress = vm.VM.NetworkConnectionSection.NetworkConnection[0].ExternalIPAddress
				VM.IPAddress = vm.VM.NetworkConnectionSection.NetworkConnection[0].IPAddress
				VM.UserName = "root"
				VM.PassWord = vm.VM.GuestCustomizationSection.AdminPassword
			}
			db.Save(&VM)
		}
		db.Save(&vappDB)
	}

	for _, VM := range vappDB.VMs {
		err := ExposePorts(vdc, int(VM.ID), vappDB.Name)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	return err
}

// Customize the VM
func CusVM (vApp *govcd.VApp, VM *models.VMTemp, script string) error {
	vm, err := vApp.GetVMByName(VM.VMName, true)
	if err != nil {
		return err
	}

	task, err := vm.Undeploy()
	if err != nil {
		return err
	}
	task.WaitTaskCompletion()

	task, err = vm.ChangeMemorySize(VM.VMem * 1024)
	if err != nil {
		return err
	}
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
	if script != "" {
		cus.CustomizationScript = script
	}
	vm.SetGuestCustomizationSection(cus)
	err = vm.PowerOnAndForceCustomization()
	return err
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
func DestroyVAPP(vdc *govcd.Vdc, vAppName string, vApp *models.Vapp) (err error) {
	defer func() {
		db := config.GetDB()
		if err != nil {
			appMonitor, _ := controller.VappMonitors.LoadVapp(vApp.ID)
			appMonitor.Lock.Lock()
			vApp.Status = "Error"
			db.Save(vApp)
			appMonitor.Lock.Unlock()
		}
		vappQuery, _ := vdc.GetVAppByName(vAppName, true)
		if vappQuery == nil {
			for _, VM := range vApp.VMs {
				err := DeletePorts(vdc, VM.ID)
				if err != nil {
					appMonitor, _ := controller.VappMonitors.LoadVapp(vApp.ID)
					appMonitor.Lock.Lock()
					vApp.Status = "Error"
					db.Save(vApp)
					appMonitor.Lock.Unlock()
					return
				}
				db.Unscoped().Delete(&VM)
				if monitor, ok := controller.VMMonitors.Load(VM.ID); ok {
					monitor.Task.Stop()
					controller.VMMonitors.Delete(VM.ID)
				}
			}
			if _, ok := controller.VappMonitors.LoadVapp(vApp.ID); ok {
				controller.VappMonitors.Delete(vApp.ID)
			}
			db.Unscoped().Delete(vApp)
		}
	}()
	vappQuery, err := vdc.GetVAppByName(vAppName, true)
	if vappQuery == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return nil
	}
	if vappQuery.VApp.Deployed {
		task, err := vappQuery.Undeploy()
		if err != nil {
			fmt.Println(err)
			return err
		}
		err = task.WaitTaskCompletion()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	task, err := vappQuery.Delete()
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = task.WaitTaskCompletion()
	if err != nil {
		fmt.Println(err)
		return err
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

	conf := config.VcdConfig{
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

	vappMonitor := controller.NewVAppMonitor(newVapp.ID)
	controller.VappMonitors.StoreVapp(newVapp.ID, vappMonitor)

	for _, VM := range tem.VMTemps {
		ports, err := CheckPorts(VM.Ports)
		if err != nil {
			fmt.Println("Wrong format of Ports: " + err.Error())
			return vapp.NewAddVappNotFound().WithPayload(&vapp.AddVappNotFoundBody{
				Message: "Create vdc failed: ports format wrong",
			})
		}
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
		monitor := controller.NewVMMonitor(newVM.ID, &conf)
		controller.VMMonitors.Store(newVM.ID, monitor)
		for _, port := range ports {
			prefix := newVapp.Name + "-" + newVM.Name + "-" + strconv.Itoa(port)
			newPort := models.Port{
				Port: uint(port),
				URL: prefix + "." + config.URLSuffix ,
				VMachineID: newVM.ID,
			}
			db.Create(&newPort)
		}
	}

	go DeployVAPP(client, org, vdc, tem.Name, tem.VMTemps, res.Catalog, body.Name, res.Network, newVapp.ID)
	return vapp.NewAddVappOK().WithPayload(&vapp.AddVappOKBody{
		ID: int64(uid),
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

	conf := config.VcdConfig{
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

	appMonitor, _ := controller.VappMonitors.LoadVapp(vApp.ID)
	appMonitor.Lock.Lock()

	if db.Preload("VMs").Where("id = ?", params.ID).First(&vApp).RowsAffected == 0 {
		return vapp.NewDeleteVappNotFound()
	}

	if(vApp.Status == "Creating" || vApp.Status == "Deleting") {
		appMonitor.Lock.Unlock()
		return vapp.NewDeleteVappForbidden()
	}

	vApp.Status = "Deleting"
	db.Save(&vApp)

	appMonitor.Lock.Unlock()

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

func CheckPorts(ports string) ([]int, error){
	strings := strings.Split(ports, ",")
	var vals []int
	if ports == "" {
		return vals, nil
	}
	for _, port := range strings {
		val, err := strconv.Atoi(port)
		if err != nil || val < 200 || val > 65535{
			return vals, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

func ExposePorts(vdc *govcd.Vdc, vmID int, VAppName string) error {
	db := config.GetDB()
	var vm models.VMachine
	db.Preload("Ports").Where("id = ?", vmID).First(&vm)
	var ruleName string
	if vm.Name == "Party1" || vm.Name == "Party2" {
		ruleName = "kubefate"
	} else {
		ruleName = "cloudtides"
	}
	for _, port := range vm.Ports {
		prefix := VAppName + "-" + vm.Name + "-" + strconv.Itoa(int(port.Port))
		gateway, err := vdc.GetEdgeGatewayByName("edge-cn-bj", true)
		if err != nil {
			fmt.Println(err)
			return err
		}
		rule, err := gateway.GetLbAppRuleByName(ruleName)
		if err != nil {
			fmt.Println(err)
			return err
		} else {
			_, err := gateway.CreateLbServerPool(&types.LbPool{
				Name:        prefix,
				Algorithm:   "round-robin",
				Transparent: false,
				MonitorId:   "monitor-1",
				Members: []types.LbPoolMember{
					{
						Name:        prefix,
						MonitorPort: int(port.Port),
						Port:        int(port.Port),
						Weight:      1,
						IpAddress:	 vm.IPAddress,
						Condition:   "enabled",
					},
				},
			})
			if err != nil {
				fmt.Println(err)
				return err
			}
			rule.Script += fmt.Sprintf("\nacl is_%s hdr(host) -i %s \nuse_backend %s if is_%s", prefix,
				port.URL, prefix, prefix)
			fmt.Println(rule.Script)
			_, err = gateway.UpdateLbAppRule(rule)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return nil
}

func DeletePorts(vdc *govcd.Vdc, vmID uint) error {
	db := config.GetDB()
	var vm models.VMachine
	var VApp models.Vapp
	var ruleName string
	db.Preload("Ports").Where("id = ?", vmID).First(&vm)
	db.Where("id = ?", vm.VappID).First(&VApp)
	if vm.Name == "Party1" || vm.Name == "Party2" {
		ruleName = "kubefate"
	} else {
		ruleName = "cloudtides"
	}
	gateway, err := vdc.GetEdgeGatewayByName("edge-cn-bj", true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	rule, err := gateway.GetLbAppRuleByName(ruleName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	lines := strings.Split(rule.Script, "\n")
	newScript := ""
	outloop:
	for _, line := range lines {
		for _, port := range vm.Ports {
			prefix := VApp.Name + "-" + vm.Name + "-" + strconv.Itoa(int(port.Port))
			if strings.Contains(line, prefix) {
				continue outloop
			}
		}
		newScript += line + "\n"
	}
	rule.Script = newScript
	_, err = gateway.UpdateLbAppRule(rule)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, port := range vm.Ports {
		prefix := VApp.Name + "-" + vm.Name + "-" + strconv.Itoa(int(port.Port))
		_, err = gateway.GetLbServerPoolByName(prefix)
		if err != nil {
			continue
		}
		err = gateway.DeleteLbServerPoolByName(prefix)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	for _, port := range vm.Ports {
		db.Unscoped().Delete(&port)
	}
	return nil
}
