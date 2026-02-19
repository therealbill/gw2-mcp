---
title: Gem Exchange Rates
weight: 12
---

# How to Check Gem Exchange Rates

**Goal**: Look up current gem-to-gold and gold-to-gem conversion rates.

**Time**: Approximately 3 minutes

## Prerequisites

None -- gem exchange rates are public data and do not require an API key.

## Steps

### 1. Check how many gems you can buy with gold

> Ask your AI: "How many gems can I get for 100 gold?"

Your assistant uses the `get_gem_exchange` tool with direction `coins` and the quantity in copper. Your assistant handles the conversion automatically -- you just ask in gold.

Behind the scenes, the API works in copper: 1 gold = 10,000 copper, so 100 gold = 1,000,000 copper. Your assistant converts for you and returns the number of gems you would receive.

### 2. Check how much gold you get for gems

> Ask your AI: "How much gold would 400 gems give me?"

This uses `get_gem_exchange` with direction `gems` and a quantity of 400. The response tells you how much gold (in gold/silver/copper format) you would receive for exchanging those gems.

400 gems is a common amount to check because many gem store items cost 400 or 800 gems.

### 3. Compare rates for planning purchases

Gem exchange rates change dynamically based on player activity. Rates tend to rise during gem store sales and new releases when more players are buying gems with gold.

> Ask your AI: "How much gold would it cost to buy 800 gems right now?"

Use this to decide whether to convert gold to gems now or wait for a better rate. There is no historical rate data available through the API, but you can check periodically and compare.

## Verify it works

Cross-check the rate your assistant reports against the in-game gem exchange panel:

1. Open the gem exchange in-game (accessible from the gem store or Trading Post)
2. Enter the same gold or gem amount
3. Compare the conversion result

The numbers should be very close. Minor differences can occur because rates fluctuate continuously.

## Troubleshooting

### Problem: Unexpected result for small amounts
**Symptom**: Asking about very small amounts (under 1 gold or under 10 gems) returns unusual results.
**Cause**: The gem exchange has minimum thresholds and the rate calculation can produce odd results at very low quantities.
**Solution**: Use realistic amounts -- 10 gold or more for coins-to-gems, and 100 gems or more for gems-to-coins.

### Problem: Rate seems very different from yesterday
**Symptom**: The exchange rate changed significantly since you last checked.
**Cause**: Gem prices are driven by player supply and demand. Major game events, gem store sales, and expansion launches cause large rate swings.
**Solution**: This is expected behavior. The rate you see is the live market rate at the moment you ask.

### Problem: "coins" vs "gems" confusion
**Symptom**: You get a result that does not make sense for your question.
**Cause**: The two directions can be confusing. "Coins" means you are spending coins (gold) to buy gems. "Gems" means you are spending gems to get coins (gold).
**Solution**: Phrase your question clearly. "How many gems for 100 gold?" uses the coins direction. "How much gold for 400 gems?" uses the gems direction. Your assistant interprets natural language, so a clear question produces the right result.

## See also

- [Trading Post Mastery](../tutorials/trading-post/) -- full tutorial including gem exchange basics
- [Check Item Profitability](check-item-profitability/) -- evaluate TP flips
- [Tools reference](../reference/tools/#get_gem_exchange) -- `get_gem_exchange` specification
