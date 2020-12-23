package swaggos

// Option option is swaggos option
type Option func(o *Swaggos)

// DefaultOptions is the default option with:
// request Content-Type - application/json
// response Content-Type - application/json
// default schemas http/https
func DefaultOptions() []Option {
	return []Option{
		WithJSON(),
		WithHTTP(),
		WithHTTPS(),
	}
}

// WithJSON provide a json response
func WithJSON() Option {
	return func(o *Swaggos) {
		o.Produces(applicationJSON)
		o.Consumes(applicationJSON)
	}
}

// WithHTTP provide a HTTP protocol
func WithHTTP() Option {
	return func(o *Swaggos) {
		o.schemas = append(o.schemas, "http")
	}
}

// WithHTTPS provide a HTTPS protocol
func WithHTTPS() Option {
	return func(o *Swaggos) {
		o.schemas = append(o.schemas, "https")
	}
}
