// Code generated by go-swagger; DO NOT EDIT.

package vm

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"context"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ListVMHandlerFunc turns a function with the right signature into a list VM handler
type ListVMHandlerFunc func(ListVMParams) middleware.Responder

// Handle executing the request and returning a response
func (fn ListVMHandlerFunc) Handle(params ListVMParams) middleware.Responder {
	return fn(params)
}

// ListVMHandler interface for that can handle valid list VM params
type ListVMHandler interface {
	Handle(ListVMParams) middleware.Responder
}

// NewListVM creates a new http.Handler for the list VM operation
func NewListVM(ctx *middleware.Context, handler ListVMHandler) *ListVM {
	return &ListVM{Context: ctx, Handler: handler}
}

/* ListVM swagger:route GET /vapp/vm/{id} vm listVm

list VMachine

*/
type ListVM struct {
	Context *middleware.Context
	Handler ListVMHandler
}

func (o *ListVM) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewListVMParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}

// ListVMOKBodyItems0 list VM o k body items0
//
// swagger:model ListVMOKBodyItems0
type ListVMOKBodyItems0 struct {

	// IP address
	IPAddress string `json:"IPAddress,omitempty"`

	// disk
	Disk int64 `json:"disk,omitempty"`

	// external IP address
	ExternalIPAddress string `json:"externalIPAddress,omitempty"`

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// password
	Password string `json:"password,omitempty"`

	// status
	Status string `json:"status,omitempty"`

	// used money
	UsedMoney float64 `json:"usedMoney,omitempty"`

	// username
	Username string `json:"username,omitempty"`

	// vcpu
	Vcpu int64 `json:"vcpu,omitempty"`

	// vmem
	Vmem int64 `json:"vmem,omitempty"`
}

// Validate validates this list VM o k body items0
func (o *ListVMOKBodyItems0) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this list VM o k body items0 based on context it is used
func (o *ListVMOKBodyItems0) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *ListVMOKBodyItems0) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *ListVMOKBodyItems0) UnmarshalBinary(b []byte) error {
	var res ListVMOKBodyItems0
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}