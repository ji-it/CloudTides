// Code generated by go-swagger; DO NOT EDIT.

package policy

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

// RemovePolicyHandlerFunc turns a function with the right signature into a remove policy handler
type RemovePolicyHandlerFunc func(RemovePolicyParams) middleware.Responder

// Handle executing the request and returning a response
func (fn RemovePolicyHandlerFunc) Handle(params RemovePolicyParams) middleware.Responder {
	return fn(params)
}

// RemovePolicyHandler interface for that can handle valid remove policy params
type RemovePolicyHandler interface {
	Handle(RemovePolicyParams) middleware.Responder
}

// NewRemovePolicy creates a new http.Handler for the remove policy operation
func NewRemovePolicy(ctx *middleware.Context, handler RemovePolicyHandler) *RemovePolicy {
	return &RemovePolicy{Context: ctx, Handler: handler}
}

/*RemovePolicy swagger:route DELETE /policy/remove policy removePolicy

remove a policy

*/
type RemovePolicy struct {
	Context *middleware.Context
	Handler RemovePolicyHandler
}

func (o *RemovePolicy) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewRemovePolicyParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// RemovePolicyBody remove policy body
//
// swagger:model RemovePolicyBody
type RemovePolicyBody struct {

	// id
	ID int64 `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this remove policy body
func (o *RemovePolicyBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *RemovePolicyBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RemovePolicyBody) UnmarshalBinary(b []byte) error {
	var res RemovePolicyBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// RemovePolicyNotFoundBody remove policy not found body
//
// swagger:model RemovePolicyNotFoundBody
type RemovePolicyNotFoundBody struct {

	// message
	// Enum: [no matching objects]
	Message string `json:"message,omitempty"`
}

// Validate validates this remove policy not found body
func (o *RemovePolicyNotFoundBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var removePolicyNotFoundBodyTypeMessagePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["no matching objects"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		removePolicyNotFoundBodyTypeMessagePropEnum = append(removePolicyNotFoundBodyTypeMessagePropEnum, v)
	}
}

const (

	// RemovePolicyNotFoundBodyMessageNoMatchingObjects captures enum value "no matching objects"
	RemovePolicyNotFoundBodyMessageNoMatchingObjects string = "no matching objects"
)

// prop value enum
func (o *RemovePolicyNotFoundBody) validateMessageEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, removePolicyNotFoundBodyTypeMessagePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *RemovePolicyNotFoundBody) validateMessage(formats strfmt.Registry) error {

	if swag.IsZero(o.Message) { // not required
		return nil
	}

	// value enum
	if err := o.validateMessageEnum("removePolicyNotFound"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RemovePolicyNotFoundBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RemovePolicyNotFoundBody) UnmarshalBinary(b []byte) error {
	var res RemovePolicyNotFoundBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// RemovePolicyOKBody remove policy o k body
//
// swagger:model RemovePolicyOKBody
type RemovePolicyOKBody struct {

	// message
	// Enum: [success]
	Message string `json:"message,omitempty"`
}

// Validate validates this remove policy o k body
func (o *RemovePolicyOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateMessage(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var removePolicyOKBodyTypeMessagePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["success"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		removePolicyOKBodyTypeMessagePropEnum = append(removePolicyOKBodyTypeMessagePropEnum, v)
	}
}

const (

	// RemovePolicyOKBodyMessageSuccess captures enum value "success"
	RemovePolicyOKBodyMessageSuccess string = "success"
)

// prop value enum
func (o *RemovePolicyOKBody) validateMessageEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, removePolicyOKBodyTypeMessagePropEnum); err != nil {
		return err
	}
	return nil
}

func (o *RemovePolicyOKBody) validateMessage(formats strfmt.Registry) error {

	if swag.IsZero(o.Message) { // not required
		return nil
	}

	// value enum
	if err := o.validateMessageEnum("removePolicyOK"+"."+"message", "body", o.Message); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RemovePolicyOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RemovePolicyOKBody) UnmarshalBinary(b []byte) error {
	var res RemovePolicyOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
