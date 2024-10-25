package servers

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kubevm.io/vink/internal/pkg/servers/route"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
	"golang.org/x/net/http2"
)

type httpServer struct {
	LogLevel   string
	ListenAddr string
	service    string
	register route.HTTPRouterRegister
	server   *http.Server
}

func NewHTTPServer(service, listenAddr string, register route.HTTPRouterRegister) Server {
	return &httpServer{
		ListenAddr: listenAddr,
		register:   register,
		service:    service,
	}
}

func (m *httpServer) Run() error {
	router := mux.NewRouter()

	// Create http metrics  middleware.
	// ref: https://github.com/slok/go-http-metrics
	httpMetricMiddleware := middleware.New(middleware.Config{
		//Recorder: prometheus.NewRecorder(prometheus.Config{
		//	Prefix: metrics.MetricsPrefix,
		//}),
		Service:                m.service,
		GroupedStatus:          false,
		DisableMeasureSize:     false,
		DisableMeasureInflight: false,
	})
	// Wrap our main handler, we pass empty handler ID so the middleware inferes
	// the handler label from the URL.

	m.register(router)
	h := std.Handler("", httpMetricMiddleware, router)

	m.server = &http.Server{
		Addr:    m.ListenAddr,
		Handler: h,
	}
	if err := http2.ConfigureServer(m.server, &http2.Server{}); err != nil {
		return err
	}

	log.Infof("Starting listening at %s", m.ListenAddr)
	if err := m.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Infof("Failed to listen and serve: %v", err)
		return err
	}

	return nil
}

func (m *httpServer) Stop() error {
	log.Info("Shutting down the http server")
	if err := m.server.Shutdown(context.Background()); err != nil {
		log.Errorf("Failed to shutdown http server: %v", err)
	}

	return nil
}

func httpErrorHandler(_ context.Context, _ *runtime.ServeMux, _ runtime.Marshaler, writer http.ResponseWriter, _ *http.Request, err error) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write([]byte(err.Error()))
}
