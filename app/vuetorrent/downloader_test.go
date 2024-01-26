package vuetorrent

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestDownload(t *testing.T) {
	// Setup
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("file content"))
	}))
	defer server.Close()

	vtRelease := Release{
		Version:     "1.2",
		DownloadUrl: server.URL,
	}
	outputDir := t.TempDir()

	downloader := HttpDownloader{}

	// Run
	downloadedFilePath, err := downloader.Download(vtRelease, outputDir)
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
	vtRelease := Release{
		Version:     "1.2",
		DownloadUrl: "should not be used",
	}
	outputDir := t.TempDir()
	expectedFilename := "vuetorrent-1.2.zip"
	os.Create(filepath.Join(outputDir, expectedFilename))

	downloader := HttpDownloader{}

	// Run
	_, err := downloader.Download(vtRelease, outputDir)
	if err != nil {
		t.Fatalf("download failed. error: %s", err.Error())
	}
}
