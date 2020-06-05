// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	// "strings"
	// "fmt"
	// "io/ioutil"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"tides-server/pkg/restapi/operations"
	"tides-server/pkg/restapi/operations/policy"
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

	if api.UsageAddHostUsageHandler == nil {
		api.UsageAddHostUsageHandler = usage.AddHostUsageHandlerFunc(func(params usage.AddHostUsageParams) middleware.Responder {
			return middleware.NotImplemented("operation usage.AddHostUsage has not yet been implemented")
		})
	} else {
		api.UsageAddHostUsageHandler = usage.AddHostUsageHandlerFunc(handler.AddHostUsageHandler)
	}
	if api.PolicyAddPolicyHandler == nil {
		api.PolicyAddPolicyHandler = policy.AddPolicyHandlerFunc(func(params policy.AddPolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation policy.AddPolicy has not yet been implemented")
		})
	} else {
		api.PolicyAddPolicyHandler = policy.AddPolicyHandlerFunc(handler.AddPolicyHandler)
	}
	if api.ResourceAddResourceHandler == nil {
		api.ResourceAddResourceHandler = resource.AddResourceHandlerFunc(func(params resource.AddResourceParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.AddResource has not yet been implemented")
		})
	} else {
		api.ResourceAddResourceHandler = resource.AddResourceHandlerFunc(handler.AddResourceHandler)
	}
	if api.TemplateAddTemplateHandler == nil {
		api.TemplateAddTemplateHandler = template.AddTemplateHandlerFunc(func(params template.AddTemplateParams) middleware.Responder {
			return middleware.NotImplemented("operation template.AddTemplate has not yet been implemented")
		})
	} else {
		api.TemplateAddTemplateHandler = template.AddTemplateHandlerFunc(handler.AddTemplateHandler)
	}
	if api.UsageAddVMUsageHandler == nil {
		api.UsageAddVMUsageHandler = usage.AddVMUsageHandlerFunc(func(params usage.AddVMUsageParams) middleware.Responder {
			return middleware.NotImplemented("operation usage.AddVMUsage has not yet been implemented")
		})
	} else {
		api.UsageAddVMUsageHandler = usage.AddVMUsageHandlerFunc(handler.AddVMUsageHandler)
	}
	if api.ResourceAssignPolicyHandler == nil {
		api.ResourceAssignPolicyHandler = resource.AssignPolicyHandlerFunc(func(params resource.AssignPolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.AssignPolicy has not yet been implemented")
		})
	} else {
		api.ResourceAssignPolicyHandler = resource.AssignPolicyHandlerFunc(handler.AssignPolicyHandler)
	}
	if api.UsageDeleteHostUsageHandler == nil {
		api.UsageDeleteHostUsageHandler = usage.DeleteHostUsageHandlerFunc(func(params usage.DeleteHostUsageParams) middleware.Responder {
			return middleware.NotImplemented("operation usage.DeleteHostUsage has not yet been implemented")
		})
	} else {
		api.UsageDeleteHostUsageHandler = usage.DeleteHostUsageHandlerFunc(handler.DeleteHostUsageHandler)
	}
	if api.ResourceDeleteResourceHandler == nil {
		api.ResourceDeleteResourceHandler = resource.DeleteResourceHandlerFunc(func(params resource.DeleteResourceParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.DeleteResource has not yet been implemented")
		})
	} else {
		api.ResourceDeleteResourceHandler = resource.DeleteResourceHandlerFunc(handler.DeleteResourceHandler)
	}
	if api.TemplateDeleteTemplateHandler == nil {
		api.TemplateDeleteTemplateHandler = template.DeleteTemplateHandlerFunc(func(params template.DeleteTemplateParams) middleware.Responder {
			return middleware.NotImplemented("operation template.DeleteTemplate has not yet been implemented")
		})
	}
	if api.ResourceDestroyVMHandler == nil {
		api.ResourceDestroyVMHandler = resource.DestroyVMHandlerFunc(func(params resource.DestroyVMParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.DestroyVM has not yet been implemented")
		})
	} else {
		api.ResourceDestroyVMHandler = resource.DestroyVMHandlerFunc(handler.DestroyVMHandler)
	}
	if api.UsageHostPastUsageHandler == nil {
		api.UsageHostPastUsageHandler = usage.HostPastUsageHandlerFunc(func(params usage.HostPastUsageParams) middleware.Responder {
			return middleware.NotImplemented("operation usage.HostPastUsage has not yet been implemented")
		})
	}
	if api.PolicyListPolicyHandler == nil {
		api.PolicyListPolicyHandler = policy.ListPolicyHandlerFunc(func(params policy.ListPolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation policy.ListPolicy has not yet been implemented")
		})
	} else {
		api.PolicyListPolicyHandler = policy.ListPolicyHandlerFunc(handler.ListPolicyHandler)
	}
	if api.ResourceListResourceHandler == nil {
		api.ResourceListResourceHandler = resource.ListResourceHandlerFunc(func(params resource.ListResourceParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.ListResource has not yet been implemented")
		})
	} else {
		api.ResourceListResourceHandler = resource.ListResourceHandlerFunc(handler.ListResourceHandler)
	}
	if api.TemplateListTemplateHandler == nil {
		api.TemplateListTemplateHandler = template.ListTemplateHandlerFunc(func(params template.ListTemplateParams) middleware.Responder {
			return middleware.NotImplemented("operation template.ListTemplate has not yet been implemented")
		})
	}
	if api.ResourceOverviewStatsHandler == nil {
		api.ResourceOverviewStatsHandler = resource.OverviewStatsHandlerFunc(func(params resource.OverviewStatsParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.OverviewStats has not yet been implemented")
		})
	}
	if api.UserRegisterUserHandler == nil {
		api.UserRegisterUserHandler = user.RegisterUserHandlerFunc(func(params user.RegisterUserParams) middleware.Responder {
			return middleware.NotImplemented("operation user.RegisterUser has not yet been implemented")
		})
	} else {
		api.UserRegisterUserHandler = user.RegisterUserHandlerFunc(handler.RegisterUserHandler)
	}
	if api.PolicyRemovePolicyHandler == nil {
		api.PolicyRemovePolicyHandler = policy.RemovePolicyHandlerFunc(func(params policy.RemovePolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation policy.RemovePolicy has not yet been implemented")
		})
	} else {
		api.PolicyRemovePolicyHandler = policy.RemovePolicyHandlerFunc(handler.RemovePolicyHandler)
	}
	if api.ResourceResourceInfoHandler == nil {
		api.ResourceResourceInfoHandler = resource.ResourceInfoHandlerFunc(func(params resource.ResourceInfoParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.ResourceInfo has not yet been implemented")
		})
	} else {
		api.ResourceResourceInfoHandler = resource.ResourceInfoHandlerFunc(handler.ResourceInfoHandler)
	}
	if api.ResourceResourceVMsInfoHandler == nil {
		api.ResourceResourceVMsInfoHandler = resource.ResourceVMsInfoHandlerFunc(func(params resource.ResourceVMsInfoParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.ResourceVMsInfo has not yet been implemented")
		})
	} else {
		api.ResourceResourceVMsInfoHandler = resource.ResourceVMsInfoHandlerFunc(handler.ResourceVMsInfoHandler)
	}
	if api.ResourceToggleActiveHandler == nil {
		api.ResourceToggleActiveHandler = resource.ToggleActiveHandlerFunc(func(params resource.ToggleActiveParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.ToggleActive has not yet been implemented")
		})
	} else {
		api.ResourceToggleActiveHandler = resource.ToggleActiveHandlerFunc(handler.ToggleActiveHandler)
	}
	if api.ResourceUpdateHostHandler == nil {
		api.ResourceUpdateHostHandler = resource.UpdateHostHandlerFunc(func(params resource.UpdateHostParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.UpdateHost has not yet been implemented")
		})
	}
	if api.UsageUpdateHostUsageHandler == nil {
		api.UsageUpdateHostUsageHandler = usage.UpdateHostUsageHandlerFunc(func(params usage.UpdateHostUsageParams) middleware.Responder {
			return middleware.NotImplemented("operation usage.UpdateHostUsage has not yet been implemented")
		})
	} else {
		api.UsageUpdateHostUsageHandler = usage.UpdateHostUsageHandlerFunc(handler.UpdateHostUsageHandler)
	}
	if api.PolicyUpdatePolicyHandler == nil {
		api.PolicyUpdatePolicyHandler = policy.UpdatePolicyHandlerFunc(func(params policy.UpdatePolicyParams) middleware.Responder {
			return middleware.NotImplemented("operation policy.UpdatePolicy has not yet been implemented")
		})
	} else {
		api.PolicyUpdatePolicyHandler = policy.UpdatePolicyHandlerFunc(handler.UpdatePolicyHandler)
	}
	if api.UserUserDetailsHandler == nil {
		api.UserUserDetailsHandler = user.UserDetailsHandlerFunc(func(params user.UserDetailsParams) middleware.Responder {
			return middleware.NotImplemented("operation user.UserDetails has not yet been implemented")
		})
	} else {
		api.UserUserDetailsHandler = user.UserDetailsHandlerFunc(handler.UserDetailsHandler)
	}
	if api.UserUserLoginHandler == nil {
		api.UserUserLoginHandler = user.UserLoginHandlerFunc(func(params user.UserLoginParams) middleware.Responder {
			return middleware.NotImplemented("operation user.UserLogin has not yet been implemented")
		})
	} else {
		api.UserUserLoginHandler = user.UserLoginHandlerFunc(handler.UserLoginHandler)
	}
	if api.ResourceValidateResourceHandler == nil {
		api.ResourceValidateResourceHandler = resource.ValidateResourceHandlerFunc(func(params resource.ValidateResourceParams) middleware.Responder {
			return middleware.NotImplemented("operation resource.ValidateResource has not yet been implemented")
		})
	} else {
		api.ResourceValidateResourceHandler = resource.ValidateResourceHandlerFunc(handler.ValidateResourceHandler)
	}

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
	return handler
}
