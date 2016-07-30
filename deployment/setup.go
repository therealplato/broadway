package deployment

import "github.com/namely/broadway/cfg"

// Setup configures deployment package with an injected configuration
func Setup(cfg cfg.Type) {
	SetupPlaybook(cfg)
	SetupKubernetes(cfg)
}
