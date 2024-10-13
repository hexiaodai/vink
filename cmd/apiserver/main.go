package main

import (
	"context"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/kubevm.io/vink/config"
	"github.com/spf13/cobra"

	interalcmd "github.com/kubevm.io/vink/internal/cmd"
	ctrlvm "github.com/kubevm.io/vink/internal/controller/virtualmachine"
	"github.com/kubevm.io/vink/internal/management"
	"github.com/kubevm.io/vink/internal/management/virtualmachine"

	resource_event_listener "github.com/kubevm.io/vink/internal/pkg/resource-event-listener"
	"github.com/kubevm.io/vink/internal/pkg/servers"
	"github.com/kubevm.io/vink/pkg/clients"

	"github.com/kubevm.io/vink/pkg/informer"
	"github.com/kubevm.io/vink/pkg/log"
)

func main() {
	configFile := ""

	root := &cobra.Command{
		Use:     "apiserver",
		Aliases: []string{"api"},
		Short:   "Virtual Machines in Kubernetes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.TODO()

			config.ParseConfigFromFile(configFile)
			if config.Instance.Debug {
				log.SetDebug()
			}

			if err := clients.InitClients(config.Instance.KubeConfig); err != nil {
				return err
			}

			kubeInformerFactory := informer.NewKubeInformerFactory(clients.GetClients().GetVinkRestClient(), clients.GetClients().GetKubeVirtClient().RestClient(), clients.GetClients().GetKubeovnRestClient(), clients.GetClients().GetKubeVirtClient())

			kubeInformerFactory.VirtualMachine()
			kubeInformerFactory.VirtualMachineInstances()
			kubeInformerFactory.DataVolume()
			kubeInformerFactory.VirtualMachineSummary()
			kubeInformerFactory.Subnet()

			kubeInformerFactory.Start(ctx.Done())
			kubeInformerFactory.WaitForCacheSync(ctx.Done())

			resourceEventListener := resource_event_listener.NewResourceEventListener(kubeInformerFactory)
			go resourceEventListener.StartListening(ctx)

			vmCtl := ctrlvm.New(kubeInformerFactory.VirtualMachine())
			go vmCtl.Run(context.TODO())

			mgr, err := interalcmd.NewCRDManager()
			if err != nil {
				return err
			}
			if err := interalcmd.Register(mgr); err != nil {
				return err
			}
			go func() {
				if err := mgr.Start(ctx); err != nil {
					panic(err)
				}
			}()

			register, err := management.RegisterGRPCRoutes(kubeInformerFactory, resourceEventListener)
			if err != nil {
				return err
			}

			httpAddress := fmt.Sprintf("%v:%v", config.Instance.APIServer.Address, config.Instance.APIServer.HTTP)
			grpcAddress := fmt.Sprintf("%v:%v", config.Instance.APIServer.Address, config.Instance.APIServer.GRPC)

			grpcServer := servers.NewGRPCServer(grpcAddress, register)
			log.Infof("Starting grpc server at: %s", grpcAddress)
			go grpcServer.Run()

			gwServer := servers.NewGatewayServer(
				"apiserver",
				httpAddress,
				grpcAddress,
				[]func(router *mux.Router){
					// router hooks
					virtualmachine.RegisterSerialConsole,
				},
				management.RegisterHTTPRoutes(),
			)

			log.Infof("Starting gateway server at: %s -> %s", httpAddress, grpcAddress)
			return gwServer.Run()
		},
	}

	root.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "Config file path.")

	root.Execute()
}
