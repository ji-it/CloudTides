package handler

import (
	// "crypto/tls"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// "reflect"

	// "github.com/go-openapi/errors"
	// "github.com/go-openapi/runtime"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-openapi/runtime/middleware"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"

	// "tides-server/pkg/restapi/operations"
	// "tides-server/pkg/restapi/operations/policy"

	"tides-server/pkg/restapi/operations/policy"
	"tides-server/pkg/restapi/operations/resource"
	"tides-server/pkg/restapi/operations/template"

	"tides-server/pkg/restapi/operations/usage"
	"tides-server/pkg/restapi/operations/user"

	"tides-server/pkg/config"
	"tides-server/pkg/models"
)

func ParseUserIdFromToken(req *http.Request) (uint, error) {
	reqToken := req.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) < 2 {
		return 0, errors.New("Token not supplied in request!")
	}
	stringToken := splitToken[1]

	claims := &Claims{}
	_, err := jwt.ParseWithClaims(
		stringToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		},
	)
	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return 0, errors.New("JWT is expired")
	}
	return claims.Id, err
}

func VerifyUser(req *http.Request) bool {
	id, err := ParseUserIdFromToken(req)
	if err != nil {
		return false
	}
	db := config.GetDB()
	var queryUser models.User
	db.Where("id = ?", id).First(&queryUser)
	if queryUser.Username == "" {
		return false
	}
	return true
}

func RegisterUserHandler(params user.RegisterUserParams) middleware.Responder {
	body := params.ReqBody
	db := config.GetDB()
	var queryUser models.User
	db.Where("username = ?", body.Username).First(&queryUser)
	if queryUser.Username != "" {
		return user.NewRegisterUserBadRequest().WithPayload(&user.RegisterUserBadRequestBody{Message: "Username already used!"})
	}

	res := &user.RegisterUserOKBodyUserInfo{
		CompanyName: body.CompanyName,
		Password:    body.Password,
		Priority:    body.Priority,
		Username:    body.Username,
	}
	newUser := models.User{Username: body.Username, Password: body.Password, CompanyName: body.CompanyName, Priority: body.Priority}

	err := db.Create(&newUser).Error
	if err != nil {
		return user.NewRegisterUserBadRequest()
	}

	return user.NewRegisterUserOK().WithPayload(&user.RegisterUserOKBody{UserInfo: res})
}

func UserLoginHandler(params user.UserLoginParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var queryUser models.User
	db.Where("Username = ?", body.Username).First(&queryUser)
	if queryUser.Username == "" || queryUser.Password != body.Password {
		return user.NewUserLoginUnauthorized()
	}
	expirationTime := time.Now().Add(expireTime)
	claims := Claims{
		Id: queryUser.Model.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secretKey))

	res := user.UserLoginOKBodyUserInfo{Priority: queryUser.Priority, Username: queryUser.Username}
	return user.NewUserLoginOK().WithPayload(&user.UserLoginOKBody{Token: signedToken, UserInfo: &res})
}

func UserDetailsHandler(params user.UserDetailsParams) middleware.Responder {
	id, err := ParseUserIdFromToken(params.HTTPRequest)
	if err != nil {
		return user.NewUserDetailsUnauthorized()
	}
	db := config.GetDB()
	var queryUser models.User
	db.Where("id = ?", id).First(&queryUser)
	if queryUser.Username == "" {
		return user.NewUserDetailsUnauthorized()
	}

	res := user.UserDetailsOKBodyResults{
		City:        queryUser.City,
		CompanyName: queryUser.CompanyName,
		Country:     queryUser.Country,
		Email:       queryUser.Email,
		FirstName:   queryUser.FirstName,
		LastName:    queryUser.LastName,
		Position:    queryUser.Position,
	}

	return user.NewUserDetailsOK().WithPayload(&user.UserDetailsOKBody{Message: "success", Results: &res})
}

// Still needs modification
func AddTemplateHandler(params template.AddTemplateParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return template.NewAddTemplateUnauthorized()
	}

	req := params.HTTPRequest
	req.ParseMultipartForm(32 << 20)
	file, handler, err := req.FormFile("File")
	if err != nil {
		fmt.Println(err)
		return template.NewAddTemplateUnauthorized()
	}
	defer file.Close()
	f, err := os.OpenFile("./"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()
	io.Copy(f, file)
	// var name string
	/*if req.Form["Source"] == "upload" {
		fmt.Println(req.Form["File"])
	}*/
	return template.NewAddTemplateOK()
}

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
			PolicyRef:    uint(body.Policy),
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

	res.PolicyRef = uint(pid)
	db.Save(&res)

	return resource.NewAssignPolicyOK().WithPayload(&resource.AssignPolicyOKBody{
		Message: "success",
	})
}

