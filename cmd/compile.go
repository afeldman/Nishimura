package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

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
			var kpc_file string
			var path string
			if len(args) == 0 {
				p, err := os.Getwd()
				if err != nil {
					log.Panic(err)
				}
				path = p
			} else {
				p := args[0]
				path = p
			}

			if _, err := os.Stat(path); os.IsNotExist(err) {
				log.Fatalln(err)
			}

			files := checkExt(path, ".kpc")

			if len(files) == 0 {
				log.Fatalln("no kpc files")
			}

			// found kpc
			kpc_file = files[0]
			fmt.Println(kpc_file)
		},
	}
)

func checkExt(path, ext string) []string {
	var files []string
	filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, f.Name())
			}
		}
		return nil
	})
	return files
}
