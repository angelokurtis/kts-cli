package tags

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/andanhm/go-prettytime"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/angelokurtis/kts-cli/internal/log"
)

var (
	brazil  *time.Location
	printer *message.Printer
)

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err)
	}
	brazil = loc
	printer = message.NewPrinter(language.BrazilianPortuguese)
}

func list(cmd *cobra.Command, args []string) {
	dockerhub := newDockerhubClient()
	repo := args[0]
	tags, total, err := dockerhub.ListTags(repo)
	if err != nil {
		return
	}
	arch := runtime.GOARCH
	_ = arch
	imgMap := make(map[string]*Image, 0)
	for _, tag := range tags {
		for _, image := range tag.Images {
			img := imgMap[image.Digest]
			if img == nil {
				img = new(Image)
			}
			img.Pushed = image.LastPushed
			img.Size = image.Size
			img.Architecture = image.Architecture
			img.Digest = image.Digest
			img.Add(&Tag{Name: tag.Name, Updated: tag.LastUpdated})
			imgMap[image.Digest] = img
		}
	}

	images := make([]*Image, 0, len(imgMap))
	for _, image := range imgMap {
		if image.Architecture == runtime.GOARCH {
			images = append(images, image)
		}
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

	table.SetHeader([]string{"IMAGE", "TAG", "DIGEST", "SIZE", "UPDATED"})
	for _, img := range images {
		table.Append([]string{repo, strings.Join(img.TagNames(), ", "), img.Digest, ByteCount(img.Size), prettytime.Format(img.Pushed)})
	}

	table.Render()
	link := fmt.Sprintf("https://hub.docker.com/%s?tab=tags", func() string {
		if strings.Contains(repo, "/") {
			return "r/" + repo
		} else {
			return "_/" + repo
		}
	}())
	fmt.Printf("\nfound %d on %s\n", total, link)
}

func ByteCount(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

type Image struct {
	Tags         []*Tag
	Pushed       time.Time
	Size         int64
	Architecture string
	Digest       string
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
