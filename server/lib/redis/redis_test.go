package redis

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	redismock "bitbucket.bri.co.id/scm/bricams-addons/transaction-status/server/lib/redis/mock"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	ctx        context.Context
	mockServer *miniredis.Miniredis
}

func (s *RedisTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.mockServer = redismock.MockRedis()
}

func TestInitRedis(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}

func (s *RedisTestSuite) TestRedis_NewRedis() {
	type expectation struct {
		out *Redis
		err error
	}

	tests := map[string]struct {
		addr     string
		expected expectation
	}{
		"Success": {
			addr: s.mockServer.Addr(),
			expected: expectation{
				out: &Redis{},
				err: nil,
			},
		},
		"Failed": {
			addr: "",
			expected: expectation{
				out: nil,
				err: errors.New("address cannot be empty or null"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			redis, redisErr := NewRedis(tt.addr, "", "", 1, 72*time.Hour)

			if redis == nil {
				if redisErr.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, redis)
				}
			}
		})
	}
}

func (s *RedisTestSuite) TestRedis_GetClient() {
	type expectation struct {
		out *Redis
	}

	tests := map[string]struct {
		expected expectation
	}{
		"Success": {
			expected: expectation{
				out: &Redis{},
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			redis := &Redis{
				client: &redis.Client{},
			}

			redisClient := redis.GetClient()

			if redisClient == nil && tt.expected.out != nil {
				t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, redis)
			}
		})
	}
}

func (s *RedisTestSuite) TestRedis_Get() {
	type expectation struct {
		out string
		err error
	}

	type mockData struct {
		key   string
		value string
	}

	tests := map[string]struct {
		addr     string
		keyRedis string
		data     *mockData
		expected expectation
	}{
		"Success": {
			addr:     s.mockServer.Addr(),
			keyRedis: "data_found",
			data: &mockData{
				key:   "data_found",
				value: "found",
			},
			expected: expectation{
				out: "found",
				err: nil,
			},
		},
		"Failed": {
			addr:     s.mockServer.Addr(),
			keyRedis: "",
			data:     nil,
			expected: expectation{
				out: "",
				err: fmt.Errorf("failed to get key: %s", redis.Nil),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			redis, _ := NewRedis(tt.addr, "", "", 1, 72*time.Hour)
			if tt.data != nil {
				redis.client.Set(s.ctx, tt.data.key, tt.data.value, 30*time.Second)
			}

			result, getErr := redis.Get(s.ctx, tt.keyRedis)

			if getErr != nil {
				if getErr.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, result)
				}
			} else {
				if result != tt.expected.out {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, result)
				}
			}
		})
	}
}

func (s *RedisTestSuite) TestRedis_Set() {
	type expectation struct {
		value string
		err   error
	}

	tests := map[string]struct {
		addr     string
		key      string
		value    interface{}
		expected expectation
	}{
		"Success": {
			addr:  s.mockServer.Addr(),
			key:   "data_valid",
			value: "valid",
			expected: expectation{
				value: "valid",
				err:   nil,
			},
		},
		"Failed": {
			addr:  s.mockServer.Addr(),
			key:   "data_invalid",
			value: make(chan int),
			expected: expectation{
				value: "",
				err:   errors.New("failed to set key: redis: can't marshal chan int (implement encoding.BinaryMarshaler)"),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			redis, _ := NewRedis(tt.addr, "", "", 1, 72*time.Hour)

			getErr := redis.Set(s.ctx, tt.key, tt.value, 30*time.Second)

			if getErr != nil {
				if getErr.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.err.Error(), getErr.Error())
				}
			} else {
				result, _ := redis.Get(s.ctx, tt.key)
				if result != tt.expected.value {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.value, result)
				}
			}
		})
	}
}

func (s *RedisTestSuite) TestRedis_Del() {
	type expectation struct {
		out int64
		err error
	}

	type mockData struct {
		key   string
		value string
	}

	tests := map[string]struct {
		addr     string
		keyRedis string
		data     *mockData
		expected expectation
	}{
		"Success": {
			addr:     s.mockServer.Addr(),
			keyRedis: "data_found",
			data: &mockData{
				key:   "data_found",
				value: "found",
			},
			expected: expectation{
				out: 1,
				err: nil,
			},
		},
		"Failed": {
			addr:     s.mockServer.Addr(),
			keyRedis: "data_not_found",
			data:     nil,
			expected: expectation{
				out: 0,
				err: fmt.Errorf("failed to del key: %s", redis.Nil),
			},
		},
	}

	for scenario, tt := range tests {
		s.T().Run(scenario, func(t *testing.T) {
			redis, _ := NewRedis(tt.addr, "", "", 1, 72*time.Hour)
			if tt.data != nil {
				redis.client.Set(s.ctx, tt.data.key, tt.data.value, 30*time.Second)
			}

			result, getErr := redis.Del(s.ctx, tt.keyRedis)

			if getErr != nil {
				if getErr.Error() != tt.expected.err.Error() {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, result)
				}
			} else {
				if result != tt.expected.out {
					t.Errorf("Out -> \nWant: %v\nGot : %v", tt.expected.out, result)
				}
			}
		})
	}
}
