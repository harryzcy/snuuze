package auth

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harryzcy/snuuze/config"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func GitHubPATClient(token string) *githubv4.Client {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	httpClient := oauth2.NewClient(context.Background(), src)
	return githubv4.NewClient(httpClient)
}

func GithubAppInstallationClient() (*githubv4.Client, error) {
	privateKeyFile := config.GetHostingConfig().GitHub.PEMFile
	privateKey, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}

	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		Issuer:    config.GetHostingConfig().GitHub.AppID,
	}

	privateKeyParsed, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	signedToken, err := token.SignedString(privateKeyParsed)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %v", err)
	}

	ctx := context.Background()
	src := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: signedToken,
	})
	httpClient := oauth2.NewClient(ctx, src)

	client := githubv4.NewClient(httpClient)

	return client, nil
}
