---
title: Wiki Search
weight: 22
---

# How to Search the GW2 Wiki

**Goal**: Search the official Guild Wars 2 wiki for game information directly through your AI assistant.

**Time**: Approximately 3 minutes

## Prerequisites

None. Wiki search is a public tool that works without an API key.

## Steps

### 1. Ask a wiki search question

> Ask your AI: "Search the wiki for Dragon Bash"

Your assistant calls the `wiki_search` tool with your query. The search runs against the official Guild Wars 2 wiki at wiki.guildwars2.com.

### 2. Read the results

For each matching page, your assistant shows:

- **Page title** -- the wiki article name
- **Snippet** -- a short excerpt showing where your search terms appear
- **URL** -- a direct link to the full wiki page

Some results also include **infobox data** -- structured information pulled from the sidebar of the wiki page. For items, this might include rarity, type, or acquisition method. For events, it might include location and frequency.

### 3. Browse the full article

If a result looks relevant, follow the URL to read the full wiki page in your browser. The wiki has far more detail than the search snippet can show.

### 4. Search for different types of content

Wiki search covers everything on the GW2 wiki. Here are examples across different categories:

**Game mechanics:**

> Ask your AI: "Search the wiki for condition damage"

Returns pages explaining how condition damage works, including the damage formula and relevant stats.

**Festivals and events:**

> Ask your AI: "Search the wiki for Dragon Bash"

Returns the Dragon Bash festival page with event schedules, achievements, and rewards.

**NPCs and characters:**

> Ask your AI: "Search the wiki for Taimi"

Returns pages about Taimi, including her role in the story and where she appears.

**Lore:**

> Ask your AI: "Search the wiki for Elder Dragons"

Returns pages about the Elder Dragons, their history, and their role in Tyria.

**Game systems:**

> Ask your AI: "Search the wiki for Mastery system"

Returns pages explaining Mastery tracks, experience requirements, and unlockable abilities.

## Verify it works

> Ask your AI: "Search the wiki for Guild Wars 2"

You should see results including the main Guild Wars 2 wiki page. If you get results with titles and URLs pointing to wiki.guildwars2.com, wiki search is working.

## Tips

**You do not need wiki search for item lookups.** The composite tools `get_item_by_name`, `get_tp_price_by_name`, and `get_item_recipe_by_name` use wiki search internally to resolve item names to IDs. Instead of searching the wiki for an item and then asking about its price, just ask directly:

> Ask your AI: "What's Mystic Coin selling for?"

The AI handles the wiki lookup behind the scenes.

**Fuzzy matching works.** You do not need exact page titles. Natural phrasing like "how does WvW work" or "what is the Mystic Forge" will return relevant results. That said, more specific queries produce more focused results.

**Use wiki search for non-item questions.** The composite tools handle items, prices, and recipes. For everything else -- game mechanics, lore, events, NPCs, achievements, locations -- wiki search is the right tool.

## Troubleshooting

### Problem: No results or irrelevant results
**Symptom**: The search returns nothing, or the results do not match what you were looking for.
**Cause**: The query terms did not match any wiki content closely enough.
**Solution**: Try different wording. Use the official in-game name when possible. For example, "WvW" may not work as well as "World vs World".

### Problem: Results are outdated
**Symptom**: Information in the results does not match the current state of the game.
**Cause**: The wiki is community-maintained and some pages lag behind game updates.
**Solution**: Check the wiki page's revision history to see when it was last updated. For very recent changes, in-game sources may be more accurate until the wiki is updated.

### Problem: Too many results
**Symptom**: The search returns many pages and it is hard to find the right one.
**Cause**: The query is too broad (for example, searching for "armor" returns hundreds of pages).
**Solution**: Be more specific. Instead of "armor", try "Blossoming Mist Shard armor" or "ascended heavy armor recipes".

## See also

- [Crafting Assistant](../tutorials/crafting/) -- uses wiki search internally for recipe and item lookups
- [Use Without an API Key](no-api-key/) -- full list of tools that work without authentication
- [Tools reference](../reference/tools/) -- specification for the `wiki_search` tool
