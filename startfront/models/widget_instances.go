package models

import (
	"time"

	"gorm.io/datatypes"
)

// Widget_instances represents widget instances
type Widget_instances struct {
	ID       int64   `gorm:"primaryKey;autoIncrement" json:"id"`
	WidgetID int64   `gorm:"not null;index" json:"widget_id"` // FK to widgets.id
	Title    *string `gorm:"type:text" json:"title,omitempty"`

	Scope    string `gorm:"type:text;not null;default:global" json:"scope"`
	ClientID *int64 `gorm:"index" json:"client_id,omitempty"` // FK to clients.id
	ScreenID *int64 `gorm:"index" json:"screen_id,omitempty"` // FK to screens.id

	Config datatypes.JSON `gorm:"type:jsonb;not null;default:'{}'" json:"config"`

	SortOrder int  `gorm:"not null;default:0" json:"sort_order"`
	IsActive  bool `gorm:"not null;default:true" json:"is_active"`

	// CreatedBy *int64    `gorm:"index" json:"created_by,omitempty"` // FK to users.id
	CreatedAt time.Time `gorm:"not null;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null;default:now()" json:"updated_at"`

	// ==== Relationships ====
	Widget *Widgets `gorm:"foreignKey:WidgetID;references:ID" json:"widget,omitempty"`
	Client *Clients `gorm:"foreignKey:ClientID;references:ID" json:"client,omitempty"`
	Screen *Screens `gorm:"foreignKey:ScreenID;references:ID" json:"screen,omitempty"`
	// User   *Users   `gorm:"foreignKey:CreatedBy;references:ID" json:"created_by_user,omitempty"`
}

func (Widget_instances) TableName() string { return "widget_instances" }
