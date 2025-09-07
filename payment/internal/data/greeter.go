package data

import (
	"context"
	"gorm.io/gorm"
	"payment/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
)

type PaymentRepo struct {
	data  *Data
	log   *log.Helper
	table *gorm.DB
}

func NewPaymentRepo(data *Data, logger log.Logger) biz.PaymentRepo {
	return &PaymentRepo{
		data:  data,
		log:   log.NewHelper(logger),
		table: data.db.Table("payments"),
	}
}

func (r *PaymentRepo) Create(ctx context.Context, p *biz.Payment) (*biz.Payment, error) {
	r.log.WithContext(ctx).Infof("Creating Payment: PaymentId=%d, FromAccountId=%d", p.PaymentId, p.FromAccountId)
	if p.PaymentId == 0 {
		r.log.WithContext(ctx).Error("payment_id is zero")
		return nil, gorm.ErrInvalidData
	}
	
	if p.FromAccountId == 0 || p.ToAccountId == 0 || p.PaymentType == "" || p.Amount == 0 || p.Currency == "" || p.Status == "" || p.PaymentMethod == "" {
		r.log.WithContext(ctx).Error("required fields are missing")
		return nil, gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Create(p)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to create payment: %v", result.Error)
		return nil, result.Error
	}
	return p, nil
}

func (r *PaymentRepo) Update(ctx context.Context, p *biz.Payment) (*biz.Payment, error) {
	r.log.WithContext(ctx).Infof("Updating Payment: %+v", p)
	if p.PaymentId == 0 {
		r.log.WithContext(ctx).Error("payment_id is zero")
		return nil, gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Where("payment_id = ?", p.PaymentId).Updates(p)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to update payment: %v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		r.log.WithContext(ctx).Errorf("payment not found: %d", p.PaymentId)
		return nil, gorm.ErrRecordNotFound
	}
	var updatedPayment biz.Payment
	if err := r.table.WithContext(ctx).Where("payment_id = ?", p.PaymentId).First(&updatedPayment).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to retrieve updated payment: %v", err)
		return nil, err
	}
	return &updatedPayment, nil
}

func (r *PaymentRepo) Delete(ctx context.Context, id int64) error {
	r.log.WithContext(ctx).Infof("Deleting Payment ID: %d", id)
	if id == 0 {
		r.log.WithContext(ctx).Error("payment_id is zero")
		return gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Where("payment_id = ?", id).Delete(&biz.Payment{})
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to delete payment: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.log.WithContext(ctx).Errorf("payment not found: %d", id)
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PaymentRepo) FindById(ctx context.Context, id int64) (*biz.Payment, error) {
	r.log.WithContext(ctx).Infof("Finding Payment by ID: %d", id)
	if id == 0 {
		r.log.WithContext(ctx).Error("payment_id is zero")
		return nil, gorm.ErrInvalidData
	}
	var payment biz.Payment
	result := r.table.WithContext(ctx).Where("payment_id = ?", id).First(&payment)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to find payment: %v", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &payment, nil
}