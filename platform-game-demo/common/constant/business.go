package constant

// 业务常量定义

// 游戏状态
const (
	GameStatusEnabled  = 1
	GameStatusDisabled = 2
)

// 同步操作类型
const (
	SyncActionCreate   = "create"
	SyncActionUpdate   = "update"
	SyncActionDelete   = "delete"
	SyncActionNoop     = "noop"
	SyncActionConflict = "conflict"
	SyncActionError    = "error"
)

// 同步阶段
const (
	SyncPhasePreview = "preview"
	SyncPhaseRun     = "run"
)

// 对象类型
const (
	ObjectTypeCategory = "category"
	ObjectTypeProvider = "provider"
	ObjectTypeChannel  = "channel"
	ObjectTypeGame     = "game"
	ObjectTypeCurrency = "currency"
)

// 数据来源
const (
	DataSourceRemote = "remote"
	DataSourceLocal  = "local"
)

// 同步冲突类型
const (
	ConflictTypeCodeDuplicate  = "code_duplicate"
	ConflictTypeIdDuplicate    = "id_duplicate"
	ConflictTypeDataMismatch   = "data_mismatch"
	ConflictTypeStatusMismatch = "status_mismatch"
)
