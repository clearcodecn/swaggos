package yidoc

type Option func(d *YiDoc)

func Produces(s []string) Option {
	return func(d *YiDoc) {
		d.produces = s
	}
}

func Consumes(s []string) Option {
	return func(d *YiDoc) {
		d.consumes = s
	}
}

func DefaultOptions() []Option {
	return []Option{
		Produces([]string{"application/json"}),
		Consumes([]string{"application/json"}),
	}
}


