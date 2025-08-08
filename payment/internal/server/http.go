package server

import (
	v1 "payment/api/helloworld/v1"
	"payment/internal/conf"
	"payment/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"

	h "github.com/gorilla/handlers"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.PaymentService, logger log.Logger) *kratoshttp.Server {
	var opts = []kratoshttp.ServerOption{
		kratoshttp.Middleware(
			recovery.Recovery(),
		),
	}

	if c.Http.Network != "" {
		opts = append(opts, kratoshttp.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, kratoshttp.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, kratoshttp.Timeout(c.Http.Timeout.AsDuration()))
	}

	// Step 3: Wrap Kratos server with CORS middleware
	cors := http.Filter(h.CORS(
		h.AllowedOrigins([]string{"http://localhost:5173"}),
		h.AllowedMethods([]string{"GET", "POST", "PATCH", "DELETE"}),
		h.AllowedHeaders([]string{"Content-Type", "Content-Disposition", ""}),
		h.ExposedHeaders([]string{"Content-Disposition", "Content-Type"}),
	))
	opts = append(opts, cors)

	// Step 1: Create Kratos HTTP server
	srv := kratoshttp.NewServer(opts...)

	// Step 2: Register the service
	v1.RegisterPaymentHTTPServer(srv, greeter)

	return srv
}
