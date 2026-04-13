# VR API Client

Go client skeleton for fetching and parsing law-project HTML into stable, app-friendly JSON models.

## Scope

- No authentication (public pages only)
- Simple client entrypoint: `client.New(url)`
- `Search(params)` for law-project search
- `Get(id)` for law-project details
- HTML parsing helpers (`[]byte` and `string`)
- Versioned selector profiles for quick layout-fix updates

## Module

`github.com/supperdoggy/vr_api`

## Default Target

`https://itd.rada.gov.ua`

The HTTP client currently assumes:

- `GET /search` for search results
- `GET /bill/{id}` for law-project details
- HTML compatible with the selector profiles in `internal/profiles/`

## Architecture

- `client/` public API (single consumer entrypoint)
- `schema/` stable normalized JSON contracts
- `internal/parser/` HTML loader + parser implementations
- `internal/profiles/` versioned selector profiles + registry
- `internal/apperr/` typed parser errors

## Quick Start

```go
package main

import (
	"log"

	"github.com/supperdoggy/vr_api/client"
)

func main() {
	c := client.New("https://example.com")

	resp, err := c.Search(client.SearchParams{
		Term: "budget",
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("search results: %d", resp.Count)

	details, err := c.Get(resp.Items[0].ID)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("law project: %s", details.Title)
}
```

## Development

```bash
make test
make fmt
make lint
```
