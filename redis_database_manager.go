package dhasar

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisDatabaseManager interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Key(ctx context.Context, key string, value any) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
}

type RedisDatabaseManagerImpl struct {
	redisClient *redis.Client
}

func (m *RedisDatabaseManagerImpl) Delete(ctx context.Context, key string) error {
	return m.redisClient.Del(ctx, key).Err()
}

func (m *RedisDatabaseManagerImpl) Key(ctx context.Context, key string, value any) (string, error) {
	keyByte, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	keyHash := md5.Sum([]byte(fmt.Sprintf("%s/%s", key, string(keyByte))))
	keySum := hex.EncodeToString(keyHash[:])

	return keySum, nil
}

func (m *RedisDatabaseManagerImpl) Get(ctx context.Context, key string) ([]byte, error) {
	valStr, err := m.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return []byte(valStr), nil
}

func (m *RedisDatabaseManagerImpl) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	valueByte, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return m.redisClient.Set(ctx, key, valueByte, expiration).Err()
}

func NewRedisDatabaseManager(redisClient *redis.Client) RedisDatabaseManager {
	return &RedisDatabaseManagerImpl{
		redisClient: redisClient,
	}
}
