package deployment

import (
	"io/ioutil"
	"os"

	"github.com/golang/glog"
	"github.com/namely/broadway/cfg"

	"k8s.io/kubernetes/pkg/client/restclient"
)

// IsKubernetesEnv returns true if necessary Kubernetes environment variables
// and files are available
func IsKubernetesEnv(cfg cfg.CommonCfgType) bool {
	files := []string{
		"/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		"/var/run/secrets/kubernetes.io/serviceaccount/token",
	}

	for _, file := range files {
		fi, err := os.Stat(file)
		if err != nil || fi.IsDir() {
			glog.Infof("Not using kubernetes cfg because a file is missing: %s", file)
			return false
		}
	}

	cfgs := []string{
		cfg.K8sServicePort,
		cfg.K8sServiceHost,
	}

	for _, cfg := range cfgs {
		if cfg == "" {
			glog.Info("Not running in kub environment because an cfg var is missing.")
			return false
		}
	}

	return true
}

// Config returns a kubernetes configuration
func Config(cfg cfg.CommonCfgType) (*restclient.Config, error) {
	config := LocalConfig(cfg)
	if IsKubernetesEnv(cfg) {
		var err error
		config, err = KubernetesConfig(cfg)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

// KubernetesConfig returns Kubernetes configuration for native Kubernetes environment
func KubernetesConfig(cfg cfg.CommonCfgType) (*restclient.Config, error) {
	token, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, err
	}
	return &restclient.Config{
		Host:        "https://" + cfg.K8sServiceHost + ":" + cfg.K8sServicePort,
		BearerToken: string(token),
		TLSClientConfig: restclient.TLSClientConfig{
			CAFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
		},
	}, nil
}

// LocalConfig returns a configuration for local development
func LocalConfig(cfg cfg.CommonCfgType) *restclient.Config {
	if cfg.K8sServiceHost != "" && cfg.K8sCertFile != "" &&
		cfg.K8sKeyFile != "" && cfg.K8sCAFile != "" {

		return &restclient.Config{
			Host: cfg.K8sServiceHost,
			TLSClientConfig: restclient.TLSClientConfig{
				CertFile: cfg.K8sCertFile,
				KeyFile:  cfg.K8sKeyFile,
				CAFile:   cfg.K8sCAFile,
			},
		}
	}

	return &restclient.Config{
		Host:     "http://localhost:8080",
		Insecure: true,
	}
}
