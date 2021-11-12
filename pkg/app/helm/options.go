package helm

type ListReleasesOptionFunc func(*Option) error

func OnNamespace(namespace string) ListReleasesOptionFunc {
	return func(o *Option) error {
		o.Namespace = namespace
		return nil
	}
}

func OnAnyNamespace() ListReleasesOptionFunc {
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

func (o *Option) apply(options ...ListReleasesOptionFunc) error {
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
