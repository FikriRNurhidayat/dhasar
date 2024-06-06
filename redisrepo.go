package dhasar

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fikrirnurhidayat/x/logger"
)

var (
	GET_ACTION   = "GET"
	EXIST_ACTION = "EXIST"
	LIST_ACTION  = "LIST"
	SIZE_ACTION  = "SIZE"
)

type Decoder[Entity any, EntityJSON any] func(value EntityJSON) (Entity, error)
type Encoder[Entity any, EntityJSON any] func(value Entity) (EntityJSON, error)

type RedisRepository[Entity any, Specification any, EntityJSON any, FallbackRepository Repository[Entity, Specification]] struct {
	Resource           string
	RDBM               RedisDatabaseManager
	Expiration         time.Duration
	Logger             logger.Logger
	FallbackRepository FallbackRepository
	Decode             Decoder[Entity, EntityJSON]
	Encode             Encoder[Entity, EntityJSON]
	NoEntities         []Entity
	NoEntity           Entity
	Key                string
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) GetKey(ctx context.Context, action string, value any) (string, error) {
	return r.RDBM.Key(ctx, fmt.Sprintf("%s.%s", r.Key, action), value)
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) Delete(ctx context.Context, specs ...Specification) error {
	if err := r.FallbackRepository.Delete(ctx, specs...); err != nil {
		return err
	}

	key, err := r.GetKey(ctx, GET_ACTION, specs)
	if err != nil {
		return nil
	}

	if err := r.RDBM.Delete(ctx, key); err != nil {
		return nil
	}

	return nil
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) Each(ctx context.Context, args ListArgs[Specification]) (Iterator[Entity], error) {
	return r.FallbackRepository.Each(ctx, args)
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) Exist(ctx context.Context, specs ...Specification) (bool, error) {
	key, err := r.GetKey(ctx, EXIST_ACTION, specs)
	if err != nil {
		return r.FallbackRepository.Exist(ctx, specs...)
	}

	if valByte, err := r.RDBM.Get(ctx, key); err != nil {
		return r.FallbackRepository.Exist(ctx, specs...)
	} else if valByte == nil {
		r.Logger.Debug("redis.repository/EXIST", logger.String("cache_key", key), logger.Any("cache_hit", false))
		exist, err := r.FallbackRepository.Exist(ctx, specs...)
		if err != nil {
			return exist, err
		}

		if err := r.RDBM.Set(ctx, key, exist, r.Expiration); err != nil {
			return exist, err
		}

		return exist, nil
	} else {
		var val bool

		if err := json.Unmarshal(valByte, &val); err != nil {
			r.Logger.Debug("redis.repository/EXIST", logger.String("cache_key", key), logger.Any("cache_hit", err.Error()))
			exist, err := r.FallbackRepository.Exist(ctx, specs...)
			if err != nil {
				return exist, err
			}

			if err := r.RDBM.Set(ctx, key, exist, r.Expiration); err != nil {
				return exist, err
			}

			return exist, nil
		}

		r.Logger.Debug("redis.repository/EXIST", logger.String("cache_key", key), logger.Any("cache_hit", true))

		return val, nil
	}
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) Get(ctx context.Context, specs ...Specification) (Entity, error) {
	key, err := r.GetKey(ctx, GET_ACTION, specs)
	if err != nil {
		return r.FallbackRepository.Get(ctx, specs...)
	}

	if valByte, err := r.RDBM.Get(ctx, key); err != nil {
		return r.FallbackRepository.Get(ctx, specs...)
	} else if valByte == nil {
		r.Logger.Debug("redis.repository/GET", logger.String("cache_key", key), logger.Any("cache_hit", false))
		entity, err := r.FallbackRepository.Get(ctx, specs...)
		if err != nil {
			return entity, err
		}

		entityJSON, err := r.Encode(entity)
		if err != nil {
			return entity, err
		}

		if err := r.RDBM.Set(ctx, key, entityJSON, r.Expiration); err != nil {
			return entity, err
		}

		return entity, nil
	} else {
		var val EntityJSON

		if err := json.Unmarshal(valByte, &val); err != nil {
			r.Logger.Debug("redis.repository/GET", logger.String("cache_key", key), logger.Any("cache_hit", err.Error()))
			entity, err := r.FallbackRepository.Get(ctx, specs...)
			if err != nil {
				return entity, err
			}

			entityJSON, err := r.Encode(entity)
			if err != nil {
				return entity, err
			}

			if err := r.RDBM.Set(ctx, key, entityJSON, r.Expiration); err != nil {
				return entity, err
			}

			return entity, nil
		}

		r.Logger.Debug("redis.repository/GET", logger.String("cache_key", key), logger.Any("cache_hit", true))

		return r.Decode(val)
	}
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) List(ctx context.Context, args ListArgs[Specification]) ([]Entity, error) {
	key, err := r.GetKey(ctx, LIST_ACTION, args)
	if err != nil {
		return r.FallbackRepository.List(ctx, args)
	}

	if valByte, err := r.RDBM.Get(ctx, key); err != nil {
		return r.FallbackRepository.List(ctx, args)
	} else if valByte == nil {
		r.Logger.Debug("redis.repository/LIST", logger.String("cache_key", key), logger.Any("cache_hit", false))
		entities, err := r.FallbackRepository.List(ctx, args)
		if err != nil {
			return entities, err
		}

		entitiesJSON := []EntityJSON{}

		for _, entity := range entities {
			entityJSON, err := r.Encode(entity)
			if err != nil {
				r.Logger.Debug("redis.repository/LIST", logger.String("cache_key", key), logger.Any("cache_hit", err.Error()))
				return entities, nil
			}

			entitiesJSON = append(entitiesJSON, entityJSON)
		}

		if err := r.RDBM.Set(ctx, key, entitiesJSON, r.Expiration); err != nil {
			r.Logger.Debug("redis.repository/LIST", logger.String("cache_key", key), logger.Any("cache_hit", err.Error()))
			return entities, nil
		}

		return entities, nil
	} else {
		var entitiesJSON []EntityJSON

		if err := json.Unmarshal(valByte, &entitiesJSON); err != nil {
			r.Logger.Debug("redis.repository/LIST", logger.String("cache_key", key), logger.Any("cache_hit", err.Error()))
			return r.FallbackRepository.List(ctx, args)
		}

		entities := []Entity{}

		for _, entityJSON := range entitiesJSON {
			entity, err := r.Decode(entityJSON)
			if err != nil {
				r.Logger.Debug("redis.repository/LIST", logger.String("cache_key", key), logger.Any("cache_hit", err.Error()))
				return r.FallbackRepository.List(ctx, args)
			}

			entities = append(entities, entity)
		}

		r.Logger.Debug("redis.repository/LIST", logger.String("cache_key", key), logger.Any("cache_hit", true))

		return entities, nil
	}
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) Save(ctx context.Context, entity Entity) error {
	return r.FallbackRepository.Save(ctx, entity)
}

