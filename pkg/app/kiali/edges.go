package kiali

type (
	Edges []*Edge
	Edge  struct {
		ID      string `json:"id"`
		Source  string `json:"source"`
		Target  string `json:"target"`
		Traffic struct {
			Protocol string `json:"protocol"`
			Rates    struct {
				HTTP           string `json:"http"`
				HTTPPercentReq string `json:"httpPercentReq"`
			} `json:"rates"`
			Responses struct {
				Num202 struct {
					Flags interface{} `json:"flags"`
					Hosts interface{} `json:"hosts"`
				} `json:"202"`
			} `json:"responses"`
		} `json:"traffic"`
	}
)

func (e Edges) Inbound(n *Node) []string {
	ids := make([]string, 0, 0)
	ids = append(ids, n.ID)
	for _, edge := range e {
		if edge.HasTargets(n.ID) {
			ids = dedupe(ids, edge.Source)
		}
	}
	return ids
}

func (e Edges) Outbound(n *Node) []string {
	ids := make([]string, 0, 0)
	ids = append(ids, n.ID)
	for _, edge := range e {
		if edge.HasSources(n.ID) {
			ids = dedupe(ids, edge.Target)
		}
	}
	return ids
}

func (e *Edge) HasTargets(ids ...string) bool {
	return contains(ids, e.Target)
}

func (e *Edge) HasSources(ids ...string) bool {
	return contains(ids, e.Source)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func dedupe(a []string, b ...string) []string {
	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}
	for letter := range check {
		res = append(res, letter)
	}
	return res
}
