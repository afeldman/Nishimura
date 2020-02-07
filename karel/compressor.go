package karel

import (
	"log"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"

	kpc "github.com/afeldman/kpc"
)

func BuildKarel(info kpc.KPC, folder string) {

	err := os.Chdir(folder)
	if err != nil {
		panic(err)
	}

	name := *(info.GetName())
	version := *(info.GetVersion())

	var files []string

	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			//fmt.Println(path, info.Size())
			files = append(files, path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		log.Println(file)
	}

	err = archiver.TarXZ.Make(name+"-"+version+".karel", files)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

func OpenKarel(name, folder string) {
	err := archiver.TarXZ.Open(name+".karel", folder)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
