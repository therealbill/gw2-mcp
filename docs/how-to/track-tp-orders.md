---
title: Track Your TP Orders
weight: 11
---

# How to Track Your TP Orders

**Goal**: View your current buy orders, sell listings, and recent transaction history on the Trading Post.

**Time**: Approximately 5 minutes

## Prerequisites

- GW2 MCP Server connected to your AI assistant -- see [Getting Started](../tutorials/getting-started/)
- An API key with the **tradingpost** scope enabled -- see [API Key Scopes](../reference/api-scopes/)

## Steps

### 1. Check your current buy orders

> Ask your AI: "Show my current buy orders"

Your assistant uses the `get_tp_transactions` tool with type `current/buys`. For each open buy order, you will see:

- The **item name** you are trying to buy
- The **price** you offered per unit
- The **quantity** you requested
- When the order was **created**

These are orders you have placed that have not yet been filled. They will remain active until someone sells to you at your price, or you cancel them in-game.

### 2. Check your current sell listings

> Ask your AI: "Show my sell listings"

This uses `get_tp_transactions` with type `current/sells`. For each active listing, you will see:

- The **item name** you listed for sale
- Your **listing price** per unit
- The **quantity** listed
- When you **created** the listing

These are items sitting on the Trading Post waiting for a buyer. Remember, the 5% listing fee was already deducted when you posted them.

### 3. Review recent transaction history

Check what has sold or been purchased recently:

> Ask your AI: "Show my recent completed sells"

This uses `get_tp_transactions` with type `history/sells` and shows every sale that completed. You will see the item name, sale price, quantity, and the date it sold.

For purchases:

> Ask your AI: "Show my recent completed buys"

This uses type `history/buys` and shows every buy order that was filled.

Transaction history covers the **past 90 days**. Older transactions are no longer available through the API.

## Verify it works

Compare your results against the in-game Trading Post:

1. Open the Trading Post panel in-game and go to the "My Transactions" tab
2. Check that your open buy orders match what your assistant reported
3. Check that your sell listings match as well

The numbers should align. If you placed or cancelled orders very recently, give it a moment for the API to update.

## Troubleshooting

### Problem: "GW2_API_KEY environment variable not configured"
**Symptom**: Your assistant returns this error when you ask about orders.
**Cause**: The API key is not set in your MCP server configuration.
**Solution**: Add your API key to the server configuration. See [Configure MCP Clients](configure-mcp-clients/) for setup instructions.

### Problem: Authorization or permissions error
**Symptom**: The tool returns a permissions error from the GW2 API.
**Cause**: Your API key does not have the `tradingpost` scope enabled.
**Solution**: Go to [Guild Wars 2 API Key Management](https://account.arena.net/applications), edit or recreate your key, and enable the **tradingpost** permission. Then update the key in your MCP server configuration.

### Problem: Empty results but you have active orders
**Symptom**: Your assistant says you have no orders, but you know you do.
**Cause**: The API key may belong to a different account, or the orders were placed on a different account.
**Solution**: Verify the API key belongs to the correct account by asking your assistant to run `get_token_info` or `get_account`. Check that the account name matches.

## See also

- [Trading Post Mastery](../tutorials/trading-post/) -- full tutorial covering all TP tools
- [Check Item Profitability](check-item-profitability/) -- evaluate whether a flip is worth it
- [API Key Scopes](../reference/api-scopes/) -- required permissions for each tool
- [Tools reference](../reference/tools/#get_tp_transactions) -- `get_tp_transactions` specification
