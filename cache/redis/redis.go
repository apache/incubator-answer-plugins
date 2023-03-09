package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/answerdev/answer/plugin"
	"github.com/answerdev/plugins/cache/redis/i18n"
	"github.com/go-redis/redis/v8"
)

type Cache struct {
	Config      *CacheConfig
	RedisClient *redis.Client
}

type CacheConfig struct {
	Endpoint string `json:"endpoint"`
}

func init() {
	plugin.Register(&Cache{
		Config: &CacheConfig{},
	})
}

func (c *Cache) Info() plugin.Info {
	return plugin.Info{
		Name:        plugin.MakeTranslator(i18n.InfoName),
		SlugName:    "redis_cache",
		Description: plugin.MakeTranslator(i18n.InfoDescription),
		Author:      "answerdev",
		Version:     "0.0.1",
		Link:        "https://github.com/answerdev/plugins/tree/main/cache/redis",
	}
}

func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	if c.RedisClient == nil {
		return "", fmt.Errorf("redis is not configed yet")
	}
	resp := c.RedisClient.Get(ctx, key)
	val, err := resp.Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (c *Cache) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	if c.RedisClient == nil {
		return fmt.Errorf("redis is not configed yet")
	}
	return c.RedisClient.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) GetInt64(ctx context.Context, key string) (int64, error) {
	if c.RedisClient == nil {
		return 0, fmt.Errorf("redis is not configed yet")
	}
	resp := c.RedisClient.Get(ctx, key)
	val, err := resp.Int64()
	if err == redis.Nil {
		return val, nil
	}
	return val, err
}

func (c *Cache) SetInt64(ctx context.Context, key string, value int64, ttl time.Duration) error {
	if c.RedisClient == nil {
		return fmt.Errorf("redis is not configed yet")
	}
	return c.RedisClient.Set(ctx, key, value, ttl).Err()
}

func (c *Cache) Del(ctx context.Context, key string) error {
	if c.RedisClient == nil {
		return fmt.Errorf("redis is not configed yet")
	}
	return c.RedisClient.Del(ctx, key).Err()
}

func (c *Cache) Flush(ctx context.Context) error {
	if c.RedisClient == nil {
		return fmt.Errorf("redis is not configed yet")
	}
	return c.RedisClient.FlushDB(ctx).Err()
}

func (c *Cache) ConfigFields() []plugin.ConfigField {
	return []plugin.ConfigField{
		{
			Name:        "endpoint",
			Type:        plugin.ConfigTypeInput,
			Title:       plugin.MakeTranslator(i18n.ConfigEndpointTitle),
			Description: plugin.MakeTranslator(i18n.ConfigEndpointDescription),
			Required:    true,
			UIOptions: plugin.ConfigFieldUIOptions{
				InputType: plugin.InputTypeText,
			},
			Value: c.Config.Endpoint,
		},
	}
}

func (c *Cache) ConfigReceiver(config []byte) error {
	conf := &CacheConfig{}
	_ = json.Unmarshal(config, conf)
	c.Config = conf

	c.RedisClient = redis.NewClient(&redis.Options{
		Addr: conf.Endpoint,
	})
	return nil
}
