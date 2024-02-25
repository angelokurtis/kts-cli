package helm

type OptionFunc func(*Option) error

func OnNamespace(namespace string) OptionFunc {
	return func(o *Option) error {
		o.Namespace = namespace
		return nil
	}
}

func OnAnyNamespace() OptionFunc {
	return func(o *Option) error {
		o.Namespace = ""
		o.AllNamespaces = true

		return nil
	}
}

type Option struct {
	Namespace     string
	AllNamespaces bool
}

func (o *Option) apply(options ...OptionFunc) error {
	for _, fn := range options {
		if fn == nil {
			continue
		}

		if err := fn(o); err != nil {
			return err
		}
	}

	return nil
}
