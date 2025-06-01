package domain

type AuthBinding struct {
	ID                int    `db:"id" json:"id"`
	ApplicationID     int    `db:"application_id" json:"application_id"`
	Method            string `db:"method" json:"method"`
	LoginEndpoint     string `db:"login_endpoint" json:"login_endpoint"`
	CredentialsSchema string `db:"credentials_schema" json:"credentials_schema"`
	TokenPath         string `db:"token_path" json:"token_path"`
	ResponseUserPath  string `db:"response_user_path" json:"response_user_path"`
}
