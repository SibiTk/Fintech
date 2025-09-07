package server

import (
//	v1 "cache_service/api/helloworld/v1"
	"cache_service/internal/conf"
	"cache_service/internal/service"

	//"github.com/go-kratos/kratos/v2/log"
	//"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.BalanceCacheServiceService) *http.Server {

	srv := http.NewServer(
		http.Address(c.Http.Addr),
	)

	// Register your service handler here
	// e.g., pb.RegisterCacheHTTPServer(srv, cacheService)

	return srv
}
