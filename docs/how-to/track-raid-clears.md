---
title: Track Raid and Dungeon Clears
weight: 19
---

# How to Track Raid and Dungeon Clears

**Goal**: See which raids, dungeons, and world bosses you have completed this reset period so you know what still has rewards available.

## Prerequisites

- GW2 MCP Server running and connected to your AI assistant -- see [Getting Started](../tutorials/getting-started/)
- A GW2 API key with **account** and **progression** scopes -- see [API Scopes](../reference/api-scopes/)
- All three tools in this guide require an API key. Without one, these lookups will not work.

## Steps

### 1. Check your raid clears

> Ask your AI: "Which raids have I cleared this week?"

Your assistant retrieves a list of raid encounters you have completed since the last weekly reset. Raid rewards reset every Monday. If you have cleared Wing 1 of Forsaken Thicket, for example, you would see encounters like Vale Guardian, Gorseval, and Sabetha listed.

An empty list means you have not cleared any raid encounters this week.

### 2. Check your dungeon paths

> Ask your AI: "Which dungeon paths have I done today?"

Your assistant retrieves your completed dungeon paths for the current day. Each dungeon has one story path and multiple explorable paths, and you earn bonus rewards the first time you complete each path per day. Dungeon rewards reset daily.

### 3. Check your world boss completions

> Ask your AI: "Which world bosses have I done today?"

Your assistant retrieves a list of world bosses you have defeated since the last daily reset. World bosses like Tequatl, Shadow Behemoth, and the Claw of Jormag each give bonus rewards once per day. Use this to decide whether it is worth chasing the next boss on the timer.

**Note:** The API returns internal identifiers for encounters (like `"forsaken_thicket"` for raids or `"claw_of_jormag"` for world bosses). Your AI assistant translates these into readable names automatically.

**Tip:** Use this before starting your play session to know what still has rewards available. Combine it with the [Wizard's Vault check](wizards-vault-daily/) to plan an efficient session that covers both your Vault objectives and your instanced content rewards.

## See also

- [Daily Checklist tutorial](../tutorials/daily-checklist/) -- full walkthrough of all daily tracking features
- [Tools reference](../reference/tools/) -- complete parameter details for `get_account_dailies`
- [API Scopes reference](../reference/api-scopes/) -- which permissions your API key needs
