package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/afeldman/Gakutensoku/ktrans"
	"github.com/afeldman/Makoto/kpc"
	"github.com/afeldman/Nishimura/nishimura"
	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compile the project via Gakutensoku ktrans wrapper",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Manifest laden
		var manifest kpc.KPC
		if _, err := toml.DecodeFile("nishimura.kpc", &manifest); err != nil {
			return fmt.Errorf("cannot read nishimura.kpc: %w", err)
		}

		if manifest.Main == "" {
			return fmt.Errorf("no main file defined in nishimura.kpc")
		}

		// Build-Verzeichnis anlegen
		wd, _ := os.Getwd()
		buildDir := filepath.Join(wd, "build")
		if err := os.MkdirAll(buildDir, 0755); err != nil {
			return err
		}

		// Dependency-Include-Pfade sammeln
		root, err := nishimura.ResolveDependencies("nishimura.kpc")
		if err != nil {
			return err
		}

		includePaths := []string{}
		home, _ := os.UserHomeDir()
		for _, dep := range root.Deps {
			cacheManifest := filepath.Join(
				home,
				".nishimura",
				"src",
				fmt.Sprintf("%s-%s", dep.Name, dep.Version),
				"nishimura.kpc",
			)
			if _, err := os.Stat(cacheManifest); err == nil {
				var depManifest kpc.KPC
				if _, err := toml.DecodeFile(cacheManifest, &depManifest); err == nil {
					if depManifest.IncludeDir != "" {
						includePaths = append(includePaths, filepath.Join(
							filepath.Dir(cacheManifest),
							depManifest.IncludeDir,
						))
					}
				}
			}
		}

		// Output-Datei bestimmen (.pc im build Ordner)
		outputFile := filepath.Join(buildDir,
			strings.TrimSuffix(filepath.Base(manifest.Main), filepath.Ext(manifest.Main))+".pc",
		)

		// ktrans über Wrapper aufrufen
		fmt.Println("Compiling using ktrans...")

		err = ktrans.Run(ktrans.Options{
			Main:        manifest.Main,
			Output:      outputFile,
			Prefix:      manifest.Prefix,
			IncludeDirs: includePaths,
		})
		if err != nil {
			return fmt.Errorf("ktrans compile failed: %w", err)
		}

		fmt.Printf("✅ compiled %s → %s\n", manifest.Main, outputFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
