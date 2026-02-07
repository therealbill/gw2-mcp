// Package cache provides caching functionality for the GW2 MCP server.
package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

// Manager handles caching for the GW2 MCP server
type Manager struct {
	cache *cache.Cache
}

// Key represents different types of cache keys
type Key string

const (
	// CurrencyListKey is the cache key for the list of all currencies
	CurrencyListKey Key = "currencies:list"
	// CurrencyDetailKey is the cache key template for individual currency details
	CurrencyDetailKey Key = "currency:detail:%d"
	// WikiSearchKey is the cache key template for wiki search results
	WikiSearchKey Key = "wiki:search:%s"
	// WikiPageKey is the cache key template for wiki page content
	WikiPageKey Key = "wiki:page:%s"

	// WalletKey is the cache key template for wallet data (short TTL)
	WalletKey Key = "wallet:%s" // %s = hashed API key

	// ItemDetailKey is the cache key template for individual item details
	ItemDetailKey Key = "item:detail:%d" // %d = item ID

	// Trading Post cache keys
	TPPriceKey       Key = "tp:price:%d"          // %d = item ID
	TPListingKey     Key = "tp:listing:%d"         // %d = item ID
	TPExchangeKey    Key = "tp:exchange:%s:%d"     // %s = direction, %d = quantity
	TPDeliveryKey    Key = "tp:delivery:%s"        // %s = hashed API key
	TPTransactionKey Key = "tp:transactions:%s:%s" // %s = hashed API key, %s = type

	// Account cache keys
	AccountKey         Key = "account:%s"            // %s = hashed API key
	BankKey            Key = "bank:%s"               // %s = hashed API key
	MaterialsKey       Key = "materials:%s"          // %s = hashed API key
	SharedInventoryKey Key = "inventory:%s"          // %s = hashed API key
	CharactersKey      Key = "characters:%s"         // %s = hashed API key
	CharacterKey       Key = "character:%s:%s"       // %s = hashed API key, %s = name
	UnlocksKey         Key = "unlocks:%s:%s"         // %s = hashed API key, %s = type
	ProgressKey        Key = "progress:%s:%s"        // %s = hashed API key, %s = type
	DailiesKey         Key = "dailies:%s:%s"         // %s = hashed API key, %s = type

	// Wizard's Vault cache keys
	WizardsVaultSeasonKey     Key = "wv:season"           // no params
	WizardsVaultObjectivesKey Key = "wv:objectives:%s:%s"  // %s = hashed API key (or "public"), %s = type
	WizardsVaultListingsKey   Key = "wv:listings:%s"       // %s = hashed API key (or "public")

	// Game data cache keys
	SkinDetailKey      Key = "skin:detail:%d"       // %d = skin ID
	RecipeDetailKey    Key = "recipe:detail:%d"      // %d = recipe ID
	RecipeSearchKey    Key = "recipe:search:%s:%d"   // %s = direction (input/output), %d = item ID
	AchievementKey     Key = "achievement:detail:%d" // %d = achievement ID
	DailyAchievementKey Key = "achievements:daily"

	// Guild cache keys
	GuildInfoKey    Key = "guild:info:%s"       // %s = guild ID
	GuildSearchKey  Key = "guild:search:%s"     // %s = guild name
	GuildDetailKey  Key = "guild:detail:%s:%s"  // %s = guild ID, %s = type

	// Metadata cache keys
	ColorDetailKey     Key = "color:detail:%d"      // %d = color ID
	MiniDetailKey      Key = "mini:detail:%d"       // %d = mini ID
	MountDetailKey     Key = "mount:%s:detail:%d"   // %s = type (skins/types), %d = mount ID
	GameBuildKey       Key = "game:build"
	TokenInfoKey       Key = "tokeninfo:%s"         // %s = hashed API key
	DungeonDetailKey   Key = "dungeon:detail:%s"    // %s = dungeon/raid ID
)

// Cache durations
const (
	// Static data - cache for very long periods
	StaticDataTTL = 24 * time.Hour * 365 // 1 year for currencies
	WikiDataTTL   = 24 * time.Hour       // 1 day for wiki content

	// Dynamic data - shorter cache periods
	WalletDataTTL = 5 * time.Minute // 5 minutes for wallet data

	// Item metadata - semi-static
	ItemDataTTL = 24 * time.Hour // 1 day for item metadata

	// Trading Post data - dynamic
	TPPriceTTL       = 5 * time.Minute  // Prices change frequently
	TPListingTTL     = 5 * time.Minute  // Listings change frequently
	TPExchangeTTL    = 10 * time.Minute // Exchange rates less volatile
	TPDeliveryTTL    = 2 * time.Minute  // Users want fresh delivery info
	TPTransactionTTL = 5 * time.Minute  // Matches API-side cache

	// Account data
	AccountDataTTL = 5 * time.Minute
	UnlocksTTL     = 10 * time.Minute
	ProgressTTL    = 5 * time.Minute
	DailiesTTL     = 2 * time.Minute

	// Wizard's Vault
	WVSeasonTTL          = 24 * time.Hour
	WVObjectivesAuthTTL  = 5 * time.Minute
	WVObjectivesPublicTTL = 1 * time.Hour
	WVListingsTTL        = 1 * time.Hour

	// Game data
	RecipeDataTTL         = 24 * time.Hour
	AchievementDataTTL    = 24 * time.Hour
	DailyAchievementTTL   = 1 * time.Hour

	// Guild
	GuildInfoTTL   = 1 * time.Hour
	GuildSearchTTL = 1 * time.Hour
	GuildDetailTTL = 5 * time.Minute

	// Metadata
	ColorDataTTL    = 24 * time.Hour
	MiniDataTTL     = 24 * time.Hour
	MountDataTTL    = 24 * time.Hour
	GameBuildTTL    = 1 * time.Hour
	TokenInfoTTL    = 10 * time.Minute
	DungeonDataTTL  = 24 * time.Hour

	// Default cleanup interval
	CleanupInterval = 10 * time.Minute
)

