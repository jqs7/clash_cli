package model

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/manifoldco/promptui"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
	T "github.com/Dreamacro/clash/tunnel"
	humanize "github.com/dustin/go-humanize"
)

type Config struct {
	AllowLan  *bool   `json:"allow-lan"`
	LogLevel  *string `json:"log-level"`
	Mode      *T.Mode `json:"mode"`
	Port      *int    `json:"port"`
	RedirPort *int    `json:"redir-port"`
	SocksPort *int    `json:"socks-port"`
}

type Proxies struct {
	Proxies map[string]Proxy `json:"proxies"`
}

type Proxy struct {
	All  []ProxyName `json:"all"`
	Type AdapterType `json:"type,omitempty"`
	Now  string      `json:"now,omitempty"`
}

const (
	ItemTypeProxy = iota
	ItemTypeDelayTest
)

var Delays sync.Map

type ProxyName struct {
	Now       bool
	ItemType  int
	Name      string
	ExtraInfo string
}

func (p ProxyName) String() string {
	s := p.Name
	if p.ExtraInfo != "" {
		s += " [" + p.ExtraInfo + "]"
	}
	if v, ok := Delays.Load(p.Name); ok {
		s += " [" + v.(string) + "]"
	}
	if p.Now {
		s += " " + promptui.IconGood
	}
	return s
}

func (p *ProxyName) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	p.Name = s
	return nil
}

type AdapterType C.AdapterType

func (a *AdapterType) UnmarshalJSON(b []byte) error {
	switch strings.Trim(string(b), "\"") {
	case C.Direct.String():
		*a = AdapterType(C.Direct)
	case C.Fallback.String():
		*a = AdapterType(C.Fallback)
	case C.Reject.String():
		*a = AdapterType(C.Reject)
	case C.Selector.String():
		*a = AdapterType(C.Selector)
	case C.Shadowsocks.String():
		*a = AdapterType(C.Shadowsocks)
	case C.Socks5.String():
		*a = AdapterType(C.Socks5)
	case C.Http.String():
		*a = AdapterType(C.Http)
	case C.URLTest.String():
		*a = AdapterType(C.URLTest)
	case C.Vmess.String():
		*a = AdapterType(C.Vmess)
	case C.LoadBalance.String():
		*a = AdapterType(C.LoadBalance)
	default:
		return errors.New("unknown adapter type")
	}
	return nil
}

func (a AdapterType) Is(v C.AdapterType) bool {
	return C.AdapterType(a) == v
}

type HumanBytes string

func (h *HumanBytes) UnmarshalJSON(d []byte) error {
	var i uint64
	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}
	*h = HumanBytes(humanize.Bytes(i))
	return nil
}

type Traffic struct {
	Up   HumanBytes `json:"up"`
	Down HumanBytes `json:"down"`
}

type Log struct {
	Type    log.LogLevel `json:"type"`
	Payload string       `json:"payload"`
}

type HumanDelay string

func (h *HumanDelay) UnmarshalJSON(d []byte) error {
	var i int
	if err := json.Unmarshal(d, &i); err != nil {
		return err
	}
	*h = HumanDelay((time.Duration(i) * time.Millisecond).String())
	return nil
}

type Delay struct {
	Delay HumanDelay
}
