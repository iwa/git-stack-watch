package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
)

// enum ChangeType
type ChangeType string

const (
	Created ChangeType = "created"
	Updated ChangeType = "updated"
	Deleted ChangeType = "deleted"
)

type Change struct {
	StackName  string
	FilePath   string
	ChangeType ChangeType
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: git-stack-watch <repository-path>")
		fmt.Println("Example: git-stack-watch /path/to/repo")
		os.Exit(1)
	}

	repoPath := os.Args[1]

	// Open the git repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatalf("Failed to open repository: %v", err)
	}

	// Get the worktree
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatalf("Failed to get worktree: %v", err)
	}

	// Get the current status
	status, err := worktree.Status()
	if err != nil {
		log.Fatalf("Failed to get status: %v", err)
	}

	// Find all compose file changes
	changes := findComposeChanges(status)

	if len(changes) == 0 {
		fmt.Println("No compose file changes detected.")
		return
	}

	fmt.Printf("Found %d stack change(s):\n", len(changes))
	for _, change := range changes {
		fmt.Printf("  - %s %s (%s)\n", change.ChangeType, change.StackName, change.FilePath)
	}

	// Create a commit for each stack change
	for _, change := range changes {
		fmt.Println()
		err := commitStackChange(worktree, repo, change)
		if err != nil {
			log.Printf("Failed to commit %s: %v", change.StackName, err)
			continue
		}
	}
}

// findComposeChanges scans the git status for compose.yml/compose.yaml changes
func findComposeChanges(status git.Status) []Change {
	var changes []Change

	for filePath, fileStatus := range status {
		// Check if the file is a compose file
		fileName := filepath.Base(filePath)
		if fileName != "compose.yml" && fileName != "compose.yaml" {
			continue
		}

		// Determine the stack name (parent directory name)
		stackName := getStackName(filePath)

		// Determine the change type
		var changeType ChangeType
		switch {
		case fileStatus.Staging == git.Added || fileStatus.Worktree == git.Untracked:
			changeType = "created"
		case fileStatus.Staging == git.Deleted || fileStatus.Worktree == git.Deleted:
			changeType = "deleted"
		case fileStatus.Staging == git.Modified || fileStatus.Worktree == git.Modified:
			changeType = "updated"
		default:
			// Skip if no relevant change
			continue
		}

		changes = append(changes, Change{
			StackName:  stackName,
			FilePath:   filePath,
			ChangeType: changeType,
		})
	}

	return changes
}

// getStackName extracts the stack name from the file path
// For example: "docker/komodo/compose.yml" -> "komodo"
func getStackName(filePath string) string {
	dir := filepath.Dir(filePath)
	// Get the last directory component
	stackName := filepath.Base(dir)

	// If the stack is in root, use the parent directory name
	if stackName == "." || stackName == "/" {
		stackName = "root"
	}

	return stackName
}

// commitStackChange creates a commit for a single stack change
func commitStackChange(worktree *git.Worktree, repo *git.Repository, change Change) error {
	if change.ChangeType == "deleted" {
		_, err := worktree.Remove(change.FilePath)
		if err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}
	} else {
		_, err := worktree.Add(change.FilePath)
		if err != nil {
			return fmt.Errorf("failed to add file: %w", err)
		}
	}

	// Create the commit
	commitMsg := fmt.Sprintf("%s %s", change.ChangeType, change.StackName)

	commit, err := worktree.Commit(commitMsg, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	// Log the commit hash
	log.Printf("✓ Created commit %s: %s", commit.String()[:7], commitMsg)

	return nil
}
