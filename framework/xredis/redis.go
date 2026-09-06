// Package xredis
// @author: fengyi
// @date: 2024/10/12
// @note:
package xredis

import (
	"context"
	"fmt"
	"general-agent/config"
	"general-agent/extension/logz"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

var Cli *Client

type Client struct {
	*redis.Client
}

func NewRedisClient(conf *config.Config) *Client {
	option := conf.Redis

	var redisCli *redis.Client

	if option.Mode == "sentinel" {
		logz.InfoNoCtx("init redis client...", "mode", "sentinel", logz.Any("master", option.Master), logz.String("nodes", fmt.Sprintf("%v", option.Nodes)))
		redisCli = redis.NewFailoverClient(&redis.FailoverOptions{
			MasterName:    conf.Redis.Master,
			SentinelAddrs: strings.Split(conf.Redis.Nodes, ","),
			Password:      option.Password,
			DB:            option.DB,
			PoolSize:      option.PoolSize,
		})
	} else {
		logz.InfoNoCtx("init redis client...", "mode", "single", "host", option.Host, "port", option.Port)
		redisCli = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", option.Host, option.Port),
			Password: option.Password,
			DB:       option.DB,
			PoolSize: option.PoolSize,
		})
	}

	if err := redisCli.Ping(context.Background()).Err(); err != nil {
		logz.ErrorNoCtx("redis connect error", logz.Err(err))
	} else {
		logz.InfoNoCtx("redis connected.", logz.String("addr", redisCli.Options().Addr))
	}
	Cli = &Client{Client: redisCli}
	return Cli
}

const clientLockKey = "agentcc"

// TryLock 获取锁的函数
func (c *Client) TryLock(lockKey string, expiration time.Duration) (bool, error) {
	// 尝试使用 SETNX 获取锁
	success, err := c.SetNX(context.Background(), lockKey, clientLockKey, expiration).Result()
	if err != nil {
		logz.ErrorNoCtx("try lock error", "key", lockKey, logz.Err(err))
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	return success, nil
}

// Locked 查看锁的状态
func (c *Client) Locked(lockKey string) error {
	_, err := c.Get(context.Background(), lockKey).Result()
	return err
}

// UnLock 释放锁的函数

func (c *Client) UnLock(lockKey string) error {
	// Lua 脚本，确保只有拥有锁的客户端才能释放锁
	script := `
    if redis.call("GET", KEYS[1]) == ARGV[1] then
        return redis.call("DEL", KEYS[1])
    else
        return 0
    end
    `
	_, err := c.Eval(context.Background(), script, []string{lockKey}, clientLockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return err
}

var Module = fx.Module("redis", fx.Provide(NewRedisClient))
