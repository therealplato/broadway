package cfg

// GlobalCfg is this deployment's global configuration object
// Only touch this in main.go and cmd/*.go, inject cfg dependency into all other code
var GlobalCfg Type

// Type declares what the common config looks like
type Type struct {
	K8sServiceHost         string // the Kubernetes host
	K8sServicePort         string // the Kubernetes port
	K8sNamespace           string // the namespace used by Broadway's deployments
	K8sCertFile            string // the cert file setting for local development
	K8sKeyFile             string // the key file setting for local development
	K8sCAFile              string // the CA file setting for local development
	EtcdEndpoints          string // the list Etcd hosts separated by comma
	EtcdPath               string // the root directory for Broadway objects
	PlaybooksPath          string // the folder where playbooks are found
	ManifestsPath          string // the folder where manifests are found
	ManifestsExtension     string // .yml or .yaml
	AuthBearerToken        string // a global token required for all requests except GET/POST command/
	SlackToken             string // the expected Slack custom command token.
	ServerHost             string // passed to gin and configures the listen address of the server
	SlackWebhook           string // your team's slack incoming message webhook URL
	InstanceExpirationDays int    // the amount of time in days for expiring an Instance
	InstanceCleanup        int    // the amount of time in seconds for doing the expired instances cleanup
}
