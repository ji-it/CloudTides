// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"
	"tides-server/pkg/restapi/operations/port"
	"tides-server/pkg/restapi/operations/vapp"
	"tides-server/pkg/restapi/operations/vendor_swagger"
	"tides-server/pkg/restapi/operations/vm"
	"tides-server/pkg/restapi/operations/vmtemp"

	interpose "github.com/carbocation/interpose/middleware"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/rs/cors"

	"tides-server/pkg/restapi/operations"
	"tides-server/pkg/restapi/operations/policy"
	"tides-server/pkg/restapi/operations/project"
	"tides-server/pkg/restapi/operations/resource"
	"tides-server/pkg/restapi/operations/template"
	"tides-server/pkg/restapi/operations/usage"
	"tides-server/pkg/restapi/operations/user"

	"tides-server/pkg/handler"
)

//go:generate swagger generate server --target ../../pkg --name CloudTides --spec ../../swagger/swagger.yml --exclude-main

func configureFlags(api *operations.CloudTidesAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CloudTidesAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.UsageAddResourceUsageHandler = usage.AddResourceUsageHandlerFunc(handler.AddResourceUsageHandler)

	api.PolicyAddPolicyHandler = policy.AddPolicyHandlerFunc(handler.AddPolicyHandler)

	api.ResourceAddVsphereResourceHandler = resource.AddVsphereResourceHandlerFunc(handler.AddVsphereResourceHandler)

	api.TemplateAddTemplateHandler = template.AddTemplateHandlerFunc(handler.AddTemplateHandler)

	api.UsageAddVMUsageHandler = usage.AddVMUsageHandlerFunc(handler.AddVMUsageHandler)

	api.UsageDeleteResourceUsageHandler = usage.DeleteResourceUsageHandlerFunc(handler.DeleteResourceUsageHandler)

	api.TemplateDeleteTemplateHandler = template.DeleteTemplateHandlerFunc(handler.DeleteTemplateHandler)

	api.PolicyListPolicyHandler = policy.ListPolicyHandlerFunc(handler.ListPolicyHandler)

	api.ResourceListVsphereResourceHandler = resource.ListVsphereResourceHandlerFunc(handler.ListVsphereResourceHandler)

	api.TemplateListTemplateHandler = template.ListTemplateHandlerFunc(handler.ListTemplateHandler)

	api.UserRegisterUserHandler = user.RegisterUserHandlerFunc(handler.RegisterUserHandler)

	api.PolicyRemovePolicyHandler = policy.RemovePolicyHandlerFunc(handler.RemovePolicyHandler)

	api.UsageUpdateResourceUsageHandler = usage.UpdateResourceUsageHandlerFunc(handler.UpdateResourceUsageHandler)

	api.PolicyUpdatePolicyHandler = policy.UpdatePolicyHandlerFunc(handler.UpdatePolicyHandler)

	api.UserGetUserProfileHandler = user.GetUserProfileHandlerFunc(handler.GetUserProfileHandler)

	api.UserUpdateUserProfileHandler = user.UpdateUserProfileHandlerFunc(handler.UpdateUserProfileHandler)

	api.UserUserLoginHandler = user.UserLoginHandlerFunc(handler.UserLoginHandler)

	api.ResourceValidateVsphereResourceHandler = resource.ValidateVsphereResourceHandlerFunc(handler.ValidateVsphereResourceHandler)

	api.ProjectAddProjectHandler = project.AddProjectHandlerFunc(handler.AddProjectHandler)

	api.ProjectListProjectHandler = project.ListProjectHandlerFunc(handler.ListProjectHandler)

	api.ProjectUpdateProjectHandler = project.UpdateProjectHandlerFunc(handler.UpdateProjectHandler)

	api.ProjectDeleteProjectHandler = project.DeleteProjectHandlerFunc(handler.DeleteProjectHandler)

	api.ResourceValidateVcdResourceHandler = resource.ValidateVcdResourceHandlerFunc(handler.ValidateVcdResourceHandler)

	api.ResourceAddVcdResourceHandler = resource.AddVcdResourceHandlerFunc(handler.AddVcdResourceHandler)

	api.ResourceListVcdResourceHandler = resource.ListVcdResourceHandlerFunc(handler.ListVcdResourceHandler)

	api.ResourceGetVcdResourceHandler = resource.GetVcdResourceHandlerFunc(handler.GetVcdResourceHandler)

	api.ResourceDeleteVcdResourceHandler = resource.DeleteVcdResourceHandlerFunc(handler.DeleteVcdResourceHandler)

	api.ResourceAssignPolicyHandler = resource.AssignPolicyHandlerFunc(handler.AssignPolicyHandler)

	api.PolicyGetPolicyHandler = policy.GetPolicyHandlerFunc(handler.GetPolicyHandler)

	api.UsageGetResourceUsageHandler = usage.GetResourceUsageHandlerFunc(handler.GetResourceUsageHandler)

	api.UsageGetPastUsageHandler = usage.GetPastUsageHandlerFunc(handler.GetPastUsageHandler)

	api.ResourceActivateResourceHandler = resource.ActivateResourceHandlerFunc(handler.ActivateResourceHandler)

	api.ResourceContributeResourceHandler = resource.ContributeResourceHandlerFunc(handler.ContributeResourceHandler)

	api.VendorSwaggerListVendorHandler = vendor_swagger.ListVendorHandlerFunc(handler.ListVendorsHandler)

	api.VendorSwaggerAddVendorHandler = vendor_swagger.AddVendorHandlerFunc(handler.AddVendorHandler)

	api.VendorSwaggerDeleteVendorHandler = vendor_swagger.DeleteVendorHandlerFunc((handler.DeleteVendorHandler))

	api.VappAddVappHandler = vapp.AddVappHandlerFunc(handler.AddVAPPHandler)

	api.VappListVappsHandler = vapp.ListVappsHandlerFunc(handler.ListVappHandler)

	api.VappDeleteVappHandler = vapp.DeleteVappHandlerFunc(handler.DeleteVAPPHandler)

	api.VmtempAddVMTempHandler = vmtemp.AddVMTempHandlerFunc(handler.AddVMTemplateHandler)

	api.VmtempListVMTempHandler = vmtemp.ListVMTempHandlerFunc(handler.ListVMTemplateHandler)

	api.VmtempDeleteVMTempHandler = vmtemp.DeleteVMTempHandlerFunc(handler.DeleteVMTemplateHandler)

	api.VmtempUpdateVMTempHandler = vmtemp.UpdateVMTempHandlerFunc(handler.UpdateVMTemplateHandler)

	api.VMListVMHandler = vm.ListVMHandlerFunc(handler.ListVMHandler)

	api.PortListPortsHandler = port.ListPortsHandlerFunc(handler.ListPortsHandler)

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200", "http://www.cloudtides.org.cn", "http://106.15.92.155"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Access-Control-Allow-Origin", "Authorization", "Content-Type"},
		AllowedMethods:   []string{"PUT", "GET", "POST", "PATCH", "DELETE"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	// Insert the middleware
	handler = c.Handler(handler)
	logViaLogrus := interpose.NegroniLogrus()

	return logViaLogrus(handler)
}
