package handler

import (
	"context"
	"fmt"
	

	"github.com/go-kratos/kratos/v2/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "cache_service/api/helloworld/v1"
	"cache_service/internal/mapper"
)

func (h *CacheHandler) SaveCache(ctx context.Context, req *pb.CacheRequest) (*emptypb.Empty, error) {

    key := fmt.Sprintf("balance:%d", req.AccountNumber)

   
    if err := h.redis.HSet(
        ctx,
        key,
        "availableBalance", req.AvailableBalance,
        "currency", req.Currency,
    ).Err(); err != nil {
        zap.S().Errorf("SaveCache HSET key [%s] error [%v]", key, err)
        return mapper.CreateCacheEmptyResponse(errors.New(500, "redis_hset_error", "failed to save cache"))
    }

   

    return mapper.CreateCacheEmptyResponse(errors.New(200, "SUCCESS", ""))
}
