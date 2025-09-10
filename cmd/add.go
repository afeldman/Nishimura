package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/afeldman/Makoto/kpc"
	"github.com/afeldman/Makoto/makoto"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name>@<version> <git-url>",
	Short: "Add a dependency to the project",
	Long: `Add a dependency to the local Nishimura manifest (default: nishimura.kpc)
and fetch it into ~/.nishimura/src if not already present.
Also registers the package in the Makoto database.

Example:
  nishimura add motion_lib@1.2.0 https://github.com/yourname/motion_lib.git`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return cmd.Usage()
		}

		// split name@version
		parts := strings.Split(args[0], "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format: use name@version")
		}
		pkgName, pkgVersion := parts[0], parts[1]
		pkgURL := args[1]

		// --- clone into global Nishimura cache (~/.nishimura/src) ---
		homedir, homeErr := os.UserHomeDir()
		if homeErr != nil {
			log.Fatal(homeErr)
		}
		cacheDir := filepath.Join(homedir, ".nishimura", "src")
		target := filepath.Join(cacheDir, fmt.Sprintf("%s-%s", pkgName, pkgVersion))

		var repo *git.Repository
		var err error

		if _, statErr := os.Stat(target); os.IsNotExist(statErr) {
			if err := os.MkdirAll(cacheDir, 0755); err != nil {
				return err
			}
			fmt.Printf("📥 cloning %s@%s into %s\n", pkgName, pkgVersion, target)
			repo, err = git.PlainClone(target, false, &git.CloneOptions{
				URL:      pkgURL,
				Progress: os.Stdout,
			})
			if err != nil {
				return fmt.Errorf("git clone failed: %w", err)
			}
		} else {
			fmt.Printf("ℹ️ %s@%s already present in %s\n", pkgName, pkgVersion, target)
			repo, err = git.PlainOpen(target)
			if err != nil {
				return fmt.Errorf("could not open repo: %w", err)
			}
		}

		// --- checkout version (Tag bevorzugt) ---
		w, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("cannot get worktree: %w", err)
		}

		// versuchen Tag zu checken
		hash, err := repo.ResolveRevision(plumbing.Revision(pkgVersion))
		if err == nil {
			err = w.Checkout(&git.CheckoutOptions{
				Hash:  *hash,
				Force: true,
			})
			if err != nil {
				return fmt.Errorf("checkout %s failed: %w", pkgVersion, err)
			}
		}

		// --- register in Makoto DB ---
		kpcFile := filepath.Join(target, "nishimura.kpc")
		if _, err := os.Stat(kpcFile); err == nil {
			k, err := kpc.ReadKPCFile(kpcFile)
			if err != nil {
				return fmt.Errorf("cannot read %s: %w", kpcFile, err)
			}

			m := makoto.InitMakoto()
			m.DBInit()

			if err := makoto.Append(k); err != nil {
				return fmt.Errorf("failed to register in Makoto DB: %w", err)
			}
			fmt.Printf("📦 registered %s@%s in Makoto DB\n", k.Name, k.Version)
		}

		// --- update project manifest (default: nishimura.kpc) ---
		var manifest kpc.KPC
		if _, err := os.Stat(manifestFile); err == nil {
			if _, err := toml.DecodeFile(manifestFile, &manifest); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("no %s found, run 'nishimura init' first", manifestFile)
		}

		// add or update dependency
		req := kpc.Requirement{
			Name:    pkgName,
			Version: pkgVersion,
			Source:  kpc.Repository{URL: pkgURL, Tag: pkgVersion},
		}
		found := false
		for i, dep := range manifest.Requirements {
			if dep.Name == pkgName {
				manifest.Requirements[i] = req
				found = true
				break
			}
		}
		if !found {
			manifest.Requirements = append(manifest.Requirements, req)
		}

		// write manifest back
		f, err := os.Create(manifestFile)
		if err != nil {
			return err
		}
		defer f.Close()
		if err := toml.NewEncoder(f).Encode(manifest); err != nil {
			return err
		}

		fmt.Printf("✅ added %s@%s (%s) to project (%s)\n", pkgName, pkgVersion, pkgURL, manifestFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(
		&manifestFile,
		"file",
		"f",
		"nishimura.kpc",
		"Path to the Nishimura manifest (default: nishimura.kpc)",
	)
}
