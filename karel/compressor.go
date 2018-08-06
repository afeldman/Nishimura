package karel

import (
	"log"

	"github.com/mholt/archiver"

	kpc "github.com/afeldman/Makoto/kpc"
)

func BuildKarel(info kpc.KPC, folder string) {
	name := *(info.GetName())
	version := *(info.GetVersion())
	err := archiver.TarXZ.Make(name+"-"+version+".karel", []string{folder})
	if err != nil {
		log.Prindln(err)
		panic(err)
	}
}

func OpenKarel(name, folder string) {
	err := archiver.TarXZ.Open(name+".karel", folder)
	if err != nil {
		log.Prindln(err)
		panic(err)
	}
}
