package main

import (
	"log"

	"github.com/afeldman/Nishimura/src/cmd"
)

func main() {
	if err := cmd.Nishimura.Execute(); err != nil {
		log.Fatal(err)
	}
}
