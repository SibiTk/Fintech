package biz

import (
	"context"
	"errors"
	// v1 "customermanager/api/helloworld/v1"
	// "github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)
var ErrCustomerNotFound = errors.New("customer not found")


type Customer struct {
	CustomerId     int64  `gorm:"column:customer_id;type:bigint;primaryKey"`
	CustomerNumber string `gorm:"column:customer_number;type:varchar(255)"`
	FirstName      string `gorm:"column:first_name;type:varchar(255)"`
	LastName       string `gorm:"column:last_name;type:varchar(255)"`
	Email          string `gorm:"column:email;type:varchar(255)"`
	Phone          string `gorm:"column:phone;type:varchar(50)"`
	DateOfBirth    string `gorm:"column:date_of_birth;type:varchar"`
	Status         string `gorm:"column:status;type:varchar(50)"`
	KycStatus      string `gorm:"column:kyc_status;type:varchar(50)"`
}

type CustomerRepo interface {
	Create(context.Context, *Customer) (*Customer, error)
	Update(context.Context, *Customer) (*Customer, error)
	Delete(context.Context, int64) error
	FindById(context.Context, int64) (*Customer, error)
}

type CustomerUsecase struct {
	repo CustomerRepo
	log  *log.Helper
}



func NewCustomerUsecase(repo CustomerRepo, logger log.Logger) *CustomerUsecase {
	return &CustomerUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *CustomerUsecase) CreateCustomer(ctx context.Context, c *Customer) (*Customer, error) {
	uc.log.WithContext(ctx).Infof("Creating customer: %+v", c)
	return uc.repo.Create(ctx, c)
}

func (uc *CustomerUsecase) UpdateCustomer(ctx context.Context, c *Customer) (*Customer, error) {
	uc.log.WithContext(ctx).Infof("Updating customer: %+v", c)
	return uc.repo.Update(ctx, c)
}

func (uc *CustomerUsecase) DeleteCustomer(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("Deleting customer ID: %d", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *CustomerUsecase) DisplayCustomer(ctx context.Context, id int64) (*Customer, error) {
	uc.log.WithContext(ctx).Infof("Getting customer by ID: %d", id)
	return uc.repo.FindById(ctx, id)
}

