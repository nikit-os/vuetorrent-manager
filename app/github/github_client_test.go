package github

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestGetReleases(t *testing.T) {
	server := mockServerWithResponce(t, "testdata/releases.json")
	defer server.Close()

	githubClient := createGithubClient(server)

	releases, err := githubClient.GetReleases()
	if err != nil {
		t.Error(err.Error())
	}

	expectedReleases := []Release{
		{
			TagName: "v2.3.0",
			Assets: []Asset{
				{
					Name:        "vuetorrent.zip",
					DownloadUrl: "https://github.com/WDaan/VueTorrent/releases/download/v2.3.0/vuetorrent.zip",
				},
			},
		},
		{
			TagName: "v2.2.0",
			Assets: []Asset{
				{
					Name:        "vuetorrent.zip",
					DownloadUrl: "https://github.com/WDaan/VueTorrent/releases/download/v2.2.0/vuetorrent.zip",
				},
			},
		},
		{
			TagName: "v2.1.1",
			Assets: []Asset{
				{
					Name:        "vuetorrent.zip",
					DownloadUrl: "https://github.com/WDaan/VueTorrent/releases/download/v2.1.1/vuetorrent.zip",
				},
			},
		},
	}

	for i, receivedRelease := range releases {
		if !reflect.DeepEqual(receivedRelease, expectedReleases[i]) {
			t.Errorf("\nGot: %+v \nExp: %+v", receivedRelease, expectedReleases[i])
		}
	}
}

func TestGetReleaseByTag(t *testing.T) {
	server := mockServerWithResponce(t, "testdata/release_by_tag.json")
	defer server.Close()

	githubClient := createGithubClient(server)

	receivedRelease, err := githubClient.GetReleaseByTag("v2.3.0")
	if err != nil {
		t.Error(err.Error())
	}

	expectedRelease := Release{
		TagName: "v2.3.0",
		Assets: []Asset{
			{
				Name:        "vuetorrent.zip",
				DownloadUrl: "https://github.com/WDaan/VueTorrent/releases/download/v2.3.0/vuetorrent.zip",
			},
		},
	}

	if !reflect.DeepEqual(receivedRelease, expectedRelease) {
		t.Errorf("\nGot: %+v \nExp: %+v", receivedRelease, expectedRelease)
	}
}

func mockServerWithResponce(t *testing.T, fileWithResponce string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := readFileContent(t, fileWithResponce)
		w.Write(resp)
	}))
}

func createGithubClient(server *httptest.Server) Client {
	return &DefaultClient{
		ApiKey:  "foo",
		Client:  server.Client(),
		BaseUrl: server.URL,
	}
}

func readFileContent(t *testing.T, filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err.Error())
	}

	return bytes
}
