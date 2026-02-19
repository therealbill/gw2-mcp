# GW2 MCP Server

[![Add MCP Server gw2-mcp to LM Studio](https://files.lmstudio.ai/deeplink/mcp-install-light.svg#gh-light-mode-only)](https://lmstudio.ai/install-mcp?name=gw2-mcp&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJhbHl4cGluay9ndzItbWNwOnYxIl19#gh-light-mode-only)
[![Add MCP Server gw2-mcp to LM Studio](https://files.lmstudio.ai/deeplink/mcp-install-dark.svg#gh-dark-mode-only)](https://lmstudio.ai/install-mcp?name=gw2-mcp&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJhbHl4cGluay9ndzItbWNwOnYxIl19#gh-dark-mode-only)

A Model Context Protocol (MCP) server for Guild Wars 2 that bridges Large Language Models (LLMs) with Guild Wars 2 data sources.

> **[Read the full documentation](https://therealbill.github.io/gw2-mcp/)**

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
| [Getting Started](docs/tutorials/getting-started/) | Install, configure, and make your first query |
| [How-To Guides](docs/how-to/) | Configure clients, add tools, contribute |
| [Reference](docs/reference/) | Tools, API scopes, caching, configuration |
| [Architecture](docs/explanation/) | System design and decision rationale |

## License

GNU Affero General Public License v3.0 -- see LICENSE file for details.

## Acknowledgments

- [Guild Wars 2 API](https://wiki.guildwars2.com/wiki/API:Main) for providing comprehensive game data
- [Guild Wars 2 Wiki](https://wiki.guildwars2.com/) for extensive game documentation
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) for the official Model Context Protocol implementation
