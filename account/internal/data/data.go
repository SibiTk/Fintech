package data

import (
    "context"

    "account/internal/conf"

    "github.com/go-kratos/kratos/v2/log"
    "github.com/google/wire"
    "github.com/redis/go-redis/v9"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// Providers
var ProviderSet = wire.NewSet(NewData, NewAccountRepo)

type Data struct {
    db  *gorm.DB
    rdb *redis.Client
}

func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
    helper := log.NewHelper(logger)

    // Postgres
    db, err := gorm.Open(postgres.Open(c.GetDatabase().GetSource()), &gorm.Config{})
    if err != nil {
        return nil, nil, err
    }

    // Redis: only use fields that exist; default others
    addr := "localhost:6379"
    if rc := c.GetRedis(); rc != nil && rc.GetAddr() != "" {
        addr = rc.GetAddr()
    }
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: "", 
        DB:       0, 
    })
   if err := rdb.Ping(context.Background()).Err(); err != nil {
    helper.Warnf("redis ping failed: %v", err)
}

    cleanup := func() {
        helper.Info("closing the data resources")
        if sqlDB, err := db.DB(); err == nil {
            _ = sqlDB.Close()
        }
        _ = rdb.Close()
    }

    return &Data{db: db, rdb: rdb}, cleanup, nil
}
