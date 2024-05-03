package kubectl

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

func ListEvents(namespace string, allNamespaces bool) (*Events, error) {
	cmd := []string{"get", "events", "-o=json"}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}

	out, err := runAndLogRead(cmd...)
	if err != nil {
		return nil, err
	}

	var events *Events
	if err = json.Unmarshal(out, &events); err != nil {
		return nil, errors.WithStack(err)
	}

	return events, nil
}

type Events struct {
	Items []*Event `json:"items"`
}

type Event struct {
	APIVersion         string         `json:"apiVersion"`
	Count              *int64         `json:"count,omitempty"`
	EventTime          time.Time      `json:"eventTime"`
	FirstTimestamp     time.Time      `json:"firstTimestamp"`
	InvolvedObject     InvolvedObject `json:"involvedObject"`
	Kind               string         `json:"kind"`
	LastTimestamp      time.Time      `json:"lastTimestamp"`
	Message            string         `json:"message"`
	Metadata           ItemMetadata   `json:"metadata"`
	Reason             string         `json:"reason"`
	ReportingComponent string         `json:"reportingComponent"`
	ReportingInstance  string         `json:"reportingInstance"`
	Source             Source         `json:"source"`
	Type               string         `json:"type"`
	Action             *string        `json:"action,omitempty"`
	Series             *Series        `json:"series,omitempty"`
}

type InvolvedObject struct {
	APIVersion      string  `json:"apiVersion"`
	Kind            string  `json:"kind"`
	Name            string  `json:"name"`
	Namespace       string  `json:"namespace"`
	ResourceVersion *string `json:"resourceVersion,omitempty"`
	Uid             string  `json:"uid"`
	FieldPath       *string `json:"fieldPath,omitempty"`
}

type ItemMetadata struct {
	CreationTimestamp time.Time `json:"creationTimestamp"`
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	ResourceVersion   string    `json:"resourceVersion"`
	Uid               string    `json:"uid"`
}

type Series struct {
	Count            int64     `json:"count"`
	LastObservedTime time.Time `json:"lastObservedTime"`
}

type Source struct {
	Component *string `json:"component,omitempty"`
	Host      *string `json:"host,omitempty"`
}
