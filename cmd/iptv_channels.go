package cmd

import (
	"bufio"
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/iptv"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// iptvChannelsCmd represents the iptv channels command
var iptvChannelsCmd = &cobra.Command{
	Use:   "channels",
	Short: "List all available IPTV channels",
	RunE: func(cmd *cobra.Command, args []string) error {
		// creating temp file
		file, err := ioutil.TempFile("", "channels.*.m3u")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		// getting credentials from env var
		user, ok := os.LookupEnv("IPTV_USERNAME")
		if !ok {
			return errors.New("the environment 'IPTV_USERNAME' should be set")
		}
		pass, ok := os.LookupEnv("IPTV_PASSWORD")
		if !ok {
			return errors.New("the environment 'IPTV_PASSWORD' should be set")
		}

		// downloading file
		res, err := http.Get(fmt.Sprintf("http://srvx.io/get.php?username=%s&password=%s&type=m3u_plus&output=ts", user, pass))
		if err != nil {
			return err
		}

		// reading response's body
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		// writing to the temp file
		if _, err = file.Write(body); err != nil {
			return err
		}

		// scanning channels
		scanner := bufio.NewScanner(file)
		extinf := ""
		channels := make(iptv.Channels, 0, 0)
		for scanner.Scan() {
			txt := scanner.Text()
			if strings.HasPrefix(txt, "#EXTINF:-1") {
				extinf = txt
			} else if extinf != "" {
				channels = append(channels, iptv.NewChannel(extinf, txt))
				extinf = ""
			}
		}
		if err = scanner.Err(); err != nil {
			return err
		}

		// selecting channels
		channels, err = channels.SelectMany()
		if err != nil {
			return err
		}

		// writing selected channels
		err = channels.Write("/home/tiagoangelo/Downloads/channels.m3u")
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	iptvCmd.AddCommand(iptvChannelsCmd)
}
