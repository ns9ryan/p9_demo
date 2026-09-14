package model

import "time"

const (
	StatusNormal   int16 = 1
	StatusDisabled int16 = 2

	MenuTypeDir    int16 = 0
	MenuTypeMenu   int16 = 1
	MenuTypeButton int16 = 2

	LoginResultSuccess int16 = 1
	LoginResultFail    int16 = 2

	ActionResultSuccess int16 = 1
	ActionResultFail    int16 = 2

	RoleSuperAdmin = "super_admin"
)

type Operator struct {
	ID                     int64      `json:"id"`
	OperatorCode           string     `json:"operator_code"`
	TimezoneCode           string     `json:"timezone_code"`
	SettlementCurrencyCode string     `json:"settlement_currency_code"`
	Status                 int16      `json:"status"`
	RequiredConfigVersion  int        `json:"required_config_version"`
	CompletedConfigVersion int        `json:"completed_config_version"`
	ConfigCompletedAt      *time.Time `json:"config_completed_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
}

type User struct {
	ID                 int64      `json:"id"`
	UserCode           string     `json:"user_code"`
	OperatorID         *int64     `json:"operator_id,omitempty"`
	Username           string     `json:"username"`
	PasswordHash       string     `json:"-"`
	Salt               string     `json:"-"`
	DisplayName        string     `json:"display_name"`
	Mobile             *string    `json:"mobile,omitempty"`
	Email              *string    `json:"email,omitempty"`
	Status             int16      `json:"status"`
	IsSuperAdmin       bool       `json:"is_super_admin"`
	LastLoginAt        *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP        *string    `json:"last_login_ip,omitempty"`
	IPWhitelistEnabled int16      `json:"ip_whitelist_enabled"`
	IPWhitelist        []string   `json:"ip_whitelist"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type Role struct {
	ID          int64      `json:"id"`
	OperatorID  *int64     `json:"operator_id,omitempty"`
	RoleCode    string     `json:"role_code"`
	RoleName    string     `json:"role_name"`
	Description *string    `json:"description,omitempty"`
	Status      int16      `json:"status"`
	IsSystem    bool       `json:"is_system"`
	SortNo      int        `json:"sort_no"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type UserRole struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RoleID    int64     `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Menu struct {
	ID         int64     `json:"id"`
	ParentID   int64     `json:"parent_id"`
	MenuType   int16     `json:"menu_type"`
	Path       string    `json:"path"`
	Name       string    `json:"name"`
	Component  string    `json:"component"`
	Redirect   string    `json:"redirect"`
	Title      string    `json:"title"`
	Icon       string    `json:"icon"`
	Permission string    `json:"permission"`
	HideMenu   int16     `json:"hide_menu"`
	Sort       int       `json:"sort"`
	Disabled   int16     `json:"disabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type RoleMenu struct {
	ID        int64     `json:"id"`
	RoleID    int64     `json:"role_id"`
	MenuID    int64     `json:"menu_id"`
	CreatedAt time.Time `json:"created_at"`
}

type API struct {
	ID          int64     `json:"id"`
	Description string    `json:"description"`
	APIGroup    string    `json:"api_group"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	IsRequired  int16     `json:"is_required"`
	ServiceName string    `json:"service_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type I18n struct {
	ID        int64     `json:"id"`
	I18nCode  string    `json:"i18n_code"`
	I18nGroup string    `json:"i18n_group"`
	TransKey  string    `json:"trans_key"`
	Lang      string    `json:"lang"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type I18nLang struct {
	ID        int64     `json:"id"`
	Lang      string    `json:"lang"`
	Name      string    `json:"name"`
	Disabled  int16     `json:"disabled"`
	SortNo    int       `json:"sort_no"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginLog struct {
	ID            int64     `json:"id"`
	UserID        *int64    `json:"user_id,omitempty"`
	OperatorID    *int64    `json:"operator_id,omitempty"`
	Username      string    `json:"username"`
	LoginResult   int16     `json:"login_result"`
	FailureReason *string   `json:"failure_reason,omitempty"`
	LoginIP       string    `json:"login_ip"`
	DeviceID      *int64    `json:"device_id,omitempty"`
	UserAgent     *string   `json:"user_agent,omitempty"`
	LoginAt       time.Time `json:"login_at"`
}

type AdminActionLog struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	OperatorID     *int64    `json:"operator_id,omitempty"`
	Username       string    `json:"username"`
	RequestMethod  string    `json:"request_method"`
	RequestPath    string    `json:"request_path"`
	RequestQuery   *string   `json:"request_query,omitempty"`
	RequestBody    *string   `json:"request_body,omitempty"`
	ActionResult   int16     `json:"action_result"`
	ResponseStatus int       `json:"response_status"`
	ResponseBody   *string   `json:"response_body,omitempty"`
	DurationMS     int       `json:"duration_ms"`
	ClientIP       string    `json:"client_ip"`
	UserAgent      *string   `json:"user_agent,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ErrorLog struct {
	ID             int64     `json:"id"`
	UserID         *int64    `json:"user_id,omitempty"`
	OperatorID     *int64    `json:"operator_id,omitempty"`
	Username       string    `json:"username"`
	RequestMethod  string    `json:"request_method"`
	RequestPath    string    `json:"request_path"`
	RequestQuery   *string   `json:"request_query,omitempty"`
	RequestBody    *string   `json:"request_body,omitempty"`
	ServiceName    string    `json:"service_name"`
	ResponseStatus int       `json:"response_status"`
	ResponseBody   *string   `json:"response_body,omitempty"`
	Subject        *string   `json:"subject,omitempty"`
	Detail         *string   `json:"detail,omitempty"`
	DurationMS     int       `json:"duration_ms"`
	ClientIP       string    `json:"client_ip"`
	UserAgent      *string   `json:"user_agent,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
