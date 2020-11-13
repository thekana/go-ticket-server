package redis_cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	customError "ticket-reservation/custom_error"
)

// All returned errors are properly labeled
type CacheEvent interface {
	GetEventQuota(refID int) (int, error)
	IncEventQuota(refID int, val int) error
	DecEventQuota(refID int, val int) error
	SetEventQuota(refID int, val int) error
	SetNXEventQuota(refID int, val int) error
}

func (r *RedisCache) GetEventQuota(refID int) (int, error) {
	key := strconv.Itoa(refID)
	intCmd, err := r.Redis.Get(context.Background(), key).Int()
	if err != nil {
		if err == redis.Nil {
			return -1, nil
		}
		return -2, &customError.InternalError{
			Code:    customError.RedisError,
			Message: "Unable to get quota from cache",
		}
	}
	return intCmd, nil
}
func (r *RedisCache) DecEventQuota(refID int, val int) error {
	key := strconv.Itoa(refID)
	lua := "if tonumber(redis.call('get',KEYS[1])) >= tonumber(ARGV[1]) then redis.call('decrby',KEYS[1], ARGV[1]) return 'true' else return 'false' end"
	cmd, err := r.Redis.Eval(context.Background(), lua, []string{key}, val).Result()
	if err != nil {
		return err
	}
	if cmd.(string) == "false" {
		return &customError.UserError{
			Code:           customError.InsufficientQuota,
			Message:        "Not enough quota",
			HTTPStatusCode: http.StatusBadRequest,
		}
	}
	return nil
}
func (r *RedisCache) IncEventQuota(refID int, val int) error {
	key := strconv.Itoa(refID)
	intCmd := r.Redis.IncrBy(context.Background(), key, int64(val))
	if err := intCmd.Err(); err != nil {
		return errors.Wrap(err, "Unable to increment event quota")
	}
	return nil
}
func (r *RedisCache) SetEventQuota(refID int, val int) error {
	key := strconv.Itoa(refID)
	setCmd := r.Redis.Set(context.Background(), key, val, 0)
	if err := setCmd.Err(); err != nil {
		return &customError.InternalError{
			Code:    customError.RedisError,
			Message: "Unable to set quota in cache",
		}
	}
	return nil
}

func (r *RedisCache) SetNXEventQuota(refID int, val int) error {
	key := strconv.Itoa(refID)
	setNXCmd := r.Redis.SetNX(context.Background(), key, val, 0)
	if err := setNXCmd.Err(); err != nil {
		return &customError.InternalError{
			Code:    customError.RedisError,
			Message: "Unable to set quota in cache",
		}
	}
	if set := setNXCmd.Val(); !set {
		return &customError.InternalError{
			Code:    customError.RedisError,
			Message: "Already set",
		}
	}
	return nil
}
