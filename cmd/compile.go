package cmd

import (
	"github.com/spf13/cobra"
)

var (
	compile = &cobra.Command{
		Use:   "compile",
		Short: "compile the project",
		Long: `
Build the project using Gakutensoku

AUTHOR:
	Anton Feldmann <anton.feldmann@gmail.com>
`,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)
