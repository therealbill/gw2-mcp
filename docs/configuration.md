# GW2 MCP Server Configuration Guide

The GW2 MCP server bridges Large Language Models with Guild Wars 2 data. It provides 34 tools covering wiki search, account data, Trading Post, achievements, guilds, Wizard's Vault, and game metadata — all accessible through the [Model Context Protocol](https://modelcontextprotocol.io/).

## Prerequisites

### GW2 API Key

Many tools require a Guild Wars 2 API key, configured via the `GW2_API_KEY` environment variable. To create one:

1. Log in at [Guild Wars 2 API Key Management](https://account.arena.net/applications)
2. Click **New Key**
3. Select the following permissions:
   - **account** — Required for wallet, bank, characters, unlocks, progress, dailies
   - **wallet** — Required for currency information
   - **tradingpost** — Required for delivery box and transaction history
   - **characters** — Required for character details
   - **inventories** — Required for bank, material storage, and shared inventory
   - **unlocks** — Required for account unlocks (skins, dyes, minis, etc.)
   - **progression** — Required for achievements, masteries, and Wizard's Vault
   - **guilds** — Required for guild detail endpoints (log, members, stash, etc.)
4. Copy the generated key and keep it safe

## Installation

### Method 1: Docker (Recommended)

No build tools required. Pull and run directly:

```bash
docker run --rm -i -e GW2_API_KEY="YOUR_KEY_HERE" alyxpink/gw2-mcp:v1
```

The image is also available from GitHub Container Registry:

```bash
docker run --rm -i -e GW2_API_KEY="YOUR_KEY_HERE" ghcr.io/alyxpink/gw2-mcp:v1
```

You can omit `-e GW2_API_KEY` to run without authentication — unauthenticated tools (wiki search, item lookups, achievements, etc.) will still work.

### Method 2: Direct Binary

**Install with `go install`:**

```bash
go install github.com/AlyxPink/gw2-mcp@latest
```

**Build from source:**

```bash
git clone https://github.com/AlyxPink/gw2-mcp.git
cd gw2-mcp
make build
```

The binary is output to `bin/gw2-mcp` (or `bin/gw2-mcp.exe` on Windows).

**Download from GitHub Releases:**

Download a prebuilt binary from the [Releases page](https://github.com/AlyxPink/gw2-mcp/releases).

## Configuration

### API Key

The API key is configured via the `GW2_API_KEY` environment variable at server startup:

```bash
GW2_API_KEY="YOUR_KEY_HERE" ./gw2-mcp
```

If `GW2_API_KEY` is not set:
- The server starts normally
- Unauthenticated tools work as expected
- Authenticated tools return a clear error: `"GW2_API_KEY environment variable not configured"`
- A warning is logged at startup

## MCP Client Configuration

The server communicates over **stdio** (standard input/output). Configure your MCP client to launch it as a subprocess.

### Claude Desktop

Add to your Claude Desktop settings JSON (Settings > Developer > Edit Config):

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

### Claude Code

Add to your project `.mcp.json`:

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

### LM Studio

Click the badge on the [project README](https://github.com/AlyxPink/gw2-mcp) to install automatically, or add the server manually with command `docker` and args `run --rm -i -e GW2_API_KEY=YOUR_KEY_HERE alyxpink/gw2-mcp:v1`.

### Other MCP Clients (Cursor, Windsurf, etc.)

Any MCP client that supports stdio transports can use the same configuration pattern — set the command to `docker` with the args shown above, or point directly to the binary path. Pass `GW2_API_KEY` via the `env` block.

## Available Tools

### Tools That Require `GW2_API_KEY`

| Tool | Description |
|------|-------------|
| `get_wallet` | Get wallet currencies and balances |
| `get_account` | Get account info (name, world, guilds, access) |
| `get_bank` | Get bank vault contents with item names |
| `get_materials` | Get material storage contents with item names |
| `get_inventory` | Get shared inventory slot contents |
| `get_characters` | List characters or get details for a specific character |
| `get_account_unlocks` | Get unlocked IDs (skins, dyes, minis, titles, etc.) |
| `get_account_progress` | Get progress data (achievements, masteries, luck, etc.) |
| `get_account_dailies` | Get completed dailies (crafting, dungeons, raids, etc.) |
| `get_tp_delivery` | Get Trading Post delivery box contents |
| `get_tp_transactions` | Get Trading Post buy/sell history |
| `get_token_info` | Get API key name and permission scopes |
| `get_guild_details` | Get guild details (log, members, ranks, stash, etc.) |

### Tools That Optionally Use `GW2_API_KEY`

| Tool | Description |
|------|-------------|
| `get_wizards_vault_objectives` | Wizard's Vault objectives (authenticated shows account progress) |
| `get_wizards_vault_listings` | Wizard's Vault reward listings (authenticated shows purchase status) |

### Tools That Do NOT Require `GW2_API_KEY`

| Tool | Description |
|------|-------------|
| `wiki_search` | Search the Guild Wars 2 wiki |
| `get_currencies` | Get currency metadata (names, descriptions) |
| `get_tp_prices` | Get Trading Post item prices |
| `get_tp_listings` | Get Trading Post order book listings |
| `get_gem_exchange` | Get gem-to-coin or coin-to-gem exchange rates |
| `get_items` | Get item metadata (name, type, rarity, level) |
| `get_skins` | Get skin metadata (name, type, icon) |
| `get_recipes` | Get recipe details (output, ingredients, disciplines) |
| `search_recipes` | Search recipes by input or output item ID |
| `get_achievements` | Get achievement details (name, description, requirements) |
| `get_daily_achievements` | Get today's and tomorrow's daily achievements |
| `get_wizards_vault` | Get current Wizard's Vault season info |
| `get_guild` | Get public guild info (name, tag, level) |
| `search_guild` | Search for a guild by name |
| `get_colors` | Get dye color metadata |
| `get_minis` | Get miniature metadata |
| `get_mounts_info` | Get mount skin or type metadata |
| `get_game_build` | Get current game build number |
| `get_dungeons_and_raids` | Get dungeon or raid metadata |

### Example Tool Invocations

**Wiki search (no API key):**

```json
{
  "tool": "wiki_search",
  "arguments": {
    "query": "Dragon Bash",
    "limit": 3
  }
}
```

**Check wallet (requires GW2_API_KEY):**

```json
{
  "tool": "get_wallet",
  "arguments": {}
}
```

**Get character details (requires GW2_API_KEY):**

```json
{
  "tool": "get_characters",
  "arguments": {
    "name": "My Character Name"
  }
}
```

**Get account unlocks (requires GW2_API_KEY):**

```json
{
  "tool": "get_account_unlocks",
  "arguments": {
    "type": "skins"
  }
}
```

**Trading Post prices (no API key):**

```json
{
  "tool": "get_tp_prices",
  "arguments": {
    "item_ids": [19976, 19721]
  }
}
```

**Search recipes by output item (no API key):**

```json
{
  "tool": "search_recipes",
  "arguments": {
    "output": 19976
  }
}
```

**Trading Post transactions (requires GW2_API_KEY):**

```json
{
  "tool": "get_tp_transactions",
  "arguments": {
    "type": "history/sells"
  }
}
```

## API Key Scopes Reference

Each authenticated tool requires specific API key scopes:

| Tool | Required Scopes |
|------|----------------|
| `get_wallet` | `account`, `wallet` |
| `get_account` | `account` |
| `get_bank` | `account`, `inventories` |
| `get_materials` | `account`, `inventories` |
| `get_inventory` | `account`, `inventories` |
| `get_characters` | `account`, `characters` |
| `get_account_unlocks` | `account`, `unlocks` |
| `get_account_progress` | `account`, `progression` |
| `get_account_dailies` | `account`, `progression` |
| `get_tp_delivery` | `account`, `tradingpost` |
| `get_tp_transactions` | `account`, `tradingpost` |
| `get_token_info` | (any valid key) |
| `get_guild_details` | `account`, `guilds` |
| `get_wizards_vault_objectives` | `account`, `progression` |
| `get_wizards_vault_listings` | `account`, `progression` |

If your API key is missing a required scope, the GW2 API will return an authorization error.

## Security Notes

- **API key is configured once at startup.** Set `GW2_API_KEY` as an environment variable — it is never passed as a tool parameter.
- **API keys are hashed before caching.** The server uses SHA-256 to hash your API key and uses only the first 8 bytes of the hash as a cache key. Your raw API key is never stored in the cache.
- **No persistent storage.** All cached data is held in memory and lost when the server process exits.
- **stdio-only communication.** The server does not open any network ports or HTTP listeners. It communicates exclusively over standard input/output with the MCP client.
- **Never share your API key.** Treat it like a password. The server sends it to the GW2 API as a Bearer token over HTTPS.

## Troubleshooting

### GW2_API_KEY not configured

```
GW2_API_KEY environment variable not configured
```

Set the `GW2_API_KEY` environment variable before starting the server. If you only need unauthenticated tools, this warning can be ignored.

### Docker not running

```
error: Cannot connect to the Docker daemon
```

Make sure Docker Desktop (or the Docker daemon) is running before launching the MCP server.

### Binary not found

```
error: spawn gw2-mcp ENOENT
```

Ensure the binary path in your MCP client configuration is correct and the file is executable. If you used `go install`, the binary is in your `$GOPATH/bin` (or `$HOME/go/bin`).

### Invalid API key

```
API request failed with status 401
```

Double-check that you copied the full API key from [account.arena.net/applications](https://account.arena.net/applications). Keys are long alphanumeric strings with dashes.

### Missing API key scopes

```
API request failed with status 403
```

Your API key exists but lacks the required permissions. Create a new key with all needed scopes enabled.

### Trading Post transaction type invalid

```
invalid transaction type "...": must be one of current/buys, current/sells, history/buys, history/sells
```

The `type` parameter for `get_tp_transactions` must be exactly one of: `current/buys`, `current/sells`, `history/buys`, `history/sells`.

### Gem exchange direction invalid

```
invalid direction "...": must be "coins" or "gems"
```

The `direction` parameter for `get_gem_exchange` must be `coins` (to convert coins into gems) or `gems` (to convert gems into coins).
