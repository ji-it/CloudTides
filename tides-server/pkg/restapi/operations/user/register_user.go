// Code generated by go-swagger; DO NOT EDIT.

package user

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

// RegisterUserHandlerFunc turns a function with the right signature into a register user handler
type RegisterUserHandlerFunc func(RegisterUserParams) middleware.Responder

// Handle executing the request and returning a response
func (fn RegisterUserHandlerFunc) Handle(params RegisterUserParams) middleware.Responder {
	return fn(params)
}

// RegisterUserHandler interface for that can handle valid register user params
type RegisterUserHandler interface {
	Handle(RegisterUserParams) middleware.Responder
}

// NewRegisterUser creates a new http.Handler for the register user operation
func NewRegisterUser(ctx *middleware.Context, handler RegisterUserHandler) *RegisterUser {
	return &RegisterUser{Context: ctx, Handler: handler}
}

/*RegisterUser swagger:route POST /users/register user registerUser

user registration

*/
type RegisterUser struct {
	Context *middleware.Context
	Handler RegisterUserHandler
}

func (o *RegisterUser) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewRegisterUserParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}

// RegisterUserBadRequestBody register user bad request body
//
// swagger:model RegisterUserBadRequestBody
type RegisterUserBadRequestBody struct {

	// message
	Message string `json:"message,omitempty"`
}

// Validate validates this register user bad request body
func (o *RegisterUserBadRequestBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *RegisterUserBadRequestBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterUserBadRequestBody) UnmarshalBinary(b []byte) error {
	var res RegisterUserBadRequestBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// RegisterUserBody register user body
//
// swagger:model RegisterUserBody
type RegisterUserBody struct {

	// city
	City string `json:"city,omitempty"`

	// company name
	CompanyName string `json:"companyName,omitempty"`

	// country
	Country string `json:"country,omitempty"`

	// email
	Email string `json:"email,omitempty"`

	// first name
	FirstName string `json:"firstName,omitempty"`

	// last name
	LastName string `json:"lastName,omitempty"`

	// password
	Password string `json:"password,omitempty"`

	// phone
	Phone string `json:"phone,omitempty"`

	// position
	Position string `json:"position,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this register user body
func (o *RegisterUserBody) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (o *RegisterUserBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterUserBody) UnmarshalBinary(b []byte) error {
	var res RegisterUserBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// RegisterUserOKBody register user o k body
//
// swagger:model RegisterUserOKBody
type RegisterUserOKBody struct {

	// user info
	UserInfo *RegisterUserOKBodyUserInfo `json:"userInfo,omitempty"`
}

// Validate validates this register user o k body
func (o *RegisterUserOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateUserInfo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *RegisterUserOKBody) validateUserInfo(formats strfmt.Registry) error {

	if swag.IsZero(o.UserInfo) { // not required
		return nil
	}

	if o.UserInfo != nil {
		if err := o.UserInfo.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("registerUserOK" + "." + "userInfo")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RegisterUserOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterUserOKBody) UnmarshalBinary(b []byte) error {
	var res RegisterUserOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}

// RegisterUserOKBodyUserInfo register user o k body user info
//
// swagger:model RegisterUserOKBodyUserInfo
type RegisterUserOKBodyUserInfo struct {

	// city
	City string `json:"city,omitempty"`

	// company name
	CompanyName string `json:"companyName,omitempty"`

	// country
	Country string `json:"country,omitempty"`

	// email
	Email string `json:"email,omitempty"`

	// first name
	FirstName string `json:"firstName,omitempty"`

	// last name
	LastName string `json:"lastName,omitempty"`

	// password
	Password string `json:"password,omitempty"`

	// phone
	Phone string `json:"phone,omitempty"`

	// position
	Position string `json:"position,omitempty"`

	// priority
	// Enum: [Low Medium High]
	Priority string `json:"priority,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this register user o k body user info
func (o *RegisterUserOKBodyUserInfo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validatePriority(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var registerUserOKBodyUserInfoTypePriorityPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Low","Medium","High"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		registerUserOKBodyUserInfoTypePriorityPropEnum = append(registerUserOKBodyUserInfoTypePriorityPropEnum, v)
	}
}

const (

	// RegisterUserOKBodyUserInfoPriorityLow captures enum value "Low"
	RegisterUserOKBodyUserInfoPriorityLow string = "Low"

	// RegisterUserOKBodyUserInfoPriorityMedium captures enum value "Medium"
	RegisterUserOKBodyUserInfoPriorityMedium string = "Medium"

	// RegisterUserOKBodyUserInfoPriorityHigh captures enum value "High"
	RegisterUserOKBodyUserInfoPriorityHigh string = "High"
)

// prop value enum
func (o *RegisterUserOKBodyUserInfo) validatePriorityEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, registerUserOKBodyUserInfoTypePriorityPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (o *RegisterUserOKBodyUserInfo) validatePriority(formats strfmt.Registry) error {

	if swag.IsZero(o.Priority) { // not required
		return nil
	}

	// value enum
	if err := o.validatePriorityEnum("registerUserOK"+"."+"userInfo"+"."+"priority", "body", o.Priority); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (o *RegisterUserOKBodyUserInfo) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *RegisterUserOKBodyUserInfo) UnmarshalBinary(b []byte) error {
	var res RegisterUserOKBodyUserInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
