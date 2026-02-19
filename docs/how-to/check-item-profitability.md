---
title: Check Item Profitability
weight: 10
---

# How to Check Item Profitability

**Goal**: Determine whether flipping an item on the Trading Post is profitable after fees.

**Time**: Approximately 5 minutes

## Prerequisites

None -- Trading Post prices are public data and do not require an API key.

You should be familiar with looking up item prices. If not, see the [Trading Post Mastery](../tutorials/trading-post/) tutorial first.

## Steps

### 1. Look up the item's buy and sell prices

> Ask your AI: "What are the Trading Post prices for Mystic Coin?"

Your assistant uses the `get_tp_price_by_name` tool to fetch live prices. You will see two numbers:

- **Buy price** -- the highest price a buyer is currently offering (what you pay if you place a buy order and wait)
- **Sell price** -- the lowest price a seller is currently listing (what buyers pay for instant purchase)

To flip an item, you buy at the buy price and list at the sell price.

### 2. Understand the Trading Post tax

The Trading Post takes a 15% cut on every sale, split into two fees:

| Fee | Rate | When charged |
|-----|------|--------------|
| Listing fee | 5% of your listing price | When you post the sell listing |
| Exchange fee | 10% of your listing price | When the item sells |

The listing fee is non-refundable -- you pay it even if you cancel the listing. Both fees are deducted automatically.

### 3. Calculate your profit

Use this formula:

```
Profit = (Sell Price x 0.85) - Buy Price
```

Multiply the sell price by 0.85 to account for the 15% total fee. Subtract the buy price (your cost). If the result is positive, the flip is profitable.

**Example with Mystic Coin:**

Suppose the current prices are:

- Buy price: 1g 50s 0c
- Sell price: 1g 85s 0c

The calculation:

```
Revenue after fees = 1g 85s x 0.85 = 1g 57s 25c
Profit = 1g 57s 25c - 1g 50s 0c = 7s 25c per coin
```

That is a profit of 7 silver 25 copper per Mystic Coin, or about a 4.8% return on investment.

### 4. Ask your AI to do the math

You do not need to calculate manually. Ask your assistant directly:

> Ask your AI: "Is Mystic Coin profitable to flip on the TP? Show the math with the 15% fee."

Your assistant will fetch the current prices, apply the 15% fee, and tell you the profit (or loss) per unit.

Try other items:

> Ask your AI: "What's the profit margin on Glob of Ectoplasm?"

> Ask your AI: "Compare the flip profit on Mystic Coin vs Amalgamated Gemstone"

## Verify it works

Pick any item and run the numbers yourself against what your assistant reports:

1. Note the buy and sell prices your assistant shows
2. Multiply the sell price by 0.85
3. Subtract the buy price
4. Confirm the result matches what your assistant calculated

If the numbers align, you are reading profitability correctly.

## Troubleshooting

### Problem: Prices look different from in-game
**Symptom**: The buy/sell prices your assistant shows do not match what you see in the Trading Post panel.
**Cause**: Prices change constantly as players post and cancel orders. There is a short delay between the API and the in-game display.
**Solution**: This is normal. Use the API prices as a close approximation. For high-volume items like Mystic Coin, prices are very close to real-time.

### Problem: Profit looks good but you lose money in practice
**Symptom**: You flip an item expecting profit but end up with less gold.
**Cause**: Prices shifted between when you placed your buy order and when your sell listing was purchased. High-volume items are safer because their prices move less.
**Solution**: Focus on items with consistent spreads. Avoid items where the price swings wildly within hours.

### Problem: Item not found
**Symptom**: Your assistant says it cannot find the item.
**Cause**: The item name does not match the wiki exactly, or the item is not tradeable on the Trading Post.
**Solution**: Use the exact in-game item name. Account-bound, soulbound, and other untradeable items do not appear on the Trading Post.

## See also

- [Trading Post Mastery](../tutorials/trading-post/) -- full tutorial on using TP tools
- [Track Your TP Orders](track-tp-orders/) -- monitor your active buy and sell orders
- [Tools reference](../reference/tools/#get_tp_price_by_name) -- `get_tp_price_by_name` specification
