package models

import (
	"time"

	"gorm.io/datatypes"
)

// Themes is a sample model you can modify or remove.
type Themes struct {
	ID        int64          `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Name      string         `json:"name" db:"name" gorm:"not null"`
	Scope     string         `json:"scope" db:"scope" gorm:"not null;default:global"` // global | client
	ClientID  *int64         `json:"client_id,omitempty" db:"client_id"`
	Variables datatypes.JSON `json:"variables" db:"variables" gorm:"type:jsonb;not null;default:'{}'"`
	IsActive  bool           `json:"is_active" db:"is_active" gorm:"not null;default:false"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null;default:now()"`

	Client *Clients `json:"client,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;"`
}

// Force GORM to use the singular table name "themes".
func (Themes) TableName() string { return "themes" }
