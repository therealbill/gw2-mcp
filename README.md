# GW2 MCP Server

[![Add MCP Server gw2-mcp to LM Studio](https://files.lmstudio.ai/deeplink/mcp-install-light.svg#gh-light-mode-only)](https://lmstudio.ai/install-mcp?name=gw2-mcp&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJhbHl4cGluay9ndzItbWNwOnYxIl19#gh-light-mode-only)
[![Add MCP Server gw2-mcp to LM Studio](https://files.lmstudio.ai/deeplink/mcp-install-dark.svg#gh-dark-mode-only)](https://lmstudio.ai/install-mcp?name=gw2-mcp&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJhbHl4cGluay9ndzItbWNwOnYxIl19#gh-dark-mode-only)

A Model Context Provider (MCP) server for Guild Wars 2 that bridges Large Language Models (LLMs) with Guild Wars 2 data sources.

## Features

- **34 Tools** covering account, Trading Post, achievements, guilds, Wizard's Vault, game metadata, and more
- **Wiki Search**: Search and retrieve content from the Guild Wars 2 wiki
- **Account Data**: Wallet, bank, materials, shared inventory, characters, unlocks, progress, dailies
- **Trading Post**: Prices, listings, gem exchange, delivery box, transaction history
- **Game Data**: Items, skins, recipes, achievements, colors, minis, mounts, dungeons, raids
- **Guild Tools**: Guild search, info, and authenticated detail endpoints
- **Wizard's Vault**: Season info, objectives, and reward listings
- **Smart Caching**: Efficient caching with appropriate TTL for static and dynamic data
- **Extensible Architecture**: Modular design for easy feature additions

## Requirements

- Go 1.24 or higher (to build from source)
- Guild Wars 2 API key (for authenticated tools — set via `GW2_API_KEY` environment variable)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/AlyxPink/gw2-mcp.git
cd gw2-mcp
```

2. Install dependencies:
```bash
go mod tidy
```

3. Build the server:
```bash
make build
```

## Configuration

### API Key

Set the `GW2_API_KEY` environment variable to enable authenticated tools:

```bash
export GW2_API_KEY="YOUR_KEY_HERE"
```

Create an API key at [Guild Wars 2 API Key Management](https://account.arena.net/applications) with the scopes you need (account, wallet, tradingpost, characters, inventories, unlocks, progression, guilds).

If `GW2_API_KEY` is not set, the server still starts — unauthenticated tools work normally, and authenticated tools return a clear error message.

See [docs/configuration.md](docs/configuration.md) for the full configuration guide.

## Usage

### Running the Server

[![Add MCP Server gw2-mcp to LM Studio](https://files.lmstudio.ai/deeplink/mcp-install-light.svg#gh-light-mode-only)](https://lmstudio.ai/install-mcp?name=gw2-mcp&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJhbHl4cGluay9ndzItbWNwOnYxIl19#gh-light-mode-only)
[![Add MCP Server gw2-mcp to LM Studio](https://files.lmstudio.ai/deeplink/mcp-install-dark.svg#gh-dark-mode-only)](https://lmstudio.ai/install-mcp?name=gw2-mcp&config=eyJjb21tYW5kIjoiZG9ja2VyIiwiYXJncyI6WyJydW4iLCItLXJtIiwiLWkiLCJhbHl4cGluay9ndzItbWNwOnYxIl19#gh-dark-mode-only)

The MCP server communicates via stdio (standard input/output):

```bash
GW2_API_KEY="YOUR_KEY_HERE" ./bin/gw2-mcp
```

Configure Claude Desktop, Claude Code, LM Studio, or other MCP clients:

**Docker:**
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

**Direct binary:**
```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "/path/to/gw2-mcp",
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

### MCP Tools

The server provides 34 tools for LLM interaction:

#### Wiki

| Tool | Auth | Description |
|------|------|-------------|
| `wiki_search` | No | Search the Guild Wars 2 wiki |

#### Account

| Tool | Auth | Description |
|------|------|-------------|
| `get_account` | Yes | Account info (name, world, guilds, access) |
| `get_wallet` | Yes | Wallet currencies and balances |
| `get_bank` | Yes | Bank vault contents with item names |
| `get_materials` | Yes | Material storage contents with item names |
| `get_inventory` | Yes | Shared inventory slot contents |
| `get_characters` | Yes | List characters or get details for one (optional `name` param) |
| `get_account_unlocks` | Yes | Unlocked IDs by type (skins, dyes, minis, titles, recipes, etc.) |
| `get_account_progress` | Yes | Progress data by type (achievements, masteries, luck, etc.) |
| `get_account_dailies` | Yes | Completed dailies by type (crafting, dungeons, raids, etc.) |
| `get_token_info` | Yes | API key name and permission scopes |

#### Trading Post

| Tool | Auth | Description |
|------|------|-------------|
| `get_currencies` | No | Currency metadata (names, descriptions) |
| `get_tp_prices` | No | Item prices (best buy/sell with coin formatting) |
| `get_tp_listings` | No | Order book listings (all price tiers) |
| `get_gem_exchange` | No | Gem-to-coin or coin-to-gem exchange rates |
| `get_tp_delivery` | Yes | Delivery box contents (pending pickups) |
| `get_tp_transactions` | Yes | Transaction history (current/buys, current/sells, history/buys, history/sells) |

#### Game Data

| Tool | Auth | Description |
|------|------|-------------|
| `get_items` | No | Item metadata (name, type, rarity, level) |
| `get_skins` | No | Skin metadata (name, type, icon) |
| `get_recipes` | No | Recipe details (output, ingredients, disciplines) |
| `search_recipes` | No | Search recipes by input or output item ID |
| `get_achievements` | No | Achievement details (name, description, requirements) |
| `get_daily_achievements` | No | Today's and tomorrow's daily achievements |

#### Wizard's Vault

| Tool | Auth | Description |
|------|------|-------------|
| `get_wizards_vault` | No | Current season info |
| `get_wizards_vault_objectives` | Optional | Objectives (authenticated shows account progress) |
| `get_wizards_vault_listings` | Optional | Reward listings (authenticated shows purchase status) |

#### Guilds

| Tool | Auth | Description |
|------|------|-------------|
| `get_guild` | No | Public guild info (name, tag, level) |
| `search_guild` | No | Search for a guild by name |
| `get_guild_details` | Yes | Guild details (log, members, ranks, stash, storage, treasury, teams, upgrades) |

#### Game Metadata

| Tool | Auth | Description |
|------|------|-------------|
| `get_colors` | No | Dye color metadata |
| `get_minis` | No | Miniature metadata |
| `get_mounts_info` | No | Mount skin or type metadata |
| `get_game_build` | No | Current game build number |
| `get_dungeons_and_raids` | No | Dungeon or raid metadata |

### MCP Resources

The server provides the following resources:

#### Currency List (`gw2://currencies`)

Complete list of all Guild Wars 2 currencies with metadata.

## Caching Strategy

The server implements intelligent caching:

- **Static Data** (currencies, items, skins, recipes, achievements, colors, minis, mounts, dungeons): 1 day
- **Account Data** (wallet, bank, materials, inventory, characters, progress): 5 minutes
- **Trading Post** (prices, listings): 5 minutes; exchange rates: 10 minutes; delivery: 2 minutes
- **Unlocks**: 10 minutes
- **Dailies**: 2 minutes
- **Guild Info/Search**: 1 hour; Guild Details: 5 minutes
- **Wizard's Vault**: Season: 1 day; Objectives (auth): 5 min, (public): 1 hour; Listings: 1 hour
- **Game Build**: 1 hour; Token Info: 10 minutes
- **Wiki/Search Results**: 24 hours

## Architecture

The project follows Clean Architecture principles:

```
internal/
├── server/          # MCP server implementation
├── cache/           # Caching layer
├── gw2api/          # GW2 API client
└── wiki/            # Wiki API client
```

## Development

### Code Standards

- Format code with `gofumpt`
- Lint with `golangci-lint`
- Write unit tests for core functionality
- Follow conventional commit messages

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

### Formatting

```bash
gofumpt -w .
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Run linting and formatting
6. Submit a pull request

## License

GNU Affero General Public License v3.0 - see LICENSE file for details.

## Acknowledgments

- [Guild Wars 2 API](https://wiki.guildwars2.com/wiki/API:Main) for providing comprehensive game data
- [Guild Wars 2 Wiki](https://wiki.guildwars2.com/) for extensive game documentation
- [MCP Go](https://github.com/mark3labs/mcp-go) for the MCP implementation framework
