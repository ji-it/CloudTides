package models

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
	"gorm.io/gorm"
)

// Vsphere schema
type Vsphere struct {
	gorm.Model

	// cluster
	Cluster string `json:"cluster,omitempty"`

	// is resource pool
	IsResourcePool bool `json:"isResourcePool,omitempty"`

	ResourceID uint

	// owner resource
	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Vcd schema
type Vcd struct {
	gorm.Model

	// allocation model
	AllocationModel string `json:"allocationModel,omitempty"`

	// organization
	Organization string `json:"organization,omitempty"`

	ResourceID uint

	// owner resource
	Resource Resource `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// Resource schema
type Resource struct {
	gorm.Model

	// type of the resource, fixed or dynamic
	Type string `json:"Type,omitempty"`

	// activated, controlled by admin, indicate whether the resource is qualified for contribution
	Activated bool `json:"Activated,omitempty"`

	// datacenter
	Datacenter string `json:"datacenter,omitempty"`

	// host address
	HostAddress string `json:"hostAddress,omitempty"`

	// is active, controlled by user, indicate whether user is willing to contribute the resource
	IsActive bool `json:"isActive,omitempty"`

	// job completed
	JobCompleted int64 `json:"jobCompleted,omitempty"`

	// monitored
	Monitored bool `json:"monitored,omitempty"`

	// name
	Name string `json:"name,omitempty" gorm:"unique"`

	// password
	Password string `json:"password,omitempty"`

	// platform type
	PlatformType string `json:"platformType,omitempty"`

	// policy foreign key
	PolicyID *uint

	Policy Policy `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// setup status
	SetupStatus string `json:"setupStatus,omitempty"`

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

	Catalog string `json:"catalog,omitempty"`

	Network string `json:"network,omitempty"`
}

// Validate validates this resource
func (m *Resource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateStatus(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
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

	// ResourcePlatformTypeVsphere = "vsphere"
	ResourcePlatformTypeVsphere string = "vsphere"

	// ResourcePlatformTypeVcd = "vcd"
	ResourcePlatformTypeVcd string = "vcd"
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
