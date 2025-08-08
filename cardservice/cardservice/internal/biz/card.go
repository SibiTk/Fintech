package biz

import (
	"context"
//	v1 "card/api/helloworld/v1"
	//"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// var (
// 	ErrCardNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "card not found")
// )
type Card struct {
	CardId       int64   `gorm:"column:card_id;type:bigint;primaryKey"`
	AccountId    int64   `gorm:"column:account_id"`
	CardNumber   string  `gorm:"column:card_number"`
	CardType     string  `gorm:"column:card_type"`
	ExpiryDate   string  `gorm:"column:expiry_date"`
	DailyLimit   float64 `gorm:"column:daily_limit"`
	MonthlyLimit float64 `gorm:"column:monthly_limit"`
	PinAttempts  int     `gorm:"column:pin_attempts"`
}




type CardRepo interface {
	Create(context.Context, *Card) (*Card, error)
	Update(context.Context, *Card) (*Card, error)
	Delete(context.Context, int64) error
	FindById(context.Context, int64) (*Card, error)
}

type CardUsecase struct {
	repo CardRepo
	log  *log.Helper
}

func NewCardUsecase(repo CardRepo, logger log.Logger) *CardUsecase {
	return &CardUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (uc *CardUsecase) CreateCard(ctx context.Context, c *Card) (*Card, error) {
	uc.log.WithContext(ctx).Infof("Creating card: %+v", c)
	return uc.repo.Create(ctx, c)
}

func (uc *CardUsecase) UpdateCard(ctx context.Context, c *Card) (*Card, error) {
	uc.log.WithContext(ctx).Infof("Updating card: %+v", c)
	return uc.repo.Update(ctx, c)
}

func (uc *CardUsecase) DeleteCard(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("Deleting card ID: %d", id)
	return uc.repo.Delete(ctx, id)
}

func (uc *CardUsecase) GetCardById(ctx context.Context, id int64) (*Card, error) {
	uc.log.WithContext(ctx).Infof("Getting card by ID: %d", id)
	return uc.repo.FindById(ctx, id)
}


