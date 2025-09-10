package cmd

import (
	"fmt"

	"github.com/afeldman/Nishimura/nishimura"
	"github.com/spf13/cobra"
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Show dependency graph",
	RunE: func(cmd *cobra.Command, args []string) error {
		root, err := nishimura.ResolveDependencies("nishimura.kpc")
		if err != nil {
			return err
		}

		var printTree func(*nishimura.Node, string)
		printTree = func(n *nishimura.Node, indent string) {
			line := fmt.Sprintf("%s- %s@%s", indent, n.Name, n.Version)
			if n.Repo.URL != "" {
				line += fmt.Sprintf(" (%s)", n.Repo.URL)
			}
			fmt.Println(line)

			for _, d := range n.Deps {
				printTree(d, indent+"  ")
			}
		}

		printTree(root, "")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
