package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/afeldman/Gakutensoku/upload"
	"github.com/afeldman/Nishimura/karel"
	"github.com/afeldman/go-util/env"
	kpc "github.com/afeldman/kpc"
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
			gakutensoku_url := env.GetEnvOrDefault("GAKUTENSOKU_URL", "http://localhost:2510")

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
				log.Fatalln("use etleased one kpc file in your path")
			}

			wg := &sync.WaitGroup{}
			wg.Add(len(files))

			// found kpc
			for _, kpc_file := range files {
				go func(wg *sync.WaitGroup, kpc_file string) {
					defer wg.Done()

					project_path := filepath.Dir(kpc_file)
					fmt.Println(project_path)
					//read kpc
					kpc_data := kpc.ReadKPCFile(kpc_file)

					filepath := karel.BuildKarel(kpc_data, project_path)

					// send file to gakutensoku
					client := upload.NewClient(gakutensoku_url, filepath)
					upload.SendData(client)

					// delete karelfile
				}(wg, kpc_file)

			}

			wg.Wait()

		},
	}
)

func checkExt(path, ext string) []string {
	var files []string
	filepath.Walk(path, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
				files = append(files, path)
			}
		}
		return nil
	})
	return files
}
