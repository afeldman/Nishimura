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
	RootDir string `yaml:"nishimura_root"`
	Version string   `yaml:"nishimura_version"`
}

var ncft NishimuraConfig

func (r *NishimuraConfig) initNishimura(nishimurapath string) {
	r.RootDir = nishimurapath
	r.Version = NISHIMURA_VERSION
}

func (r *NishimuraConfig) save() error{
	d, err := yaml.Marshal(r)
	if err != nil {
		log.Fatal("cannot yamalize Nishimura config")
		return err
	}

	if err := ioutil.WriteFile(r.RootDir, d, 0640); err != nil {
		log.Fatal("can not write configuration into Nishimura configuration file")
		return err
	}
	return nil
}

func loadConfig(path string) (*NishimuraConfig,error){

	var config NishimuraConfig

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("can not read the file")
		return &config, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return &config, err
	}

	return &config, nil
}

func (r *NishimuraConfig)build_file() bool{
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
		return false
	}

	if !str_util.StringEmpty(r.RootDir){
		if _, err := os.Stat(r.RootDir); os.IsNotExist(err){
			log.Println("the requested file is not reachable")
			//build config file
			if err = buildconfig(r.RootDir); err != nil{
				log.Println(err)
				return false
			}
		}
	}else if len(os.Getenv("NISHIMURA_HOME")) > 0 {
		r.RootDir = os.Getenv("NISHIMURA_HOME")
		if err = buildconfig(r.RootDir); err != nil {
			log.Println(err)
			return false
		}
	}else{
		r.RootDir = path.Join(home, ".config", "nishimura", "nishimura.yaml")
		if err = buildconfig(r.RootDir); err != nil {
			log.Println(err)
			return false
		}
	}

	return true
}

func buildconfig(path string) error{
	if err := fs.MkDir(filepath.Dir(path),0764); err != nil {
		return err
	}

	ncft.initNishimura(path)
	ncft.save()
	return nil
}
