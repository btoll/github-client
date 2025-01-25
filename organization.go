package githubclient

import (
	"github.com/google/go-github/v68/github"
)

func GetOrganization(g *GitHubClient) (*github.Organization, error) {
	org, _, err := client.Organizations.Get(ctx, g.Upstream)
	return org, err
}

func ListInstallations(g *GitHubClient) (*github.OrganizationInstallations, error) {
	orgInstallations, _, err := client.Organizations.ListInstallations(ctx, g.Upstream, nil)
	return orgInstallations, err
}
