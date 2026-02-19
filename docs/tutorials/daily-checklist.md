---
title: Daily Checklist
weight: 5
---

# Daily Checklist

Use your AI assistant as a daily companion to track Wizard's Vault objectives, achievements, and instanced content clears -- all without leaving your conversation.

By the end of this tutorial, you will have:
- Checked your daily and weekly Wizard's Vault objectives and tracked Astral Acclaim progress
- Browsed the Wizard's Vault reward listings
- Looked up today's daily achievements
- Reviewed your raid, dungeon, and world boss completion for the current reset period

**Estimated time:** 20 minutes

## Prerequisites

Before starting, make sure you have:
- Completed the [Getting Started](../getting-started/) tutorial (GW2 MCP server running, Claude Desktop configured)
- A GW2 API key with **account** and **progression** scopes enabled (for personal progress tracking)

If your API key does not have the `progression` scope, you can still follow along -- the Wizard's Vault sections will show objective lists without your personal completion status, and the raid/dungeon/world boss sections will not work. See the [API Scopes reference](../reference/api-scopes/) for details.

## Step 1: Check your daily Wizard's Vault objectives

The Wizard's Vault replaced the old daily achievement system. Each day you get a set of objectives worth Astral Acclaim -- the currency you spend on Wizard's Vault rewards.

Open a conversation with your AI assistant and ask:

> "What are today's Wizard's Vault objectives?"

Your assistant calls the `get_wizards_vault_objectives` tool with type `daily` and returns your daily objective list. Because you have an API key configured, the response includes your personal progress: which objectives you have already completed, how many Astral Acclaim each one is worth, and your progress on partially completed objectives.

If you did not have an API key configured, you would still see the objective list -- but without any completion tracking. The objectives themselves are public data.

### Checkpoint

Your assistant responded with a list of daily objectives. Each objective has a name, an Astral Acclaim reward value, and a completion status (complete or in-progress). Compare this to the Wizard's Vault panel in-game (open it with the default keybind or through the hero panel) to confirm the objectives match.

## Step 2: Check your weekly Wizard's Vault objectives

Weekly objectives work the same way as dailies but reset once per week and tend to be worth more Astral Acclaim. These are the longer-term tasks like completing meta events, fractals, or PvP matches.

> "What are this week's Wizard's Vault objectives?"

Your assistant calls `get_wizards_vault_objectives` with type `weekly`. The response shows your weekly objectives with the same completion tracking as the dailies.

Weekly objectives reset on Monday. If you are checking mid-week, you will see partial progress on objectives you have started but not finished.

### Checkpoint

You see a list of weekly objectives with their Astral Acclaim values and your completion progress. The weekly objectives are different from the daily ones and generally require more effort.

## Step 3: Browse Wizard's Vault rewards

Now that you know how much Astral Acclaim you are earning, check what you can spend it on.

> "What can I buy from the Wizard's Vault?"

Your assistant calls the `get_wizards_vault_listings` tool. Because you have an API key, the response shows the full reward catalog along with your purchase status -- which items you have already bought and how many times (some rewards are repeatable).

Without an API key, the same tool returns the reward list and prices, but without any purchase history.

### Checkpoint

You see a list of Wizard's Vault rewards with their Astral Acclaim costs. The list includes a mix of items: legendary crafting materials, skins, and other account upgrades. Items you have already purchased this season are marked as such.

## Step 4: Check the current Wizard's Vault season

Each Wizard's Vault season runs for a set period and has its own reward track. Check what season is active and when it ends.

> "What's the current Wizard's Vault season?"

Your assistant calls `get_wizards_vault` and returns the current season name, start date, end date, and any other season metadata. This is public data and does not require an API key.

### Checkpoint

You see the current season name and its date range. This tells you how much time you have left to earn Astral Acclaim and claim rewards before the season rotates.

## Step 5: Look up today's daily achievements

Beyond the Wizard's Vault, Guild Wars 2 still has daily achievements that rotate each day. These are organized by game mode: PvE, PvP, WvW, fractals, and special events.

> "What are today's daily achievements?"

