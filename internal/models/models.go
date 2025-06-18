// models/user.go
package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    gorm.Model
    Email     string    `gorm:"unique;not null"`
    Username  string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Role      string    `gorm:"not null;check:role IN ('admin','client')"`
    DeviceID  *string
    Location  *string
    IPAddress *string
    CreatedAt time.Time
    UpdatedAt time.Time
    Sessions  []Session
}

type Session struct {
    gorm.Model
    UserID           uint
    JWTAccessToken   string
    RefreshToken     string
    Expiration       time.Time
    LoginTime        time.Time
    ActiveTokens     int
    Revocations      int
}

type Application struct {
    gorm.Model
    Name        string
    Description string
    OwnerID     uint
    IsTemplate  bool
    Collaborators []Collaborator
}

type Collaborator struct {
    gorm.Model
    ApplicationID uint
    UserID        uint
    Role          string `gorm:"not null;check:role IN ('owner','editor','viewer')"`
}

type AppConnection struct {
    gorm.Model
    ApplicationID uint
    EndpointName  string
    Method        string
    URL           string
    Headers       map[string]interface{} `gorm:"type:jsonb"`
    Auth          map[string]interface{} `gorm:"type:jsonb"`
    Params        map[string]interface{} `gorm:"type:jsonb"`
    ResponseMapping map[string]interface{} `gorm:"type:jsonb"`
}

type Screen struct {
    gorm.Model
    ApplicationID uint
    Title         string
    Route         string
    ScreenType    string
    DeviceType    string
    Description   *string
    Status        string `gorm:"not null;check:status IN ('published','draft')"`
}

type Widget struct {
    gorm.Model
    ScreenID  uint
    ParentID  *uint
    Properties map[string]interface{} `gorm:"type:jsonb"`
    Actions    map[string]interface{} `gorm:"type:jsonb"`
}

type AdminDashboard struct {
    gorm.Model
    Metrics        map[string]interface{} `gorm:"type:jsonb"`
    SystemSettings map[string]interface{} `gorm:"type:jsonb"`
}

type ClientDashboard struct {
    gorm.Model
    UserID         uint
    ProjectsMetadata map[string]interface{} `gorm:"type:jsonb"`
    ThemeConfig      map[string]interface{} `gorm:"type:jsonb"`
    APIManager       map[string]interface{} `gorm:"type:jsonb"`
    Versioning       map[string]interface{} `gorm:"type:jsonb"`
    Feedback         map[string]interface{} `gorm:"type:jsonb"`
    ProfileSettings  map[string]interface{} `gorm:"type:jsonb"`
}
