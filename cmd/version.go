package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const NISHIMURA_VERSION = "0.2.0"

var version = &cobra.Command{
	Use:   "version",
	Short: "print version of Nishimura",
	Long: `
Nishimura has a changing versions. This means during the work on Nishimura,
the version number will change.
I look forward to be down compatible with all NISHIMURA_VERSIONS.
AUTHOR:
	Anton Feldmann <anton.feldmann@gmail.com>
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Nishimura's Version is: ", NISHIMURA_VERSION)
	},
}

func GetVersion() string {
	return NISHIMURA_VERSION
}

func init() {
	Nishimura.AddCommand(version)
}
