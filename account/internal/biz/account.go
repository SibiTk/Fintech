package biz

import (
	"context"

	v1 "account/api/helloworld/v1"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type Account struct {
    AccountId         int64  `gorm:"primaryKey;column:account_id"`
    CustomerId        int64  `gorm:"column:customer_id"`
    AccountNumber     int64  `gorm:"column:account_number"`
    AccountType       string `gorm:"column:account_type"`
    Currency          string `gorm:"column:currency"`
    Status            string `gorm:"column:status"`
    AvailableBalance  int64  `gorm:"column:available_balance"`
    PendingBalance    int64  `gorm:"column:pending_balance"`
    CreditLimit       string `gorm:"column:credit_limit"`
    LastTransactionAt string `gorm:"column:last_transaction_at"`
}

// AccountRepo is the interface that must be implemented by the data layer.
type AccountRepo interface {
	Create(ctx context.Context, acc *Account) (*Account, error)
	Update(ctx context.Context, acc *Account) (*Account, error)
	Delete(ctx context.Context, customerId int64) error
	GetByCustomerID(ctx context.Context, customerId int64) ([]*Account, error)
	GetByAccountID(ctx context.Context, accountId int64) ([]*Account, error) 
}

// AccountUsecase handles business logic for accounts.
type AccountUsecase struct {
	repo AccountRepo
	log  *log.Helper
}

func NewAccountUsecase(repo AccountRepo, logger log.Logger) *AccountUsecase {
	return &AccountUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *AccountUsecase) Create(ctx context.Context, acc *Account) (*Account, error) {
	uc.log.WithContext(ctx).Infof("Create Account: %+v", acc)
	return uc.repo.Create(ctx, acc)
}

func (uc *AccountUsecase) Update(ctx context.Context, acc *Account) (*Account, error) {
	uc.log.WithContext(ctx).Infof("Update Account: %+v", acc)
	return uc.repo.Update(ctx, acc)
}

func (uc *AccountUsecase) Delete(ctx context.Context, customerId int64) error {
	uc.log.WithContext(ctx).Infof("Delete Account for customerId: %d", customerId)
	return uc.repo.Delete(ctx, customerId)
}

func (uc *AccountUsecase) GetByCustomerID(ctx context.Context, customerId int64) ([]*Account, error) {
	uc.log.WithContext(ctx).Infof("Get accounts for customerId: %d", customerId)
	return uc.repo.GetByCustomerID(ctx, customerId)
}
func (uc *AccountUsecase) GetByAccountID(ctx context.Context, accountId int64) ([]*Account, error) {
    uc.log.WithContext(ctx).Infof("Get account(s) for accountId: %d", accountId)
    return uc.repo.GetByAccountID(ctx, accountId)
}
