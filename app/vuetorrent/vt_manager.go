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
	GetLatestVuetorrentRelease() (VueTorrentRelease, error)
	GetVuetorrentRelease(tag string) (VueTorrentRelease, error)
	GetAllReleases() ([]VueTorrentRelease, error)
	Install(release VueTorrentRelease, outputDir string) error
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

func (mng *vtManager) Install(release VueTorrentRelease, outputDir string) error {
	log.Printf("[INFO] Start downloading %v", release)
	filePath, err := download(release, os.TempDir())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Downloaded release into %s", filePath)

	var backupedDir, backupErr = backupPreviousVersion(outputDir)

	err = unzip(filePath, outputDir, release.Version)
	if err != nil {
		return err
	}

	if backupErr == nil {
		os.RemoveAll(backupedDir)
		log.Printf("[INFO] Removed old dir %s", backupedDir)
	}

	err = createVersionFile(release.Version, outputDir)
	if err != nil {
		log.Printf("[WARN] Can't create version file. Error: %s", err.Error())
	}

	return nil
}

func download(release VueTorrentRelease, outputDir string) (filePath string, err error) {
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

func unzip(filePath string, outputDir string, version string) error {
	log.Printf("[INFO] Extracting %s into %s \n", filePath, outputDir)

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

	for _, file := range archive.File {
		fileName, _ := strings.CutPrefix(file.Name, "vuetorrent/")
		if fileName == "" {
			continue
		}

		filePath = filepath.Join(filepath.Clean(outputDir), fileName)

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

	return nil
}

func createVersionFile(version string, outputDir string) error {
	filePath := filepath.Join(filepath.Clean(outputDir), "version.txt")

	if _, err := os.Stat(filePath); err == nil {
		log.Printf("[INFO] %s already exists", filePath)
		return nil
	}

	log.Println("[INFO] Creating missed version.txt file")
	versionData := []byte(version)
	err := os.WriteFile(filePath, versionData, 0777)
	if err != nil {
		return err
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
