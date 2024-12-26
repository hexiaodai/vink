package grpcwebproxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/kubevm.io/vink/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/net/trace" // register in DefaultServerMux

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/mwitkow/go-conntrack"
	"github.com/mwitkow/go-conntrack/connhelpers"
	"github.com/mwitkow/grpc-proxy/proxy"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

var (
	DefaultMaxCallRecvMsgSize = 1024 * 1024 * 4
)

func NewDetaaultProxy() *Proxy {
	return &Proxy{
		BindAddr:           "0.0.0.0",
		HttpPort:           config.Instance.APIServer.GRPCWeb,
		RunHttpServer:      true,
		HealthEndpointName: "_health",
		HealthServiceName:  "health",
	}
}

type Proxy struct {
	AllowAllOrigins     bool
	HttpMaxWriteTimeout time.Duration
	HttpMaxReadTimeout  time.Duration
	BindAddr            string

	TlsServerCert                   string
	TlsServerKey                    string
	TlsServerClientCertVerification string
	TlsServerClientCAFiles          []string

	AllowedOrigins []string

	AllowedHeaders []string

	RunHttpServer bool
	RunTlsServer  bool

	EnableHealthEndpoint     bool
	HealthEndpointName       string
	EnableHealthCheckService bool
	HealthServiceName        string

	EnableRequestDebug bool

	HttpPort    int
	HttpTlsPort int
}

func (p *Proxy) Run(ctx context.Context) error {
	backend := NewDetaultBackend()

	backendConn, err := backend.Dial()
	if err != nil {
		return err
	}
	grpcServer := p.buildGrpcProxyServer(backendConn)
	errChan := make(chan error)

	allowedOrigins := makeAllowedOrigins(p.AllowedOrigins)

	options := []grpcweb.Option{
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(p.makeHttpOriginFunc(allowedOrigins)),
	}

	if len(p.AllowedHeaders) > 0 {
		options = append(
			options,
			grpcweb.WithAllowedRequestHeaders(p.AllowedHeaders),
		)
	}

	wrappedGrpc := grpcweb.WrapServer(grpcServer, options...)

	if !p.RunHttpServer && !p.RunTlsServer {
		return fmt.Errorf("both run_http_server and run_tls_server are set to false")
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/", wrappedGrpc)

	if p.EnableHealthEndpoint {
		logrus.Printf("health endpoint enabled on /%v", p.HealthEndpointName)
		if p.EnableHealthCheckService {
			logrus.Printf("health checking enabled for service '%v'", p.HealthServiceName)
			// Health checking endpoint set up
			healthCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			healthChecker := runHealthChecker(healthCtx, backendConn, p.HealthServiceName)
			serveMux.HandleFunc("/"+p.HealthEndpointName, func(resp http.ResponseWriter, req *http.Request) {
				status := healthChecker.GetStatus()
				resp.WriteHeader(status)
			})
		} else {
			// Health endpoint always returns HTTP status 200 if service is disabled
			serveMux.HandleFunc("/"+p.HealthEndpointName, func(resp http.ResponseWriter, req *http.Request) {
				resp.WriteHeader(http.StatusOK)
			})
		}
	}

	if p.RunHttpServer {
		// Debug server.
		if p.EnableRequestDebug {
			serveMux.Handle("/metrics", promhttp.Handler())
			serveMux.HandleFunc("/debug/requests", func(resp http.ResponseWriter, req *http.Request) {
				trace.Traces(resp, req)
			})
			serveMux.HandleFunc("/debug/events", func(resp http.ResponseWriter, req *http.Request) {
				trace.Events(resp, req)
			})
		}

		debugServer := p.buildServer(wrappedGrpc, serveMux)
		debugListener, err := p.buildListener("http", p.HttpPort)
		if err != nil {
			return err
		}
		p.serveServer(debugServer, debugListener, "http", errChan)
	}

	if p.RunTlsServer {
		servingServer := p.buildServer(wrappedGrpc, serveMux)
		servingListener, err := p.buildListener("http", p.HttpTlsPort)
		if err != nil {
			return err
		}
		tlsConfig, err := p.buildServerTLS()
		if err != nil {
			return err
		}
		servingListener = tls.NewListener(servingListener, tlsConfig)
		p.serveServer(servingServer, servingListener, "http_tls", errChan)
	}

	return nil
}

func (p *Proxy) buildGrpcProxyServer(backendConn *grpc.ClientConn) *grpc.Server {
	// gRPC-wide changes.
	grpc.EnableTracing = true
	// grpc_logrus.ReplaceGrpcLogger(logger)

	// gRPC proxy logic.
	director := func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		outCtx, _ := context.WithCancel(ctx)
		mdCopy := md.Copy()
		delete(mdCopy, "user-agent")
		// If this header is present in the request from the web client,
		// the actual connection to the backend will not be established.
		// https://github.com/improbable-eng/grpc-web/issues/568
		delete(mdCopy, "connection")
		outCtx = metadata.NewOutgoingContext(outCtx, mdCopy)
		return outCtx, backendConn, nil
	}

	// Server with logging and monitoring enabled.
	return grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()), // needed for proxy to function.
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
		grpc.MaxRecvMsgSize(DefaultMaxCallRecvMsgSize),
		grpc_middleware.WithUnaryServerChain(
			// grpc_logrus.UnaryServerInterceptor(logger),
			grpc_prometheus.UnaryServerInterceptor,
		),
		grpc_middleware.WithStreamServerChain(
			// grpc_logrus.StreamServerInterceptor(logger),
			grpc_prometheus.StreamServerInterceptor,
		),
	)
}

