---
title: Find Valuable Items in Your Bank
weight: 13
---

# How to Find Valuable Items in Your Bank

**Goal**: Identify items in your bank vault that are worth gold on the Trading Post, so you can sell them or make informed decisions about your storage.

## Prerequisites

- GW2 MCP Server running and connected to your AI assistant -- see [Getting Started](../tutorials/getting-started/) if you need setup help
- A GW2 API key with **account** and **inventories** scopes enabled

## Steps

### 1. Ask your AI to check your bank contents

> Ask your AI: "What's the most valuable stuff in my bank?"

Your AI calls the `get_bank` tool to retrieve the contents of your bank vault. This returns every occupied slot with item names and quantities, just like opening your bank at a bank NPC in-game.

### 2. Let the AI cross-reference with Trading Post prices

Your AI chains the bank results into the `get_tp_prices` tool, passing the item IDs from your bank to look up current buy and sell prices on the Trading Post. This happens automatically as part of answering your question -- you do not need to make a separate request.

The AI then sorts and presents the most valuable items by their sell price. You should see something like a ranked list of your bank items with their estimated gold value based on current TP listings.

### 3. Review the results

Look through the list your AI provides. For each valuable item, consider:

- **Sell price** -- the current highest buy order or lowest sell listing on the TP
- **Quantity** -- how many you have stacked in your bank
- **Total value** -- quantity multiplied by the sell price

This gives you a clear picture of where your stored wealth is sitting.

### 4. Watch for account-bound items

Some items in your bank may be **account-bound** or **soulbound** -- these cannot be sold on the Trading Post regardless of their listed value. Common examples include:

- Ascended equipment
- Items from achievement rewards
- Account-bound crafting materials like Dragonite Ore or Bloodstone Dust
- Legendary weapons and armor

Your AI may flag these items, but if you see a surprisingly high-value item in the list, double-check whether it is tradeable before planning to sell it. You can ask your AI for clarification:

> Ask your AI: "Which of those items are actually tradeable on the TP?"

## Verify it works

After reviewing the list, spot-check a few items against the Trading Post in-game. Open the TP panel (default keybind: O), search for one of the items your AI listed, and compare the price. The values should be close, though they may fluctuate slightly as the market moves.

## Troubleshooting

### Problem: Authorization error when checking the bank
**Symptom**: Your AI reports a permissions or authorization error.
**Cause**: Your API key is missing the **inventories** scope.
**Solution**: Go to [Guild Wars 2 API Key Management](https://account.arena.net/applications) and create a new key with **account** and **inventories** scopes enabled. Update your `GW2_API_KEY` environment variable and restart the server.

### Problem: Bank appears empty
**Symptom**: The AI reports no items or very few items in your bank.
**Cause**: Your bank vault genuinely has few items stored, or you may be checking an account that has not been played much.
**Solution**: This is normal for newer accounts. Try asking about material storage instead: "What valuable materials do I have?"

### Problem: Prices seem outdated or missing
**Symptom**: Some items show no Trading Post price or the values seem off.
**Cause**: Certain items are not listed on the Trading Post (account-bound items, discontinued items), or the TP data may be cached.
**Solution**: Items with no TP listing simply cannot be sold. For items that do have listings, the prices reflect recent TP data and should be reasonably current.

## See also

- [Know Your Account tutorial](../tutorials/account-overview/) -- full walkthrough of exploring your bank, wallet, characters, and more
- [API Scopes reference](../reference/api-scopes/) -- which permissions each tool requires
