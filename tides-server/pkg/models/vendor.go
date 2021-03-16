package models

import "gorm.io/gorm"

type Vendor struct {
	gorm.Model

	Id int `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	URL string `json:"url,omitempty"`

	Type string `json:"type,omitempty"`

	Version string `json:"version,omitempty"`
}
