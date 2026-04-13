package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/supperdoggy/vr_api/internal/parser"
	"github.com/supperdoggy/vr_api/internal/profiles"
)

type Client struct {
	baseURL string

	httpClient *http.Client
	loader     parser.DocumentLoader
	registry   *profiles.Registry

	searchProfileVersion     string
	lawProjectProfileVersion string
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
		loader:     parser.NewGoqueryLoader(),
		registry:   profiles.NewRegistry(),
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	if c.registry == nil {
		c.registry = profiles.NewRegistry()
	}

	return c
}

func (c *Client) ParseSearchHTML(ctx context.Context, html []byte) (SearchResponse, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	profile, err := c.registry.Search(c.searchProfileVersion)
	if err != nil {
		return SearchResponse{}, err
	}

	p := parser.NewSearchHTMLParser(c.loader, profile.Selectors())
	resp, err := p.Parse(ctx, html)
	if err != nil {
		return SearchResponse{}, err
	}

	for i := range resp.Items {
		resp.Items[i].URL = c.normalizeURL(resp.Items[i].URL)
	}

	return resp, nil
}

func (c *Client) ParseSearchHTMLString(ctx context.Context, html string) (SearchResponse, error) {
	return c.ParseSearchHTML(ctx, []byte(html))
}

func (c *Client) ParseLawProjectHTML(ctx context.Context, html []byte) (LawProjectDetails, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	profile, err := c.registry.LawProject(c.lawProjectProfileVersion)
	if err != nil {
		return LawProjectDetails{}, err
	}

	p := parser.NewLawProjectHTMLParser(c.loader, profile.Selectors())
	resp, err := p.Parse(ctx, html)
	if err != nil {
		return LawProjectDetails{}, err
	}

	if resp.SourceURL == "" {
		resp.SourceURL = c.baseURL
	} else {
		resp.SourceURL = c.normalizeURL(resp.SourceURL)
	}
	for i := range resp.Documents {
		resp.Documents[i].URL = c.normalizeURL(resp.Documents[i].URL)
	}

	return resp, nil
}

func (c *Client) ParseLawProjectHTMLString(ctx context.Context, html string) (LawProjectDetails, error) {
	return c.ParseLawProjectHTML(ctx, []byte(html))
}

func (c *Client) Search(ctx context.Context, query SearchQuery) (SearchResponse, error) {
	_ = ctx
	_ = query
	return SearchResponse{}, ErrNotImplemented
}

func (c *Client) LawProjectDetails(ctx context.Context, projectID string) (LawProjectDetails, error) {
	_ = ctx
	_ = projectID
	return LawProjectDetails{}, ErrNotImplemented
}

func (c *Client) normalizeURL(raw string) string {
	if raw == "" {
		return ""
	}

	base, err := url.Parse(c.baseURL)
	if err != nil {
		return raw
	}
	target, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return base.ResolveReference(target).String()
}
