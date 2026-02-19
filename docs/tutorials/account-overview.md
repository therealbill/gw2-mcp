---
title: Know Your Account
weight: 3
---

# Know Your Account

Explore your Guild Wars 2 account through your AI assistant in about 30 minutes.

By the end of this tutorial, you will have:
- Checked your wallet balances across all currencies
- Browsed your bank vault, material storage, and shared inventory
- Listed your characters and inspected one in detail
- Looked up your account unlocks
- Viewed your overall account information

## Prerequisites

Before starting, make sure you have:
- The [Getting Started](../getting-started/) tutorial complete (GW2 MCP server running with Claude Desktop or another AI client)
- A GW2 API key with these scopes enabled: **account**, **characters**, **inventories**, **wallet**, **unlocks**

If you are not sure which scopes your key has, check at [Guild Wars 2 API Key Management](https://account.arena.net/applications). For details on what each scope enables, see the [API Scopes reference](../../reference/api-scopes/).

## What you will build

This is not a coding project -- it is a guided tour of your own GW2 account data. You will learn the prompts that let your AI assistant pull up every major piece of account information: currencies, storage, characters, and unlocks. By the end, you will know exactly how to ask your AI about any aspect of your account.

## Step 1: Check your wallet

Your wallet holds every currency you have earned in Guild Wars 2 -- not just gold, but karma, gems, laurels, spirit shards, badges of honor, and dozens more.

> Ask your AI: "What's in my wallet?"

Your AI calls the `get_wallet` tool and returns a list of all your currencies with their current balances. You should see entries like:

- **Coin** -- your gold, silver, and copper (the main in-game currency)
- **Karma** -- earned from events and hearts, spent at karma merchants
- **Gem** -- the premium currency, convertible to gold
- **Laurel** -- earned from daily login rewards, used to buy ascended materials and other items
- **Spirit Shard** -- earned after reaching level 80, used in Mystic Forge recipes
- **Badge of Honor** -- earned from World vs. World
- **Guild Commendation** -- earned from guild missions

You may see many more currencies depending on which content you play. The wallet tool returns all of them at once.

### Checkpoint

You should see a list of currencies with names and amounts. The gold amount should match what you see on your character's wallet panel in-game. If you get an authorization error, verify your API key has the **wallet** scope enabled.

## Step 2: Browse your bank

Your bank vault is the shared storage accessible from any bank NPC or crafting station in-game. It holds up to 735 item slots across multiple tabs.

> Ask your AI: "What's in my bank?"

Your AI calls the `get_bank` tool and returns the contents of every bank slot. Each item comes back with its name, so you do not need to look up item IDs yourself. Empty slots are skipped.

The response lists items with their names and quantities. You might see things like stacks of crafting materials you stashed, equipment you are saving for another character, or consumables you bought from the Trading Post.

### Checkpoint

You should see a list of items with names and quantities that matches your in-game bank. If you have a lot of items, the list will be long -- that is expected. If you get an authorization error, verify your API key has the **inventories** scope enabled.

## Step 3: Check material storage

Material storage is a separate bank section dedicated to crafting materials. It holds up to 250 of each material (or 2000 with an upgrade) and keeps them organized by category: basic crafting materials, fine crafting materials, common ascended materials, and so on.

> Ask your AI: "What materials do I have?"

Your AI calls the `get_materials` tool. The response lists every material slot that contains at least one item, showing the material name and how many you have stored.

This is useful for checking whether you have enough materials before starting a crafting project, or for spotting valuable materials you have been accumulating.

### Checkpoint

You should see a list of crafting materials with names and counts. Compare a few entries to your in-game material storage panel to confirm the numbers match.

## Step 4: Look at shared inventory

Shared inventory slots appear at the top of every character's inventory. They are account-wide slots that hold items you want accessible on all characters -- things like salvage kits, gathering tools, or utility items.

> Ask your AI: "Show my shared inventory slots"

Your AI calls the `get_inventory` tool. The response shows what is in each of your shared inventory slots, with item names and quantities. If a slot is empty, it may be omitted from the results.

Not every account has shared inventory slots -- they are purchased from the Gem Store. If you have none, you will see an empty result. That is normal.

### Checkpoint

You should see a list of items in your shared slots, or an empty response if you have not purchased any shared inventory slots. If you do have shared slots, verify the items match what you see at the top of any character's inventory in-game.

## Step 5: List your characters

Time to look at your roster. Every GW2 account has one or more characters, each with their own profession, level, equipment, and crafting disciplines.

> Ask your AI: "List my characters"

Your AI calls the `get_characters` tool without specifying a character name. The response is a list of all your character names. This gives you a quick overview of your roster.

Take note of one character name from the list -- you will use it in the next step.

### Checkpoint

You should see a list of character names that matches your in-game character selection screen. If you get an authorization error, verify your API key has the **characters** scope enabled.

## Step 6: Inspect a character

Now pick one character from the list and ask for the full details. Replace the name below with one of your actual character names.

> Ask your AI: "Show me details for Zojja the Brave"

(Use your own character's name instead of "Zojja the Brave.")

Your AI calls the `get_characters` tool with your character's name. The response contains a wealth of information:

- **Basic info** -- name, race, profession, level, gender, creation date, play time
- **Crafting disciplines** -- which crafting professions this character has trained and their current levels (for example, Armorsmith 500, Artificer 400)
- **Equipment** -- the gear currently equipped in each slot, with item names
- **Build tabs** -- your saved builds, including specialization and skill selections
- **Skills and specializations** -- the active build configuration

This is like opening your Hero panel in-game, but all at once and in a format your AI can analyze or compare.

### Checkpoint

You should see detailed information about the character you named. Verify that the profession, level, and crafting disciplines match what you see in-game. If the character name does not match exactly (including spaces and capitalization), you may get an error -- try copying the exact name from the list in Step 5.

## Step 7: Check your unlocks

Over the course of playing GW2, you unlock skins, dyes, miniatures, titles, recipes, and more. These unlocks are account-wide and permanent.

> Ask your AI: "What skins have I unlocked?"

Your AI calls the `get_account_unlocks` tool with the type set to `skins`. The response is a list of unlocked skin IDs. Depending on how long you have played, this list could be very large -- veteran players often have hundreds or thousands of skins unlocked.

The tool returns numeric IDs rather than names. Your AI can look up specific IDs if you want to know what a particular skin is, but the main value here is seeing the total count and checking whether a specific skin is in your collection.

You can check other unlock types by changing what you ask for. Try any of these:

- "What dyes have I unlocked?" -- dye colors for the wardrobe
- "What minis have I unlocked?" -- miniature pets
- "What titles have I unlocked?" -- titles for your nameplate
- "What recipes have I unlocked?" -- crafting recipes your account knows
- "What mount skins have I unlocked?" -- appearance options for your mounts
- "What mount types have I unlocked?" -- the base mount types (raptor, springer, skimmer, etc.)

### Checkpoint

You should see a list of numeric IDs for the unlock type you requested. If you asked for skins, the count should roughly match the number shown in your in-game wardrobe. If you get an authorization error, verify your API key has the **unlocks** scope enabled.

## Step 8: View your account overview

Finally, let's pull up the top-level account information that ties everything together.

> Ask your AI: "Tell me about my account"

Your AI calls the `get_account` tool. The response includes:

- **Account name** -- your display name (the one with the four-digit number, like "PlayerName.1234")
- **World** -- your home world for World vs. World
- **Access level** -- which editions of the game you own (core, Heart of Thorns, Path of Fire, End of Dragons, Secrets of the Obscure, Janthir Wilds)
- **Guilds** -- the guild IDs your account belongs to
- **Other details** -- account age, commander tag status, and more

This is useful for confirming which expansions you have access to or checking which guilds your account is part of.

### Checkpoint

You should see your account name, home world, and access level. Verify the account name matches what you see in-game (press F11 or check the top of your Friends panel). The access level should list the expansions you own.

## What you learned

You now know how to ask your AI assistant about every major aspect of your GW2 account:

| Prompt | What it does |
|--------|-------------|
| "What's in my wallet?" | Shows all currency balances |
| "What's in my bank?" | Lists bank vault contents with item names |
| "What materials do I have?" | Shows material storage quantities |
| "Show my shared inventory slots" | Lists shared inventory items |
| "List my characters" | Shows all character names |
| "Show me details for [name]" | Returns full character info including gear, crafting, and builds |
| "What skins have I unlocked?" | Lists unlocked IDs for any unlock type |
| "Tell me about my account" | Shows account name, world, access, and guilds |

These are the building blocks. Once your AI has this data, you can ask follow-up questions like "Which of my characters has the highest crafting level?" or "How much karma do I have?" and your AI will use these same tools to answer.

## Next steps

Now that you know how to explore your account, try these task-focused guides:

- **Compare your characters side by side** -- See [Compare Characters](../how-to/compare-characters/) for a focused guide on inspecting gear and builds across your roster
- **Find valuable items in your bank** -- See [Find Valuable Items in Your Bank](../how-to/find-bank-valuables/) to cross-reference your bank contents with Trading Post prices
- **Browse all available tools** -- See the [Tools reference](../../reference/tools/) for the complete list of 37 tools
- **Understand API key permissions** -- See the [API Scopes reference](../../reference/api-scopes/) for which scopes each tool requires

## Troubleshooting

### Problem: Authorization error on a tool
**Symptom**: Your AI reports an error about missing permissions or invalid API key.
**Solution**: Go to [Guild Wars 2 API Key Management](https://account.arena.net/applications) and verify your key has the required scopes. This tutorial needs **account**, **characters**, **inventories**, **wallet**, and **unlocks**. If any are missing, create a new key with all scopes enabled.

### Problem: Character name not found
**Symptom**: You ask for character details but get an error about the character not existing.
**Solution**: Character names are case-sensitive and must include spaces exactly as shown. Run "List my characters" first, then copy the exact name from the results.

### Problem: Empty bank or material storage
**Symptom**: The bank or materials response comes back empty or very short.
**Solution**: This is normal if you have not stored items in those locations yet. Try checking on an account that has been played for a while. Shared inventory slots in particular require a Gem Store purchase.

### Need more help?
- Browse the [Tools reference](../../reference/tools/) for details on each tool's parameters and behavior
- Check the [API Scopes reference](../../reference/api-scopes/) for scope requirements
