package models

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"gorm.io/gorm"
)

// Policy policy

type Policy struct {
	gorm.Model

	// account type
	// Enum: [accManager boinc]
	AccountType string `json:"accountType,omitempty"`

	// boinc password
	BoincPassword string `json:"BoincPassword,omitempty"`

	// boinc username
	BoincUsername string `json:"BoincUsername,omitempty"`

	// deploy type
	// Enum: [K8S VM]
	DeployType string `json:"deployType,omitempty"`

	// idle policy
	IdlePolicy string `json:"idlePolicy,omitempty"`

	// is destroy
	IsDestroy bool `json:"isDestroy,omitempty"`

	// name
	Name string `json:"name,omitempty" gorm:"unique"`

	// project foreign key
	ProjectID uint

	Project Project `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// template foreign key
	TemplateID uint

	Template Template `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// threshold policy
	ThresholdPolicy string `json:"thresholdPolicy,omitempty"`

	// user foreign key
	UserID uint

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// resource foreign key
	ResourceID uint
}

// Validate validates this policy
func (m *Policy) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccountType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateDeployType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var policyTypeAccountTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["accManager","boinc"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		policyTypeAccountTypePropEnum = append(policyTypeAccountTypePropEnum, v)
	}
}

const (

	// PolicyAccountTypeAccManager captures enum value "accManager"
	PolicyAccountTypeAccManager string = "accManager"

	// PolicyAccountTypeBoinc captures enum value "boinc"
	PolicyAccountTypeBoinc string = "boinc"
)

// prop value enum
func (m *Policy) validateAccountTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, policyTypeAccountTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Policy) validateAccountType(formats strfmt.Registry) error {

	if swag.IsZero(m.AccountType) { // not required
		return nil
	}

	// value enum
	if err := m.validateAccountTypeEnum("accountType", "body", m.AccountType); err != nil {
		return err
	}

	return nil
}

var policyTypeDeployTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["K8S","VM"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		policyTypeDeployTypePropEnum = append(policyTypeDeployTypePropEnum, v)
	}
}

const (

	// PolicyDeployTypeK8S captures enum value "K8S"
	PolicyDeployTypeK8S string = "K8S"

	// PolicyDeployTypeVM captures enum value "VM"
	PolicyDeployTypeVM string = "VM"
)

// prop value enum
func (m *Policy) validateDeployTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, policyTypeDeployTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Policy) validateDeployType(formats strfmt.Registry) error {

	if swag.IsZero(m.DeployType) { // not required
		return nil
	}

	// value enum
	if err := m.validateDeployTypeEnum("deployType", "body", m.DeployType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Policy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Policy) UnmarshalBinary(b []byte) error {
	var res Policy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
