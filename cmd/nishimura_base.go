package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	confFile   string
	configRoot *NishimuraConfig

	Nishimura = &cobra.Command{
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

	Nishimura.PersistentFlags().StringVar(
		&confFile,
		"config",
		DefaultConfPath(),
		"path to config file (default: ~/.nishimura/config.toml)",
	)
}

func Execute() {
	if err := Nishimura.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	// Initialisiere globale Config
	configRoot = InitNishimura(confFile)

	if confFile != "" {
		viper.SetConfigFile(confFile)
	} else {
		viper.AddConfigPath(configRoot.RootDir)
		viper.SetConfigName("nishimura") // ohne Suffix
	}
	viper.AutomaticEnv()
}
