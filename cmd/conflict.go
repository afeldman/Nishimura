package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/afeldman/Makoto/kpc"
	"github.com/afeldman/Nishimura/nishimura"
	"github.com/spf13/cobra"
)

// Conflict root command
var conflictCmd = &cobra.Command{
	Use:   "conflict",
	Short: "Manage conflicts for the project",
	Long: `Conflicts describe which packages are incompatible with this project.
They are stored in the project manifest (default: nishimura.kpc).`,
}

// Add a new conflict
var conflictAddCmd = &cobra.Command{
	Use:   "add <name>@<version>",
	Short: "Add a conflict to the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		parts := strings.Split(args[0], "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format: use name@version")
		}
		pkgName, pkgVersion := parts[0], parts[1]

		var manifest kpc.KPC
		if _, err := toml.DecodeFile(manifestFile, &manifest); err != nil {
			return fmt.Errorf("cannot read %s: %w", manifestFile, err)
		}

		// prüfen ob es schon einen Conflict-Eintrag für dieses Package gibt
		found := false
		for i, c := range manifest.Conflicts {
			if c.Name == pkgName {
				// prüfen ob Version schon drin ist
				for _, v := range c.Versions {
					if v == pkgVersion {
						fmt.Printf("ℹ️ conflict %s@%s already listed\n", pkgName, pkgVersion)
						return nil
					}
				}
				// Version hinzufügen
				manifest.Conflicts[i].Versions = append(manifest.Conflicts[i].Versions, pkgVersion)
				found = true
				break
			}
		}

		if !found {
			manifest.Conflicts = append(manifest.Conflicts, kpc.Conflict{
				Name:     pkgName,
				Versions: []string{pkgVersion},
			})
		}

		return writeManifest(&manifest)
	},
}

// Remove a conflict
var conflictRmCmd = &cobra.Command{
	Use:   "rm <name>@<version>",
	Short: "Remove a conflict from the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return cmd.Usage()
		}

		parts := strings.Split(args[0], "@")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format: use name@version")
		}
		pkgName, pkgVersion := parts[0], parts[1]

		var manifest kpc.KPC
		if _, err := toml.DecodeFile(manifestFile, &manifest); err != nil {
			return fmt.Errorf("cannot read %s: %w", manifestFile, err)
		}

		removed := false
		newConflicts := []kpc.Conflict{}

		for _, c := range manifest.Conflicts {
			if c.Name == pkgName {
				// Versionsliste neu bauen ohne pkgVersion
				newVers := []string{}
				for _, v := range c.Versions {
					if v == pkgVersion {
						removed = true
						continue
					}
					newVers = append(newVers, v)
				}
				if len(newVers) > 0 {
					newConflicts = append(newConflicts, kpc.Conflict{Name: c.Name, Versions: newVers})
				}
			} else {
				newConflicts = append(newConflicts, c)
			}
		}

		manifest.Conflicts = newConflicts

		if !removed {
			fmt.Printf("ℹ️ no conflict %s@%s found\n", pkgName, pkgVersion)
			return nil
		}

		return writeManifest(&manifest)
	},
}

// List all conflicts
var conflictListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all conflicts in the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		var manifest kpc.KPC
		if _, err := toml.DecodeFile(manifestFile, &manifest); err != nil {
			return fmt.Errorf("cannot read %s: %w", manifestFile, err)
		}

		if len(manifest.Conflicts) == 0 {
			fmt.Println("ℹ️ no conflicts defined")
			return nil
		}

		for _, c := range manifest.Conflicts {
			for _, v := range c.Versions {
				fmt.Printf("- %s@%s\n", c.Name, v)
			}
		}
		return nil
	},
}

// Check all conflicts against current requirements
var conflictCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Check recursively if requirements conflict with project conflicts",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := nishimura.ResolveDependencies("nishimura.kpc")
		if err != nil {
			return err
		}

		// Root-Manifest nochmal parsen für die Conflicts
		var manifest kpc.KPC
		if _, err := toml.DecodeFile("nishimura.kpc", &manifest); err != nil {
			return fmt.Errorf("cannot read manifest: %w", err)
		}

		all := nishimura.Flatten(root)
		hasConflict := false

		for _, dep := range all[1:] { // Root überspringen
			parts := strings.Split(dep, "@")
			name, ver := parts[0], parts[1]

			for _, c := range manifest.Conflicts {
				if c.Name != name {
					continue
				}
				for _, v := range c.Versions {
					if v == ver {
						fmt.Printf("❌ conflict: requirement %s@%s is incompatible\n", name, ver)
						hasConflict = true
					}
				}
			}
		}

		if hasConflict {
			return fmt.Errorf("conflicts detected")
		}

		fmt.Println("✅ no conflicts with current requirements")
		return nil
	},
}

// Helper: write manifest back to file
func writeManifest(manifest *kpc.KPC) error {
	f, err := os.Create(manifestFile)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(manifest); err != nil {
		return err
	}

	fmt.Printf("✅ updated %s\n", manifestFile)
	return nil
}

func init() {
	rootCmd.AddCommand(conflictCmd)
	conflictCmd.AddCommand(conflictAddCmd)
	conflictCmd.AddCommand(conflictRmCmd)
	conflictCmd.AddCommand(conflictListCmd)
	conflictCmd.AddCommand(conflictCheckCmd)
}
