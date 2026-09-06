package model

// SyncDiff 同步差异记录
type SyncDiff struct {
	ObjectType   string      `json:"object_type" comment:"对象类型：category/provider/channel/game/currency"`
	ObjectID     int64       `json:"object_id" comment:"P9对象ID"`
	ObjectCode   string      `json:"object_code" comment:"P9对象编码"`
	Action       string      `json:"action" comment:"同步操作：create/update/delete/noop/conflict/error"`
	RemoteID     int64       `json:"remote_id" comment:"上游对象ID"`
	RemoteCode   string      `json:"remote_code" comment:"上游对象编码"`
	Reason       string      `json:"reason" comment:"原因说明"`
	ConflictType string      `json:"conflict_type" comment:"冲突类型"`
	Details      interface{} `json:"details" comment:"详细信息"`
}

// SyncStats 同步统计
type SyncStats struct {
	RemoteTotal   int64 `json:"remote_total" comment:"上游总数"`
	LocalTotal    int64 `json:"local_total" comment:"本地总数"`
	CreateTotal   int64 `json:"create_total" comment:"新增数量"`
	UpdateTotal   int64 `json:"update_total" comment:"更新数量"`
	DeleteTotal   int64 `json:"delete_total" comment:"删除数量"`
	NoopTotal     int64 `json:"noop_total" comment:"无操作数量"`
	ConflictTotal int64 `json:"conflict_total" comment:"冲突数量"`
	ErrorTotal    int64 `json:"error_total" comment:"错误数量"`
}

// SyncPreviewResp 同步预检查响应
type SyncPreviewResp struct {
	Stats SyncStats  `json:"stats" comment:"统计信息"`
	Diffs []SyncDiff `json:"diffs" comment:"差异列表"`
}

// SyncApplyResult 同步执行结果
type SyncApplyResult struct {
	Created int64 `json:"created" comment:"成功新增数量"`
	Updated int64 `json:"updated" comment:"成功更新数量"`
	Deleted int64 `json:"deleted" comment:"成功删除数量"`
	Failed  int64 `json:"failed" comment:"失败数量"`
	Skipped int64 `json:"skipped" comment:"跳过数量"`
}

// SyncRunResp 同步执行响应
type SyncRunResp struct {
	Preview SyncPreviewResp `json:"preview" comment:"预检查结果"`
	Apply   SyncApplyResult `json:"apply" comment:"执行结果"`
}
