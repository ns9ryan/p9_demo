package svc

import (
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	// GameService game.GameService
	DB *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	var db *gorm.DB
	var err error

	// 初始化数据库连接
	fmt.Printf("[RPC ServiceContext] 数据库配置 - Driver: %s, DSN: %s\n", c.Database.Driver, c.Database.DSN)

	if c.Database.DSN == "" {
		fmt.Println("[RPC ServiceContext] ❌ 数据库DSN为空")
	} else {
		fmt.Printf("[RPC ServiceContext] 开始连接数据库...\n")
		switch c.Database.Driver {
		case "mysql":
			db, err = gorm.Open(mysql.Open(c.Database.DSN), &gorm.Config{})
		case "postgres":
			db, err = gorm.Open(postgres.Open(c.Database.DSN), &gorm.Config{})
		default:
			err = fmt.Errorf("unsupported database driver: %s", c.Database.Driver)
		}

		if err != nil {
			fmt.Printf("[RPC ServiceContext] ❌ 数据库连接失败: %v\n", err)
		} else {
			fmt.Println("[RPC ServiceContext] ✓ 数据库连接成功")
			// 设置连接池
			sqlDB, _ := db.DB()
			sqlDB.SetMaxOpenConns(c.Database.MaxConns)
			sqlDB.SetMaxIdleConns(c.Database.MaxIdle)
			fmt.Printf("[RPC ServiceContext] ✓ 数据库连接池设置完成: MaxOpenConns=%d, MaxIdleConns=%d\n",
				c.Database.MaxConns, c.Database.MaxIdle)
		}
	}

	return &ServiceContext{
		Config: c,
		// GameService: game.NewGameService(),
		DB: db,
	}
}
