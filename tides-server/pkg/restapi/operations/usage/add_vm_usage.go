// Code generated by go-swagger; DO NOT EDIT.

package usage

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AddVMUsageHandlerFunc turns a function with the right signature into a add VM usage handler
type AddVMUsageHandlerFunc func(AddVMUsageParams) middleware.Responder

// Handle executing the request and returning a response
func (fn AddVMUsageHandlerFunc) Handle(params AddVMUsageParams) middleware.Responder {
	return fn(params)
}

// AddVMUsageHandler interface for that can handle valid add VM usage params
type AddVMUsageHandler interface {
	Handle(AddVMUsageParams) middleware.Responder
}

// NewAddVMUsage creates a new http.Handler for the add VM usage operation
func NewAddVMUsage(ctx *middleware.Context, handler AddVMUsageHandler) *AddVMUsage {
	return &AddVMUsage{Context: ctx, Handler: handler}
}

/*AddVMUsage swagger:route POST /usage/addVM usage addVmUsage

add VM usage info into database

*/
type AddVMUsage struct {
	Context *middleware.Context
	Handler AddVMUsageHandler
}

func (o *AddVMUsage) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewAddVMUsageParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// AddVMUsageBody add VM usage body
//
// swagger:model AddVMUsageBody
type AddVMUsageBody struct {

	// v ms
	VMs map[string]VMsAnon `json:"VMs,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this add VM usage body
func (o *AddVMUsageBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateVMs(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *AddVMUsageBody) validateVMs(formats strfmt.Registry) error {

	if swag.IsZero(o.VMs) { // not required
		return nil
	}

	for k := range o.VMs {

		if swag.IsZero(o.VMs[k]) { // not required
			continue
		}
		if val, ok := o.VMs[k]; ok {
			if err := val.Validate(formats); err != nil {
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *AddVMUsageBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddVMUsageBody) UnmarshalBinary(b []byte) error {
	var res AddVMUsageBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// AddVMUsageOKBody add VM usage o k body
//
// swagger:model AddVMUsageOKBody
type AddVMUsageOKBody struct {

	// message
	// Enum: [success]
	Message string `json:"message,omitempty"`
}

// Validate validates this add VM usage o k body
func (o *AddVMUsageOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var addVmUsageOKBodyTypeMessagePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["success"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		addVmUsageOKBodyTypeMessagePropEnum = append(addVmUsageOKBodyTypeMessagePropEnum, v)
	}
}

const (

	// AddVMUsageOKBodyMessageSuccess captures enum value "success"
	AddVMUsageOKBodyMessageSuccess string = "success"
)

// prop value enum
func (o *AddVMUsageOKBody) validateMessageEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, addVmUsageOKBodyTypeMessagePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *AddVMUsageOKBody) validateMessage(formats strfmt.Registry) error {

	if swag.IsZero(o.Message) { // not required
		return nil
	}

	// value enum
	if err := o.validateMessageEnum("addVmUsageOK"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *AddVMUsageOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *AddVMUsageOKBody) UnmarshalBinary(b []byte) error {
	var res AddVMUsageOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// VMsAnon v ms anon
//
// swagger:model VMsAnon
type VMsAnon struct {

	// boinc start time
	BoincStartTime string `json:"boincStartTime,omitempty"`

	// current CPU
	CurrentCPU float64 `json:"currentCPU,omitempty"`

	// current RAM
	CurrentRAM float64 `json:"currentRAM,omitempty"`

	// datacenter
	Datacenter string `json:"datacenter,omitempty"`

	// guest o s
	GuestOS string `json:"guestOS,omitempty"`

	// host name
	HostName string `json:"hostName,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// num CPU
	NumCPU int64 `json:"numCPU,omitempty"`

	// powered on
	PoweredOn bool `json:"poweredOn,omitempty"`

	// total CPU
	TotalCPU float64 `json:"totalCPU,omitempty"`

	// total RAM
	TotalRAM float64 `json:"totalRAM,omitempty"`
}

// Validate validates this v ms anon
func (o *VMsAnon) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *VMsAnon) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *VMsAnon) UnmarshalBinary(b []byte) error {
	var res VMsAnon
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
