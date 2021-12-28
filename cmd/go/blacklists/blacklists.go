package blacklists

import "encoding/json"

type Blacklists struct {
	Set []*Blacklist
}

func UnmarshalBlacklists(data []byte) (*Blacklists, error) {
	var r []*Blacklist
	err := json.Unmarshal(data, &r)
	return &Blacklists{Set: r}, err
}

func (r *Blacklists) Marshal() ([]byte, error) {
	return json.MarshalIndent(r.Set, "", "  ")
}

func (r *Blacklists) Allow(imports []string, on string) {
	for _, imp := range imports {
		blacklist := r.getOrCreate(imp)
		blacklist.AddExcept(on)
	}
}

func (r *Blacklists) getOrCreate(imp string) *Blacklist {
	for _, blacklist := range r.Set {
		if blacklist.Import == imp {
			return blacklist
		}
	}
	newone := &Blacklist{Import: imp}
	r.Set = append(r.Set, newone)
	return newone
}

type Blacklist struct {
	Import string   `json:"Import"`
	Except []string `json:"Except"`
}

func (r *Blacklist) AddExcept(dir string) {
	r.Except = dedupe(r.Except, dir)
}
