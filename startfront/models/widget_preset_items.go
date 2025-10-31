package models

import (
	"time"

	"gorm.io/datatypes"
)

// Widget_preset_items is a sample model you can modify or remove.
type Widget_preset_items struct {
	ID            int64          `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	PresetID      int64          `json:"preset_id" db:"preset_id" gorm:"not null;index"`
	WidgetID      int64          `json:"widget_id" db:"widget_id" gorm:"not null"`
	DefaultConfig datatypes.JSON `json:"default_config" db:"default_config" gorm:"type:jsonb;not null;default:'{}'"`
	SortOrder     int            `json:"sort_order" db:"sort_order" gorm:"not null;default:0"`
	CreatedAt     time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null;default:now()" json:"updated_at"`

	// Relations
	Preset *Widget_presets `json:"preset,omitempty" gorm:"foreignKey:PresetID;constraint:OnDelete:CASCADE;"`
	Widget *Widgets        `json:"widget,omitempty" gorm:"foreignKey:WidgetID;constraint:OnDelete:RESTRICT;"`
}

// Force GORM to use the singular table name "widget_preset_items".
func (Widget_preset_items) TableName() string { return "widget_preset_items" }
