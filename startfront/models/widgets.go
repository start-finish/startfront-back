package models

import (
	"time"

	"gorm.io/datatypes"
)

// Widgets is a sample model you can modify or remove.
type Widgets struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Key          string         `gorm:"type:text;unique;not null" json:"key"`                  // unique code name (e.g. kpi_card, chart_bar)
	IsBuiltin    bool           `gorm:"not null;default:false" json:"is_builtin"`              // true = built-in, false = custom
	Version      *string        `gorm:"type:text" json:"version,omitempty"`                    // semantic version (NULL if builtin)
	ConfigSchema datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"config_schema"` // JSON schema for config validation
	CreatedAt    time.Time      `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"not null;default:now()" json:"updated_at"`
}

// Force GORM to use the singular table name "widgets".
func (Widgets) TableName() string { return "widgets" }
