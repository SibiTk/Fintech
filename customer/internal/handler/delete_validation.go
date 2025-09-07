package handler

import (
    "context"
    "fmt"

    accountpb "account/api/helloworld/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

// DeleteCustomerWithCascade deletes customer accounts via gRPC 
func DeleteCustomerWithCascade(ctx context.Context, customerId int64, deleteCustomerFunc func(context.Context, int64) error) error {
  
    conn, err := grpc.Dial("localhost:9013", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return fmt.Errorf("failed to connect to AccountService: %w", err)
    }
    defer conn.Close()

    accountClient := accountpb.NewAccountClient(conn)

  
    _, err = accountClient.DeleteAccount(ctx, &accountpb.DeleteRequest{
        CustomerId: customerId,
    })
    if err != nil {
        return fmt.Errorf("failed to delete accounts for customer %d: %w", customerId, err)
    }

    // Delete customer record
    err = deleteCustomerFunc(ctx, customerId)
    if err != nil {
        return fmt.Errorf("failed to delete customer %d: %w", customerId, err)
    }

    return nil
}
