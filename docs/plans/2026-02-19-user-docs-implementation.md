# User-Facing Docs Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create 4 tutorials and 13 how-to guides for GW2 players using AI assistants.

**Architecture:** Each doc is a standalone markdown file in `docs/tutorials/` or `docs/how-to/`. Tutorials are conversation-driven walkthroughs (show natural language prompts and describe responses). How-to guides are short task-focused recipes. All files use YAML frontmatter with `title` and `weight` for sidebar ordering.

**Tech Stack:** Markdown, Hugo Book theme (frontmatter: `title`, `weight`)

**Design doc:** `docs/plans/2026-02-19-user-docs-design.md`

**Reference for tool details:** `docs/reference/tools.md` — contains all 37 tools with parameters, examples, and descriptions. Use this as the source of truth for tool names, parameter formats, and behavior.

**Writing style:**
- Prompts shown as: `> Ask your AI: "natural language question"`
- Responses described in plain terms, not raw JSON
- Include GW2 context (what the data means in-game)
- Tutorials have checkpoints; how-to guides are goal → steps → done
- Keep cross-references as relative markdown links (e.g., `../reference/tools/`)

---

### Task 1: Tutorial — Trading Post Mastery

**Agent:** `diataxis-docs:doc-tutorial-writer`

**Files:**
- Create: `docs/tutorials/trading-post.md`

**Step 1: Write the tutorial**

Create `docs/tutorials/trading-post.md` with:
- Frontmatter: `title: Trading Post Mastery`, `weight: 2`
- Prerequisites: Getting Started complete, API key with `tradingpost` scope
- Sections walking through:
  1. **Look up item prices by name** — "What's the TP price for Mystic Coin?" → explains buy/sell prices, spread. Uses `get_tp_price_by_name`.
  2. **Compare multiple items** — "Compare prices for Mystic Coin and Glob of Ectoplasm" → shows how AI can use the tool multiple times.
  3. **Check gem exchange rates** — "How much gold for 400 gems?" and "How many gems can I get for 100 gold?" → explains both directions of `get_gem_exchange`, notes quantity is in copper for coins.
  4. **Review your open orders** — "Show my current buy orders" and "Show my sell listings" → `get_tp_transactions` with `current/buys` and `current/sells`.
  5. **Check your delivery box** — "What's waiting at the Trading Post?" → `get_tp_delivery`.
  6. **Transaction history** — "Show my recent sells" → `get_tp_transactions` with `history/sells`.
- Checkpoint after each section
- "What you learned" summary
- Next steps linking to: check-item-profitability how-to, gem-exchange how-to

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`
Expected: Clean build, no errors, page count increases

**Step 3: Commit**

```bash
git add docs/tutorials/trading-post.md
git commit -m "docs: add Trading Post Mastery tutorial"
```

---

### Task 2: Tutorial — Know Your Account

**Agent:** `diataxis-docs:doc-tutorial-writer`

**Files:**
- Create: `docs/tutorials/account-overview.md`

**Step 1: Write the tutorial**

Create `docs/tutorials/account-overview.md` with:
- Frontmatter: `title: Know Your Account`, `weight: 3`
- Prerequisites: Getting Started complete, API key with `account`, `characters`, `inventories`, `wallet`, `unlocks` scopes
- Sections:
  1. **Check your wallet** — "What's in my wallet?" → `get_wallet`, explain gold/karma/gems/laurels etc.
  2. **Browse your bank** — "What's in my bank?" → `get_bank`, note items are returned with names.
  3. **Material storage** — "What materials do I have?" → `get_materials`.
  4. **Shared inventory** — "Show my shared inventory slots" → `get_inventory`.
  5. **List your characters** — "List my characters" → `get_characters` without name.
  6. **Inspect a character** — "Show me details for [character name]" → `get_characters` with name, explains equipment, crafting disciplines, build tabs.
  7. **Check unlocks** — "What skins have I unlocked?" → `get_account_unlocks` with type `skins`. Mention other types: dyes, minis, titles, recipes, mounts.
  8. **Account overview** — "Tell me about my account" → `get_account`, shows account name, world, access level, guilds.
- Checkpoints throughout
- Next steps linking to: compare-characters how-to, find-bank-valuables how-to

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`
Expected: Clean build

**Step 3: Commit**

```bash
git add docs/tutorials/account-overview.md
git commit -m "docs: add Know Your Account tutorial"
```

---

### Task 3: Tutorial — Crafting Assistant

**Agent:** `diataxis-docs:doc-tutorial-writer`

