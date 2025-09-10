package nishimura

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/afeldman/Makoto/kpc"
)

// ResolveDependencies lädt rekursiv alle Dependencies eines Projekts
func ResolveDependencies(manifestPath string) (*Node, error) {
	var manifest kpc.KPC
	if _, err := toml.DecodeFile(manifestPath, &manifest); err != nil {
		return nil, fmt.Errorf("cannot read manifest %s: %w", manifestPath, err)
	}

	root := &Node{
		Name:    manifest.Name,
		Version: manifest.Version,
		Repo:    manifest.Source, // WICHTIG
	}
	visited := make(map[string]*Node)

	var loadDeps func(*Node, []kpc.Requirement) error
	loadDeps = func(parent *Node, deps []kpc.Requirement) error {
		for _, dep := range deps {
			key := fmt.Sprintf("%s@%s", dep.Name, dep.Version)
			if existing, ok := visited[key]; ok {
				parent.Deps = append(parent.Deps, existing)
				continue
			}

			node := &Node{
				Name:    dep.Name,
				Version: dep.Version,
				Repo:    dep.Source, // WICHTIG
			}
			visited[key] = node
			parent.Deps = append(parent.Deps, node)

			// rekursiv: falls das Dependency im Cache liegt
			home, _ := os.UserHomeDir()
			cachePath := filepath.Join(
				home, ".nishimura", "src",
				fmt.Sprintf("%s-%s", dep.Name, dep.Version),
				"nishimura.kpc",
			)
			if _, err := os.Stat(cachePath); err == nil {
				var depManifest kpc.KPC
				if _, err := toml.DecodeFile(cachePath, &depManifest); err == nil {
					loadDeps(node, depManifest.Requirements)
				}
			}
		}
		return nil
	}

	if err := loadDeps(root, manifest.Requirements); err != nil {
		return nil, err
	}

	return root, nil
}

// Flatten liefert eine Liste aller Nodes
func Flatten(root *Node) []string {
	var result []string
	seen := make(map[string]bool)

	var walk func(*Node)
	walk = func(n *Node) {
		key := fmt.Sprintf("%s@%s", n.Name, n.Version)
		if seen[key] {
			return
		}
		seen[key] = true
		result = append(result, key)
		for _, dep := range n.Deps {
			walk(dep)
		}
	}

	walk(root)
	return result
}
