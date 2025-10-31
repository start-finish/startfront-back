package models

import "time"

// Widget_presets is a sample model you can modify or remove.
type Widget_presets struct {
	ID          int64     `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" db:"name" gorm:"not null"`
	Description *string   `json:"description,omitempty" db:"description"`
	Scope       string    `json:"scope" db:"scope" gorm:"default:global;not null"` // global | client
	ClientID    *int64    `json:"client_id,omitempty" db:"client_id"`
	CreatedBy   *int64    `json:"created_by,omitempty" db:"created_by"`
	CreatedAt   time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null;default:now()" json:"updated_at"`
}

// Force GORM to use the singular table name "widget_presets".
func (Widget_presets) TableName() string { return "widget_presets" }
