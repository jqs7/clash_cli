package step

import (
	"fmt"

	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/tunnel"
	"github.com/jqs7/clash_cli/api"
	"github.com/jqs7/clash_cli/model"
	"github.com/manifoldco/promptui"
)

type SwitchProxy struct {
	*api.Client
	LastStep Step
}

func (sp SwitchProxy) Run() error {
	configs, err := sp.GetConfigs()
	if err != nil {
		return err
	}
	proxies, err := sp.GetProxies()
	if err != nil {
		return err
	}
	switch *configs.Mode {
	case tunnel.Direct:
		fmt.Println(promptui.IconWarn, "当前为直连模式，无需选择代理")
		return sp.LastStep.Run()
	case tunnel.Global:
		if err := sp.UpdateProxy(proxies.Proxies, proxies.Proxies["GLOBAL"],
			"GLOBAL", "选择全局代理"); err != nil {
			return err
		}
		return sp.LastStep.Run()
	default: // Rule
		var selectors []struct {
			model.Proxy
			Name string
		}
		for name, group := range proxies.Proxies {
			if group.Type.Is(C.Selector) && name != "GLOBAL" {
				selectors = append(selectors, struct {
					model.Proxy
					Name string
				}{Name: name, Proxy: group})
			}
		}
		switch len(selectors) {
		case 0:
			fmt.Println(promptui.IconWarn, "当前没有可选择代理")
			return sp.LastStep.Run()
		case 1:
			if err := sp.UpdateProxy(proxies.Proxies, selectors[0].Proxy,
				selectors[0].Name, selectors[0].Name); err != nil {
				return err
			}
			return sp.LastStep.Run()
		default:
		}
	}
	return nil
}

func (sp SwitchProxy) UpdateProxy(proxies map[string]model.Proxy,
	selector model.Proxy, groupName, label string) error {
	items := []model.ProxyName{{Name: "测试延迟", ItemType: model.ItemTypeDelayTest}}
	for _, v := range selector.All {
		switch C.AdapterType(proxies[v.Name].Type) {
		case C.Selector:
			continue
		case C.Fallback, C.URLTest:
			v.ExtraInfo = proxies[v.Name].Now
		}
		if v.Name == selector.Now {
			v.Now = true
		}
		v.ItemType = model.ItemTypeProxy
		items = append(items, v)
	}
	prompt := promptui.Select{
		Label: label,
		Size:  10,
		Items: items,
	}
	result, _, err := prompt.Run()
	if err != nil {
		return err
	}
	if items[result].ItemType == model.ItemTypeDelayTest {
		return DelayTest{
			Proxies:  items,
			Client:   sp.Client,
			LastStep: sp,
		}.Run()
	}
	if err := sp.Client.UpdateProxy(groupName, items[result].Name); err != nil {
		return err
	}
	return db.UpdateProxy(groupName, items[result].Name)
}
