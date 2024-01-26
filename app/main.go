package main

import (
	"n1kit0s/vt-manager/app/cmd"
	"os"

	"log/slog"

	"github.com/jessevdk/go-flags"
)

type Opts struct {
	InstallCmd  cmd.InstallCommand  `command:"install"`
	InfoCmd     cmd.InfoCommand     `command:"info"`
	ListCmd     cmd.ListCommand     `command:"list"`
	RevisionCmd cmd.RevisionCommand `command:"revision"`
}

func main() {
	var opts Opts
	_, err := flags.Parse(&opts)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
