package models

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"gorm.io/gorm"
)

// Template schema
type Template struct {
	gorm.Model

	// compatibility
	Compatibility string `json:"compatibility,omitempty"`

	// guest o s
	GuestOS string `json:"guestOS,omitempty"`

	// memory size
	MemorySize float64 `json:"memorySize,omitempty"`

	// cpu number
	VCPUSize int64 `json:"vcpuSize,omitempty"`

	// name
	Name string `json:"name,omitempty" gorm:"unique"`

	// provisioned space
	ProvisionedSpace float64 `json:"provisionedSpace,omitempty"`

	// template type
	// Enum: [datastore upload]
	TemplateType string `json:"templateType,omitempty"`

	// vm name
	VMName string `json:"vmName,omitempty"`

	// tag
	Tag string `json:"tag,omitempty"`

	// description
	Description string `json:"description,omitempty"`

	VMTemps []VMTemp

	// resource id
	ResourceID uint `json:"resourceID,omitempty" gorm:"default:1"`
}

type VMTemp struct {
	gorm.Model

	VMName string `json:"vmName,omitempty"`

	// number of vCPU
	VCPU int `json:"vCPU,omitempty"`

	// number of memeory unit GB
	VMem int `json:"vMem,omitempty"`

	// number of disk unit GB
	Disk int `json:"disk,omitempty"`

	// the ports wanted to be exposed
	Ports string `json:"ports,omitempty"`

	// foreign key for Template
	TemplateID uint `json:"templateID,omitempty"`
}

// Validate validates this template
func (m *Template) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTemplateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var templateTypeTemplateTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["datastore","upload"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		templateTypeTemplateTypePropEnum = append(templateTypeTemplateTypePropEnum, v)
	}
}

const (

	// TemplateTemplateTypeDatastore captures enum value "datastore"
	TemplateTemplateTypeDatastore string = "datastore"

	// TemplateTemplateTypeUpload captures enum value "upload"
	TemplateTemplateTypeUpload string = "upload"
)

// prop value enum
func (m *Template) validateTemplateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, templateTypeTemplateTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Template) validateTemplateType(formats strfmt.Registry) error {

	if swag.IsZero(m.TemplateType) { // not required
		return nil
	}

	// value enum
	if err := m.validateTemplateTypeEnum("templateType", "body", m.TemplateType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Template) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Template) UnmarshalBinary(b []byte) error {
	var res Template
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