func (r *RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]) Size(ctx context.Context, specs ...Specification) (uint32, error) {
	key, err := r.GetKey(ctx, SIZE_ACTION, specs)
	if err != nil {
		return r.FallbackRepository.Size(ctx, specs...)
	}

	if valByte, err := r.RDBM.Get(ctx, key); err != nil {
		return r.FallbackRepository.Size(ctx, specs...)
	} else if valByte == nil {
		r.Logger.Debug("redis.repository/SIZE", logger.String("cache_key", key), logger.Any("cache_hit", false))
		exist, err := r.FallbackRepository.Size(ctx, specs...)
		if err != nil {
			return exist, err
		}

		if err := r.RDBM.Set(ctx, key, exist, r.Expiration); err != nil {
			return exist, err
		}

		return exist, nil
	} else {
		var val uint32

		if err := json.Unmarshal(valByte, &val); err != nil {
			r.Logger.Debug("redis.repository/SIZE", logger.String("cache_key", key), logger.Any("cache_hit", err.Error()))
			exist, err := r.FallbackRepository.Size(ctx, specs...)
			if err != nil {
				return exist, err
			}

			if err := r.RDBM.Set(ctx, key, exist, r.Expiration); err != nil {
				return exist, err
			}

			return exist, nil
		}

		r.Logger.Debug("redis.repository/SIZE", logger.String("cache_key", key), logger.Any("cache_hit", true))

		return val, nil
	}
}

type RedisRepositoryOption[Entity any, Specification any, EntityJSON any, FallbackRepository Repository[Entity, Specification]] struct {
	Resource             string
	RedisDatabaseManager RedisDatabaseManager
	Logger               logger.Logger
	FallbackRepository   FallbackRepository
	Decode               Decoder[Entity, EntityJSON]
	Encode               Encoder[Entity, EntityJSON]
}

func NewRedisRepository[Entity any, Specification any, EntityJSON any, FallbackRepository Repository[Entity, Specification]](opt RedisRepositoryOption[Entity, Specification, EntityJSON, FallbackRepository]) Repository[Entity, Specification] {
	return &RedisRepository[Entity, Specification, EntityJSON, FallbackRepository]{
		Resource:           opt.Resource,
		RDBM:               opt.RedisDatabaseManager,
		Logger:             opt.Logger,
		FallbackRepository: opt.FallbackRepository,
		Encode:             opt.Encode,
		Decode:             opt.Decode,
	}
}
