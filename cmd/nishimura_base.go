package cmd

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	Nishimura.PersistentFlags().StringVar(&confFile, "config", DefaultConfPath(), "the config file to set the Path.")

	Nishimura.AddCommand(version)
	Nishimura.AddCommand(build)
}

func Execute() {
	if err := Nishimura.Execute(); err != nil {
		log.Fatal(err)
	}

}

func initConfig(){
	ncft.initNishimura(confFile)
	ncft.build_file()

	if confFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(confFile)
	} else {

		viper.AddConfigPath(ncft.RootDir)
		viper.SetConfigName(ncft.ConfFile)
	}
	viper.AutomaticEnv() // read in environment variables that match
}
