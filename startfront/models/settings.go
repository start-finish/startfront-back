package models

import (
	"time"

	"gorm.io/datatypes"
)

// Settings represents configuration entries with global or client scope.
type Settings struct {
	ID        int64          `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	Key       string         `json:"key" db:"key" gorm:"not null"`
	Scope     string         `json:"scope" db:"scope" gorm:"not null;default:global"` // global | client
	ClientID  *int64         `json:"client_id,omitempty" db:"client_id"`
	Value     datatypes.JSON `json:"value" db:"value" gorm:"type:jsonb;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"not null;default:now()"`

	Client *Clients `json:"client,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:CASCADE;"`
}

func (Settings) TableName() string { return "settings" }
