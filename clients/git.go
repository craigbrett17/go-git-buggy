package clients

import (
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

type GitClient struct {
	Token string
}

func (c *GitClient) auth() transport.AuthMethod {
	return &http.BasicAuth{
		Username: "anything",
		Password: c.Token,
	}
}

func (c *GitClient) Clone(dir string, organisation string, repo string) (*git.Repository, error) {
	repositoryName := fmt.Sprintf("https://github.com/%s/%s.git", organisation, repo)
	repository, err := git.PlainClone(dir, false, &git.CloneOptions{
		Auth:     c.auth(),
		URL:      repositoryName,
		Progress: os.Stdout,
		Depth:    1,
	})

	return repository, err
}

func (c *GitClient) Commit(workTree *git.Worktree, commitMessage string) (plumbing.Hash, error) {
	return workTree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Author McAuthorface",
			Email: "author@github.com",
			When:  time.Now(),
		},
	})
}

func (c *GitClient) GetFromDirectory(dir string) (*git.Repository, error) {
	return git.PlainOpen(dir)
}

func (c *GitClient) GetBranchName(repository *git.Repository) (string, error) {
	head, err := repository.Head()
	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}

func (c *GitClient) CreateBranch(worktree *git.Worktree, checkoutHash plumbing.Hash, branchName string, force bool) error {
	return worktree.Checkout(&git.CheckoutOptions{
		Hash:   checkoutHash,
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
		Create: true,
		Force:  force,
	})
}

func (c *GitClient) DeleteBranchRemote(repository *git.Repository, branchName string) error {
	return repository.Push(
		&git.PushOptions{
			Auth:     c.auth(),
			RefSpecs: []config.RefSpec{config.RefSpec(":refs/heads/" + branchName)},
			Progress: os.Stdout,
		})

}

func (c *GitClient) Push(repository *git.Repository) error {
	return repository.Push(&git.PushOptions{
		Auth: c.auth(),
	})
}

func (c *GitClient) CheckoutBranch(repository *git.Repository, branchName string, force bool) error {
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
		Force:  force,
	})
}

func (c *GitClient) DeleteBranch(repository *git.Repository, branchName string) error {
	branchRef := plumbing.NewBranchReferenceName(branchName)
	return repository.Storer.RemoveReference(branchRef)
}

func (c *GitClient) AddFile(repository *git.Repository, path string) error {
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}

	_, err = worktree.Add(path)
	return err
}

func (c *GitClient) GetStatus(repository *git.Repository) (git.Status, error) {
	worktree, err := repository.Worktree()
	if err != nil {
		return nil, err
	}

	return worktree.Status()
}
