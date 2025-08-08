package service

import (
	"context"



	pb "BlockAccount/api/helloworld/v1"
	"BlockAccount/internal/handler"
	accountpb "account/api/helloworld/v1"
)


type AccountBlockServiceService struct {
	AccountClient  accountpb.AccountClient
	pb.UnimplementedAccountBlockServiceServer
	accountblock *handler.AccountBlockHandler
}

func NewAccountBlockServiceService(accountblock *handler.AccountBlockHandler) *AccountBlockServiceService {
	return &AccountBlockServiceService{accountblock: accountblock}
}

func (s *AccountBlockServiceService) SaveAccBlock(ctx context.Context, req *pb.SaveAccBlockRequest) (*pb.SaveAccBlockReply, error) {
	return s.accountblock.AccountBlock(ctx, req)
}

func (s *AccountBlockServiceService) GetAccBlock(ctx context.Context, req *pb.GetAccBlockRequest) (*pb.GetAccBlockReply, error) {
	return s.accountblock.GetAccBlock(ctx, req)

}
func (s *AccountBlockServiceService) UpdateAccBlock(ctx context.Context, req *pb.UpdateAccBlockRequest) (*pb.UpdateAccBlockReply, error) {
	return s.accountblock.UpdateAccBlock(ctx, req)
}
