---
title: GW2 MCP Server
bookCollapseSection: false
---

# GW2 MCP Server

A Model Context Protocol (MCP) server for Guild Wars 2 that bridges Large Language Models (LLMs) with Guild Wars 2 data sources.

## Quick Start

1. Create a GW2 API key at [account.arena.net/applications](https://account.arena.net/applications) with the scopes you need.

2. Add to your MCP client config (Claude Desktop, Claude Code, LM Studio, etc.):

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "-e", "GW2_API_KEY", "alyxpink/gw2-mcp:v1"],
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

3. Ask your LLM about Guild Wars 2 -- "Check my wallet", "Look up Mystic Coin", "What are today's dailies?"

## Features

- **37 MCP tools** covering account data, Trading Post, achievements, guilds, Wizard's Vault, wiki search, and game metadata
- **Composite tools** that chain wiki search with API lookups in a single call (`get_item_by_name`, `get_tp_price_by_name`, `get_item_recipe_by_name`)
- **Smart caching** with per-data-type TTLs (2 minutes for live data up to 1 year for static metadata)
- **Graceful degradation** -- works without an API key; authenticated tools return clear errors
- **Docker and binary** distribution options

## Documentation

| Section | Description |
|---------|-------------|
| [Getting Started](tutorials/getting-started/) | Install, configure, and make your first query |
| [How-To Guides](how-to/) | Configure clients, add tools, contribute |
| [Reference](reference/) | Tools, API scopes, caching, configuration |
| [Architecture](explanation/) | System design and decision rationale |
