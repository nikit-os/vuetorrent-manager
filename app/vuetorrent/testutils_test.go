package vuetorrent

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type TestFile struct {
	Path    string
	Content string
}

func createZip(t *testing.T, files []TestFile) string {
	tempDir := t.TempDir()
	archivePath := filepath.Join(tempDir, "test_archive.zip")

	archive, err := os.Create(archivePath)
	if err != nil {
		t.Fatalf("Can't create test archive. Error: %s", err.Error())
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()

	for _, testFile := range files {
		reader := strings.NewReader(testFile.Content)
		writer, _ := zipWriter.Create(testFile.Path)

		_, err := io.Copy(writer, reader)
		if err != nil {
			t.Fatalf("Can't add test file [%s] to archive. Error: %s", testFile.Path, err.Error())
		}
	}

	return archivePath
}
