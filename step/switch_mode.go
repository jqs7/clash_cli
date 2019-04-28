package step

import (
	T "github.com/Dreamacro/clash/tunnel"
	"github.com/jqs7/clash_cli/api"
	"github.com/jqs7/clash_cli/storage"
	"github.com/manifoldco/promptui"
)

type SwitchMode struct {
	*api.Client
	LastStep Step
}

func (sm SwitchMode) Run() error {
	configs, err := sm.GetConfigs()
	if err != nil {
		return err
	}

	items := make([]string, 3)
	for i, v := range []T.Mode{T.Global, T.Rule, T.Direct} {
		if v != *configs.Mode {
			items[i] = v.String()
			continue
		}
		items[i] = v.String() + promptui.IconGood
	}
	prompt := promptui.Select{
		Label: "选择出站模式",
		Items: items,
	}
	result, _, err := prompt.Run()

	if err != nil {
		return err
	}
	if err := sm.UpdateMode(T.Mode(result)); err != nil {
		return err
	}
	if err := sm.SaveMode(T.Mode(result)); err != nil {
		return err
	}
	return sm.LastStep.Run()
}

func (sm SwitchMode) SaveMode(mode T.Mode) error {
	db, err := storage.Open()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.SaveMode(mode)
}
