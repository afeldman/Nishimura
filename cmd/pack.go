package cmd

import (
	"log"
	"io/ioutil"
	"strings"
	"github.com/spf13/cobra"
	"path/filepath"

	kpc "github.com/afeldman/Makoto/kpc"
	krl "github.com/afeldman/Nishimura/karel"
)

var (
	pack = &cobra.Command{
		Use:   "pack [project_path]",
		Short: "pack creates a karel package",
		Long: `

AUTHOR:
	Anton Feldmann <anton.feldmann@gmail.com>
`,
		Run: func(cmd *cobra.Command, args []string) {
			PackPackage(args[0])
		},
	}
)

func PackPackage(path string){

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(strings.TrimSpace(file.Name())) == ".kpc" {

			data, file_err := ioutil.ReadFile(filepath.Join(path,file.Name()))
			if file_err != nil {
				log.Fatalln(err)
			}

			if kpc_err, kpc_ := kpc.FromYAML(data); kpc_err != nil {
				log.Fatalln(kpc_err)
			}else{
				krl.BuildKarel(*kpc_, path)
			}
			break
		}
	}

}
