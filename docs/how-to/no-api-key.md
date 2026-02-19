---
title: Use Without an API Key
weight: 21
---

# How to Use the Server Without an API Key

**Goal**: Understand which tools work without an API key, so you can try the server before setting up authentication.

**Time**: Approximately 5 minutes to read; ongoing use

## Prerequisites

None. That is the point of this guide.

You need the GW2 MCP Server connected to an AI client (Claude Desktop, Claude Code, LM Studio, etc.), but no API key configured. See [Getting Started](../tutorials/getting-started/) for installation steps -- you can skip Step 2 (creating an API key) and leave `GW2_API_KEY` empty or remove the `env` block entirely.

## Steps

### 1. Know what works without a key

A large portion of the server's tools query public GW2 data that requires no authentication. These are great for trying the server out or for players who prefer not to share their API key.

### 2. Browse available tools by category

**Item and Price Lookups**

These tools let you research items, prices, and recipes -- the core of Trading Post gameplay:

- Look up any item by name -- `get_item_by_name`
- Check Trading Post buy/sell prices -- `get_tp_price_by_name`, `get_tp_prices`, `get_tp_listings`
- Look up crafting recipes by item name -- `get_item_recipe_by_name`
- Search recipes by input or output ID -- `get_recipes`, `search_recipes`
- Check the gem-to-gold and gold-to-gem exchange rates -- `get_gem_exchange`

**Game Information**

General game data and metadata available to everyone:

- Search the official GW2 wiki -- `wiki_search`
- Browse item metadata by ID -- `get_items`
- Browse skins, dye colors, miniatures, and mounts -- `get_skins`, `get_colors`, `get_minis`, `get_mounts_info`
- Check daily achievement rotations -- `get_daily_achievements`
- Browse currency definitions -- `get_currencies`
- Check dungeon and raid metadata -- `get_dungeons_and_raids`
- Check the current game build number -- `get_game_build`

**Wizard's Vault (public)**

The Wizard's Vault seasonal system exposes some data publicly:

- Current season info -- `get_wizards_vault`
- Objective lists (without personal progress) -- `get_wizards_vault_objectives`
- Reward listings (without purchase status) -- `get_wizards_vault_listings`

**Guilds (public)**

Basic guild lookups work without a key:

- Search for guilds by name -- `search_guild`
- View public guild info (name, tag, level) -- `get_guild`

### 3. Try some example prompts

These prompts all work without an API key:

> Ask your AI: "What's Mystic Coin selling for on the Trading Post?"

Your assistant fetches live buy and sell prices for Mystic Coin. Useful for checking market prices before logging in.

> Ask your AI: "How do I craft Dawn?"

Your assistant looks up the full recipe for the precursor greatsword Dawn, including all ingredient names and quantities.

> Ask your AI: "Search the wiki for Dragon Bash"

Your assistant searches the official GW2 wiki and returns page titles, snippets, and links to the full articles.

> Ask your AI: "How much gold for 400 gems?"

Your assistant checks the current gem exchange rate and tells you the gold cost to convert 400 gems.

### 4. Understand what an API key unlocks

Adding a GW2 API key gives access to your personal account data. Here is what becomes available:

| Category | What you get | Required scopes |
|----------|-------------|-----------------|
| Wallet and currencies | Your gold, karma, gems, and all other currency balances | `wallet` |
| Bank and materials | Contents of your bank vault and material storage | `inventories` |
| Characters | Character list, equipment, builds, crafting levels | `characters` |
| Trading Post history | Your active buy/sell orders and transaction history | `tradingpost` |
| Unlocks | Skins, dyes, minis, titles, and recipes you have unlocked | `unlocks` |
| Achievement progress | Your progress on achievements and dailies | `progression` |
| Wizard's Vault progress | Personal objective completion and astral acclaim balance | `progression` |
| Guild details | Members, stash, treasury, logs (guild leader only) | `guilds` |
| Account info | Account name, world, access level, guild memberships | `account` |

When you are ready to set up a key, follow the [Getting Started](../tutorials/getting-started/) tutorial from Step 2.

## Verify it works

Test that keyless access is working:

> Ask your AI: "What is the current GW2 game build number?"

Your assistant calls `get_game_build` and returns a numeric build ID. This is the simplest public tool -- if it works, the server is connected and ready.

Then try a price lookup:

> Ask your AI: "What's Glob of Ectoplasm selling for?"

If you see buy and sell prices, the public tools are working correctly.

## Troubleshooting

### Problem: "GW2_API_KEY environment variable not configured"
**Symptom**: You get this error when trying a tool.
**Cause**: The tool you called requires an API key. Not all tools are keyless.
**Solution**: Stick to the tools listed in this guide. If you want access to personal account data, set up an API key by following the [Getting Started](../tutorials/getting-started/) tutorial.

### Problem: Tools are not loading at all
**Symptom**: Your AI client does not show any GW2 tools.
**Cause**: The MCP server is not connected to your AI client.
**Solution**: Verify the server is configured in your client. See [Configure MCP Clients](configure-mcp-clients/) for setup instructions. The server itself starts without a key -- the key is only needed for authenticated tools.

### Problem: Item not found
**Symptom**: The AI says it cannot find an item you asked about.
**Cause**: The wiki search could not match the name you used.
**Solution**: Use the full in-game item name. Common abbreviations (like "Ecto" for Glob of Ectoplasm or "MC" for Mystic Coin) may not resolve correctly.

## See also

- [Getting Started](../tutorials/getting-started/) -- full tutorial including API key setup
- [Configure MCP Clients](configure-mcp-clients/) -- connect Claude Desktop, Claude Code, LM Studio, and other clients
- [API Scopes reference](../reference/api-scopes/) -- which permissions each tool requires
