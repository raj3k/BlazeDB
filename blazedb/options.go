package blazedb

type OptFunc func(opts *Options)

type Options struct {
	DBName string
}

func WithDBName(name string) OptFunc {
	return func(opts *Options) {
		opts.DBName = name
	}
}
