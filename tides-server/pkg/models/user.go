package models

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"gorm.io/gorm"
)

// User user

type User struct {
	gorm.Model

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

	// position
	Position string `json:"position,omitempty"`

	// priority
	// Enum: [Low Medium High]
	Priority string `json:"priority,omitempty"`

	// username
	Username string `gorm:"unique;not null"`
}

// Validate validates this user
func (m *User) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePriority(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var userTypePriorityPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["Low","Medium","High"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		userTypePriorityPropEnum = append(userTypePriorityPropEnum, v)
	}
}

const (

	// UserPriorityLow captures enum value "Low"
	UserPriorityLow string = "Low"

	// UserPriorityMedium captures enum value "Medium"
	UserPriorityMedium string = "Medium"

	// UserPriorityHigh captures enum value "High"
	UserPriorityHigh string = "High"
)

// prop value enum
func (m *User) validatePriorityEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, userTypePriorityPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *User) validatePriority(formats strfmt.Registry) error {

	if swag.IsZero(m.Priority) { // not required
		return nil
	}

	// value enum
	if err := m.validatePriorityEnum("priority", "body", m.Priority); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *User) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *User) UnmarshalBinary(b []byte) error {
	var res User
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
