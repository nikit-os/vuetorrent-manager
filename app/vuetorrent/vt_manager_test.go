package vuetorrent

import (
	"fmt"
	"n1kit0s/vt-manager/app/github"
	"os"
	"path/filepath"
	"testing"
)

type mockGithubClient struct {
}

func (c *mockGithubClient) GetReleases() ([]github.Release, error) {
	return []github.Release{
		{TagName: "v1.1.3", Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent-113.zip"}}},
		{TagName: "v1.1.2", Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent-112.zip"}}},
		{TagName: "v1.1.1", Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent-111.zip"}}},
	}, nil
}
func (c *mockGithubClient) GetReleaseByTag(tag string) (github.Release, error) {
	if tag == "v0.0.0" {
		return github.Release{}, fmt.Errorf("tag %s not found", tag)
	}
	return github.Release{
		TagName: tag, Assets: []github.Asset{{Name: "vuetorrent.zip", DownloadUrl: "http://localhost:9876/dw/vuetorrent.zip"}},
	}, nil
}

func TestGetAllReleases(t *testing.T) {
	// Setup
	githubClient := &mockGithubClient{}
	vtManager := NewVTManager(githubClient)
	expectedReleases := []Release{
		{Version: "1.1.3", DownloadUrl: "http://localhost:9876/dw/vuetorrent-113.zip"},
		{Version: "1.1.2", DownloadUrl: "http://localhost:9876/dw/vuetorrent-112.zip"},
		{Version: "1.1.1", DownloadUrl: "http://localhost:9876/dw/vuetorrent-111.zip"},
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
		Version:     "1.1.3",
		DownloadUrl: "http://localhost:9876/dw/vuetorrent-113.zip",
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

func TestGetReleaseForVersion(t *testing.T) {
	tests := map[string]struct {
		targetVersion   string
		expectedRelease Release
		expectedError   error
	}{
		"version specified and exists": {
			targetVersion: "1.1.1",
			expectedRelease: Release{
				Version:     "1.1.1",
				DownloadUrl: "http://localhost:9876/dw/vuetorrent.zip",
			},
			expectedError: nil,
		},
		"version specified but not exists": {
			targetVersion:   "0.0.0",
			expectedRelease: Release{},
			expectedError:   fmt.Errorf("tag v0.0.0 not found"),
		},
		"version has not specified": {
			targetVersion: "",
			expectedRelease: Release{
				Version:     "1.1.3",
				DownloadUrl: "http://localhost:9876/dw/vuetorrent-113.zip",
			},
			expectedError: nil,
		},
	}

	vtManager := vtManager{
		githubClient: &mockGithubClient{},
		downloader:   mockDownloader{},
		unzipper:     mockUnziper{},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualRelease, err := vtManager.GetReleaseForVersion(test.targetVersion)
			if err != nil && test.expectedError == nil {
				t.Fatalf("GetReleaseForVersion failed. Error: %s", err.Error())
			}

			if test.expectedError != nil && test.expectedError.Error() != err.Error() {
				t.Fatalf("GetReleaseForVersion error doesn't match. Expected: %s | Actual: %s", test.expectedError.Error(), err.Error())
			}

			if actualRelease != test.expectedRelease {
				t.Fatalf("Releases don't match. Expected: %s | Actual: %s", test.expectedRelease, actualRelease)
			}
		})
	}
}

func TestInstall(t *testing.T) {
	// Setup
	vtManager := vtManager{
		githubClient: &mockGithubClient{},
		downloader:   mockDownloader{},
		unzipper:     mockUnziper{},
	}

	expectedVersion := "1.1.1"

	tempDir := t.TempDir()
	outputDir := filepath.Join(tempDir, "vuetorrent")
	expectedVersionFilePath := filepath.Join(outputDir, "version.txt")

	// Run
	err := vtManager.Install(expectedVersion, outputDir)
	if err != nil {
		t.Fatalf("Installation failed. Error: %s", err.Error())
	}

	versionBytes, err := os.ReadFile(expectedVersionFilePath)
	if err != nil {
		t.Fatalf("Can't read %s file. Error: %s", expectedVersionFilePath, err.Error())
	}
	version := string(versionBytes)

	if version != expectedVersion {
		t.Fatalf("Version doesn't match. Expected %s | Actual %s", expectedVersion, version)
	}

}
