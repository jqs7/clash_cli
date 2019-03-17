package step

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/jqs7/clash_cli/api"
	"github.com/jqs7/clash_cli/model"
)

type Log struct {
	*api.Client
	LastStep Step
}

func (l Log) Run() error {
	if err := handleReaderWithInterrupt(l.GetLogs, func(b []byte) error {
		v := &model.Log{}
		if err := json.Unmarshal(b, v); err != nil {
			return err
		}
		log.Printf("[%s] %s\n", v.Type, v.Payload)
		return nil
	}); err != nil && !IsCanceled(err) {
		return err
	}

	return l.LastStep.Run()
}

func handleReaderWithInterrupt(reader func(context.Context) (io.ReadCloser, error), f func([]byte) error) error {
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	c, cancel := context.WithCancel(context.Background())
	go func() {
		<-interrupt
		signal.Stop(interrupt)
		cancel()
	}()
	fmt.Println("Ctrl + C 退出")
	r, err := reader(c)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		select {
		case <-interrupt:
			signal.Stop(interrupt)
			return r.Close()
		default:
			if err := f(scanner.Bytes()); err != nil {
				return err
			}
		}
	}
	return nil
}

func IsCanceled(err error) bool {
	return strings.HasSuffix(err.Error(), context.Canceled.Error())
}
