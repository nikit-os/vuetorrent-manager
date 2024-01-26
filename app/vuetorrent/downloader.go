package vuetorrent

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader interface {
	Download(release Release, outputDir string) (filePath string, err error)
}

type HttpDownloader struct{}

func (d HttpDownloader) Download(release Release, outputDir string) (filePath string, err error) {
	var filename = fmt.Sprintf("vuetorrent-%s.zip", release.Version)
	filePath = filepath.Join(outputDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		log.Printf("[INFO] %s already exists here %s. skipping download", filename, filePath)
		return filePath, nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	resp, err := http.Get(release.DownloadUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
