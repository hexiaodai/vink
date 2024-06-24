package servers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hexiaodai/vink/pkg/log"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type gatewayServer struct {
	service     string
	ctx         context.Context
	listenAddr  string
	grpcAddr    string
	handlers    []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error
	routerHooks []func(router *mux.Router)
}

func (g *gatewayServer) Run() error {
	ctx, cancel := context.WithCancel(g.ctx)
	defer cancel()

	conn, err := grpc.DialContext(ctx, g.grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err := conn.Close(); err != nil {
			log.Errorf("Failed to close a client connection to the gPRC server: %v", err)
		}
	}()

	gw, err := NewGateway(ctx, conn, g.handlers...)
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	if len(g.routerHooks) > 0 {
		for _, hook := range g.routerHooks {
			hook(router)
		}
	}
	router.PathPrefix("/").Handler(gw)
	router.Use(otelmux.Middleware(g.service))

	// Create http metrics  middleware.
	// ref: https://github.com/slok/go-http-metrics
	httpMetricMiddleware := middleware.New(middleware.Config{
		//Recorder: prometheus.NewRecorder(prometheus.Config{
		//	Prefix: metrics.MetricsPrefix,
		//}),
		Service:                g.service,
		GroupedStatus:          false,
		DisableMeasureSize:     false,
		DisableMeasureInflight: false,
	})
	// Wrap our main handler, we pass empty handler ID so the middleware inferes
	// the handler label from the URL.
	h := std.Handler("", httpMetricMiddleware, router)

	s := &http.Server{
		Addr:    g.listenAddr,
		Handler: h,
	}
	if err := http2.ConfigureServer(s, &http2.Server{}); err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		log.Info("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Errorf("Failed to shutdown http server: %v", err)
		}
	}()
	log.Infof("Starting listening at %s", g.listenAddr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Infof("Failed to listen and serve: %v", err)
		return err
	}
	return nil
}

func NewGateway(ctx context.Context, conn *grpc.ClientConn,
	registerFuncs ...func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error,
) (http.Handler, error) {
	gw := runtime.NewServeMux(
		// runtime.WithErrorHandler(httpErrorHandler),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			md := map[string]string{}
			md["path"] = request.URL.Path
			return metadata.New(md)
		}),
		runtime.WithForwardResponseOption(func(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
			return nil
		}),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					// true 枚举字段的值使用数字
					UseEnumNumbers: false,
					// 传给 clients 的 json key 使用下划线 `_`
					UseProtoNames: false,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					// 忽略 client 发送的不存在的 poroto 字段
					DiscardUnknown: true,
				},
			},
		),
	)
	for _, f := range registerFuncs {
		if err := f(ctx, gw, conn); err != nil {
			return nil, err
		}
	}
	return gw, nil
}

func (g gatewayServer) Stop() error {
	return nil
}

func NewGatewayServer(service string, listenAddr, grpcAddr string,
	routerHooks []func(router *mux.Router),
	handlers []func(ctx context.Context, serveMux *runtime.ServeMux, clientConn *grpc.ClientConn) error) Server {
	gw := &gatewayServer{
		service:     service,
		ctx:         context.Background(),
		listenAddr:  listenAddr,
		grpcAddr:    grpcAddr,
		routerHooks: routerHooks,
		handlers:    handlers,
	}
	return gw
}
