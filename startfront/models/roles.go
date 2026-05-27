package models

// Roles is a sample model you can modify or remove.
type Roles struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	RoleName    string `json:"role_name" gorm:"uniqueIndex"`
	Description string `json:"description"`
	UserCount   int    `json:"user_count"`
	Permissions string `json:"permissions"` // Serialized JSON array of strings
}

// Force GORM to use the singular table name "roles".
func (Roles) TableName() string { return "roles" }
