package apiserver

import (
	"context"
	"fmt"

	"github.com/kubevm.io/vink/config"

	ctrlvm "github.com/kubevm.io/vink/internal/controller/virtualmachine"
	"github.com/kubevm.io/vink/internal/management"

	"github.com/kubevm.io/vink/internal/pkg/cache"
	resource_event_listener "github.com/kubevm.io/vink/internal/pkg/resource-event-listener"
	"github.com/kubevm.io/vink/internal/pkg/servers"
	"github.com/kubevm.io/vink/pkg/clients"

	"github.com/kubevm.io/vink/pkg/log"
)

func New(config *config.Configuration, clients clients.Clients, kubeCache *cache.KubeCache) *Daemon {
	return &Daemon{
		config:    config,
		clients:   clients,
		kubeCache: kubeCache,
	}
}

type Daemon struct {
	config    *config.Configuration
	clients   clients.Clients
	kubeCache *cache.KubeCache

	grpcServer servers.Server

	httpServer servers.Server
}

func (dm *Daemon) Execute(ctx context.Context) error {
	resourceEventListener := resource_event_listener.NewResourceEventListener(dm.kubeCache.InformerFactory)
	go resourceEventListener.StartListening(ctx)

	vmCtl := ctrlvm.New(dm.clients, dm.kubeCache.InformerFactory.VirtualMachine())
	go vmCtl.Run(ctx)

	register, err := management.RegisterGRPCRoutes(dm.clients, dm.kubeCache.InformerFactory, resourceEventListener)
	if err != nil {
		return err
	}

	httpAddress := fmt.Sprintf("%v:%v", dm.config.APIServer.Address, dm.config.APIServer.HTTP)
	grpcAddress := fmt.Sprintf("%v:%v", dm.config.APIServer.Address, dm.config.APIServer.GRPC)

	dm.grpcServer = servers.NewGRPCServer(grpcAddress, register)
	log.Infof("Starting grpc server at: %s", grpcAddress)

	go func() {
		if err := dm.grpcServer.Run(); err != nil {
			panic(err)
		}
	}()

	httpRegister, err := management.RegisterHTTPRoutes(dm.clients)

	dm.httpServer = servers.NewHTTPServer("apiserver", httpAddress, httpRegister)
	log.Infof("Starting http server at: %s", httpAddress)
	return dm.httpServer.Run()
}

func (dm *Daemon) Stop() error {
	if err := dm.grpcServer.Stop(); err != nil {
		log.Errorf("failed to stop grpc server: %v", err)
	}
	return dm.httpServer.Stop()
}
