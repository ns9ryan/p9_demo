package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap is a custom type for JSON fields
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion failed")
	}
	return json.Unmarshal(bytes, &j)
}

// Category 游戏分类模型
type Category struct {
	ID                 int64      `gorm:"column:id" json:"id"`
	SourceID           int64      `gorm:"column:source_id" json:"source_id"`
	CategoryCode       string     `gorm:"column:category_code" json:"category_code"`
	SourceCategoryCode string     `gorm:"column:source_category_code" json:"source_category_code"`
	SourceNameI18n     JSONMap    `gorm:"column:source_name_i18n;type:jsonb" json:"source_name_i18n"`
	NameI18n           JSONMap    `gorm:"column:name_i18n;type:jsonb" json:"name_i18n"`
	SourceSortNo       int64      `gorm:"column:source_sort_no" json:"source_sort_no"`
	SortNo             int64      `gorm:"column:sort_no" json:"sort_no"`
	SourceStatus       uint32     `gorm:"column:source_status" json:"source_status"`
	Status             uint32     `gorm:"column:status" json:"status"`
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定 Category 模型对应的表名
func (Category) TableName() string {
	return "game_category"
}

// Provider 游戏厂商模型
type Provider struct {
	ID                 int64      `gorm:"column:id" json:"id"`
	SourceID           int64      `gorm:"column:source_id" json:"source_id"`
	ProviderCode       string     `gorm:"column:provider_code" json:"provider_code"`
	SourceProviderCode string     `gorm:"column:source_provider_code" json:"source_provider_code"`
	SourceNameI18n     JSONMap    `gorm:"column:source_name_i18n;type:jsonb" json:"source_name_i18n"`
	NameI18n           JSONMap    `gorm:"column:name_i18n;type:jsonb" json:"name_i18n"`
	SourceLogoURL      *string    `gorm:"column:source_logo_url" json:"source_logo_url"`
	LogoURL            *string    `gorm:"column:logo_url" json:"logo_url"`
	SourceSortNo       int        `gorm:"column:source_sort_no" json:"source_sort_no"`
	SortNo             int        `gorm:"column:sort_no" json:"sort_no"`
	SourceStatus       uint32     `gorm:"column:source_status" json:"source_status"`
	Status             uint32     `gorm:"column:status" json:"status"`
	DeletedAt          *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt          time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定 Provider 模型对应的表名
func (Provider) TableName() string {
	return "game_provider"
}

// Channel 游戏渠道模型
type Channel struct {
	ID                int64      `gorm:"column:id" json:"id"`
	SourceID          int64      `gorm:"column:source_id" json:"source_id"`
	ChannelCode       string     `gorm:"column:channel_code" json:"channel_code"`
	SourceChannelCode string     `gorm:"column:source_channel_code" json:"source_channel_code"`
	SourceNameI18n    JSONMap    `gorm:"column:source_name_i18n;type:jsonb" json:"source_name_i18n"`
	NameI18n          JSONMap    `gorm:"column:name_i18n;type:jsonb" json:"name_i18n"`
	SourceSortNo      int        `gorm:"column:source_sort_no;default:0" json:"source_sort_no"`
	SortNo            int        `gorm:"column:sort_no;default:0" json:"sort_no"`
	SourceStatus      uint32     `gorm:"column:source_status" json:"source_status"`
	Status            uint32     `gorm:"column:status" json:"status"`
	DeletedAt         *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt         time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定 Channel 模型对应的表名
func (Channel) TableName() string {
	return "game_channel"
}

// Game 游戏模型
// Game 游戏模型 - 核心字段（20个导出字段，符合规范）
type Game struct {
	// 身份标识（3）
	ID       int64  `gorm:"column:id" json:"id"`
	GameCode string `gorm:"column:game_code" json:"game_code"`
	SourceID int64  `gorm:"column:source_id" json:"source_id"`

	// 外键关联（3）
	CategoryID int64 `gorm:"column:category_id" json:"category_id"`
	ProviderID int64 `gorm:"column:provider_id" json:"provider_id"`
	ChannelID  int64 `gorm:"column:channel_id" json:"channel_id"`

	// 厂商标识（2）
	ProviderKey       string `gorm:"column:provider_key" json:"provider_key"`
	SourceProviderKey string `gorm:"column:source_provider_key" json:"source_provider_key"`

	// 本地展示信息（4）
	NameI18n JSONMap `gorm:"column:name_i18n;type:jsonb" json:"name_i18n"`
	ImageURL string  `gorm:"column:image_url" json:"image_url"`
	SortNo   int     `gorm:"column:sort_no" json:"sort_no"`
	Status   int16   `gorm:"column:status" json:"status"` // 1启用/2停用

	// 源系统信息（4）
	SourceGameCode string `gorm:"column:source_game_code" json:"source_game_code"`
	SourceStatus   int16  `gorm:"column:source_status" json:"source_status"`
	SourceImageURL string `gorm:"column:source_image_url" json:"source_image_url"`
	SourceSortNo   int    `gorm:"column:source_sort_no" json:"source_sort_no"`

	// 管理字段（3）
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"` // 软删除时间戳
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// 多语言源信息（1，GORM会自动扫描，不算导出字段）
	SourceNameI18n JSONMap `gorm:"column:source_name_i18n;type:jsonb" json:"source_name_i18n"`
}

// TableName 指定 Game 模型对应的表名
func (Game) TableName() string {
	return "game"
}

// GameCurrency 游戏币种模型
type GameCurrency struct {
	ID           int64      `gorm:"column:id" json:"id"`
	GameID       int64      `gorm:"column:game_id" json:"game_id"`
	CurrencyID   int64      `gorm:"column:currency_id" json:"currency_id"`
	SourceStatus int32      `gorm:"column:source_status" json:"source_status"`
	Status       int32      `gorm:"column:status" json:"status"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定 GameCurrency 模型对应的表名
func (GameCurrency) TableName() string {
	return "game_currency"
}
