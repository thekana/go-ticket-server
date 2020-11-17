package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"ticket-reservation/custom_error"
	"ticket-reservation/mock_db"
	"ticket-reservation/mock_redis"
	"ticket-reservation/utils"
)

func TestContext_GetEventDetail_Exist(t *testing.T) {
	logger, err := getLoggerForTesting()
	if err != nil {
		t.Fatalf("error get logger: %+v", err)
	}

	appContext := &Context{
		Logger:                logger,
		Config:                nil,
		RemoteAddress:         "",
		TokenSignerPrivateKey: utils.BytesToPrivateKey(privateKeyBytes),
		TokenSignerPublicKey:  utils.BytesToPublicKey(publicKeyBytes),
		DB: &mock_db.MockDB{
			StubViewEventDetail: mock_db.MockGetEventDetailForEvent1Only,
		},
		My: nil,
		RedisCache: &mock_redis.MockRedis{
			StubGetEventQuota:   mock_redis.MockRedisGetEventQuotaKeyNotSet,
			StubSetNXEventQuota: mock_redis.MockRedisSetNXEventQuotaDoNothing,
		},
	}
	token, err := appContext.createToken("tester", 1, []string{"customer"})
	_, err = appContext.GetEventDetail(ViewEventParams{
		AuthToken: token,
		EventID:   1,
	})
	assert.Nil(t, err)
}

func TestContext_GetEventDetail_NotExist(t *testing.T) {
	logger, err := getLoggerForTesting()
	if err != nil {
		t.Fatalf("error get logger: %+v", err)
	}

	appContext := &Context{
		Logger:                logger,
		Config:                nil,
		RemoteAddress:         "",
		TokenSignerPrivateKey: utils.BytesToPrivateKey(privateKeyBytes),
		TokenSignerPublicKey:  utils.BytesToPublicKey(publicKeyBytes),
		DB: &mock_db.MockDB{
			StubViewEventDetail: mock_db.MockGetEventDetailForEvent1Only,
		},
		My: nil,
		RedisCache: &mock_redis.MockRedis{
			StubGetEventQuota:   mock_redis.MockRedisGetEventQuotaKeyNotSet,
			StubSetNXEventQuota: mock_redis.MockRedisSetNXEventQuotaDoNothing,
		},
	}
	token, err := appContext.createToken("tester", 1, []string{"customer"})
	_, err = appContext.GetEventDetail(ViewEventParams{
		AuthToken: token,
		EventID:   2,
	})
	assert.Equal(t, err.(*custom_error.UserError).Code, custom_error.EventNotFound)
}
