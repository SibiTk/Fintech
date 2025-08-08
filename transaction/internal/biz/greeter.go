package biz

import (
	"context"
	//v1 "transaction/api/helloworld/v1"
	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)



type Transaction struct {
	TransactionId        int64   `gorm:"column:transaction_id;type:bigint;primaryKey"`
	AccountId           int64  `gorm:"column:account_id;type:varchar"`
	RelatedTransactionId int64 `gorm:"column:related_transaction_id;type:varchar"`
	TransactionType     string  `gorm:"column:transaction_type;type:varchar"`
	Amount              float64 `gorm:"column:amount;type:numeric(20,4)"`
	Currency            string  `gorm:"column:currency;type:varchar(3)"`
	Status              string  `gorm:"column:status;type:varchar"`
	Description         string  `gorm:"column:description;type:text"`
	ReferenceNumber     string  `gorm:"column:reference_number;type:varchar"`
	PostingDate         string  `gorm:"column:posting_date;type:varchar"`
}

type TransactionRepo interface {
	Create(context.Context, *Transaction) (*Transaction, error)
	Update(context.Context, *Transaction) (*Transaction, error)
	Delete(context.Context, int64) error
	FindById(context.Context, int64) (*Transaction, error)
}

type TransactionUsecase struct {
	repo TransactionRepo
	log  *log.Helper
}

func NewTransactionUsecase(repo TransactionRepo, logger log.Logger) *TransactionUsecase {
	return &TransactionUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *TransactionUsecase) CreateTransaction(ctx context.Context, t *Transaction) (*Transaction, error) {
	uc.log.WithContext(ctx).Infof("Creating transaction: %+v", t)
	return uc.repo.Create(ctx, t)
}

func (uc *TransactionUsecase) UpdateTransaction(ctx context.Context, t *Transaction) (*Transaction, error) {
	uc.log.WithContext(ctx).Infof("Updating transaction: %+v", t)
	return uc.repo.Update(ctx, t)
}

func (uc *TransactionUsecase) DeleteTransaction(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("Deleting transaction ID: %d", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *TransactionUsecase) GetTransactionById(ctx context.Context, id int64) (*Transaction, error) {
	uc.log.WithContext(ctx).Infof("Getting transaction by ID: %d", id)
	return uc.repo.FindById(ctx, id)
}