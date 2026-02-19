---
title: Trading Post Mastery
weight: 2
---

# Trading Post Mastery

Turn your AI assistant into a Trading Post companion in about 30 minutes.

By the end of this tutorial, you will have:
- Looked up live buy and sell prices for any item by name
- Compared prices across multiple items in a single conversation
- Checked gem-to-gold and gold-to-gem exchange rates
- Reviewed your open buy orders and sell listings
- Checked your Trading Post delivery box for uncollected items and coins
- Browsed your recent transaction history

## Prerequisites

Before starting, make sure you have:
- Completed the [Getting Started](../getting-started/) tutorial (server running, AI assistant connected)
- An API key with the **tradingpost** scope enabled (the Getting Started tutorial has you enable all scopes, so you are covered if you followed it)

## What you will build

You will learn to use your AI assistant as a real-time Trading Post dashboard. Instead of tabbing into the game to check prices, review orders, or calculate gem exchange costs, you will ask your assistant in plain language and get answers immediately.

The first three sections (price lookups, comparisons, and gem exchange) work without an API key because they use public market data. Sections 4 through 6 (orders, deliveries, and transaction history) require your API key with the `tradingpost` scope.

---

## Section 1: Look up item prices by name

The most common Trading Post question is "what's this item worth right now?" Your assistant can answer that directly.

> Ask your AI: "What's the TP price for Mystic Coin?"

Your assistant will use the `get_tp_price_by_name` tool to search for Mystic Coin, find its item ID, and fetch the current Trading Post prices. You should see a response that includes:

- **Buy price** -- the highest price a buyer is currently offering. This is what you would receive if you sold instantly (minus the 15% listing fee and sales tax).
- **Sell price** -- the lowest price a seller is currently listing. This is what you would pay if you bought instantly.
- **Item name** -- confirming which item was matched.

Prices are displayed in the gold/silver/copper format you are familiar with from the game (for example, `2g 35s 10c`).

> **What is the spread?**
>
> The difference between the buy price and the sell price is called the spread. A large spread means there is a bigger gap between what buyers want to pay and what sellers are asking. Flippers profit by buying at the buy price and listing at the sell price, but remember the Trading Post takes a 15% cut (5% listing fee when you post, 10% sales tax when it sells). A flip is only profitable if the spread is greater than 15% of the sell price.

### Checkpoint

You should see both a buy price and a sell price for Mystic Coin, displayed in gold, silver, and copper. The buy price is always lower than the sell price. If the prices look reasonable compared to what you see in the in-game Trading Post, everything is working.

---

## Section 2: Compare multiple items

You are not limited to one item at a time. Your assistant can look up several items in the same conversation so you can compare them side by side.

> Ask your AI: "Compare the Trading Post prices for Mystic Coin and Glob of Ectoplasm"

Your assistant will call `get_tp_price_by_name` once for each item and present both results together. You should see buy and sell prices for both items, making it straightforward to compare their current market value.

Try adding more items to the comparison:

> Ask your AI: "What are the current TP prices for Mystic Coin, Glob of Ectoplasm, and Amalgamated Gemstone?"

Your assistant handles each lookup and summarizes all three. This is useful when you are deciding which materials to farm, which to buy, or which to sell from your material storage.

### Checkpoint

You should see buy and sell prices for each item you asked about. The items are identified by name, and each has its own buy and sell price pair. If your assistant reports that an item was not found, double-check the spelling -- it uses the exact item name from the GW2 Wiki.

---

## Section 3: Check gem exchange rates

The gem exchange lets you convert gold to gems or gems to gold. Your assistant can check the current rate in both directions.

**Gold to gems:**

> Ask your AI: "How much gold would it cost to buy 400 gems?"

Your assistant will use the `get_gem_exchange` tool with the `gems` direction and a quantity of `400`. You should see:

- How many coins (in gold/silver/copper) you would need to spend
- The effective rate per gem

**Gems to gold:**

> Ask your AI: "How many gems would I need to sell to get 100 gold?"

This time, the tool is called with the `coins` direction. Behind the scenes, the quantity is specified in copper (100 gold = 1,000,000 copper), but your assistant handles that conversion for you. You should see:

- How many gems you would need to exchange
- The rate per gem in gold/silver/copper

> **Why do the rates differ?**
>
> The gem exchange works like a market. The price you pay to buy gems is higher than what you receive when selling gems. This is similar to the buy/sell spread on the Trading Post -- ArenaNet takes a cut, and the rates fluctuate based on player activity. Gem prices tend to spike during gem store sales and new expansion releases.

