package models

import (
	"time"

	"gorm.io/datatypes"
)

// Analytics_events is a sample model you can modify or remove.
type Analytics_events struct {
	ID               int64          `json:"id" db:"id" gorm:"primaryKey;autoIncrement"`
	OccurredAt       time.Time      `json:"occurred_at" db:"occurred_at" gorm:"not null;default:now()"`
	UserID           *int64         `json:"user_id,omitempty" db:"user_id"`
	ClientID         *int64         `json:"client_id,omitempty" db:"client_id"`
	EventName        string         `json:"event_name" db:"event_name" gorm:"not null"`
	ScreenID         *int64         `json:"screen_id,omitempty" db:"screen_id"`
	WidgetInstanceID *int64         `json:"widget_instance_id,omitempty" db:"widget_instance_id"`
	Meta             datatypes.JSON `json:"meta" db:"meta" gorm:"type:jsonb;not null;default:'{}'"`

	User           *Users            `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`
	Client         *Clients          `json:"client,omitempty" gorm:"foreignKey:ClientID;constraint:OnDelete:SET NULL;"`
	Screen         *Screens          `json:"screen,omitempty" gorm:"foreignKey:ScreenID;constraint:OnDelete:SET NULL;"`
	WidgetInstance *Widget_instances `json:"widget_instance,omitempty" gorm:"foreignKey:WidgetInstanceID;constraint:OnDelete:SET NULL;"`
}

// Force GORM to use the singular table name "analytics_events".
func (Analytics_events) TableName() string { return "analytics_events" }
