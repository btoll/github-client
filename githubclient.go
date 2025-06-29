package githubclient

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v68/github"
)

var (
	client *github.Client
	ctx    context.Context = context.Background()
)

type GitHubClient struct {
	Upstream   string
	Owner      string
	Repository string
	Branch     string
	Path       string
	Message    string
	Paginate   bool
}

func init() {
	apiToken, isSet := os.LookupEnv("GITHUB_TOKEN")
	if apiToken == "" || !isSet {
		panic("[ERROR] Must set $GITHUB_TOKEN!")
	}
	client = github.NewClient(nil).WithAuthToken(apiToken)
}

func GetPullRequestTemplate(g *GitHubClient, sha string) ([]byte, error) {
	tree, err := GetTree(g, sha)
	if err != nil {
		return nil, err
	}
	f := ".github/PULL_REQUEST_TEMPLATE.md"
	for _, entry := range tree.Entries {
		if strings.Contains(*entry.Path, f) {
			return GetBlobContent(g, *entry.SHA)
		}
	}
	return nil, errors.New(fmt.Sprintf("[ERROR] `%s` could not be found.\n", f))
}
