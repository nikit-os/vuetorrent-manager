package vuetorrent

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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
		Version: "1.2",
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
		Version: "1.2",
		DownloadUrl: "shoul not be used",
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
