package Nishimura

import (
	"strings"

	"github.com/afeldman/kpc"
)

type Nishimura struct {
	Project_Name string
	Version      string
	Description  string
	Mainfile     string
	Parser_ver   string
	Repo_type    string
	Repo_add     string
	Keywords     string
	Author       string
	Email        string
	License      string
}

func (nishimura *Nishimura) To_KPC() *kpc.KPC {

	kpc_ := kpc.InitKPC(nishimura.Project_Name)

	kpc_.SetVersion(nishimura.Version)
	kpc_.SetDescription(nishimura.Description)
	kpc_.SetMainSourceFile(nishimura.Mainfile)

	kpc_.AddAuthor(kpc.Author{
		Name:  nishimura.Author,
		Email: nishimura.Email,
	})
	kpc_.Keywords = strings.Split(nishimura.Keywords, ",")

	return kpc_
}
