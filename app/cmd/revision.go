package cmd

import (
	"log"
	"fmt"
	"runtime/debug"
)

type RevisionCommand struct {}

func (c *RevisionCommand) Execute(args []string) error {
	log.Printf("[INFO] Revision: %s", revision())
	return nil
}

func revision() string {
	var revision string
	var modified bool
	
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}
	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}
	
	return revision
	}