**Files:**
- Create: `docs/tutorials/crafting.md`

**Step 1: Write the tutorial**

Create `docs/tutorials/crafting.md` with:
- Frontmatter: `title: Crafting Assistant`, `weight: 4`
- Prerequisites: Getting Started complete (no API key required for crafting lookups)
- Sections:
  1. **Look up a recipe by name** — "How do I craft Dawn?" → `get_item_recipe_by_name`, explain the response: output item, ingredients with names and quantities, crafting discipline.
  2. **Look up item details** — "Tell me about Dusk" → `get_item_by_name`, explain item metadata: rarity, type, level, description, vendor value.
  3. **Find recipes that use a material** — "What can I craft with Glob of Ectoplasm?" → explain that this requires the item ID first, then `search_recipes` with input. Show the two-step flow: ask AI to find recipes using ectoplasm, it calls `get_item_by_name` to get the ID then `search_recipes`.
  4. **Estimate crafting cost** — "How much would it cost to craft [item] on the Trading Post?" → combines `get_item_recipe_by_name` + `get_tp_price_by_name` for ingredients. Show how AI chains these naturally.
  5. **Compare crafting vs buying** — "Is it cheaper to craft or buy [item]?" → AI combines recipe cost estimate with TP price of finished item.
- Checkpoint after each section
- Note: composite tools use wiki search internally, so name matching is fuzzy — exact names work best
- Next steps: crafting-vs-buying how-to, find-recipes-for-material how-to

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`
Expected: Clean build

**Step 3: Commit**

```bash
git add docs/tutorials/crafting.md
git commit -m "docs: add Crafting Assistant tutorial"
```

---

### Task 4: Tutorial — Daily Checklist

**Agent:** `diataxis-docs:doc-tutorial-writer`

**Files:**
- Create: `docs/tutorials/daily-checklist.md`

**Step 1: Write the tutorial**

Create `docs/tutorials/daily-checklist.md` with:
- Frontmatter: `title: Daily Checklist`, `weight: 5`
- Prerequisites: Getting Started complete, API key with `account` and `progression` scopes (for personal progress tracking)
- Sections:
  1. **Wizard's Vault dailies** — "What are today's Wizard's Vault objectives?" → `get_wizards_vault_objectives` with type `daily`. Explain: with API key, shows completion progress; without, shows objectives only.
  2. **Wizard's Vault weeklies** — "What are this week's Wizard's Vault objectives?" → type `weekly`.
  3. **Wizard's Vault rewards** — "What can I buy from the Wizard's Vault?" → `get_wizards_vault_listings`. With API key shows what you've already purchased.
  4. **Wizard's Vault season** — "What's the current Wizard's Vault season?" → `get_wizards_vault`.
  5. **Daily achievements** — "What are today's daily achievements?" → `get_daily_achievements`.
  6. **Raid completion** — "Which raids have I cleared this week?" → `get_account_dailies` with type `raids`.
  7. **Dungeon completion** — "Which dungeon paths have I done today?" → `get_account_dailies` with type `dungeons`.
  8. **World bosses** — "Which world bosses have I done today?" → `get_account_dailies` with type `worldbosses`.
- Checkpoints
- Next steps: wizards-vault-daily how-to, track-raid-clears how-to

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`
Expected: Clean build

**Step 3: Commit**

```bash
git add docs/tutorials/daily-checklist.md
git commit -m "docs: add Daily Checklist tutorial"
```

---

### Task 5: How-To Guides — Trading Post (3 guides)

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/check-item-profitability.md`
- Create: `docs/how-to/track-tp-orders.md`
- Create: `docs/how-to/gem-exchange.md`

**Step 1: Write the three guides**

Each guide follows the pattern: title, goal line, prerequisites, numbered steps, see-also links.

**`check-item-profitability.md`** (`weight: 10`):
- Goal: Determine if an item is worth flipping on the Trading Post
- Steps: Look up buy/sell prices → calculate listing fee (15% of sell price) → subtract buy price → determine profit
- Prompt: "Is Mystic Coin profitable to flip?" — explain that AI calls `get_tp_price_by_name` and can calculate the margin
- Note the 15% TP tax (5% listing fee + 10% exchange fee)

**`track-tp-orders.md`** (`weight: 11`):
- Goal: View your current and recent Trading Post activity
- Steps: Check current buy orders → check current sell listings → review recent history
- Prompts: "Show my current buy orders", "Show my recent completed sells"
- Requires `tradingpost` API scope

**`gem-exchange.md`** (`weight: 12`):
- Goal: Check gem-to-gold and gold-to-gem conversion rates
- Steps: coins→gems query → gems→coins query
- Note: for coins direction, quantity is in copper (1 gold = 10000 copper, so 100 gold = 1000000)
- Note: for gems direction, quantity is number of gems

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`
Expected: Clean build

