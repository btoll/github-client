package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/go-github/v41/github"
)

type RepositoryService struct {
	*BaseService
}

type Repository struct {
	Name          *string  `json:"name,omitempty" yaml:"name,omitempty"`
	Owner         *string  `json:"owner,omitempty" yaml:"owner,omitempty"`
	TplName       *string  `json:"tpl_name,omitempty" yaml:"tpl_name,omitempty"`
	Archived      bool     `json:"archived,omitempty" yaml:"archived,omitempty"`
	Collaborators []string `json:"collaborators,omitempty" yaml:"collaborators,omitempty"`
}

func NewRepositoryService(p *Provisioner) *RepositoryService {
	return &RepositoryService{
		BaseService: &BaseService{
			provisioner: p,
		},
	}
}

type RepositoryCtx struct {
	Organization *github.Organization
	Repository   Repository
}

func (r *RepositoryService) AddCollaborators(rc *RepositoryCtx, repository *github.Repository) {
	for _, collaborator := range rc.Repository.Collaborators {
		_, _, err := r.provisioner.Client.Repositories.AddCollaborator(getContext(), *rc.Organization.Login, *rc.Repository.Name, collaborator, &github.RepositoryAddCollaboratorOptions{
			Permission: "pull",
		})
		if err != nil {
			fmt.Printf("[ERROR] Collaborator `%s` could not be added to repository `%s` could not be updated -- %s\n", collaborator, *rc.Repository.Name, err)
		} else {
			fmt.Printf("[SUCCESS] Added collaborator `%s` to repository `%s`\n", collaborator, *rc.Repository.Name)
		}
	}
}

// For our purposes, always create from a template.
func (r *RepositoryService) Create(rc *RepositoryCtx) {
	if rc.Repository.TplName == nil {
		fmt.Printf("[ERROR] Repository `%s` could not be created, no template name\n", *rc.Repository.Name)
		return
	}
	err := r.IsTemplateRepository(rc)
	if err != nil {
		fmt.Println(err)
		return
	}
	repo, _, err := r.provisioner.Client.Repositories.CreateFromTemplate(getContext(), r.GetTemplateRepositoryOwner(rc), *rc.Repository.TplName, &github.TemplateRepoRequest{
		Name:  rc.Repository.Name,
		Owner: rc.Organization.Login,
	})
	if err != nil {
		fmt.Printf("[ERROR] Repository `%s` could not be created -- %s\n", *rc.Repository.Name, err)
	} else {
		fmt.Printf("[SUCCESS] Created repository `%s` from template repository `%s`.\n", *rc.Repository.Name, *rc.Repository.TplName)

		if rc.Repository.Archived {
			b := true
			repo.Archived = &b
			// Slight delay to ensure the repo has been created before doing more operations.
			time.Sleep(300 * time.Millisecond)
			r.Update(rc, repo)
		}

		if len(rc.Repository.Collaborators) > 0 {
			r.AddCollaborators(rc, repo)
		}
	}
}

func (r *RepositoryService) Delete(rc *RepositoryCtx) {
	_, err := r.provisioner.Client.Repositories.Delete(getContext(), *rc.Organization.Login, *rc.Repository.Name)
	if err != nil {
		fmt.Printf("[ERROR] Repository `%s` could not be deleted -- %s\n", *rc.Repository.Name, err)
	} else {
		fmt.Printf("[SUCCESS] Deleted repository `%s`.\n", *rc.Repository.Name)
	}
}

func (r *RepositoryService) GetTemplateRepositoryOwner(rc *RepositoryCtx) string {
	if rc.Repository.Owner != nil {
		return *rc.Repository.Owner
	}
	return *rc.Organization.Login
}

func (r *RepositoryService) IsTemplateRepository(rc *RepositoryCtx) error {
	repo, _, err := r.provisioner.Client.Repositories.Get(getContext(), r.GetTemplateRepositoryOwner(rc), *rc.Repository.TplName)
	if err != nil {
		return errors.New(fmt.Sprintf("[ERROR] There was a problem accessing repository `%s` -- %s", *rc.Repository.Name, err))
	} else {
		if !*repo.IsTemplate {
			return errors.New(fmt.Sprintf("[ERROR] The repository `%s` is not a template repository.", *rc.Repository.Name))
		}
	}
	return nil
}

func (r *RepositoryService) Update(rc *RepositoryCtx, repository *github.Repository) {
	_, _, err := r.provisioner.Client.Repositories.Edit(getContext(), *rc.Organization.Login, *rc.Repository.Name, repository)
	if err != nil {
		fmt.Printf("[ERROR] Repository `%s` could not be updated -- %s\n", *rc.Repository.Name, err)
	} else {
		fmt.Printf("[SUCCESS] Updated repository `%s`\n", *rc.Repository.Name)
	}
}
