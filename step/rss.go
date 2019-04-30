package step

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/levigross/grequests"
	"github.com/manifoldco/promptui"
)

type RSS struct {
	LastStep Step
}

type VMessRSS struct {
	Add  string `json:"add"`
	Aid  string `json:"aid"`
	Host string `json:"host"`
	ID   string `json:"id"`
	Net  string `json:"net"`
	Path string `json:"path"`
	Port string `json:"port"`
	Ps   string `json:"ps"`
	TLS  string `json:"tls"`
	Type string `json:"type"`
	V    string `json:"v"`
}

func (r RSS) Run() error {
	prompt := promptui.Prompt{
		Label: "地址",
		Validate: func(s string) error {
			u, err := url.Parse(s)
			if err != nil {
				return err
			}
			if u.Scheme != "http" && u.Scheme != "https" {
				return errors.New("不支持的地址类型")
			}
			return nil
		},
	}
	result, err := prompt.Run()
	if err != nil {
		return err
	}
	resp, err := grequests.Get(result, nil)
	if err != nil {
		return r.Run()
	}
	bs, err := TryBase64Decode(resp.String())
	if err != nil {
		log.Println("Base64 解码失败")
		return r.Run()
	}
	scanner := bufio.NewScanner(bytes.NewReader(bs))
	names := map[string]struct{}{}
	for scanner.Scan() {
		parsed, err := url.Parse(scanner.Text())
		if err != nil {
			log.Println("链接解析失败", scanner.Text())
			return r.Run()
		}
		switch parsed.Scheme {
		case "vmess":
			bs, err := TryBase64Decode(parsed.Host)
			if err != nil {
				log.Println("Base64 解码失败")
				return r.Run()
			}
			vmess := &VMessRSS{}
			if err := json.Unmarshal(bs, vmess); err != nil {
				log.Println("JSON 解码失败")
				return r.Run()
			}
			fmt.Printf("- name: %s\n  type: %s\n  server: %s\n  port: %s\n  "+
				"uuid: %s\n  alterId: %s\n  cipher: %s\n",
				vmess.Ps, "vmess", vmess.Add, vmess.Port,
				vmess.ID, vmess.V, "auto")
			names[vmess.Ps] = struct{}{}
		case "ssr":
			bs, err := TryBase64Decode(parsed.Host)
			if err != nil {
				log.Println("Base64 解码失败")
				return r.Run()
			}
			infos := strings.SplitN(string(bs), ":", 6)
			if len(infos) < 6 {
				log.Println("ss 格式错误")
				return r.Run()
			}
			paq := strings.Split(infos[5], "?")

			ssPwd, err := TryBase64Decode(string(strings.TrimRight(paq[0], "/")))
			if err != nil {
				log.Println("密码 Base64 解码失败")
				return r.Run()
			}

			query, err := url.ParseQuery(paq[1])
			if err != nil {
				log.Println("URL 解析失败")
				return r.Run()
			}
			remarks, err := TryBase64Decode(query.Get("remarks"))
			if err != nil {
				log.Println("备注 Base64 解码失败")
				return r.Run()
			}
			fmt.Printf("- name: \"%s\"\n  type: %s\n  server: %s\n  port: %s\n  "+
				"cipher: %s\n  password: %s\n",
				string(remarks), "ss", infos[0], infos[1],
				infos[3], string(ssPwd))
			names[string(remarks)] = struct{}{}
		}
	}
	fmt.Println()
	for name := range names {
		fmt.Printf("- \"%s\"\n", name)
	}

	return r.LastStep.Run()
}

func TryBase64Decode(s string) (bs []byte, err error) {
	bs, err = base64.RawURLEncoding.DecodeString(s)
	if err == nil {
		return
	}
	bs, err = base64.RawStdEncoding.DecodeString(s)
	if err == nil {
		return
	}
	bs, err = base64.URLEncoding.DecodeString(s)
	if err == nil {
		return
	}
	bs, err = base64.StdEncoding.DecodeString(s)
	if err == nil {
		return
	}
	return nil, err
}
