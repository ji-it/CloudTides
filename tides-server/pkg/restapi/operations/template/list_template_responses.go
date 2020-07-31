// Code generated by go-swagger; DO NOT EDIT.

package template

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"
)

// ListTemplateOKCode is the HTTP code returned for type ListTemplateOK
const ListTemplateOKCode int = 200

/*ListTemplateOK OK

swagger:response listTemplateOK
*/
type ListTemplateOK struct {

	/*
	  In: Body
	*/
	Payload *ListTemplateOKBody `json:"body,omitempty"`
}

// NewListTemplateOK creates ListTemplateOK with default headers values
func NewListTemplateOK() *ListTemplateOK {

	return &ListTemplateOK{}
}

// WithPayload adds the payload to the list template o k response
func (o *ListTemplateOK) WithPayload(payload *ListTemplateOKBody) *ListTemplateOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the list template o k response
func (o *ListTemplateOK) SetPayload(payload *ListTemplateOKBody) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *ListTemplateOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// ListTemplateBadRequestCode is the HTTP code returned for type ListTemplateBadRequest
const ListTemplateBadRequestCode int = 400

/*ListTemplateBadRequest bad request

swagger:response listTemplateBadRequest
*/
type ListTemplateBadRequest struct {
}

// NewListTemplateBadRequest creates ListTemplateBadRequest with default headers values
func NewListTemplateBadRequest() *ListTemplateBadRequest {

	return &ListTemplateBadRequest{}
}

// WriteResponse to the client
func (o *ListTemplateBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(400)
}

// ListTemplateUnauthorizedCode is the HTTP code returned for type ListTemplateUnauthorized
const ListTemplateUnauthorizedCode int = 401

/*ListTemplateUnauthorized Unauthorized

swagger:response listTemplateUnauthorized
*/
type ListTemplateUnauthorized struct {
}

// NewListTemplateUnauthorized creates ListTemplateUnauthorized with default headers values
func NewListTemplateUnauthorized() *ListTemplateUnauthorized {

	return &ListTemplateUnauthorized{}
}

// WriteResponse to the client
func (o *ListTemplateUnauthorized) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(401)
}
