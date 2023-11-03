package grpcsrv

import (
	"fmt"
	"runtime/debug"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGrpcServer(auth grpc_auth.AuthFunc) *grpc.Server {
	server := grpc.NewServer(
		stdUnaryMiddleware(UnaryReflectionFilter(grpc_auth.UnaryServerInterceptor(auth))),
		stdStreamMiddleware(StreamReflectionFilter(grpc_auth.StreamServerInterceptor(auth))),
	)

	stdRegister(server)

	return server
}

func stdUnaryMiddleware(interceptors ...grpc.UnaryServerInterceptor) grpc.ServerOption {
	arr := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_prometheus.UnaryServerInterceptor,
		grpc_recovery.UnaryServerInterceptor(
			grpc_recovery.WithRecoveryHandler(
				func(i interface{}) error {
					log.Error().Str("stack", string(debug.Stack())).Msg("grpc panic")

					return fmt.Errorf("%#v", i)
				},
			),
		),
	}
	arr = append(arr, interceptors...)

	return grpc.UnaryInterceptor(
		grpc_middleware.ChainUnaryServer(arr...),
	)
}

func stdStreamMiddleware(interceptors ...grpc.StreamServerInterceptor) grpc.ServerOption {
	arr := []grpc.StreamServerInterceptor{
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_recovery.StreamServerInterceptor(),
	}
	arr = append(arr, interceptors...)

	return grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(arr...),
	)
}

func stdRegister(s *grpc.Server) {
	reflection.Register(s)
	grpc_prometheus.EnableHandlingTimeHistogram(
		grpc_prometheus.WithHistogramBuckets([]float64{0.02, 0.05, 0.1, 0.2, 0.3, 0.4, 0.5, 0.8, 1, 1.2, 1.5, 2, 4, 8}),
	)
}
