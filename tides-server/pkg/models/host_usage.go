package models

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/jinzhu/gorm"
)

// HostUsage host usage

type HostUsage struct {
	gorm.Model

	// current CPU
	CurrentCPU float64 `json:"currentCPU,omitempty"`

	// current RAM
	CurrentRAM float64 `json:"currentRAM,omitempty"`

	// host address
	HostAddress string `json:"hostAddress,omitempty"`

	// host name
	HostName string `json:"hostName,omitempty"`

	// percent CPU
	PercentCPU float64 `json:"percentCPU,omitempty"`

	// percent RAM
	PercentRAM float64 `json:"percentRAM,omitempty"`

	// total CPU
	TotalCPU float64 `json:"totalCPU,omitempty"`

	// total RAM
	TotalRAM float64 `json:"totalRAM,omitempty"`

	// resource foreign key
	ResourceRef uint
}

// Validate validates this host usage
func (m *HostUsage) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *HostUsage) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HostUsage) UnmarshalBinary(b []byte) error {
	var res HostUsage
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
