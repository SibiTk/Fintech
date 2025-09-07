package service

import (
	"context"
	"encoding/json"
	"fmt"

	accountpb "account/api/helloworld/v1"
	v1 "customer/api/helloworld/v1"
	"customer/internal/biz"
	"customer/internal/handler"


	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	//"customermanager/internal/validation"
)

type CustomerManagerService struct {
	v1.UnimplementedCustomerManagerServer
	uc *biz.CustomerUsecase
}

type Message1 struct{
	CustomerNumber string
	FirstName string
	Email string `json:"email"`
	status string

}
func connectToNATS() (*nats.Conn, error) {
	fmt.Println("nats is Connected Successfully")
	return nats.Connect(nats.DefaultURL)
}


func NewCustomerManagerService(uc *biz.CustomerUsecase) *CustomerManagerService {
	return &CustomerManagerService{uc: uc}
}

func (s *CustomerManagerService) CreateCustomer(ctx context.Context, req *v1.CreateCustomerRequest) (*v1.CreateCustomerReply, error) {

	if err := handler.ValidateCustomerFields(
		req.CustomerNumber,
		req.FirstName,
		req.LastName,
		req.DateOfBirth,
		req.Status,
		req.KycStatus,
		
	); err != nil {
		return &v1.CreateCustomerReply{
			Message: "Validation failed FOR Customers: " + err.Error(),
		}, nil
	}
   if err := handler.ValidateEmail(
	req.Email,
   ); err != nil {
		return &v1.CreateCustomerReply{
			Message: "Validation failed FOR EMAIL: " + err.Error(),
		}, nil
	}
	if err:=handler.ValidatePhoneNumber(req.Phone,); err != nil {
		return &v1.CreateCustomerReply{
			Message: "Validation failed for PHONE NUMBER MUST BE IN 10 DIGIT: " + err.Error(),
		}, nil
	}
	nc,_:=connectToNATS()
	fmt.Println("Nats Connecting............")
	customer := &biz.Customer{
		CustomerNumber: req.CustomerNumber,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		DateOfBirth:    req.DateOfBirth,
		Status:         req.Status,
		KycStatus:      req.KycStatus,
	}

	c, err := s.uc.CreateCustomer(ctx, customer)
	if err != nil {
		return &v1.CreateCustomerReply{
			Message: "Failed to create customer: " + err.Error(),
		}, nil
	}


	result := Message1{
		CustomerNumber: req.CustomerNumber,
		FirstName:req.FirstName,
		Email:req.Email,
		status:req.Status,

	}
	data,err := json.Marshal(result)
	if err!=nil{
		return nil,fmt.Errorf("Failed:%w",err)
	}
	// nc,_:=connectToNATS()
	nc.Publish("Create",data)
	fmt.Println("EMail")
	return &v1.CreateCustomerReply{
		Message:        "Customer Created Successfully",
		CustomerId:     c.CustomerId,
		CustomerNumber: c.CustomerNumber,
		FirstName:      c.FirstName,
		LastName:       c.LastName,
		Email:          c.Email,
		Phone:          c.Phone,
		DateOfBirth:    c.DateOfBirth,
		Status:         c.Status,
		KycStatus:      c.KycStatus,
	}, nil
}

func (s *CustomerManagerService) UpdateCustomer(ctx context.Context, req *v1.UpdateCustomerRequest) (*v1.UpdateCustomerReply, error) {
	customer := &biz.Customer{
		CustomerId:     req.CustomerId,
		CustomerNumber: req.CustomerNumber,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Email:          req.Email,
		Phone:          req.Phone,
		DateOfBirth:    req.DateOfBirth,
		Status:         req.Status,
		KycStatus:      req.KycStatus,
	}

	updatedCustomer, err := s.uc.UpdateCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateCustomerReply{
		Message: "Customer Details are updated",
		Display: []*v1.CreateCustomerReply{
			{
				CustomerId:     updatedCustomer.CustomerId,
				CustomerNumber: updatedCustomer.CustomerNumber,
				FirstName:      updatedCustomer.FirstName,
				LastName:       updatedCustomer.LastName,
				Email:          updatedCustomer.Email,
				Phone:          updatedCustomer.Phone,
				DateOfBirth:    updatedCustomer.DateOfBirth,
				Status:         updatedCustomer.Status,
				KycStatus:      updatedCustomer.KycStatus,
			},
		},
	}, nil
}

func (s *CustomerManagerService) DeleteCustomer(ctx context.Context, req *v1.DeleteCustomerRequest) (*v1.DeleteCustomerReply, error) {
    err := handler.DeleteCustomerWithCascade(ctx, req.CustomerId, s.uc.DeleteCustomer)
    if err != nil {
        return &v1.DeleteCustomerReply{
            Message: err.Error(),
        }, nil
    }

    return &v1.DeleteCustomerReply{
        Message: "Customer and all associated accounts deleted successfully",
    }, nil
}


func (s *CustomerManagerService) DisplayCustomer(ctx context.Context, req *v1.FindCustomerByIdRequest) (*v1.FindCustomerByIdReply, error) {
    // Input validation
    if req.CustomerId == 0 {
        return nil, status.Error(codes.InvalidArgument, "customer ID is required")
    }

    cust, err := s.uc.DisplayCustomer(ctx, req.CustomerId)
    if err != nil {
        
        return nil, status.Errorf(codes.Internal, "failed to fetch customer: %v", err)
    }

    conn, err := grpc.Dial("localhost:9013", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to connect to AccountService: %v", err)
    }
    defer conn.Close()

    accountClient := accountpb.NewAccountClient(conn)
    accountResp, err := accountClient.GetCustomerWithId(ctx, &accountpb.GetCustomerWithIdRequest{
        CustomerId: req.CustomerId,
    })
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to fetch accounts for customer: %v", err)
    }

    var accountList []*v1.AccountDetails
    for _, acc := range accountResp.Accounts {
        accountList = append(accountList, &v1.AccountDetails{
            AccountId:         acc.AccountId,
            CustomerId:        acc.CustomerId,
            AccountNumber:     acc.AccountNumber,
            AccountType:       acc.AccountType,
            Currency:          acc.Currency,
            Status:            acc.Status,
            AvailableBalance:  acc.AvailableBalance,
            PendingBalance:    acc.PendingBalance,
            CreditLimit:       acc.CreditLimit,
            LastTransactionAt: acc.LastTransactionAt,
        })
    }

    return &v1.FindCustomerByIdReply{
        Display: []*v1.CreateCustomerReply{
            {
                CustomerId:     cust.CustomerId,
                CustomerNumber: cust.CustomerNumber,
                FirstName:      cust.FirstName,
                LastName:       cust.LastName,
                Email:          cust.Email,
                Phone:          cust.Phone,
                DateOfBirth:    cust.DateOfBirth,
                Status:         cust.Status,
                KycStatus:      cust.KycStatus,
            },
        },
        Acc: accountList,
    }, nil
}
