package dockerhub

import (
	"net/http"
)

const baseURL = "https://hub.docker.com"

type Client struct {
	client *http.Client
}

func NewClient(client *http.Client) *Client {
	return &Client{client: client}
}