### Checkpoint

You should see a clear answer for each direction: a gold cost when buying gems, and a gem count when converting gold. The rates reflect the live exchange and will change over time. If you see an error, try a different quantity -- very small amounts (under 1 gold or under 10 gems) may not return useful results.

---

## Section 4: Review your open orders

When you place buy orders or sell listings on the Trading Post, they sit there until someone fills them. Your assistant can show you all your active orders.

**Check buy orders:**

> Ask your AI: "Show my current buy orders"

Your assistant will use the `get_tp_transactions` tool with the type `current/buys`. For each open buy order, you should see:

- The **item name** you are trying to buy
- The **price** you offered per unit (in gold/silver/copper)
- The **quantity** you requested
- When you **created** the order

**Check sell listings:**

> Ask your AI: "Show my sell listings"

This time the tool uses the type `current/sells`. For each active sell listing, you should see:

- The **item name** you are selling
- The **listing price** per unit
- The **quantity** listed
- When you **created** the listing

If you have no active orders or listings, your assistant will tell you the list is empty. That is normal if you have not placed any orders recently.

### Checkpoint

You should see your current buy orders and sell listings, or a message confirming you have none. Compare the results to what you see on the Trading Post "My Transactions" tab in-game to verify they match. If you see an authentication error, make sure your API key has the `tradingpost` scope enabled.

---

## Section 5: Check your delivery box

When a buy order is filled or a sell listing is purchased, the items or coins go to your Trading Post delivery box. You need to visit a Trading Post NPC or use the Trading Post panel to pick them up.

> Ask your AI: "What's waiting at the Trading Post?"

Your assistant will use the `get_tp_delivery` tool. You should see:

- **Coins** awaiting pickup (if any of your sell listings have sold), displayed in gold/silver/copper
- **Items** awaiting pickup (if any of your buy orders have been filled), with item names and quantities

If both are empty, you have already collected everything. Time to place some new orders.

### Checkpoint

You should see your pending deliveries, or confirmation that your delivery box is empty. If you know you have uncollected gold or items from recent sales, verify that the amounts shown match what the in-game Trading Post panel reports. If you see an authentication error, confirm your API key includes the `tradingpost` scope.

---

## Section 6: Transaction history

Want to see what you have bought or sold recently? The Trading Post keeps a history of your completed transactions for the past 90 days.

> Ask your AI: "Show my recent sells"

Your assistant will use the `get_tp_transactions` tool with the type `history/sells`. For each completed sale, you should see:

- The **item name** that sold
- The **price** it sold at
- The **quantity** that sold
- When it was **created** (listed) and when it was **purchased** (sold)

You can also check your purchase history:

> Ask your AI: "Show my recent buys"

This uses the type `history/buys` and shows every buy order that was filled in the past 90 days.

> **Tip: Ask follow-up questions**
>
> Once your assistant has your transaction history, try asking questions about it: "How much gold did I spend on Mystic Coins this month?" or "What was my most expensive sale?" Your assistant can analyze the data it already retrieved.

### Checkpoint

You should see a list of your recent completed transactions with item names, prices, quantities, and dates. If you have been active on the Trading Post, you may see a long list. If you have not traded recently, the list may be short or empty. Both are expected.

---

## What you learned

In this tutorial, you used your AI assistant to:

- **Look up item prices by name** using `get_tp_price_by_name`, and learned the difference between buy price, sell price, and the spread
- **Compare prices for multiple items** by asking about several items in one conversation
- **Check gem exchange rates** in both directions using `get_gem_exchange`, and learned that the buy and sell rates differ
- **Review open orders** using `get_tp_transactions` with `current/buys` and `current/sells`
- **Check your delivery box** using `get_tp_delivery` to see uncollected items and coins
- **Browse transaction history** using `get_tp_transactions` with `history/buys` and `history/sells`

You now have a complete picture of how to monitor the Trading Post through your AI assistant, without needing to be logged into the game.

## Next steps

Now that you are comfortable with Trading Post queries, explore these guides:

- **Evaluate whether flipping an item is profitable**: See [Check Item Profitability](../how-to/check-item-profitability/)
- **Convert between gold and gems efficiently**: See [Gem Exchange](../how-to/gem-exchange/)
- **Browse all available tools**: See the [Tools Reference](../reference/tools/) for the full list of 34 tools
