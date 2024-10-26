package main

import (
	"github.com/kubevm.io/vink/config"
	"github.com/kubevm.io/vink/internal/cmd/apiserver"
	cmdctl "github.com/kubevm.io/vink/internal/cmd/ctrl"
	"github.com/kubevm.io/vink/internal/pkg/cache"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	var configFile string

	root := &cobra.Command{
		Use:     "vink",
		Aliases: []string{"vink"},
		Short:   "Virtual Machines in Kubernetes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := ctrl.SetupSignalHandler()

			config, err := config.ParseConfigFromFile(configFile)
			if err != nil {
				return err
			}
			if config.Debug {
				log.SetDebug()
			}

			clients, err := clients.NewClients(config.KubeConfig)
			if err != nil {
				return err
			}

			kubeCache := cache.NewKubeCache(clients)
			if err := kubeCache.Start(ctx); err != nil {
				return err
			}

			ctrl := cmdctl.New(config, clients)
			go func() {
				if err := ctrl.Execute(ctx); err != nil {
					panic(err)
				}
			}()

			apiserver := apiserver.New(config, clients, kubeCache)
			go func() {
				if err := apiserver.Execute(ctx); err != nil {
					panic(err)
				}
			}()

			<-ctx.Done()
			if err := ctrl.Stop(); err != nil {
				log.Error(err)
			}
			if err := apiserver.Stop(); err != nil {
				log.Error(err)
			}

			return nil
		},
	}

	root.PersistentFlags().StringVarP(&configFile, "config", "c", "config/config.yaml", "Config file path.")

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
