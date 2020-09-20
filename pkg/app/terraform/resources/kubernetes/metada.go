package kubernetes

type Metadata struct {
	Name      string `hcl:"name"`
	Namespace string `hcl:"namespace"`
}
