package route

import (
	"github.com/gorilla/mux"
	"google.golang.org/grpc/reflection"
)

type (
	HTTPRouterRegister func(r *mux.Router)
	GRPCRouterRegister func(server reflection.GRPCServer)
)
