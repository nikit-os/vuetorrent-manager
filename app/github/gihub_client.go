package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Asset struct {
	Name        string `json:"name"`
	DownloadUrl string `json:"browser_download_url"`
}

type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Client interface {
	GetReleases() ([]Release, error)
	GetReleaseByTag(tag string) (Release, error)
}

type DefaultClient struct {
	ApiKey  string
	Client  *http.Client
	BaseUrl string
}

func NewClient(apiKey string) Client {
	return &DefaultClient{
		ApiKey:  apiKey,
		Client:  &http.Client{},
		BaseUrl: "https://api.github.com",
	}
}

func (github *DefaultClient) GetReleases() ([]Release, error) {
	var releasesUrl = fmt.Sprintf("%s/%s", github.BaseUrl, "repos/WDaan/VueTorrent/releases")
	req, err := http.NewRequest("GET", releasesUrl, nil)
	if err != nil {
		return []Release{}, fmt.Errorf("failed to create releases request. %s", err.Error())
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", github.ApiKey))

	resp, err := github.Client.Do(req)
	if err != nil {
		return []Release{}, fmt.Errorf("failed to retrieve releases from github. %s", err.Error())
	}
	defer resp.Body.Close()

	releasesBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return []Release{}, fmt.Errorf("failed to read releases responce. %s", err.Error())
	}

	var githubReleases []Release
	err = json.Unmarshal(releasesBody, &githubReleases)
	if err != nil {
		return []Release{}, fmt.Errorf("failed to decode releases. response: [%s]. %s", string(releasesBody), err.Error())
	}

	return githubReleases, nil
}

func (github *DefaultClient) GetReleaseByTag(tag string) (Release, error) {
	var releasesUrl = fmt.Sprintf("%s/%s/%s", github.BaseUrl, "repos/WDaan/VueTorrent/releases/tags", tag)
	req, err := http.NewRequest("GET", releasesUrl, nil)
	if err != nil {
		return Release{}, fmt.Errorf("failed to create get release by tag request. %s", err.Error())
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", github.ApiKey))

	resp, err := github.Client.Do(req)
	if err != nil {
		return Release{}, fmt.Errorf("failed to retrieve release by tag from github. %s", err.Error())
	}
	defer resp.Body.Close()

	releaseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return Release{}, fmt.Errorf("failed to read release by tag responce. %s", err.Error())
	}

	var githubRelease Release
	err = json.Unmarshal(releaseBody, &githubRelease)
	if err != nil {
		return Release{}, fmt.Errorf("failed to decode release by tag. response: [%s]. %s", string(releaseBody), err.Error())
	}

	return githubRelease, nil
}
