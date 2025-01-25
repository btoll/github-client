package githubclient

import (
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v68/github"
)

func CreateBlob(g *GitHubClient, content string) (*github.Blob, error) {
	blob, _, err := client.Git.CreateBlob(ctx, g.Owner, g.Repository, &github.Blob{
		Content:  &content,
		Encoding: github.Ptr("utf-8"),
	})
	return blob, err
}

func CreateCommit(g *GitHubClient, parent *github.RepositoryCommit, tree *github.Tree) (*github.Commit, error) {
	commit, _, err := client.Git.CreateCommit(ctx, g.Owner, g.Repository, &github.Commit{
		//		SHA:     &head,
		Tree:    tree,
		Parents: []*github.Commit{parent.Commit},
		Message: github.Ptr("this is a foo"),
	}, &github.CreateCommitOptions{})
	return commit, err
}

func CreateReference(g *GitHubClient, baseRef *github.Reference) (*github.Reference, error) {
	//	newRef := &github.Reference{Ref: github.Ptr("refs/heads/" + *commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err := client.Git.CreateRef(ctx, g.Owner, g.Repository, &github.Reference{
		Ref: github.Ptr(fmt.Sprintf("refs/heads/%s", g.Branch)),
		Object: &github.GitObject{
			Type: github.Ptr("branch"),
			SHA:  baseRef.Object.SHA,
		},
	})
	return ref, err
}

func CreateTree(g *GitHubClient, sha string, entries []*github.TreeEntry) (*github.Tree, error) {
	tree, _, err := client.Git.CreateTree(ctx, g.Owner, g.Repository, sha, entries)
	return tree, err
}

func DeleteReference(g *GitHubClient) error {
	_, err := client.Git.DeleteRef(ctx, g.Owner, g.Repository, fmt.Sprintf("refs/heads/%s", g.Branch))
	return err
}

func GetBlob(g *GitHubClient, sha string) (*github.Blob, error) {
	blob, _, err := client.Git.GetBlob(ctx, g.Owner, g.Repository, sha)
	return blob, err
}

func GetBlobContent(g *GitHubClient, sha string) ([]byte, error) {
	blob, err := GetBlob(g, sha)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(*blob.Content)
}

func GetReference(g *GitHubClient) (*github.Reference, error) {
	ref, _, err := client.Git.GetRef(ctx, g.Owner, g.Repository, fmt.Sprintf("refs/heads/%s", g.Branch))
	return ref, err
}

func GetTree(g *GitHubClient, sha string) (*github.Tree, error) {
	tree, _, err := client.Git.GetTree(ctx, g.Owner, g.Repository, sha, true)
	return tree, err
}

func UpdateReference(g *GitHubClient, baseRef *github.Reference) (*github.Reference, error) {
	//	   _, _, err = client.Git.UpdateRef(ctx, *sourceOwner, *sourceRepo, ref, false)
	//	ref, _, err := client.Git.UpdateRef(ctx, g.Owner, g.Repository, &github.Reference{
	//		Ref: &head,
	//		Object: &github.GitObject{
	//			Type: github.Ptr("branch"),
	//			SHA:  &tree,
	//		},
	//	}, true)
	ref, _, err := client.Git.UpdateRef(ctx, g.Owner, g.Repository, baseRef, false)
	return ref, err
}
