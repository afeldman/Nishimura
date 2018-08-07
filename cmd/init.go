package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	kpc "github.com/afeldman/Makoto/kpc"
	license "github.com/afeldman/Nishimura/licenses"
	nishi "github.com/afeldman/Nishimura/nishimura"
	"github.com/afeldman/go-util/string"
)

var (
	build = &cobra.Command{
		Use:   "build",
		Short: "print initilize a Karel project",
		Long: `
Initilize a karel project.
AUTHOR:
	Anton Feldmann <anton.feldmann@gmail.com>
`,
		Run: func(cmd *cobra.Command, args []string) {
			InitProject()
		},
	}
)

func getcurrentpathname() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func console_input(consoletext, def_str string) string {
	in := bufio.NewReader(os.Stdin)
	fmt.Print(consoletext)
	text, err := in.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	text = strings.TrimSuffix(text, "\n")
	if str_util.StringEmpty(text) {
		text = def_str
	}
	return text
}

func InitProject() {

	var data = &nishi.Nishimura{}

	kpc_ := kpc.KPC_Init(data.Project_Name)

	for {
		project_name := filepath.Base(getcurrentpathname())
		data.Project_Name = console_input("package name ("+project_name+"):", project_name)
		data.Version = console_input("package version (0.1.0):", "0.1.0")
		data.Description = console_input("package description:", "")
		data.Mainfile = console_input("main file ("+data.Project_Name+".kl):", data.Project_Name+".kl")
		data.Parser_ver = console_input("paser version (v9.10):", "v9.10")
		data.Repo_type = console_input("repository type (git):", "git")
		data.Repo_add = console_input("repository address:", "")
		data.Keywords = console_input("package keywords:", "")
		data.Author = console_input("author:", os.Getenv("USER"))
		data.Email = console_input("author's email:", "")
		data.License = console_input("license (MIT):", "MIT")

		kpc_ = data.To_KPC()
		fmt.Println("\n", string(kpc_.ToYAML()))

		ok := console_input("OK (yes)?:", "no")
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(ok)), "y") {

			dir_path := createProjectDirectory(data.Project_Name)
			makeLicense(kpc_, data, dir_path)

			break
		} else {
			fmt.Println("")
			fmt.Println("you decided to do not start a new project")
			fmt.Println("")
		}
	}
}

func makeLicense(kpc_ *kpc.KPC, data *nishi.Nishimura, path string) {
	log.Println(path)

	authors := kpc_.Authors
	lic := license.GetLicense(data.License,
		*((authors[0]).GetEmail()),
		*((authors[0]).GetName()),
		data.Project_Name)

	file, err := os.Create(filepath.Join(path, "LICENSE"))
	if err != nil {
		return
	}
	defer file.Close()

	if !str_util.StringEmpty(lic) {
		file.WriteString(lic)
	} else {
		file.WriteString("")
	}
}

func createProjectDirectory(project_name string) string {
	path := getcurrentpathname()
	base_path := ""

	if filepath.Base(path) == project_name {
		isempty, err := isEmpty(path)
		if err != nil {
			log.Fatal(err)
		}
		if !isempty {
			base_path = project_name
		}
	} else {
		base_path = project_name
	}

	directoryPath := filepath.Join(path, base_path)

	log.Println(directoryPath)

	//choose your permissions well
	pathErr := os.MkdirAll(directoryPath, 0764)

	//check if you need to panic, fallback or report
	if pathErr != nil {
		fmt.Println(pathErr)
	}

	return directoryPath
}

func isEmpty(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
