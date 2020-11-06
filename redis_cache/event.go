package redis_cache

import (
	"context"
	"github.com/pkg/errors"
	"strconv"
)

type CacheEvent interface {
	GetEventQuota(refID int) (int, error)
	IncEventQuota(refID int, val int) (int, error)
	DecEventQuota(refID int, val int) (int, error)
	SetEventQuota(refID int, val int) error
}

func (r *RedisCache) GetEventQuota(refID int) (int, error) {
	key := "event-" + strconv.Itoa(refID)
	intCmd, err := r.Redis.Get(context.Background(), key).Int()
	if err != nil {
		return -1, errors.Wrap(err, "Unable to get quota")
	}
	return intCmd, nil
}
func (r *RedisCache) DecEventQuota(refID int, val int) (int, error) {
	key := "event-" + strconv.Itoa(refID)
	currentQuota, err := r.GetEventQuota(refID)
	if err != nil {
		return -1, err
	}
	if currentQuota < val {
		return -1, errors.New("Not enough quota")
	}
	intCmd := r.Redis.DecrBy(context.Background(), key, int64(val))
	if err := intCmd.Err(); err != nil {
		return -1, errors.Wrap(err, "Unable to decrement event quota")
	}
	return int(intCmd.Val()), nil
}
func (r *RedisCache) IncEventQuota(refID int, val int) (int, error) {
	key := "event-" + strconv.Itoa(refID)
	intCmd := r.Redis.IncrBy(context.Background(), key, int64(val))
	if err := intCmd.Err(); err != nil {
		return -1, errors.Wrap(err, "Unable to increment event quota")
	}
	return int(intCmd.Val()), nil
}
func (r *RedisCache) SetEventQuota(refID int, val int) error {
	key := "event-" + strconv.Itoa(refID)
	setCmd := r.Redis.Set(context.Background(), key, val, 0)
	if err := setCmd.Err(); err != nil {
		return errors.Wrap(err, "Unable to set event quota")
	}
	return nil
}
