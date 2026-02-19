---
title: Getting Started
---

# Getting Started

Connect Claude Desktop to your Guild Wars 2 account in about 15 minutes.

By the end of this tutorial, you will have:
- A running GW2 MCP server inside Docker
- Claude Desktop configured to use it
- Queried your GW2 wallet, looked up an item, and checked Trading Post prices -- all through natural language

## Prerequisites

Before starting, make sure you have:
- **Docker Desktop** installed and running ([download here](https://www.docker.com/products/docker-desktop/))
- **Claude Desktop** installed ([download here](https://claude.ai/download))
- A **Guild Wars 2** account

## Step 1: Pull the server image

Open a terminal and pull the GW2 MCP server Docker image:

```bash
docker pull alyxpink/gw2-mcp:v1
```

This downloads the pre-built server image. You should see output ending with:

```
Status: Downloaded newer image for alyxpink/gw2-mcp:v1
```

### Checkpoint

Verify the image is available:

```bash
docker images alyxpink/gw2-mcp
```

You should see a row listing `alyxpink/gw2-mcp` with the tag `v1`. If you see it, the image is ready.

## Step 2: Create a GW2 API key

Open your browser and go to [https://account.arena.net/applications](https://account.arena.net/applications). Log in with your Guild Wars 2 account.

1. Click **New Key**.
2. Give the key a name (for example, `mcp-server`).
3. Check **all** permission boxes: account, characters, inventories, progression, guilds, tradingpost, unlocks, wallet.
4. Click **Create API Key**.
5. Copy the generated key. It is a long string of letters, numbers, and dashes.

Save this key somewhere safe. You will use it in the next step.

### Checkpoint

You have an API key string that looks similar to this format:

```
XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXXXXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
```

If the key management page did not load, make sure you are logged in to the correct Guild Wars 2 account. For details on what each scope enables, see the [API Scopes reference](../reference/api-scopes/).

## Step 3: Configure Claude Desktop

Open Claude Desktop. Navigate to the MCP server configuration:

1. Open **Settings** (click the gear icon or use the menu).
2. Go to the **Developer** section.
3. Click **Edit Config**. This opens a JSON file in your text editor.

Replace the file contents with the following JSON, substituting your actual API key for `YOUR_KEY_HERE`:

```json
{
  "mcpServers": {
    "gw2-mcp": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "-e", "GW2_API_KEY", "alyxpink/gw2-mcp:v1"],
      "env": {
        "GW2_API_KEY": "YOUR_KEY_HERE"
      }
    }
  }
}
```

If the file already has other MCP servers configured, add the `"gw2-mcp"` block inside the existing `"mcpServers"` object rather than replacing the entire file.

Save the file and **restart Claude Desktop** (quit and reopen it).

### Checkpoint

After restarting, open a new conversation in Claude Desktop. Look for the MCP tools icon (a hammer or plug icon near the text input). Click it and confirm that GW2 tools appear in the list -- you should see tools like `get_wallet`, `get_items`, and `wiki_search`.

If the tools do not appear, double-check that:
- Docker Desktop is running
- The JSON is valid (no trailing commas, correct bracket matching)
- You saved the config file before restarting Claude Desktop

## Step 4: Make your first query

In Claude Desktop, start a new conversation and type:

```
Check my GW2 wallet
```

Claude will call the `get_wallet` tool through the MCP server. After a moment, you should see a response listing your in-game currencies -- gold, karma, gems, laurels, and more -- with their current balances.

### Checkpoint

You see your wallet data displayed in the conversation. The response includes currency names and amounts that match what you would see in-game. If you see an error about `GW2_API_KEY`, go back to Step 3 and verify that you replaced `YOUR_KEY_HERE` with your actual API key.

## Step 5: Look up game data

Now try looking up an item. Type:

```
Look up the item Mystic Coin
```

Claude will use the `get_item_by_name` tool to search for the item and return its details. You should see information including the item name, rarity, type, level, description, and vendor value.

### Checkpoint

You see detailed item information for Mystic Coin, including that it is an Exotic rarity currency item. This lookup works without an API key -- it queries public game data.

## Step 6: Check Trading Post prices

Type:

```
What's the Trading Post price for Glob of Ectoplasm?
```

Claude will look up the item and retrieve current Trading Post prices using the `get_tp_price_by_name` tool. You should see buy and sell prices displayed in gold, silver, and copper.

### Checkpoint

You see current buy and sell prices for Glob of Ectoplasm. The prices reflect the live Trading Post market. Compare them to the in-game Trading Post to confirm they match.

## What you accomplished

You now have a working setup where Claude Desktop can access Guild Wars 2 data on your behalf:

- **Docker** runs the GW2 MCP server as a container, with no build tools needed
- **Claude Desktop** launches the container automatically when you start a conversation
- **Your API key** lets the server access your personal account data (wallet, characters, bank)
- **Public tools** like item lookups and Trading Post prices work for any game data

The server provides 34 tools in total, covering account data, the Trading Post, achievements, guilds, the Wizard's Vault, and general game metadata.

## Next steps

Now that your setup is working, explore further:

- **Set up other MCP clients** -- See the [Configure MCP Clients](../how-to/configure-mcp-clients/) guide for Claude Code, LM Studio, and other editors
- **Browse all 34 tools** -- See the [Tools reference](../reference/tools/) for the complete list of available tools
- **Understand API key permissions** -- See the [API Scopes reference](../reference/api-scopes/) for which scopes each tool requires

Try asking Claude about your characters, bank contents, achievement progress, or the current Wizard's Vault objectives. The server handles caching automatically, so repeated queries are fast.
