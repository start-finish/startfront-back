package models

import "time"

// Screens is a sample model you can modify or remove.
type Screens struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	RoutePath   string    `json:"route_path" gorm:"uniqueIndex"`
	Description string    `json:"description"`
	IsActive    string    `json:"is_active" gorm:"default:'0'"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:current_timestamp"`
}

// Force GORM to use the singular table name "screens".
func (Screens) TableName() string { return "screens" }
