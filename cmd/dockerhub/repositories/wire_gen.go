// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package repositories

import (
	"github.com/angelokurtis/kts-cli/pkg/app/dockerhub"
	"net/http"
)

// Injectors from wire.go:

func newDockerhubClient() *dockerhub.Client {
	client := &http.Client{}
	dockerhubClient := dockerhub.NewClient(client)
	return dockerhubClient
}
