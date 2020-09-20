package gcloud

import (
	"encoding/json"
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
)

func ListNodePools(cluster string) (interface{}, error) {
	out, err := bash.RunAndLogRead(fmt.Sprintf("gcloud container node-pools list --cluster %s", cluster))
	if err != nil {
		return nil, err
	}

	var projects []*Project
	if err := json.Unmarshal(out, &projects); err != nil {
		return nil, errors.WithStack(err)
	}

	return projects, nil
}
