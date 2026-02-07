package cache

import (
	"testing"
	"time"
)

func TestManager_SetAndGet(t *testing.T) {
	m := NewManager()

	// Test basic set and get
	key := "test_key"
	value := "test_value"
	m.Set(key, value, time.Minute)

	retrieved, found := m.Get(key)
	if !found {
		t.Error("Expected to find cached value")
	}

	if retrieved != value {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}
}

func TestManager_GetString(t *testing.T) {
	m := NewManager()

	// Test string retrieval
	key := "string_key"
	value := "string_value"
	m.Set(key, value, time.Minute)

	retrieved, found := m.GetString(key)
	if !found {
		t.Error("Expected to find cached string value")
	}

	if retrieved != value {
		t.Errorf("Expected %s, got %s", value, retrieved)
	}

	// Test non-string value
	m.Set("non_string", 123, time.Minute)
	_, found = m.GetString("non_string")
	if found {
		t.Error("Expected not to find non-string value as string")
	}
}

func TestManager_SetAndGetJSON(t *testing.T) {
	m := NewManager()

	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	original := TestStruct{
		Name:  "test",
		Value: 42,
	}

	// Test JSON set and get
	key := "json_key"
	err := m.SetJSON(key, original, time.Minute)
	if err != nil {
		t.Fatalf("Failed to set JSON: %v", err)
	}

	var retrieved TestStruct
	found := m.GetJSON(key, &retrieved)
	if !found {
		t.Error("Expected to find cached JSON value")
	}

	if retrieved.Name != original.Name || retrieved.Value != original.Value {
		t.Errorf("Expected %+v, got %+v", original, retrieved)
	}
}

func TestManager_Delete(t *testing.T) {
	m := NewManager()

	key := "delete_key"
	value := "delete_value"
	m.Set(key, value, time.Minute)

	// Verify it exists
	_, found := m.Get(key)
	if !found {
		t.Error("Expected to find value before deletion")
	}

	// Delete and verify it's gone
	m.Delete(key)
	_, found = m.Get(key)
	if found {
		t.Error("Expected value to be deleted")
	}
}

func TestManager_Flush(t *testing.T) {
	m := NewManager()

	// Add multiple items
	m.Set("key1", "value1", time.Minute)
	m.Set("key2", "value2", time.Minute)
	m.Set("key3", "value3", time.Minute)

	if m.ItemCount() != 3 {
		t.Errorf("Expected 3 items, got %d", m.ItemCount())
	}

	// Flush all
	m.Flush()

	if m.ItemCount() != 0 {
		t.Errorf("Expected 0 items after flush, got %d", m.ItemCount())
	}
}

func TestManager_CacheKeys(t *testing.T) {
	m := NewManager()

	// Test currency list key
	key := m.GetCurrencyListKey()
	expected := string(CurrencyListKey)
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test currency detail key
	currencyID := 1
	key = m.GetCurrencyDetailKey(currencyID)
	expected = "currency:detail:1"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test wiki search key
	query := "test query"
	key = m.GetWikiSearchKey(query)
	expected = "wiki:search:test query"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test wiki page key
	pageTitle := "Test Page"
	key = m.GetWikiPageKey(pageTitle)
	expected = "wiki:page:Test Page"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test wallet key
	apiKeyHash := "abcd1234"
	key = m.GetWalletKey(apiKeyHash)
	expected = "wallet:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test item detail key
	itemID := 19976
	key = m.GetItemDetailKey(itemID)
	expected = "item:detail:19976"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test TP price key
	key = m.GetTPPriceKey(19976)
	expected = "tp:price:19976"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test TP listing key
	key = m.GetTPListingKey(19976)
	expected = "tp:listing:19976"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test TP exchange key
	key = m.GetTPExchangeKey("coins", 100000)
	expected = "tp:exchange:coins:100000"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test TP delivery key
	key = m.GetTPDeliveryKey("abcd1234")
	expected = "tp:delivery:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test TP transaction key
	key = m.GetTPTransactionKey("abcd1234", "current/buys")
	expected = "tp:transactions:abcd1234:current/buys"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test account key
	key = m.GetAccountKey("abcd1234")
	expected = "account:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test bank key
	key = m.GetBankKey("abcd1234")
	expected = "bank:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test materials key
	key = m.GetMaterialsKey("abcd1234")
	expected = "materials:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test shared inventory key
	key = m.GetSharedInventoryKey("abcd1234")
	expected = "inventory:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test characters key
	key = m.GetCharactersKey("abcd1234")
	expected = "characters:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test character key
	key = m.GetCharacterKey("abcd1234", "My Character")
	expected = "character:abcd1234:My Character"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test unlocks key
	key = m.GetUnlocksKey("abcd1234", "skins")
	expected = "unlocks:abcd1234:skins"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test progress key
	key = m.GetProgressKey("abcd1234", "achievements")
	expected = "progress:abcd1234:achievements"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test dailies key
	key = m.GetDailiesKey("abcd1234", "dungeons")
	expected = "dailies:abcd1234:dungeons"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test wizard's vault season key
	key = m.GetWizardsVaultSeasonKey()
	expected = "wv:season"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test wizard's vault objectives key
	key = m.GetWizardsVaultObjectivesKey("abcd1234", "daily")
	expected = "wv:objectives:abcd1234:daily"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test wizard's vault listings key
	key = m.GetWizardsVaultListingsKey("abcd1234")
	expected = "wv:listings:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test skin detail key
	key = m.GetSkinDetailKey(100)
	expected = "skin:detail:100"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test recipe detail key
	key = m.GetRecipeDetailKey(200)
	expected = "recipe:detail:200"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test recipe search key
	key = m.GetRecipeSearchKey("input", 19976)
	expected = "recipe:search:input:19976"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test achievement key
	key = m.GetAchievementKey(300)
	expected = "achievement:detail:300"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test daily achievement key
	key = m.GetDailyAchievementKey()
	expected = "achievements:daily"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test guild info key
	key = m.GetGuildInfoKey("guild-uuid-123")
	expected = "guild:info:guild-uuid-123"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test guild search key
	key = m.GetGuildSearchKey("Test Guild")
	expected = "guild:search:Test Guild"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test guild detail key
	key = m.GetGuildDetailKey("guild-uuid-123", "members")
	expected = "guild:detail:guild-uuid-123:members"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test color detail key
	key = m.GetColorDetailKey(400)
	expected = "color:detail:400"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test mini detail key
	key = m.GetMiniDetailKey(500)
	expected = "mini:detail:500"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test mount detail key
	key = m.GetMountDetailKey("skins", 600)
	expected = "mount:skins:detail:600"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test game build key
	key = m.GetGameBuildKey()
	expected = "game:build"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test token info key
	key = m.GetTokenInfoKey("abcd1234")
	expected = "tokeninfo:abcd1234"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}

	// Test dungeon detail key
	key = m.GetDungeonDetailKey("ascalonian_catacombs")
	expected = "dungeon:detail:ascalonian_catacombs"
	if key != expected {
		t.Errorf("Expected %s, got %s", expected, key)
	}
}

func TestManager_TTLExpiration(t *testing.T) {
	m := NewManager()

	key := "ttl_key"
	value := "ttl_value"
	ttl := 100 * time.Millisecond

	m.Set(key, value, ttl)

	// Should exist immediately
	_, found := m.Get(key)
	if !found {
		t.Error("Expected to find value immediately after setting")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	_, found = m.Get(key)
	if found {
		t.Error("Expected value to be expired")
	}
}
