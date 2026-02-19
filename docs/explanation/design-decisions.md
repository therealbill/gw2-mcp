---
title: "Design Decisions"
linkTitle: "Design Decisions"
weight: 20
description: >
  Why the GW2 MCP Server is built the way it is: the reasoning behind
  key architectural choices, what alternatives were considered, and
  what trade-offs were accepted.
---

This page explains the major design decisions in the GW2 MCP Server and the
reasoning behind them. Understanding these decisions will help you make sense of
the codebase's structure and anticipate how changes in one area affect others.

For the structural overview of the system itself, see
[Architecture](../explanation/architecture/).

## json.RawMessage for variable structures

### The problem

The Guild Wars 2 API returns deeply nested, polymorphic JSON. The `/v2/items`
endpoint alone has over 20 item subtypes -- Weapon, Armor, Consumable, Trinket,
UpgradeComponent, and many more -- each with a different `details` object. A
Weapon's details include `damage_type`, `min_power`, and `max_power`; an Armor's
details include `weight_class` and `defense`. Character data is similarly
variable: equipment, skills, specializations, build tabs, and equipment tabs all
carry deeply nested structures that differ by profession and game mode.

A traditional Go approach would define typed structs for every variant. For the
item `details` field alone, that would mean 20+ struct definitions, a
discriminator-based unmarshalling layer, and a matching re-serialization path to
produce the JSON output the LLM ultimately reads.

### The decision

Fields with variable or deeply nested schemas use `json.RawMessage` instead of
typed Go structs. In the `Item` struct, the `Details` field is declared as:

```go
Details json.RawMessage `json:"details,omitempty"`
```

The same pattern applies to the `CharacterInfo` struct, where six fields use
`json.RawMessage`:

```go
Equipment       json.RawMessage `json:"equipment,omitempty"`
Skills          json.RawMessage `json:"skills,omitempty"`
Specializations json.RawMessage `json:"specializations,omitempty"`
Training        json.RawMessage `json:"training,omitempty"`
BuildTabs       json.RawMessage `json:"build_tabs,omitempty"`
EquipmentTabs   json.RawMessage `json:"equipment_tabs,omitempty"`
```

Other types that use this approach include `Skin` (its `Details` field),
`DailyAchievements` (today/tomorrow data), and all the "pass-through" endpoints
like account unlocks, progress, dailies, guild details, and Wizard's Vault data.

### Why it works here

The consumer of this data is an LLM, not application code. The JSON travels a
short path: GW2 API response, into the Go process, back out as an MCP tool
result. The Go code never inspects weapon damage ranges or armor defense values.
It would be pure overhead to deserialize the JSON into typed structs only to
serialize it right back.

`json.RawMessage` preserves the original bytes. The JSON passes through without
a round-trip through Go's type system, which means zero risk of dropping unknown
fields that ArenaNet might add in future API versions.

### The trade-off

There is no compile-time validation of these fields. If the GW2 API changed its
`details` schema, the server would pass the new shape through without error --
which is actually the desired behavior for a passthrough proxy. The downside is
that Go code cannot access the contents of these fields without an explicit
`json.Unmarshal` call. In practice this has not been needed, because the server
never needs to interpret item subtypes.

## Typed structs where they earn their keep

Not everything uses `json.RawMessage`. Several types are defined as full Go
structs with named fields.

### Where typed structs are used

**CraftingDiscipline** -- three fields (`Discipline`, `Rating`, `Active`), a
stable schema that has not changed since the API's inception:

```go
type CraftingDiscipline struct {
    Discipline string `json:"discipline"`
    Rating     int    `json:"rating"`
    Active     bool   `json:"active"`
}
```

**ColorComponent** -- shared across four material types (`cloth`, `leather`,
`metal`, `fur`) within each `Color` entry. The struct has six numeric fields for
brightness, contrast, hue, saturation, lightness, and RGB. Defining it once
eliminates duplication across four materials.

**Achievement sub-types** -- `AchievementTier` (2 fields), `AchievementReward`
(4 fields), and `AchievementBit` (3 fields) are all small, stable structures.
The `Achievement` struct references them with typed slices, making the code
self-documenting.

