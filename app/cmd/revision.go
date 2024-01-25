package cmd

import (
	"log"
)

type RevisionCommand struct {}

var version string

func (c *RevisionCommand) Execute(args []string) error {
	log.Printf("[INFO] Revision: %s", version)
	return nil
}