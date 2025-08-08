package handler

import (
    "context"
    "fmt"

    customerpb "customer/api/helloworld/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func CheckCustomerExistsViaGRPC(customerID int64) (bool, error) {
    conn, err := grpc.Dial("localhost:9022", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return false, fmt.Errorf("failed to connect to customer service: %w", err)
    }
    defer conn.Close()

    client := customerpb.NewCustomerManagerClient(conn)

    resp, err := client.DisplayCustomer(context.Background(), &customerpb.FindCustomerByIdRequest{
        CustomerId: customerID,
    })
    if err != nil {
        return false, fmt.Errorf("error calling DisplayCustomer: %w", err)
    }

    // Determine existence based on Display slice length
    exists := len(resp.Display) > 0
    return exists, nil
}

