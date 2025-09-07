package mapper

import (
    kerrors "github.com/go-kratos/kratos/v2/errors"
    "google.golang.org/protobuf/types/known/emptypb"

    pb "cache_service/api/helloworld/v1"
)


func CreateCacheFailureResponse(_ *pb.CacheRequest, err *kerrors.Error) (*pb.CacheResponse, error) {
    return nil, err
}


func CreateCacheResponse(req *pb.CacheRequest, _  *kerrors.Error) *pb.CacheResponse {
    if req == nil {
        return &pb.CacheResponse{}
    }
    return &pb.CacheResponse{
        AccountNumber: req.AccountNumber,      
        AvailableBalance: req.AvailableBalance, 
    }
}


func CreateCacheEmptyResponse(reason *kerrors.Error) (*emptypb.Empty, error) {

    if reason == nil || reason.Reason == "SUCCESS" {
        return &emptypb.Empty{}, nil
    }
    return nil, reason
}
