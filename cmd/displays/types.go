package displays

import (
	"encoding/json"
	"fmt"
	"sort"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

type DisplayList []Display

func UnmarshalDisplays(data []byte) (DisplayList, error) {
	var r DisplayList
	err := json.Unmarshal(data, &r)

	return r, err
}

func (r DisplayList) SelectMany() (DisplayList, error) {
	if len(r) == 0 {
		return DisplayList{}, nil
	}

	options := lo.KeyBy(r, func(item Display) string {
		status := ""
		if !item.Active() && !item.Connected {
			status = " " + color.Gray.Text("disconnected")
		}

		return item.Name + status
	})
	keys := lo.Keys(options)
	sort.Strings(keys)
	prompt := &survey.MultiSelect{
		Message: "Select displays:",
		Options: keys,
		Default: lo.Filter(keys, func(key string, index int) bool {
			item := options[key]
			return item.Active()
		}),
	}

	var selects []string
	if err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true)); err != nil {
		return nil, errors.WithStack(err)
	}

	return lo.Values(lo.PickByKeys(options, selects)), nil
}

func (r DisplayList) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Display struct {
	Name        string      `json:"name"`
	Connected   bool        `json:"connected"`
	Primary     *bool       `json:"primary,omitempty"`
	Resolution  *Resolution `json:"resolution,omitempty"`
	Position    *Position   `json:"position,omitempty"`
	RefreshRate *int64      `json:"refresh_rate,omitempty"`
	Rotation    *string     `json:"rotation,omitempty"`
	Modes       []Mode      `json:"modes,omitempty"`
}

type Mode struct {
	Width        int64     `json:"width"`
	Height       int64     `json:"height"`
	RefreshRates []float64 `json:"refresh_rates"`
	Preferred    bool      // true if any refresh rate in this mode has '+'
	Active       bool      // true if any refresh rate in this mode has '*'
}

type Position struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type Resolution struct {
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

func (d Display) Active() bool {
	_, ok := d.ActiveMode()
	return ok
}

func (d Display) PreferredMode() (*Mode, bool) {
	for i := range d.Modes {
		if d.Modes[i].Preferred {
			return &d.Modes[i], true
		}
	}

	return nil, false
}

func (d Display) ActiveMode() (*Mode, bool) {
	for i := range d.Modes {
		if d.Modes[i].Active {
			return &d.Modes[i], true
		}
	}

	return nil, false
}

func (d Display) String() string {
	if !d.Connected {
		return fmt.Sprintf("%s (disconnected)", d.Name)
	}

	// Select the best mode: Active > Preferred > First
	var mode *Mode
	if m, ok := d.ActiveMode(); ok {
		mode = m
	} else if m, ok := d.PreferredMode(); ok {
		mode = m
	} else if len(d.Modes) > 0 {
		mode = &d.Modes[0]
	}

	// Resolution
	resStr := "unknown"
	if mode != nil {
		resStr = fmt.Sprintf("%dx%d", mode.Width, mode.Height)
	} else if d.Resolution != nil {
		resStr = fmt.Sprintf("%dx%d", d.Resolution.Width, d.Resolution.Height)
	}

	// Refresh Rate
	rateStr := "unknown"
	if mode != nil && len(mode.RefreshRates) > 0 {
		rateStr = fmt.Sprintf("%.0fHz", mode.RefreshRates[0])
	} else if d.RefreshRate != nil {
		rateStr = fmt.Sprintf("%dHz", *d.RefreshRate)
	}

	// Rotation
	rotStr := ""
	if d.Rotation != nil {
		rotStr = fmt.Sprintf(" (%s)", *d.Rotation)
	}

	return fmt.Sprintf("%s %s @ %s%s", d.Name, resStr, rateStr, rotStr)
}
