---
title: "Architecture"
summary: "How the GW2 MCP Server is structured and why"
prerequisites: []
est_time: "15 min"
roles: ["developer", "contributor"]
stability: "stable"
---

# Understanding the Architecture

The GW2 MCP Server bridges two very different systems: the [Model Context Protocol](https://modelcontextprotocol.io/) (MCP), which lets AI assistants discover and call tools over JSON-RPC, and the [Guild Wars 2 API](https://wiki.guildwars2.com/wiki/API:Main), which exposes game data over HTTPS REST endpoints. This document explains how the server is organized, how data flows through it, and the reasoning behind the key design choices.

## The problem

An MCP client (such as Claude Desktop or another AI host) needs structured access to Guild Wars 2 game data. The GW2 API is powerful but presents several challenges for direct AI consumption:

- **Multi-step lookups are common.** Finding an item's Trading Post price by name requires searching the wiki for the item ID, then querying the GW2 API with that ID. An AI would need multiple tool calls and intermediate reasoning to accomplish this.
- **The API uses numeric IDs everywhere.** Endpoints return item IDs, recipe IDs, and currency IDs that are meaningless without a follow-up lookup for names and metadata.
- **Rate limiting and latency matter.** The GW2 API is a shared public resource, and every unnecessary request adds latency to what should feel like an instant tool call.
- **Authentication is per-user.** Account-specific data (wallet, bank, characters) requires an API key, and the server must handle this cleanly.

The architecture exists to solve these problems: collapse multi-step operations into single tool calls, enrich numeric IDs with human-readable names, cache aggressively to reduce latency and API load, and isolate per-user data safely.

## Package structure

The codebase is organized into four internal packages, each with a single clear responsibility. The separation follows Go convention (the `internal/` directory prevents external imports) and keeps the dependency graph shallow.

```
main.go                     Reads config, wires dependencies, starts server
internal/
  server/
    server.go               MCPServer struct, tool registration, arg structs
    handlers.go             Handler implementations, composite tool logic
  gw2api/
    client.go               GW2 API client, struct definitions, caching
  wiki/
    client.go               Wiki search, infobox parsing, recipe extraction
  cache/
    manager.go              In-memory cache with per-key TTLs
```

### `internal/server/` -- MCP protocol layer

This package is the glue between the MCP protocol and the domain clients. It contains two files:

**`server.go`** defines the `MCPServer` struct, which holds references to the MCP server, the GW2 API client, the wiki client, and the cache manager. It also defines all the argument structs (like `WikiSearchArgs`, `GetItemsArgs`, `GetTPPriceByNameArgs`) and the `registerTools()` method that wires each tool name to its handler. This is where you look to understand what tools exist and what parameters they accept.

**`handlers.go`** implements the handler functions. Most handlers follow a simple pattern: validate input, call a domain client method, return the result as JSON. The more interesting handlers are the composite tools at the bottom of the file, which orchestrate calls across both the wiki and GW2 API clients.

The separation between these two files is deliberate. `server.go` is a declaration of the server's surface area -- its tools, their schemas, and their wiring. `handlers.go` is the implementation of behavior. When adding a new tool, you touch both files: define the struct and register the tool in `server.go`, implement the handler in `handlers.go`.

### `internal/gw2api/` -- GW2 API client

This single-file package (`client.go`) handles all communication with `https://api.guildwars2.com/v2`. It is responsible for:

- **Struct definitions.** All the Go types that model GW2 API responses (`Item`, `Recipe`, `PriceInfo`, `AccountInfo`, `WalletInfo`, and many more) live here.
- **HTTP request execution.** Helper methods like `fetchPublic()`, `fetchAuthenticated()`, `fetchPublicRaw()`, and `fetchAuthenticatedRaw()` handle the mechanics of building requests, setting headers, checking status codes, and decoding JSON.
- **Cache integration.** Every public method (like `GetItems`, `GetPrices`, `GetWallet`) checks the cache before making an HTTP request, and populates the cache after a successful fetch.
- **Data enrichment.** Methods like `GetPrices` and `GetBank` automatically resolve item IDs to names by calling `GetItems` internally, so callers always receive human-readable results.
- **Authentication.** The client stores the API key at construction time and uses it for authenticated endpoints. The key is never logged or cached directly; instead, a SHA-256 hash of the key is used for cache key namespacing.

Everything lives in one file because the package has a single type (`Client`) with a consistent pattern across all its methods. Adding a new endpoint means adding a struct, a public method, and a fetch helper -- all in the same file, following the same pattern.

### `internal/wiki/` -- Wiki integration

This single-file package (`client.go`) communicates with the Guild Wars 2 Wiki's MediaWiki API at `https://wiki.guildwars2.com/api.php`. It provides:

- **Full-text search.** The `Search` method queries the wiki and returns results with titles, snippets, and URLs.
- **Infobox parsing.** For each search result, the client fetches the page's wikitext and parses `{{Infobox}}` templates to extract structured key-value data (item IDs, rarity, level, etc.). This is what enables the composite tools -- the item ID extracted from a wiki infobox is the bridge to the GW2 API.
- **Recipe template extraction.** The `parseRecipes` function finds all `{{Recipe}}` templates in a page's wikitext and extracts their fields (ingredients, quantities, crafting disciplines). This allows `get_item_recipe_by_name` to return recipe data directly from the wiki when available.
- **Markup cleaning.** Wiki text contains MediaWiki markup (`[[links]]`, `'''bold'''`). The `cleanWikiMarkup` function strips this to produce clean text for the structured results.

The wiki client exists as a separate package from the GW2 API client because it talks to a completely different service (MediaWiki vs. the GW2 REST API), uses different request patterns, and has its own parsing logic. The server package is what brings them together.

### `internal/cache/` -- Caching layer

This single-file package (`manager.go`) wraps the `patrickmn/go-cache` library to provide typed, TTL-aware caching. It defines:

- **Cache key templates.** Structured key patterns like `"item:detail:%d"`, `"wallet:%s"`, `"tp:price:%d"` ensure predictable, collision-free key generation.
- **TTL constants.** Every data category has its own TTL, from 1 year for truly static data (currency metadata) down to 2 minutes for rapidly changing data (trading post delivery).
- **JSON serialization helpers.** `GetJSON` and `SetJSON` methods handle marshal/unmarshal so callers work with Go structs rather than raw bytes.

The cache manager is a shared dependency -- both the GW2 API client and the wiki client receive a reference to the same `Manager` at construction time. This means cached item metadata is available regardless of which code path fetched it first. For a detailed breakdown of TTL values, see the [caching reference](../reference/caching/).

## Request flow

To understand how these packages work together, trace a tool call from start to finish. Here is the path for `get_tp_price_by_name`, one of the composite tools:

```
MCP Client (e.g. Claude Desktop)
  |
  | JSON-RPC over stdio: tools/call { name: "get_tp_price_by_name", args: { name: "Mystic Coin" } }
  v
mcp.Server (go-sdk)
  |
  | Routes to registered handler based on tool name
  v
MCPServer.handleGetTPPriceByName()       [internal/server/handlers.go]
  |
  | 1. Validate input: name must not be empty
  | 2. Search wiki for "Mystic Coin"
  v
wiki.Client.Search()                     [internal/wiki/client.go]
  |
  | Check cache for "wiki:search:mystic coin"
  | On miss: HTTP GET to wiki API, parse results, extract infobox, cache result
  | Return: SearchResponse with Infobox containing { "id": "19976", ... }
  v
handleGetTPPriceByName() continued
  |
  | 3. Extract item ID 19976 from wiki infobox
  | 4. Fetch trading post prices for item 19976
  v
gw2api.Client.GetPrices()               [internal/gw2api/client.go]
  |
  | Check cache for "tp:price:19976"
  | On miss: HTTP GET to api.guildwars2.com/v2/commerce/prices?ids=19976
  | Cache result with 5-minute TTL
  | Enrich with item name via GetItems() (also cached)
  v
handleGetTPPriceByName() continued
  |
  | 5. Serialize PriceInfo to JSON
  v
mcp.Server
  |
  | JSON-RPC response over stdio
  v
MCP Client
```

A simpler tool like `get_wallet` has a shorter path -- the handler calls `gw2API.GetWallet()` directly, which checks the cache and falls back to the GW2 API. But the overall shape is the same: MCP client sends JSON-RPC over stdio, the SDK routes to a handler, the handler calls domain clients, domain clients use the cache, and the result flows back as serialized JSON over stdio.

The key architectural point here is that **stdio is the only transport**. The server communicates exclusively through standard input and output. There is no HTTP server, no TCP listener, no port to configure. The MCP host process (Claude Desktop, for example) spawns the server binary and communicates via its stdin/stdout pipes. This makes deployment simple and security straightforward -- the server process inherits its permissions from the host.

## Caching architecture

Caching is central to the server's performance model. Without it, a single composite tool call like `get_item_recipe_by_name` could trigger four or more HTTP requests (wiki search, wiki page fetch, recipe API call, item name resolution). With caching, repeated queries for the same data are essentially free.

### How it works

The cache is an in-memory key-value store with per-key expiration. When a domain client method is called:

1. It constructs a cache key from a template (e.g., `"item:detail:19976"`).
2. It calls `cache.GetJSON(key, &dest)` to check for a cached value.
3. On hit: return the cached value immediately.
4. On miss: fetch from the external API, then call `cache.SetJSON(key, value, ttl)` to store the result.

A background goroutine runs every 10 minutes (`CleanupInterval`) to evict expired entries.

### Per-user cache isolation

Account-specific data must not leak between users. The server achieves this by incorporating a SHA-256 hash of the API key into cache keys for authenticated data:

```
wallet:<sha256-hash>
bank:<sha256-hash>
characters:<sha256-hash>
```

The hash serves two purposes: it prevents the raw API key from appearing in memory as a cache key, and it naturally partitions data per user. Public data (items, recipes, currencies) uses shared keys without a user prefix, since this data is the same for all users.

### TTL strategy

Different data categories have different volatility, and the TTL values reflect this:

| Category | TTL | Rationale |
|----------|-----|-----------|
| Currency metadata | 1 year | Currencies are added with game expansions, not daily |
| Item/recipe/skin metadata | 24 hours | Game data changes only with patches |
| Wiki content | 24 hours | Wiki edits are infrequent for most pages |
| Trading Post prices | 5 minutes | Prices fluctuate constantly |
| Account data (wallet, bank) | 5 minutes | Players actively changing inventory |
| TP delivery box | 2 minutes | Users want to see pickups promptly |
| Wizard's Vault (public) | 1 hour | Season data changes slowly |

This strategy means the first call for a given item's metadata pays the full latency cost, but subsequent calls within the TTL window return instantly. For exact TTL values for every data category, see the [caching reference](../reference/caching/).

### Trade-offs

The in-memory cache is simple and fast, but it comes with trade-offs:

| Benefit | Trade-off |
|---------|-----------|
| Zero external dependencies | Cache is lost on process restart |
| Microsecond lookup latency | Memory usage grows with cache size |
| No configuration required | No sharing between server instances |
| Simple implementation | No persistence across sessions |

For a single-user MCP server running as a subprocess, these trade-offs are appropriate. The server starts fresh when launched, quickly warms its cache through normal usage, and the memory overhead for cached game data is modest (typically a few megabytes at most).

## Tool registration pattern

The server uses the official `modelcontextprotocol/go-sdk` library for MCP protocol handling. Tool registration follows a specific pattern that is worth understanding because it drives how the server's API surface is defined.

### How `mcp.AddTool` works

Each tool is registered by calling `mcp.AddTool` with three arguments: the MCP server instance, a tool descriptor (name and description), and a handler function:

```go
mcp.AddTool(s.mcp, &mcp.Tool{
    Name:        "get_items",
    Description: "Get item metadata for given item IDs.",
}, s.handleGetItems)
```

The handler function's signature includes a typed argument struct:

```go
func (s *MCPServer) handleGetItems(
    ctx context.Context,
    _ *mcp.CallToolRequest,
    args GetItemsArgs,
) (*mcp.CallToolResult, any, error)
```

The SDK uses the `GetItemsArgs` struct to **automatically generate the JSON Schema** for the tool's parameters. The struct's `json` tags define parameter names, and `jsonschema` tags provide descriptions:

```go
type GetItemsArgs struct {
    IDs []int `json:"ids" jsonschema:"Array of item IDs to look up"`
}
```

This means the tool's parameter schema is derived directly from the Go type system. There is no separate schema file to maintain, no OpenAPI spec to keep in sync. When a client calls `tools/list`, the SDK reflects the struct and returns the generated schema. When a client calls `tools/call`, the SDK deserializes the JSON arguments into the struct and passes it to the handler.

### Why this matters

This pattern has a significant effect on how AI clients interact with the server. The MCP protocol includes a `tools/list` endpoint that returns all registered tools with their names, descriptions, and parameter schemas. AI clients use this to understand what tools are available and how to call them -- this is the auto-discovery mechanism.

The quality of the tool descriptions and parameter schema descriptions directly affects how well an AI client uses the tools. The descriptions in the `mcp.Tool` struct and the `jsonschema` tags are not just documentation -- they are the interface contract that AI clients rely on to generate correct tool calls.

### Tools without parameters

Some tools take no parameters (e.g., `get_wallet`, `get_account`). These use `struct{}` as the args type:

```go
func (s *MCPServer) handleGetWallet(
    ctx context.Context,
    _ *mcp.CallToolRequest,
    _ struct{},
) (*mcp.CallToolResult, any, error)
```

The SDK sees the empty struct and generates a schema with no required parameters. For the practical steps of adding a new tool using this pattern, see [How to add a new tool](../how-to/add-a-new-tool/).

## Wiki integration and composite tools

The most architecturally interesting aspect of the server is how it combines wiki data and API data into composite tools. This is the feature that transforms the server from a simple API proxy into something genuinely more useful.

### The problem composite tools solve

Consider what an AI client would need to do to find the Trading Post price of "Mystic Coin" using only basic tools:

1. Call `wiki_search` with query "Mystic Coin" to find the wiki page.
2. Parse the returned infobox data to find the item ID (19976).
3. Call `get_tp_prices` with item ID 19976.

That is three steps requiring intermediate reasoning. With the composite tool `get_tp_price_by_name`, it becomes a single call: the server handles the wiki search, ID extraction, and price lookup internally.

### How wiki data enables this

The wiki client does not just return search results -- it fetches each result page's full wikitext and parses it for structured data. Two parsing functions make this work:

**`parseInfobox`** finds `{{Infobox}}` templates (like `{{Item infobox}}`, `{{NPC infobox}}`) and extracts their key-value parameters. For an item page, this typically includes `id`, `name`, `rarity`, `type`, and other fields. The extracted `id` field is the numeric GW2 API item ID -- the bridge between wiki knowledge and API data.

**`parseRecipes`** finds all `{{Recipe}}` templates on a page and extracts their parameters (output item, ingredients, quantities, crafting disciplines). This provides recipe data that might otherwise require multiple API calls to discover.

Both parsers handle nested templates (templates within templates) by tracking brace depth, and they strip wiki markup from values to return clean text.

### The three composite tools

The server currently provides three composite tools:

- **`get_item_by_name`** -- Searches the wiki for an item name, extracts the item ID from the infobox, then calls `GetItems` to return full item metadata from the API.
- **`get_item_recipe_by_name`** -- Searches the wiki, extracts recipe IDs from `{{Recipe}}` templates (falling back to the API's recipe search endpoint if the wiki does not have them), fetches full recipe details, and resolves all ingredient item IDs to names. The result is a fully enriched recipe with human-readable ingredient names.
- **`get_tp_price_by_name`** -- Searches the wiki, extracts the item ID, and fetches current Trading Post buy/sell prices.

Each of these collapses what would be a multi-step, multi-tool interaction into a single call. This is possible because the server can hold context across internal operations that an MCP client would otherwise need to manage externally.

### Fallback strategy

The composite tools include a fallback path for recipe lookups. If the wiki page does not contain `{{Recipe}}` templates (which happens for some items), `handleGetItemRecipeByName` falls back to the GW2 API's `/v2/recipes/search?output=<itemID>` endpoint. This means the tool still works even when the wiki lacks recipe templates, at the cost of one additional API call.

## Related topics

### Practical guides
- [How to add a new tool](../how-to/add-a-new-tool/) -- Step-by-step process for extending the server

### Reference documentation
- [Caching reference](../reference/caching/) -- Exact TTL values, key templates, and cache behavior

### Deeper rationale
- [Design decisions](../explanation/design-decisions/) -- Why specific technical choices were made
