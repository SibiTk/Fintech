package service

import (
    v1 "card/api/helloworld/v1"
    "card/internal/biz"
    "card/internal/validation"
    "context"
    "fmt"
    clientv3 "go.etcd.io/etcd/client/v3"
    "time"
)

type CardService struct {
    v1.UnimplementedCardServer
    uc         *biz.CardUsecase
    etcdClient *clientv3.Client
}

// NewEtcdClient creates an etcd client connected to the provided endpoints
func NewEtcdClient(endpoints []string) (*clientv3.Client, error) {
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   endpoints,
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create etcd client: %w", err)
    }
    return cli, nil
}

// registerService registers the service in etcd with a lease and keepalive
func (s *CardService) registerService(name, addr string, leaseTTL int64) error {
    leaseResp, err := s.etcdClient.Grant(context.Background(), leaseTTL)
    if err != nil {
        return err
    }
    _, err = s.etcdClient.Put(context.Background(), "/services/"+name, addr, clientv3.WithLease(leaseResp.ID))
    if err != nil {
        return err
    }
    ch, err := s.etcdClient.KeepAlive(context.Background(), leaseResp.ID)
    if err != nil {
        return err
    }
    go func() {
        for range ch {
           
        }
    }()
    return nil
}


func NewCardService(uc *biz.CardUsecase, etcdEndpoints []string) (*CardService, error) {
    etcdClient, err := NewEtcdClient(etcdEndpoints)
    if err != nil {
        return nil, err
    }

    s := &CardService{
        uc:         uc,
        etcdClient: etcdClient,
    }

    // Register the card service with etcd for service discovery
    if err := s.registerService("card-service", "localhost:9070", 10); err != nil {
        return nil, err
    }

    return s, nil
}

// CreateCard performs validation, existence check, generate card number, and persists card details
func (s *CardService) CreateCard(ctx context.Context, req *v1.CreateRequest) (*v1.CreateReply, error) {
    err := validation.ValidateCreateCardRequest(
        req.AccountId,
        req.Cardtype,
        req.Expirydate,
        float64(req.Dailylimit),
        float64(req.Monthlylimit),
        int(req.Pinattempts),
    )
    if err != nil {
        return &v1.CreateReply{
            Message: "Validation Failed for CardService: " + err.Error(),
        }, nil
    }
    exists, err := validation.CheckAccountExistsViaGRPC(req.AccountId)
    if err != nil {
        return &v1.CreateReply{
            Message: "Failed to verify customer existence: " + err.Error(),
        }, nil
    }
    if !exists {
        return &v1.CreateReply{
            Message: "AccountId does not exist",
        }, nil
    }
    cardNumber, err := validation.GenerateCardNumber()
    if err != nil {
        return &v1.CreateReply{
            Message: "Failed to generate account number: " + err.Error(),
        }, nil
    }
    fmt.Println("Generated account number:", cardNumber)

    card := &biz.Card{
        AccountId:    req.AccountId,
        CardNumber:   cardNumber,
        CardType:     req.Cardtype,
        ExpiryDate:   req.Expirydate,
        DailyLimit:   float64(req.Dailylimit),
        MonthlyLimit: float64(req.Monthlylimit),
        PinAttempts:  int(req.Pinattempts),
    }

    a, err := s.uc.CreateCard(ctx, card)
    if err != nil {
        return nil, err
    }

    return &v1.CreateReply{
        CardId:       a.CardId,
        AccountId:    a.AccountId,
        Cardnumber:   a.CardNumber,
        Cardtype:     a.CardType,
        Expirydate:   a.ExpiryDate,
        Dailylimit:   int64(a.DailyLimit),
        Monthlylimit: int64(a.MonthlyLimit),
        Pinattempts:  int64(a.PinAttempts),
        Message:      "Card Details are stored",
    }, nil
}


func (s *CardService) UpdateCard(ctx context.Context, req *v1.UpdateCardRequest) (*v1.UpdateCardReply, error) {
    // Validate request using validation package
    err := validation.ValidateUpdateCardRequest(
        
        req.AccountId,
        req.Cardnumber,
        req.Cardtype,
        float64(req.Dailylimit),
        float64(req.Monthlylimit),
    )
  if err != nil {
        return &v1.UpdateCardReply{
            Message:"Validation Failed for CardService:"+err.Error(),
        },nil
    }
    
    
    card := &biz.Card{
        CardId:       req.CardId,
        AccountId:    req.AccountId,
        CardNumber:   req.Cardnumber,
        CardType:     req.Cardtype,
        DailyLimit:   float64(req.Dailylimit),
        MonthlyLimit: float64(req.Monthlylimit),
    }

    updatedCard, err := s.uc.UpdateCard(ctx, card)
    if err != nil {
        return nil, err
    }

    return &v1.UpdateCardReply{
        Message: "Card Details are updated",
        Display: []*v1.CreateReply{{
            CardId:       updatedCard.CardId,
            AccountId:    updatedCard.AccountId,
            Cardnumber:   updatedCard.CardNumber,
            Cardtype:     updatedCard.CardType,
            Expirydate:   updatedCard.ExpiryDate,
            Dailylimit:   int64(updatedCard.DailyLimit),
            Monthlylimit: int64(updatedCard.MonthlyLimit),
            Pinattempts:  int64(updatedCard.PinAttempts),
        }},
    }, nil
}

func (s *CardService) DeleteCard(ctx context.Context, req *v1.DeleteCardRequest) (*v1.DeleteCardReply, error) {
    // Validate card ID using validation package
   
    
    err := s.uc.DeleteCard(ctx, req.CardId)
    if err != nil {
        return nil, err
    }

    return &v1.DeleteCardReply{
        Message: "Card Details are deleted",
    }, nil
}

func (s *CardService) DisplayCard(ctx context.Context, req *v1.FindByIdRequest) (*v1.FindByIdReply, error) {
    // Validate card ID using validation package
  
    
    card, err := s.uc.GetCardById(ctx, req.CardId)
    if err != nil {
        return nil, err
    }

    reply := &v1.CreateReply{
        CardId:       card.CardId,
        AccountId:    card.AccountId,
        Cardnumber:   card.CardNumber,
        Cardtype:     card.CardType,
        Expirydate:   card.ExpiryDate,
        Dailylimit:   int64(card.DailyLimit),
        Monthlylimit: int64(card.MonthlyLimit),
        Pinattempts:  int64(card.PinAttempts),
    }

    return &v1.FindByIdReply{
        Display: []*v1.CreateReply{reply},
    }, nil
}
