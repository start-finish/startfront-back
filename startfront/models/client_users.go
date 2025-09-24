package models

// Client_users is a sample model you can modify or remove.
type Client_users struct {
	ID       uint  `json:"id" gorm:"primaryKey"`
	ClientId int64 `gorm:"primaryKey;not null" json:"client_id"`
	UserId   int64 `gorm:"primaryKey;not null" json:"user_id"`

	Client Clients `gorm:"foreignKey:ClientId;onDelete:CASCADE" json:"clients"`
	User   Users   `gorm:"foreignKey:UserId;onDelete:CASCADE" json:"users"`
}

// Force GORM to use the singular table name "client_users".
func (Client_users) TableName() string { return "client_users" }
