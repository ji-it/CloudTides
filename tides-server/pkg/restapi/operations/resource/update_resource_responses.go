// Code generated by go-swagger; DO NOT EDIT.

package resource

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// UpdateResourceOKCode is the HTTP code returned for type UpdateResourceOK
const UpdateResourceOKCode int = 200

/*UpdateResourceOK returns success message

swagger:response updateResourceOK
*/
type UpdateResourceOK struct {

	/*
	  In: Body
	*/
	Payload *UpdateResourceOKBody `json:"body,omitempty"`
}

// NewUpdateResourceOK creates UpdateResourceOK with default headers values
func NewUpdateResourceOK() *UpdateResourceOK {

	return &UpdateResourceOK{}
}

// WithPayload adds the payload to the update resource o k response
func (o *UpdateResourceOK) WithPayload(payload *UpdateResourceOKBody) *UpdateResourceOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update resource o k response
func (o *UpdateResourceOK) SetPayload(payload *UpdateResourceOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateResourceOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// UpdateResourceNotFoundCode is the HTTP code returned for type UpdateResourceNotFound
const UpdateResourceNotFoundCode int = 404

/*UpdateResourceNotFound resource not found

swagger:response updateResourceNotFound
*/
type UpdateResourceNotFound struct {

	/*
	  In: Body
	*/
	Payload *UpdateResourceNotFoundBody `json:"body,omitempty"`
}

// NewUpdateResourceNotFound creates UpdateResourceNotFound with default headers values
func NewUpdateResourceNotFound() *UpdateResourceNotFound {

	return &UpdateResourceNotFound{}
}

// WithPayload adds the payload to the update resource not found response
func (o *UpdateResourceNotFound) WithPayload(payload *UpdateResourceNotFoundBody) *UpdateResourceNotFound {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the update resource not found response
func (o *UpdateResourceNotFound) SetPayload(payload *UpdateResourceNotFoundBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *UpdateResourceNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(404)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
