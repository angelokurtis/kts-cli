//go:build wireinject
// +build wireinject

package repositories

import (
	"github.com/google/wire"

	"github.com/angelokurtis/kts-cli/internal/provider"
	"github.com/angelokurtis/kts-cli/pkg/app/dockerhub"
)

func newDockerhubClient() *dockerhub.Client {
	wire.Build(provider.Set)
	return &dockerhub.Client{}
}
