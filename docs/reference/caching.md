---
title: "Caching"
---

The GW2 MCP Server maintains an in-memory cache to reduce redundant calls to the Guild Wars 2 API. Each data category has a fixed time-to-live (TTL) after which the cached entry expires and the next request fetches fresh data from the API.

## Cache TTL Values

All TTL constants are defined in `internal/cache/manager.go`.

### Static Game Data

| Constant | TTL | Applies To |
|----------|-----|------------|
| `StaticDataTTL` | 365 days | Currency definitions |
| `ItemDataTTL` | 24 hours | Item metadata, skin metadata |
| `RecipeDataTTL` | 24 hours | Recipe details, recipe search results |
| `AchievementDataTTL` | 24 hours | Achievement details |
| `ColorDataTTL` | 24 hours | Dye color definitions |
| `MiniDataTTL` | 24 hours | Miniature definitions |
| `MountDataTTL` | 24 hours | Mount skin and type definitions (defined but not currently used in client) |
| `DungeonDataTTL` | 24 hours | Dungeon and raid definitions (defined but not currently used in client) |
| `WikiDataTTL` | 24 hours | Wiki search results, wiki page content |

### Account Data

All account data keys include a SHA-256 hash of the API key for per-key isolation.

| Constant | TTL | Applies To |
|----------|-----|------------|
| `AccountDataTTL` | 5 minutes | Account info, bank contents, material storage, shared inventory, character list, character details |
| `WalletDataTTL` | 5 minutes | Wallet balances |
| `ProgressTTL` | 5 minutes | Account progress (achievements, masteries, mastery points, luck, legendary armory, progression) |
| `UnlocksTTL` | 10 minutes | Account unlocks (skins, dyes, minis, titles, recipes, finishers, outfits, gliders, mail carriers, novelties, emotes, mounts, skiffs, jade bots) |
| `DailiesTTL` | 2 minutes | Completed dailies (daily crafting, dungeons, raids, map chests, world bosses) |

### Trading Post

| Constant | TTL | Applies To |
|----------|-----|------------|
| `TPPriceTTL` | 5 minutes | Item buy/sell prices |
| `TPListingTTL` | 5 minutes | Item order book listings |
| `TPTransactionTTL` | 5 minutes | Transaction history (current buys, current sells, history buys, history sells) |
| `TPExchangeTTL` | 10 minutes | Gem-to-coin and coin-to-gem exchange rates |
| `TPDeliveryTTL` | 2 minutes | Trading Post delivery box contents |

### Guild

| Constant | TTL | Applies To |
|----------|-----|------------|
| `GuildInfoTTL` | 1 hour | Public guild info (name, tag, level) |
| `GuildSearchTTL` | 1 hour | Guild name search results |
| `GuildDetailTTL` | 5 minutes | Guild detail data (log, members, ranks, stash, storage, treasury, teams, upgrades) |

### Wizard's Vault

| Constant | TTL | Applies To |
|----------|-----|------------|
| `WVSeasonTTL` | 24 hours | Current season information |
| `WVObjectivesAuthTTL` | 5 minutes | Objectives (authenticated, per-account) |
| `WVObjectivesPublicTTL` | 1 hour | Objectives (public, unauthenticated) |
| `WVListingsTTL` | 1 hour | Reward listings (authenticated and public) |

### Game Metadata

| Constant | TTL | Applies To |
|----------|-----|------------|
| `DailyAchievementTTL` | 1 hour | Today's and tomorrow's daily achievements |
| `GameBuildTTL` | 1 hour | Current game build number |
| `TokenInfoTTL` | 10 minutes | API key token info (name, permissions) |

## Cache Behavior

- **Storage**: In-memory only, using `github.com/patrickmn/go-cache`.
- **Persistence**: The cache is not persisted to disk. All cached data is lost when the process exits.
- **Cleanup interval**: Expired entries are purged every 10 minutes (`CleanupInterval`).
- **Default TTL**: The underlying cache instance is created with `StaticDataTTL` (365 days) as the default expiration; individual entries override this with their specific TTL at write time.
- **Per-key isolation**: Account-specific data (wallet, bank, materials, inventory, characters, unlocks, progress, dailies, trading post delivery, trading post transactions, Wizard's Vault objectives, Wizard's Vault listings, token info) is keyed by a SHA-256 hash of the API key. Different API keys produce separate cache entries.
- **Public data sharing**: Non-authenticated data (items, skins, currencies, recipes, achievements, colors, minis, wiki content, guild info, game build, daily achievements, trading post prices, trading post listings, gem exchange rates) is shared across all users.

## See Also

- [Architecture](../explanation/architecture/) -- cache layer design
- [Design Decisions](../explanation/design-decisions/) -- rationale for TTL values
