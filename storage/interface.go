package storage

import "github.com/Dreamacro/clash/tunnel"

type IStorage interface {
	SaveMode(tunnel.Mode) error
	GetMode() (tunnel.Mode, error)
	Close()
	UpdateProxy(group, proxy string) error
	GetProxies() (map[string]string, error)
}
