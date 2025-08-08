package data

import (
	"context"
	"gorm.io/gorm"
	"transaction/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type TransactionRepo struct {
	data  *Data
	log   *log.Helper
	table *gorm.DB
}

func NewTransactionRepo(data *Data, logger log.Logger) biz.TransactionRepo {
	return &TransactionRepo{
		data:  data,
		log:   log.NewHelper(logger),
		table: data.db.Table("transactions"),
	}
}

func (r *TransactionRepo) Create(ctx context.Context, t *biz.Transaction) (*biz.Transaction, error) {
	r.log.WithContext(ctx).Infof("Creating Transaction: %+v", t)
	if t.TransactionId == 0 {
		r.log.WithContext(ctx).Error("transaction_id is zero")
		return nil, gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Create(t)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to create transaction: %v", result.Error)
		return nil, result.Error
	}
	return t, nil
}

func (r *TransactionRepo) Update(ctx context.Context, t *biz.Transaction) (*biz.Transaction, error) {
	r.log.WithContext(ctx).Infof("Updating Transaction: %+v", t)
	if t.TransactionId == 0 {
		r.log.WithContext(ctx).Error("transaction_id is zero")
		return nil, gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Where("transaction_id = ?", t.TransactionId).Updates(t)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to update transaction: %v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		r.log.WithContext(ctx).Errorf("transaction not found: %d", t.TransactionId)
		return nil, gorm.ErrRecordNotFound
	}
	var updatedTransaction biz.Transaction
	if err := r.table.WithContext(ctx).Where("transaction_id = ?", t.TransactionId).First(&updatedTransaction).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to retrieve updated transaction: %v", err)
		return nil, err
	}
	return &updatedTransaction, nil
}

func (r *TransactionRepo) Delete(ctx context.Context, id int64) error {
	r.log.WithContext(ctx).Infof("Deleting Transaction ID: %d", id)
	if id == 0 {
		r.log.WithContext(ctx).Error("transaction_id is zero")
		return gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Where("transaction_id = ?", id).Delete(&biz.Transaction{})
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to delete transaction: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.log.WithContext(ctx).Errorf("transaction not found: %d", id)
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TransactionRepo) FindById(ctx context.Context, id int64) (*biz.Transaction, error) {
	r.log.WithContext(ctx).Infof("Finding Transaction by ID: %d", id)
	if id == 0 {
		r.log.WithContext(ctx).Error("transaction_id is zero")
		return nil, gorm.ErrInvalidData
	}
	var transaction biz.Transaction
	result := r.table.WithContext(ctx).Where("transaction_id = ?", id).First(&transaction)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to find transaction: %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &transaction, nil
}