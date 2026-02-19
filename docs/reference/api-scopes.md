---
title: API Key Scopes
---

# API Key Scopes

GW2 MCP Server tools that access account-specific data require a Guild Wars 2 API key with specific permission scopes. This page lists the scopes required by each authenticated tool and describes each available scope.

## Tool Scope Requirements

Every authenticated tool requires the `account` scope plus zero or more additional scopes. Tools not listed here do not require an API key.

| Tool | Required Scopes |
|------|-----------------|
| `get_account` | `account` |
| `get_account_dailies` | `account`, `progression` |
| `get_account_progress` | `account`, `progression` |
| `get_account_unlocks` | `account`, `unlocks` |
| `get_bank` | `account`, `inventories` |
| `get_characters` | `account`, `characters` |
| `get_guild_details` | `account`, `guilds` |
| `get_inventory` | `account`, `inventories` |
| `get_materials` | `account`, `inventories` |
| `get_token_info` | any valid key |
| `get_tp_delivery` | `account`, `tradingpost` |
| `get_tp_transactions` | `account`, `tradingpost` |
| `get_wallet` | `account`, `wallet` |
| `get_wizards_vault_listings` | `account`, `progression` |
| `get_wizards_vault_objectives` | `account`, `progression` |

If the API key is missing a required scope, the GW2 API returns an authorization error.

## Available Scopes

| Scope | Description |
|-------|-------------|
| `account` | Basic account information. Required by all authenticated tools. |
| `characters` | Character names, equipment, builds, and crafting disciplines. |
| `guilds` | Guild detail endpoints: log, members, ranks, stash, storage, treasury, teams, upgrades. Requires guild leader permissions on the key's account. |
| `inventories` | Bank vault, material storage, and shared inventory slots. |
| `progression` | Achievements, masteries, mastery points, luck, daily completions, and Wizard's Vault objectives and listings. |
| `tradingpost` | Trading Post delivery box and transaction history (current and past 90 days). |
| `unlocks` | Account unlocks: skins, dyes, minis, titles, recipes, finishers, outfits, gliders, mail carriers, novelties, emotes, mount skins, mount types, skiffs, and jade bots. |
| `wallet` | Currency balances (gold, karma, gems, etc.). |

## Key Management

API keys are created and managed at [Guild Wars 2 API Key Management](https://account.arena.net/applications).

The API key is configured via the `GW2_API_KEY` environment variable. See [Configuration](configuration/) for details on setting environment variables.

## See Also

- [Tools](tools/) -- complete reference for all MCP tools
- [Configuration](configuration/) -- environment variable and startup reference
- [How to Configure MCP Clients](../how-to/configure-mcp-clients/) -- setup instructions for Claude Desktop, VS Code, and other clients
