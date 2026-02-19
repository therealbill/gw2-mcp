---
title: Tools Reference
---

# Tools Reference

Complete specification for all 37 MCP tools exposed by the GW2 MCP Server. Each tool is invoked via the MCP `tools/call` method over stdio. For authentication requirements, see [API Scopes](../api-scopes/). For cache behavior, see [Caching](../caching/). For client setup, see [How to Configure MCP Clients](../../how-to/configure-mcp-clients/).

## Overview

### Wiki

| Tool | Auth | Description |
|------|------|-------------|
| [`wiki_search`](#wiki_search) | None | Search Guild Wars 2 wiki for information about game content |

### Account

| Tool | Auth | Description |
|------|------|-------------|
| [`get_account`](#get_account) | `GW2_API_KEY` | Get account information including name, world, guilds, and access |
| [`get_wallet`](#get_wallet) | `GW2_API_KEY` | Get wallet information including all currencies |
| [`get_bank`](#get_bank) | `GW2_API_KEY` | Get bank vault contents with item names |
| [`get_materials`](#get_materials) | `GW2_API_KEY` | Get material storage contents with item names |
| [`get_inventory`](#get_inventory) | `GW2_API_KEY` | Get shared inventory slot contents with item names |
| [`get_characters`](#get_characters) | `GW2_API_KEY` | List characters or get details for a specific character |
| [`get_account_unlocks`](#get_account_unlocks) | `GW2_API_KEY` | Get account unlocked IDs by type |
| [`get_account_progress`](#get_account_progress) | `GW2_API_KEY` | Get account progress data by type |
| [`get_account_dailies`](#get_account_dailies) | `GW2_API_KEY` | Get completed daily content IDs by type |
| [`get_token_info`](#get_token_info) | `GW2_API_KEY` | Get API key name and permission scopes |

### Trading Post

| Tool | Auth | Description |
|------|------|-------------|
| [`get_currencies`](#get_currencies) | None | Get information about Guild Wars 2 currencies |
| [`get_tp_prices`](#get_tp_prices) | None | Get aggregated best buy/sell prices for items |
| [`get_tp_listings`](#get_tp_listings) | None | Get full order book listings for items |
| [`get_gem_exchange`](#get_gem_exchange) | None | Get gem exchange rates between coins and gems |
| [`get_tp_delivery`](#get_tp_delivery) | `GW2_API_KEY` | Get items and coins awaiting pickup from the Trading Post |
| [`get_tp_transactions`](#get_tp_transactions) | `GW2_API_KEY` | Get current orders or completed transactions from the past 90 days |

### Game Data

| Tool | Auth | Description |
|------|------|-------------|
| [`get_items`](#get_items) | None | Get item metadata for given item IDs |
| [`get_skins`](#get_skins) | None | Get skin metadata for given skin IDs |
| [`get_recipes`](#get_recipes) | None | Get recipe details for given recipe IDs |
| [`search_recipes`](#search_recipes) | None | Search for recipe IDs by input or output item ID |
| [`get_achievements`](#get_achievements) | None | Get achievement details for given achievement IDs |
| [`get_daily_achievements`](#get_daily_achievements) | None | Get today's and tomorrow's daily achievements |

### Wizard's Vault

| Tool | Auth | Description |
|------|------|-------------|
| [`get_wizards_vault`](#get_wizards_vault) | None | Get current Wizard's Vault season information |
| [`get_wizards_vault_objectives`](#get_wizards_vault_objectives) | Optional | Get Wizard's Vault objectives; authenticated endpoint shows account progress |
| [`get_wizards_vault_listings`](#get_wizards_vault_listings) | Optional | Get Wizard's Vault reward listings; authenticated endpoint shows purchase status |

### Guilds

| Tool | Auth | Description |
|------|------|-------------|
| [`get_guild`](#get_guild) | None | Get public guild information (name, tag, level) |
| [`search_guild`](#search_guild) | None | Search for a guild by name |
| [`get_guild_details`](#get_guild_details) | `GW2_API_KEY` | Get detailed guild data (log, members, ranks, stash, etc.) |

### Game Metadata

| Tool | Auth | Description |
|------|------|-------------|
| [`get_colors`](#get_colors) | None | Get dye color metadata for given color IDs |
| [`get_minis`](#get_minis) | None | Get miniature metadata for given mini IDs |
| [`get_mounts_info`](#get_mounts_info) | None | Get mount skin or mount type metadata for given IDs |
| [`get_game_build`](#get_game_build) | None | Get the current Guild Wars 2 game build number |
| [`get_dungeons_and_raids`](#get_dungeons_and_raids) | None | Get dungeon or raid metadata for given IDs |

### Composite Tools

| Tool | Auth | Description |
|------|------|-------------|
| [`get_item_by_name`](#get_item_by_name) | None | Look up a GW2 item by name via wiki search, then return full item details |
| [`get_item_recipe_by_name`](#get_item_recipe_by_name) | None | Find crafting recipes for a GW2 item by name via wiki search |
| [`get_tp_price_by_name`](#get_tp_price_by_name) | None | Get Trading Post prices for an item by name via wiki search |

---

## Wiki

### wiki_search

Search Guild Wars 2 wiki for information about game content.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `query` | string | Yes | -- | Search query for wiki content (e.g. `"Dragon Bash"`, `"currencies"`, `"wallet"`) |
| `limit` | integer | No | `5` | Maximum number of results to return |

#### Example

```json
{
  "tool": "wiki_search",
  "arguments": {
    "query": "Dragon Bash",
    "limit": 3
  }
}
```

---

## Account

### get_account

Get account information including name, world, guilds, and access. Requires `GW2_API_KEY`.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_account",
  "arguments": {}
}
```

### get_wallet

Get wallet information including all currencies. Requires `GW2_API_KEY`.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_wallet",
  "arguments": {}
}
```

### get_bank

Get bank vault contents with item names. Requires `GW2_API_KEY`.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_bank",
  "arguments": {}
}
```

### get_materials

Get material storage contents with item names. Requires `GW2_API_KEY`.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_materials",
  "arguments": {}
}
```

### get_inventory

Get shared inventory slot contents with item names. Requires `GW2_API_KEY`.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_inventory",
  "arguments": {}
}
```

### get_characters

Get list of character names, or detailed info for a specific character including crafting disciplines, equipment, skills, specializations, and build tabs. Requires `GW2_API_KEY`.

When `name` is omitted, returns a list of all character names. When `name` is provided, returns detailed information for that character.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `name` | string | No | -- | Character name to get details for; omit to list all characters |

#### Examples

List all characters:

```json
{
  "tool": "get_characters",
  "arguments": {}
}
```

Get details for a specific character:

```json
{
  "tool": "get_characters",
  "arguments": {
    "name": "My Character Name"
  }
}
```

### get_account_unlocks

Get account unlocked IDs. Requires `GW2_API_KEY`.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `type` | string | Yes | -- | Unlock type |

Valid values for `type`:

| Value | Description |
|-------|-------------|
| `skins` | Unlocked skin IDs |
| `dyes` | Unlocked dye color IDs |
| `minis` | Unlocked miniature IDs |
| `titles` | Unlocked title IDs |
| `recipes` | Unlocked recipe IDs |
| `finishers` | Unlocked finisher IDs |
| `outfits` | Unlocked outfit IDs |
| `gliders` | Unlocked glider IDs |
| `mailcarriers` | Unlocked mail carrier IDs |
| `novelties` | Unlocked novelty IDs |
| `emotes` | Unlocked emote IDs |
| `mounts/skins` | Unlocked mount skin IDs |
| `mounts/types` | Unlocked mount type IDs |
| `skiffs` | Unlocked skiff skin IDs |
| `jadebots` | Unlocked jade bot skin IDs |

#### Example

```json
{
  "tool": "get_account_unlocks",
  "arguments": {
    "type": "skins"
  }
}
```

### get_account_progress

Get account progress data. Requires `GW2_API_KEY`.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `type` | string | Yes | -- | Progress type |

Valid values for `type`:

| Value | Description |
|-------|-------------|
| `achievements` | Account achievement progress |
| `masteries` | Mastery training progress |
| `mastery/points` | Mastery point totals |
| `luck` | Account luck value |
| `legendaryarmory` | Legendary armory contents |
| `progression` | General account progression |

#### Example

```json
{
  "tool": "get_account_progress",
  "arguments": {
    "type": "achievements"
  }
}
```

### get_account_dailies

Get completed daily content IDs. Requires `GW2_API_KEY`.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `type` | string | Yes | -- | Daily type |

Valid values for `type`:

| Value | Description |
|-------|-------------|
| `dailycrafting` | Completed daily crafting IDs |
| `dungeons` | Completed dungeon paths |
| `raids` | Completed raid encounters |
| `mapchests` | Completed map chest IDs |
| `worldbosses` | Completed world boss IDs |

#### Example

```json
{
  "tool": "get_account_dailies",
  "arguments": {
    "type": "raids"
  }
}
```

### get_token_info

Get API key information including name and permission scopes. Requires `GW2_API_KEY`.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_token_info",
  "arguments": {}
}
```

---

## Trading Post

### get_currencies

Get information about Guild Wars 2 currencies. Returns currency metadata including names and descriptions. When `ids` is omitted, returns all currencies.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ids` | array of integers | No | -- | Specific currency IDs to fetch; returns all if not specified |

#### Examples

Get all currencies:

```json
{
  "tool": "get_currencies",
  "arguments": {}
}
```

Get specific currencies:

```json
{
  "tool": "get_currencies",
  "arguments": {
    "ids": [1, 4]
  }
}
```

### get_tp_prices

Get Trading Post prices for items. Returns aggregated best buy/sell prices with item names and formatted coin values.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `item_ids` | array of integers | Yes | -- | Array of item IDs to get prices for (e.g. `[19976, 19721]` for Mystic Coin and Glob of Ectoplasm) |

#### Example

```json
{
  "tool": "get_tp_prices",
  "arguments": {
    "item_ids": [19976, 19721]
  }
}
```

### get_tp_listings

Get Trading Post order book listings for items. Returns all buy/sell price tiers with quantities.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `item_ids` | array of integers | Yes | -- | Array of item IDs to get listings for |

#### Example

```json
{
  "tool": "get_tp_listings",
  "arguments": {
    "item_ids": [19976]
  }
}
```

### get_gem_exchange

Get gem exchange rates. Convert coins to gems or gems to coins.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `direction` | string | Yes | -- | Exchange direction: `"coins"` (coins to gems) or `"gems"` (gems to coins) |
| `quantity` | integer | Yes | -- | Amount to convert (coins in copper, or number of gems); must be greater than 0 |

#### Example

```json
{
  "tool": "get_gem_exchange",
  "arguments": {
    "direction": "coins",
    "quantity": 10000000
  }
}
```

### get_tp_delivery

Get items and coins awaiting pickup from the Trading Post. Requires `GW2_API_KEY` with `account` and `tradingpost` scopes.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_tp_delivery",
  "arguments": {}
}
```

### get_tp_transactions

Get Trading Post transaction history. View current orders or completed transactions from the past 90 days. Requires `GW2_API_KEY` with `account` and `tradingpost` scopes.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `type` | string | Yes | -- | Transaction type |

Valid values for `type`:

| Value | Description |
|-------|-------------|
| `current/buys` | Current buy orders |
| `current/sells` | Current sell listings |
| `history/buys` | Completed buy transactions (past 90 days) |
| `history/sells` | Completed sell transactions (past 90 days) |

#### Example

```json
{
  "tool": "get_tp_transactions",
  "arguments": {
    "type": "history/sells"
  }
}
```

---

## Game Data

### get_items

Get item metadata (name, type, rarity, level, icon, description, vendor value, flags, game types, restrictions, and type-specific details) for given item IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ids` | array of integers | Yes | -- | Array of item IDs to look up |

#### Example

```json
{
  "tool": "get_items",
  "arguments": {
    "ids": [19976, 19721]
  }
}
```

### get_skins

Get skin metadata (name, type, icon, rarity, description, flags, restrictions, and type-specific details) for given skin IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ids` | array of integers | Yes | -- | Array of skin IDs to look up |

#### Example

```json
{
  "tool": "get_skins",
  "arguments": {
    "ids": [1, 2, 3]
  }
}
```

### get_recipes

Get recipe details (type, output, ingredients, disciplines, crafting time, flags, guild ingredients, chat link) for given recipe IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ids` | array of integers | Yes | -- | Array of recipe IDs to look up |

#### Example

```json
{
  "tool": "get_recipes",
  "arguments": {
    "ids": [7259]
  }
}
```

### search_recipes

Search for recipe IDs by input or output item ID. At least one of `input` or `output` is required.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `input` | integer | No | -- | Input item ID to search for recipes that use this item |
| `output` | integer | No | -- | Output item ID to search for recipes that produce this item |

#### Examples

Search by output item:

```json
{
  "tool": "search_recipes",
  "arguments": {
    "output": 19976
  }
}
```

Search by input item:

```json
{
  "tool": "search_recipes",
  "arguments": {
    "input": 19721
  }
}
```

### get_achievements

Get achievement details (name, description, requirements, tiers, prerequisites, rewards, bits, icon) for given achievement IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ids` | array of integers | Yes | -- | Array of achievement IDs to look up |

#### Example

```json
{
  "tool": "get_achievements",
  "arguments": {
    "ids": [1, 2, 3]
  }
}
```

### get_daily_achievements

Get today's and tomorrow's daily achievements.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_daily_achievements",
  "arguments": {}
}
```

---

## Wizard's Vault

### get_wizards_vault

Get current Wizard's Vault season information.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_wizards_vault",
  "arguments": {}
}
```

### get_wizards_vault_objectives

Get Wizard's Vault objectives. Uses authenticated endpoint if `GW2_API_KEY` is set (returns account-specific progress), otherwise returns the public objective list.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `type` | string | Yes | -- | Objective type |

Valid values for `type`:

| Value | Description |
|-------|-------------|
| `daily` | Daily objectives |
| `weekly` | Weekly objectives |
| `special` | Special objectives |

#### Example

```json
{
  "tool": "get_wizards_vault_objectives",
  "arguments": {
    "type": "daily"
  }
}
```

### get_wizards_vault_listings

Get Wizard's Vault reward listings. Uses authenticated endpoint if `GW2_API_KEY` is set (returns purchase status).

#### Parameters

None.

#### Example

```json
{
  "tool": "get_wizards_vault_listings",
  "arguments": {}
}
```

---

## Guilds

### get_guild

Get public guild information (name, tag, level).

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `id` | string | Yes | -- | Guild ID (UUID) |

#### Example

```json
{
  "tool": "get_guild",
  "arguments": {
    "id": "4BBB52AA-D768-4FC6-8EDE-C299F2822F0F"
  }
}
```

### search_guild

Search for a guild by name. Returns matching guild IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Guild name to search for |

#### Example

```json
{
  "tool": "search_guild",
  "arguments": {
    "name": "My Guild"
  }
}
```

### get_guild_details

Get detailed guild data (log, members, ranks, stash, etc.). Requires `GW2_API_KEY` with guild leader permissions.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `id` | string | Yes | -- | Guild ID (UUID) |
| `type` | string | Yes | -- | Detail type |

Valid values for `type`:

| Value | Description |
|-------|-------------|
| `log` | Guild activity log |
| `members` | Guild member list |
| `ranks` | Guild rank definitions |
| `stash` | Guild stash contents |
| `storage` | Guild storage contents |
| `treasury` | Guild treasury contents |
| `teams` | Guild PvP teams |
| `upgrades` | Guild upgrades |

#### Example

```json
{
  "tool": "get_guild_details",
  "arguments": {
    "id": "4BBB52AA-D768-4FC6-8EDE-C299F2822F0F",
    "type": "members"
  }
}
```

---

## Game Metadata

### get_colors

Get dye color metadata (name, base RGB, cloth/leather/metal/fur material adjustments) for given color IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ids` | array of integers | Yes | -- | Array of color IDs to look up |

#### Example

```json
{
  "tool": "get_colors",
  "arguments": {
    "ids": [1, 2, 3]
  }
}
```

### get_minis

Get miniature metadata (name, icon, item_id) for given mini IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `ids` | array of integers | Yes | -- | Array of mini IDs to look up |

#### Example

```json
{
  "tool": "get_minis",
  "arguments": {
    "ids": [1, 2, 3]
  }
}
```

### get_mounts_info

Get mount skin or mount type metadata for given IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `type` | string | Yes | -- | Mount info type: `"skins"` or `"types"` |
| `ids` | array of integers | Yes | -- | Array of mount skin or type IDs to look up |

#### Example

```json
{
  "tool": "get_mounts_info",
  "arguments": {
    "type": "skins",
    "ids": [1, 2, 3]
  }
}
```

### get_game_build

Get the current Guild Wars 2 game build number.

#### Parameters

None.

#### Example

```json
{
  "tool": "get_game_build",
  "arguments": {}
}
```

### get_dungeons_and_raids

Get dungeon or raid metadata (paths, wings, events) for given IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `type` | string | Yes | -- | Content type: `"dungeons"` or `"raids"` |
| `ids` | array of strings | Yes | -- | Array of dungeon or raid IDs (e.g. `"ascalonian_catacombs"`, `"forsaken_thicket"`) |

#### Example

```json
{
  "tool": "get_dungeons_and_raids",
  "arguments": {
    "type": "raids",
    "ids": ["forsaken_thicket"]
  }
}
```

---

## Composite Tools

Composite tools combine a wiki search with a GW2 API lookup in a single call. They search the wiki to resolve an item name to an ID, then fetch full data from the API.

### get_item_by_name

Look up a GW2 item by name. Searches the wiki to find the item ID, then returns full item details from the API.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Item name to search for (e.g. `"Mystic Coin"`, `"Dusk"`) |

#### Example

```json
{
  "tool": "get_item_by_name",
  "arguments": {
    "name": "Mystic Coin"
  }
}
```

### get_item_recipe_by_name

Find crafting recipes for a GW2 item by name. Searches the wiki for recipe data, then returns full recipe details with resolved ingredient names. Falls back to the API recipe search if the wiki result does not contain recipe IDs.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Item name to find recipes for (e.g. `"18 Slot Silk Bag"`, `"Dawn"`) |

#### Example

```json
{
  "tool": "get_item_recipe_by_name",
  "arguments": {
    "name": "Dawn"
  }
}
```

### get_tp_price_by_name

Get Trading Post prices for an item by name. Searches the wiki to find the item ID, then returns current buy/sell prices.

#### Parameters

| Name | Type | Required | Default | Description |
|------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Item name to get trading post prices for (e.g. `"Glob of Ectoplasm"`, `"Mystic Coin"`) |

#### Example

```json
{
  "tool": "get_tp_price_by_name",
  "arguments": {
    "name": "Glob of Ectoplasm"
  }
}
```
