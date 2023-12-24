package cmd

import (
	"fmt"
	"n1kit0s/vt-manager/app/vuetorrent"
)

type InfoCommand struct {
	Directory    string `short:"d" long:"dir" required:"true" description:"VueTorrent directory"`
}

func (c *InfoCommand) Execute(args []string) error {
	version, err := vuetorrent.GetVersion(c.Directory)
	if err != nil {
		return err
	}

	fmt.Printf("Vuetorrent version: %s\n", version)
	return nil
}