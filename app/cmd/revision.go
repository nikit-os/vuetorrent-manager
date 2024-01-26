package cmd

import (
	"fmt"
	"log/slog"
)

type RevisionCommand struct{}

var version string

func (c *RevisionCommand) Execute(args []string) error {
	slog.Info(fmt.Sprintf("Revision: %s", version))
	return nil
}
