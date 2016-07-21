package cfg

// CommonCfg is this deployment's global configuration object
var CommonCfg CommonCfgType

// CommonConfigType declares what the common config looks like
type CommonCfgType struct {
	K8sServiceHost string // the Kubernetes host
	K8sServicePort string // the Kubernetes port
	K8sNamespace   string // the namespace used by Broadway's deployments
	K8sCertFile    string // the cert file setting for local development
	K8sKeyFile     string // the key file setting for local development
	K8sCAFile      string // the CA file setting for local development
	EtcdEndpoints  string // the list Etcd hosts separated by comma
	EtcdPath       string // the root directory for Broadway objects
}
