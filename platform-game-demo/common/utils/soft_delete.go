package utils

import (
	"database/sql"
	"time"
)

// IsDeleted 检查是否已被软删除（DeletedAt 不为 nil）
// 返回 1 表示已删除，0 表示未删除
func IsDeleted(deletedAt *time.Time) int32 {
	if deletedAt != nil {
		return 1
	}
	return 0
}

func IsDel(deletedAt time.Time) int32 {
	if !deletedAt.IsZero() {
		return 1
	}
	return 0
}

// IsDeletedByNullTime 检查是否已被软删除（使用 sql.NullTime）
func IsDeletedByNullTime(deletedAt sql.NullTime) int32 {
	if deletedAt.Valid {
		return 1
	}
	return 0
}
