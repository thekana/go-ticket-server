package mock_redis

import (
	"ticket-reservation/redis_cache"
)

type MockRedis struct {
	redis_cache.Cache
	StubGetEventQuota   func(refID int) (int, error)
	StubIncEventQuota   func(refID int, val int) error
	StubDecEventQuota   func(refID int, val int) error
	StubSetEventQuota   func(refID int, val int) error
	StubSetNXEventQuota func(refID int, val int) error
	StubDelEventQuota   func(refID int) error
}

func (r *MockRedis) GetEventQuota(refID int) (int, error) {
	return r.StubGetEventQuota(refID)
}
func (r *MockRedis) DecEventQuota(refID int, val int) error {
	return r.StubDecEventQuota(refID, val)
}
func (r *MockRedis) IncEventQuota(refID int, val int) error {
	return r.StubIncEventQuota(refID, val)
}
func (r *MockRedis) SetEventQuota(refID int, val int) error {
	return r.StubSetEventQuota(refID, val)
}
func (r *MockRedis) SetNXEventQuota(refID int, val int) error {
	return r.StubSetNXEventQuota(refID, val)
}
func (r *MockRedis) DelEventQuota(refID int) error {
	return r.StubDelEventQuota(refID)
}

func MockRedisGetEventQuotaKeyNotSet(refID int) (int, error) {
	return -1, nil
}

func MockRedisSetNXEventQuotaDoNothing(refID int, val int) error {
	return nil
}
