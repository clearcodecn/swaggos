package swaggos

type Option func(o *Swaggo)

func DefaultOptions() []Option {
	return []Option{
		WithJSON(),
		WithHttp(),
		WithHttps(),
	}
}

func WithJSON() Option {
	return func(o *Swaggo) {
		o.Produces(applicationJson)
		o.Consumes(applicationJson)
	}
}

func WithHttp() Option {
	return func(o *Swaggo) {
		o.schemas = append(o.schemas, "http")
	}
}

func WithHttps() Option {
	return func(o *Swaggo) {
		o.schemas = append(o.schemas, "https")
	}
}