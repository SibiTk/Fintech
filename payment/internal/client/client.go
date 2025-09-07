
package client
 
import (
    "context"
    "log"
    
   // v1 "payment/api/helloworld/v1"
	accountpb "account/api/helloworld/v1"
  cu "customer/api/helloworld/v1"
   tra "transaction/api/hello/v1"
    "github.com/go-kratos/kratos/v2/transport/grpc"
   // "github.com/google/wire"
)
 

 
//var ProviderSet = wire.NewSet(AccountClient)

 
func AccountClient() accountpb.AccountClient {
    con, err := grpc.DialInsecure(context.Background(),
        grpc.WithEndpoint("localhost:9013"))
    if err != nil {
        log.Println("CLIENT CONNECTIOn FOR Account FAILED ", err)
    }
    return accountpb.NewAccountClient(con)
}
func CustomerClient() cu.CustomerManagerClient {
    con, err := grpc.DialInsecure(context.Background(),
        grpc.WithEndpoint("localhost:9022"))
    if err != nil {
        log.Println("CLIENT CONNECTION FOR CUSTOMER FAILED ", err)
    }
    return cu.NewCustomerManagerClient(con)
}
 func TransactionnClient() tra.TransactionClient{
    con, err := grpc.DialInsecure(context.Background(),
        grpc.WithEndpoint("localhost:9029"))
    if err != nil {
        log.Println("CLIENT CONNECTION FOR TRANSACTION FAILED ", err)
    }
    return tra.NewTransactionClient(con)
}
 
 
 
	  