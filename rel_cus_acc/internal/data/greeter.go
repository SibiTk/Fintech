package data

import (
    "context"
    "time"
    "rel_cus_acc/internal/biz"
    "github.com/go-kratos/kratos/v2/log"
    "gorm.io/gorm"
)

type RelCusAccRepo struct {
    data  *Data
    log   *log.Helper
    table *gorm.DB
}

func NewRelCusAccRepo(data *Data, logger log.Logger) biz.RelCusAccRepo {
    return &RelCusAccRepo{
        data:  data,
        log:   log.NewHelper(logger),
        table: data.db.Table("rel_cus_acc"),
    }
}

// Create a new customer-account relationship
func (r *RelCusAccRepo) Create(ctx context.Context, rel *biz.RelCusAcc) (*biz.RelCusAcc, error) {
    r.log.WithContext(ctx).Infof("Creating RelCusAcc: %+v", rel)
    
    // Set timestamps if not provided
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    if rel.CustomerCreatedAt == "" {
        rel.CustomerCreatedAt = currentTime
    }
    if rel.CustomerUpdatedAt == "" {
        rel.CustomerUpdatedAt = currentTime
    }
    if rel.AccountCreatedAt == "" {
        rel.AccountCreatedAt = currentTime
    }
    if rel.LastUsedAt == "" {
        rel.LastUsedAt = currentTime
    }
    
    result := r.table.WithContext(ctx).Create(rel)
    if result.Error != nil {
        return nil, result.Error
    }
    return rel, nil
}

// Update an existing relationship by ID
func (r *RelCusAccRepo) Update(ctx context.Context, rel *biz.RelCusAcc) (*biz.RelCusAcc, error) {
    r.log.WithContext(ctx).Infof("Updating RelCusAcc: %+v", rel)
    
    // Set update timestamp
    rel.CustomerUpdatedAt = time.Now().Format("2006-01-02 15:04:05")
    rel.LastUsedAt = time.Now().Format("2006-01-02 15:04:05")
    
    result := r.table.WithContext(ctx).Where("id = ?", rel.Id).Updates(rel)
    if result.Error != nil {
        return nil, result.Error
    }
    return rel, nil
}

// Delete relationship by ID
func (r *RelCusAccRepo) Delete(ctx context.Context, id int64) error {
    r.log.WithContext(ctx).Infof("Deleting relationship for ID: %d", id)
    
    result := r.table.WithContext(ctx).Where("id = ?", id).Delete(&biz.RelCusAcc{})
    if result.Error != nil {
        return result.Error
    }
    return nil
}

// Find relationship by ID
func (r *RelCusAccRepo) FindById(ctx context.Context, id int64) (*biz.RelCusAcc, error) {
    r.log.WithContext(ctx).Infof("Fetching relationship for ID: %d", id)
    
    var relation biz.RelCusAcc
    result := r.table.WithContext(ctx).Where("id = ?", id).First(&relation)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }
    return &relation, nil
}

// Get all relationships by CustomerID
func (r *RelCusAccRepo) FindByCustomerId(ctx context.Context, customerId int64) ([]*biz.RelCusAcc, error) {
    r.log.WithContext(ctx).Infof("Fetching relationships for customer ID: %d", customerId)
    
    var relations []*biz.RelCusAcc
    result := r.table.WithContext(ctx).Where("customer_id = ?", customerId).Find(&relations)
    if result.Error != nil {
        return nil, result.Error
    }
    return relations, nil
}

// Get relationship by AccountID
func (r *RelCusAccRepo) FindByAccountId(ctx context.Context, accountId int64) (*biz.RelCusAcc, error) {
    r.log.WithContext(ctx).Infof("Fetching relationship for account ID: %d", accountId)
    
    var relation biz.RelCusAcc
    result := r.table.WithContext(ctx).Where("account_id = ?", accountId).First(&relation)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, nil
        }
        return nil, result.Error
    }
    return &relation, nil
}


