package yidoc

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

type arrayAble struct {
	itemType string
}

type queryArgument struct {
	emptyArg
	arrayAble
	name     string
	typ      string
	required bool
}

func (a *queryArgument) Name() string {
	return a.name
}

func (a *queryArgument) Type() string {
	if a.typ == "" {
		a.typ = "string"
	}
	return a.typ
}

func (a *queryArgument) Required() bool {
	return a.required
}

func (a *queryArgument) In() string {
	return "query"
}

func QueryString(name string, ) Arg {
	return &queryArgument{name: name, typ: "string", required: getBool(required...)}
}

func QueryNumber(name string, required ...bool) Arg {
	return &queryArgument{name: name, typ: "number", required: getBool(required...)}
}

// query array string
func QueryStringArray(name string, required ...bool) Arg {
	return &queryArgument{name: name, typ: "array", arrayAble: arrayAble{itemType: "string"}, required: getBool(required...)}
}

// query array string
func QueryNumberArray(name string, required ...bool) Arg {
	return &queryArgument{name: name, typ: "array", arrayAble: arrayAble{itemType: "number"}, required: getBool(required...)}
}

func getBool(args ...bool) bool {
	if len(args) > 0 {
		return args[0]
	}
	return false
}
