package iptv

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	changecase "github.com/ku/go-change-case"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var Command = &cobra.Command{
	Use:   "iptv",
	Short: "IPTV functions utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "channels", Run: channels})
}

// iptv channels
func channels(cmd *cobra.Command, args []string) {
	file, err := os.Open("/home/tiagoangelo/Downloads/tv_channels_jrGF1_plus.m3u")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	extinf := ""
	channels := make(Channels, 0, 0)
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "#EXTINF:-1") {
			extinf = txt
		} else if extinf != "" {
			channels = append(channels, NewChannel(extinf, txt))
			extinf = ""
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	channels, err = channels.SelectMany()
	if err != nil {
		log.Fatal(err)
	}
	err = channels.Write("/home/tiagoangelo/Downloads/channels.m3u")
	if err != nil {
		log.Fatal(err)
	}
}

type Channels []*Channel

func (c Channels) Write(path string) error {
	var b strings.Builder
	fmt.Fprint(&b, "#EXTM3U\n")
	for _, channel := range c {
		fmt.Fprintf(&b, "#EXTINF:-1 tvg-logo=\"%s\" group-title=\"%s\",%s\n%s\n", channel.Logo, channel.Group, channel.Name, channel.Address)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(b.String())
	return err
}

func (c Channels) Get(id string) *Channel {
	for _, channel := range c {
		if channel.ID == id {
			return channel
		}
	}
	return nil
}

func (c Channels) IDs() []string {
	n := make([]string, 0, 0)
	for _, channel := range c {
		n = append(n, channel.ID)
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
	}

	var selects []string
	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	channels := make(Channels, 0, len(selects))
	for _, name := range selects {
		channels = append(channels, c.Get(name))
	}
	return channels, nil
}

type Channel struct {
	Logo    string
	Address string
	Group   string
	Name    string
	ID      string
}

func NewChannel(extinf string, addr string) *Channel {
	logo := strings.Split(extinf, "tvg-logo=\"")[1]
	logo = strings.Split(logo, "\"")[0]
	group := strings.Split(extinf, "group-title=\"")[1]
	group = strings.Split(group, "\"")[0]
	name := strings.Split(extinf, ",")[1]
	return &Channel{
		Logo:    logo,
		Address: addr,
		Group:   group,
		Name:    name,
		ID:      changecase.Snake(name),
	}
}
