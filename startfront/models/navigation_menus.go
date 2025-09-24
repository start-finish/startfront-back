package models

import (
	"time"

	"gorm.io/gorm"
)

// Navigation_menus is a sample model you can modify or remove.
type Navigation_menus struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Name      string    `json:"name" gorm:"not null"`
	Scope     string    `json:"scope" gorm:"not null;default:'global'"`
	ClientID  string    `json:"client_id" gorm:"default:null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;default:current_timestamp"`
}

// Define the check constraint in GORM using the `scope` and `client_id` logic
func (Navigation_menus) TableName() string {
	return "navigation_menus"
}

// GORM model hooks (optional) for auto-updating timestamps
func (menu *Navigation_menus) BeforeSave(tx *gorm.DB) (err error) {
	if menu.CreatedAt.IsZero() {
		menu.CreatedAt = time.Now()
	}
	menu.UpdatedAt = time.Now()
	return nil
}
