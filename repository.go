package githubclient

import (
	"github.com/google/go-github/v68/github"
)

func CreateFork(g *GitHubClient) (*github.Repository, error) {
	repo, _, err := client.Repositories.CreateFork(ctx, g.Upstream, g.Repository, nil)
	return repo, err
}

func CreateRepository(g *GitHubClient) (*github.Repository, error) {
	repo, _, err := client.Repositories.Create(ctx, g.Owner, &github.Repository{
		Name: &g.Repository,
	})
	return repo, err
}

func DeleteRepository(g *GitHubClient) error {
	_, err := client.Repositories.Delete(ctx, g.Owner, g.Repository)
	return err
}

func GetBranch(g *GitHubClient, branch string) (*github.Branch, error) {
	b, _, err := client.Repositories.GetBranch(ctx, g.Owner, g.Repository, branch, 3)
	return b, err
}

func GetCommit(g *GitHubClient, sha string) (*github.RepositoryCommit, error) {
	commit, _, err := client.Repositories.GetCommit(ctx, g.Owner, g.Repository, sha, nil)
	return commit, err
}

func ListByOrg(g *GitHubClient) ([]*github.Repository, error) {
	options := &github.RepositoryListByOrgOptions{
		Type:      "all",
		Sort:      "full_name",
		Direction: "asc",
	}
	var repos []*github.Repository
	for true {
		page, resp, err := client.Repositories.ListByOrg(ctx, g.Upstream, options)
		if err != nil {
			// Print something here so it's known an error occurred.
			return repos, err
		}
		repos = append(repos, page...)
		if !g.Paginate || resp.NextPage == 0 {
			break
		}
		options.Page += 1
	}
	return repos, nil
}

func ListForks(g *GitHubClient) ([]*github.Repository, error) {
	repos, _, err := client.Repositories.ListForks(ctx, g.Owner, g.Repository, nil)
	return repos, err
}

func ListRepositories(g *GitHubClient) ([]*github.Repository, error) {
	//	https://api.github.com/users/btoll/repos
	repos, _, err := client.Repositories.ListByUser(ctx, g.Owner, nil)
	return repos, err
}
