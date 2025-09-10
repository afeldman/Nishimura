package gitutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/afeldman/Makoto/kpc"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneAndCheckout klont ein Repository und checkt Tag/Rev/Branch aus
func CloneAndCheckout(repo kpc.Repository, destDir string) error {
	if repo.URL == "" {
		return fmt.Errorf("repository URL is empty")
	}

	// Zielpfad bestimmen
	name := filepath.Base(repo.URL)
	if destDir == "" {
		destDir = ".nishimura"
	}
	target := filepath.Join(destDir, name)

	var repository *git.Repository
	var err error

	if _, statErr := os.Stat(target); statErr == nil {
		// Repo existiert schon → öffnen und fetch
		repository, err = git.PlainOpen(target)
		if err != nil {
			return fmt.Errorf("could not open repo: %w", err)
		}
		wt, err := repository.Worktree()
		if err != nil {
			return fmt.Errorf("could not get worktree: %w", err)
		}
		err = repository.Fetch(&git.FetchOptions{
			RemoteName: "origin",
			Tags:       git.AllTags,
			Force:      true,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("git fetch failed: %w", err)
		}
		_ = wt // optional, für späteren Checkout
	} else {
		// frisches Clone
		fmt.Printf("📥 cloning %s into %s\n", repo.URL, target)
		repository, err = git.PlainClone(target, false, &git.CloneOptions{
			URL:      repo.URL,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}
	}

	// Checkout vorbereiten
	ref := ""
	switch {
	case repo.Tag != "":
		ref = repo.Tag
	case repo.Rev != "":
		ref = repo.Rev
	case repo.Branch != "":
		ref = repo.Branch
	}

	if ref != "" {
		wt, err := repository.Worktree()
		if err != nil {
			return fmt.Errorf("could not get worktree: %w", err)
		}

		var hash plumbing.Hash
		var referenceName plumbing.ReferenceName

		// Tag
		if repo.Tag != "" {
			referenceName = plumbing.NewTagReferenceName(repo.Tag)
		}
		// Branch
		if repo.Branch != "" {
			referenceName = plumbing.NewBranchReferenceName(repo.Branch)
		}
		// Rev (direkt als Commit)
		if repo.Rev != "" {
			hash = plumbing.NewHash(repo.Rev)
		}

		if hash != (plumbing.Hash{}) {
			err = wt.Checkout(&git.CheckoutOptions{Hash: hash})
		} else if referenceName != "" {
			err = wt.Checkout(&git.CheckoutOptions{Branch: referenceName})
		}

		if err != nil {
			return fmt.Errorf("git checkout %s failed: %w", ref, err)
		}
	}

	return nil
}
