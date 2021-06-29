package handler

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vmware/go-vcloud-director/v2/govcd"
	"github.com/vmware/go-vcloud-director/v2/types/v56"

	"tides-server/pkg/config"
	"tides-server/pkg/logger"
	"tides-server/pkg/models"
)

// ParseUserIDFromToken parses user ID in database from JWT supplied by client
func ParseUserIDFromToken(req *http.Request) (uint, error) {
	reqToken := req.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		logger.SetLogLevel("ERROR")
		logger.Error("Token not supplied in request")
		return 0, errors.New("Token not supplied in request")
	}
	stringToken := splitToken[1]

	claims := &Claims{}
	secretKey := config.GetConfig().SecretKey
	_, err := jwt.ParseWithClaims(
		stringToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		logger.SetLogLevel("ERROR")
		logger.Error("JWT is expired")
		return 0, errors.New("JWT is expired")
	}
	return claims.ID, err
}

// VerifyUser checks if user is registered
func VerifyUser(req *http.Request) bool {
	id, err := ParseUserIDFromToken(req)
	if err != nil {
		return false
	}
	db := config.GetDB()
	var queryUser models.User
	if db.Where("id = ?", id).First(&queryUser).Error != nil {
		return false
	}

	return true
}

// VerifyAdmin checks if user is admin
func VerifyAdmin(req *http.Request) bool {
	id, err := ParseUserIDFromToken(req)
	if err != nil {
		return false
	}
	db := config.GetDB()
	var queryUser models.User
	if db.Where("id = ?", id).First(&queryUser).Error != nil {
		return false
	}
	if queryUser.Priority == models.UserPriorityHigh {
		return true
	}
	return false
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
		// fmt.Printf("Token: %s\n", resp.Header[govcd.AuthorizationHeader])
	}
	return vcdClient, nil
}*/

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func initValidation(conf *config.VcdConfig, catalog string, network string, res *models.Resource) {
	db := config.GetDB()
	client, err := conf.Client() // We now have a client
	if err != nil {
		fmt.Println(err)
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	org, err := client.GetOrgByName(conf.Org)
	if err != nil {
		fmt.Println(err)
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	vdc, err := org.GetVDCByName(conf.VDC, false)
	if err != nil {
		fmt.Println(err)
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}

	deployVapp(org, vdc, catalog, vappName, network, res)
}

// Deploy VAPP
func deployVapp(org *govcd.Org, vdc *govcd.Vdc, cataName string, vAppName string, netName string, res *models.Resource) {
	db := config.GetDB()

	catalog, err := org.GetCatalogByName(cataName, true)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	cataItem, err := catalog.GetCatalogItemByName(temName, true)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	vappTem, err := cataItem.GetVAppTemplate()
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	net, err := vdc.GetOrgVdcNetworkByName(netName, true)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	networks := []*types.OrgVDCNetwork{}

	networks = append(networks, net.OrgVDCNetwork)

	storageProf := vdc.Vdc.VdcStorageProfiles.VdcStorageProfile[0]

	task, err := vdc.ComposeVApp(networks, vappTem, *storageProf, vAppName, "test purpose", true)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	task.WaitTaskCompletion()

	vapp, err := vdc.GetVAppByName(vAppName, true)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	task, err = vapp.PowerOn()
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	task.WaitTaskCompletion()

	vm, err := vapp.GetVMByName(vmName, true)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}

	task, err = vm.Undeploy()
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	task.WaitTaskCompletion()

	task, err = vm.ChangeCPUCount(2)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	task.WaitTaskCompletion()
	task, err = vm.ChangeMemorySize(2048)
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}
	task.WaitTaskCompletion()

	cus, _ := vm.GetGuestCustomizationSection()
	cus.Enabled = new(bool)
	*cus.Enabled = true
	// cus.ComputerName = "tides-" + randSeq(5)
	vm.SetGuestCustomizationSection(cus)
	err = vm.PowerOnAndForceCustomization()
	if err != nil {
		res.SetupStatus = err.Error()
		db.Save(&res)
		return
	}

	res.SetupStatus = "Waiting for Validation"
	db.Save(&res)
}

func initDestruction(conf *config.VcdConfig) {
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

	for i := 0; i < 3; i++ {
		err = destroyVapp(vdc, vappName)
		vappQuery, _ := vdc.GetVAppByName(vappName, true)
		if err == nil && vappQuery == nil {
			break
		}
	}
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

// PolicySetup sets up a default policy for vcd resources
func PolicySetup(rid uint, uid uint, network string, catalog string) {
	db := config.GetDB()
	var tem models.Template
	if db.Where("name = ?", "tides-boinc-attached").First(&tem).RowsAffected == 0 {
		return
	}
	var pol models.Policy
	if db.Where("name = ?", "tides-policy").First(&pol).RowsAffected == 0 {
		pol = models.Policy{
			IdlePolicy:      idlePolicy,
			IsDestroy:       true,
			Name:            "tides-policy",
			PlatformType:    models.ResourcePlatformTypeVcd,
			TemplateID:      tem.ID,
			ThresholdPolicy: thresholdPolicy,
			UserID:          uid,
		}
		db.Create(&pol)
		vcdPolicy := models.VcdPolicy{
			Catalog:  catalog,
			Network:  network,
			PolicyID: pol.ID,
		}
		db.Create(&vcdPolicy)
	}
	var res models.Resource
	db.Where("id = ?", rid).First(&res)
	if res.PolicyID == nil {
		res.PolicyID = new(uint)
		*res.PolicyID = pol.ID
	} else {
		*res.PolicyID = pol.ID
	}
	db.Save(&res)
}
