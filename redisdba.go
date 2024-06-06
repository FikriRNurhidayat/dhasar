package dhasar

import (
	"context"

	"github.com/fikrirnurhidayat/x/logger"
	"github.com/redis/go-redis/v9"
)

type RedisDatabaseAdapterOption struct {
	Network  string
	Addr     string
	Username string
	Password string
	DB       int
}

type RedisDatabaseAdapter struct {
	db     *redis.Client
	logger logger.Logger
}

func (r *RedisDatabaseAdapter) Close() error {
	if err := r.db.Close(); err != nil {
		r.logger.Error("redis/CLOSE", logger.String("error", err.Error()))
		return err
	}

	r.logger.Debug("redis/CLOSE", logger.String("status", "OK!"))

	return nil
}

func (r *RedisDatabaseAdapter) Connect(opt *RedisDatabaseAdapterOption) (*redis.Client, error) {
	ctx := context.Background()

	db := redis.NewClient(&redis.Options{
		Network:  opt.Network,
		Addr:     opt.Addr,
		Username: opt.Username,
		Password: opt.Password,
		DB:       opt.DB,
	})

	if _, err := db.Ping(ctx).Result(); err != nil {
		r.logger.Error("redis/CONNECT", logger.String("error", err.Error()))
		return nil, err
	}

	r.logger.Debug("redis/CONNECT", logger.String("status", "OK!"))

	return db, nil
}

func NewRedisDatabaseAdapter(logger logger.Logger) Adapter[*RedisDatabaseAdapterOption, *redis.Client] {
	return &RedisDatabaseAdapter{
		logger: logger,
	}
}
