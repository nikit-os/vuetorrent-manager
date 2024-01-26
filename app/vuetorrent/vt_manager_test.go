package vuetorrent

import (
	"n1kit0s/vt-manager/app/github"
	"testing"
)

type mockGithubClient struct {
}

func (c *mockGithubClient) GetReleases() ([]github.Release, error) {
	return []github.Release{
		{TagName: "v1.1.1", Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent-111.zip"}}},
		{TagName: "v1.1.2", Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent-112.zip"}}},
		{TagName: "v1.1.3", Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent-113.zip"}}},
	}, nil
}
func (c *mockGithubClient) GetReleaseByTag(tag string) (github.Release, error) {
	return github.Release{
		TagName: tag, Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent.zip"}},
	}, nil
}

func TestGetAllReleases(t *testing.T) {
	// Setup
	githubClient := &mockGithubClient{}
	vtManager := NewVTManager(githubClient)
	expectedReleases := []Release{
		{Version: "1.1.1", DownloadUrl: "http://localhost:9876/dw/vuetorrent-111.zip"},
		{Version: "1.1.2", DownloadUrl: "http://localhost:9876/dw/vuetorrent-112.zip"},
		{Version: "1.1.3", DownloadUrl: "http://localhost:9876/dw/vuetorrent-113.zip"},
	}

	// Run
	releases, _ := vtManager.GetAllReleases()
	for i, release := range releases {
		if release != expectedReleases[i] {
			t.Errorf("Actual: %+v | Expected: %+v", release, expectedReleases[i])
		}
	}
}

func TestGetLatestRelease(t *testing.T) {
	// Setup
	githubClient := &mockGithubClient{}
	vtManager := NewVTManager(githubClient)
	expectedRelease := Release{
		Version:     "1.1.1",
		DownloadUrl: "http://localhost:9876/dw/vuetorrent-111.zip",
	}

	// Run
	release, _ := vtManager.GetLatestRelease()
	if release != expectedRelease {
		t.Errorf("Actual: %+v | Expected: %+v", release, expectedRelease)
	}
}

func TestGetReleaseByTag(t *testing.T) {
	// Setu
	githubClient := &mockGithubClient{}
	vtManager := NewVTManager(githubClient)
	expectedRelease := Release{
		Version:     "1.1.2",
		DownloadUrl: "http://localhost:9876/dw/vuetorrent.zip",
	}

	// Run
	release, _ := vtManager.GetReleaseByTag("1.1.2")
	if release != expectedRelease {
		t.Errorf("Actual: %+v | Expected: %+v", release, expectedRelease)
	}
}
