package etcd

import "time"

// conf etcd配置
type conf struct {
	Endpoints   []string      `mapstructure:"endpoints"`    // etcd地址列表
	TTL         int64         `mapstructure:"ttl"`          // 租约时间，单位秒
	DialTimeout time.Duration `mapstructure:"dial-timeout"` // 连接超时时间
}

func (c *conf) verify() {
	if len(c.Endpoints) == 0 {
		c.Endpoints = []string{"http://localhost:2379"}
	}
	if c.TTL == 0 {
		c.TTL = 10
	}
	if c.DialTimeout == 0 {
		c.DialTimeout = 3 * time.Second
	}
}
