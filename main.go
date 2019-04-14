package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/jqs7/clash_cli/api"
	"github.com/jqs7/clash_cli/step"
	"github.com/jqs7/clash_cli/storage"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	kingpin.HelpFlag.Short('h')
	quit := kingpin.Flag("quit", "apply latest setting and quit").Short('q').Bool()
	urlArg := kingpin.Arg("endpoint", "clash api endpoint").Default("http://localhost:9090").String()
	kingpin.Parse()

	db, err := storage.Open()
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	client := &api.Client{
		BaseURL: getBaseURL(*urlArg),
	}

	time.Sleep(time.Millisecond * 100)
	mode, err := db.GetMode()
	if err == nil {
		if err := client.UpdateMode(mode); err != nil {
			log.Println(err)
		}
	}

	proxies, err := db.GetProxies()
	if err == nil {
		for group, name := range proxies {
			if err := client.UpdateProxy(group, name); err != nil {
				log.Println(err)
			}
		}
	}

	if *quit {
		os.Exit(0)
	}

	root := step.Root{
		Client: client,
	}
	root.SetupDB(db)
	if err := root.Run(); err != nil {
		log.Fatalln(err)
	}
}

func getBaseURL(arg string) string {
	if strings.HasPrefix(arg, "http") {
		return arg
	}
	return "http://" + arg
}
