package data

import (
	"context"

	"message-service/internal/biz"
"fmt"
	"github.com/go-kratos/kratos/v2/log"
)

type greeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) biz.GreeterRepo {
	return &greeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}


func (r *greeterRepo) PaymentNotification(ctx context.Context,g *biz.Payment) (*biz.Payment, error) {
	return nil, nil
}

func (r *greeterRepo) CreateNotification(ctx context.Context, g *biz.Notification) (*biz.Notification, error) {
r.log.WithContext(ctx).Infof("Saving email to repo: %s", g.Email)
fmt.Println("Email id id:",g.Email)
	return g, nil
}

