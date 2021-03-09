package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/andanhm/go-prettytime"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

// dockerImagesCmd represents the docker command
var dockerImagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Search for images and its tagging timeline",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("please, inform the image name")
		}
		repo := args[0]
		res, err := searchForImage(repo)
		if err != nil {
			return err
		}

		imgMap := make(map[string]*Image, 0)
		for _, tag := range res.Results {
			for _, image := range tag.Images {
				img := imgMap[image.Digest]
				if img == nil {
					img = new(Image)
				}
				img.Pushed = image.LastPushed
				img.Size = int64(image.Size)
				img.Architecture = image.Architecture
				img.Add(&Tag{Name: tag.Name, Updated: tag.LastUpdated})
				imgMap[image.Digest] = img
			}
		}

		images := make([]*Image, 0, len(imgMap))
		for _, image := range imgMap {
			images = append(images, image)
		}

		sort.Slice(images, func(i, j int) bool {
			it := images[i].Pushed
			jt := images[j].Pushed
			return it.After(jt)
		})

		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetColumnSeparator("")
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.SetColWidth(100)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

		table.SetHeader([]string{"IMAGE", "ARCHITECTURE", "TAG", "UPDATED", "SIZE"})
		for _, img := range images {
			table.Append([]string{repo, img.Architecture, strings.Join(img.TagNames(), ", "), prettytime.Format(img.Pushed), ByteCountSI(img.Size)})
		}

		table.Render()
		link := fmt.Sprintf("https://hub.docker.com/%s?tab=tags", func() string {
			if strings.Contains(repo, "/") {
				return "r/" + repo
			} else {
				return "_/" + repo
			}
		}())
		fmt.Printf("\nfound %d on %s\n", res.Count, link)
		return nil
	},
}

func init() {
	dockerCmd.AddCommand(dockerImagesCmd)
}

func searchForImage(img string) (*ImageSearchResult, error) {
	if !strings.Contains(img, "/") {
		img = "library/" + img
	}
	url := fmt.Sprintf("https://hub.docker.com/v2/repositories/%s/tags/?page_size=1000&page=1&ordering=last_updated", img)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	target := new(ImageSearchResult)
	if err = json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}

	return target, nil
}

type Image struct {
	Tags         []*Tag
	Pushed       time.Time
	Size         int64
	Architecture string
}

func (i *Image) Add(t *Tag) {
	if i.Tags == nil {
		i.Tags = make([]*Tag, 0, 0)
	}
	i.Tags = append(i.Tags, t)
}

func (i *Image) TagNames() []string {
	names := make([]string, 0, len(i.Tags))
	for _, tag := range i.Tags {
		names = append(names, tag.Name)
	}
	return names
}

type Tag struct {
	Name    string
	Updated time.Time
}

type ImageSearchResult struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Creator int         `json:"creator"`
		ID      int         `json:"id"`
		ImageID interface{} `json:"image_id"`
		Images  []struct {
			Architecture string      `json:"architecture"`
			Features     string      `json:"features"`
			Variant      interface{} `json:"variant"`
			Digest       string      `json:"digest"`
			Os           string      `json:"os"`
			OsFeatures   string      `json:"os_features"`
			OsVersion    interface{} `json:"os_version"`
			Size         int         `json:"size"`
			Status       string      `json:"status"`
			LastPulled   time.Time   `json:"last_pulled"`
			LastPushed   time.Time   `json:"last_pushed"`
		} `json:"images"`
		LastUpdated         time.Time `json:"last_updated"`
		LastUpdater         int       `json:"last_updater"`
		LastUpdaterUsername string    `json:"last_updater_username"`
		Name                string    `json:"name"`
		Repository          int       `json:"repository"`
		FullSize            int       `json:"full_size"`
		V2                  bool      `json:"v2"`
		TagStatus           string    `json:"tag_status"`
		TagLastPulled       time.Time `json:"tag_last_pulled"`
		TagLastPushed       time.Time `json:"tag_last_pushed"`
	} `json:"results"`
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
