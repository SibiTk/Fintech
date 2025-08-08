package biz

import (
	"context"

	v1 "message-service/api/helloworld/v1"
	//"message-service/internal/data"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Notification struct {
	 CustomerNumber string
	FirstName string
	Email string
	status string

}
type Payment struct {
	PaymentId        int64   `gorm:"column:payment_id;type:bigint;primaryKey"`
	FromAccountId    int64   `gorm:"column:from_account_id;type:bigint"`
	ToAccountId      int64   `gorm:"column:to_account_id;type:bigint"`
	PaymentType      string  `gorm:"column:payment_type;type:varchar"`
	Amount           int64 `gorm:"column:amount;type:bigint"`
	Currency         string  `gorm:"column:currency;type:varchar(3)"`
	Status           string  `gorm:"column:status;type:varchar"`
	PaymentMethod    string  `gorm:"column:payment_method;type:varchar"`
	ReferenceNumber  string  `gorm:"column:reference_number;type:varchar"`
	ExternalReference string  `gorm:"column:external_reference;type:varchar"`
}


// GreeterRepo is a Greater repo.
type GreeterRepo interface {
	Save(context.Context, *Notification) (*Notification, error)
	Update(context.Context, *Notification) (*Notification, error)
	FindByID(context.Context, int64) (*Notification, error)
	ListByHello(context.Context, string) ([]*Notification, error)
	ListAll(context.Context) ([]*Notification, error)
	PaymentNotification(context.Context,*Payment)(*Payment,error)
}

// GreeterUsecase is a Greeter usecase.
type NotificationUsecase struct {
	repo GreeterRepo
	log  *log.Helper
}


// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(repo GreeterRepo, logger log.Logger) *NotificationUsecase {
	return &NotificationUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *NotificationUsecase) CreateNotification(ctx context.Context, g *Notification) (*Notification, error) {
	uc.log.WithContext(ctx).Infof("CreateNotification: %v", g.Email)
	return uc.repo.Save(ctx, g)
}


func (uc *NotificationUsecase) PaymentNotification(ctx context.Context, g *Payment) (*Payment, error) {
	uc.log.WithContext(ctx).Infof("PaymentNotification: %+v ",g )
	return uc.repo.PaymentNotification(ctx, g)
}
