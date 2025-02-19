package config

import (
	"time"

	"github.com/spf13/viper"
)

const (
	DebugDefault = false

	APIServerHTTPDefault = "9090"

	APIServerGRPCDefault = "9091"

	APIServerGRPCWebDefault = "8080"

	PrometheusDefault = "http://monitoring-kube-prometheus-prometheus.monitoring.svc.cluster.local:9090"
	// PrometheusDefault = "http://192.168.18.199:31511"

	// CephDefault = "https://192.168.18.199:8443"
	CephDefault = "https://rook-ceph-mgr-dashboard.rook-ceph.svc.cluster.local:8443"

	CephUsernameDefault = "admin"

	CephPasswordDefault = "dpDloRDp7'lol=NI;/RO"

	MonitorIntervalDefault = "1m"
)

const (
	Debug = "debug"

	APIServerHTTP = "apiserver-http"

	APIServerGRPC = "apiserver-grpc"

	APIServerGRPCWeb = "apiserver-grpc-web"

	Prometheus = "prometheus"

	Ceph = "ceph"

	CephUsername = "ceph-username"

	CephPassword = "ceph-password"

	MonitorInterval = "monitor-interval"
)

type Config struct {
	Debug bool

	APIServerHTTP int

	APIServerGRPC int

	APIServerGRPCWeb int

	Prometheus string

	Ceph string

	CephUsername string

	CephPassword string

	MonitorInterval time.Duration
}

func (c *Config) Populate() {
	c.Debug = viper.GetBool(Debug)
	c.APIServerHTTP = viper.GetInt(APIServerHTTP)
	c.APIServerGRPC = viper.GetInt(APIServerGRPC)
	c.APIServerGRPCWeb = viper.GetInt(APIServerGRPCWeb)
	c.Prometheus = viper.GetString(Prometheus)
	c.Ceph = viper.GetString(Ceph)
	c.CephUsername = viper.GetString(CephUsername)
	c.CephPassword = viper.GetString(CephPassword)
	c.MonitorInterval = viper.GetDuration(MonitorInterval)
}
