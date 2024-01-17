package vuetorrent

import (
	"n1kit0s/vt-manager/app/github"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
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
	expectedReleases := []VueTorrentRelease{
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

func TestGetLatestVuetorrentRelease(t *testing.T) {
	// Setup
	githubClient := &mockGithubClient{}
	vtManager := NewVTManager(githubClient)
	expectedRelease := VueTorrentRelease{
		Version:     "1.1.1",
		DownloadUrl: "http://localhost:9876/dw/vuetorrent-111.zip",
	}

	// Run
	release, _ := vtManager.GetLatestVuetorrentRelease()
	if release != expectedRelease {
		t.Errorf("Actual: %+v | Expected: %+v", release, expectedRelease)
	}
}

func TestGetVuetorrentRelease(t *testing.T) {
	// Setu
	githubClient := &mockGithubClient{}
	vtManager := NewVTManager(githubClient)
	expectedRelease := VueTorrentRelease{
		Version:     "1.1.2",
		DownloadUrl: "http://localhost:9876/dw/vuetorrent.zip",
	}

	// Run
	release, _ := vtManager.GetVuetorrentRelease("1.1.2")
	if release != expectedRelease {
		t.Errorf("Actual: %+v | Expected: %+v", release, expectedRelease)
	}
}

func TestUnzip(t *testing.T) {
	// Setup
	files := []TestFile{
		{Path: "vuetorrent/version.txt", Content: "1.2.3"},
		{Path: "vuetorrent/folder1/file1.js", Content: "file1 content"},
		{Path: "vuetorrent/folder1/file2.js", Content: "file2 content"},
		{Path: "vuetorrent/folder2/file3.js", Content: "file3 content"},
		{Path: "vuetorrent/folder2/file4.js", Content: "file4 content"},
	}

	archivePath := createZip(t, files)
	t.Logf("Created archive [%s]", archivePath)

	outputDir := filepath.Join(t.TempDir(), "test_out")

	// Run
	err := unzip(archivePath, outputDir, "0.0.0")
	if err != nil {
		t.Fatalf("Failed to unzip. Error: %s", err.Error())
	}

	for _, file := range files {
		expectedFilePath, _ := strings.CutPrefix(file.Path, "vuetorrent/")
		_, err := os.Stat(filepath.Join(outputDir, expectedFilePath))
		if err != nil {
			t.Errorf("File [%s] was not unarchived into [%s]", expectedFilePath, outputDir)
		}
	}
}

func TestDownload(t *testing.T) {
	// Setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("file content"))
	}))
	defer server.Close()

	vtRelease := VueTorrentRelease{
		Version:     "1.2",
		DownloadUrl: server.URL,
	}
	outputDir := t.TempDir()

	// Run
	downloadedFilePath, err := download(vtRelease, outputDir)
	if err != nil {
		t.Fatalf("download failed. error: %s", err.Error())
	}

	fileInfo, err := os.Stat(downloadedFilePath)
	if err != nil {
		t.Fatalf("File was not downloaded into [%s]. Error: %s", outputDir, err.Error())
	}

	expectedFilename := "vuetorrent-1.2.zip"
	if fileInfo.Name() != expectedFilename {
		t.Fatalf("Expected filename is [%s]. Actual filename is [%s]", expectedFilename, fileInfo.Name())
	}
}

func TestDownloadShouldSkipIfFileExists(t *testing.T) {
	// Setup
	vtRelease := VueTorrentRelease{
		Version:     "1.2",
		DownloadUrl: "should not be used",
	}
	outputDir := t.TempDir()
	expectedFilename := "vuetorrent-1.2.zip"
	os.Create(filepath.Join(outputDir, expectedFilename))

	// Run
	_, err := download(vtRelease, outputDir)
	if err != nil {
		t.Fatalf("download failed. error: %s", err.Error())
	}

}
