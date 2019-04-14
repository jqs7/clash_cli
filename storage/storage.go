package storage

import (
	"log"

	"github.com/Dreamacro/clash/tunnel"

	"github.com/asdine/storm"
)

const (
	Bucket  = "CLASH_CLI"
	Mode    = "CLASH_CLI_MODE"
	Proxies = "CLASH_CLI_PROXIES"
)

type Client struct {
	db *storm.DB
}

func Open() (*Client, error) {
	db, err := storm.Open("clash_cli.db")
	if err != nil {
		return nil, err
	}
	return &Client{db: db}, nil
}

func (c *Client) SaveMode(mode tunnel.Mode) error {
	return c.db.Set(Bucket, Mode, mode)
}

func (c *Client) GetMode() (tunnel.Mode, error) {
	var mode tunnel.Mode
	if err := c.db.Get(Bucket, Mode, &mode); err != nil {
		return -1, err
	}
	return mode, nil
}

func (c *Client) Close() {
	if err := c.db.Close(); err != nil {
		log.Println(err)
	}
}

func (c *Client) UpdateProxy(group, proxy string) error {
	proxies, err := c.GetProxies()
	if err != nil || proxies == nil {
		proxies = make(map[string]string)
	}
	proxies[group] = proxy
	return c.db.Set(Bucket, Proxies, proxies)
}

func (c *Client) GetProxies() (map[string]string, error) {
	var m map[string]string
	if err := c.db.Get(Bucket, Proxies, &m); err != nil {
		return nil, err
	}
	return m, nil
}
