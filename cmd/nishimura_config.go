package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"gopkg.in/yaml.v2"

	homedir "github.com/atrox/homedir"
	fs "github.com/afeldman/go-util/fs"
	"github.com/afeldman/go-util/string"
)

type NishimuraConfig struct {
	RootDir string `yaml:"path"`
	ConfFile string  `yaml:"file"`
	Version string   `yaml:"version"`
	TemplateDir string `yaml:"template"`
	PluginDir string `yaml:"plugin"`
}

var ncft NishimuraConfig

func DefaultConfPath() string {
	if len(os.Getenv("NISHIMURA_HOME")) > 0 {
		return os.Getenv("NISHIMURA_HOME")
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		return path.Join(home, ".config", "nishimura", "nishimura.yaml")
	}
}

func (r *NishimuraConfig) initNishimura(nishimurapath string) {
	dir, file := filepath.Split(nishimurapath)
	r.RootDir = dir
	if str_util.StringEmpty(file) {
		r.ConfFile = "nishimura.yaml"
	}else{
		r.ConfFile = file
	}
	r.Version = NISHIMURA_VERSION

	r.TemplateDir = filepath.Join(r.RootDir,"template")
	r.PluginDir = filepath.Join(r.RootDir,"plugin")

}

func (r *NishimuraConfig) save() error{
	d, err := yaml.Marshal(r)
	if err != nil {
		log.Println("cannot yamalize Nishimura config")
		return err
	}

	if err := ioutil.WriteFile(path.Join(r.RootDir,r.ConfFile), d, 0640); err != nil {
		log.Println("can not write configuration into Nishimura configuration file")
		return err
	}
	return nil
}

func (r *NishimuraConfig)build_file() bool{
	if _, err := os.Stat(path.Join(r.RootDir,r.ConfFile)); os.IsNotExist(err){
		log.Println("the requested file is not reachable")

		//build config file
		if err = r.buildconfig(); err != nil{
			log.Println(err)
			return false
		}
	}

	return true
}

func (r *NishimuraConfig)buildconfig() error{
	if err := fs.MkDir(r.RootDir,0764); err != nil {
		return err
	}

	ncft.save()
	return nil
}
