package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"gorm.io/gorm"
)

// Project project

type Project struct {
	gorm.Model

	// has account manager
	HasAccountManager bool `json:"hasAccountManager,omitempty"`

	// project name
	ProjectName string `json:"projectName,omitempty"`

	// url
	URL string `json:"url,omitempty" gorm:"unique"`
}

// Validate validates this project
func (m *Project) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Project) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Project) UnmarshalBinary(b []byte) error {
	var res Project
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
