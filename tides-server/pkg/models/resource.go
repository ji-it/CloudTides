package models

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"gorm.io/gorm"
)

// Resource resource

type Resource struct {
	gorm.Model

	// cluster
	Cluster string `json:"cluster,omitempty"`

	// datacenter
	Datacenter string `json:"datacenter,omitempty"`

	// host address
	HostAddress string `json:"hostAddress,omitempty"`

	// is active
	IsActive bool `json:"isActive,omitempty"`

	// is resource pool
	IsResourcePool bool `json:"isResourcePool,omitempty"`

	// job completed
	JobCompleted int64 `json:"jobCompleted,omitempty"`

	// monitored
	Monitored bool `json:"monitored,omitempty"`

	// name
	Name string `json:"name,omitempty" gorm:"unique"`

	// password
	Password string `json:"password,omitempty"`

	// platform type
	// Enum: [vsphere kvm hyper-v]
	PlatformType string `json:"platformType,omitempty"`

	// policy foreign key
	Policy Policy `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// status
	// Enum: [idle normal busy unknown]
	Status string `json:"status,omitempty"`

	// total jobs
	TotalJobs int64 `json:"totalJobs,omitempty"`

	// total v ms
	TotalVMs int64 `json:"totalVMs,omitempty"`

	// user foreign key
	UserID uint

	User User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this resource
func (m *Resource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validatePlatformType(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var resourceTypePlatformTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["vsphere","kvm","hyper-v"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		resourceTypePlatformTypePropEnum = append(resourceTypePlatformTypePropEnum, v)
	}
}

const (

	// ResourcePlatformTypeVsphere captures enum value "vsphere"
	ResourcePlatformTypeVsphere string = "vsphere"

	// ResourcePlatformTypeKvm captures enum value "kvm"
	ResourcePlatformTypeKvm string = "kvm"

	// ResourcePlatformTypeHyperv captures enum value "hyper-v"
	ResourcePlatformTypeHyperv string = "hyper-v"
)

// prop value enum
func (m *Resource) validatePlatformTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, resourceTypePlatformTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Resource) validatePlatformType(formats strfmt.Registry) error {

	if swag.IsZero(m.PlatformType) { // not required
		return nil
	}

	// value enum
	if err := m.validatePlatformTypeEnum("platformType", "body", m.PlatformType); err != nil {
		return err
	}

	return nil
}

var resourceTypeStatusPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["idle","normal","busy","unknown"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		resourceTypeStatusPropEnum = append(resourceTypeStatusPropEnum, v)
	}
}

const (

	// ResourceStatusIdle captures enum value "idle"
	ResourceStatusIdle string = "idle"

	// ResourceStatusNormal captures enum value "normal"
	ResourceStatusNormal string = "normal"

	// ResourceStatusBusy captures enum value "busy"
	ResourceStatusBusy string = "busy"

	// ResourceStatusUnknown captures enum value "unknown"
	ResourceStatusUnknown string = "unknown"
)

// prop value enum
func (m *Resource) validateStatusEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, resourceTypeStatusPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *Resource) validateStatus(formats strfmt.Registry) error {

	if swag.IsZero(m.Status) { // not required
		return nil
	}

	// value enum
	if err := m.validateStatusEnum("status", "body", m.Status); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *Resource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Resource) UnmarshalBinary(b []byte) error {
	var res Resource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
