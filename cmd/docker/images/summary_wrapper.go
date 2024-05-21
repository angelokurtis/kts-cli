package images

import (
	"fmt"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/docker/docker/api/types/image"
	"github.com/pkg/errors"
)

type imageSummariesWrapper struct {
	wrapped []*imageSummaryWrapper
}

func (w *imageSummariesWrapper) FilterTagged() *imageSummariesWrapper {
	summaries := make([]*imageSummaryWrapper, 0)

	for _, summary := range w.wrapped {
		if len(summary.wrapped.RepoTags) > 0 {
			summaries = append(summaries, summary)
		}
	}

	return &imageSummariesWrapper{wrapped: summaries}
}

func (w *imageSummariesWrapper) Select() (*imageSummariesWrapper, error) {
	options := make([]string, 0)

	for _, summary := range w.wrapped {
		options = append(options, summary.String())
	}

	prompt := &survey.MultiSelect{
		Message: "Select images:",
		Options: options,
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	summaries := make([]*imageSummaryWrapper, 0, len(selects))

	for _, summary := range w.wrapped {
		str := summary.String()
		for _, s := range selects {
			if str == s {
				summaries = append(summaries, summary)
			}
		}
	}

	return &imageSummariesWrapper{wrapped: summaries}, nil
}

func wrapImageSummaries(summaries []image.Summary) *imageSummariesWrapper {
	wrapped := make([]*imageSummaryWrapper, 0, len(summaries))
	for _, summary := range summaries {
		wrapped = append(wrapped, &imageSummaryWrapper{wrapped: &summary})
	}

	return &imageSummariesWrapper{wrapped: wrapped}
}

type imageSummaryWrapper struct {
	wrapped *image.Summary
}

func (i *imageSummaryWrapper) String() string {
	str := fmt.Sprintf("%s | %s", byteCount(i.wrapped.Size), i.wrapped.ID)
	if len(i.wrapped.RepoTags) > 0 {
		str = str + " | " + strings.Join(i.wrapped.RepoTags, ",")
	}

	return str
}
