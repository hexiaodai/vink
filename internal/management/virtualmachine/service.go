package virtualmachine

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/virtualmachine/business"
	"github.com/kubevm.io/vink/pkg/clients"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/cert"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func NewVirtualMachineManagement() vmv1alpha1.VirtualMachineManagementServer {
	return &virtualMachineManagement{}
}

type virtualMachineManagement struct {
	vmv1alpha1.UnimplementedVirtualMachineManagementServer
}

func (m *virtualMachineManagement) VirtualMachinePowerState(ctx context.Context, request *vmv1alpha1.VirtualMachinePowerStateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, business.VirtualMachinePowerState(ctx, request.NamespaceName, request.PowerState)
}

func RegisterSerialConsole(router *mux.Router) {
	router.PathPrefix(business.SerialConsoleRequestPathTmpl).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		namespace, name := vars["namespace"], vars["name"]

		parse, err := url.Parse(clients.Clients.Config().Host)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		ws := fmt.Sprintf("wss://%s/apis/subresources.kubevirt.io/v1/namespaces/%s/virtualmachineinstances/%s/console", parse.Host, namespace, name)

		dialer := websocket.Dialer{
			HandshakeTimeout: 15 * time.Second,
			TLSClientConfig:  generateSerialConsoleTLSConfig(clients.Clients.Config()),
		}

		serverConn, _, err := dialer.Dial(ws, http.Header{})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer serverConn.Close()

		upgrader := websocket.Upgrader{
			HandshakeTimeout: 15 * time.Second,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		clientConn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer clientConn.Close()

		eg := errgroup.Group{}
		eg.Go(func() error {
			_, err := io.Copy(clientConn.UnderlyingConn(), serverConn.UnderlyingConn())
			return err
		})
		eg.Go(func() error {
			_, err := io.Copy(serverConn.UnderlyingConn(), clientConn.UnderlyingConn())
			return err
		})

		if err := eg.Wait(); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}

func generateSerialConsoleTLSConfig(restConfig *rest.Config) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		GetClientCertificate: func(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return func(restConfig *rest.Config) (*tls.Certificate, error) {
				certBytes := restConfig.CertData
				keyBytes := restConfig.KeyData

				crt, err := tls.X509KeyPair(certBytes, keyBytes)
				if err != nil {
					return nil, fmt.Errorf("failed to load certificate: %v", err)
				}
				leaf, err := cert.ParseCertsPEM(certBytes)
				if err != nil {
					return nil, fmt.Errorf("failed to load leaf certificate: %v", err)
				}
				crt.Leaf = leaf[0]
				return &crt, nil
			}(restConfig)
		},
	}
}
