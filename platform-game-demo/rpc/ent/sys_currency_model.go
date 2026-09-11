package ent

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SysCurrencyModel = (*customSysCurrencyModel)(nil)

type (
	// SysCurrencyModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSysCurrencyModel.
	SysCurrencyModel interface {
		sysCurrencyModel
		withSession(session sqlx.Session) SysCurrencyModel
	}

	customSysCurrencyModel struct {
		*defaultSysCurrencyModel
	}
)

// NewSysCurrencyModel returns a model for the database table.
func NewSysCurrencyModel(conn sqlx.SqlConn) SysCurrencyModel {
	return &customSysCurrencyModel{
		defaultSysCurrencyModel: newSysCurrencyModel(conn),
	}
}

func (m *customSysCurrencyModel) withSession(session sqlx.Session) SysCurrencyModel {
	return NewSysCurrencyModel(sqlx.NewSqlConnFromSession(session))
}

func (SysCurrency) TableName() string {
	return "sys_currency"
}
