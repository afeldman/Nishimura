package karel

import (
	"os"

	"github.com/mholt/archiver"

	kpc "github.com/afeldman/kpc"
)

func BuildKarel(info *kpc.KPC, folder string) {
	name := *(info.GetName())
	version := *(info.GetVersion())

	err := archiver.Archive([]string{folder}, name+".tar.sz")
	if err != nil {
		panic(err)
	}
	err = os.Rename(name+".tar.sz", name+"@"+version+".karel")
	if err != nil {
		panic(err)
	}
}

func OpenKarel(name, folder string) {
	err := os.Rename(name+".karel", name+".tar.sz")
	if err != nil {
		panic(err)
	}

	err = archiver.Unarchive(name+".tar.sz", folder)
	if err != nil {
		panic(err)
	}
}
