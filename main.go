package main

import (
	"fmt"
	"os"

	"github.com/craigbrett17/go-git-buggy/clients"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func main() {
	gitClient := clients.GitClient{}

	fmt.Println("Cloning the repo")
	_, err := gitClient.Clone("temp", "craigbrett17", "go-git-buggy")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Cloned successfully")

	repository, err := gitClient.LoadFromDirectory("temp")
	if err != nil {
		fmt.Println(err)
		return
	}

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

	// echo out useful info for people looking for the problem
	fmt.Println("You'll notice that a file is modified that hasn't been touched.")
	fmt.Println("Try running `git status` in the temp directory to see the same thing.")
	fmt.Println("Run git diff --cached to see the differences.")
	fmt.Println("When done, remove the temp directory before running again.")
}
