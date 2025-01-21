package apiserver

import (
	"context"
	"fmt"

	"github.com/kubevm.io/vink/config"
	"github.com/kubevm.io/vink/internal/management"

	grpcwebproxy "github.com/kubevm.io/vink/internal/pkg/grpc-web-proxy"
	"github.com/kubevm.io/vink/internal/pkg/servers"

	"github.com/kubevm.io/vink/pkg/informer"
	"github.com/kubevm.io/vink/pkg/log"
)

func New() *Daemon {
	return &Daemon{}
}

type Daemon struct {
	informerFactory informer.KubeInformerFactory

	grpcServer servers.Server
	httpServer servers.Server
}

func (dm *Daemon) Execute(ctx context.Context) error {
	dm.informerFactory = informer.NewKubeInformerFactory()
	_ = dm.informerFactory.VirtualMachine()
	_ = dm.informerFactory.VirtualMachineInstances()
	_ = dm.informerFactory.DataVolume()
	_ = dm.informerFactory.VirtualMachineSummary()
	_ = dm.informerFactory.Multus()
	_ = dm.informerFactory.Subnet()
	_ = dm.informerFactory.IPPool()
	_ = dm.informerFactory.VPC()
	_ = dm.informerFactory.Event()
	_ = dm.informerFactory.Namespace()

	dm.informerFactory.Start(ctx.Done())
	dm.informerFactory.WaitForCacheSync(ctx.Done())

	register, err := management.RegisterGRPCRoutes(dm.informerFactory)
	if err != nil {
		return err
	}

	httpAddress := fmt.Sprintf(":%v", config.Instance.APIServer.HTTP)
	grpcAddress := fmt.Sprintf(":%v", config.Instance.APIServer.GRPC)

	dm.grpcServer = servers.NewGRPCServer(grpcAddress, register)
	log.Infof("Starting grpc server at: %s", grpcAddress)

	go func() {
		if err := dm.grpcServer.Run(); err != nil {
			panic(err)
		}
	}()

	httpRegister, err := management.RegisterHTTPRoutes()
	dm.httpServer = servers.NewHTTPServer("apiserver", httpAddress, httpRegister)
	log.Infof("Starting http server at: %s", httpAddress)
	go func() {
		if err := dm.httpServer.Run(); err != nil {
			panic(err)
		}
	}()

	grpcweb := grpcwebproxy.NewDetaultProxy()

	go func() {
		if err := grpcweb.Run(ctx); err != nil {
			panic(err)
		}
	}()

	<-ctx.Done()
	return nil
}

func (dm *Daemon) Shutdown() error {
	if err := dm.grpcServer.Stop(); err != nil {
		log.Errorf("failed to stop grpc server: %v", err)
	}
	return dm.httpServer.Stop()
}
