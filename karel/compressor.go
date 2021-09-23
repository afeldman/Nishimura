package karel

import (
	"os"
	"strings"

	"github.com/mholt/archiver"

	kpc "github.com/afeldman/kpc"
)

func BuildKarel(info kpc.KPC, folder string) {
	name := *(info.GetName())
	version := *(info.GetVersion())

	err := archiver.Archive([]string{folder}, name+".sz")
	if err != nil {
		panic(err)
	}
	err = os.Rename(name+".sz", name+"@"+version+".karel")
	if err != nil {
		panic(err)
	}
}

func OpenKarel(name, folder string) {
	// name is "name@version"
	purename := strings.Split(name, "@")
	err := os.Rename(name+".karel", purename[0]+".sz")
	if err != nil {
		panic(err)
	}

	err = archiver.Unarchive(purename[0]+".sz", folder)
	if err != nil {
		panic(err)
	}
}
