package service

import (
	v1 "account/api/helloworld/v1"
	"account/internal/biz"
	"account/internal/handler"
	"context"
	"fmt"
	//"strconv"
)

type AccountService struct {
	v1.UnimplementedAccountServer
	uc *biz.AccountUsecase
}

func NewAccountService(uc *biz.AccountUsecase) *AccountService {
	return &AccountService{uc: uc}
}

func (s *AccountService) CreateAccount(ctx context.Context, req *v1.AccountRequest) (*v1.AccountReply, error) {
  
    if err := handler.ValidateAccountFields(
        req.AccountType,
        req.Currency,
        req.Status,
        req.AvailableBalance,
        req.PendingBalance,
    ); err != nil {
        return &v1.AccountReply{
            Message: "Validation failed for Account Fields: " + err.Error(),
        }, nil
    }

 
   
    accountNumber, err := handler.GenerateAccountNumber()
    if err != nil {
        return &v1.AccountReply{
            Message: "Failed to generate account number: " + err.Error(),
        }, nil
    }
    fmt.Println("Generated account number:", accountNumber)
	

exists, err := handler.CheckCustomerExistsViaGRPC(req.CustomerId)
if err != nil {
    return &v1.AccountReply{
        Message: "Failed to verify customer existence: " + err.Error(),
    }, nil
}
if !exists {
    return &v1.AccountReply{
        Message: "Customer does not exist",
    }, nil
}
   
    acc := &biz.Account{
        CustomerId:        req.CustomerId,
        AccountNumber:     accountNumber,
        AccountType:       req.AccountType,
        Currency:          req.Currency,
        Status:            req.Status,
        AvailableBalance:  req.AvailableBalance,
        PendingBalance:    req.PendingBalance,
        CreditLimit:       req.CreditLimit,
        LastTransactionAt: req.LastTransactionAt,
    }

   
    a, err := s.uc.Create(ctx, acc)
    if err != nil {
        return &v1.AccountReply{
            Message: "Failed to create account: " + err.Error(),
        }, nil
    }

 
    return &v1.AccountReply{
        AccountId:         a.AccountId,
        CustomerId:        a.CustomerId,
        AccountNumber:     a.AccountNumber,
        AccountType:       a.AccountType,
        Currency:          a.Currency,
        Status:            a.Status,
        AvailableBalance:  a.AvailableBalance,
        PendingBalance:    a.PendingBalance,
        CreditLimit:       a.CreditLimit,
        LastTransactionAt: a.LastTransactionAt,
        Message:           "Account created successfully",
    }, nil
}

// UpdateAccount updates account type and status only (as per proto)
func (s *AccountService) UpdateAccount(ctx context.Context, req *v1.UpdateRequest) (*v1.UpdateReply, error) {
	acc := &biz.Account{
		CustomerId:  req.CustomerId,
		AccountType: req.AccountType,
		Status:      req.Status,
		AccountId: req.AccountId,
		AccountNumber: req.AccountNumber,
		Currency: req.Currency,
		AvailableBalance: req.AvailableBalance,
		PendingBalance: req.PendingBalance,
		CreditLimit: req.CreditLimit,
		
	}
	a, err := s.uc.Update(ctx, acc)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateReply{
		AccountId:     a.AccountId,
		CustomerId:    a.CustomerId,
		AccountNumber: a.AccountNumber,
		AccountType:   a.AccountType,
		Currency:      a.Currency,
		Status:        a.Status,
		AvailableBalance: a.AvailableBalance,
		PendingBalance: a.PendingBalance,
		CreditLimit: a.CreditLimit,
		LastTransactionAt: a.LastTransactionAt,
		Message:           "Fetched successfully",

	}, nil
}

// DeleteAccount deletes an account by customer ID
func (s *AccountService) DeleteAccount(ctx context.Context, req *v1.DeleteRequest) (*v1.DeleteReply, error) {
	err := s.uc.Delete(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteReply{Message: "Account deleted"}, nil
}

// GetCustomerWithId returns all accounts linked to a given customer ID
func (s *AccountService) GetCustomerWithId(ctx context.Context, req *v1.GetCustomerWithIdRequest) (*v1.GetCustomerWithIdReply, error) {
	accounts, err := s.uc.GetByCustomerID(ctx, req.CustomerId)
	if err != nil {
		return nil, err
	}

	var accountList []*v1.AccountReply
	for _, a := range accounts {
		accountList = append(accountList, &v1.AccountReply{
			AccountId:         a.AccountId,
			CustomerId:        a.CustomerId,
			AccountNumber:     a.AccountNumber,
			AccountType:       a.AccountType,
			Currency:          a.Currency,
			Status:            a.Status,
			AvailableBalance:  a.AvailableBalance,
			PendingBalance:    a.PendingBalance,
			CreditLimit:       a.CreditLimit,
			LastTransactionAt: a.LastTransactionAt,
			Message:           "Fetched successfully",
		})
	}

	return &v1.GetCustomerWithIdReply{
		Accounts: accountList,
	}, nil
}
func (s *AccountService) GetAccountWithId(ctx context.Context, req *v1.AccountIdRequest) (*v1.AccountIdReply, error) {
    accounts, err := s.uc.GetByAccountID(ctx, req.AccountId)
    if err != nil {
        return nil, err
    }

    var accountReplies []*v1.AccountReply
    for _, a := range accounts {
        accountReplies = append(accountReplies, &v1.AccountReply{
            AccountId:         a.AccountId,
            CustomerId:        a.CustomerId,
            AccountNumber:     a.AccountNumber,
            AccountType:       a.AccountType,
            Currency:          a.Currency,
            Status:            a.Status,
            AvailableBalance:  a.AvailableBalance,
            PendingBalance:    a.PendingBalance,
            CreditLimit:       a.CreditLimit,
            LastTransactionAt: a.LastTransactionAt,
            Message:           "Fetched successfully",
        })
    }

    return &v1.AccountIdReply{
        Accounts: accountReplies,
    }, nil
}