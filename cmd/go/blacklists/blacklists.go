package blacklists

import "encoding/json"

type Blacklists []Blacklist

func UnmarshalBlacklists(data []byte) (Blacklists, error) {
	var r Blacklists
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Blacklists) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Blacklist struct {
	Import string   `json:"Import"`
	Except []string `json:"Except"`
}
