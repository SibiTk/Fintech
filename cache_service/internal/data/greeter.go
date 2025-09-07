package data

import (
	"context"

	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"cache_service/internal/conf"
)

func NewRedisServer(r *conf.InternalConf) redis.UniversalClient {

	if r.Redis.Mode == "ring" {
		ringClient := redis.NewRing(&redis.RingOptions{
			Addrs:    r.Redis.Ring.Address,
			Username: r.Redis.Username,
			Password: r.Redis.Password,
		})
		zap.S().Info(ringClient.Ping(context.Background()))
		return ringClient

	}
	return nil
}
