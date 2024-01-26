package vuetorrent

import (
	"n1kit0s/vt-manager/app/github"
	"os"
	"path/filepath"
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
	// Setup
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

type mockDownloader struct{}

func (m mockDownloader) Download(release Release, outputDir string) (filePath string, err error) {
	return "/some/file/path", nil
}

type mockUnziper struct{}

func (m mockUnziper) Unzip(filePath string, outputDir string) error {
	return os.MkdirAll(outputDir, os.ModePerm)
}

func TestInstall(t *testing.T) {
	// Setup
	vtManager := vtManager{
		githubClient: &mockGithubClient{},
		downloader:   mockDownloader{},
		unzipper:     mockUnziper{},
	}

	targetRelease := Release{
		Version:     "1.1.1",
		DownloadUrl: "http://localhost:9876/dw/vuetorrent-111.zip",
	}

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "vuetorrent")
	expectedVersionFilePath := filepath.Join(outputDir, "version.txt")

	// Run
	err := vtManager.Install(targetRelease, outputDir)
	if err != nil {
		t.Fatalf("Installation failed. Error: %s", err.Error())
	}

	versionBytes, err := os.ReadFile(expectedVersionFilePath)
	if err != nil {
		t.Fatalf("Can't read %s file. Error: %s", expectedVersionFilePath, err.Error())
	}
	version := string(versionBytes)

	if version != targetRelease.Version {
		t.Fatalf("Version doesn't match. Expected %s | Actual %s", targetRelease.Version, version)
	}

}
