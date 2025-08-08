package client

import (
	accountpb "account/api/helloworld/v1"
	"context"
	"log"

	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(AccountClient)

func AccountClient() accountpb.AccountClient {
	con, err := grpc.DialInsecure(context.Background(),
		grpc.WithEndpoint("localhost:9013"))
	if err != nil {
		log.Println("CLIENT CONNECTIOn FOR Account FAILED ", err)
	}
	return accountpb.NewAccountClient(con)
}
