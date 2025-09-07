package handler

import (
    "context"
    "fmt"

    kerrors "github.com/go-kratos/kratos/v2/errors"
    "go.uber.org/zap"
    "google.golang.org/protobuf/types/known/emptypb"

    "cache_service/internal/mapper"
    pb "cache_service/api/helloworld/v1"
)

func balanceKey(accountID int64) string {
    return fmt.Sprintf("balance:%d", accountID)
}

func (h *CacheHandler) DeleteCache(ctx context.Context, req *pb.CacheRequest) (*emptypb.Empty, error) {
    key := balanceKey(req.AccountNumber)

    if err := h.redis.Del(ctx, key).Err(); err != nil {
        zap.S().Errorf("DeleteCache DEL key [%s] error [%v]", key, err)
        return mapper.CreateCacheEmptyResponse(kerrors.New(500, "redis_del_error", "failed to delete cache"))
    }

    return mapper.CreateCacheEmptyResponse(kerrors.New(200, "SUCCESS", ""))
}
