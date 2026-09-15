package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// GameSyncCheckpoint holds the schema definition for the GameSyncCheckpoint entity.
type GameSyncCheckpoint struct {
	ent.Schema
}

// Table of the GameSyncCheckpoint.
func (GameSyncCheckpoint) Table() string {
	return "game_sync_checkpoint"
}

// Fields of the GameSyncCheckpoint.
func (GameSyncCheckpoint) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("sync_scope"),
		field.String("checkpoint_value"),
		field.Int64("remote_total"),
		field.Int64("local_total"),
		field.Int64("created_count"),
		field.Int64("updated_count"),
		field.Int64("deleted_count"),
		field.Int64("failed_count"),
		field.Int64("sync_status"),
		field.String("last_error_message").Optional(),
		field.Time("last_success_at"),
		field.Time("last_sync_at"),
		field.Time("created_at"),
		field.Time("updated_at"),
		field.Int64("progress"),
	}
}

// Edges of the GameSyncCheckpoint.
func (GameSyncCheckpoint) Edges() []ent.Edge {
	return nil
}

func (GameSyncCheckpoint) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.WithComments(true),                        // 启用数据库字段注释
		schema.Comment("系统游戏同步检查点表"),                     // 设置数据库表注释
		entsql.Annotation{Table: "game_sync_checkpoint"}, // 设置数据库表名
	}
}
