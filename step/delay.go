package step

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/jqs7/clash_cli/api"
	"github.com/jqs7/clash_cli/model"
	"golang.org/x/sync/errgroup"
)

type DelayTest struct {
	*api.Client
	LastStep Step
	Proxies  []model.ProxyName
}

func (dt DelayTest) Run() error {
	fmt.Print("Ctrl + C 退出")
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-interrupt
		signal.Stop(interrupt)
		cancel()
		fmt.Print("\r")
	}()
	e, _ := errgroup.WithContext(ctx)
	for _, v := range dt.Proxies {
		if v.ItemType == model.ItemTypeProxy {
			e.Go(getDelay(ctx, dt.Client, v.Name))
		}
	}
	if err := e.Wait(); err != nil && !IsCanceled(err) {
		return err
	}
	return dt.LastStep.Run()
}

func getDelay(ctx context.Context, client *api.Client, name string) func() error {
	return func() error {
		delay, err := client.GetDelay(ctx, name)
		switch err {
		case api.DelayTestError:
			model.Delays.Store(name, "错误")
		case api.DelayTestTimeout:
			model.Delays.Store(name, "超时")
		case nil:
			model.Delays.Store(name, string(delay.Delay))
		default:
			return err
		}
		return nil
	}
}
