package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"gorm.io/gorm"
)

// ResourceUsage host usage

type ResourceUsage struct {
	gorm.Model

	// current CPU
	CurrentCPU float64 `json:"currentCPU,omitempty"`

	// current RAM
	CurrentRAM float64 `json:"currentRAM,omitempty"`

	// host address
	HostAddress string `json:"hostAddress,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// percent CPU
	PercentCPU float64 `json:"percentCPU,omitempty"`

	// percent RAM
	PercentRAM float64 `json:"percentRAM,omitempty"`

	// total CPU
	TotalCPU float64 `json:"totalCPU,omitempty"`

	// total RAM
	TotalRAM float64 `json:"totalRAM,omitempty"`

	// resource foreign key
	ResourceID uint

	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Validate validates this host usage
func (m *ResourceUsage) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ResourceUsage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ResourceUsage) UnmarshalBinary(b []byte) error {
	var res ResourceUsage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
