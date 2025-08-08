package data

import (
	"context"
	"account/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type AccountRepo struct {
	data  *Data
	log   *log.Helper
	table *gorm.DB
}

func NewAccountRepo(data *Data, logger log.Logger) biz.AccountRepo {
	return &AccountRepo{
		data:  data,
		log:   log.NewHelper(logger),
		table: data.db.Table("accountservice"),
	}
}

// Create a new account
func (r *AccountRepo) Create(ctx context.Context, g *biz.Account) (*biz.Account, error) {
	r.log.WithContext(ctx).Infof("Creating Account: %+v", g)
	result := r.table.WithContext(ctx).Create(g)
	if result.Error != nil {
		return nil, result.Error
	}
	return g, nil
}

// Update an existing account by CustomerID
func (r *AccountRepo) Update(ctx context.Context, c *biz.Account) (*biz.Account, error) {
	r.log.WithContext(ctx).Infof("Updating account: %+v", c)
	result := r.table.WithContext(ctx).Where("customer_id = ?", c.CustomerId).Updates(c)
	if result.Error != nil {
		return nil, result.Error
	}
	return c, nil
}


func (r *AccountRepo) Delete(ctx context.Context, customerId int64) error {
	r.log.WithContext(ctx).Infof("Deleting account for customer: %d", customerId)

	result := r.table.WithContext(ctx).Where("customer_id = ?", customerId).Delete(&biz.Account{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Get all accounts by CustomerID
func (r *AccountRepo) GetByCustomerID(ctx context.Context, customerId int64) ([]*biz.Account, error) {
	r.log.WithContext(ctx).Infof("Fetching accounts for customer ID: %d", customerId)

	var accounts []*biz.Account
	result := r.table.WithContext(ctx).Where("customer_id = ?", customerId).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}
func (r *AccountRepo) GetByAccountID(ctx context.Context, accountId int64) ([]*biz.Account, error) {
    r.log.WithContext(ctx).Infof("Fetching account(s) for account ID: %d", accountId)
    var accounts []*biz.Account
    result := r.table.WithContext(ctx).Where("account_id = ?", accountId).Find(&accounts)
    if result.Error != nil {
        return nil, result.Error
    }
    return accounts, nil
}
