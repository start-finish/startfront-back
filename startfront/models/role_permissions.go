package models

// Role_permissions is a sample model you can modify or remove.
type Role_permissions struct {
	ID           uint  `json:"id" gorm:"primaryKey"`
	RoleID       int64 `gorm:"primaryKey;not null" json:"role_id"`
	PermissionID int64 `gorm:"primaryKey;not null" json:"permission_id"`

	Role       Roles       `gorm:"foreignKey:RoleID;onDelete:CASCADE" json:"roles"`
	Permission Permissions `gorm:"foreignKey:RoleID;onDelete:CASCADE" json:"permissions"`
}

// Force GORM to use the singular table name "role_permissions".
func (Role_permissions) TableName() string { return "role_permissions" }