Your assistant calls `get_daily_achievements` and returns today's daily achievements along with tomorrow's, so you can plan ahead. This is public data and does not require an API key.

The response includes achievement names and the game mode they belong to. Your assistant may also show the level ranges required for PvE dailies, since some are restricted to specific level brackets.

### Checkpoint

You see two lists: today's dailies and tomorrow's dailies. Each entry includes the achievement name and its category (PvE, PvP, WvW, fractals, or special). Compare the PvE dailies to the in-game achievement panel to confirm they match.

## Step 6: Check your raid clears for the week

Raid encounters in Guild Wars 2 reset weekly on Monday. Each wing has multiple encounters, and you can only earn rewards from each encounter once per week. Keeping track of which encounters you have cleared is essential for planning your raid week.

> "Which raids have I cleared this week?"

Your assistant calls `get_account_dailies` with type `raids`. The response is a list of encounter IDs that you have completed since the last weekly reset. Your assistant will translate these IDs into readable encounter names so you can see exactly which bosses you have already downed.

If the list is empty, you have not cleared any raid encounters this week.

### Checkpoint

You see a list of completed raid encounters (or an empty list if you have not raided this week). If you have been raiding, the encounter names should match what you have cleared in-game. For example, if you cleared the first wing of Forsaken Thicket, you would see encounters like Vale Guardian, Gorseval, and Sabetha listed.

## Step 7: Check your dungeon paths for the day

Dungeon paths reset daily. Each dungeon has one story path and multiple explorable paths, and you get bonus rewards the first time you complete each path per day.

> "Which dungeon paths have I done today?"

Your assistant calls `get_account_dailies` with type `dungeons`. The response lists the dungeon path IDs you have completed since the last daily reset.

### Checkpoint

You see your completed dungeon paths for today. If you have not run any dungeons today, the list is empty. If you ran Ascalonian Catacombs path 1, you would see that specific path listed.

## Step 8: Check your world boss completions for the day

World bosses are open-world encounters on a fixed schedule. You get bonus rewards from each world boss once per day. Tracking which ones you have done helps you decide whether to chase the next one on the timer.

> "Which world bosses have I done today?"

Your assistant calls `get_account_dailies` with type `worldbosses`. The response lists the world boss IDs you have defeated since the last daily reset.

### Checkpoint

You see a list of world bosses you have completed today. This might include bosses like Tequatl, Shadow Behemoth, or the Claw of Jormag. If the list is empty, you have not defeated any world bosses since today's reset.

## What you learned

In this tutorial, you used your AI assistant to run through a full daily checklist:

- **Wizard's Vault dailies and weeklies** -- Tracked your Astral Acclaim progress using `get_wizards_vault_objectives` with both `daily` and `weekly` types
- **Wizard's Vault rewards** -- Browsed what you can buy and what you have already purchased using `get_wizards_vault_listings`
- **Wizard's Vault season** -- Checked the current season timeline using `get_wizards_vault`
- **Daily achievements** -- Reviewed today's and tomorrow's rotating achievements using `get_daily_achievements`
- **Raid clears** -- Checked your weekly raid encounter completions using `get_account_dailies` with type `raids`
- **Dungeon paths** -- Checked your daily dungeon path completions using `get_account_dailies` with type `dungeons`
- **World bosses** -- Checked your daily world boss completions using `get_account_dailies` with type `worldbosses`

You also learned the distinction between tools that require authentication (raid/dungeon/world boss clears), tools that are enhanced by authentication (Wizard's Vault objectives and listings), and tools that work without any API key (daily achievements, season info).

## Next steps

Now that you know how to check your daily and weekly progress, explore further:

- **Automate your Wizard's Vault routine** -- See the [Wizard's Vault Daily](../how-to/wizards-vault-daily/) how-to guide for tips on building this into a daily habit
- **Track raid clears across the week** -- See the [Track Raid Clears](../how-to/track-raid-clears/) how-to guide for organizing your weekly raid schedule
- **Browse all available tools** -- See the [Tools reference](../reference/tools/) for the full list of 37 tools
- **Understand API key permissions** -- See the [API Scopes reference](../reference/api-scopes/) for which scopes each tool requires
