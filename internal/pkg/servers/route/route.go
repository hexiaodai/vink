package route

import (
	"context"

	"google.golang.org/grpc/reflection"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	HTTPRouterRegister func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error
	GRPCRouterRegister func(server reflection.GRPCServer)
)
