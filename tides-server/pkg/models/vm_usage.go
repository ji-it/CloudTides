package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"gorm.io/gorm"
)

// VMUsage VM usage

type VMUsage struct {
	gorm.Model

	// current CPU
	CurrentCPU float64 `json:"currentCPU,omitempty"`

	// current RAM
	CurrentRAM float64 `json:"currentRAM,omitempty"`

	// total CPU
	TotalCPU float64 `json:"totalCPU,omitempty"`

	// total RAM
	TotalRAM float64 `json:"totalRAM,omitempty"`

	// vm foreign key
	VmID uint

	Vm VM `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Validate validates this VM usage
func (m *VMUsage) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *VMUsage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *VMUsage) UnmarshalBinary(b []byte) error {
	var res VMUsage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
