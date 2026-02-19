---
title: Crafting Assistant
weight: 4
---

# Crafting Assistant

Use your AI assistant as a crafting research companion -- look up recipes, check material prices, and figure out whether to craft or buy from the Trading Post.

By the end of this tutorial, you will have:
- Looked up a crafting recipe by item name
- Retrieved detailed item metadata
- Found recipes that use a specific material
- Estimated the cost to craft an item
- Compared crafting cost against the Trading Post buy price

**Estimated time**: 30 minutes

## Prerequisites

Before starting, make sure you have completed the [Getting Started](../getting-started/) tutorial. Your AI client should be connected to the GW2 MCP server and able to call tools.

None of the tools in this tutorial require a GW2 API key. Crafting data, item metadata, and Trading Post prices are all public endpoints.

## What you will build

You will walk through five crafting research tasks, each building on the last. By the end, you will know how to use natural language prompts to drive a full crafting cost analysis -- the kind of research that normally means switching between the wiki, the Trading Post, and a spreadsheet.

The item we will investigate is **Dawn**, one of the two precursor weapons needed to craft the legendary greatsword Sunrise.

## Section 1: Look up a recipe by name

Let's start with the most common crafting question: what ingredients does this item need?

> Ask your AI: "How do I craft Dawn?"

The AI calls the `get_item_recipe_by_name` tool with the name "Dawn". Behind the scenes, the server searches the GW2 wiki to find Dawn's item page, extracts recipe data from it, fetches full recipe details from the API, and resolves every ingredient ID to a human-readable name.

The response tells you:

- **Output item**: Dawn (the precursor greatsword)
- **Crafting discipline**: Weaponsmith
- **Minimum rating**: 400
- **Ingredients**: A list of components with names and quantities -- for example, items like Gift of Metal, Gift of Light, an Icy Runestone, and other materials that vary by recipe

Each ingredient shows both its item ID and its name, so the AI can present them in plain language rather than raw numbers.

> **What just happened?**
>
> The `get_item_recipe_by_name` tool did three things in one call: searched the wiki for "Dawn", found its recipe IDs, then fetched complete recipe details with resolved ingredient names. Without this composite tool, you would need to search the wiki manually, extract the recipe ID, call the recipes endpoint, then look up each ingredient ID separately.

### Checkpoint

Your AI responded with a recipe for Dawn that includes a list of named ingredients, the crafting discipline (Weaponsmith), and a minimum crafting rating. If the AI says "No recipes found," try the exact name "Dawn" -- the wiki search works best with precise item names.

## Section 2: Look up item details

Now let's look at a related item. Dusk is Dawn's counterpart -- the precursor for the legendary greatsword Twilight.

> Ask your AI: "Tell me about Dusk"

The AI calls the `get_item_by_name` tool with the name "Dusk". The server searches the wiki, finds the item ID, and returns full item metadata from the API.

The response includes:

- **Name**: Dusk
- **Rarity**: Exotic
- **Type**: Weapon (Greatsword)
- **Level**: 80 (required level to equip)
- **Description**: The item's in-game flavor text
- **Vendor value**: The amount you would get selling it to a vendor (in gold, silver, and copper)
- **Chat link**: A GW2 chat code you can paste in-game to link the item

This is the same data you see on an item's wiki page or in-game tooltip, but available through a single natural language prompt.

### Checkpoint

Your AI responded with Dusk's metadata showing it as an Exotic rarity, level 80 greatsword. If you get a different item (the wiki search is fuzzy), try "Dusk" by itself -- shorter, exact names produce the best matches.

## Section 3: Find recipes that use a material

Sometimes you want to search in the other direction: not "what does this item need?" but "what can I make with this material?"

> Ask your AI: "What can I craft with Glob of Ectoplasm?"

This is a more complex question. The GW2 API lets you search for recipes by input item ID, but it needs the numeric ID, not the name. Your AI handles this in two steps:

1. **Resolve the name to an ID**: The AI calls `get_item_by_name` with "Glob of Ectoplasm" to get its item ID (19721).
2. **Search for recipes**: The AI calls `search_recipes` with that ID as the input parameter. The API returns a list of recipe IDs that use Glob of Ectoplasm as an ingredient.

The result is a list of recipe IDs. Depending on your AI client, it may then look up details for some of those recipes to show you what they produce. Glob of Ectoplasm is used in hundreds of recipes, so the AI will likely summarize rather than list every single one.

> **What just happened?**
>
> Unlike the previous sections where one tool did everything, this task required the AI to chain two tools together. It recognized that `search_recipes` needs a numeric item ID, so it first called `get_item_by_name` to resolve the name. This kind of multi-step reasoning is where an AI assistant adds the most value -- it handles the plumbing between tools automatically.

### Checkpoint

Your AI responded with recipe IDs (or details about recipes) that use Glob of Ectoplasm as an ingredient. You should see a substantial number of results -- Glob of Ectoplasm is one of the most widely used crafting materials in the game.

## Section 4: Estimate crafting cost