**Recipe and RecipeIngredient** -- recipes have a fixed structure with output
item, ingredients, disciplines, and flags. The composite tool
`get_item_recipe_by_name` needs to iterate over ingredients to resolve their
names, which requires access to `Ingredients[i].ItemID`. A typed struct makes
this natural.

### The rule of thumb

Use typed structs when:
- The schema is stable and unlikely to gain new polymorphic variants.
- The field count is small (roughly fewer than 5 domain-specific fields).
- The Go code needs to read or manipulate individual fields (e.g., iterating
  recipe ingredients to resolve item names).

Use `json.RawMessage` when:
- The structure is polymorphic or deeply nested.
- The server is acting as a passthrough -- it never inspects the data.
- Supporting unknown future fields without code changes is valuable.

## Composite tools (Wiki + API)

### The problem

Consider what an LLM must do to answer "What's the Trading Post price of a
Mystic Coin?" without composite tools:

1. Call `wiki_search` with query "Mystic Coin" to find the wiki page.
2. Parse the result to extract the item ID (19976) from the infobox.
3. Call `get_tp_prices` with item ID 19976 to get the price data.

Each tool call consumes LLM context tokens (the request, the response, and the
reasoning between them). Three calls also create three opportunities for the LLM
to make a mistake: it might extract the wrong ID, call the wrong tool, or lose
track of intermediate state. In practice, multi-step chaining is one of the most
common failure modes in LLM tool use.

### The decision

The server provides three composite tools that collapse multi-step workflows
into a single call:

| Composite tool | What it replaces |
|---|---|
| `get_item_by_name` | wiki_search + extract ID + get_items |
| `get_item_recipe_by_name` | wiki_search + extract ID + search_recipes + get_recipes + get_items (for ingredient names) |
| `get_tp_price_by_name` | wiki_search + extract ID + get_tp_prices |

The `get_item_recipe_by_name` handler is the most involved. It searches the wiki
for the item, extracts recipe IDs from the wiki's recipe template data if
available, falls back to the API's `/v2/recipes/search` endpoint by output item
ID if not, fetches full recipe details, resolves all ingredient and output item
names, and returns an enriched result with human-readable names alongside IDs.
Without this composite tool, the LLM would need to orchestrate up to five
sequential tool calls.

### Why not make everything composite?

The underlying ID-based tools (`get_items`, `get_tp_prices`, `get_recipes`)
remain available for cases where the caller already knows the item ID. Wrapping
every tool with a wiki lookup would add unnecessary latency when IDs are already
at hand. The composite tools are an optimization for the common case of
name-based lookups, not a replacement for the ID-based tools.

### The trade-off

Composite tools are less flexible than their building blocks. If the LLM needs
to search the wiki for context beyond item IDs, it still needs `wiki_search`
directly. The composite tools also depend on wiki infobox data being present and
correctly formatted -- if a wiki page lacks an infobox `id` field, the composite
tool returns an error rather than a partial result.

## In-memory cache only

### The problem

The GW2 API has rate limits and non-trivial response times. Caching is
essential. The question is where to cache.

### The decision

All caching lives in process memory using the `patrickmn/go-cache` library.
There is no Redis, no SQLite, no disk persistence.

```go
func NewManager() *Manager {
    return &Manager{
        cache: cache.New(StaticDataTTL, CleanupInterval),
    }
}
```

### Why this is sufficient

The MCP server runs as a subprocess of the MCP client (Claude Desktop, an IDE
plugin, etc.). Its process lifetime matches the MCP session lifetime. When the
user closes their client, the server process exits. When the user opens a new
session, a fresh server process starts.

Given this lifecycle, persistence would add complexity with no benefit. Cached
data from a previous session is likely stale anyway -- trading post prices change
by the minute, wallet balances change between sessions, and daily objectives
reset. The only data that would genuinely benefit from persistence is static
metadata like item names and currency definitions, but even those are fast to
re-fetch from the GW2 API.

