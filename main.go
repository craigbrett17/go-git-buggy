package main

import (
	"fmt"
	"os"

	"github.com/craigbrett17/go-git-buggy/clients"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	token := ""
	// replace token with the contents of GITHUB_TOKEN if in the env
	if os.Getenv("GITHUB_TOKEN") != "" {
		token = os.Getenv("GITHUB_TOKEN")
	}

	gitClient := clients.GitClient{
		Token: token,
	}

	fmt.Println("Cloning the repo")
	repository, err := gitClient.Clone("temp", "craigbrett17", "go-git-buggy")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Cloned successfully")

	// create a new branch
	fmt.Println("Creating a new branch")
	workTree, err := repository.Worktree()
	if err != nil {
		fmt.Println(err)
		return
	}

	branchName := "test-branch"
	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Branch created and checked out successfully")

	// write a new test file and add it to the repo
	fmt.Println("Creating a new file")
	file, err := os.Create("temp/test_files/test_3.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	file.WriteString("This is a test file")

	_, err = workTree.Add("test_files/test_3.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	// print the status at this point
	status, err := workTree.Status()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(status)
}
