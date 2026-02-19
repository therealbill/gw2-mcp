---
title: Find Recipes for a Material
weight: 17
---

# Find Recipes for a Material

**Goal**: Find out what you can craft with a specific material you have in your inventory.

## Prerequisites

None -- recipe and item lookups use public GW2 API endpoints. No API key is required.

## Steps

### 1. Ask your AI what uses a material

> Ask your AI: "What can I craft with Glob of Ectoplasm?"

Use the exact in-game name of the material you want to look up. The AI handles the multi-step lookup automatically.

### 2. AI resolves the item name to an ID

The AI calls `get_item_by_name` to find the item and retrieve its numeric ID. The GW2 API uses numeric IDs internally, so this name-to-ID resolution step is required before searching recipes. This happens automatically -- you do not need to know the ID yourself.

### 3. AI searches for recipes using that ingredient

With the item ID in hand, the AI calls `search_recipes` with the `input` parameter set to that ID. This returns a list of recipe IDs for every recipe that uses the material as an ingredient.

### 4. AI fetches recipe details

The AI calls `get_recipes` with the returned recipe IDs to get the full details: what each recipe produces, what discipline crafts it (Armorsmith, Weaponsmith, etc.), the required rating, and the other ingredients involved.

### 5. AI presents the results

The AI gives you a readable list of items you can craft with that material. A typical response looks like:

> "Glob of Ectoplasm is used in 200+ recipes. Here are some notable ones:
> - **Deldrimor Steel Ingot** (Weaponsmith 500) -- 1 Glob of Ectoplasm + other materials
> - **Elonian Leather Square** (Leatherworker 500) -- 1 Glob of Ectoplasm + other materials
> - **Bolt of Damask** (Tailor 500) -- 1 Glob of Ectoplasm + other materials
> - ..."

For common materials like Glob of Ectoplasm, the list can be long. Narrow it down by asking a follow-up question:

> Ask your AI: "Which of those are Huntsman recipes at rating 400 or above?"

## Important: search_recipes takes IDs, not names

The `search_recipes` tool only accepts numeric item IDs, not item names. This is why the AI always calls `get_item_by_name` first to resolve the name. This two-step process happens behind the scenes -- you just ask your question in plain language and the AI chains the tools together.

## Troubleshooting

### Problem: AI returns no recipes for a material
**Symptom**: The AI says no recipes were found.
**Cause**: The material may not be used as a crafting input, or the item name may be slightly off. Some items that look like crafting materials (such as certain trophies) are only used in Mystic Forge recipes, which are not included in the standard recipe search.
**Solution**: Verify the exact item name. If the material is a Mystic Forge ingredient rather than a standard crafting ingredient, ask the AI to search the wiki instead:
> Ask your AI: "Search the wiki for Mystic Forge recipes using Gift of Exploration"

### Problem: Too many results to be useful
**Symptom**: The AI returns hundreds of recipes and the list is overwhelming.
**Cause**: Common materials like ore, wood, or cloth are used in a large number of recipes across all disciplines.
**Solution**: Add constraints to your question to filter the results:
> Ask your AI: "What Artificer recipes above rating 400 use Orichalcum Ore?"

### Problem: AI shows recipe IDs instead of item names
**Symptom**: The response contains numeric IDs rather than readable item names.
**Cause**: The AI may not have fetched the output item details for all recipes.
**Solution**: Ask the AI to look up the item names:
> Ask your AI: "Can you show me the item names for those recipes?"

## See also

- [Crafting Assistant tutorial](../tutorials/crafting/) -- full walkthrough of crafting research with your AI
- [Tools reference](../reference/tools/) -- details on `get_item_by_name`, `search_recipes`, and `get_recipes`
