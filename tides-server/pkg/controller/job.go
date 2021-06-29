package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"tides-server/pkg/config"
	"tides-server/pkg/models"
	"time"

	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"
	"gopkg.in/yaml.v2"
)

func randSeq(n int) string {
	b := make([]rune, n)
	t := time.Now()
	rand.Seed(int64(t.Nanosecond()))
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Checks that a configuration structure is complete
func checkConfiguration(conf config.VcdConfig) {
	willExit := false
	abort := func(s string) {
		fmt.Printf("configuration field '%s' empty or missing\n", s)
		willExit = true
	}
	if conf.Org == "" {
		abort("org")
	}
	if conf.Href == "" || conf.Href == "https://YOUR_VCD_IP/api" {
		abort("href")
	}
	if conf.VDC == "" {
		abort("vdc")
	}
	if conf.Token != "" {
		return
	}
	if conf.User == "" {
		abort("user")
	}
	if conf.Password == "" {
		abort("password")
	}
	if willExit {
		return
	}
}

// Retrieves the configuration from a Json or Yaml file
func getConfig(configFile string) (config.VcdConfig, error) {
	var configuration config.VcdConfig
	buffer, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("Configuration file %s not found\n%s\n", configFile, err)
		return configuration, err
	}
	err = yaml.Unmarshal(buffer, &configuration)
	if err != nil {
		fmt.Printf("Error retrieving configuration from file %s\n%s\n", configFile, err)
		return configuration, err
	}
	checkConfiguration(configuration)

	// If something goes wrong, rerun the program after setting
	// the environment variable SAMPLES_DEBUG, and you can check how the
	// configuration was read
	if os.Getenv("SAMPLES_DEBUG") != "" {
		fmt.Printf("configuration text: %s\n", buffer)
		fmt.Printf("configuration rec: %#v\n", configuration)
		newConf, _ := yaml.Marshal(configuration)
		fmt.Printf("YAML configuration: \n%s\n", newConf)
	}
	return configuration, nil
}

// Client creates a vCD client
/*func (c *VcdConfig) Client() (*govcd.VCDClient, error) {
	u, err := url.ParseRequestURI(c.Href)
	if err != nil {
		return nil, fmt.Errorf("unable to pass url: %s", err)
	}

	vcdClient := govcd.NewVCDClient(*u, c.Insecure)
	if c.Token != "" {
		_ = vcdClient.SetToken(c.Org, govcd.AuthorizationHeader, c.Token)
	} else {
		_, err := vcdClient.GetAuthResponse(c.User, c.Password, c.Org)
		if err != nil {
			return nil, fmt.Errorf("unable to authenticate: %s", err)
		}
	}
	return vcdClient, nil
}*/

