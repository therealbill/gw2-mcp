---
title: Compare Characters
weight: 15
---

# How to Compare Characters

**Goal**: View your character roster and compare details like gear, crafting levels, and builds across two or more characters.

## Prerequisites

- GW2 MCP Server running and connected to your AI assistant -- see [Getting Started](../tutorials/getting-started/) if you need setup help
- A GW2 API key with **account** and **characters** scopes enabled

## Steps

### 1. List all your characters

> Ask your AI: "List my characters"

Your AI calls the `get_characters` tool without specifying a name. The response is a list of every character on your account -- their names, just like the character selection screen in-game.

Take note of the character names you want to compare. Names are case-sensitive and must match exactly when you use them in later steps.

### 2. Inspect a specific character

Pick a character from the list and ask for details.

> Ask your AI: "Show me details for Zojja the Brave"

(Replace "Zojja the Brave" with your actual character name.)

Your AI calls `get_characters` with that character's name. The response includes:

- **Basic info** -- race, profession, level, gender, creation date, total play time
- **Crafting disciplines** -- which crafting professions are trained and their current levels (for example, Weaponsmith 500, Tailor 400)
- **Equipment** -- the gear equipped in each slot, with item names
- **Build tabs** -- saved builds including specializations and skill selections

This is equivalent to opening the Hero panel for that character, but in a format your AI can analyze.

### 3. Inspect a second character

Repeat the same question for another character.

> Ask your AI: "Now show me details for Rytlock Firebrand"

Your AI calls `get_characters` again with the second name, giving you the same level of detail.

### 4. Ask for a comparison

Now ask your AI to compare the two characters directly.

> Ask your AI: "Compare the gear on Zojja the Brave and Rytlock Firebrand"

Your AI calls `get_characters` for both names (if it has not already cached the results from the earlier steps) and presents a side-by-side comparison. Depending on what you ask, the comparison might cover:

- **Equipment** -- armor stats, weapon types, rarity (exotic vs. ascended vs. legendary)
- **Crafting** -- which disciplines each character has and at what levels
- **Builds** -- specialization and skill differences
- **Level and play time** -- who has been played more

You can focus the comparison on whatever aspect matters to you:

> Ask your AI: "Which of my characters has the highest crafting levels?"

> Ask your AI: "Do any of my characters have ascended gear equipped?"

## Verify it works

Pick one of the characters from the comparison and log into that character in-game. Open the Hero panel (default keybind: H) and check that the profession, level, and equipped gear match what your AI reported. The equipment names and crafting levels should be consistent.

## Troubleshooting

### Problem: Character name not found
**Symptom**: Your AI reports that the character does not exist.
**Cause**: The name does not match exactly. Character names are case-sensitive and include spaces.
**Solution**: Run "List my characters" first, then copy the exact name from the results. Make sure you include any spaces and use the correct capitalization.

### Problem: Authorization error
**Symptom**: Your AI reports a permissions or authorization error when requesting character details.
**Cause**: Your API key is missing the **characters** scope.
**Solution**: Go to [Guild Wars 2 API Key Management](https://account.arena.net/applications) and create a new key with **account** and **characters** scopes enabled. Update your `GW2_API_KEY` and restart the server.

### Problem: Equipment or builds appear incomplete
**Symptom**: Some gear slots or build tabs show as empty.
**Cause**: The character genuinely has empty equipment slots or unused build tabs. This is normal for characters that are not fully geared or have not unlocked additional build template slots.
**Solution**: This is expected behavior, not an error. Only occupied slots and configured build tabs are returned.

## See also

- [Know Your Account tutorial](../tutorials/account-overview/) -- full walkthrough of exploring your characters, bank, wallet, and more
- [API Scopes reference](../reference/api-scopes/) -- which permissions each tool requires
