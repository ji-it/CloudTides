package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
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
func check_configuration(conf VcdConfig) {
	will_exit := false
	abort := func(s string) {
		fmt.Printf("configuration field '%s' empty or missing\n", s)
		will_exit = true
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
	if will_exit {
		os.Exit(1)
	}
}

// Retrieves the configuration from a Json or Yaml file
func getConfig(config_file string) VcdConfig {
	var configuration VcdConfig
	buffer, err := ioutil.ReadFile(config_file)
	if err != nil {
		fmt.Printf("Configuration file %s not found\n%s\n", config_file, err)
		os.Exit(1)
	}
	err = yaml.Unmarshal(buffer, &configuration)
	if err != nil {
		fmt.Printf("Error retrieving configuration from file %s\n%s\n", config_file, err)
		os.Exit(1)
	}
	check_configuration(configuration)

	// If something goes wrong, rerun the program after setting
	// the environment variable SAMPLES_DEBUG, and you can check how the
	// configuration was read
	if os.Getenv("SAMPLES_DEBUG") != "" {
		fmt.Printf("configuration text: %s\n", buffer)
		fmt.Printf("configuration rec: %#v\n", configuration)
		new_conf, _ := yaml.Marshal(configuration)
		fmt.Printf("YAML configuration: \n%s\n", new_conf)
	}
	return configuration
}

// Creates a vCD client
func (c *VcdConfig) Client() (*govcd.VCDClient, error) {
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
}

// Deploy VAPP
func deployVapp(org *govcd.Org, vdc *govcd.Vdc, temName string, cataName string, vAppName string, netName string, storageName string) *govcd.VApp {

	catalog, _ := org.GetCatalogByName(cataName, true)
	cataItem, _ := catalog.GetCatalogItemByName(temName, true)
	vappTem, _ := cataItem.GetVAppTemplate()
	net, err := vdc.GetOrgVdcNetworkByName(netName, true)
	networks := []*types.OrgVDCNetwork{}

	networks = append(networks, net.OrgVDCNetwork)

	storageProf := vdc.Vdc.VdcStorageProfiles.VdcStorageProfile[0]

	task, err := vdc.ComposeVApp(networks, vappTem, *storageProf, vAppName, "test purpose", true)
	task.WaitTaskCompletion()
	fmt.Println(err)
	vapp, err := vdc.GetVAppByName(vAppName, true)

	task, err = vapp.PowerOn()
	task.WaitTaskCompletion()

	return vapp
}

// Suspend VAPP
func suspendVapp(vdc *govcd.Vdc, vAppName string) {
	vapp, _ := vdc.GetVAppByName(vAppName, true)
	if vapp == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return
	}
	task, _ := vapp.Suspend()
	task.WaitTaskCompletion()
}

// Destroy VAPP
func destroyVapp(vdc *govcd.Vdc, vAppName string) {
	vapp, _ := vdc.GetVAppByName(vAppName, true)
	if vapp == nil {
		fmt.Println("Vapp " + vAppName + " not found")
		return
	}
	task, _ := vapp.Undeploy()
	task.WaitTaskCompletion()
	task, _ = vapp.Delete()
	task.WaitTaskCompletion()
}

// Cronjob for resource. Query usage, update status, deploy/destroy/suspend Vapps.
func RunJob(configFile string) {

	// Reads the configuration file
	conf := getConfig(configFile)

	client, err := conf.Client() // We now have a client
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	vdc, err := org.GetVDCByName(conf.VDC, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
	resUsage.CurrentCPU = currentCPU
	resUsage.CurrentRAM = currentRAM
	resUsage.TotalCPU = totalCPU
	resUsage.TotalRAM = totalRAM
	resUsage.PercentCPU = currentCPU / totalCPU
	resUsage.PercentRAM = currentRAM / totalRAM
	db.Save(&resUsage)

	newVcdPastUsage := models.ResourcePastUsage{
		CurrentCPU: currentCPU,
		CurrentRAM: currentRAM,
		PercentCPU: currentCPU / totalCPU,
		PercentRAM: currentRAM / totalRAM,
		TotalCPU:   totalCPU,
		TotalRAM:   totalRAM,
		ResourceID: res.ID,
	}
	db.Create(&newVcdPastUsage)

	var pol models.Policy
	db.Where("id = ?", res.PolicyID).First(&pol)
	idle := Policy{}
	thres := Policy{}
	json.Unmarshal([]byte(pol.IdlePolicy), &idle)
	json.Unmarshal([]byte(pol.ThresholdPolicy), &thres)

	if resUsage.PercentCPU < idle.CPU && resUsage.PercentRAM < idle.RAM {
		res.Status = "idle"
		db.Save(&res)
		if pol.PlatformType == models.ResourcePlatformTypeVcd {
			var vcdPol models.VcdPolicy
			db.Where("policy_id = ?", pol.ID).First(&vcdPol)
			var tem models.Template
			db.Where("id = ?", pol.TemplateID).First(&tem)
			vapp := deployVapp(org, vdc, tem.Name, vcdPol.Catalog, "cloudtides-vapp-"+randSeq(6), vcdPol.Network, vcdPol.Storage)
			newVapp := models.VM{
				IPAddress:   vapp.VApp.HREF,
				IsDestroyed: false,
				Name:        vapp.VApp.Name,
				PoweredOn:   true,
				ResourceID:  res.ID,
			}
			db.Create(&newVapp)
		}
	} else if resUsage.PercentCPU > thres.CPU && resUsage.PercentRAM > thres.RAM {
		res.Status = "busy"
		db.Save(&res)
		if pol.PlatformType == models.ResourcePlatformTypeVcd {
			var vapp models.VM
			db.Where("resource_id = ? AND powered_on = ?", res.ID, true).Last(&vapp)
			if pol.IsDestroy {
				destroyVapp(vdc, vapp.Name)
				db.Unscoped().Delete(&vapp)
			} else {
				suspendVapp(vdc, vapp.Name)
				vapp.PoweredOn = false
				db.Save(&vapp)
			}
		}
	} else {
		res.Status = "normal"
		db.Save(&res)
	}
}
