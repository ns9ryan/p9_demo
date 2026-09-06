### 1. N/A

1. route definition

- Url: /game-category/get
- Method: POST
- Request: `GameCategoryGetReq`
- Response: `GameCategoryResp`

2. request definition



```golang
type GameCategoryGetReq struct {
	ID int64 `json:"id" binding:"required" comment:"分类ID"`
}
```


3. response definition



```golang
type GameCategoryResp struct {
	ID int64 `json:"id" comment:"分类ID"`
	SourceID int64 `json:"source_id" comment:"上游分类ID"`
	CategoryCode string `json:"category_code" comment:"分类编码"`
	SourceCategoryCode string `json:"source_category_code" comment:"上游分类编码"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游分类名称（多语言JSON）"`
	NameI18n string `json:"name_i18n" comment:"分类名称（多语言JSON）"`
	SourceSortNo int32 `json:"source_sort_no,omitempty" comment:"上游排序号"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 2. N/A

1. route definition

- Url: /game-category/list
- Method: POST
- Request: `GameCategoryListReq`
- Response: `GameCategoryListResp`

2. request definition



```golang
type GameCategoryListReq struct {
	Page int64 `json:"page" comment:"页码"`
	PageSize int64 `json:"page_size" comment:"每页大小"`
	Status int16 `json:"status,omitempty" comment:"状态筛选"`
}
```


3. response definition



```golang
type GameCategoryListResp struct {
	Items []GameCategoryResp `json:"items" comment:"分类列表"`
	Total int64 `json:"total" comment:"总数"`
}
```

### 3. N/A

1. route definition

- Url: /game-category/update
- Method: POST
- Request: `GameCategoryUpdateReq`
- Response: `GameCategoryResp`

2. request definition



```golang
type GameCategoryUpdateReq struct {
	ID int64 `json:"id" binding:"required" comment:"分类ID"`
	NameI18n string `json:"name_i18n,omitempty" comment:"分类名称（多语言JSON）"`
	SortNo int32 `json:"sort_no,omitempty" comment:"排序号"`
	Status int16 `json:"status,omitempty" comment:"状态：1启用/2禁用"`
}
```


3. response definition



```golang
type GameCategoryResp struct {
	ID int64 `json:"id" comment:"分类ID"`
	SourceID int64 `json:"source_id" comment:"上游分类ID"`
	CategoryCode string `json:"category_code" comment:"分类编码"`
	SourceCategoryCode string `json:"source_category_code" comment:"上游分类编码"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游分类名称（多语言JSON）"`
	NameI18n string `json:"name_i18n" comment:"分类名称（多语言JSON）"`
	SourceSortNo int32 `json:"source_sort_no,omitempty" comment:"上游排序号"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 4. N/A

1. route definition

- Url: /game-channel/get
- Method: POST
- Request: `GameChannelGetReq`
- Response: `GameChannelResp`

2. request definition



```golang
type GameChannelGetReq struct {
	ID int64 `json:"id" binding:"required" comment:"渠道ID"`
}
```


3. response definition



```golang
type GameChannelResp struct {
	ID int64 `json:"id" comment:"渠道ID"`
	SourceID int64 `json:"source_id" comment:"上游渠道ID"`
	ChannelCode string `json:"channel_code" comment:"渠道编码"`
	SourceChannelCode string `json:"source_channel_code" comment:"上游渠道编码"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游渠道名称（多语言JSON）"`
	NameI18n string `json:"name_i18n" comment:"渠道名称（多语言JSON）"`
	SourceSortNo int32 `json:"source_sort_no" comment:"上游排序号"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 5. N/A

1. route definition

- Url: /game-channel/list
- Method: POST
- Request: `GameChannelListReq`
- Response: `GameChannelListResp`

2. request definition



```golang
type GameChannelListReq struct {
	Page int64 `json:"page" comment:"页码"`
	PageSize int64 `json:"page_size" comment:"每页大小"`
	Status int16 `json:"status,omitempty" comment:"状态筛选"`
}
```


3. response definition



```golang
type GameChannelListResp struct {
	Items []GameChannelResp `json:"items" comment:"渠道列表"`
	Total int64 `json:"total" comment:"总数"`
}
```

### 6. N/A

1. route definition

- Url: /game-channel/update
- Method: POST
- Request: `GameChannelUpdateReq`
- Response: `GameChannelResp`

2. request definition



```golang
type GameChannelUpdateReq struct {
	ID int64 `json:"id" binding:"required" comment:"渠道ID"`
	NameI18n string `json:"name_i18n,omitempty" comment:"渠道名称（多语言JSON）"`
	SortNo int32 `json:"sort_no,omitempty" comment:"排序号"`
	Status int16 `json:"status,omitempty" comment:"状态：1启用/2禁用"`
}
```


3. response definition



```golang
type GameChannelResp struct {
	ID int64 `json:"id" comment:"渠道ID"`
	SourceID int64 `json:"source_id" comment:"上游渠道ID"`
	ChannelCode string `json:"channel_code" comment:"渠道编码"`
	SourceChannelCode string `json:"source_channel_code" comment:"上游渠道编码"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游渠道名称（多语言JSON）"`
	NameI18n string `json:"name_i18n" comment:"渠道名称（多语言JSON）"`
	SourceSortNo int32 `json:"source_sort_no" comment:"上游排序号"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 7. N/A

1. route definition

- Url: /game-currency/get
- Method: POST
- Request: `GameCurrencyGetReq`
- Response: `GameCurrencyResp`

2. request definition



```golang
type GameCurrencyGetReq struct {
	ID int64 `json:"id" binding:"required" comment:"游戏货币ID"`
}
```


3. response definition



```golang
type GameCurrencyResp struct {
	ID int64 `json:"id" comment:"游戏货币ID"`
	GameID int64 `json:"game_id" comment:"游戏ID"`
	CurrencyID int64 `json:"currency_id" comment:"货币ID"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 8. N/A

1. route definition

- Url: /game-currency/list
- Method: POST
- Request: `GameCurrencyListReq`
- Response: `GameCurrencyListResp`

2. request definition



```golang
type GameCurrencyListReq struct {
	Page int64 `json:"page" comment:"页码"`
	PageSize int64 `json:"page_size" comment:"每页大小"`
	GameID int64 `json:"game_id,omitempty" comment:"游戏ID筛选"`
	Status int16 `json:"status,omitempty" comment:"状态筛选"`
}
```


3. response definition



```golang
type GameCurrencyListResp struct {
	Items []GameCurrencyResp `json:"items" comment:"游戏货币列表"`
	Total int64 `json:"total" comment:"总数"`
}
```

### 9. N/A

1. route definition

- Url: /game-currency/update
- Method: POST
- Request: `GameCurrencyUpdateReq`
- Response: `GameCurrencyResp`

2. request definition



```golang
type GameCurrencyUpdateReq struct {
	ID int64 `json:"id" binding:"required" comment:"游戏货币ID"`
	Status int16 `json:"status,omitempty" comment:"状态：1启用/2禁用"`
}
```


3. response definition



```golang
type GameCurrencyResp struct {
	ID int64 `json:"id" comment:"游戏货币ID"`
	GameID int64 `json:"game_id" comment:"游戏ID"`
	CurrencyID int64 `json:"currency_id" comment:"货币ID"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 10. N/A

1. route definition

- Url: /game-provider/get
- Method: POST
- Request: `GameProviderGetReq`
- Response: `GameProviderResp`

2. request definition



```golang
type GameProviderGetReq struct {
	ID int64 `json:"id" binding:"required" comment:"厂商ID"`
}
```


3. response definition



```golang
type GameProviderResp struct {
	ID int64 `json:"id" comment:"厂商ID"`
	SourceID int64 `json:"source_id" comment:"上游厂商ID"`
	ProviderCode string `json:"provider_code" comment:"厂商编码"`
	SourceProviderCode string `json:"source_provider_code" comment:"上游厂商编码"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游厂商名称（多语言JSON）"`
	NameI18n string `json:"name_i18n" comment:"厂商名称（多语言JSON）"`
	SourceLogoUrl string `json:"source_logo_url,omitempty" comment:"上游厂商Logo URL"`
	LogoUrl string `json:"logo_url,omitempty" comment:"厂商Logo URL"`
	SourceSortNo int32 `json:"source_sort_no,omitempty" comment:"上游排序号"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 11. N/A

1. route definition

- Url: /game-provider/list
- Method: POST
- Request: `GameProviderListReq`
- Response: `GameProviderListResp`

2. request definition



```golang
type GameProviderListReq struct {
	Page int64 `json:"page" comment:"页码"`
	PageSize int64 `json:"page_size" comment:"每页大小"`
	Status int16 `json:"status,omitempty" comment:"状态筛选"`
}
```


3. response definition



```golang
type GameProviderListResp struct {
	Items []GameProviderResp `json:"items" comment:"厂商列表"`
	Total int64 `json:"total" comment:"总数"`
}
```

### 12. N/A

1. route definition

- Url: /game-provider/update
- Method: POST
- Request: `GameProviderUpdateReq`
- Response: `GameProviderResp`

2. request definition



```golang
type GameProviderUpdateReq struct {
	ID int64 `json:"id" binding:"required" comment:"厂商ID"`
	NameI18n string `json:"name_i18n,omitempty" comment:"厂商名称（多语言JSON）"`
	LogoUrl string `json:"logo_url,omitempty" comment:"厂商Logo URL"`
	SortNo int32 `json:"sort_no,omitempty" comment:"排序号"`
	Status int16 `json:"status,omitempty" comment:"状态：1启用/2禁用"`
}
```


3. response definition



```golang
type GameProviderResp struct {
	ID int64 `json:"id" comment:"厂商ID"`
	SourceID int64 `json:"source_id" comment:"上游厂商ID"`
	ProviderCode string `json:"provider_code" comment:"厂商编码"`
	SourceProviderCode string `json:"source_provider_code" comment:"上游厂商编码"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游厂商名称（多语言JSON）"`
	NameI18n string `json:"name_i18n" comment:"厂商名称（多语言JSON）"`
	SourceLogoUrl string `json:"source_logo_url,omitempty" comment:"上游厂商Logo URL"`
	LogoUrl string `json:"logo_url,omitempty" comment:"厂商Logo URL"`
	SourceSortNo int32 `json:"source_sort_no,omitempty" comment:"上游排序号"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 13. N/A

1. route definition

- Url: /game-sync-checkpoint/get
- Method: POST
- Request: `GameSyncCheckpointGetReq`
- Response: `GameSyncCheckpointResp`

2. request definition



```golang
type GameSyncCheckpointGetReq struct {
	SyncScope string `json:"sync_scope" binding:"required" comment:"同步范围"`
}
```


3. response definition



```golang
type GameSyncCheckpointResp struct {
	ID int64 `json:"id" comment:"检查点ID"`
	SyncScope string `json:"sync_scope" comment:"同步范围"`
	CheckpointValue string `json:"checkpoint_value" comment:"检查点值"`
	RemoteTotal int64 `json:"remote_total" comment:"上游总数"`
	LocalTotal int64 `json:"local_total" comment:"本地总数"`
	CreatedCount int64 `json:"created_count" comment:"创建数量"`
	UpdatedCount int64 `json:"updated_count" comment:"更新数量"`
	DeletedCount int64 `json:"deleted_count" comment:"删除数量"`
	FailedCount int64 `json:"failed_count" comment:"失败数量"`
	SyncStatus int16 `json:"sync_status" comment:"同步状态：0未进行/1进行中"`
	LastErrorMessage string `json:"last_error_message,omitempty" comment:"最后错误信息"`
	LastSuccessAt int64 `json:"last_success_at" comment:"最后成功时间戳"`
	LastSyncAt int64 `json:"last_sync_at" comment:"最后同步时间戳"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}
```

### 14. N/A

1. route definition

- Url: /game-sync-checkpoint/list
- Method: POST
- Request: `GameSyncCheckpointListReq`
- Response: `GameSyncCheckpointListResp`

2. request definition



```golang
type GameSyncCheckpointListReq struct {
	Page int64 `json:"page" comment:"页码"`
	PageSize int64 `json:"page_size" comment:"每页大小"`
}
```


3. response definition



```golang
type GameSyncCheckpointListResp struct {
	Items []GameSyncCheckpointResp `json:"items" comment:"检查点列表"`
	Total int64 `json:"total" comment:"总数"`
}
```

### 15. N/A

1. route definition

- Url: /game/get
- Method: POST
- Request: `GameGetReq`
- Response: `GameResp`

2. request definition



```golang
type GameGetReq struct {
	ID int64 `json:"id" binding:"required" comment:"游戏ID"`
}
```


3. response definition



```golang
type GameResp struct {
	ID int64 `json:"id" comment:"游戏ID"`
	GameCode string `json:"game_code" comment:"游戏编码"`
	CategoryID int64 `json:"category_id" comment:"分类ID"`
	ProviderID int64 `json:"provider_id" comment:"厂商ID"`
	ChannelID int64 `json:"channel_id,omitempty" comment:"渠道ID"`
	ProviderKey string `json:"provider_key,omitempty" comment:"厂商Key"`
	NameI18n string `json:"name_i18n" comment:"游戏名称（多语言JSON）"`
	ImageUrl string `json:"image_url,omitempty" comment:"游戏图片URL"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SupportsEmbed bool `json:"supports_embed" comment:"是否支持嵌入"`
	SupportsRedirect bool `json:"supports_redirect" comment:"是否支持重定向"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	SourceInfo GameSourceInfo `json:"source_info" comment:"上游游戏信息"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}

type GameSourceInfo struct {
	SourceID int64 `json:"source_id" comment:"上游游戏ID"`
	SourceGameCode string `json:"source_game_code" comment:"上游游戏编码"`
	SourceProviderKey string `json:"source_provider_key,omitempty" comment:"上游厂商Key"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游游戏名称（多语言JSON）"`
	SourceImageUrl string `json:"source_image_url,omitempty" comment:"上游游戏图片URL"`
	SourceSortNo int32 `json:"source_sort_no,omitempty" comment:"上游排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
}
```

### 16. N/A

1. route definition

- Url: /game/list
- Method: POST
- Request: `GameListReq`
- Response: `GameListResp`

2. request definition



```golang
type GameListReq struct {
	Page int64 `json:"page" comment:"页码"`
	PageSize int64 `json:"page_size" comment:"每页大小"`
	CategoryID int64 `json:"category_id,omitempty" comment:"分类ID筛选"`
	ProviderID int64 `json:"provider_id,omitempty" comment:"厂商ID筛选"`
	ChannelID int64 `json:"channel_id,omitempty" comment:"渠道ID筛选"`
	Status int16 `json:"status,omitempty" comment:"状态筛选"`
}
```


3. response definition



```golang
type GameListResp struct {
	Items []GameResp `json:"items" comment:"游戏列表"`
	Total int64 `json:"total" comment:"总数"`
}
```

### 17. N/A

1. route definition

- Url: /game/update
- Method: POST
- Request: `GameUpdateReq`
- Response: `GameResp`

2. request definition



```golang
type GameUpdateReq struct {
	ID int64 `json:"id" binding:"required" comment:"游戏ID"`
	NameI18n string `json:"name_i18n,omitempty" comment:"游戏名称（多语言JSON）"`
	ImageUrl string `json:"image_url,omitempty" comment:"游戏图片URL"`
	SortNo int32 `json:"sort_no,omitempty" comment:"排序号"`
	ProviderKey string `json:"provider_key,omitempty" comment:"厂商Key"`
	SupportsEmbed bool `json:"supports_embed,omitempty" comment:"是否支持嵌入"`
	SupportsRedirect bool `json:"supports_redirect,omitempty" comment:"是否支持重定向"`
	Status int16 `json:"status,omitempty" comment:"状态：1启用/2禁用"`
}
```


3. response definition



```golang
type GameResp struct {
	ID int64 `json:"id" comment:"游戏ID"`
	GameCode string `json:"game_code" comment:"游戏编码"`
	CategoryID int64 `json:"category_id" comment:"分类ID"`
	ProviderID int64 `json:"provider_id" comment:"厂商ID"`
	ChannelID int64 `json:"channel_id,omitempty" comment:"渠道ID"`
	ProviderKey string `json:"provider_key,omitempty" comment:"厂商Key"`
	NameI18n string `json:"name_i18n" comment:"游戏名称（多语言JSON）"`
	ImageUrl string `json:"image_url,omitempty" comment:"游戏图片URL"`
	SortNo int32 `json:"sort_no" comment:"排序号"`
	SupportsEmbed bool `json:"supports_embed" comment:"是否支持嵌入"`
	SupportsRedirect bool `json:"supports_redirect" comment:"是否支持重定向"`
	Status int16 `json:"status" comment:"状态：1启用/2禁用"`
	IsDeleted int16 `json:"is_deleted" comment:"软删除：0正常/1已删除"`
	SourceInfo GameSourceInfo `json:"source_info" comment:"上游游戏信息"`
	CreatedAt int64 `json:"created_at" comment:"创建时间戳"`
	UpdatedAt int64 `json:"updated_at" comment:"更新时间戳"`
}

type GameSourceInfo struct {
	SourceID int64 `json:"source_id" comment:"上游游戏ID"`
	SourceGameCode string `json:"source_game_code" comment:"上游游戏编码"`
	SourceProviderKey string `json:"source_provider_key,omitempty" comment:"上游厂商Key"`
	SourceNameI18n string `json:"source_name_i18n" comment:"上游游戏名称（多语言JSON）"`
	SourceImageUrl string `json:"source_image_url,omitempty" comment:"上游游戏图片URL"`
	SourceSortNo int32 `json:"source_sort_no,omitempty" comment:"上游排序号"`
	SourceStatus int16 `json:"source_status" comment:"上游状态：1启用/2禁用"`
}
```

### 18. N/A

1. route definition

- Url: /ping
- Method: GET
- Request: `-`
- Response: `PingResponse`

2. request definition



3. response definition



```golang
type PingResponse struct {
	Message string `json:"message"`
}
```

### 19. N/A

1. route definition

- Url: /sync/game-vendor/categories/preview
- Method: POST
- Request: `SyncPreviewReq`
- Response: `SyncPreviewResp`

2. request definition



```golang
type SyncPreviewReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
	RemoteTotal int64 `json:"remote_total" comment:"上游总数"`
	LocalTotal int64 `json:"local_total" comment:"本地总数"`
	CreateTotal int64 `json:"create_total" comment:"新增数量"`
	UpdateTotal int64 `json:"update_total" comment:"更新数量"`
	DeleteTotal int64 `json:"delete_total" comment:"删除数量"`
	NoopTotal int64 `json:"noop_total" comment:"无操作数量"`
	ConflictTotal int64 `json:"conflict_total" comment:"冲突数量"`
	ErrorTotal int64 `json:"error_total" comment:"错误数量"`
}
```

### 20. N/A

1. route definition

- Url: /sync/game-vendor/categories/run
- Method: POST
- Request: `SyncRunReq`
- Response: `SyncRunResp`

2. request definition



```golang
type SyncRunReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	AutoApply bool `json:"auto_apply" comment:"是否自动应用"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncRunResp struct {
	Preview SyncPreviewResp `json:"preview" comment:"预检查结果"`
	Apply SyncApplyResult `json:"apply" comment:"执行结果"`
}

type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
}

type SyncApplyResult struct {
	Created int64 `json:"created" comment:"新增数量"`
	Updated int64 `json:"updated" comment:"更新数量"`
	Deleted int64 `json:"deleted" comment:"删除数量"`
	Failed int64 `json:"failed" comment:"失败数量"`
	Skipped int64 `json:"skipped" comment:"跳过数量"`
}
```

### 21. N/A

1. route definition

- Url: /sync/game-vendor/channel/preview
- Method: POST
- Request: `SyncPreviewReq`
- Response: `SyncPreviewResp`

2. request definition



```golang
type SyncPreviewReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
	RemoteTotal int64 `json:"remote_total" comment:"上游总数"`
	LocalTotal int64 `json:"local_total" comment:"本地总数"`
	CreateTotal int64 `json:"create_total" comment:"新增数量"`
	UpdateTotal int64 `json:"update_total" comment:"更新数量"`
	DeleteTotal int64 `json:"delete_total" comment:"删除数量"`
	NoopTotal int64 `json:"noop_total" comment:"无操作数量"`
	ConflictTotal int64 `json:"conflict_total" comment:"冲突数量"`
	ErrorTotal int64 `json:"error_total" comment:"错误数量"`
}
```

### 22. N/A

1. route definition

- Url: /sync/game-vendor/channel/run
- Method: POST
- Request: `SyncRunReq`
- Response: `SyncRunResp`

2. request definition



```golang
type SyncRunReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	AutoApply bool `json:"auto_apply" comment:"是否自动应用"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncRunResp struct {
	Preview SyncPreviewResp `json:"preview" comment:"预检查结果"`
	Apply SyncApplyResult `json:"apply" comment:"执行结果"`
}

type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
}

type SyncApplyResult struct {
	Created int64 `json:"created" comment:"新增数量"`
	Updated int64 `json:"updated" comment:"更新数量"`
	Deleted int64 `json:"deleted" comment:"删除数量"`
	Failed int64 `json:"failed" comment:"失败数量"`
	Skipped int64 `json:"skipped" comment:"跳过数量"`
}
```

### 23. N/A

1. route definition

- Url: /sync/game-vendor/currencies/preview
- Method: POST
- Request: `SyncPreviewReq`
- Response: `SyncPreviewResp`

2. request definition



```golang
type SyncPreviewReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
	RemoteTotal int64 `json:"remote_total" comment:"上游总数"`
	LocalTotal int64 `json:"local_total" comment:"本地总数"`
	CreateTotal int64 `json:"create_total" comment:"新增数量"`
	UpdateTotal int64 `json:"update_total" comment:"更新数量"`
	DeleteTotal int64 `json:"delete_total" comment:"删除数量"`
	NoopTotal int64 `json:"noop_total" comment:"无操作数量"`
	ConflictTotal int64 `json:"conflict_total" comment:"冲突数量"`
	ErrorTotal int64 `json:"error_total" comment:"错误数量"`
}
```

### 24. N/A

1. route definition

- Url: /sync/game-vendor/currencies/run
- Method: POST
- Request: `SyncRunReq`
- Response: `SyncRunResp`

2. request definition



```golang
type SyncRunReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	AutoApply bool `json:"auto_apply" comment:"是否自动应用"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncRunResp struct {
	Preview SyncPreviewResp `json:"preview" comment:"预检查结果"`
	Apply SyncApplyResult `json:"apply" comment:"执行结果"`
}

type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
}

type SyncApplyResult struct {
	Created int64 `json:"created" comment:"新增数量"`
	Updated int64 `json:"updated" comment:"更新数量"`
	Deleted int64 `json:"deleted" comment:"删除数量"`
	Failed int64 `json:"failed" comment:"失败数量"`
	Skipped int64 `json:"skipped" comment:"跳过数量"`
}
```

### 25. N/A

1. route definition

- Url: /sync/game-vendor/games/preview
- Method: POST
- Request: `SyncPreviewReq`
- Response: `SyncPreviewResp`

2. request definition



```golang
type SyncPreviewReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
	RemoteTotal int64 `json:"remote_total" comment:"上游总数"`
	LocalTotal int64 `json:"local_total" comment:"本地总数"`
	CreateTotal int64 `json:"create_total" comment:"新增数量"`
	UpdateTotal int64 `json:"update_total" comment:"更新数量"`
	DeleteTotal int64 `json:"delete_total" comment:"删除数量"`
	NoopTotal int64 `json:"noop_total" comment:"无操作数量"`
	ConflictTotal int64 `json:"conflict_total" comment:"冲突数量"`
	ErrorTotal int64 `json:"error_total" comment:"错误数量"`
}
```

### 26. N/A

1. route definition

- Url: /sync/game-vendor/games/run
- Method: POST
- Request: `SyncRunReq`
- Response: `SyncRunResp`

2. request definition



```golang
type SyncRunReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	AutoApply bool `json:"auto_apply" comment:"是否自动应用"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncRunResp struct {
	Preview SyncPreviewResp `json:"preview" comment:"预检查结果"`
	Apply SyncApplyResult `json:"apply" comment:"执行结果"`
}

type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
}

type SyncApplyResult struct {
	Created int64 `json:"created" comment:"新增数量"`
	Updated int64 `json:"updated" comment:"更新数量"`
	Deleted int64 `json:"deleted" comment:"删除数量"`
	Failed int64 `json:"failed" comment:"失败数量"`
	Skipped int64 `json:"skipped" comment:"跳过数量"`
}
```

### 27. N/A

1. route definition

- Url: /sync/game-vendor/vendors/preview
- Method: POST
- Request: `SyncPreviewReq`
- Response: `SyncPreviewResp`

2. request definition



```golang
type SyncPreviewReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
	RemoteTotal int64 `json:"remote_total" comment:"上游总数"`
	LocalTotal int64 `json:"local_total" comment:"本地总数"`
	CreateTotal int64 `json:"create_total" comment:"新增数量"`
	UpdateTotal int64 `json:"update_total" comment:"更新数量"`
	DeleteTotal int64 `json:"delete_total" comment:"删除数量"`
	NoopTotal int64 `json:"noop_total" comment:"无操作数量"`
	ConflictTotal int64 `json:"conflict_total" comment:"冲突数量"`
	ErrorTotal int64 `json:"error_total" comment:"错误数量"`
}
```

### 28. N/A

1. route definition

- Url: /sync/game-vendor/vendors/run
- Method: POST
- Request: `SyncRunReq`
- Response: `SyncRunResp`

2. request definition



```golang
type SyncRunReq struct {
	ObjectType string `json:"object_type" binding:"required" comment:"对象类型"`
	SelectedIds []int64 `json:"selected_ids" comment:"选中的ID列表"`
	SelectedCodes []string `json:"selected_codes" comment:"选中的编码列表"`
	AutoApply bool `json:"auto_apply" comment:"是否自动应用"`
	StrictConflict bool `json:"strict_conflict" comment:"是否严格处理冲突"`
}
```


3. response definition



```golang
type SyncRunResp struct {
	Preview SyncPreviewResp `json:"preview" comment:"预检查结果"`
	Apply SyncApplyResult `json:"apply" comment:"执行结果"`
}

type SyncPreviewResp struct {
	Stats SyncStats `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

type SyncStats struct {
}

type SyncApplyResult struct {
	Created int64 `json:"created" comment:"新增数量"`
	Updated int64 `json:"updated" comment:"更新数量"`
	Deleted int64 `json:"deleted" comment:"删除数量"`
	Failed int64 `json:"failed" comment:"失败数量"`
	Skipped int64 `json:"skipped" comment:"跳过数量"`
}
```

