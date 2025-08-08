package service

import (
	"context"
	"time"
	//"encoding/json"
	"fmt"
	"log"

	v1 "payment/api/helloworld/v1"
	"payment/internal/biz"

	//"payment/internal/client"

	accountpb "account/api/helloworld/v1"
	cu "customer/api/helloworld/v1"
	tra "transaction/api/helloworld/v1"

	"github.com/nats-io/nats.go"

)

type PaymentService struct {
	AccountClient  accountpb.AccountClient
	CustomerClient cu.CustomerManagerClient
	transactionClient tra.CardClient
	v1.UnimplementedPaymentServer
	uc *biz.PaymentUsecase
}

type PaymentStatus struct {
	FromEmail     string
	ToEmail       string
	FromBalance   int64
	ToBalance     int64
	FromFirstName string
	ToFirstName   string
	Amount        int64
	AccountNumber int64
}

func NewPaymentService(uc *biz.PaymentUsecase, accClient accountpb.AccountClient, cusClient cu.CustomerManagerClient,traClient tra.CardClient) *PaymentService {
	return &PaymentService{uc: uc, AccountClient: accClient, CustomerClient: cusClient,transactionClient:traClient}
}

// Connect to NATS
func connectToNats() (*nats.Conn, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return nc, nil
}



func (s *PaymentService) CreatePayment(ctx context.Context, req *v1.CreatePaymentRequest) (*v1.CreatePaymentReply, error) {






	// Create Payment from the request
	payment := &biz.Payment{
		PaymentId:         req.PaymentId,
		FromAccountId:     req.FromAccountId,
		ToAccountId:       req.ToAccountId,
		PaymentType:       req.PaymentType,
		Amount:            int64(req.Amount),
		Currency:          req.Currency,
		Status:            req.Status,
		PaymentMethod:     req.PaymentMethod,
		ReferenceNumber:   req.ReferenceNumber,
		ExternalReference: req.ExternalReference,
	}


	// p, err := s.uc.CreatePayment(ctx, payment)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create payment: %w", err)
	// }
 
    
    trans1,err:=s.transactionClient.CreateTransaction(ctx,&tra.CreateTransactionRequest{Transactionid: req.PaymentId,AccountId: req.FromAccountId,RelatedTransactionId: req.ToAccountId,Amount: req.Amount,TransactionType: req.PaymentType,Currency: req.Currency,Status: req.Status,Description:"transaction details is stored successfully",ReferenceNumber: req.ReferenceNumber,PostingDate: time.DateTime})
    if err!=nil{
		fmt.Print(trans1)
        return &v1.CreatePaymentReply{
            Message: "There is an issue with creating a transaction",
        },nil
    }
	// Fetch From Account
	fromResp, err := s.AccountClient.GetAccountWithId(ctx, &accountpb.AccountIdRequest{AccountId: req.FromAccountId})
	if err != nil {
		log.Println("Failed to get From Account:", err)
		return nil, fmt.Errorf("failed to get FromAccount: %w", err)
	}
	if len(fromResp.Accounts) == 0 {
		return nil, fmt.Errorf("no account found for FromAccountId: %v", req.FromAccountId)
	}
	fromAcc := fromResp.Accounts[0]

	// Fetch To Account
	toResp, err := s.AccountClient.GetAccountWithId(ctx, &accountpb.AccountIdRequest{AccountId: req.ToAccountId})
	if err != nil {
		log.Println("Failed to get To Account:", err)
		return nil, fmt.Errorf("failed to get ToAccount: %w", err)
	}
	if len(toResp.Accounts) == 0 {
		return nil, fmt.Errorf("no account found for ToAccountId: %v", req.ToAccountId)
	}
	toAcc := toResp.Accounts[0]

	if fromAcc.AvailableBalance < int64(req.Amount) {
		return nil, fmt.Errorf("insufficient funds in the FromAccount")
	}

	
	newFromBalance := fromAcc.AvailableBalance - int64(req.Amount)
	newToBalance := toAcc.AvailableBalance + int64(req.Amount)

	_, err = s.AccountClient.UpdateAccount(ctx, &accountpb.UpdateRequest{
		CustomerId:        fromAcc.CustomerId,
		AccountId:         fromAcc.AccountId,
		AvailableBalance:  newFromBalance,
		Status:            fromAcc.Status,
		AccountNumber: fromAcc.AccountNumber,
		AccountType: fromAcc.AccountType,
		PendingBalance: fromAcc.PendingBalance,
		CreditLimit: fromAcc.CreditLimit,
		Currency: fromAcc.Currency,

		
		
	})
	if err != nil {
		log.Println("Failed to update From Account:", err)
		return nil, fmt.Errorf("failed to update FromAccount: %w", err)
	}

	// Update the To Account balance
	_, err = s.AccountClient.UpdateAccount(ctx, &accountpb.UpdateRequest{
		CustomerId:        toAcc.CustomerId,
		AccountId:         toAcc.AccountId,
		AvailableBalance:  newToBalance,
		Status:            toAcc.Status,
		AccountNumber: toAcc.AccountNumber,
		AccountType: toAcc.AccountType,
		Currency: toAcc.Currency,
		CreditLimit: toAcc.CreditLimit,
		PendingBalance: toAcc.PendingBalance,
		
	})
	if err != nil {
		log.Println("Failed to update To Account:", err)
		return nil, fmt.Errorf("failed to update ToAccount: %w", err)
	}

	// Fetch Customer Information for both sender and receiver
	// fromCus, err := s.CustomerClient.DisplayCustomer(ctx, &cu.FindCustomerByIdRequest{CustomerId: fromAcc.CustomerId})
	// if err != nil {
	// 	log.Println("Failed to get From Customer:", err)
	// 	return nil, fmt.Errorf("failed to fetch FromCustomer: %w", err)
	// }

	// toCus, err := s.CustomerClient.DisplayCustomer(ctx, &cu.FindCustomerByIdRequest{CustomerId: toAcc.CustomerId})
	// if err != nil {
	// 	log.Println("Failed to get To Customer:", err)
	// 	return nil, fmt.Errorf("failed to fetch ToCustomer: %w", err)
	// }

	// Build the notification
	// notification := PaymentStatus{
	// 	FromEmail:     fromCus.Display[0].Email,
	// 	ToEmail:       toCus.Display[0].Email,
	// 	FromBalance:   newFromBalance,
	// 	ToBalance:     newToBalance,
	// 	FromFirstName: fromCus.Display[0].FirstName,
	// 	ToFirstName:   toCus.Display[0].FirstName,
	// 	Amount:        int64(req.Amount),
	// 	AccountNumber: req.FromAccountId,
	// }

	// Send the notification to NATS
	nc, err := connectToNats()
	if err != nil {
		log.Println("Failed to connect to NATS:", err)
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	defer nc.Close()

	// data, err := json.Marshal(notification)
	// if err != nil {
	// 	log.Println("Failed to marshal notification:", err)
	// 	return nil, fmt.Errorf("failed to marshal notification: %w", err)
	// }

	// if err := nc.Publish("PaymentNotification", data); err != nil {
	// 	log.Println("Failed to publish to NATS:", err)
	// 	return nil, fmt.Errorf("failed to publish notification to NATS: %w", err)
	// }

	// Return the payment details as the response


	p, err := s.uc.CreatePayment(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	return &v1.CreatePaymentReply{
		PaymentId:         p.PaymentId,
		FromAccountId:     p.FromAccountId,
		ToAccountId:       p.ToAccountId,
		PaymentType:       p.PaymentType,
		Amount:            float64(p.Amount),
		Currency:          p.Currency,
		Status:            p.Status,
		PaymentMethod:     p.PaymentMethod,
		ReferenceNumber:   p.ReferenceNumber,
		ExternalReference: p.ExternalReference,
		Message:           "Payment created successfully",
	}, nil
}
