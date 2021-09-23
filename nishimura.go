package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
)

func main() {
	/*if err := cmd.Nishimura.Execute(); err != nil {
		log.Fatal(err)
	}*/

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	p := filepath.FromSlash(dir)

	name := "nishimura"
	version := "0.1.0"

	err = archiver.Archive([]string{filepath.Base(p)}, name+".sz")
	if err != nil {
		panic(err)
	}
	err = os.Rename(name+".sz", name+"@"+version+".karel")
	if err != nil {
		panic(err)
	}

}
