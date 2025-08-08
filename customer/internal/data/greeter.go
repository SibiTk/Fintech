package data

import (
	"context"
	"customer/internal/biz"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type CustomerRepo struct {
	data  *Data
	log   *log.Helper
	table *gorm.DB
}

func NewCustomerRepo(data *Data, logger log.Logger) biz.CustomerRepo {
	return &CustomerRepo{
		data:  data,
		log:   log.NewHelper(logger),
		table: data.db.Table("customermanager"),
	}
}

func (r *CustomerRepo) Create(ctx context.Context, c *biz.Customer) (*biz.Customer, error) {
	
	
	result := r.table.WithContext(ctx).Create(c)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to create customer: %v", result.Error)
		return nil, result.Error
	}
	return c, nil
}

func (r *CustomerRepo) Update(ctx context.Context, c *biz.Customer) (*biz.Customer, error) {
	r.log.WithContext(ctx).Infof("Updating Customer: %+v", c)
	if c.CustomerId == 0 {
		r.log.WithContext(ctx).Error("customer_id is zero")
		return nil, gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Where("customer_id = ?", c.CustomerId).Updates(c)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to update customer: %v", result.Error)
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		r.log.WithContext(ctx).Errorf("customer not found: %d", c.CustomerId)
		return nil, gorm.ErrRecordNotFound
	}
	var updatedCustomer biz.Customer
	if err := r.table.WithContext(ctx).Where("customer_id = ?", c.CustomerId).First(&updatedCustomer).Error; err != nil {
		r.log.WithContext(ctx).Errorf("failed to retrieve updated customer: %v", err)
		return nil, err
	}
	return &updatedCustomer, nil
}

func (r *CustomerRepo) Delete(ctx context.Context, id int64) error {
	r.log.WithContext(ctx).Infof("Deleting Customer ID: %d", id)
	if id == 0 {
		r.log.WithContext(ctx).Error("customer_id is zero")
		return gorm.ErrInvalidData
	}
	result := r.table.WithContext(ctx).Where("customer_id = ?", id).Delete(&biz.Customer{})
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to delete customer: %v", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		r.log.WithContext(ctx).Errorf("customer not found: %d", id)
		return gorm.ErrRecordNotFound
	}
	return nil
}
func (r *CustomerRepo) FindById(ctx context.Context, id int64) (*biz.Customer, error) {
	r.log.WithContext(ctx).Infof("Finding Customer by ID: %d", id)
	if id == 0 {
		r.log.WithContext(ctx).Error("customer_id is zero")
		return nil, gorm.ErrInvalidData
	}
	var customer biz.Customer
	result := r.table.WithContext(ctx).Where("customer_id = ?", id).First(&customer)
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("failed to find customer: %v", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, biz.ErrCustomerNotFound // <- return explicit error here
		}
		return nil, result.Error
	}
	return &customer, nil
}