package etcdstore

import "github.com/namely/broadway/cfg"

func Setup(cfg cfg.Type) {
	SetupEtcd(cfg)
}
