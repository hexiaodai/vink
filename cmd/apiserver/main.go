package main

import (
	"fmt"

	"github.com/gorilla/mux"
	"github.com/kubevm.io/vink/config"
	"github.com/spf13/cobra"

	"github.com/kubevm.io/vink/internal/management"
	"github.com/kubevm.io/vink/internal/pkg/servers"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/log"
)

func main() {
	configFile := ""

	root := &cobra.Command{
		Use:     "apiserver",
		Aliases: []string{"api"},
		Short:   "Virtual Machines in Kubernetes",
		RunE: func(cmd *cobra.Command, args []string) error {
			config.ParseConfigFromFile(configFile)
			if config.Instance.Debug {
				log.SetDebug()
			}

			if err := clients.InitClients(config.Instance.KubeConfig); err != nil {
				return err
			}

			register, err := management.RegisterGRPCRoutes()
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
