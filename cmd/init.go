package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/afeldman/Makoto/kpc"
	"github.com/spf13/cobra"
)

var manifestFile string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Nishimura project",
	Long: `Create a new Nishimura project in the current directory.
It generates a manifest file (default: nishimura.kpc) with basic fields
(name, version, authors, dependencies, conflicts).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		projectName := filepath.Base(wd)

		// Manifest über kpc.KPC erzeugen
		manifest := kpc.InitKPC(projectName)
		manifest.Version = "0.1.0"

		// prüfen, ob Datei schon existiert
		if _, err := os.Stat(manifestFile); err == nil {
			return fmt.Errorf("%s already exists", manifestFile)
		}

		// schreiben
		f, err := os.Create(manifestFile)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := toml.NewEncoder(f).Encode(manifest); err != nil {
			return err
		}

		fmt.Printf("✅ Initialized new Nishimura project: %s (%s)\n", projectName, manifestFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVarP(
		&manifestFile,
		"file",
		"f",
		"nishimura.kpc",
		"Path to the Nishimura manifest (default: nishimura.kpc)",
	)
}
