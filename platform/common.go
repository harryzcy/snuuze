package platform

type Client interface {
	// ListTags returns a sorted list of tags for the given repo
	ListTags(params *ListTagsInput) ([]string, error)

	CreatePullRequest(input *CreatePullRequestInput) error
}

type ListTagsInput struct {
	Owner  string
	Repo   string
	Prefix string // optional
}

type CreatePullRequestInput struct {
	Title string
	Body  string
	Base  string
	Head  string
	Owner string
	Repo  string
}
