# VR API HTML Client (Scaffold)

Go client scaffold for parsing Verkhovna Rada ITD HTML into stable, app-friendly JSON models.

## Scope (Current)

- No authentication (public pages only)
- HTML-in parsing API (`[]byte` and `string`)
- Search and law-project detail parsing entrypoints
- Versioned selector profiles for quick layout-fix updates
- Future HTTP fetch methods declared but intentionally not implemented yet

## Module

`github.com/supperdoggy/vr_api`

## Default Target

`https://itd.rada.gov.ua`

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
	"context"
	"log"

	"github.com/supperdoggy/vr_api/client"
)

func main() {
	c := client.NewClient()
	resp, err := c.ParseSearchHTMLString(context.Background(), `<div id="SearchResultContainer"></div>`)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("search results: %d", resp.Count)
}
```

## Development

```bash
make test
make fmt
make lint
```
