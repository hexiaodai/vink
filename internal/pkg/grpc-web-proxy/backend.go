package grpcwebproxy

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/mwitkow/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewDetaultBackend(port int) *Backend {
	return &Backend{
		BackendBackoffMaxDelay: grpc.DefaultBackoffConfig.MaxDelay,
		BackendHostPort:        fmt.Sprintf("0.0.0.0:%s", strconv.Itoa(port)),
	}
}

type Backend struct {
	BackendHostPort         string
	BackendIsUsingTLS       bool
	BackendTlsNoVerify      bool
	BackendTlsClientCert    string
	BackendTlsClientKey     string
	BackendTlsCa            []string
	BackendDefaultAuthority string
	BackendBackoffMaxDelay  time.Duration
}

func (b *Backend) Dial() (*grpc.ClientConn, error) {
	opt := []grpc.DialOption{}
	opt = append(opt, grpc.WithCodec(proxy.Codec()))

	if len(b.BackendDefaultAuthority) > 0 {
		opt = append(opt, grpc.WithAuthority(b.BackendDefaultAuthority))
	}

	if b.BackendIsUsingTLS {
		tlsConfig, err := b.buildBackendTLS()
		if err != nil {
			return nil, fmt.Errorf("failed to build backend TLS config: %v", err)
		}
		opt = append(opt, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	} else {
		opt = append(opt, grpc.WithInsecure())
	}

	opt = append(opt,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(DefaultMaxCallRecvMsgSize)),
		grpc.WithBackoffMaxDelay(b.BackendBackoffMaxDelay),
	)

	cc, err := grpc.Dial(b.BackendHostPort, opt...)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (b *Backend) buildBackendTLS() (*tls.Config, error) {
	tlsConfig := &tls.Config{}
	tlsConfig.MinVersion = tls.VersionTLS12
	if b.BackendTlsNoVerify {
		tlsConfig.InsecureSkipVerify = true
	} else if len(b.BackendTlsCa) > 0 {
		tlsConfig.RootCAs = x509.NewCertPool()
		for _, path := range b.BackendTlsCa {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("failed reading backend CA file %v: %w", path, err)
			}
			if ok := tlsConfig.RootCAs.AppendCertsFromPEM(data); !ok {
				return nil, fmt.Errorf("failed processing backend CA file %v", path)
			}
		}
	}
	if len(b.BackendTlsClientCert) > 0 || len(b.BackendTlsClientKey) > 0 {
		if len(b.BackendTlsClientCert) == 0 {
			return nil, fmt.Errorf("'backend_client_tls_cert_file' must be set when 'backend_client_tls_key_file' is set")
		}
		if len(b.BackendTlsClientKey) == 0 {
			return nil, fmt.Errorf("'backend_client_tls_key_file' must be set when 'backend_client_tls_cert_file' is set")
		}
		cert, err := tls.LoadX509KeyPair(b.BackendTlsClientCert, b.BackendTlsClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed reading TLS client keys: %v", err)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	}
	return tlsConfig, nil
}
