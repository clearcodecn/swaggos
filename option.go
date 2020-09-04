package swaggos

type Option func(o *Swaggos)

// DefaultOptions is the default option with:
// request Content-Type - application/json
// response Content-Type - application/json
// default schemas http/https
func DefaultOptions() []Option {
	return []Option{
		WithJSON(),
		WithHttp(),
		WithHttps(),
	}
}

func WithJSON() Option {
	return func(o *Swaggos) {
		o.Produces(applicationJson)
		o.Consumes(applicationJson)
	}
}

func WithHttp() Option {
	return func(o *Swaggos) {
		o.schemas = append(o.schemas, "http")
	}
}

func WithHttps() Option {
	return func(o *Swaggos) {
		o.schemas = append(o.schemas, "https")
	}
}