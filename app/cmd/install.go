package cmd

import (
	"log"
	"n1kit0s/vt-manager/app/github"
	"n1kit0s/vt-manager/app/vuetorrent"
)

type InstallCommand struct {
	Version      string `short:"v" long:"version" optional:"true" description:"VueTorrent version to install"`
	Directory    string `short:"d" long:"dir" required:"true" description:"VueTorrent directory"`
	GithubApiKey string `short:"k" long:"api-key" required:"true" description:"Github API key"`
}

func (c *InstallCommand) Execute(args []string) error {
	log.Println("[INFO] Start installing")
	var githubClient = github.NewClient(c.GithubApiKey)
	var vtManager = vuetorrent.NewVTManager(githubClient)

	var vtRelease vuetorrent.VueTorrentRelease
	if c.Version == "" {
		release, err := vtManager.GetLatestVuetorrentRelease()
		if err != nil {
			return err
		}
		vtRelease = release
	} else {
		tag := vuetorrent.MakeTagName(c.Version)
		release, err := vtManager.GetVuetorrentRelease(tag)
		if err != nil {
			return err
		}
		vtRelease = release
	}

	err := vtManager.Install(vtRelease, c.Directory)
	if err != nil {
		return err
	}

	return nil
}
