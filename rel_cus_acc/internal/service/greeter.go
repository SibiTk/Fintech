package service

import (
    "context"
    "fmt"
    "strconv"
    "time"
    
    v1 "rel_cus_acc/api/helloworld/v1"
    "rel_cus_acc/internal/biz"
    
   
    accountpb "account/api/helloworld/v1"
    customerpb "customer/api/helloworld/v1"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type RelCusAccService struct {
    v1.UnimplementedRelCusAccServer
    uc *biz.RelCusAccUsecase
}

func NewRelCusAccService(uc *biz.RelCusAccUsecase) *RelCusAccService {
    return &RelCusAccService{uc: uc}
}

// connect to Customer Service
func (s *RelCusAccService) connectToCustomerService() (customerpb.CustomerManagerClient, *grpc.ClientConn, error) {
    conn, err := grpc.Dial("localhost:9012", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect to CustomerService: %w", err)
    }
    client := customerpb.NewCustomerManagerClient(conn)
    return client, conn, nil
}

// connect to account service
func (s *RelCusAccService) connectToAccountService() (accountpb.AccountClient, *grpc.ClientConn, error) {
    conn, err := grpc.Dial("localhost:9013", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect to AccountService: %w", err)
    }
    client := accountpb.NewAccountClient(conn)
    return client, conn, nil
}

//  customer-account relationship by fetching data from both services
func (s *RelCusAccService) CreateRelation(ctx context.Context, req *v1.CreateRelationRequest) (*v1.CreateRelationReply, error) {
   
    customerClient, customerConn, err := s.connectToCustomerService()
    if err != nil {
        return nil, err
    }
    defer customerConn.Close()

    // Connect to account service
    accountClient, accountConn, err := s.connectToAccountService()
    if err != nil {
        return nil, err
    }
    defer accountConn.Close()

    // Fetch customer details
    customerResp, err := customerClient.DisplayCustomer(ctx, &customerpb.FindCustomerByIdRequest{
        CustomerId: req.CustomerId,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get customer details: %w", err)
    }

    if len(customerResp.Display) == 0 {
        return nil, fmt.Errorf("customer not found: %d", req.CustomerId)
    }

    customer := customerResp.Display[0]

    // Fetch account details
    accountResp, err := accountClient.GetAccountWithId(ctx, &accountpb.AccountIdRequest{
        AccountId: req.AccountId,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get account details: %w", err)
    }

    if len(accountResp.Accounts) == 0 {
        return nil, fmt.Errorf("account not found: %d", req.AccountId)
    }

    account := accountResp.Accounts[0]

    // Create relationship with combined data
    currentTime := time.Now().Format("2006-01-02 15:04:05")
    relation := &biz.RelCusAcc{
        CustomerId:          customer.CustomerId,
        CustomerNumber:      customer.CustomerNumber,
        FirstName:           customer.FirstName,
        LastName:            customer.LastName,
        Email:               customer.Email,
        Phone:               customer.Phone,
        DateOfBirth:         customer.DateOfBirth,
        CustomerStatus:      customer.Status,
        KycStatus:           customer.KycStatus,
        CustomerCreatedAt:   currentTime,
        CustomerUpdatedAt:   currentTime,
        AccountId:           account.AccountId,
        AccountNumber:       strconv.FormatInt(account.AccountNumber, 10),
        AccountType:         account.AccountType,
        Currency:            account.Currency,
        AccountStatus:       account.Status,
        AvailableBalance:    account.AvailableBalance,
        PendingBalance:      account.PendingBalance,
        CreditLimit:         account.CreditLimit,
        LastTransactionAt:   account.LastTransactionAt,
        AccountCreatedAt:    currentTime,
        LastUsedAt:          currentTime,
    }

    createdRelation, err := s.uc.CreateRelation(ctx, relation)
    if err != nil {
        return nil, err
    }

    return &v1.CreateRelationReply{
        Id:                  createdRelation.Id,
        CustomerId:          createdRelation.CustomerId,
        CustomerNumber:      createdRelation.CustomerNumber,
        FirstName:           createdRelation.FirstName,
        LastName:            createdRelation.LastName,
        Email:               createdRelation.Email,
        Phone:               createdRelation.Phone,
        DateOfBirth:         createdRelation.DateOfBirth,
        CustomerStatus:      createdRelation.CustomerStatus,
        KycStatus:           createdRelation.KycStatus,
        CustomerCreatedAt:   createdRelation.CustomerCreatedAt,
        CustomerUpdatedAt:   createdRelation.CustomerUpdatedAt,
        AccountId:           createdRelation.AccountId,
        AccountNumber:       createdRelation.AccountNumber,
        AccountType:         createdRelation.AccountType,
        Currency:            createdRelation.Currency,
        AccountStatus:       createdRelation.AccountStatus,
        AvailableBalance:    createdRelation.AvailableBalance,
        PendingBalance:      createdRelation.PendingBalance,
        CreditLimit:         createdRelation.CreditLimit,
        LastTransactionAt:   createdRelation.LastTransactionAt,
        AccountCreatedAt:    createdRelation.AccountCreatedAt,
        LastUsedAt:          createdRelation.LastUsedAt,
        Message:             "Relationship created successfully",
    }, nil
}

// UpdateRelation updates existing relationship with fresh data from both services
func (s *RelCusAccService) UpdateRelation(ctx context.Context, req *v1.UpdateRelationRequest) (*v1.UpdateRelationReply, error) {
    // Connect to Customer Service
    customerClient, customerConn, err := s.connectToCustomerService()
    if err != nil {
        return nil, err
    }
    defer customerConn.Close()

    // Connect to Account Service
    accountClient, accountConn, err := s.connectToAccountService()
    if err != nil {
        return nil, err
    }
    defer accountConn.Close()

    // Fetch updated customer details
    customerResp, err := customerClient.DisplayCustomer(ctx, &customerpb.FindCustomerByIdRequest{
        CustomerId: req.CustomerId,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get customer details: %w", err)
    }

    if len(customerResp.Display) == 0 {
        return nil, fmt.Errorf("customer not found: %d", req.CustomerId)
    }

    customer := customerResp.Display[0]

    // Fetch updated account details
    accountResp, err := accountClient.GetAccountWithId(ctx, &accountpb.AccountIdRequest{
        AccountId: req.AccountId,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to get account details: %w", err)
    }

    if len(accountResp.Accounts) == 0 {
        return nil, fmt.Errorf("account not found: %d", req.AccountId)
    }

    account := accountResp.Accounts[0]

    // Update relationship with fresh data
    relation := &biz.RelCusAcc{
        Id:                  req.Id,
        CustomerId:          customer.CustomerId,
        CustomerNumber:      customer.CustomerNumber,
        FirstName:           customer.FirstName,
        LastName:            customer.LastName,
        Email:               customer.Email,
        Phone:               customer.Phone,
        DateOfBirth:         customer.DateOfBirth,
        CustomerStatus:      customer.Status,
        KycStatus:           customer.KycStatus,
        CustomerUpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
        AccountId:           account.AccountId,
        AccountNumber:       strconv.FormatInt(account.AccountNumber, 10),
        AccountType:         account.AccountType,
        Currency:            account.Currency,
        AccountStatus:       account.Status,
        AvailableBalance:    account.AvailableBalance,
        PendingBalance:      account.PendingBalance,
        CreditLimit:         account.CreditLimit,
        LastTransactionAt:   account.LastTransactionAt,
        LastUsedAt:          time.Now().Format("2006-01-02 15:04:05"),
    }

    updatedRelation, err := s.uc.UpdateRelation(ctx, relation)
    if err != nil {
        return nil, err
    }

    relationReply := &v1.RelationReply{
        Id:                  updatedRelation.Id,
        CustomerId:          updatedRelation.CustomerId,
        CustomerNumber:      updatedRelation.CustomerNumber,
        FirstName:           updatedRelation.FirstName,
        LastName:            updatedRelation.LastName,
        Email:               updatedRelation.Email,
        Phone:               updatedRelation.Phone,
        DateOfBirth:         updatedRelation.DateOfBirth,
        CustomerStatus:      updatedRelation.CustomerStatus,
        KycStatus:           updatedRelation.KycStatus,
        CustomerCreatedAt:   updatedRelation.CustomerCreatedAt,
        CustomerUpdatedAt:   updatedRelation.CustomerUpdatedAt,
        AccountId:           updatedRelation.AccountId,
        AccountNumber:       updatedRelation.AccountNumber,
        AccountType:         updatedRelation.AccountType,
        Currency:            updatedRelation.Currency,
        AccountStatus:       updatedRelation.AccountStatus,
        AvailableBalance:    updatedRelation.AvailableBalance,
        PendingBalance:      updatedRelation.PendingBalance,
        CreditLimit:         updatedRelation.CreditLimit,
        LastTransactionAt:   updatedRelation.LastTransactionAt,
        AccountCreatedAt:    updatedRelation.AccountCreatedAt,
        LastUsedAt:          updatedRelation.LastUsedAt,
    }

    return &v1.UpdateRelationReply{
        Message:   "Relationship updated successfully",
        Relations: []*v1.RelationReply{relationReply},
    }, nil
}

// DeleteRelation deletes a relationship by ID
func (s *RelCusAccService) DeleteRelation(ctx context.Context, req *v1.DeleteRelationRequest) (*v1.DeleteRelationReply, error) {
    err := s.uc.DeleteRelation(ctx, req.Id)
    if err != nil {
        return nil, err
    }

    return &v1.DeleteRelationReply{
        Message: "Relationship deleted successfully",
    }, nil
}

// GetRelationsByCustomer returns all relationships for a customer
func (s *RelCusAccService) GetRelationsByCustomer(ctx context.Context, req *v1.GetRelationsByCustomerRequest) (*v1.GetRelationsByCustomerReply, error) {
    relations, err := s.uc.GetRelationsByCustomer(ctx, req.CustomerId)
    if err != nil {
        return nil, err
    }

    var relationReplies []*v1.RelationReply
    for _, rel := range relations {
        relationReplies = append(relationReplies, &v1.RelationReply{
            Id:                  rel.Id,
            CustomerId:          rel.CustomerId,
            CustomerNumber:      rel.CustomerNumber,
            FirstName:           rel.FirstName,
            LastName:            rel.LastName,
            Email:               rel.Email,
            Phone:               rel.Phone,
            DateOfBirth:         rel.DateOfBirth,
            CustomerStatus:      rel.CustomerStatus,
            KycStatus:           rel.KycStatus,
            CustomerCreatedAt:   rel.CustomerCreatedAt,
            CustomerUpdatedAt:   rel.CustomerUpdatedAt,
            AccountId:           rel.AccountId,
            AccountNumber:       rel.AccountNumber,
            AccountType:         rel.AccountType,
            Currency:            rel.Currency,
            AccountStatus:       rel.AccountStatus,
            AvailableBalance:    rel.AvailableBalance,
            PendingBalance:      rel.PendingBalance,
            CreditLimit:         rel.CreditLimit,
            LastTransactionAt:   rel.LastTransactionAt,
            AccountCreatedAt:    rel.AccountCreatedAt,
            LastUsedAt:          rel.LastUsedAt,
        })
    }

    return &v1.GetRelationsByCustomerReply{
        Relations: relationReplies,
    }, nil
}

// GetRelationByAccount returns relationship for a specific account
func (s *RelCusAccService) GetRelationByAccount(ctx context.Context, req *v1.GetRelationByAccountRequest) (*v1.GetRelationByAccountReply, error) {
    relation, err := s.uc.GetRelationByAccount(ctx, req.AccountId)
    if err != nil {
        return nil, err
    }

    if relation == nil {
        return &v1.GetRelationByAccountReply{
            Relation: nil,
        }, nil
    }

    relationReply := &v1.RelationReply{
        Id:                  relation.Id,
        CustomerId:          relation.CustomerId,
        CustomerNumber:      relation.CustomerNumber,
        FirstName:           relation.FirstName,
        LastName:            relation.LastName,
        Email:               relation.Email,
        Phone:               relation.Phone,
        DateOfBirth:         relation.DateOfBirth,
        CustomerStatus:      relation.CustomerStatus,
        KycStatus:           relation.KycStatus,
        CustomerCreatedAt:   relation.CustomerCreatedAt,
        CustomerUpdatedAt:   relation.CustomerUpdatedAt,
        AccountId:           relation.AccountId,
        AccountNumber:       relation.AccountNumber,
        AccountType:         relation.AccountType,
        Currency:            relation.Currency,
        AccountStatus:       relation.AccountStatus,
        AvailableBalance:    relation.AvailableBalance,
        PendingBalance:      relation.PendingBalance,
        CreditLimit:         relation.CreditLimit,
        LastTransactionAt:   relation.LastTransactionAt,
        AccountCreatedAt:    relation.AccountCreatedAt,
        LastUsedAt:          relation.LastUsedAt,
    }

    return &v1.GetRelationByAccountReply{
        Relation: relationReply,
    }, nil
}

