package main

import (
	"github.com/kubevm.io/vink/config"
	vinkcmd "github.com/kubevm.io/vink/internal/cmd"
	"github.com/kubevm.io/vink/internal/cmd/apiserver"
	cmdctl "github.com/kubevm.io/vink/internal/cmd/ctrl"
	"github.com/kubevm.io/vink/pkg/clients"
	"github.com/kubevm.io/vink/pkg/log"
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {
	root := &cobra.Command{
		Use:     "vink",
		Aliases: []string{"vink"},
		Short:   "Virtual Machines in Kubernetes",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := ctrl.SetupSignalHandler()

			config := &config.Config{}
			config.Populate()

			log.InitEngine(&log.Config{Debug: config.Debug, Output: "stdout"})

			if err := clients.InitClients(ctx, config); err != nil {
				return err
			}

			ctrl := cmdctl.New(config)
			go func() {
				if err := ctrl.Execute(ctx); err != nil {
					panic(err)
				}
			}()

			apiserver := apiserver.New(config)
			go func() {
				if err := apiserver.Execute(ctx); err != nil {
					panic(err)
				}
			}()

			<-ctx.Done()
			if err := ctrl.Shutdown(); err != nil {
				log.Error(err)
			}
			if err := apiserver.Shutdown(); err != nil {
				log.Error(err)
			}

			return nil
		},
	}

	vinkcmd.InitFlags(root)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
