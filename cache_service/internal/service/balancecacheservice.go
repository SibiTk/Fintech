package service

import (
	"context"
	//"fmt"

	pb "cache_service/api/helloworld/v1"
	"cache_service/internal/handler"

	"google.golang.org/protobuf/types/known/emptypb"
)

type BalanceCacheServiceService struct {
	pb.UnimplementedBalanceCacheServiceServer
	CacheHandler *handler.CacheHandler
}

func NewBalanceCacheServiceService(CacheHandler *handler.CacheHandler) *BalanceCacheServiceService {
	return &BalanceCacheServiceService{CacheHandler: CacheHandler}
}

func (s *BalanceCacheServiceService) GetBalance(ctx context.Context, req *pb.CacheRequest) (*pb.CacheResponse, error) {
	return s.CacheHandler.GetBalance(ctx,req)
}
func (s *BalanceCacheServiceService) SetBalance(ctx context.Context, req *pb.CacheRequest) (*emptypb.Empty, error) {
	
	return s.CacheHandler.SaveCache(ctx,req)
}
func (s *BalanceCacheServiceService) DeleteBalance(ctx context.Context, req *pb.CacheRequest) (*emptypb.Empty, error) {
		return s.CacheHandler.DeleteCache(ctx,req)
}
