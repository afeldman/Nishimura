package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/afeldman/Makoto/kpc"
	"github.com/afeldman/Makoto/makoto"
	"github.com/afeldman/Nishimura/gitutil"
	"github.com/afeldman/Nishimura/nishimura"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install all dependencies for the current project",
	Long:  `Reads 'nishimura.kpc' and installs all dependencies listed under [[requirements]].`,
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := nishimura.ResolveDependencies("nishimura.kpc")
		if err != nil {
			return err
		}

		if len(root.Deps) == 0 {
			fmt.Println("ℹ️ no dependencies to install")
			return nil
		}

		// Init DB
		m := makoto.InitMakoto()
		m.DBInit()

		fmt.Println("📦 Installing dependencies:")
		for _, dep := range root.Deps {
			fmt.Printf("  - %s@%s\n", dep.Name, dep.Version)

			// Git holen
			req := kpc.Requirement{
				Name:    dep.Name,
				Version: dep.Version,
				Source:  dep.Repo,
			}
			if err := gitutil.Fetch(req); err != nil {
				return err
			}

			// KPC ins DB übernehmen, falls vorhanden
			home, _ := os.UserHomeDir()
			cachePath := filepath.Join(home, ".nishimura", "src",
				fmt.Sprintf("%s-%s", dep.Name, dep.Version), "nishimura.kpc")

			if _, err := os.Stat(cachePath); err == nil {
				k, err := kpc.ReadKPCFile(cachePath)
				if err == nil {
					if err := makoto.Append(k); err != nil {
						return fmt.Errorf("DB append failed for %s@%s: %w", dep.Name, dep.Version, err)
					}
					fmt.Printf("    ↳ registered in DB\n")
				}
			}
		}

		fmt.Println("✅ install finished")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
