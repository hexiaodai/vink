package cmd

import (
	"strings"
	"time"

	"github.com/kubevm.io/vink/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	envPrefix = "VINK_"
)

func InitFlags(cmd *cobra.Command) error {
	flags := cmd.Flags()

	flags.Bool(config.Debug, config.DebugDefault, "Debug mode")
	bindEnv(config.Debug)

	flags.String(config.APIServerHTTP, config.APIServerHTTPDefault, "APIServer HTTP port")
	bindEnv(config.APIServerHTTPDefault)

	flags.String(config.APIServerGRPC, config.APIServerGRPCDefault, "APIServer GRPC port")
	bindEnv(config.APIServerGRPC)

	flags.String(config.APIServerGRPCWeb, config.APIServerGRPCWebDefault, "APIServer GRPC Web port")
	bindEnv(config.APIServerGRPCWeb)

	flags.String(config.Prometheus, config.PrometheusDefault, "Prometheus port")
	bindEnv(config.Prometheus)

	flags.String(config.Ceph, config.CephDefault, "Ceph port")
	bindEnv(config.Ceph)

	flags.String(config.CephUsername, config.CephUsernameDefault, "Ceph username")
	bindEnv(config.CephUsername)

	flags.String(config.CephPassword, config.CephPasswordDefault, "Ceph password")
	bindEnv(config.CephPassword)

	flags.String(config.CephPasswordSecretName, config.CephPasswordSecretNameDefault, "Ceph password secret name")
	bindEnv(config.CephPasswordSecretName)

	flags.String(config.CephPasswordSecretNamespace, config.CephPasswordSecretNamespaceDefault, "Ceph password secret namespace")
	bindEnv(config.CephPasswordSecretNamespace)

	monitorInterval, err := time.ParseDuration(config.MonitorIntervalDefault)
	if err != nil {
		return err
	}
	flags.Duration(config.MonitorInterval, monitorInterval, "Monitor interval")

	return viper.BindPFlags(cmd.Flags())
}

func bindEnv(option string) {
	viper.BindEnv(option, getEnvName(option))
}

func getEnvName(option string) string {
	under := strings.Replace(option, "-", "_", -1)
	upper := strings.ToUpper(under)
	return envPrefix + upper
}
