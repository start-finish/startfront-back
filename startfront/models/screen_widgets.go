package models

import (
	"time"

	"gorm.io/datatypes"
)

// Screen_widgets is a sample model you can modify or remove.
type Screen_widgets struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement" json:"id"`                  // New primary key field
	ScreenID         uint64         `gorm:"column:screen_id" json:"screen_id"`                   // FK to screens
	WidgetInstanceID uint64         `gorm:"column:widget_instance_id" json:"widget_instance_id"` // FK to widget_instances
	Layout           datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"layout"`               // JSON layout data
	SortOrder        int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	Screen         *Screens          `gorm:"foreignKey:ScreenID;references:ID;constraint:OnDelete:CASCADE" json:"screen,omitempty"`
	WidgetInstance *Widget_instances `gorm:"foreignKey:WidgetInstanceID;references:ID;constraint:OnDelete:CASCADE" json:"widget_instance,omitempty"`
}

// Force GORM to use the singular table name "screen_widgets".
func (Screen_widgets) TableName() string { return "screen_widgets" }
