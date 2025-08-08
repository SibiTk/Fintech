package biz

import (
    "context"
  //  v1 "rel_cus_acc/api/helloworld/v1"
   
    "github.com/go-kratos/kratos/v2/log"
)



type RelCusAcc struct {
    Id                  int64  `gorm:"column:id;type:bigserial;primaryKey;autoIncrement"`
    CustomerId          int64  `gorm:"column:customer_id;type:bigint;not null"`
    CustomerNumber      string `gorm:"column:customer_number;type:varchar(255);not null"`
    FirstName           string `gorm:"column:first_name;type:varchar(255);not null"`
    LastName            string `gorm:"column:last_name;type:varchar(255);not null"`
    Email               string `gorm:"column:email;type:varchar(255);not null"`
    Phone               string `gorm:"column:phone;type:varchar(50);not null"`
    DateOfBirth         string `gorm:"column:date_of_birth;type:varchar"`
    CustomerStatus      string `gorm:"column:customer_status;type:varchar(50);not null"`
    KycStatus           string `gorm:"column:kyc_status;type:varchar(50);not null"`
    CustomerCreatedAt   string `gorm:"column:customer_created_at;type:timestamp;not null"`
    CustomerUpdatedAt   string `gorm:"column:customer_updated_at;type:timestamp;not null"`
    AccountId           int64  `gorm:"column:account_id;type:bigint;not null"`
    AccountNumber       string `gorm:"column:account_number;type:varchar(16);not null"`
    AccountType         string `gorm:"column:account_type;type:varchar(50);not null"`
    Currency            string `gorm:"column:currency;type:varchar(32);not null"`
    AccountStatus       string `gorm:"column:account_status;type:varchar(20);not null"`
    AvailableBalance    int64  `gorm:"column:available_balance;type:bigint;default:0"`
    PendingBalance      int64  `gorm:"column:pending_balance;type:bigint;default:0"`
    CreditLimit         string `gorm:"column:credit_limit;type:varchar(20)"`
    LastTransactionAt   string `gorm:"column:last_transaction_at;type:timestamp"`
    AccountCreatedAt    string `gorm:"column:account_created_at;type:timestamp;not null"`
    LastUsedAt          string `gorm:"column:last_used_at;type:timestamp;not null"`
}

// RelCusAccRepo is the interface that must be implemented by the data layer.
type RelCusAccRepo interface {
    Create(ctx context.Context, rel *RelCusAcc) (*RelCusAcc, error)
    Update(ctx context.Context, rel *RelCusAcc) (*RelCusAcc, error)
    Delete(ctx context.Context, id int64) error
    FindById(ctx context.Context, id int64) (*RelCusAcc, error)
    FindByCustomerId(ctx context.Context, customerId int64) ([]*RelCusAcc, error)
    FindByAccountId(ctx context.Context, accountId int64) (*RelCusAcc, error)
  
}

// RelCusAccUsecase handles business logic for customer-account relationships.
type RelCusAccUsecase struct {
    repo RelCusAccRepo
    log  *log.Helper
}

func NewRelCusAccUsecase(repo RelCusAccRepo, logger log.Logger) *RelCusAccUsecase {
    return &RelCusAccUsecase{
        repo: repo,
        log:  log.NewHelper(logger),
    }
}

func (uc *RelCusAccUsecase) CreateRelation(ctx context.Context, rel *RelCusAcc) (*RelCusAcc, error) {
    uc.log.WithContext(ctx).Infof("Creating customer-account relationship: %+v", rel)
    return uc.repo.Create(ctx, rel)
}

func (uc *RelCusAccUsecase) UpdateRelation(ctx context.Context, rel *RelCusAcc) (*RelCusAcc, error) {
    uc.log.WithContext(ctx).Infof("Updating customer-account relationship: %+v", rel)
    return uc.repo.Update(ctx, rel)
}

func (uc *RelCusAccUsecase) DeleteRelation(ctx context.Context, id int64) error {
    uc.log.WithContext(ctx).Infof("Deleting relationship ID: %d", id)
    return uc.repo.Delete(ctx, id)
}

func (uc *RelCusAccUsecase) GetRelationById(ctx context.Context, id int64) (*RelCusAcc, error) {
    uc.log.WithContext(ctx).Infof("Getting relation by ID: %d", id)
    return uc.repo.FindById(ctx, id)
}

func (uc *RelCusAccUsecase) GetRelationsByCustomer(ctx context.Context, customerId int64) ([]*RelCusAcc, error) {
    uc.log.WithContext(ctx).Infof("Getting relations for customer ID: %d", customerId)
    return uc.repo.FindByCustomerId(ctx, customerId)
}

func (uc *RelCusAccUsecase) GetRelationByAccount(ctx context.Context, accountId int64) (*RelCusAcc, error) {
    uc.log.WithContext(ctx).Infof("Getting relation for account ID: %d", accountId)
    return uc.repo.FindByAccountId(ctx, accountId)
}


