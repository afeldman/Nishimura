package vcs_plugin

import (
	"fmt"
	"sync"
)

type manager struct{
	sync.Mutex
	plugins []VCS_Plugin
	registered map[string]bool
}

var (
	default_Manager = newManager()
)

func newManager() *manager{
	return &manager{
		registered make(map[string]bool)
	}
}

func (this *manager) Plugins() []VCS_Plugin{
	this.Lock()
	defer this.Unlock()

	return this.plugins
}

func (this *manager) Register(plugin VCS_Plugin) error{
	this.Lock()
	defer this.Unlock()

	vcs_type := plugin.Type()

	if this.registered[vcs_type] {
		return fmt.Errorf("This plugin of type %s with name %s is already registered", vcs_type, plugin.Name())
	}

	this.registered[vcs_type] = true
	this.plugins = append(this.plugins, plugin)
	return nil
}
