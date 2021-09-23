package main

import (
	"log"

	"github.com/afeldman/Nishimura/cmd"
)

func main() {
	if err := cmd.Nishimura.Execute(); err != nil {
		log.Fatal(err)
	}
}
