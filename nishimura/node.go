package nishimura

import "github.com/afeldman/Makoto/kpc"

type Node struct {
	Name    string
	Version string
	Repo    kpc.Repository
	Deps    []*Node
}
