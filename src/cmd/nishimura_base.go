package cmd

import (
	"log"
	"path"

	homedir "github.com/atrox/homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/afeldman/go-util/env"
)

var (
	Nishimura = &cobra.Command{
		Use:   "Nishimura",
		Short: "Nishimura is a FANUC Karel package configurator",
		Long: `
	   	  Nishimura Karel package manager
	   	  ===============================
Nishimura is a Karel package manager. The idea is that a Karel program
can put in an all containing package. This makes it more flexible to write
Karel Code and share the code with other users.
AUTHOR:
	Anton Feldmann <anton.feldmann@gmail.com>
\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\`,
	}

	confFile string
)

func init() {
	cobra.OnInitialize(initConfig)

	Nishimura.PersistentFlags().StringVar(&confFile, "config", "", "config file (default $HOME/.config/nishimura/nishimura.yaml)")

	Nishimura.AddCommand(version)
	Nishimura.AddCommand(init)
}

func Execute() {
	Nishimura.Execute()
}

func initConfig() {
	if confFile != "" {
		viper.SetConfigFile(confFile)
	} else {

		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		confFile = path.Join(home, ".config", "nishimura", "nishimura.yaml")

		viper.AddConfigPath(path.Join(home, ".config", "nishimura"))
		viper.SetConfigName("nishimura")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
	} else {
		if err := viper.Unmarshal(&rfg); err != nil {
			log.Fatal("unable to decode into the Nishimura configuration structure, %v", err)
		}
	}

	if len(rfg.RootDir) > 0 {
		rfg.init(env.GetEnv("NISHIMURA_PATH"))
	}
	rfg.save(confFile)

}
