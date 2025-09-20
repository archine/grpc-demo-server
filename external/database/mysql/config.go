package mysql

import (
	"errors"
	"time"
)

type conf struct {
	URL        string        `mapstructure:"url"`           // Mysql地址, 例: localhost:3306
	UserName   string        `mapstructure:"username"`      // 用户名
	Password   string        `mapstructure:"password"`      // 密码
	Database   string        `mapstructure:"database"`      // 数据库名称
	LogLevel   string        `mapstructure:"log-level"`     // 日志级别, 例: debug、error
	MaxIdle    int           `mapstructure:"max_idle"`      // 最大空闲连接数
	MaxConnect int           `mapstructure:"max_connect"`   // 最大连接数
	IdleTime   time.Duration `mapstructure:"max_idle_time"` // 最大空闲时间
}

func (c *conf) verify() error {
	if c.URL == "" {
		return errors.New("mysql url is empty")
	}
	if c.UserName == "" {
		c.UserName = "root"
	}
	if c.Password == "" {
		c.Password = "root"
	}
	if c.Database == "" {
		return errors.New("mysql database is empty")
	}
	if c.LogLevel == "" {
		c.LogLevel = "debug"
	}
	if c.MaxIdle == 0 {
		c.MaxIdle = 10
	}
	if c.MaxConnect == 0 {
		c.MaxConnect = 100
	}
	if c.IdleTime == 0 {
		c.IdleTime = time.Minute * 5
	}

	return nil
}
