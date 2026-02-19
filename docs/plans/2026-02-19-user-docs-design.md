# User-Facing Tutorials and How-To Guides

**Date**: 2026-02-19
**Status**: Approved

## Context

The existing documentation is developer-oriented (add tools, contribute, configure clients). GW2 players using AI assistants need docs that teach them how to use the server for common gameplay tasks. The composite tools (name-based item/recipe/price lookups) are the highest-value features and should be front and center.

## Audience

GW2 players using AI assistants (Claude Desktop, LM Studio, etc.). They know the game well but may not be technical. Docs should be conversation-driven — showing natural language prompts and describing expected responses.

## Design

### Tutorials (workflow-based, conversation-driven)

Each tutorial is a learn-by-doing walkthrough of a natural play session. Assumes Getting Started is complete.

#### 1. Trading Post Mastery (`tutorials/trading-post.md`)

**Goal**: Learn to use the AI as a Trading Post companion.

Walk through:
- Look up an item price by name ("What's Mystic Coin selling for?")
- Compare buy vs sell prices to understand the spread
- Check gem-to-gold exchange rates both directions
- Review your open buy/sell orders
- Check your TP delivery box for items to pick up
- Review recent transaction history

**Tools exercised**: `get_tp_price_by_name`, `get_tp_prices`, `get_gem_exchange`, `get_tp_transactions`, `get_tp_delivery`

#### 2. Know Your Account (`tutorials/account-overview.md`)

**Goal**: Use the AI to explore your GW2 account data.

Walk through:
- Check wallet balances (gold, karma, gems, etc.)
- Browse bank vault contents
- Check material storage
- List characters and inspect one in detail (equipment, crafting, build)
- Check what you've unlocked (skins, dyes, minis, etc.)

**Tools exercised**: `get_wallet`, `get_bank`, `get_materials`, `get_characters`, `get_account_unlocks`

#### 3. Crafting Assistant (`tutorials/crafting.md`)

**Goal**: Use the AI to plan crafting projects.

Walk through:
- Look up how to craft a specific item ("How do I craft Dawn?")
- Understand recipe output (ingredients, disciplines, quantities)
- Find what recipes use a particular material ("What can I craft with Glob of Ectoplasm?")
- Combine with TP prices to estimate crafting cost

**Tools exercised**: `get_item_recipe_by_name`, `get_item_by_name`, `search_recipes`, `get_tp_price_by_name`

#### 4. Daily Checklist (`tutorials/daily-checklist.md`)

**Goal**: Use the AI to track daily/weekly progress.

Walk through:
- Check today's Wizard's Vault daily and weekly objectives
- See which dailies you've already completed (authenticated)
- Check Wizard's Vault reward listings and what you can afford
- Review daily achievements
- Check raid/dungeon completion for the week

**Tools exercised**: `get_wizards_vault_objectives`, `get_wizards_vault_listings`, `get_daily_achievements`, `get_account_dailies`

### How-To Guides (task-specific, tool-category oriented)

Short, goal-oriented recipes. Each solves a specific problem.

#### Trading Post
- **Check if an item is profitable to flip** (`how-to/check-item-profitability.md`) — Look up buy/sell prices, calculate the 15% listing fee, determine profit margin
- **Track your open orders** (`how-to/track-tp-orders.md`) — View current buy/sell orders, check filled transactions
- **Convert between gems and gold** (`how-to/gem-exchange.md`) — Check rates both directions, understand the quantity parameter

#### Account & Inventory
- **Find valuable items in your bank** (`how-to/find-bank-valuables.md`) — Pull bank contents, cross-reference with TP prices
- **Audit your API key permissions** (`how-to/audit-api-key.md`) — Use `get_token_info` to check scopes, understand tool requirements
- **Compare characters** (`how-to/compare-characters.md`) — List characters, inspect gear and build tabs

#### Crafting
- **Estimate crafting cost vs buying** (`how-to/crafting-vs-buying.md`) — Look up recipe, get ingredient prices, compare to TP price of finished item
- **Find all recipes that use a material** (`how-to/find-recipes-for-material.md`) — Use `search_recipes` with input item

#### Daily Play
- **Plan your daily Wizard's Vault session** (`how-to/wizards-vault-daily.md`) — Check objectives, progress, available rewards
- **Track weekly raid and dungeon clears** (`how-to/track-raid-clears.md`) — Check completion for current reset

#### Guild
- **Look up guild information** (`how-to/guild-lookup.md`) — Search by name, view members/stash/treasury

#### General
- **Use the server without an API key** (`how-to/no-api-key.md`) — What works unauthenticated: wiki, items, TP prices, metadata
- **Look up anything on the GW2 wiki** (`how-to/wiki-search.md`) — Search for mechanics, events, NPCs, lore

### Writing style

- Conversation-driven: show prompts as "Ask your AI:" followed by natural language
- Describe expected responses in plain terms, not raw JSON
- Include GW2 context (what the data means in-game)
- Each tutorial has checkpoints to confirm progress
- How-to guides are concise — goal, prerequisites, steps, done

### File structure

```
docs/
  tutorials/
    _index.md           (update: add new entries)
    getting-started.md  (existing)
    trading-post.md
    account-overview.md
    crafting.md
    daily-checklist.md
  how-to/
    _index.md           (update: add new entries)
    configure-mcp-clients.md  (existing)
    add-a-new-tool.md         (existing)
    contribute.md             (existing)
    check-item-profitability.md
    track-tp-orders.md
    gem-exchange.md
    find-bank-valuables.md
    audit-api-key.md
    compare-characters.md
    crafting-vs-buying.md
    find-recipes-for-material.md
    wizards-vault-daily.md
    track-raid-clears.md
    guild-lookup.md
    no-api-key.md
    wiki-search.md
```

### Summary

- 4 new tutorials (workflow-based, conversation-driven)
- 13 new how-to guides (task-specific, tool-category organized)
- Update `tutorials/_index.md` and `how-to/_index.md` with new entries
- All 37 tools get coverage across the new docs
