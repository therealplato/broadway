package cfg

// Broadway configuration object contains all global configs and will be passed to dependencies
type Broadway struct {
	// AuthBearerToken is a global token required for all requests except GET/POST command/
	AuthBearerToken string
	// SlackToken contains the expected Slack custom command token.
	SlackToken string

	// ServerHost is passed to gin and configures the listen address of the server
	ServerHost string

	// K8sServiceHost is the Kubernetes host
	K8sServiceHost string

	// K8sServicePort is the Kubernetes port
	K8sServicePort string

	// K8sNamespace is the namespace used by Broadway's deployments
	K8sNamespace string

	// K8sCertFile is the cert file setting for local development
	K8sCertFile string

	// K8sKeyFile is the key file setting for local development
	K8sKeyFile string

	// K8sCAFile is the CA file setting for local development
	K8sCAFile string

	// EtcdEndpoints is the list Etcd hosts separated by comma
	EtcdEndpoints string

	// EtcdPath is the root directory for Broadway objects
	EtcdPath string

	// SlackWebhook is the Slack incoming webhook
	SlackWebhook string

	// ManifestsPath is the absolute path where manifest files are read from
	ManifestsPath string

	// PlaybooksPath is the absolute path where playbook files are read from
	PlaybooksPath string
}
