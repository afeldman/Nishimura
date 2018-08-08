package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"io/ioutil"

	kpc "github.com/afeldman/Makoto/kpc"
	license "github.com/afeldman/Nishimura/licenses"
	nishi "github.com/afeldman/Nishimura/nishimura"
	gakuten "github.com/afeldman/Gakutensoku/ktrans"

	"github.com/afeldman/go-util/string"
	"github.com/afeldman/go-util/fs"
	"github.com/afeldman/go-util/file"

	"github.com/spf13/cobra"
	"github.com/vigneshuvi/GoDateFormat"
	"github.com/afeldman/go-util/time"
	"html/template"
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
			InitPackage()
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
	return strings.TrimSpace(text)
}

func InitPackage() {

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

		kpc_ = data.To_KPC();
		if kpc_err, kpc_data := kpc_.ToYAML(); kpc_err != nil {
			log.Fatal(kpc_err)
		}else {
			fmt.Println("\n", string(kpc_data))
		}

		ok := console_input("OK (yes)?:", "no")
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(ok)), "y") {

			dir_path := createProjectDirectory(data.Project_Name)
			if err := makeLicense(kpc_, data, dir_path); err != nil {
				log.Println(err)
			}

			if data.Repo_type == "git" {
				if err := copy_git(dir_path); err != nil {
					log.Println(err)
				}
			}

			if err := make_kpc(kpc_,data.Project_Name,dir_path); err != nil {
				log.Fatal(err)
			}

			if err := make_compiler_conf(data, dir_path); err != nil {
				log.Fatal(err)
			}

			if err := build_start_file(data, dir_path); err != nil {
				log.Fatal(err)
			}

			break
		} else {
			fmt.Println("")
			fmt.Println("you decided to do not start a new project")
			fmt.Println("")
		}
	}
}

func make_compiler_conf(data *nishi.Nishimura, path string) error{
	compiler_info := gakuten.KtransInit()
	compiler_info.Version = data.Parser_ver
	compiler_info.Input = data.Mainfile

	file, err_ := os.Create(filepath.Join(path, ".ktrans.conf"))
	if err_ != nil {
		return err_
	}
	defer file.Close()

	err, file_containt := compiler_info.ToJSON()
	if err != nil {
		return err
	}

	file.WriteString(string(file_containt))
	file.Sync()
	return nil
}

func make_kpc(kpc_ *kpc.KPC, target, path string) error{
	tmpfile, err := ioutil.TempFile("", "nishimura-")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	kpc_err, outputdata := kpc_.ToJSON()
	if kpc_err != nil {
		return kpc_err
	}

	if _, err := tmpfile.WriteString(string(outputdata)); err != nil {
		return err
	}
	tmpfile.Sync()

	if err := fileinfo.Fcopy(tmpfile.Name(), filepath.Join(path, target+".kpc")); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	return nil
}

func copy_git(path string) error {
	nishimura_home_template := filepath.Join(ncft.RootDir, "template")

	if err := fileinfo.Fcopy(filepath.Join(nishimura_home_template, "Karel.gitignore"),
		filepath.Join(path, ".gitignore")); err != nil {
		return err
	}
	if err := fileinfo.Fcopy(filepath.Join(nishimura_home_template, "Karel.gitattribute"),
		filepath.Join(path, ".gitattribute")); err != nil {
		return err
	}
	return nil
}

func makeLicense(kpc_ *kpc.KPC, data *nishi.Nishimura, path string) error{
	authors := kpc_.Authors
	lic, err := license.GetLicense(data.License,
		*((authors[0]).GetEmail()),
		*((authors[0]).GetName()),
		data.Project_Name)
	if err != nil {
		return err
	}

	file, err_ := os.Create(filepath.Join(path, "LICENSE"))
	if err_ != nil {
		return err_
	}
	defer file.Close()

	if !str_util.StringEmpty(lic) {
		file.WriteString(lic)
	} else {
		file.WriteString("")
	}

	return nil
}

func createProjectDirectory(project_name string) string {
	path := getcurrentpathname()
	base_path := ""

	if filepath.Base(path) == project_name {
		isempty, err := filesystem.IsEmpty(path)
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

	pathErr := filesystem.MkDir(directoryPath, 0764)
	if pathErr != nil {
		log.Println(pathErr)
	}

	return directoryPath
}

func build_start_file(data *nishi.Nishimura, path string) error{

	type page_data struct {
		FileName    string
		SmallDesc   string
		Desc        string
		Copyright   string
		Author      string
		Today       string
		License     string
		Projectname string
	}

	today := time_util.GetToday(GoDateFormat.ConvertFormat("dd-MMM-yyyy"))

	page_ := page_data{
		FileName: data.Mainfile,
		SmallDesc: data.Description,
		Desc: data.Description,
		Copyright: data.License,
		Author: data.Author,
		Today: today,
		License: data.License,
		Projectname: data.Project_Name,
	}

	nishimura_home_template := filepath.Join(ncft.RootDir, "template", "project.kl")
	tmpl, err := template.ParseFiles(nishimura_home_template)
	if err != nil {
		return err
	}

	file, err_ := os.Create(filepath.Join(path, data.Mainfile))
	if err_ != nil {
		return err_
	}
	defer file.Close()

	tmpl.Execute(file, page_)

	return nil
}
