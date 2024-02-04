package vuetorrent

import (
	"archive/zip"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Unzipper interface {
	Unzip(filePath string, outputDir string) error
}

type DefaultUnzipper struct{}

func (u DefaultUnzipper) Unzip(filePath string, outputDir string) error {
	slog.Info(fmt.Sprintf("Extracting %s into %s", filePath, outputDir))

	_, err := os.Open(outputDir)
	if err != nil && os.IsNotExist(err) {
		slog.Info(fmt.Sprintf("Output direcrory %s doesn't exists. Creating...", outputDir))
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return err
		}
	}

	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		fileName, _ := strings.CutPrefix(file.Name, "vuetorrent/")
		if fileName == "" {
			continue
		}

		filePath = filepath.Join(filepath.Clean(outputDir), fileName)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := file.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