func DestroyVMHandler(params resource.DestroyVMParams) middleware.Responder {
	ip := params.ReqBody.IPAddress
	var vm models.VM

	db := config.GetDB()
	if db.Where("ip_address = ?", ip).First(&vm).RecordNotFound() {
		return resource.NewDestroyVMNotFound()
	}

	vm.IsDestroyed = true
	vm.PoweredOn = false
	db.Save(&vm)

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
	cost := 0
	totalActiveVM := 0
	for _, res := range resources {
		var activeVMs int
		vms := []*models.VM{}
		totalVMs += int(res.TotalVMs)
		db.Where("resource_ref = ? AND is_destroyed = ? AND powered_on = ?",
			res.Model.ID, false, true).Find(&vms).Count(&activeVMs)
		totalActiveVM += activeVMs
	}

}
*/

func AddPolicyHandler(params policy.AddPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewAddPolicyUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)
	body := params.ReqBody

	newPolicy := models.Policy{
		AccountType:     body.AccountType,
		DeployType:      body.DeployType,
		IdlePolicy:      body.Idle,
		IsDestroy:       body.IsDestroy,
		Name:            body.Name,
		ProjectRef:      uint(body.ProjectID),
		TemplateRef:     uint(body.TemplateID),
		ThresholdPolicy: body.Threshold,
		UserRef:         uid,
	}

	db := config.GetDB()
	err := db.Create(&newPolicy).Error

	if err != nil {
		return policy.NewAddPolicyBadRequest()
	}

	return policy.NewAddPolicyOK().WithPayload(&policy.AddPolicyOKBody{
		Message: "success",
	})
}

func UpdatePolicyHandler(params policy.UpdatePolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewAddPolicyUnauthorized()
	}

	body := params.ReqBody
	var pol models.Policy
	db := config.GetDB()
	if db.Where("id = ?", body.ID).First(&pol).RecordNotFound() {
		return policy.NewUpdatePolicyNotFound()
	}

	pol.AccountType = body.AccountType
	pol.DeployType = body.DeployType
	pol.IdlePolicy = body.Idle
	pol.IsDestroy = body.IsDestroy
	pol.Name = body.Name
	pol.ProjectRef = uint(body.ProjectID)
	pol.TemplateRef = uint(body.TemplateID)
	pol.ThresholdPolicy = body.Threshold

	err := db.Save(&pol).Error
	if err != nil {
		return policy.NewUpdatePolicyBadRequest()
	}

	return policy.NewUpdatePolicyOK().WithPayload(&policy.UpdatePolicyOKBody{
		Message: "success",
	})
}

func ListPolicyHandler(params policy.ListPolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewListPolicyUnauthorized()
	}

	uid, _ := ParseUserIdFromToken(params.HTTPRequest)

	policies := []*models.Policy{}
	db := config.GetDB()

	db.Where("user_ref = ?", uid).Find(&policies)
	if len(policies) == 0 {
		return policy.NewListPolicyNotFound()
	}

	results := []*policy.ResultsItems0{}
	for _, pol := range policies {
		resources := []*models.Resource{}
		var pro models.Project
		db.Where("policy_ref = ?", pol.Model.ID).Find(&resources)
		db.Where("id = ?", pol.ProjectRef).First(&pro)

		newResult := policy.ResultsItems0{
			DeployType:      pol.DeployType,
			HostsAssigned:   int64(len(resources)),
			ID:              int64(pol.Model.ID),
			IdlePolicy:      pol.IdlePolicy,
			IsDestroy:       pol.IsDestroy,
			Name:            pol.Name,
			ProjectName:     pro.ProjectName,
			ThresholdPolicy: pol.ThresholdPolicy,
		}

		results = append(results, &newResult)
	}

	return policy.NewListPolicyOK().WithPayload(&policy.ListPolicyOKBody{
		Message: "success",
		Results: results,
	})
}

