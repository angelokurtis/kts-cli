package m3u

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/jamesnetherton/m3u"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type Channels []*Channel

func ListChannels(filedir string) (Channels, error) {
	playlist, err := m3u.Parse(filedir)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	channels := make(Channels, 0, 0)
	for _, track := range playlist.Tracks {
		channels = append(channels, &Channel{
			filename: filedir,
			Track: &m3u.Track{
				Name:   track.Name,
				Length: track.Length,
				URI:    track.URI,
				Tags:   track.Tags,
			},
		})
	}

	return channels, nil
}

func (c Channels) tracks() []m3u.Track {
	tracks := make([]m3u.Track, 0, 0)
	for _, channel := range c {
		tracks = append(tracks, *channel.Track)
	}

	return tracks
}

func (c Channels) Write(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if err = m3u.MarshallInto(m3u.Playlist{Tracks: c.tracks()}, bufio.NewWriter(file)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c Channels) Get(id string) *Channel {
	for _, channel := range c {
		if strings.ReplaceAll(channel.Name, " ", ".") == id {
			return &Channel{
				filename: channel.filename,
				Track: &m3u.Track{
					Name:   channel.Name,
					Length: channel.Length,
					URI:    channel.URI,
					Tags:   channel.Tags,
				},
			}
		}
	}

	return nil
}

func (c Channels) IDs() []string {
	n := make([]string, 0, 0)
	for _, channel := range c {
		n = append(n, strings.ReplaceAll(channel.Name, " ", "."))
	}

	return n
}

func (c Channels) SelectMany() (Channels, error) {
	if len(c) == 0 {
		return c, nil
	}

	prompt := &survey.MultiSelect{
		Message: "Select Channels:",
		Options: c.IDs(),
		Default: c.defaults(),
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	channels := make(Channels, 0, len(selects))

	for _, name := range selects {
		channel := c.Get(name)
		if channel != nil {
			channels = append(channels, channel)
		}
	}

	return channels, nil
}

func (c Channels) Groups() *Groups {
	groups := make([]string, 0, 0)

	for _, channel := range c {
		for _, tag := range channel.Tags {
			if tag.Name == "group-title" {
				groups = append(groups, tag.Value)
			}
		}
	}

	groups = lo.Uniq(groups)
	sort.Strings(groups)

	return &Groups{
		Items:    groups,
		filename: c.FileName(),
	}
}

func (c Channels) FileName() string {
	for _, channel := range c {
		return channel.filename
	}

	return ""
}

func (c Channels) FilterByGroups(groups *Groups) Channels {
	channels := make(Channels, 0, 0)

	for _, channel := range c {
		present := lo.ContainsBy(groups.Items, func(x string) bool {
			return channel.Group() == x
		})
		if present {
			channels = append(channels, channel)
		}
	}

	return channels
}

func (c Channels) defaults() []string {
	filename := c.FileName()
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	channels, err := ListChannels(name + "[edited]" + ext)
	if err != nil {
		channels, err = ListChannels(filename)
		if err != nil {
			return nil
		}
	}

	return channels.IDs()
}
