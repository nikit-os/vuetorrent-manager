package main

import (
	"n1kit0s/vt-manager/app/cmd"

	"github.com/jessevdk/go-flags"
	"log"
)

type Opts struct {
	InstallCmd cmd.InstallCommand `command:"install"`
	InfoCmd cmd.InfoCommand `command:"info"`
	ListCmd cmd.ListCommand `command:"list"`
	RevisionCmd cmd.RevisionCommand `command:"revision"`
}

func main() {
	log.Println("[INFO] Starting vt-manager")
	var opts Opts
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalf("[ERROR] %s", err.Error())
	}
}
