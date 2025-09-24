package models

import (
	"time"

	"gorm.io/datatypes"
)

// Navigation_items is a sample model you can modify or remove.
type Navigation_items struct {
	ID         uint           `json:"id" gorm:"primaryKey;column:id"`
	MenuID     uint           `json:"menu_id" gorm:"not null;column:menu_id;index:idx_navigation_items_menu_parent"`
	ParentID   *uint          `json:"parent_id,omitempty" gorm:"column:parent_id;index:idx_navigation_items_menu_parent"`
	Label      string         `json:"label" gorm:"not null;column:label"`
	Path       *string        `json:"path,omitempty" gorm:"column:path"`
	ScreenID   *uint          `json:"screen_id,omitempty" gorm:"column:screen_id"`
	Icon       *string        `json:"icon,omitempty" gorm:"column:icon"`
	SortOrder  int            `json:"sort_order" gorm:"not null;default:0;column:sort_order;index:idx_navigation_items_menu_parent"`
	Properties datatypes.JSON `json:"properties" gorm:"type:jsonb;default:'{}';column:properties"`
	CreatedAt  time.Time      `json:"created_at" gorm:"not null;default:current_timestamp"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"not null;default:current_timestamp"`

	Menu     *Navigation_menus  `json:"menu" gorm:"foreignKey:MenuID;references:ID"`
	Parent   *Navigation_items  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []Navigation_items `json:"children" gorm:"foreignKey:ParentID"`
}

// Force GORM to use the singular table name "navigation_items".
func (Navigation_items) TableName() string { return "navigation_items" }
