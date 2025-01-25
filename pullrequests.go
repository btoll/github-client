package githubclient

import (
	"fmt"

	"github.com/google/go-github/v68/github"
)

var upstreamBase string = "master"

func CreatePullRequest(g *GitHubClient, title, body string) (*github.PullRequest, error) {
	pr, _, err := client.PullRequests.Create(ctx, g.Upstream, g.Repository, &github.NewPullRequest{
		Title: &title,
		Head:  github.Ptr(fmt.Sprintf("%s:%s", g.Owner, g.Branch)),
		Base:  &upstreamBase,
		Body:  &body,
	})
	return pr, err
}