**Step 3: Commit**

```bash
git add docs/how-to/check-item-profitability.md docs/how-to/track-tp-orders.md docs/how-to/gem-exchange.md
git commit -m "docs: add Trading Post how-to guides"
```

---

### Task 6: How-To Guides — Account & Inventory (3 guides)

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/find-bank-valuables.md`
- Create: `docs/how-to/audit-api-key.md`
- Create: `docs/how-to/compare-characters.md`

**Step 1: Write the three guides**

**`find-bank-valuables.md`** (`weight: 13`):
- Goal: Find items in your bank that are worth gold on the Trading Post
- Steps: Ask AI to check bank → ask it to look up prices for the items → identify the most valuable
- Prompt: "What's the most valuable stuff in my bank?"
- The AI can chain `get_bank` → `get_tp_prices` for the item IDs

**`audit-api-key.md`** (`weight: 14`):
- Goal: Check what permissions your API key has and understand what tools need which scopes
- Steps: Ask "What permissions does my API key have?" → `get_token_info`
- Include table mapping scopes to tool categories
- Link to API Scopes reference for full details

**`compare-characters.md`** (`weight: 15`):
- Goal: View and compare your characters' details
- Steps: List characters → inspect specific ones → compare gear/builds
- Prompt: "List my characters", "Show me details for [name]"

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`

**Step 3: Commit**

```bash
git add docs/how-to/find-bank-valuables.md docs/how-to/audit-api-key.md docs/how-to/compare-characters.md
git commit -m "docs: add Account & Inventory how-to guides"
```

---

### Task 7: How-To Guides — Crafting (2 guides)

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/crafting-vs-buying.md`
- Create: `docs/how-to/find-recipes-for-material.md`

**Step 1: Write the two guides**

**`crafting-vs-buying.md`** (`weight: 16`):
- Goal: Determine whether it's cheaper to craft an item or buy it on the TP
- Steps: Look up recipe → get ingredient prices → sum total → compare to finished item TP price
- Prompt: "Is it cheaper to craft or buy [item]?"
- Note: AI naturally chains `get_item_recipe_by_name` + `get_tp_price_by_name`

**`find-recipes-for-material.md`** (`weight: 17`):
- Goal: Find out what you can craft with a specific material
- Steps: Get item ID (via `get_item_by_name`) → search recipes by input → get recipe details
- Prompt: "What can I craft with Glob of Ectoplasm?"
- Note: `search_recipes` takes item IDs not names, so the AI uses the composite tool first

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`

**Step 3: Commit**

```bash
git add docs/how-to/crafting-vs-buying.md docs/how-to/find-recipes-for-material.md
git commit -m "docs: add Crafting how-to guides"
```

---

### Task 8: How-To Guides — Daily Play (2 guides)

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/wizards-vault-daily.md`
- Create: `docs/how-to/track-raid-clears.md`

**Step 1: Write the two guides**

**`wizards-vault-daily.md`** (`weight: 18`):
- Goal: Plan your daily Wizard's Vault session
- Steps: Check daily objectives → check weekly objectives → see available rewards → check what you can afford
- Note: authenticated vs unauthenticated differences

**`track-raid-clears.md`** (`weight: 19`):
- Goal: See which raids and dungeons you've completed this reset
- Steps: Check raid clears → check dungeon paths → check world bosses
- Uses `get_account_dailies` with types `raids`, `dungeons`, `worldbosses`
- Note: resets weekly (raids) and daily (dungeons, world bosses)

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`

**Step 3: Commit**

```bash
git add docs/how-to/wizards-vault-daily.md docs/how-to/track-raid-clears.md
git commit -m "docs: add Daily Play how-to guides"
```

---

### Task 9: How-To Guides — Guild + General (3 guides)

**Agent:** `diataxis-docs:doc-howto-writer`

**Files:**
- Create: `docs/how-to/guild-lookup.md`
- Create: `docs/how-to/no-api-key.md`
- Create: `docs/how-to/wiki-search.md`

**Step 1: Write the three guides**

**`guild-lookup.md`** (`weight: 20`):
- Goal: Find and view guild information
- Steps: Search by name → get public info → get detailed data (if guild leader)
- Tools: `search_guild`, `get_guild`, `get_guild_details`
- Note: detailed data (members, stash, treasury) requires guild leader API key permissions

