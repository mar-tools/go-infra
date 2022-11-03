package cache

import (
	"gopkg.in/redis.v3"
)

// redis initialize config
type Config struct {
	Endpoint string
	Password string
	Database int64
	PoolSize int
}

// redis failover initialize config
type ConfigFailover struct {
	MasterName        string
	SentinelEndpoints []string
	Password          string
	Database          int64
	PoolSize          int
}

type RedisCache struct {
	options         *redis.Options
	failoverOptions *redis.FailoverOptions
	client          *redis.Client
}

func Init(config *Config, configFailover *ConfigFailover) *RedisCache {
	if configFailover != nil {
		options := redis.FailoverOptions{
			MasterName:    configFailover.MasterName,
			SentinelAddrs: configFailover.SentinelEndpoints,
			Password:      configFailover.Password,
			DB:            configFailover.Database,
			PoolSize:      configFailover.PoolSize,
		}
		return &RedisCache{
			options:         nil,
			failoverOptions: &options,
			client:          redis.NewFailoverClient(&options),
		}
	} else if config != nil {
		options := redis.Options{
			Addr:     config.Endpoint,
			Password: config.Password,
			DB:       config.Database,
			PoolSize: config.PoolSize,
		}
		return &RedisCache{
			options:         &options,
			failoverOptions: nil,
			client:          redis.NewClient(&options),
		}
	} else {
		return nil
	}
}

func (r *RedisCache) Get(key []byte) ([]byte, error) {
	val, err := r.client.Get(string(key)).Result()
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r *RedisCache) Set(key []byte, value []byte) error {
	return r.client.Set(string(key), string(value), 0).Err()
}

func (r *RedisCache) Delete(key []byte) error {
	return r.client.Del(string(key)).Err()
}
