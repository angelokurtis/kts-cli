package provider

import (
	"net/http"

	"github.com/google/wire"

	"github.com/angelokurtis/kts-cli/pkg/app/dockerhub"
)

var Set = wire.NewSet(
	wire.Struct(new(http.Client)),
	dockerhub.NewClient,
)
