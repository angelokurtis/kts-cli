package displays

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var Command = &cobra.Command{
	Use: "displays",
	Run: system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "select", Run: wrapWithErrorHandler(selectMany)})
}

func selectMany(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	displays, err := listDisplays(ctx)
	if err != nil {
		return err
	}

	selected, err := displays.SelectMany()
	if err != nil {
		return err
	}

	ordered, err := promptForDisplayOrder(selected)
	if err != nil {
		return err
	}

	return configureDisplays(ctx, displays, ordered...)
}

// promptForDisplayOrder prompts the user to assign a position to each connected display
func promptForDisplayOrder(displays DisplayList) ([]string, error) {
	type entry struct {
		name string
		pos  int
	}

	sort.Slice(displays, func(i, j int) bool {
		return displays[i].Name < displays[j].Name
	})

	var entries []entry
	used := make(map[int]struct{})

	for _, d := range displays {
		var opts []string

		for i := 1; i <= len(displays); i++ {
			if _, ok := used[i]; !ok {
				opts = append(opts, strconv.Itoa(i))
			}
		}

		if len(opts) == 0 {
			continue
		}

		var pos int
		if len(opts) == 1 {
			pos, _ = strconv.Atoi(opts[0])
		} else {
			var input string

			q := &survey.Select{
				Message: fmt.Sprintf("Select position for display '%s':", d.Name),
				Options: opts,
			}
			if err := survey.AskOne(q, &input); err != nil {
				return nil, err
			}

			pos, _ = strconv.Atoi(input)
		}

		used[pos] = struct{}{}

		entries = append(entries, entry{name: d.Name, pos: pos})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].pos < entries[j].pos
	})

	ordered := make([]string, len(entries))
	for i, e := range entries {
		ordered[i] = e.name
	}

	return ordered, nil
}

func configureDisplays(ctx context.Context, displays DisplayList, ordered ...string) error {
	if len(displays) == 0 {
		return nil
	}

	m := make(map[string]Display, len(displays))
	for _, d := range displays {
		m[d.Name] = d
	}

	active := make(map[string]bool, len(ordered))

	for _, name := range ordered {
		d, ok := m[name]
		if !ok {
			return fmt.Errorf("unknown display: %s", name)
		}

		if !d.Connected {
			return fmt.Errorf("display %s is not connected", name)
		}

		active[name] = true
	}

	var x int64

	for i, name := range ordered {
		d := m[name]

		mode, ok := d.PreferredMode()
		if !ok {
			return fmt.Errorf("no preferred mode for display: %s", name)
		}

		if err := enableDisplay(ctx, name, mode.Width, mode.Height, x, i == 0); err != nil {
			return err
		}

		x += mode.Width
	}

	for name, d := range m {
		if !active[name] && d.Connected {
			if err := disableDisplay(ctx, name); err != nil {
				return err
			}
		}
	}

	return nil
}
