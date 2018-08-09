package vcs

type Plugin interface {
	Version()     string
	// get a list of all files
	Files()       []string
	// get the root directory path
	Root()        string
	// init the plugin with all data
	Init(string, string) err
	// the vcs type
	Type()        string
}

func NewPlugin(typ string) Plugin{

}
