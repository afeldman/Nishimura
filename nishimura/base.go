package Nishimura

import (
	"github.com/afeldman/Makoto/kpc"
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

func (this *Nishimura) To_KPC() *kpc.KPC {

	kpc_ := kpc.KPC_Init(this.Project_Name)

	kpc_.SetVersion(this.Version)
	kpc_.SetDescription(this.Description)
	kpc_.SetMainSourceFile(this.Mainfile)

	repo := kpc.Repo_Init()
	repo.SetType(this.Repo_type)
	repo.SetURL(this.Repo_add)
	kpc_.AddRepo(*repo)

	kpc_.AddAuthor(kpc.Author{
		Name:  this.Author,
		Email: this.Email,
	})

	return kpc_
}
