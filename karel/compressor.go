package karel

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"

	kpc "github.com/afeldman/kpc"
)

func BuildKarel(info *kpc.KPC, folder string) string {
	name := *(info.GetName())
	version := *(info.GetVersion())

	err := archiver.Archive([]string{folder}, name+".tar.sz")
	if err != nil {
		panic(err)
	}

	package_name := name + "@" + version + ".karel"

	err = os.Rename(name+".tar.sz", package_name)
	if err != nil {
		panic(err)
	}

	return package_name

}

func OpenKarel(name, folder string) {
	name_without_ext := strings.TrimSuffix(name, filepath.Ext(name))
	err := os.Rename(name_without_ext+".karel", name_without_ext+".tar.sz")
	if err != nil {
		panic(err)
	}

	err = archiver.Unarchive(name_without_ext+".tar.sz", folder)
	if err != nil {
		panic(err)
	}
}
