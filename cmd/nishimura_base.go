package cmd

import (
	"log"

	"github.com/afeldman/Nishimura/nishimura"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	confFile   string
	configRoot *nishimura.NishimuraConfig

	rootCmd = &cobra.Command{
		Use:   "nishimura",
		Short: "Nishimura is a FANUC Karel package manager",
		Long: `
Nishimura Karel package manager
===============================

Nishimura is a Karel package manager. The idea is that a Karel program
can be bundled in a self-contained package. This makes it easier to write,
share and reuse Karel code.

AUTHOR:
  Anton Feldmann <anton.feldmann@gmail.com>
`,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&confFile,
		"config",
		"",
		"path to config file (default: ~/.nishimura/config.toml)",
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if confFile == "" {
		confFile = nishimura.DefaultConfPath()
	}

	// Initialisiere globale Config
	configRoot = nishimura.InitNishimura(confFile)

	viper.SetConfigFile(confFile)
	viper.AutomaticEnv()
}