**`no-api-key.md`** (`weight: 21`):
- Goal: Understand what the server can do without authentication
- List all tools that work without a key: `wiki_search`, `get_item_by_name`, `get_item_recipe_by_name`, `get_tp_price_by_name`, `get_tp_prices`, `get_tp_listings`, `get_gem_exchange`, `get_items`, `get_skins`, `get_recipes`, `search_recipes`, `get_achievements`, `get_daily_achievements`, `get_currencies`, `get_wizards_vault`, `get_wizards_vault_objectives` (public), `get_wizards_vault_listings` (public), `get_guild`, `search_guild`, `get_colors`, `get_minis`, `get_mounts_info`, `get_game_build`, `get_dungeons_and_raids`
- Show a few example prompts that work without a key
- Explain what adding a key unlocks (account data, personal progress)

**`wiki-search.md`** (`weight: 22`):
- Goal: Search the GW2 wiki for game information
- Steps: Simple search → understand results (title, snippet, URL, infobox data)
- Prompt: "Search the wiki for Dragon Bash"
- Note: wiki search is also used internally by composite tools
- Examples: mechanics, events, NPCs, lore, game systems

**Step 2: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`

**Step 3: Commit**

```bash
git add docs/how-to/guild-lookup.md docs/how-to/no-api-key.md docs/how-to/wiki-search.md
git commit -m "docs: add Guild and General how-to guides"
```

---

### Task 10: Update Section Index Pages

**Agent:** `diataxis-docs:doc-crosslink-validator` (for validation after edits)

**Files:**
- Modify: `docs/tutorials/_index.md`
- Modify: `docs/how-to/_index.md`

**Step 1: Update tutorials index**

Add the four new tutorials to the list in `docs/tutorials/_index.md`:

```markdown
- [Getting Started](getting-started/) — Install, configure, and make your first query
- [Trading Post Mastery](trading-post/) — Use your AI as a Trading Post companion
- [Know Your Account](account-overview/) — Explore your wallet, bank, characters, and unlocks
- [Crafting Assistant](crafting/) — Plan crafting projects with recipe lookups and cost estimates
- [Daily Checklist](daily-checklist/) — Track Wizard's Vault, achievements, and raid clears
```

**Step 2: Update how-to index**

Reorganize `docs/how-to/_index.md` with grouped sections:

```markdown
### For Players

**Trading Post**
- [Check Item Profitability](check-item-profitability/) — Calculate flip margins including the 15% TP tax
- [Track Your TP Orders](track-tp-orders/) — View current orders and recent transaction history
- [Gem Exchange Rates](gem-exchange/) — Convert between gems and gold

**Account & Inventory**
- [Find Valuable Items in Your Bank](find-bank-valuables/) — Cross-reference bank contents with TP prices
- [Audit Your API Key](audit-api-key/) — Check what permissions your key has
- [Compare Characters](compare-characters/) — Inspect gear, builds, and crafting disciplines

**Crafting**
- [Crafting vs Buying](crafting-vs-buying/) — Estimate whether crafting or buying is cheaper
- [Find Recipes for a Material](find-recipes-for-material/) — Discover what you can craft with an item

**Daily Play**
- [Wizard's Vault Daily Plan](wizards-vault-daily/) — Check objectives, progress, and rewards
- [Track Raid and Dungeon Clears](track-raid-clears/) — See what you've completed this reset

**Guild**
- [Guild Lookup](guild-lookup/) — Search for guilds and view member data

**General**
- [Use Without an API Key](no-api-key/) — Everything that works without authentication
- [Wiki Search](wiki-search/) — Search the GW2 wiki for game information

### For Developers

- [Configure MCP Clients](configure-mcp-clients/) — Set up Claude Desktop, Claude Code, LM Studio, and other clients
- [Add a New Tool](add-a-new-tool/) — Developer guide for extending the server with new tools
- [Contribute](contribute/) — Fork, develop, test, and submit pull requests
```

**Step 3: Verify**

Run: `cd /c/Users/ucntc/repos/gw2-mcp && bin/hugo.exe --buildDrafts 2>&1`
Expected: Clean build, final page count should be ~43 (26 existing + 17 new)

**Step 4: Full link validation**

Crawl all new pages with curl to verify they return 200 and internal links resolve.

**Step 5: Commit**

```bash
git add docs/tutorials/_index.md docs/how-to/_index.md
git commit -m "docs: update section indexes with new tutorials and how-to guides"
```
