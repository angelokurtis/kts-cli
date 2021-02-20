package kiali

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

func LoadGraphInfo(ns ...string) (*Graph, error) {
	url := "http://127.0.0.1:20001/kiali/api/namespaces/graph?appenders=deadNode,sidecarsCheck,serviceEntry,istio&namespaces=" + strings.Join(ns, ",")
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return nil, errors.WithStack(err)
	//}
	//fmt.Println(string(body))
	target := &Graph{}
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, errors.WithStack(err)
	}

	return target, nil
}

type Graph struct {
	Timestamp int    `json:"timestamp"`
	Duration  int    `json:"duration"`
	GraphType string `json:"graphType"`
	Elements  struct {
		Nodes []struct {
			Data *Node `json:"data"`
		} `json:"nodes"`
		Edges []struct {
			Data *Edge `json:"data"`
		} `json:"edges"`
	} `json:"elements"`
}

func (g *Graph) GetNodes() Nodes {
	nodes := make(Nodes, 0)
	for _, node := range g.Elements.Nodes {
		nodes[node.Data.ID] = node.Data
	}
	return nodes
}

func (g *Graph) GetEdges() Edges {
	edges := make(Edges, 0)
	for _, edge := range g.Elements.Edges {
		edges = append(edges, edge.Data)
	}
	return edges
}

func (g *Graph) Inbound(n *Node) Nodes {
	nodes := g.GetNodes()
	r := make(Nodes, 0)
	for _, id := range g.GetEdges().Inbound(n) {
		r[id] = nodes[id]
	}
	return r
}

func (g *Graph) Outbound(n *Node) Nodes {
	nodes := g.GetNodes()
	r := make(Nodes, 0)
	for _, id := range g.GetEdges().Outbound(n) {
		r[id] = nodes[id]
	}
	return r
}
