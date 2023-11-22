package auth

import (
	"context"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/harryzcy/snuuze/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func GitHubPATClient() (client *githubv4.Client, token string) {
	token = config.GitHubToken()
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return githubv4.NewClient(httpClient), token
}

func GithubAppInstallationClient() (*githubv4.Client, *ghinstallation.Transport, error) {
	conf := config.GetHostingConfig()
	appID := conf.GitHub.AppID
	installationID := conf.GitHub.InstallationID
	privateKeyFile := conf.GitHub.PEMFile

	tr := http.DefaultTransport
	transport, err := ghinstallation.NewKeyFromFile(tr, appID, installationID, privateKeyFile)
	if err != nil {
		return nil, nil, err
	}

	return githubv4.NewClient(&http.Client{Transport: transport}), transport, nil
}
