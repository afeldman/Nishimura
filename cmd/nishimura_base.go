package cmd

import (
	"log"
	"github.com/spf13/cobra"
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
	Nishimura.AddCommand(build)
}

func Execute() {
	if err := Nishimura.Execute(); err != nil {
		log.Fatal(err)
	}

}

func initConfig(){
	ncft, err := loadConfig(confFile)
	log.Println(ncft.RootDir)
	if err != nil {
		ncft.build_file()
	}
}
