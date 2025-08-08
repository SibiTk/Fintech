package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)



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

type PaymentRepo interface {
	Create(context.Context, *Payment) (*Payment, error)
	Update(context.Context, *Payment) (*Payment, error)
	Delete(context.Context, int64) error
	FindById(context.Context, int64) (*Payment, error)
}

type PaymentUsecase struct {
	repo PaymentRepo
	log  *log.Helper
}

func NewPaymentUsecase(repo PaymentRepo, logger log.Logger) *PaymentUsecase {
	return &PaymentUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *PaymentUsecase) CreatePayment(ctx context.Context, p *Payment) (*Payment, error) {
	uc.log.WithContext(ctx).Infof("Creating payment: %+v", p)
	return uc.repo.Create(ctx, p)
}

func (uc *PaymentUsecase) UpdatePayment(ctx context.Context, p *Payment) (*Payment, error) {
	uc.log.WithContext(ctx).Infof("Updating payment: %+v", p)
	return uc.repo.Update(ctx, p)
}

func (uc *PaymentUsecase) DeletePayment(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("Deleting payment ID: %d", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *PaymentUsecase) GetPaymentById(ctx context.Context, id int64) (*Payment, error) {
	uc.log.WithContext(ctx).Infof("Getting payment by ID: %d", id)
	return uc.repo.FindById(ctx, id)
}