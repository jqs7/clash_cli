package step

import (
	"github.com/jqs7/clash_cli/api"
	"github.com/manifoldco/promptui"
)

type Root struct {
	*api.Client
}

func (r Root) Run() error {
	prompt := promptui.Select{
		Label: "功能选择",
		Size:  10,
		Items: []string{"出站模式", "选择代理", "实时速率", "代理日志", "退出"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		return err
	}
	var step Step
	switch result {
	case "出站模式":
		step = SwitchMode{
			Client:   r.Client,
			LastStep: r,
		}
	case "选择代理":
		step = SwitchProxy{
			Client:   r.Client,
			LastStep: r,
		}
	case "实时速率":
		step = Traffic{
			Client:   r.Client,
			LastStep: r,
		}
	case "代理日志":
		step = Log{
			Client:   r.Client,
			LastStep: r,
		}
	default:
		step = Exit{}
	}
	return step.Run()
}