func RemovePolicyHandler(params policy.RemovePolicyParams) middleware.Responder {
	if !VerifyUser(params.HTTPRequest) {
		return policy.NewRemovePolicyUnauthorized()
	}

	body := params.ReqBody
	var pol models.Policy
	db := config.GetDB()

	err := db.Unscoped().Where("id = ?", body.ID).Delete(&pol).Error
	if err != nil {
		return policy.NewRemovePolicyNotFound()
	}

	return policy.NewRemovePolicyOK().WithPayload(&policy.RemovePolicyOKBody{
		Message: "success",
	})
}

func AddHostUsageHandler(params usage.AddHostUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var res models.Resource

	if db.Where("host_address = ? AND host_name = ?", body.HostAddress, body.HostName).First(&res).RecordNotFound() {
		return usage.NewAddHostUsageNotFound()
	}

	PercentCPU := body.CurrentCPU / body.TotalCPU
	PercentRAM := body.CurrentRAM / body.TotalRAM

	newHostUsage := models.HostUsage{
		CurrentCPU:  body.CurrentCPU,
		CurrentRAM:  body.CurrentRAM,
		HostAddress: body.HostAddress,
		HostName:    body.HostName,
		PercentCPU:  PercentCPU,
		PercentRAM:  PercentRAM,
		TotalCPU:    body.TotalCPU,
		TotalRAM:    body.TotalRAM,
		ResourceRef: res.Model.ID,
	}

	err := db.Create(&newHostUsage).Error
	if err != nil {
		return usage.NewAddHostUsageBadRequest()
	}

	return usage.NewAddHostUsageOK().WithPayload(&usage.AddHostUsageOKBody{
		Message: "success",
	})
}

func UpdateHostUsageHandler(params usage.UpdateHostUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()

	var hu models.HostUsage
	if db.Where("host_address = ? AND host_name = ?", body.HostAddress, body.HostName).First(&hu).RecordNotFound() {
		return usage.NewUpdateHostUsageNotFound()
	}

	hu.CurrentCPU = body.CurrentCPU
	hu.CurrentRAM = body.CurrentRAM
	hu.PercentCPU = body.CurrentCPU / hu.TotalCPU
	hu.PercentRAM = body.CurrentRAM / hu.TotalRAM

	db.Save(&hu)

	return usage.NewUpdateHostUsageOK().WithPayload(&usage.UpdateHostUsageOKBody{
		Message: "host usage recorded",
	})
}

func DeleteHostUsageHandler(params usage.DeleteHostUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()
	var hu models.HostUsage

	err := db.Unscoped().Where("host_address = ? AND host_name = ?", body.HostAddress, body.HostName).Delete(&hu).Error
	if err != nil {
		return usage.NewDeleteHostUsageBadRequest()
	}

	return usage.NewDeleteHostUsageOK().WithPayload(&usage.DeleteHostUsageOKBody{
		Message: "success",
	})
}

func AddVMUsageHandler(params usage.AddVMUsageParams) middleware.Responder {
	body := params.ReqBody

	db := config.GetDB()

	for ip, val := range body {
		var vm models.VM
		if db.Where("ip_address = ?", ip).First(&vm).RecordNotFound() {
			var host models.Resource
			err := db.Where("host_name = ?", val.HostName).First(&host).Error
			if err != nil {
				return usage.NewAddVMUsageBadRequest()
			}

			boincTime, _ := time.Parse("Mon Jan 2 15:04:05 -0700 MST 2006", val.BoincStartTime)

			newvm := models.VM{
				BoincTime:   boincTime,
				GuestOS:     val.GuestOS,
				IPAddress:   ip,
				IsDestroyed: false,
				Name:        val.Name,
				NumCPU:      val.NumCPU,
				PoweredOn:   val.PoweredOn,
				ResourceRef: host.Model.ID,
			}

			db.Create(&newvm)

			newvmUsage := models.VMUsage{
				CurrentCPU: val.CurrentCPU,
				CurrentRAM: val.CurrentRAM,
				TotalCPU:   val.TotalCPU,
				TotalRAM:   val.TotalRAM,
				VmRef:      newvm.Model.ID,
			}

			db.Create(&newvmUsage)
		} else {
			var vmu models.VMUsage
			db.Where("vm_ref = ?", vm.Model.ID).First(&vmu)

			vmu.CurrentCPU = val.CurrentCPU
			vmu.CurrentRAM = val.CurrentRAM

			db.Save(&vm)
		}
	}

	return usage.NewAddVMUsageOK().WithPayload(&usage.AddVMUsageOKBody{
		Message: "success",
	})
}
