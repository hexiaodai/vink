package servers

import (
	"context"
	"fmt"
	"net"
	"runtime/debug"

	grpcotel "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"github.com/kubevm.io/vink/internal/pkg/servers/route"
	"github.com/kubevm.io/vink/pkg/log"
)

type grpcServer struct {
	LogLevel   string
	ListenAddr string
	routes     route.GRPCRouterRegister
	server     *grpc.Server
}

func NewGRPCServer(listenAddr string, routes route.GRPCRouterRegister) Server {
	ser := grpcServer{
		ListenAddr: listenAddr,
		routes:     routes,
		// server:     grpc.NewServer(),
		// TODO enable auth later
		server: grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				// TODO: use StatsHandler
				grpcotel.UnaryServerInterceptor(), // nolint: staticcheck
				func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
					defer func() {
						r := recover()
						if r != nil {
							fmt.Println(string(debug.Stack()))
							//log.Errorf("stacktrace from(%s) panic: \n%s", info.FullMethod, string(debug.Stack()))
							err = fmt.Errorf("%v", r)
							resp = nil
						}
					}()
					if req, ok := req.(interface {
						ValidateAll() error
					}); ok {
						err := req.ValidateAll()
						if err != nil {
							log.Debugf("validate request %s for %+v error: %v", info.FullMethod, req, err)
							return nil, err
						}
					}

					m, err := handler(ctx, req)
					if err != nil {
						return m, grpcErrorHandler(err)
					}

					return m, nil
				},
			),
		),
	}
	return &ser
}

func (m *grpcServer) Run() error {
	m.routes(m.server)
	l, err := net.Listen("tcp", m.ListenAddr)
	if err != nil {
		return err
	}
	return m.server.Serve(l)
}

func (m *grpcServer) Stop() error {
	m.server.GracefulStop()
	return nil
}

func grpcErrorHandler(err error) error {
	log.Error(err)
	fmt.Println(string(debug.Stack()))

	return err
}
