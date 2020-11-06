package redis_cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"ticket-reservation/log"
	"time"
)

type Cache interface {
	CacheEvent
	Close() error
}

type RedisCache struct {
	logger log.Logger
	Config *Config
	Redis  *redis.Client
}

func New(config *Config, logger log.Logger) (*RedisCache, error) {
	var redisClient *redis.Client
	redisClient = redis.NewClient(&redis.Options{
		Addr:            config.RedisHost + ":" + config.RedisPort,
		Password:        config.RedisPassword,
		DB:              config.RedisDB,
		MaxRetries:      config.MaxRetries,
		MinRetryBackoff: time.Duration(config.MinRetryBackoffSeconds) * time.Second,
		MaxRetryBackoff: time.Duration(config.MaxRetryBackoffSeconds) * time.Second,
		DialTimeout:     time.Duration(config.DialTimeoutSeconds) * time.Second,
		WriteTimeout:    time.Duration(config.WriteTimeoutSeconds) * time.Second,
		PoolTimeout:     time.Duration(config.PoolTimeoutSeconds) * time.Second,
	})
	ctxRedis, cancelRedis := context.WithTimeout(context.Background(), config.RedisConnectionTimeout*time.Second)
	defer cancelRedis()

	cliCh := make(chan string)
	errCh := make(chan error)

	go func() {
		logger.Infof("Ping to redis")
		cli, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			errCh <- err
		}
		if cli == "PONG" {
			cliCh <- cli
		}
	}()

	select {
	case <-cliCh:
		logger.Infof("Redis connected")
		return &RedisCache{
			logger: logger.WithFields(log.Fields{
				"module": "redis_cache",
			}),
			Redis:  redisClient,
			Config: config,
		}, nil
	case errMsg := <-errCh:
		return nil, errors.New("Cannot connect to redis : " + errMsg.Error())
	case <-ctxRedis.Done():
		return nil, errors.New("Redis connection timeout")
	}
}

func (r *RedisCache) Close() error {
	return r.Redis.Close()
}
