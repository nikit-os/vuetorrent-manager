package vuetorrent

import (
	"fmt"
	"log"
	"n1kit0s/vt-manager/app/github"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Release struct {
	Version     string
	DownloadUrl string
}

type VTManager interface {
	GetLatestRelease() (Release, error)
	GetReleaseByTag(tag string) (Release, error)
	GetAllReleases() ([]Release, error)
	Install(release Release, outputDir string) error
}

type vtManager struct {
	githubClient github.Client
	unzipper     Unzipper
	downloader   Downloader
}

func NewVTManager(githubClient github.Client) VTManager {
	return &vtManager{
		githubClient: githubClient,
		unzipper:     DefaultUnzipper{},
		downloader:   HttpDownloader{},
	}
}

func convertToVuetorrentRelease(githubRelease github.Release) Release {
	var version, _ = strings.CutPrefix(githubRelease.TagName, "v")

	var downloadUrl string
	for _, asset := range githubRelease.Assets {
		if asset.Name == "vuetorrent.zip" {
			downloadUrl = asset.DownloadUrl
			break
		}
	}

	return Release{
		Version:     version,
		DownloadUrl: downloadUrl,
	}
}

func (mng *vtManager) GetReleaseByTag(tag string) (Release, error) {
	githubRelease, err := mng.githubClient.GetReleaseByTag(tag)
	if err != nil {
		return Release{}, err
	}

	vtRelease := convertToVuetorrentRelease(githubRelease)

	return vtRelease, nil
}

func (mng *vtManager) GetLatestRelease() (Release, error) {
	githubReleases, err := mng.githubClient.GetReleases()
	if err != nil {
		return Release{}, err
	}

	var latestRelease = githubReleases[0]
	vtRelease := convertToVuetorrentRelease(latestRelease)

	return vtRelease, nil
}

func (mng *vtManager) GetAllReleases() ([]Release, error) {
	githubReleases, err := mng.githubClient.GetReleases()
	if err != nil {
		return []Release{}, err
	}

	var vtReleases []Release

	for _, githubRelease := range githubReleases {
		vtRelease := convertToVuetorrentRelease(githubRelease)
		vtReleases = append(vtReleases, vtRelease)
	}

	return vtReleases, nil
}

func (mng *vtManager) Install(release Release, outputDir string) error {
	log.Printf("[INFO] Start downloading %v", release)
	cleanedOutputDir := filepath.Clean(outputDir)
	filePath, err := mng.downloader.Download(release, os.TempDir())
	if err != nil {
		return err
	}
	log.Printf("[INFO] Downloaded release into %s", filePath)

	var backupedDir, backupErr = backupPreviousVersion(cleanedOutputDir)

	err = mng.unzipper.Unzip(filePath, cleanedOutputDir, release.Version)
	if err != nil {
		return err
	}

	if backupErr == nil {
		os.RemoveAll(backupedDir)
		log.Printf("[INFO] Removed old dir %s", backupedDir)
	}

	err = createVersionFile(release.Version, cleanedOutputDir)
	if err != nil {
		log.Printf("[WARN] Can't create version file. Error: %s", err.Error())
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
