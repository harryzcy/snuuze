package platform

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeout = 10 * time.Second
)

type GiteaClient struct {
	client http.Client
	host   string
	token  string
}

func NewGiteaClient(host string) Client {
	token := getGiteaToken(host)
	return &GiteaClient{
		client: http.Client{
			Timeout: defaultTimeout,
		},
		host:  host,
		token: token,
	}
}

func getGiteaToken(host string) string {
	for _, giteaConfig := range getGiteaConfigs() {
		if giteaConfig.GetHost() == host {
			return giteaConfig.Token
		}
	}

	return ""
}

func (c *GiteaClient) sendRequest(method, path string, body []byte) ([]byte, error) {
	reader := bytes.NewReader(body)
	req, err := http.NewRequest(method, path, reader)
	if err != nil {
		return nil, err
	}

	if method == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", "bearer "+c.token)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type giteaTagList []struct {
	Commit struct {
		Created string `json:"created"`
		Sha     string `json:"sha"`
		URL     string `json:"url"`
	} `json:"commit"`
	ID         string `json:"id"`
	Message    string `json:"message"`
	Name       string `json:"name"`
	TarballURL string `json:"tarball_url"`
	ZipballURL string `json:"zipball_url"`
}

func (c *GiteaClient) ListTags(params *ListTagsInput) ([]string, error) {
	// TODO: support pagination
	url := c.host + "/api/v1/repos/" + params.Owner + "/" + params.Repo + "/tags"
	data, err := c.sendRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var tags giteaTagList
	err = json.Unmarshal(data, &tags)
	if err != nil {
		return nil, err
	}

	tagNames := make([]string, 0)
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	return tagNames, nil
}

type giteaCreatePullRequestOptions struct {
	Assignee  string   `json:"assignee"`
	Assignees []string `json:"assignees"`
	Base      string   `json:"base"`
	Body      string   `json:"body"`
	DueDate   string   `json:"due_date"`
	Head      string   `json:"head"`
	Labels    []int64  `json:"labels"`
	Milestone int64    `json:"milestone"`
	Title     string   `json:"title"`
}

func (c *GiteaClient) CreatePullRequest(input *CreatePullRequestInput) error {
	options := &giteaCreatePullRequestOptions{
		Base:  input.Base,
		Head:  input.Head,
		Title: input.Title,
		Body:  input.Body,
	}

	url := c.host + "/api/v1/repos/" + input.Owner + "/" + input.Repo + "/pulls"

	body, err := json.Marshal(options)
	if err != nil {
		return err
	}

	_, err = c.sendRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	return nil
}
