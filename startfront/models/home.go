package models

// Home is a sample model you can modify or remove.
type Home struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

// Force GORM to use the singular table name "home".
func (Home) TableName() string { return "home" }
