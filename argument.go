package ydoc

type Arg interface {
	Name() string
	Type() string
	Description() string
	Required() bool
	In() string
}

type emptyArg struct{}

func (e *emptyArg) Name() string { return "" }

func (e *emptyArg) Type() string { return "" }

func (e *emptyArg) Description() string { return "" }

func (e *emptyArg) Required() bool { return false }

func (e *emptyArg) In() string { return "" }
