package cmd

import (
	"fmt"
	"n1kit0s/vt-manager/app/github"
	"n1kit0s/vt-manager/app/vuetorrent"
	"os"
)

type InstallCommand struct {
	Version      string `short:"v" long:"version" optional:"true" description:"VueTorrent version to install"`
	Directory    string `short:"d" long:"dir" required:"true" description:"VueTorrent directory"`
	GithubApiKey string `short:"k" long:"api-key" required:"true" description:"Github API key"`
}

func (c *InstallCommand) Execute(args []string) error {
	fmt.Println("Start installing")
	var githubClient = github.NewClient(c.GithubApiKey)
	var vtManager = vuetorrent.NewVTManager(githubClient)

	var vtRelease vuetorrent.VueTorrentRelease
	if c.Version == "" {
		release, err := vtManager.GetLatestVuetorrentRelease()
		if err != nil {
			fmt.Printf("Failed to retrive latest vuetorrent release. Error: %s \n", err.Error())
			return err
		}
		vtRelease = release
	} else {
		tag := vuetorrent.MakeTagName(c.Version)
		release, err := vtManager.GetVuetorrentRelease(tag)
		if err != nil {
			fmt.Printf("Failed to retrive vuetorrent release %s. Error: %s \n", tag, err.Error())
			return err
		}
		vtRelease = release
	}

	fmt.Printf("Start downloading %v \n", vtRelease)
	filePath, err := vtManager.Download(vtRelease, os.TempDir())
	if err != nil {
		fmt.Printf("Failed to download latest vuetorrent release: Error: %s \n", err.Error())
		return err
	}
	fmt.Printf("Downloaded latest release into %s \n", filePath)

	err = vtManager.Unzip(filePath, c.Directory, vtRelease.Version)
	if err != nil {
		fmt.Printf("Failed to unzip vuetorrent release. Error: %s \n", err.Error())
		return err
	}

	return nil
}
