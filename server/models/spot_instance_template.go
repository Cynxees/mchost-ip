package models

import (
	"time"

	"gorm.io/gorm"
)

type Ip struct {
	gorm.Model
	Id             int `gorm:"uniqueIndex"`
	AllocationId string
	OwnerId         int
	SpotInstanceTemplateId     *int
	InstanceId     *string
	
	Name         string
	Type string
	Region string
	Address string

	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