Now let's combine recipe data with Trading Post prices to estimate what it would actually cost to craft something.

> Ask your AI: "How much would it cost to craft Dawn on the Trading Post?"

The AI needs to pull data from multiple tools to answer this:

1. **Get the recipe**: Call `get_item_recipe_by_name` with "Dawn" to get the ingredient list with names and quantities.
2. **Price each ingredient**: Call `get_tp_price_by_name` for each ingredient to get current buy and sell prices.

The AI then multiplies each ingredient's price by the required quantity and adds them up to give you a total estimated crafting cost.

The response typically shows:

- Each ingredient with its quantity and current Trading Post price
- A total estimated cost in gold
- Whether it used buy prices (instant purchase) or sell prices (placing buy orders) -- or both for comparison

Keep in mind this is an estimate for the top-level recipe. Some ingredients (like Gift of Metal) are themselves crafted from other materials. A full cost breakdown would need to recursively look up sub-recipes, which you can do by asking follow-up questions about individual components.

### Checkpoint

Your AI responded with a cost breakdown listing Dawn's ingredients and their Trading Post prices, with a total estimated crafting cost. The prices reflect the live market, so they will differ from any guide you find online. If the AI could not find a price for an ingredient, that ingredient may not be tradeable on the Trading Post (account-bound items like Gifts from the Mystic Forge have no TP listing).

## Section 5: Compare crafting vs buying

The ultimate crafting question: should you craft the item yourself or buy it outright?

> Ask your AI: "Is it cheaper to craft or buy Dawn?"

The AI combines everything from the previous sections:

1. **Crafting cost**: Calls `get_item_recipe_by_name` for the recipe, then `get_tp_price_by_name` for each tradeable ingredient to estimate the total crafting cost.
2. **Buy price**: Calls `get_tp_price_by_name` with "Dawn" to get the current Trading Post price for the finished item.
3. **Comparison**: Compares the two numbers and tells you which option is cheaper.

The response typically includes:

- The estimated crafting cost (sum of ingredient prices)
- The current Trading Post price for Dawn
- The difference between the two
- A recommendation on which is more economical

This is where the AI assistant really shines. It coordinates four or five tool calls, handles the math, and presents a clear answer -- work that would normally involve tabbing between the wiki, the Trading Post, and a calculator.

> **What just happened?**
>
> The AI chained `get_item_recipe_by_name` and multiple `get_tp_price_by_name` calls together without you having to orchestrate anything. You asked a single question in plain language, and the AI figured out the right sequence of tool calls to answer it. This is the core value of the MCP server approach -- the tools are simple building blocks, and the AI composes them into workflows.

### Checkpoint

Your AI responded with both the estimated crafting cost and the Trading Post buy price for Dawn, along with a comparison. The answer will change over time as market prices fluctuate. If some ingredients show no price, the AI should note that those items are account-bound or otherwise untradeable.

## A note on name matching

The composite tools (`get_item_by_name`, `get_item_recipe_by_name`, `get_tp_price_by_name`) use wiki search internally to resolve item names to IDs. This means:

- **Exact names work best**: "Glob of Ectoplasm" will always find the right item. "Ecto" might not.
- **Fuzzy matching has limits**: Common abbreviations that players use (like "Ecto" or "MC" for Mystic Coin) may not resolve correctly. Use the full item name when precision matters.
- **Ambiguous names**: If an item name matches multiple wiki pages, the server uses the top search result. For items that share names with other game concepts, adding context to your prompt helps the AI pick the right tool and name.

## What you learned

In this tutorial, you used five different crafting workflows:

- **Recipe lookup** (`get_item_recipe_by_name`) -- Get a complete recipe with named ingredients from a single item name.
- **Item details** (`get_item_by_name`) -- Retrieve full item metadata including rarity, type, level, and vendor value.
- **Reverse recipe search** (`get_item_by_name` + `search_recipes`) -- Find what recipes use a given material, by resolving the name to an ID first.
- **Cost estimation** (`get_item_recipe_by_name` + `get_tp_price_by_name`) -- Price out ingredients on the Trading Post to estimate crafting cost.
- **Craft vs buy** (all of the above) -- Compare total crafting cost against the finished item's Trading Post price.

The key insight is that these tools are composable. Each one does something simple -- look up an item, get a recipe, check a price -- and the AI chains them together to answer complex questions. You do not need to know which tools exist or how they connect. Ask your question in natural language, and the AI handles the rest.

## Next steps

Now that you know how the crafting tools work, explore further:

- **[Compare crafting vs buying](../how-to/crafting-vs-buying/)** -- A step-by-step how-to for running detailed cost comparisons, including sub-recipe breakdowns
- **[Find recipes for a material](../how-to/find-recipes-for-material/)** -- A how-to guide for the reverse recipe search workflow, including tips for narrowing results
- **Browse all tools** -- See the [Tools reference](../reference/tools/) for the complete list of available tools

Try asking your AI about other craftable items -- legendary weapons, ascended gear, or even basic crafting components. The same patterns you learned here apply to any item in the game.
