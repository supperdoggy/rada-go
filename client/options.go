package client

import (
	"net/http"

	"github.com/supperdoggy/vr_api/internal/profiles"
)

const defaultBaseURL = "https://itd.rada.gov.ua"

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithSearchProfileVersion(version string) Option {
	return func(c *Client) {
		c.searchProfileVersion = version
	}
}

func WithLawProjectProfileVersion(version string) Option {
	return func(c *Client) {
		c.lawProjectProfileVersion = version
	}
}

func WithProfilesRegistry(registry *profiles.Registry) Option {
	return func(c *Client) {
		c.registry = registry
	}
}
