package models

import (
	"time"
)

type AuditFields struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	CreatedBy  *int64
	ModifiedBy *int64
	SoftDelete bool
}

type Users struct {
	Id                   *int64    `json:"id"`
	Email                string    `json:"email" binding:"required"`
	Username             string    `json:"username" binding:"required"`
	PasswordHash         string    `json:"password_hash"`
	Role                 string    `json:"role"`
	FullName             string    `json:"full_name"`
	Phone                string    `json:"phone"`
	Location             string    `json:"location"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	MemberSince          time.Time `json:"member_since"`
	NotificationSettings JSONMap   `json:"notification_settings"`
	PrivacySettings      JSONMap   `json:"privacy_settings"`
	CreatedBy            *int64    `json:"created_by"`
	ModifiedBy           *int64    `json:"modified_by"`
	SoftDelete           bool      `json:"soft_delete"`
}

type UserSessions struct {
	Id             *int64    `json:"id"`
	UserId         *int64    `json:"user_id"`
	DeviceId       string    `json:"device_id"`
	Location       string    `json:"location"`
	IpAddress      string    `json:"ip_address"`
	JwtToken       string    `json:"jwt_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpirationTime time.Time `json:"expiration_time"`
	CreatedAt      time.Time `json:"created_at"`
	LastUsedAt     time.Time `json:"last_used_at"`
	RevokedAt      time.Time `json:"revoked_at"`
}

type Applications struct {
	Id          *int64    `json:"id"`
	OwnerUserId *int64    `json:"owner_user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	IsTemplate  bool      `json:"is_template"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   *int64    `json:"created_by"`
	ModifiedBy  *int64    `json:"modified_by"`
	SoftDelete  bool      `json:"soft_delete"`
}

type ApplicationCollaborators struct {
	Id            *int64    `json:"id"`
	ApplicationId *int64    `json:"application_id"`
	UserId        *int64    `json:"user_id"`
	RoleInApp     string    `json:"role_in_app"`
	CreatedAt     time.Time `json:"created_at"`
}

type AppConnections struct {
	Id              *int64    `json:"id"`
	ApplicationId   *int64    `json:"application_id"`
	EndpointName    string    `json:"endpoint_name"`
	Method          string    `json:"method"`
	Url             string    `json:"url"`
	Headers         JSONMap   `json:"headers"`
	AuthConfig      JSONMap   `json:"auth_config"`
	Params          JSONMap   `json:"params"`
	ResponseMapping JSONMap   `json:"response_mapping"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Screens struct {
	Id            *int64    `json:"id"`
	ApplicationId *int64    `json:"application_id"`
	Title         string    `json:"title"`
	Route         string    `json:"route"`
	ScreenType    string    `json:"screen_type"`
	DeviceType    string    `json:"device_type"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     *int64    `json:"created_by"`
	ModifiedBy    *int64    `json:"modified_by"`
	SoftDelete    bool      `json:"soft_delete"`
}

type Widgets struct {
	Id             *int64    `json:"id"`
	ScreenId       *int64    `json:"screen_id"`
	ParentWidgetId *int64    `json:"parent_widget_id"`
	WidgetType     string    `json:"widget_type"`
	Props          JSONMap   `json:"props"`
	Actions        JSONMap   `json:"actions"`
	OrderIndex     *int64    `json:"order_index"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      *int64    `json:"created_by"`
	ModifiedBy     *int64    `json:"modified_by"`
	SoftDelete     bool      `json:"soft_delete"`
}

type Themes struct {
	Id             *int64    `json:"id"`
	ApplicationId  *int64    `json:"application_id"`
	ColorPrimary   string    `json:"color_primary"`
	ColorSecondary string    `json:"color_secondary"`
	ColorSuccess   string    `json:"color_success"`
	ColorError     string    `json:"color_error"`
	FontFamily     string    `json:"font_family"`
	FontSize       JSONMap   `json:"font_size"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Assets struct {
	Id               *int64    `json:"id"`
	ApplicationId    *int64    `json:"application_id"`
	FileName         string    `json:"file_name"`
	FileType         string    `json:"file_type"`
	Url              string    `json:"url"`
	UploadedByUserId *int64    `json:"uploaded_by_user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type ReleaseVersions struct {
	Id            *int64    `json:"id"`
	ApplicationId *int64    `json:"application_id"`
	VersionNumber string    `json:"version_number"`
	Description   JSONMap   `json:"description"`
	ReleaseDate   time.Time `json:"release_date"`
	CreatedAt     time.Time `json:"created_at"`
}

type PlatformAnalytics struct {
	Id             *int64    `json:"id"`
	TotalUsers     *int64    `json:"total_users"`
	ActiveSessions *int64    `json:"active_sessions"`
	PageViews      *int64    `json:"page_views"`
	GrowthRate     string    `json:"growth_rate"`
	CreatedAt      time.Time `json:"created_at"`
}

type PlatformScreens struct {
	Id          *int64    `json:"id"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	ScreenType  string    `json:"screen_type"`
	DeviceType  string    `json:"device_type"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NavigationMenus struct {
	Id           *int64    `json:"id"`
	ParentMenuId *int64    `json:"parent_menu_id"`
	Title        string    `json:"title"`
	Route        string    `json:"route"`
	Icon         string    `json:"icon"`
	OrderIndex   *int64    `json:"order_index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type WidgetPresets struct {
	Id          *int64    `json:"id"`
	WidgetType  string    `json:"widget_type"`
	Props       JSONMap   `json:"props"`
	Functions   JSONMap   `json:"functions"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Roles struct {
	Id          *int64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RolePermissions struct {
	Id             *int64 `json:"id"`
	RoleId         *int64 `json:"role_id"`
	PermissionName string `json:"permission_name"`
}

type Feedback struct {
	Id           *int64    `json:"id"`
	UserId       *int64    `json:"user_id"`
	ProjectId    *int64    `json:"project_id"`
	Stars        *int64    `json:"stars"`
	FeedbackType string    `json:"feedback_type"`
	Comments     string    `json:"comments"`
	CreatedAt    time.Time `json:"created_at"`
}

type SystemSettings struct {
	Id        *int64    `json:"id"`
	Key       string    `json:"key"`
	Value     JSONMap   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfileSettings struct {
	Id                   *int64    `json:"id"`
	UserId               *int64    `json:"user_id"`
	FullName             string    `json:"full_name"`
	PhoneNumber          string    `json:"phone_number"`
	Location             string    `json:"location"`
	NotificationSettings JSONMap   `json:"notification_settings"`
	PrivacySettings      JSONMap   `json:"privacy_settings"`
	CreatedAt            time.Time `json:"created_at"`
}

type ProjectVersionHistory struct {
	Id            *int64    `json:"id"`
	ApplicationId *int64    `json:"application_id"`
	VersionNumber string    `json:"version_number"`
	Changes       JSONMap   `json:"changes"`
	CreatedAt     time.Time `json:"created_at"`
}
