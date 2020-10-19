package models

import (
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"gorm.io/gorm"
)

// VM VM

type VM struct {
	gorm.Model

	// boinc time
	// Format: date-time
	BoincTime time.Time `json:"boincTime,omitempty"`

	// guest o s
	GuestOS string `json:"guestOS,omitempty"`

	// ip address
	IPAddress string `json:"ipAddress,omitempty" gorm:"unique"`

	// is destroyed
	IsDestroyed bool `json:"isDestroyed,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// num CPU
	NumCPU int64 `json:"numCPU,omitempty"`

	// powered on
	PoweredOn bool `json:"poweredOn,omitempty"`

	// resource foreign key
	ResourceID uint

	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Validate validates this VM
func (m *VM) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *VM) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VM) UnmarshalBinary(b []byte) error {
	var res VM
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
