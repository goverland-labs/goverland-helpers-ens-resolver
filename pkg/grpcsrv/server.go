package grpcsrv

import (
	"fmt"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var prom = prometheus.NewServerMetrics(
	prometheus.WithServerCounterOptions(),
	prometheus.WithServerHandlingTimeHistogram(),
)

func NewGrpcServer(a auth.AuthFunc) *grpc.Server {
	server := grpc.NewServer(
		stdUnaryMiddleware(UnaryReflectionFilter(auth.UnaryServerInterceptor(a))),
		stdStreamMiddleware(StreamReflectionFilter(auth.StreamServerInterceptor(a))),
	)

	stdRegister(server)

	return server
}

func stdUnaryMiddleware(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	arr := []grpc.UnaryServerInterceptor{
		prom.UnaryServerInterceptor(),
		recovery.UnaryServerInterceptor(
			recovery.WithRecoveryHandler(func(p any) (err error) {
				log.Error().Str("stack", string(debug.Stack())).Msg("grpc panic")

				return fmt.Errorf("%#v", p)
			}),
		),
	}
	arr = append(arr, interceptors...)

	return grpc.ChainUnaryInterceptor(arr...)
}

func stdStreamMiddleware(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	arr := []grpc.StreamServerInterceptor{
		prom.StreamServerInterceptor(),
		recovery.StreamServerInterceptor(),
	}
	arr = append(arr, interceptors...)

	return grpc.ChainStreamInterceptor(arr...)
}

func stdRegister(s *grpc.Server) {
	reflection.Register(s)
}
