package tags

import (
	"fmt"
	log "log"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	mastermindssemver "github.com/Masterminds/semver"
	prettytime "github.com/andanhm/go-prettytime"
	"github.com/gotidy/ptr"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var brazil *time.Location

var semanticVersionRegex = regexp.MustCompile(`\d+(?:\.\d+){0,2}`)

func init() {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatal(err.Error())
	}

	brazil = loc
}

func list(cmd *cobra.Command, args []string) {
	dockerhub := newDockerhubClient()
	repo := args[0]

	tags, total, err := dockerhub.ListTags(repo)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Build map of unique images by digest
	imgMap := map[string]*Image{}

	for _, tag := range tags {
		for _, image := range tag.Images {
			img, exists := imgMap[image.Digest]
			if !exists {
				img = &Image{
					Pushed:       image.LastPushed,
					Size:         image.Size,
					Architecture: image.Architecture,
					Digest:       image.Digest,
				}
				imgMap[image.Digest] = img
			}

			img.Add(&Tag{Name: tag.Name, Updated: tag.LastUpdated})
		}
	}

	// Filter and collect images matching current architecture
	var images []*Image

	for _, img := range imgMap {
		if img.Architecture == runtime.GOARCH {
			images = append(images, img)
		}
	}

	// Sort images by last pushed date (descending)
	sort.Slice(images, func(i, j int) bool {
		return ptr.To(images[i].Pushed).After(ptr.To(images[j].Pushed))
	})

	// Setup output table
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(100)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"IMAGE", "TAG", "DIGEST", "SIZE", "UPDATED"})

	// Constraint to filter versions (e.g. semver ^1.0.0)
	var constraint *mastermindssemver.Constraints
	if len(semver) > 0 {
		constraint, err = mastermindssemver.NewConstraint(semver)
		if err != nil {
			log.Fatal(err.Error())
		}

		// Populate table with filtered and formatted image data
		for _, img := range images {
			if constraint == nil {
				continue
			}

			// At least one tag must satisfy the constraint.
			valid := false

			for _, tag := range img.TagNames() {
				version, ok := convertToSemVer(tag)
				if !ok {
					continue
				}

				v, err := mastermindssemver.NewVersion(version)
				if err != nil {
					continue
				}

				if constraint.Check(v) {
					valid = true
					break
				}
			}

			if !valid {
				continue
			}

			var updated string

			if img.Pushed != nil {
				t := *img.Pushed
				updated = fmt.Sprintf("%s (%s)", t.In(brazil).Format("02/01/2006 15:04"), prettytime.Format(t))
			}

			table.Append([]string{
				repo,
				strings.Join(img.TagNames(), ", "),
				img.Digest,
				ByteCount(img.Size),
				updated,
			})
		}
	} else {
		for _, img := range images {
			var updated string

			if img.Pushed != nil {
				t := *img.Pushed
				updated = fmt.Sprintf("%s (%s)", t.In(brazil).Format("02/01/2006 15:04"), prettytime.Format(t))
			}

			table.Append([]string{
				repo,
				strings.Join(img.TagNames(), ", "),
				img.Digest,
				ByteCount(img.Size),
				updated,
			})
		}
	}

	table.Render()

	// Print Docker Hub link for tags
	prefix := "_/"
	if strings.Contains(repo, "/") {
		prefix = "r/"
	}

	link := fmt.Sprintf("https://hub.docker.com/%s%s?tab=tags", prefix, repo)
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
	Pushed       *time.Time
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

func convertToSemVer(s string) (string, bool) {
	match := semanticVersionRegex.FindString(s)
	if match == "" {
		return "", false
	}

	parts := strings.Split(match, ".")

	for len(parts) < 3 {
		parts = append(parts, "0")
	}

	return strings.Join(parts[:3], "."), true
}
