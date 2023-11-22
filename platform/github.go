package platform

import (
	"context"
	"fmt"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/harryzcy/snuuze/config"
	"github.com/harryzcy/snuuze/platform/auth"
	"github.com/harryzcy/snuuze/types"
	"github.com/shurcooL/githubv4"
)

type GitHubClient struct {
	authType  string
	client    *githubv4.Client
	transport *ghinstallation.Transport // only used for GitHub App
	token     string                    // only used for personal access token
}

var _ Client = &GitHubClient{}

// NewGitHubClient creates a new GitHubClient authenticated with
// either a GitHub App or a personal access token.
func NewGitHubClient() (Client, error) {
	client := &GitHubClient{
		authType: config.GetHostingConfig().GitHub.AuthType,
	}

	if client.authType == types.AuthTypeGithubApp {
		var err error
		client.client, client.transport, err = auth.GithubAppInstallationClient()
		if err != nil {
			return nil, fmt.Errorf("failed to create GitHub client: %v", err)
		}
	} else if client.authType == types.AuthTypeToken {
		client.client, client.token = auth.GitHubPATClient()
	}

	return client, nil
}

func (c *GitHubClient) Token(ctx context.Context) (string, error) {
	if c.authType == types.AuthTypeGithubApp {
		return c.transport.Token(ctx)
	}
	return c.token, nil
}

func (c *GitHubClient) ListRepos() ([]Repo, error) {
	var query struct {
		Viewer struct {
			Repositories struct {
				Nodes []struct {
					Owner struct {
						Login string
					}
					Name             string
					URL              string
					IsPrivate        bool
					DefaultBranchRef struct {
						Name string
					}
				}
			} `graphql:"repositories(first: 100, isArchived: false, orderBy: {field: NAME, direction: ASC})"`
		}
	}

	err := c.client.Query(context.Background(), &query, nil)
	if err != nil {
		return nil, err
	}

	repos := make([]Repo, 0)
	for _, node := range query.Viewer.Repositories.Nodes {
		repos = append(repos, Repo{
			Server:        "github.com",
			Owner:         node.Owner.Login,
			Repo:          node.Name,
			URL:           node.URL,
			IsPrivate:     node.IsPrivate,
			DefaultBranch: node.DefaultBranchRef.Name,
		})
	}

	return repos, nil
}

func (c *GitHubClient) ListTags(params *ListTagsInput) ([]string, error) {
	var query struct {
		Repository struct {
			Refs struct {
				Edges []struct {
					Node struct {
						Name string
					}
				}
			} `graphql:"refs(refPrefix: $refPrefix, first: 100, orderBy: {field: TAG_COMMIT_DATE, direction: DESC})"`
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	refPrefix := "refs/tags/"
	if params.Prefix != "" {
		refPrefix += params.Prefix
		if !strings.HasSuffix(refPrefix, "/") {
			refPrefix += "/"
		}
	}
	variables := map[string]interface{}{
		"owner":     githubv4.String(params.Owner),
		"name":      githubv4.String(params.Repo),
		"refPrefix": githubv4.String(refPrefix),
	}

	err := c.client.Query(context.Background(), &query, variables)
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0)
	for _, edge := range query.Repository.Refs.Edges {
		tags = append(tags, edge.Node.Name)
	}

	return tags, nil
}

func (c *GitHubClient) CreatePullRequest(input *CreatePullRequestInput) error {
	repoID, err := c.lookupRepositoryID(input.Owner, input.Repo)
	if err != nil {
		return err
	}

	var mutation struct {
		CreatePullRequest struct {
			PullRequest struct {
				Permalink githubv4.URI
			}
		} `graphql:"createPullRequest(input: $input)"`
	}
	githubInput := githubv4.CreatePullRequestInput{
		RepositoryID: repoID,
		BaseRefName:  githubv4.String(input.Base),
		HeadRefName:  githubv4.String(input.Head),
		Title:        githubv4.String(input.Title),
		Body:         githubv4.NewString(githubv4.String(input.Body)),
	}

	err = c.client.Mutate(context.Background(), &mutation, githubInput, nil)
	if err != nil {
		if strings.Contains(err.Error(), "A pull request already exists") {
			fmt.Println("Pull request already exists")
			return nil
		}
		return err
	}
	fmt.Println("Pull request created: ", mutation.CreatePullRequest.PullRequest.Permalink)

	return nil
}

func (c *GitHubClient) lookupRepositoryID(owner, repo string) (string, error) {
	var query struct {
		Repository struct {
			ID string
		} `graphql:"repository(owner: $owner, name: $name)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(repo),
	}

	err := c.client.Query(context.Background(), &query, variables)
	if err != nil {
		return "", err
	}

	return query.Repository.ID, nil
}
