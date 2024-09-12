package virtualmachine

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	vmv1alpha1 "github.com/kubevm.io/vink/apis/management/virtualmachine/v1alpha1"
	"github.com/kubevm.io/vink/internal/management/virtualmachine/business"
	"github.com/kubevm.io/vink/pkg/clients"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/emptypb"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/cert"
	"kubevirt.io/client-go/kubecli"
)

func NewVirtualMachineManagement(clients clients.Clients) vmv1alpha1.VirtualMachineManagementServer {
	return &virtualMachineManagement{
		clients: clients,
	}
}

type virtualMachineManagement struct {
	clients clients.Clients

	vmv1alpha1.UnimplementedVirtualMachineManagementServer
}

func (m *virtualMachineManagement) VirtualMachinePowerState(ctx context.Context, request *vmv1alpha1.VirtualMachinePowerStateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, business.VirtualMachinePowerState(ctx, m.clients, request.NamespaceName, request.PowerState)
}

func RegisterSerialConsole(router *mux.Router) {
	router.PathPrefix(business.SerialConsoleRequestPathTmpl).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		namespace, name := vars["namespace"], vars["name"]

		kv := clients.GetClients().GetKubeVirtClient()

		parse, err := url.Parse(kv.Config().Host)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		ws := fmt.Sprintf("wss://%s/apis/subresources.kubevirt.io/v1/namespaces/%s/virtualmachineinstances/%s/console", parse.Host, namespace, name)
		result, _, err := kubecli.Dial(ws, generateSerialConsoleTLSConfig(kv.Config()))
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		serverConn := result.UnderlyingConn()

		upgrader := kubecli.NewUpgrader()
		upgrader.HandshakeTimeout = 15 * time.Second
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		clientConn := conn.UnderlyingConn()

		ctx, cancel := context.WithCancel(request.Context())
		defer cancel()
		go func() {
			<-ctx.Done()
			serverConn.Close()
			clientConn.Close()
		}()

		eg := errgroup.Group{}
		eg.Go(func() error {
			_, err := io.Copy(clientConn, serverConn)
			return err
		})
		eg.Go(func() error {
			_, err := io.Copy(serverConn, clientConn)
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
