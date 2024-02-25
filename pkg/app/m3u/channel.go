package m3u

import (
	"github.com/jamesnetherton/m3u"
)

type Channel struct {
	*m3u.Track
	filename string
}

func (c *Channel) Group() string {
	for _, tag := range c.Tags {
		if tag.Name == "group-title" {
			return tag.Value
		}
	}

	return ""
}
