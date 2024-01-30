package vuetorrent

import (
	"fmt"
	"log/slog"
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
	GetReleaseForVersion(version string) (Release, error)
	Install(version string, outputDir string) error
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

func (mng *vtManager) Install(targetVersion string, outputDir string) error {
	release, err := mng.GetReleaseForVersion(targetVersion)
	if err != nil {
		return err
	}

	installedVersion, _ := GetInstalledVersion(outputDir)

	slog.Info(fmt.Sprintf("Installed version: %s. Target version: %s", installedVersion, release.Version))

	if installedVersion == release.Version {
		slog.Info(fmt.Sprintf("Version %s already installed. Abort installation", release.Version))
		return nil
	}

	slog.Info("Start downloading", "release", release)
	cleanedOutputDir := filepath.Clean(outputDir)
	filePath, err := mng.downloader.Download(release, os.TempDir())
	if err != nil {
		return err
	}
	slog.Info("Downloaded release", "downloadPath", filePath)

	var backupedDir, backupErr = backupPreviousVersion(cleanedOutputDir)

	err = mng.unzipper.Unzip(filePath, cleanedOutputDir)
	if err != nil {
		return err
	}

	if backupErr == nil {
		os.RemoveAll(backupedDir)
		slog.Info("Removed old dir", "dir", backupedDir)
	}

	err = createVersionFile(release.Version, cleanedOutputDir)
	if err != nil {
		slog.Warn("Can't create version file", "error", err.Error())
	}

	return nil
}

func (mng *vtManager) GetReleaseForVersion(version string) (Release, error) {
	var vtRelease Release

	if version == "" {
		release, err := mng.GetLatestRelease()
		if err != nil {
			return Release{}, err
		}
		vtRelease = release
	} else {
		tag := MakeTagName(version)
		release, err := mng.GetReleaseByTag(tag)
		if err != nil {
			return Release{}, err
		}
		vtRelease = release
	}

	return vtRelease, nil
}

func createVersionFile(version string, outputDir string) error {
	filePath := filepath.Join(filepath.Clean(outputDir), "version.txt")

	if _, err := os.Stat(filePath); err == nil {
		slog.Info("Version file already exists", "file", filePath)
		return nil
	}

	slog.Info("Creating missed version.txt file")
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

func GetInstalledVersion(vtDirectory string) (string, error) {
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
		previousVersion, err := GetInstalledVersion(outputDir)
		if err != nil {
			slog.Warn("Previous version is unknown", "error", err.Error())
		}
		backupedDir = fmt.Sprintf(outputDir + "-" + previousVersion)
		slog.Info("Renaming old output directory", "renamedDir", backupedDir)
		os.Rename(outputDir, backupedDir)
	}

	return backupedDir, err
}