func makeAllowedOrigins(origins []string) *allowedOrigins {
	o := map[string]struct{}{}
	for _, allowedOrigin := range origins {
		o[allowedOrigin] = struct{}{}
	}
	return &allowedOrigins{
		origins: o,
	}
}

type allowedOrigins struct {
	origins map[string]struct{}
}

func (a *allowedOrigins) IsAllowed(origin string) bool {
	_, ok := a.origins[origin]
	return ok
}

func (p *Proxy) makeHttpOriginFunc(allowedOrigins *allowedOrigins) func(origin string) bool {
	if p.AllowAllOrigins {
		return func(origin string) bool {
			return true
		}
	}
	return allowedOrigins.IsAllowed
}

func (p *Proxy) makeWebsocketOriginFunc(allowedOrigins *allowedOrigins) func(req *http.Request) bool {
	if p.AllowAllOrigins {
		return func(req *http.Request) bool {
			return true
		}
	} else {
		return func(req *http.Request) bool {
			origin, err := grpcweb.WebsocketRequestOrigin(req)
			if err != nil {
				grpclog.Warning(err)
				return false
			}
			return allowedOrigins.IsAllowed(origin)
		}
	}
}

func (p *Proxy) buildServer(_ *grpcweb.WrappedGrpcServer, handler http.Handler) *http.Server {
	return &http.Server{
		WriteTimeout: p.HttpMaxWriteTimeout,
		ReadTimeout:  p.HttpMaxReadTimeout,
		Handler:      handler,
	}
}

func (p *Proxy) buildListener(name string, port int) (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", p.BindAddr, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen for '%v' on %v: %w", name, port, err)
	}
	return conntrack.NewListener(listener,
		conntrack.TrackWithName(name),
		conntrack.TrackWithTcpKeepAlive(20*time.Second),
		conntrack.TrackWithTracing(),
	), nil
}

func (p *Proxy) serveServer(server *http.Server, listener net.Listener, name string, errChan chan error) {
	go func() {
		logrus.Infof("listening for %s on: %v", name, listener.Addr().String())
		if err := server.Serve(listener); err != nil {
			errChan <- fmt.Errorf("%s server error: %v", name, err)
		}
	}()
}

func (p *Proxy) buildServerTLS() (*tls.Config, error) {
	if len(p.TlsServerCert) == 0 || len(p.TlsServerKey) == 0 {
		return nil, fmt.Errorf("flags server_tls_cert_file and server_tls_key_file must be set")
	}
	tlsConfig, err := connhelpers.TlsConfigForServerCerts(p.TlsServerCert, p.TlsServerKey)
	if err != nil {
		return nil, fmt.Errorf("failed building TLS config: %v", err)
	}
	tlsConfig.MinVersion = tls.VersionTLS12
	switch p.TlsServerClientCertVerification {
	case "none":
		tlsConfig.ClientAuth = tls.NoClientCert
	case "verify_if_given":
		tlsConfig.ClientAuth = tls.VerifyClientCertIfGiven
	case "require":
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	default:
		return nil, fmt.Errorf("unknown value '%v' for server_tls_client_cert_verification", p.TlsServerClientCertVerification)
	}
	if tlsConfig.ClientAuth != tls.NoClientCert {
		if len(p.TlsServerClientCAFiles) > 0 {
			tlsConfig.ClientCAs = x509.NewCertPool()
			for _, path := range p.TlsServerClientCAFiles {
				data, err := ioutil.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("failed reading client CA file %v: %w", path, err)
				}
				if ok := tlsConfig.ClientCAs.AppendCertsFromPEM(data); !ok {
					return nil, fmt.Errorf("failed processing client CA file %v", path)
				}
			}
		} else {
			var err error
			tlsConfig.ClientCAs, err = x509.SystemCertPool()
			if err != nil {
				return nil, fmt.Errorf("no client CA files specified, fallback to system CA chain failed: %v", err)
			}
		}

	}
	tlsConfig, err = connhelpers.TlsConfigWithHttp2Enabled(tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("can't configure h2 handling: %v", err)
	}
	return tlsConfig, nil
}

type healthChecker struct {
	status int
	mutex  sync.Mutex
}

func (h *healthChecker) GetStatus() int {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	return h.status
}

func (h *healthChecker) setServing(serving bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if serving {
		h.status = http.StatusOK
	} else {
		h.status = http.StatusServiceUnavailable
	}
}

// Runs health check on a backend connection for a given service name
// returns *healthChecker to get status from
func runHealthChecker(ctx context.Context, backendConn *grpc.ClientConn, service string) *healthChecker {
	h := new(healthChecker)
	h.status = http.StatusServiceUnavailable

	go func() {
		err := grpcweb.ClientHealthCheck(ctx, backendConn, service, h.setServing)
		if err != nil {
			logrus.Errorf("%s health check service error: %v", service, err)
		}
	}()

	return h
}
