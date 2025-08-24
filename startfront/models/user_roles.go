package models

// User_roles is a sample model you can modify or remove.
type User_roles struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID int64  `gorm:"primaryKey;not null" json:"user_id"`
	RoleID int64  `gorm:"primaryKey;not null" json:"role_id"`

	User Users `gorm:"foreignKey:UserID;onDelete:CASCADE" json:"users"`
	Role Roles `gorm:"foreignKey:RoleID;onDelete:CASCADE" json:"roles"`
}

// Force GORM to use the singular table name "user_roles".
func (User_roles) TableName() string { return "user_roles" }
