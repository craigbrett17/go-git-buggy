package clients

import (
	"fmt"
	"os"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/storage/filesystem"
)

type ExperimentalFuncs struct {
}

type GitClient struct {
	Experimental ExperimentalFuncs
}

func (c *GitClient) Clone(dir string, organisation string, repo string) (*git.Repository, error) {
	repositoryName := fmt.Sprintf("https://github.com/%s/%s.git", organisation, repo)
	repository, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:      repositoryName,
		Progress: os.Stdout,
		Depth:    1,
	})

	return repository, err
}

func (c *GitClient) LoadFromDirectory(dir string) (*git.Repository, error) {
	fs := osfs.New(dir, osfs.WithBoundOS())
	return git.PlainOpen(fs.Root())
}

func (c *GitClient) CreateBranch(worktree *git.Worktree, checkoutHash plumbing.Hash, branchName string, force bool) error {
	return worktree.Checkout(&git.CheckoutOptions{
		Hash:   checkoutHash,
		Branch: plumbing.ReferenceName("refs/heads/" + branchName),
		Create: true,
		Force:  force,
	})
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

func (e ExperimentalFuncs) Clone(dir string, organisation string, repo string) (*git.Repository, error) {
	fmt.Println("Using experimental clone")
	repositoryName := fmt.Sprintf("https://github.com/%s/%s.git", organisation, repo)
	wt := osfs.New(dir, osfs.WithBoundOS())
	dotfs, err := wt.Chroot(git.GitDirName)
	if err != nil {
		return nil, err
	}
	store := filesystem.NewStorage(dotfs, cache.NewObjectLRUDefault())

	repository, err := git.Clone(store, wt, &git.CloneOptions{
		URL:      repositoryName,
		Progress: os.Stdout,
		Depth:    1,
	})

	return repository, err
}

func (e ExperimentalFuncs) LoadFromDirectory(dir string) (*git.Repository, error) {
	fmt.Println("Using experimental load from directory")
	wt := osfs.New(dir, osfs.WithBoundOS())
	dotfs, err := wt.Chroot(git.GitDirName)
	if err != nil {
		return nil, err
	}
	store := filesystem.NewStorage(dotfs, cache.NewObjectLRUDefault())
	return git.Open(store, wt)
}
