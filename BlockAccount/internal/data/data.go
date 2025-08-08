package data

import (
	"BlockAccount/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	DB *gorm.DB
}

// NewData .
func NewData(c *conf.Data, Logger log.Logger) (*Data, func(), error) {
	db, err := gorm.Open(postgres.Open(c.GetDatabase().GetSource()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {

		log.NewHelper(Logger).Info("closing the data resources")
	}
	return &Data{
		DB: db,
	}, cleanup, nil
}
