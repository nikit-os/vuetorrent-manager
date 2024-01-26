package cmd

import (
	"n1kit0s/vt-manager/app/github"
	"n1kit0s/vt-manager/app/vuetorrent"
)

type InstallCommand struct {
	Version      string `short:"v" long:"version" optional:"true" description:"VueTorrent version to install" env:"VUETORRENT_INSTALL_VERSION"`
	Directory    string `short:"d" long:"dir" required:"true" description:"VueTorrent directory" env:"VUETORRENT_DIRECTORY"`
	GithubApiKey string `short:"k" long:"api-key" required:"true" description:"Github API key" env:"GITHUB_API_KEY"`
}

func (c *InstallCommand) Execute(args []string) error {
	var githubClient = github.NewClient(c.GithubApiKey)
	var vtManager = vuetorrent.NewVTManager(githubClient)

	var vtRelease vuetorrent.Release
	if c.Version == "" {
		release, err := vtManager.GetLatestRelease()
		if err != nil {
			return err
		}
		vtRelease = release
	} else {
		tag := vuetorrent.MakeTagName(c.Version)
		release, err := vtManager.GetReleaseByTag(tag)
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
