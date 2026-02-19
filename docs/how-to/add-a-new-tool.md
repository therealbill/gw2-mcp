---
title: Add a New Tool
---

# Add a New Tool

**Goal**: Add a new MCP tool to the GW2 MCP Server, from API struct to registered handler.

**Time**: Approximately 20 minutes

This guide walks through adding a hypothetical `get_titles` tool that retrieves title metadata from the GW2 `/v2/titles` endpoint. The pattern shown here applies to any new tool.

## Prerequisites

- Go 1.24+ installed
- Repository cloned and building (`make build` passes)
- Familiarity with the [project architecture](../explanation/architecture/)

## Steps

### 1. Define the response struct

Open `internal/gw2api/client.go` and add a struct that matches the GW2 API response shape. Place it near the other game metadata types (after the `Mini` struct is a good location).

```go
// Title represents a title from /v2/titles
type Title struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Achievement  int    `json:"achievement,omitempty"`
	Achievements []int  `json:"achievements,omitempty"`
	APRequired   int    `json:"ap_required,omitempty"`
}
```

Use the same `json` tag conventions as the existing structs: lowercase snake_case field names, `omitempty` for optional fields.

### 2. Add the API method

Still in `internal/gw2api/client.go`, add a `GetTitles` method on the `Client`. Follow the same cache-check-then-fetch pattern used by `GetMinis`:

```go
// GetTitles retrieves title metadata for the given IDs
func (c *Client) GetTitles(ctx context.Context, ids []int) ([]Title, error) {
	var results []Title
	var missingIDs []int

	for _, id := range ids {
		cacheKey := c.cache.GetTitleDetailKey(id)
		var title Title
		if c.cache.GetJSON(cacheKey, &title) {
			results = append(results, title)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		var fetched []Title
		if err := c.fetchPublic(ctx, "/titles?ids="+idsToParam(missingIDs), &fetched); err != nil {
			return nil, fmt.Errorf("failed to fetch titles: %w", err)
		}
		for _, title := range fetched {
			results = append(results, title)
			if err := c.cache.SetJSON(c.cache.GetTitleDetailKey(title.ID), title, cache.TitleDataTTL); err != nil {
				c.logger.Warn("Failed to cache title", "id", title.ID, "error", err)
			}
		}
	}

	return results, nil
}
```

Key decisions:
- **`fetchPublic`** because `/v2/titles` does not require authentication. Use `fetchAuthenticated` for endpoints that need an API key.
- **`idsToParam`** converts `[]int` to the comma-separated format the GW2 API expects.
- You will need to add `GetTitleDetailKey` and `TitleDataTTL` to the cache manager -- follow the existing pattern in `internal/cache/manager.go`.

### 3. Define the args struct

Open `internal/server/server.go` and add an args struct with `jsonschema` tags. Place it with the other arg structs near the top of the file:

```go
type GetTitlesArgs struct {
	IDs []int `json:"ids" jsonschema:"Array of title IDs to look up"`
}
```

The `jsonschema` tag becomes the parameter description that MCP clients display. The `json` tag determines the parameter name in the tool's input schema.

### 4. Write the handler

Open `internal/server/handlers.go` and add the handler function. Follow the same validate-call-return pattern as `handleGetMinis`:

```go
// handleGetTitles handles title lookup requests
func (s *MCPServer) handleGetTitles(ctx context.Context, _ *mcp.CallToolRequest, args GetTitlesArgs) (*mcp.CallToolResult, any, error) {
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Titles request", "ids", args.IDs)

	titles, err := s.gw2API.GetTitles(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get titles: %v", err))
	}

	return jsonResult(titles)
}
```

Three helpers are available for return values:
- **`errResult(msg)`** -- returns an error visible to the LLM
- **`jsonResult(v)`** -- marshals any value to indented JSON
- **`textResult(text)`** -- returns a plain text response

The handler signature must match `func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error)` where `T` is your args struct.

### 5. Register the tool

In `internal/server/server.go`, add a `mcp.AddTool` call inside the `registerTools()` method. Place it in the appropriate section (Game Metadata for this example):

```go
mcp.AddTool(s.mcp, &mcp.Tool{
	Name:        "get_titles",
	Description: "Get title metadata (name, achievement requirements) for given title IDs.",
}, s.handleGetTitles)
```

The `Name` field is the tool name exposed to MCP clients. The `Description` is what LLMs see when deciding which tool to call -- keep it specific and concise.

### 6. Add tests

Open `internal/server/handlers_test.go` and add test cases. For handlers that call external APIs, test the input validation logic:

```go
func TestHandleGetTitles_Validation(t *testing.T) {
	s := &MCPServer{}
	ctx := context.Background()

	// Test empty IDs
	result, _, err := s.handleGetTitles(ctx, nil, GetTitlesArgs{IDs: []int{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsError {
		t.Error("expected error result for empty IDs")
	}
}
```

Add the necessary imports (`context`, `testing`) to the test file if they are not already present.

### 7. Build and verify

Run the build and test suite to confirm everything compiles and passes:

```bash
make build && go test ./...
```

Expected output:

```
ok  	github.com/AlyxPink/gw2-mcp/internal/server	0.XXXs
ok  	github.com/AlyxPink/gw2-mcp/internal/cache	0.XXXs
ok  	github.com/AlyxPink/gw2-mcp/internal/gw2api	0.XXXs
```

## Verify it works

Start the server and use an MCP client to call your new tool:

```bash
GW2_API_KEY=your-key ./gw2-mcp
```

From Claude Desktop or another MCP client, invoke:

```
get_titles with ids [1, 2, 3]
```

You should receive a JSON response with title names and metadata.

## Troubleshooting

### Problem: "undefined: GetTitleDetailKey" compile error
**Symptom**: `make build` fails with an undefined reference in `client.go`.
**Cause**: The cache key function has not been added yet.
**Solution**: Add `GetTitleDetailKey(id int) string` to `internal/cache/manager.go` following the pattern of `GetMiniDetailKey`, and add a `TitleDataTTL` constant to the TTL block.

### Problem: Tool does not appear in MCP client
**Symptom**: The server starts but the client does not list `get_titles`.
**Cause**: The `mcp.AddTool` call is missing or the handler function signature does not match.
**Solution**: Verify the `mcp.AddTool` call is inside `registerTools()` and that the handler accepts the correct args struct type. The generic type parameter is inferred from the handler signature.

### Problem: "failed to fetch titles" at runtime
**Symptom**: The tool returns an error when called with valid IDs.
**Cause**: The API path or query parameter format is wrong.
**Solution**: Verify the path matches the GW2 API docs. Test the URL directly: `curl "https://api.guildwars2.com/v2/titles?ids=1,2,3"`.

## Next steps

- [Tools reference](../reference/tools/) -- add your new tool to the reference docs
- [Architecture](../explanation/architecture/) -- understand the layered design
- [Contribute](contribute/) -- submit your tool as a pull request

## See also

- [GW2 API documentation](https://wiki.guildwars2.com/wiki/API:Main) -- official endpoint reference
- [Tools reference](../reference/tools/) -- specifications for all existing tools
