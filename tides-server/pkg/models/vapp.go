package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"gorm.io/gorm"
)

type Vapp struct {
	gorm.Model

	// user id
	UserId uint `json:"userId,omitempty"`

	// ip address
	IPAddress string `json:"ipAddress,omitempty"`

	// name
	Name string `json:"name,omitempty" gorm:"unique"`

	// template name
	Template string `json:"template,omitempty"`

	// is destroyed
	IsDestroyed bool `json:"isDestroyed,omitempty"`

	// is powered on
	PoweredOn bool `json:"poweredOn,omitempty"`

	// resource foreign key
	ResourceID uint

	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Status string `json:"status"`

	VMs []VMachine
}

// Validate validates this VM
func (m *Vapp) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Vapp) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Vapp) UnmarshalBinary(b []byte) error {
	var res Vapp
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
