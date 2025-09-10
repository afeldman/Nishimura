package gitutil

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/afeldman/Makoto/kpc"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func FetchRequirement(name, version string, repo kpc.Repository) error {
	req := kpc.Requirement{
		Name:    name,
		Version: version,
		Source:  repo,
	}
	return Fetch(req)
}

// Fetch lädt ein Requirement aus Git in den lokalen Cache (~/.nishimura/src).
// Falls es bereits existiert, wird fetch + checkout ausgeführt.
func Fetch(req kpc.Requirement) error {
	if req.Source.URL == "" {
		return fmt.Errorf("requirement %s@%s has no source URL", req.Name, req.Version)
	}

	home, _ := os.UserHomeDir()
	cacheDir := filepath.Join(home, ".nishimura", "src")
	target := filepath.Join(cacheDir, fmt.Sprintf("%s-%s", req.Name, req.Version))

	var repo *git.Repository
	var err error

	if _, statErr := os.Stat(target); statErr == nil {
		// Repo existiert → öffnen und fetch
		repo, err = git.PlainOpen(target)
		if err != nil {
			return fmt.Errorf("could not open repo: %w", err)
		}

		err = repo.Fetch(&git.FetchOptions{
			RemoteName: "origin",
			Tags:       git.AllTags,
			Force:      true,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("git fetch failed: %w", err)
		}
	} else {
		// sicherstellen, dass cacheDir existiert
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			return err
		}

		fmt.Printf("📥 cloning %s into %s\n", req.Source.URL, target)
		repo, err = git.PlainClone(target, false, &git.CloneOptions{
			URL:      req.Source.URL,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("git clone failed: %w", err)
		}
	}

	// Checkout: Tag > Rev > Branch
	ref := ""
	switch {
	case req.Source.Tag != "":
		ref = req.Source.Tag
	case req.Source.Rev != "":
		ref = req.Source.Rev
	case req.Source.Branch != "":
		ref = req.Source.Branch
	}

	if ref != "" {
		w, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("could not get worktree: %w", err)
		}

		var checkout plumbing.ReferenceName
		if req.Source.Branch != "" {
			checkout = plumbing.NewBranchReferenceName(req.Source.Branch)
		} else if req.Source.Tag != "" {
			checkout = plumbing.NewTagReferenceName(req.Source.Tag)
		}

		if checkout != "" {
			err = w.Checkout(&git.CheckoutOptions{
				Branch: checkout,
				Force:  true,
			})
		} else {
			// Rev (Commit Hash)
			err = w.Checkout(&git.CheckoutOptions{
				Hash:  plumbing.NewHash(ref),
				Force: true,
			})
		}

		if err != nil {
			return fmt.Errorf("git checkout %s failed for %s: %w", ref, req.Name, err)
		}
	}

	return nil
}
