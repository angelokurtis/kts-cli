package m3u

import (
	"bufio"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jamesnetherton/m3u"
	"github.com/pkg/errors"
)

func ListChannels(filedir string) (Channels, error) {
	playlist, err := m3u.Parse(filedir)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return playlist.Tracks, nil
}

type Channels []m3u.Track

func (c Channels) Write(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	if err = m3u.MarshallInto(m3u.Playlist{Tracks: c}, bufio.NewWriter(file)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (c Channels) Get(name string) *m3u.Track {
	for _, channel := range c {
		if channel.Name == name {
			return &m3u.Track{
				Name:   channel.Name,
				Length: channel.Length,
				URI:    channel.URI,
				Tags:   channel.Tags,
			}
		}
	}
	return nil
}

func (c Channels) Names() []string {
	n := make([]string, 0, 0)
	for _, channel := range c {
		n = append(n, channel.Name)
	}
	return n
}

func (c Channels) SelectMany() (Channels, error) {
	if len(c) == 0 {
		return c, nil
	}
	prompt := &survey.MultiSelect{
		Message: "Select Channels:",
		Options: c.Names(),
	}

	var selects []string
	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	channels := make(Channels, 0, len(selects))
	for _, name := range selects {
		channel := c.Get(name)
		if channel != nil {
			channels = append(channels, *channel)
		}
	}
	return channels, nil
}
