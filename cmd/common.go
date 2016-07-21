package cmd

import "gopkg.in/urfave/cli.v1"

// CommonFlags declare configuration that is used by broadway's packages
var CommonFlags = []cli.Flag{
	cli.StringFlag{
		Name:        "kubernetes-host, k8s-host",
		Usage:       "the Kubernetes service host",
		EnvVar:      "KUBERNETES_SERVICE_HOST",
		Destination: &cfg.CommonCfg.K8sServiceHost,
	},
	cli.StringFlag{
		Name:        "kubernetes-port, k8s-port",
		Usage:       "the Kubernetes service port",
		EnvVar:      "KUBERNETES_PORT_443_TCP_PORT",
		Destination: &cfg.CommonCfg.K8sServicePort,
	},
	cli.StringFlag{
		Name:        "kubernetes-namespace, k8s-ns",
		Value:       "broadway",
		Usage:       "broadway's Kubernetes namespace",
		EnvVar:      "KUBERNETES_NAMESPACE",
		Destination: &cfg.CommonCfg.K8sNamespace,
	},
	cli.StringFlag{
		Name:        "kubernetes-cert-file, k8s-cert",
		Usage:       "path to Kubernetes certificate",
		EnvVar:      "KUBERNETES_CERT_FILE",
		Destination: &cfg.CommonCfg.K8sCertFile,
	},
	cli.StringFlag{
		Name:        "kubernetes-key-file, k8s-key",
		Usage:       "path to Kubernetes key file",
		EnvVar:      "KUBERNETES_KEY_FILE",
		Destination: &cfg.CommonCfg.K8sKeyFile,
	},
	cli.StringFlag{
		Name:        "kubernetes-ca-file, k8s-ca",
		Usage:       "path to Kubernetes CA file",
		EnvVar:      "KUBERNETES_CA_FILE",
		Destination: &cfg.CommonCfg.K8sCAFile,
	},
	cli.StringFlag{
		Name:        "etcd-endpoints",
		Usage:       "one or more comma separated etcd endpoints",
		Value:       "http://localhost:4001",
		EnvVar:      "ETCD_ENDPOINTS",
		Destination: &cfg.CommonCfg.EtcdEndpoints,
	},
	cli.StringFlag{
		Name:        "etcd-path",
		Usage:       "an etcd key prefix beginning with /",
		Value:       "/broadway",
		EnvVar:      "ETCD_PATH",
		Destination: &cfg.CommonCfg.EtcdPath,
	},
}
