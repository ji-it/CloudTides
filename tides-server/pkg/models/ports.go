package models

import (
	"gorm.io/gorm"
)

type Port struct {
	gorm.Model

	// exposed port for the vm
	Port uint `json:"port,omitempty"`

	// corresponding url for the exposed port
	URL string `json:"url,omitempty"`

	// foreign key for vmachine
	VMachineID uint `json:"vmachineID,omitempty"`
}
