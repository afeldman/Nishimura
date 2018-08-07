package cmd

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type NishimuraConfig struct {
	RootDir []string `yaml:"nishimura_root"`
	Version string   `yaml:"nishimura_version"`
}

var rfg NishimuraConfig

func (r *NishimuraConfig) init(nishimurapath []string) {
	r.RootDir = nishimurapath
	r.Version = NISHIMURA_VERSION
}

func (r *NishimuraConfig) save(path string) {
	d, err := yaml.Marshal(r)
	if err != nil {
		log.Fatal("cannot yamalize Nishimura config")
	}

	if err := ioutil.WriteFile(path, d, 0640); err != nil {
		log.Fatal("cannot write configuration into Nishimura configuration file")
	}
}
