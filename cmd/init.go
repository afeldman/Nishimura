package cmd

import (
	"fmt"
	"bufio"
	"strings"
	"path/filepath"
	"log"

	"github.com/spf13/cobra"

	kpc "github.com/afeldman/Makoto/kpc/"
)

var init = &cobra.Command{
	Use:   "init",
	Short: "print initilize a Karel project",
	Long:  `
Initilize the karel project.
AUTHOR:
	Anton Feldmann <anton.feldmann@gmail.com>
`,
	Run:   func(cmd *cobra.Command, args []string){
		InitProject();
	},
}

func getcurrentpathname() string{
	dir, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func consol_input(consoletext string) (string){
	reader:= bufio.newReader(os.Stdin)
	fmt.Print(consoletext)
	text, err := reader.ReadString('\n')
	if ( err != nil) {
		log.Fatal(err)
	}

	return text;
}

type NishimuraProject struct{
	name        string
	version     string
	description string
	mainfile    string
	parser_ver  string
	repo_type   string
	repo_add    string
	keywords    string
	authors     string
	license     string
}

func ask(project_data *NishimuraProject) bool{
	project_data.name        := consol_input("package name (" + getcurrentpathname() + "):")
	project_data.version     := consol_input("package version (0.1.0):")
	project_data.description := console_input("package description:")
	project_data.mainfile    := console_input("main file (" + name + ".kl):")
	project_data.parser_ver  := console_input("paser version (v9.10):")
	project_data.repo_type   := console_input("repository type (git):")
	project_data.repo_add    := console_input("repository address:")
	project_data.keywords    := console_input("package keywords:")
	project_data.authors     := console_input("author:")
	project_data.license     := console_input("license (MIT):")


	kpc := project_data.to_KPC(project_data)

	fmt.Println(string(kpc.ToYAML()))

	ok  := console_input("OK (yes)?:")
	if (strings.HasPrefix(strings.ToLower(strings.TrimSpace(ok)),'y')){
		return true
	}else{
		return false
	}
}

func InitProject() {

	project_data := NishimuraProject{
		getcurrentpathname(),
		"0.1.0",
		"",
		getcurrentpathname(),
		"v9.10",
		"git",
		"",
		"",
		os.GetEnv('USER'),
		"MIT"
	}

	for{
		if ask(project_data) {
			// make file for kpc and compiler
			break
		}else{
			fmt.Println("you decided to do not start a new project")
		}
	}

}

func ( this *NishimuraProject ) to_KPC() (*kpc.KPC){

	kpc := kpc.KPC_Init(this.name);

	kpc.GetVersion(this.version)
	kpc.SetDescription(this.description)
	kpc.main        = this.mainfile

	repo := kpc.Repo_Init();
	repo.SetType(this.repo_type)
	repo.SetUTL(this.repo_add)
	kpc.AddRepo(repo)
	kpc.Main = mainfile

	kpc.Authors := strings.Split(this.authors, ";")

	parser_ver  string
	keywords    string
	license     string

	return &kpc;
}