// Deploy VAPP
func deployVapp(org *govcd.Org, vdc *govcd.Vdc, temName string, VMName string, cataName string, vAppName string, netName string) *govcd.VApp {

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
func powerOnVapp(vdc *govcd.Vdc, vAppName string) error {
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
func undeployVapp(vdc *govcd.Vdc, vAppName string) error {
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
func destroyVapp(vdc *govcd.Vdc, vAppName string) error {
	vapp, err := vdc.GetVAppByName(vAppName, true)
	if vapp == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return err
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

// RunJob is a cronjob for resource. Query usage, update status, deploy/destroy/suspend Vapps.
func RunJob(configFile string) {

	// Reads the configuration file
	conf, err := getConfig(configFile)
	if err != nil {
		return
	}

	client, err := conf.Client() // We now have a client
	if err != nil {
		fmt.Println(err)
		return
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		fmt.Println(err)
		return
	}
	vdc, err := org.GetVDCByName(conf.VDC, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	db := config.GetDB()
	var res models.Resource
	db.Where("host_address = ?", conf.Href).First(&res)
	var resUsage models.ResourceUsage
	db.Where("resource_id = ?", res.ID).First(&resUsage)

	// Update usage
	currentCPU := float64(vdc.Vdc.ComputeCapacity[0].CPU.Used)
	currentRAM := float64(vdc.Vdc.ComputeCapacity[0].Memory.Used)
	totalCPU := float64(vdc.Vdc.ComputeCapacity[0].CPU.Limit)
	totalRAM := float64(vdc.Vdc.ComputeCapacity[0].Memory.Limit)
	storageRef := vdc.Vdc.VdcStorageProfiles.VdcStorageProfile[0].HREF
	storage, err := govcd.GetStorageProfileByHref(client, storageRef)
	currentDisk := float64(storage.StorageUsedMB)
	totalDisk := float64(storage.Limit)
	resUsage.CurrentCPU = currentCPU
	resUsage.CurrentRAM = currentRAM
	resUsage.TotalCPU = totalCPU
	resUsage.TotalRAM = totalRAM
	resUsage.PercentCPU = currentCPU / totalCPU
	resUsage.PercentRAM = currentRAM / totalRAM
	resUsage.CurrentDisk = currentDisk
	resUsage.TotalDisk = totalDisk
	resUsage.PercentDisk = currentDisk / totalDisk
	db.Save(&resUsage)

	newVcdPastUsage := models.ResourcePastUsage{
		CurrentCPU:  currentCPU,
		CurrentDisk: currentDisk,
		CurrentRAM:  currentRAM,
		PercentCPU:  currentCPU / totalCPU,
		PercentDisk: currentDisk / totalDisk,
		PercentRAM:  currentRAM / totalRAM,
		TotalCPU:    totalCPU,
		TotalDisk:   totalDisk,
		TotalRAM:    totalRAM,
		ResourceID:  res.ID,
	}
	db.Create(&newVcdPastUsage)

	var pol models.Policy
	if db.Where("id = ?", res.PolicyID).First(&pol).RowsAffected == 0 {
		return
	}
	idle := Policy{}
	thres := Policy{}
	json.Unmarshal([]byte(pol.IdlePolicy), &idle)
	json.Unmarshal([]byte(pol.ThresholdPolicy), &thres)

	if resUsage.PercentCPU < idle.CPU && resUsage.PercentRAM < idle.RAM {
		res.Status = "idle"
		db.Save(&res)
		if pol.PlatformType == models.ResourcePlatformTypeVcd {
			var susVapp models.VM
			if db.Where("resource_id = ? AND powered_on = ?", res.ID, false).First(&susVapp).RowsAffected > 0 {
				err := powerOnVapp(vdc, susVapp.Name)
				if err == nil {
					susVapp.PoweredOn = true
					db.Save(&susVapp)
					return
				}
			}
			var vcdPol models.VcdPolicy
			db.Where("policy_id = ?", pol.ID).First(&vcdPol)
			var tem models.Template
			db.Where("id = ?", pol.TemplateID).First(&tem)
			vapp := deployVapp(org, vdc, tem.Name, tem.VMName, vcdPol.Catalog, "cloudtides-vapp-"+randSeq(6), vcdPol.Network)
			if vapp != nil {
				newVapp := models.VM{
					IPAddress:   vapp.VApp.HREF,
					IsDestroyed: false,
					Name:        vapp.VApp.Name,
					PoweredOn:   true,
					ResourceID:  res.ID,
				}
				db.Create(&newVapp)
			}
		}
	} else if resUsage.PercentCPU > thres.CPU && resUsage.PercentRAM > thres.RAM {
		res.Status = "busy"
		db.Save(&res)
		if pol.PlatformType == models.ResourcePlatformTypeVcd {
			var vapp models.VM
			if db.Where("resource_id = ? AND powered_on = ?", res.ID, true).Last(&vapp).RowsAffected == 0 {
				return
			}
			if pol.IsDestroy {
				// fix destroy failure on Web UI
				for i := 0; i < 3; i++ {
					err = destroyVapp(vdc, vapp.Name)
					vappQuery, _ := vdc.GetVAppByName(vapp.Name, true)
					if err == nil && vappQuery == nil {
						db.Unscoped().Delete(&vapp)
						break
					}
				}
			} else {
				err := undeployVapp(vdc, vapp.Name)
				if err == nil {
					vapp.PoweredOn = false
					db.Save(&vapp)
				}
			}
		}
	} else {
		res.Status = "normal"
		db.Save(&res)
	}
}