// NewManager creates a new cache manager
func NewManager() *Manager {
	return &Manager{
		cache: cache.New(StaticDataTTL, CleanupInterval),
	}
}

// Set stores a value in the cache with the specified TTL
func (m *Manager) Set(key string, value interface{}, ttl time.Duration) {
	m.cache.Set(key, value, ttl)
}

// Get retrieves a value from the cache
func (m *Manager) Get(key string) (interface{}, bool) {
	return m.cache.Get(key)
}

// GetString retrieves a string value from the cache
func (m *Manager) GetString(key string) (string, bool) {
	if value, found := m.cache.Get(key); found {
		if str, ok := value.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetJSON retrieves and unmarshals a JSON value from the cache
func (m *Manager) GetJSON(key string, dest interface{}) bool {
	if value, found := m.cache.Get(key); found {
		if jsonStr, ok := value.(string); ok {
			if err := json.Unmarshal([]byte(jsonStr), dest); err == nil {
				return true
			}
		}
	}
	return false
}

// SetJSON marshals and stores a JSON value in the cache
func (m *Manager) SetJSON(key string, value interface{}, ttl time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}
	m.cache.Set(key, string(jsonData), ttl)
	return nil
}

// Delete removes a value from the cache
func (m *Manager) Delete(key string) {
	m.cache.Delete(key)
}

// Flush clears all cached data
func (m *Manager) Flush() {
	m.cache.Flush()
}

// ItemCount returns the number of items in the cache
func (m *Manager) ItemCount() int {
	return m.cache.ItemCount()
}

// GetCurrencyListKey returns the cache key for currency list
func (m *Manager) GetCurrencyListKey() string {
	return string(CurrencyListKey)
}

// GetCurrencyDetailKey returns the cache key for a specific currency
func (m *Manager) GetCurrencyDetailKey(id int) string {
	return fmt.Sprintf(string(CurrencyDetailKey), id)
}

// GetWikiSearchKey returns the cache key for wiki search results
func (m *Manager) GetWikiSearchKey(query string) string {
	return fmt.Sprintf(string(WikiSearchKey), query)
}

// GetWikiPageKey returns the cache key for a wiki page
func (m *Manager) GetWikiPageKey(title string) string {
	return fmt.Sprintf(string(WikiPageKey), title)
}

// GetWalletKey returns the cache key for wallet data
func (m *Manager) GetWalletKey(apiKeyHash string) string {
	return fmt.Sprintf(string(WalletKey), apiKeyHash)
}

// GetItemDetailKey returns the cache key for item metadata
func (m *Manager) GetItemDetailKey(id int) string {
	return fmt.Sprintf(string(ItemDetailKey), id)
}

// GetTPPriceKey returns the cache key for TP price data
func (m *Manager) GetTPPriceKey(itemID int) string {
	return fmt.Sprintf(string(TPPriceKey), itemID)
}

// GetTPListingKey returns the cache key for TP listing data
func (m *Manager) GetTPListingKey(itemID int) string {
	return fmt.Sprintf(string(TPListingKey), itemID)
}

// GetTPExchangeKey returns the cache key for gem exchange rates
func (m *Manager) GetTPExchangeKey(direction string, quantity int) string {
	return fmt.Sprintf(string(TPExchangeKey), direction, quantity)
}

// GetTPDeliveryKey returns the cache key for delivery box
func (m *Manager) GetTPDeliveryKey(apiKeyHash string) string {
	return fmt.Sprintf(string(TPDeliveryKey), apiKeyHash)
}

// GetTPTransactionKey returns the cache key for transaction history
func (m *Manager) GetTPTransactionKey(apiKeyHash string, txType string) string {
	return fmt.Sprintf(string(TPTransactionKey), apiKeyHash, txType)
}

// GetAccountKey returns the cache key for account data
func (m *Manager) GetAccountKey(apiKeyHash string) string {
	return fmt.Sprintf(string(AccountKey), apiKeyHash)
}

// GetBankKey returns the cache key for bank data
func (m *Manager) GetBankKey(apiKeyHash string) string {
	return fmt.Sprintf(string(BankKey), apiKeyHash)
}

// GetMaterialsKey returns the cache key for material storage data
func (m *Manager) GetMaterialsKey(apiKeyHash string) string {
	return fmt.Sprintf(string(MaterialsKey), apiKeyHash)
}

// GetSharedInventoryKey returns the cache key for shared inventory data
func (m *Manager) GetSharedInventoryKey(apiKeyHash string) string {
	return fmt.Sprintf(string(SharedInventoryKey), apiKeyHash)
}

// GetCharactersKey returns the cache key for characters list
func (m *Manager) GetCharactersKey(apiKeyHash string) string {
	return fmt.Sprintf(string(CharactersKey), apiKeyHash)
}

// GetCharacterKey returns the cache key for a specific character
func (m *Manager) GetCharacterKey(apiKeyHash string, name string) string {
	return fmt.Sprintf(string(CharacterKey), apiKeyHash, name)
}

// GetUnlocksKey returns the cache key for account unlocks
func (m *Manager) GetUnlocksKey(apiKeyHash string, unlockType string) string {
	return fmt.Sprintf(string(UnlocksKey), apiKeyHash, unlockType)
}

// GetProgressKey returns the cache key for account progress
func (m *Manager) GetProgressKey(apiKeyHash string, progressType string) string {
	return fmt.Sprintf(string(ProgressKey), apiKeyHash, progressType)
}

// GetDailiesKey returns the cache key for account dailies
func (m *Manager) GetDailiesKey(apiKeyHash string, dailyType string) string {
	return fmt.Sprintf(string(DailiesKey), apiKeyHash, dailyType)
}

// GetWizardsVaultSeasonKey returns the cache key for wizard's vault season info
func (m *Manager) GetWizardsVaultSeasonKey() string {
	return string(WizardsVaultSeasonKey)
}

// GetWizardsVaultObjectivesKey returns the cache key for wizard's vault objectives
func (m *Manager) GetWizardsVaultObjectivesKey(apiKeyHash string, objType string) string {
	return fmt.Sprintf(string(WizardsVaultObjectivesKey), apiKeyHash, objType)
}

// GetWizardsVaultListingsKey returns the cache key for wizard's vault listings
func (m *Manager) GetWizardsVaultListingsKey(apiKeyHash string) string {
	return fmt.Sprintf(string(WizardsVaultListingsKey), apiKeyHash)
}

// GetSkinDetailKey returns the cache key for skin metadata
func (m *Manager) GetSkinDetailKey(id int) string {
	return fmt.Sprintf(string(SkinDetailKey), id)
}

// GetRecipeDetailKey returns the cache key for recipe metadata
func (m *Manager) GetRecipeDetailKey(id int) string {
	return fmt.Sprintf(string(RecipeDetailKey), id)
}

// GetRecipeSearchKey returns the cache key for recipe search
func (m *Manager) GetRecipeSearchKey(direction string, itemID int) string {
	return fmt.Sprintf(string(RecipeSearchKey), direction, itemID)
}

// GetAchievementKey returns the cache key for achievement metadata
func (m *Manager) GetAchievementKey(id int) string {
	return fmt.Sprintf(string(AchievementKey), id)
}

// GetDailyAchievementKey returns the cache key for daily achievements
func (m *Manager) GetDailyAchievementKey() string {
	return string(DailyAchievementKey)
}

// GetGuildInfoKey returns the cache key for guild info
func (m *Manager) GetGuildInfoKey(guildID string) string {
	return fmt.Sprintf(string(GuildInfoKey), guildID)
}

// GetGuildSearchKey returns the cache key for guild search
func (m *Manager) GetGuildSearchKey(name string) string {
	return fmt.Sprintf(string(GuildSearchKey), name)
}

// GetGuildDetailKey returns the cache key for guild details
func (m *Manager) GetGuildDetailKey(guildID string, detailType string) string {
	return fmt.Sprintf(string(GuildDetailKey), guildID, detailType)
}

// GetColorDetailKey returns the cache key for color metadata
func (m *Manager) GetColorDetailKey(id int) string {
	return fmt.Sprintf(string(ColorDetailKey), id)
}

// GetMiniDetailKey returns the cache key for mini metadata
func (m *Manager) GetMiniDetailKey(id int) string {
	return fmt.Sprintf(string(MiniDetailKey), id)
}

// GetMountDetailKey returns the cache key for mount metadata
func (m *Manager) GetMountDetailKey(mountType string, id int) string {
	return fmt.Sprintf(string(MountDetailKey), mountType, id)
}

// GetGameBuildKey returns the cache key for game build number
func (m *Manager) GetGameBuildKey() string {
	return string(GameBuildKey)
}

// GetTokenInfoKey returns the cache key for token info
func (m *Manager) GetTokenInfoKey(apiKeyHash string) string {
	return fmt.Sprintf(string(TokenInfoKey), apiKeyHash)
}

// GetDungeonDetailKey returns the cache key for dungeon/raid metadata
func (m *Manager) GetDungeonDetailKey(id string) string {
	return fmt.Sprintf(string(DungeonDetailKey), id)
}
