package data

import (
	"context"

	"card/internal/biz"
	"gorm.io/gorm"
	"github.com/go-kratos/kratos/v2/log"
)

type CardRepo struct {
	data *Data
	log  *log.Helper
	table *gorm.DB
}

// NewGreeterRepo .
func NewCardRepo(data *Data, logger log.Logger) biz.CardRepo {
	return &CardRepo{
		data: data,
		log:  log.NewHelper(logger),
		table:data.db.Table("cardmanager"),
	}
}

func (r *CardRepo) Create(ctx context.Context, g *biz.Card) (*biz.Card, error) {
	r.log.WithContext(ctx).Infof("Creating Card: %+v", g)
	result := r.table.WithContext(ctx).Create(g)
	if result.Error != nil {
		return nil, result.Error
	}
	return g, nil
}
func (r *CardRepo) Update(ctx context.Context, g *biz.Card) (*biz.Card, error) {
    r.log.WithContext(ctx).Infof("Updating Card: %+v", g)
    result:=r.table.Model(g).Where("card_Id=?",g.CardId).Updates(biz.Card{})
    if result.Error != nil {
        return nil, result.Error
    }
    return g, nil
}
func (r *CardRepo) Delete(ctx context.Context, CardId int64) error {
    result := r.table.WithContext(ctx).Where("card_id = ?", CardId).Delete(&biz.Card{})
    if result.Error != nil {
        return result.Error
    }
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
    return nil
}

func (r *CardRepo) FindById(ctx context.Context, CardId int64) (*biz.Card, error) {
    result := r.table.WithContext(ctx).Where("card_id = ?",CardId).Find(&biz.Card{})
    if result.Error != nil {
        return nil, result.Error
    }
    return &biz.Card{},result.Error
}




