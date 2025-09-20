package mysql

import (
	"fmt"

	"github.com/archine/gin-plus/v4/app"
	"github.com/archine/gin-plus/v4/component/gplog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Starter MySQL 启动器
type Starter struct{}

func NewStarter() *Starter {
	return &Starter{}
}

func (s *Starter) Order() int {
	return 0
}

func (s *Starter) OnContainerRefreshBefore(ctx app.ApplicationContext) {
	var cfg conf
	if err := ctx.GetConfigProvider().Unmarshal("mysql", &cfg); err != nil {
		gplog.Fatal(fmt.Sprintf("Connecting to MySQL failed, unable to parse configuration: %v", err))
	}
	if err := cfg.verify(); err != nil {
		gplog.Fatal(fmt.Sprintf("Connecting to MySQL failed, %v", err))
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.UserName, cfg.Password, cfg.URL, cfg.Database)

	gormConf := &gorm.Config{
		PrepareStmt: true,
	}

	if cfg.LogLevel == "error" {
		gormConf.Logger = logger.Default.LogMode(logger.Error)
	} else {
		gormConf.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}), gormConf)

	if err != nil {
		gplog.Fatal(fmt.Sprintf("Connecting to MySQL failed, %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		gplog.Fatal(fmt.Sprintf("Connecting to MySQL failed, %v", err))
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.MaxConnect)
	sqlDB.SetConnMaxLifetime(cfg.IdleTime)

	ctx.RegisterBean("mysql", db)
}
