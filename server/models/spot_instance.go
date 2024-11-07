package models

import (
	"time"

	"gorm.io/gorm"
)

type SpotInstance struct {
	gorm.Model
	Id        int `gorm:"uniqueIndex"`
	UserId    int
	Name      string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
