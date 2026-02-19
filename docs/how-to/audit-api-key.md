---
title: Audit Your API Key
weight: 14
---

# How to Audit Your API Key

**Goal**: Check what permissions your GW2 API key has and understand which tools require which scopes, so you can fix authorization errors or set up a key with the right access.

## Prerequisites

- GW2 MCP Server running and connected to your AI assistant -- see [Getting Started](../tutorials/getting-started/) if you need setup help
- A GW2 API key configured (any scopes -- even a minimal key works for this check)

## Steps

### 1. Ask your AI what permissions your key has

> Ask your AI: "What permissions does my API key have?"

Your AI calls the `get_token_info` tool. This works with any valid API key, regardless of its scopes. The response includes:

- **Key name** -- the label you gave the key when you created it (for example, "GW2 MCP Key")
- **Permissions** -- the list of scopes enabled on this key (for example, `account`, `characters`, `inventories`)

### 2. Review your scopes

Compare the scopes your key has against what you need. The following table shows which scopes unlock which tool categories:

| Scope | What it enables |
|-------|----------------|
| `account` | Account info, wallet, bank, materials, inventory, characters, unlocks, progress, dailies. Required by all authenticated tools. |
| `characters` | Character details including equipment, builds, and crafting disciplines. |
| `inventories` | Bank vault, material storage, and shared inventory slots. |
| `tradingpost` | Trading Post current orders, delivery box, and transaction history. |
| `wallet` | Currency balances (gold, karma, gems, laurels, and all other currencies). |
| `unlocks` | Skins, dyes, minis, titles, recipes, mount skins, mount types, and other wardrobe unlocks. |
| `progression` | Achievement progress, masteries, luck, daily completions, and Wizard's Vault objectives. |
| `guilds` | Guild details (requires guild leader permissions on the account). |

If your key is missing a scope, any tool that requires it will return an authorization error.

### 3. Fix missing scopes if needed

If you need scopes that your current key does not have:

1. Go to [Guild Wars 2 API Key Management](https://account.arena.net/applications)
2. Create a new key with the scopes you need (or enable all scopes for full access)
3. Copy the new key
4. Update the `GW2_API_KEY` value in your MCP client configuration
5. Restart the GW2 MCP Server

You cannot edit an existing key's scopes -- you must create a new one.

### 4. Confirm the updated key works

After updating your key, verify the change took effect:

> Ask your AI: "Check my API key permissions again"

Your AI calls `get_token_info` with the new key and should now show the additional scopes.

For the full scope-to-tool mapping, see the [API Scopes reference](../reference/api-scopes/).

## Troubleshooting

### Problem: "GW2_API_KEY environment variable not configured"
**Symptom**: The tool returns an error saying no API key is set.
**Cause**: The `GW2_API_KEY` environment variable is not reaching the server process.
**Solution**: Check your MCP client configuration. The `env` block must contain `GW2_API_KEY` with your key value. See [Configure MCP Clients](configure-mcp-clients/) for setup details.

### Problem: "Invalid API key" error
**Symptom**: `get_token_info` fails with an authentication error instead of returning key details.
**Cause**: The API key is malformed, expired, or was deleted from your ArenaNet account.
**Solution**: Go to [Guild Wars 2 API Key Management](https://account.arena.net/applications) and confirm the key still exists. If not, create a new one and update your configuration.

### Problem: Key has scopes but tools still fail
**Symptom**: `get_token_info` shows the right scopes, but certain tools return authorization errors.
**Cause**: Some tools require multiple scopes. For example, `get_bank` requires both `account` and `inventories`. If your key has `inventories` but not `account`, it will fail.
**Solution**: Ensure your key includes the `account` scope -- it is required by every authenticated tool. The safest approach is to enable all scopes when creating your key.

## See also

- [API Scopes reference](../reference/api-scopes/) -- complete scope-to-tool mapping
- [How to use tools without an API key](no-api-key/) -- tools that work without authentication
