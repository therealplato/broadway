package etcdstore

import "github.com/namely/broadway/pkg/cfg"

// Setup configures etcdstore package with an injected configuration
func Setup(cfg cfg.Type) {
	SetupEtcd(cfg)
}