The cache uses different TTLs tuned to data volatility: static game data caches
for 24 hours, trading post prices for 5 minutes, delivery box contents for 2
minutes. For the full TTL schedule, see the
[Caching reference](../reference/caching/).

### The trade-off

Every new session starts cold. The first request for item metadata, currency
lists, or trading post prices hits the GW2 API. In practice, cold start latency
is not a problem because GW2 API responses typically return in under 200ms, and
the cache warms quickly as the user interacts with the server.

If the server were ever repurposed as a long-running shared service (serving
multiple users), this decision would need revisiting. But that would also
require a fundamentally different transport model.

## stdio-only transport

### The decision

The server communicates exclusively over standard input/output using the MCP
stdio transport:

```go
func (s *MCPServer) Start(ctx context.Context) error {
    s.logger.Info("Starting MCP server on stdio")
    return s.mcp.Run(ctx, &mcp.StdioTransport{})
}
```

There is no HTTP listener, no WebSocket server, no network port.

### Why this matters

This is a deliberate security decision. An MCP server that listens on a network
port is a network service -- it needs authentication, TLS, firewall rules, and
protection against the full spectrum of network attacks. The GW2 API key, which
grants access to the user's account data, would be accessible to anything that
can reach the port.

With stdio transport, the server is a subprocess with no network attack surface
beyond what the parent process exposes. The only entity that can communicate with
the server is the MCP client that spawned it. The GW2 API key passes through an
environment variable to the subprocess and never crosses a network boundary.

This also simplifies deployment. The server is a single binary with no
configuration for ports, bind addresses, TLS certificates, or CORS policies. It
runs wherever the MCP client runs.

### The trade-off

The server cannot be shared across multiple clients or accessed remotely. Each
MCP client spawns its own server process. For a personal gaming assistant, this
is the correct trade-off -- the server handles one user's data for one session.

## Excluding character bags from CharacterInfo

### The problem

The GW2 API's `/v2/characters/:name` endpoint returns everything about a
character, including the contents of every inventory bag. A character can carry
multiple bags with hundreds of item slots. Including this data in the
`CharacterInfo` response would produce a massive JSON payload that consumes
significant LLM context tokens -- mostly with inventory noise that the user
did not ask about.

### The decision

The `CharacterInfo` struct intentionally omits the `bags` field. It includes
metadata that users commonly ask about -- name, profession, level, crafting
disciplines, equipment, skills, specializations, and build tabs -- but not
per-character inventory.

For inventory access, the server provides a dedicated `get_inventory` tool that
returns the account-wide shared inventory slots, and `get_bank` / `get_materials`
for other storage. These are separate, focused tools that the LLM calls only when
the user actually asks about inventory.

### Why this matters for MCP

In an MCP context, every byte of a tool response consumes tokens from the LLM's
context window. A character info response with hundreds of inventory items could
easily reach tens of thousands of tokens, crowding out space for the
conversation. By excluding bags, the `get_characters` response stays compact and
relevant to the most common character queries (builds, equipment, crafting
progress).

### The trade-off

There is currently no tool to query a specific character's bag contents. Users
who want to know what a particular character is carrying cannot get that
information through the server. The shared inventory, bank vault, and material
storage tools cover the most common inventory questions, but per-character bag
queries remain a gap.

## Summary of trade-offs

| Decision | Benefit | Trade-off |
|---|---|---|
| `json.RawMessage` for variable data | Zero overhead passthrough, future-proof | No compile-time field validation |
| Typed structs for stable data | Self-documenting, enables field access | Must be updated if schema changes |
| Composite wiki+API tools | Fewer LLM tool calls, fewer errors | Less flexible than raw building blocks |
| In-memory cache only | Simple, no external dependencies | Cold start every session |
| stdio-only transport | No network attack surface | Single-user, single-session only |
| Exclude character bags | Compact responses, saves LLM context | No per-character inventory access |

## Related pages

- [Architecture](../explanation/architecture/) -- structural overview of the server
- [Caching reference](../reference/caching/) -- complete TTL schedule and cache key format
- [Tools reference](../reference/tools/) -- specifications for all tools including composite tools
