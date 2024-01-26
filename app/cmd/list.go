package cmd

import (
	"fmt"
	"n1kit0s/vt-manager/app/github"
	"n1kit0s/vt-manager/app/vuetorrent"
)

type ListCommand struct {
	GithubApiKey string `short:"k" long:"api-key" required:"true" description:"Github API key"`
}

func (c *ListCommand) Execute(args []string) error {
	var githubClient = github.NewClient(c.GithubApiKey)
	var vtManager = vuetorrent.NewVTManager(githubClient)

	releases, err := vtManager.GetAllReleases()
	if err != nil {
		return err
	}

	for _, release := range releases {
		fmt.Println(release.Version)
	}

	return nil
}
