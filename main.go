package main

import (
	"log"
	"os"
	"strings"

	"github.com/jqs7/clash_cli/api"
	"github.com/jqs7/clash_cli/step"
)

func main() {
	root := step.Root{
		Client: &api.Client{
			BaseURL: getBaseURL(),
		},
	}
	if err := root.Run(); err != nil {
		log.Fatalln(err)
	}
}

func getBaseURL() string {
	baseURL := "http://localhost:9090"
	if len(os.Args) <= 1 {
		return baseURL
	}
	if strings.HasPrefix(os.Args[1], "http") {
		return os.Args[1]
	}
	return "http://" + os.Args[1]
}
