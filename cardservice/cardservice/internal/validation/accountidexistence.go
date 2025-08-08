package validation



import (
    "context"
    "fmt"

    accountpb "account/api/helloworld/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func CheckAccountExistsViaGRPC(accountID int64) (bool, error) {
    conn, err := grpc.NewClient("localhost:9013", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return false, fmt.Errorf("failed to connect to customer service: %w", err)
    }
    defer conn.Close()

    client := accountpb.NewAccountClient(conn)

    resp, err := client.GetAccountWithId(context.Background(), &accountpb.AccountIdRequest{
        AccountId: accountID,
    })
    if err != nil {
        return false, fmt.Errorf("error calling AccountID or Account Id does not exist: %w", err)
    }

    // Determine existence based on Display slice length
    exists := len(resp.Accounts ) > 0
    return exists, nil
}