package handler

import (
    "context"
   // "fmt"
    "strconv"

    kerrors "github.com/go-kratos/kratos/v2/errors"
    "go.uber.org/zap"
    "google.golang.org/protobuf/types/known/emptypb"

    pb "cache_service/api/helloworld/v1"
)

// func balanceKey(accountID int64) string {
//     return fmt.Sprintf("balance:%d", accountID)
// }

// SetBalance: writes the balance fields to a Redis hash and returns Empty.
func (h *CacheHandler) SetBalance(ctx context.Context, req *pb.CacheRequest) (*emptypb.Empty, error) {
    key := balanceKey(req.AccountNumber)

    // HSET balance:{id} availableBalance currency
    if err := h.redis.HSet(
        ctx, key,
        "availableBalance", req.AvailableBalance,
        "currency", req.Currency,
    ).Err(); err != nil {
        zap.S().Errorf("SetBalance HSET key [%s] error [%v]", key, err)
        return &emptypb.Empty{}, kerrors.New(500, "redis_hset_error", "failed to write balance")
    }

    // Optionally set a TTL on the entire hash (EXPIRE applies to the key, not individual fields).
    // if h.ttl > 0 { _ = h.redis.Expire(ctx, key, h.ttl).Err() }

    return &emptypb.Empty{}, nil
}

// GetBalance: reads the balance from Redis and returns CacheResponse.
func (h *CacheHandler) GetBalance(ctx context.Context, req *pb.CacheRequest) (*pb.CacheResponse, error) {
    key := balanceKey(req.AccountNumber)

    m, err := h.redis.HGetAll(ctx, key).Result()
    if err != nil {
        zap.S().Errorf("GetBalance HGETALL key [%s] error [%v]", key, err)
        return nil, kerrors.New(500, "redis_hgetall_error", "failed to read balance")
    }
    // Cache miss: return zero values for the response.
    if len(m) == 0 {
        return &pb.CacheResponse{
         AccountNumber: req.AccountNumber,
            AvailableBalance: 0,
        }, nil
    }

    // Parse availableBalance (stored as string in Redis) into int64.
    var avail int64
    if s, ok := m["availableBalance"]; ok && s != "" {
        if parsed, perr := strconv.ParseInt(s, 10, 64); perr == nil {
            avail = parsed
        } else {
            zap.S().Errorf("GetBalance parse availableBalance [%s] error [%v]", s, perr)
            return nil, kerrors.New(500, "redis_parse_error", "invalid availableBalance in cache")
        }
    }

    return &pb.CacheResponse{
        AccountNumber:req.AccountNumber,
        AvailableBalance: avail,
    }, nil
}

// DeleteBalance: deletes the Redis hash and returns Empty.
func (h *CacheHandler) DeleteBalance(ctx context.Context, req *pb.CacheRequest) (*emptypb.Empty, error) {
    key := balanceKey(req.AccountNumber)
    if err := h.redis.Del(ctx, key).Err(); err != nil {
        zap.S().Errorf("DeleteBalance DEL key [%s] error [%v]", key, err)
        return &emptypb.Empty{}, kerrors.New(500, "redis_del_error", "failed to delete balance")
    }
    return &emptypb.Empty{}, nil
}
