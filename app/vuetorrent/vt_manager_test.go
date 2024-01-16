package vuetorrent

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnzip(t *testing.T) {
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
