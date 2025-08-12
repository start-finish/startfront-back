package models

import "time"

type Home struct {
	ID         uint       `json:"id"   gorm:"column:id;primaryKey"`
	Name       string     `json:"name" gorm:"column:name"`
	CreateTime *time.Time `json:"create_time,omitempty" gorm:"column:create_time"`
}

// Force GORM to use the singular table
func (Home) TableName() string { return "home" }
