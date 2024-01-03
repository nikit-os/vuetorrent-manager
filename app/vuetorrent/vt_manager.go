package vuetorrent

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"n1kit0s/vt-manager/app/github"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type VueTorrentRelease struct {
	Version     string
	DownloadUrl string
}

type VTManager interface {
	Download(release VueTorrentRelease, outputDir string) (filePath string, err error)
	GetLatestVuetorrentRelease() (VueTorrentRelease, error)
	GetVuetorrentRelease(tag string) (VueTorrentRelease, error)
	GetAllReleases() ([]VueTorrentRelease, error)
	Unzip(filePath string, outputDir string, version string) (err error)
}

type vtManager struct {
	githubClient github.Client
}

func NewVTManager(githubClient github.Client) VTManager {
	return &vtManager{
		githubClient: githubClient,
	}
}

func convertToVuetorrentRelease(githubRelease github.Release) VueTorrentRelease {
	var version, _ = strings.CutPrefix(githubRelease.TagName, "v")

	var downloadUrl string
	for _, asset := range githubRelease.Assets {
		if asset.Name == "vuetorrent.zip" {
			downloadUrl = asset.DownloadUrl
			break
		}
	}

	return VueTorrentRelease{
		Version:     version,
		DownloadUrl: downloadUrl,
	}
}

func (mng *vtManager) GetVuetorrentRelease(tag string) (VueTorrentRelease, error) {
	githubRelease, err := mng.githubClient.GetReleaseByTag(tag)
	if err != nil {
		return VueTorrentRelease{}, err
	}

	vtRelease := convertToVuetorrentRelease(githubRelease)

	return vtRelease, nil
}

func (mng *vtManager) GetLatestVuetorrentRelease() (VueTorrentRelease, error) {
	githubReleases, err := mng.githubClient.GetReleases()
	if err != nil {
		return VueTorrentRelease{}, err
	}

	var latestRelease = githubReleases[0]
	vtRelease := convertToVuetorrentRelease(latestRelease)

	return vtRelease, nil
}

func (mng *vtManager) GetAllReleases() ([]VueTorrentRelease, error) {
	githubReleases, err := mng.githubClient.GetReleases()
	if err != nil {
		return []VueTorrentRelease{}, err
	}

	var vtReleases []VueTorrentRelease

	for _, githubRelease := range githubReleases {
		vtRelease := convertToVuetorrentRelease(githubRelease)
		vtReleases = append(vtReleases, vtRelease)
	}

	return vtReleases, nil
}

func (mng *vtManager) Download(release VueTorrentRelease, outputDir string) (filePath string, err error) {
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

func (mng *vtManager) Unzip(filePath string, outputDir string, version string) error {
	log.Printf("[INFO] Extracting %s into %s \n", filePath, outputDir)

	var backupedDir, backupErr = backupPreviousVersion(outputDir)

	_, err := os.Open(outputDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[INFO] Output direcrory %s doesn't exists. Creating...", outputDir)
			os.MkdirAll(outputDir, os.ModePerm)
		}
	}

	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	var versionFileExists = false
	for _, file := range archive.File {
		fileName, _ := strings.CutPrefix(file.Name, "vuetorrent/")
		if fileName == "" {
			continue
		}

		filePath = filepath.Join(filepath.Clean(outputDir), fileName)
		if !versionFileExists && fileName == "version.txt" {
			versionFileExists = true
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
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

	if !versionFileExists {
		log.Println("[INFO] Creating missed version.txt file")
		filePath = filepath.Join(filepath.Clean(outputDir), "version.txt")
		versionData := []byte(version)
		err := os.WriteFile(filePath, versionData, 0777)
		if err != nil {
			return err
		}

	}

	if backupErr == nil {
		os.RemoveAll(backupedDir)
		log.Printf("[INFO] Removed old dir %s", backupedDir)
	}

	return nil
}

func MakeTagName(version string) string {
	if strings.HasPrefix(version, "v") {
		return version
	} else {
		return fmt.Sprintf("v%s", version)
	}
}

func GetVersion(vtDirectory string) (string, error) {
	var versionFilePath = path.Join(vtDirectory, "version.txt")
	_, err := os.Stat(versionFilePath)
	if err != nil {
		return "unknown", err
	}

	fileBytes, err := os.ReadFile(versionFilePath)
	if err != nil {
		return "unknown", err
	}

	return string(fileBytes), nil
}

func backupPreviousVersion(outputDir string) (string, error) {
	_, err := os.Stat(outputDir)
	var backupedDir = ""
	if err == nil {
		previousVersion, err := GetVersion(outputDir)
		if err != nil {
			log.Printf("[WARN] Previous version is unknown. Err: %s", err.Error())
		}
		backupedDir = fmt.Sprintf(outputDir + "-" + previousVersion)
		log.Printf("[INFO] Renaming old output dir to %s", backupedDir)
		os.Rename(outputDir, backupedDir)
	}

	return backupedDir, err
}
