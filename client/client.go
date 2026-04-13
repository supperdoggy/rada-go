package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strconv"

	"github.com/supperdoggy/vr_api/internal/parser"
	"github.com/supperdoggy/vr_api/internal/profiles"
	"golang.org/x/net/html/charset"
)

type Client struct {
	baseURL string

	httpClient *http.Client
	loader     parser.DocumentLoader
	registry   *profiles.Registry

	searchProfileVersion     string
	lawProjectProfileVersion string
}

func New(baseURL string, opts ...Option) *Client {
	if baseURL != "" {
		opts = append([]Option{WithBaseURL(baseURL)}, opts...)
	}

	return NewClient(opts...)
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

	c.normalizeLawProjectDetails(&resp)

	return resp, nil
}

func (c *Client) ParseLawProjectHTMLString(ctx context.Context, html string) (LawProjectDetails, error) {
	return c.ParseLawProjectHTML(ctx, []byte(html))
}

func (c *Client) Search(params SearchParams) (SearchResponse, error) {
	return c.SearchContext(context.Background(), params)
}

func (c *Client) SearchContext(ctx context.Context, params SearchParams) (SearchResponse, error) {
	body, err := c.fetch(ctx, c.searchURL(params))
	if err != nil {
		return SearchResponse{}, err
	}

	return c.ParseSearchHTML(ctx, body)
}

func (c *Client) Get(projectID string) (LawProjectDetails, error) {
	return c.GetContext(context.Background(), projectID)
}

func (c *Client) GetContext(ctx context.Context, projectID string) (LawProjectDetails, error) {
	targetURL := c.lawProjectURL(projectID)
	body, err := c.fetch(ctx, targetURL)
	if err != nil {
		return LawProjectDetails{}, err
	}

	resp, err := c.ParseLawProjectHTML(ctx, body)
	if err != nil {
		return LawProjectDetails{}, err
	}

	resp.SourceURL = targetURL
	if resp.VotingResults != nil && resp.VotingResults.URL != "" {
		summary, err := c.fetchVotingResultsSummary(ctx, resp.VotingResults.URL)
		if err == nil && len(summary) > 0 {
			resp.VotingResults.Summary = summary
		}
	}

	return resp, nil
}

func (c *Client) LawProjectDetails(ctx context.Context, projectID string) (LawProjectDetails, error) {
	return c.GetContext(ctx, projectID)
}

func (c *Client) fetch(ctx context.Context, target string) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, target)
	}

	reader, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		reader = resp.Body
	}

	return io.ReadAll(reader)
}

func (c *Client) searchURL(params SearchParams) string {
	values := url.Values{}
	if params.Term != "" {
		values.Set("term", params.Term)
	}
	if params.Page > 0 {
		values.Set("page", strconv.Itoa(params.Page))
	}
	if params.PerPage > 0 {
		values.Set("perPage", strconv.Itoa(params.PerPage))
	}

	if len(params.Filters) > 0 {
		keys := make([]string, 0, len(params.Filters))
		for key := range params.Filters {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			values.Set(key, params.Filters[key])
		}
	}

	return c.resolveURL("/search", values)
}

func (c *Client) lawProjectURL(projectID string) string {
	return c.resolveURL(path.Join("/billinfo/Bills/Card", projectID), nil)
}

func (c *Client) resolveURL(route string, values url.Values) string {
	base, err := url.Parse(c.baseURL)
	if err != nil {
		return c.baseURL
	}

	target, err := url.Parse(route)
	if err != nil {
		return c.baseURL
	}

	resolved := base.ResolveReference(target)
	if len(values) > 0 {
		resolved.RawQuery = values.Encode()
	}

	return resolved.String()
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

func (c *Client) normalizeLawProjectDetails(resp *LawProjectDetails) {
	if resp.SourceURL == "" {
		resp.SourceURL = c.baseURL
	} else {
		resp.SourceURL = c.normalizeURL(resp.SourceURL)
	}

	if resp.Act != nil {
		resp.Act.URL = c.normalizeURL(resp.Act.URL)
	}
	for i := range resp.Initiators {
		resp.Initiators[i].URL = c.normalizeURL(resp.Initiators[i].URL)
	}
	if resp.MainCommittee != nil {
		resp.MainCommittee.URL = c.normalizeURL(resp.MainCommittee.URL)
	}
	for i := range resp.OtherCommittees {
		resp.OtherCommittees[i].URL = c.normalizeURL(resp.OtherCommittees[i].URL)
	}
	for i := range resp.Committees {
		resp.Committees[i].URL = c.normalizeURL(resp.Committees[i].URL)
	}
	for i := range resp.Documents {
		resp.Documents[i].URL = c.normalizeURL(resp.Documents[i].URL)
	}
	for i := range resp.RelatedBills {
		resp.RelatedBills[i].URL = c.normalizeURL(resp.RelatedBills[i].URL)
	}
	resp.ChronologyURL = c.normalizeURL(resp.ChronologyURL)
	if resp.VotingResults != nil {
		resp.VotingResults.URL = c.normalizeURL(resp.VotingResults.URL)
	}
}

func (c *Client) fetchVotingResultsSummary(ctx context.Context, target string) ([]VoteSummaryItem, error) {
	body, err := c.fetch(ctx, c.normalizeURL(target))
	if err != nil {
		return nil, err
	}

	p := parser.NewVotingResultsHTMLParser(c.loader)
	return p.Parse(ctx, body)
}
