package cmd

import (
	"fmt"
	"log/slog"
	"n1kit0s/vt-manager/app/vuetorrent"
)

type InfoCommand struct {
	Directory string `short:"d" long:"dir" required:"true" description:"VueTorrent directory" env:"VUETORRENT_DIRECTORY"`
}

func (c *InfoCommand) Execute(args []string) error {
	version, err := vuetorrent.GetInstalledVersion(c.Directory)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Version: %s", version))
	return nil
}
