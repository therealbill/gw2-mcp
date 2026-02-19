---
title: Crafting vs Buying
weight: 16
---

# Crafting vs Buying

**Goal**: Determine whether it is cheaper to craft an item yourself or buy it directly from the Trading Post.

## Prerequisites

None -- all lookups use public GW2 API endpoints. No API key is required.

## Steps

### 1. Ask your AI for a cost comparison

> Ask your AI: "Is it cheaper to craft or buy Deldrimor Steel Ingot?"

You can use any craftable item name. The AI handles the entire multi-step analysis automatically.

### 2. AI looks up the recipe

Behind the scenes, the AI calls `get_item_recipe_by_name` to find the recipe for the item you named. This returns the list of ingredients and their quantities -- for example, Deldrimor Steel Ingot requires Iron Ingots, Steel Ingots, Darksteel Ingots, and a Lump of Mithrillium.

### 3. AI gets ingredient prices

For each ingredient, the AI calls `get_tp_price_by_name` to look up the current Trading Post price. It checks both the highest buy order and lowest sell listing so you get an accurate picture of what you would actually pay.

### 4. AI sums the crafting cost

The AI multiplies each ingredient price by the quantity required and totals them up. This gives you the cost to craft the item from materials purchased on the Trading Post.

### 5. AI gets the finished item price

The AI calls `get_tp_price_by_name` one more time for the finished item itself, giving you the current buy-it-now price on the Trading Post.

### 6. AI reports the result

The AI compares the crafting cost against the Trading Post price and tells you:

- Which option is cheaper
- The exact gold difference between the two
- The individual ingredient costs that make up the crafting total

A typical response looks like:

> "Crafting a Deldrimor Steel Ingot costs approximately 1g 20s from ingredients, while buying one on the Trading Post costs 1g 85s. Crafting saves you about 65 silver per ingot."

## Important: Trading Post tax

The comparison above assumes you are crafting or buying the item **for your own use**. If you plan to **sell** the crafted item on the Trading Post, remember to factor in the 15% listing fee (5% listing fee + 10% exchange fee). This tax is deducted from your sale proceeds.

For example, if a crafted item sells for 10 gold, you only receive 8 gold 50 silver after fees. The AI can account for this if you ask:

> Ask your AI: "Is it profitable to craft and sell Deldrimor Steel Ingots on the TP after tax?"

## Troubleshooting

### Problem: AI says the item has no recipe
**Symptom**: The AI reports that no recipe was found for the item.
**Cause**: The item is not craftable (dropped loot, gem store item, etc.) or you may have the name slightly wrong.
**Solution**: Double-check the exact in-game item name. Some items have similar names -- "Mystic Coin" is a drop, not a crafted item.

### Problem: Prices seem outdated or wrong
**Symptom**: The AI reports prices that do not match what you see in-game.
**Cause**: Trading Post prices change constantly. The API returns a snapshot at the moment of the query, which may differ by the time you check in-game.
**Solution**: Re-run the query for a fresh price check. For volatile items, ask the AI to check prices again right before you commit to buying or crafting.

### Problem: Multi-step recipes are not fully costed
**Symptom**: The AI prices an intermediate material at its TP cost instead of breaking it down further.
**Cause**: Some items require crafting sub-components that are themselves craftable. The AI may stop at one level of depth by default.
**Solution**: Ask the AI to break down the full recipe tree:
> Ask your AI: "What is the total crafting cost of Dawn, including all sub-recipes?"

## See also

- [Crafting Assistant tutorial](../tutorials/crafting/) -- full walkthrough of crafting research with your AI
- [Tools reference](../reference/tools/) -- details on `get_item_recipe_by_name` and `get_tp_price_by_name`
