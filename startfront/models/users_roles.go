package models

// Users_roles is a sample model you can modify or remove.
type Users_roles struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID int64  `gorm:"primaryKey;not null"`
	RoleID int64  `gorm:"primaryKey;not null"`
	Status string `gorm:"default:'active'"`

	User Users `gorm:"foreignKey:UserID;onDelete:CASCADE"`
	Role Roles `gorm:"foreignKey:RoleID;onDelete:CASCADE"`
}

// Force GORM to use the singular table name "users_roles".
func (Users_roles) TableName() string { return "users_roles" }
