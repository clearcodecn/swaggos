package swaggos

import (
	"strings"

	"github.com/go-openapi/spec"
)

// Group is a path group for same prefix
type Group struct {
	prefix  string
	swaggos *Swaggos
	params  []spec.Parameter
}

// NewGroup returns a new Group
func NewGroup(prefix string, swaggos *Swaggos) *Group {
	group := new(Group)
	prefix = "/" + strings.Trim(prefix, "/")
	group.prefix = prefix
	group.swaggos = swaggos
	return group
}

// Group returns a new group based on current group
func (g *Group) Group(prefix string) *Group {
	group := new(Group)
	if g.prefix == "/" {
		prefix = g.prefix + strings.Trim(prefix, "/")
	} else {
		prefix = g.prefix + "/" + strings.Trim(prefix, "/")
	}
	group.prefix = prefix
	group.swaggos = g.swaggos
	group.params = make([]spec.Parameter, 0)
	copy(group.params, g.params)
	return group
}

// Header defines a header param
func (g *Group) Header(name string, desc string, required bool) *Group {
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: _String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          _InHeader,
			Required:    required,
		},
	}
	g.params = append(g.params, param)
	return g
}

// Query add query param to the group
func (g *Group) Query(name string, desc string, required bool) *Group {
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: _String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          _InQuery,
			Required:    required,
		},
	}
	g.params = append(g.params, param)
	return g
}

// Form add form param to the group
func (g *Group) Form(name string, desc string, required bool) *Group {
	param := spec.Parameter{
		SimpleSchema: spec.SimpleSchema{
			Type: _String,
		},
		ParamProps: spec.ParamProps{
			Description: desc,
			Name:        name,
			In:          _InForm,
			Required:    required,
		},
	}
	g.params = append(g.params, param)
	return g
}

// Get create a path with group's prefix and given path of Get method
func (g *Group) Get(path string) *Path {
	p := g.swaggos.Get(g.trimPath(path))
	return p.addParam(g.params...)
}

// Post create a path with group's prefix and given path of Post method
func (g *Group) Post(path string) *Path {
	p := g.swaggos.Post(g.trimPath(path))
	return p.addParam(g.params...)
}

// Put create a path with group's prefix and given path of Put method
func (g *Group) Put(path string) *Path {
	p := g.swaggos.Put(g.trimPath(path))
	return p.addParam(g.params...)
}

// Patch create a path with group's prefix and given path of Patch method
func (g *Group) Patch(path string) *Path {
	p := g.swaggos.Patch(g.trimPath(path))
	return p.addParam(g.params...)
}

// Options create a path with group's prefix and given path of Options method
func (g *Group) Options(path string) *Path {
	p := g.swaggos.Options(g.trimPath(path))
	return p.addParam(g.params...)
}

// Delete create a path with group's prefix and given path of Delete method
func (g *Group) Delete(path string) *Path {
	p := g.swaggos.Delete(g.trimPath(path))
	return p.addParam(g.params...)
}

func (g *Group) trimPath(path string) string {
	path = "/" + strings.Trim(path, "/")
	return "/" + strings.Trim(g.prefix+path, "/")
}

// Swaggos returns instance of Swaggos
func (g *Group) Swaggos() *Swaggos {
	return g.swaggos
}

// Group setup the group of swaggos
func (swaggos *Swaggos) Group(path string) *Group {
	return NewGroup(path, swaggos)
}
