package step

import "os"

type Exit struct{}

func (Exit) Run() error {
	os.Exit(0)
	return nil
}
