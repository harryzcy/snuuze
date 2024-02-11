package platform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitClient_ListTags(t *testing.T) {
	githubClient, _ := NewGitClient("https://github.com")
	tags, err := githubClient.ListTags(&ListTagsInput{
		Owner: "harryzcy",
		Repo:  "mailbox",
	})
	assert.Nil(t, err)
	assert.Contains(t, tags, "v1.0.0")

	giteaClient, _ := NewGitClient("https://gitea.com")
	tags, err = giteaClient.ListTags(&ListTagsInput{
		Owner: "harryzcy",
		Repo:  "act_runner",
	})
	assert.Nil(t, err)
	assert.Contains(t, tags, "v0.1.0")
}
