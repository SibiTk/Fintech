package data

import (
    "context"
    "fmt"
    v1 "BlockAccount/api/helloworld/v1"
    "github.com/go-kratos/kratos/v2/log"
    "gorm.io/gorm"
)


type AccountBlock struct {
    CustomerId  int64  `gorm:"column:customer_id;primaryKey"`
    Description string `gorm:"column:description"`
    Source      string `gorm:"column:source"`      
    Status      string `gorm:"column:status"`    
}

type GreeterRepo struct {
    data  *Data
    log   *log.Helper
    table *gorm.DB
}

func NewGreeterRepo(data *Data, logger log.Logger) *GreeterRepo {
    return &GreeterRepo{
        data:  data,
        log:   log.NewHelper(logger),
        table: data.DB.Model(AccountBlock{}),
    }
}


func (r *GreeterRepo) Save(ctx context.Context, g *v1.SaveAccBlockRequest) error {
    account := &AccountBlock{
        CustomerId:  g.CustomerId,
        Description: g.Description,
        Source:      g.Source,   
        Status:      g.Status,    
    }
    result := r.table.WithContext(ctx).Create(account)
    fmt.Println("create account data is trigger")
    if result.Error != nil {
        return result.Error
    }
    return nil
}


func (r *GreeterRepo) Get(ctx context.Context, g *v1.GetAccBlockRequest) (*v1.GetAccBlockReply, error) {
    var account AccountBlock
    query := "SELECT * from account_blocks where customer_id = ?"
    result := r.table.WithContext(ctx).Raw(query, g.CustomerId).Scan(&account)
    if result.Error != nil {
        return nil, result.Error
    }
    return &v1.GetAccBlockReply{
        CustomerId:  account.CustomerId,
        Description: account.Description,
        Source:      account.Source,  
        Status:      account.Status,  
    }, nil
}


func (r *GreeterRepo) UpdateAccBlock(ctx context.Context, g *v1.UpdateAccBlockRequest) (*v1.UpdateAccBlockReply, error) {
 
    query := "UPDATE account_blocks SET description = ?, source = ?, status = ? WHERE customer_id = ?"
    result := r.table.WithContext(ctx).Exec(query, g.Description, g.Source, g.Status, g.CustomerId)
    if result.Error != nil {
        return nil, result.Error
    }
    if result.RowsAffected == 0 {
        return nil, fmt.Errorf("no rows updated for customer_id %d", g.CustomerId)
    }

    
    return &v1.UpdateAccBlockReply{
        CustomerId:  g.CustomerId,
        Description: g.Description,
        Source:      g.Source,
        Status:      g.Status,
    }, nil
}

func (r *GreeterRepo) GetByCustomerAndSource(ctx context.Context, customerId int64, source string) (*AccountBlock, error) {
    var account AccountBlock
    query := "SELECT * FROM account_blocks WHERE customer_id = ? AND source = ?"
    result := r.table.WithContext(ctx).Raw(query, customerId, source).Scan(&account)
    
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return nil, gorm.ErrRecordNotFound
        }
        return nil, result.Error
    }
    
    // If no rows found, return ErrRecordNotFound
    if result.RowsAffected == 0 {
        return nil, gorm.ErrRecordNotFound
    }
    
    return &account, nil
}
