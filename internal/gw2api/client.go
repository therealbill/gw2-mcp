// Package gw2api provides functionality for interacting with the Guild Wars 2 API.
package gw2api

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/AlyxPink/gw2-mcp/internal/cache"
)

const (
	baseURL        = "https://api.guildwars2.com/v2"
	userAgent      = "github.com/AlyxPink/gw2-mcp"
	requestTimeout = 30 * time.Second
)

// Client handles GW2 API requests
type Client struct {
	httpClient *http.Client
	cache      *cache.Manager
	logger     *log.Logger
}

// WalletEntry represents a single currency in the wallet
type WalletEntry struct {
	ID    int `json:"id"`
	Value int `json:"value"`
}

// Currency represents currency metadata
type Currency struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ID          int    `json:"id"`
	Order       int    `json:"order"`
}

// WalletInfo combines wallet entries with currency metadata
type WalletInfo struct {
	UpdatedAt  time.Time        `json:"updated_at"`
	Currencies map[int]Currency `json:"currencies"`
	Entries    []WalletEntry    `json:"entries"`
	Total      int              `json:"total_currencies"`
}

// Item represents basic item metadata from /v2/items
type Item struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Rarity string `json:"rarity"`
	Level  int    `json:"level"`
	Icon   string `json:"icon"`
}

// FormatCoins converts copper coins to a human-readable string (e.g., "1g 50s 35c")
func FormatCoins(copper int) string {
	gold := copper / 10000
	silver := (copper % 10000) / 100
	cop := copper % 100

	parts := []string{}
	if gold > 0 {
		parts = append(parts, fmt.Sprintf("%dg", gold))
	}
	if silver > 0 || gold > 0 {
		parts = append(parts, fmt.Sprintf("%ds", silver))
	}
	parts = append(parts, fmt.Sprintf("%dc", cop))

	return strings.Join(parts, " ")
}

// PricePoint represents aggregated buy or sell price data
type PricePoint struct {
	UnitPrice int `json:"unit_price"`
	Quantity  int `json:"quantity"`
}

// PriceInfo represents trading post price data for an item
type PriceInfo struct {
	ID          int        `json:"id"`
	Whitelisted bool       `json:"whitelisted"`
	Buys        PricePoint `json:"buys"`
	Sells       PricePoint `json:"sells"`
	ItemName    string     `json:"item_name,omitempty"`
	BuyPrice    string     `json:"buy_price_formatted,omitempty"`
	SellPrice   string     `json:"sell_price_formatted,omitempty"`
}

// ListingEntry represents a single price tier in the order book
type ListingEntry struct {
	Listings  int `json:"listings"`
	UnitPrice int `json:"unit_price"`
	Quantity  int `json:"quantity"`
}

// ListingInfo represents detailed trading post order book for an item
type ListingInfo struct {
	ID       int            `json:"id"`
	Buys     []ListingEntry `json:"buys"`
	Sells    []ListingEntry `json:"sells"`
	ItemName string         `json:"item_name,omitempty"`
}

// ExchangeRate represents gem exchange rate data
type ExchangeRate struct {
	CoinsPerGem   int    `json:"coins_per_gem"`
	Quantity      int    `json:"quantity"`
	Direction     string `json:"direction"`
	FormattedRate string `json:"formatted_rate"`
}

// DeliveryItem represents an item awaiting pickup from the Trading Post
type DeliveryItem struct {
	ID       int    `json:"id"`
	Count    int    `json:"count"`
	ItemName string `json:"item_name,omitempty"`
}

// DeliveryInfo represents the Trading Post delivery box contents
type DeliveryInfo struct {
	Coins          int            `json:"coins"`
	FormattedCoins string         `json:"formatted_coins"`
	Items          []DeliveryItem `json:"items"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// Transaction represents a single trading post transaction
type Transaction struct {
	ID             int    `json:"id"`
	ItemID         int    `json:"item_id"`
	Price          int    `json:"price"`
	Quantity       int    `json:"quantity"`
	Created        string `json:"created"`
	Purchased      string `json:"purchased,omitempty"`
	ItemName       string `json:"item_name,omitempty"`
	FormattedPrice string `json:"formatted_price,omitempty"`
}

// TransactionList represents a list of trading post transactions
type TransactionList struct {
	Type         string        `json:"type"`
	Transactions []Transaction `json:"transactions"`
	Total        int           `json:"total"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

// NewClient creates a new GW2 API client
func NewClient(cacheManager *cache.Manager, logger *log.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
		cache:  cacheManager,
		logger: logger,
	}
}

// GetWallet retrieves wallet information for the given API key
func (c *Client) GetWallet(ctx context.Context, apiKey string) (*WalletInfo, error) {
	// Create a hash of the API key for caching (security)
	hash := sha256.Sum256([]byte(apiKey))
	apiKeyHash := fmt.Sprintf("%x", hash[:8]) // Use first 8 bytes of hash

	cacheKey := c.cache.GetWalletKey(apiKeyHash)

	// Try to get from cache first
	var walletInfo WalletInfo
	if c.cache.GetJSON(cacheKey, &walletInfo) {
		c.logger.Debug("Wallet cache hit", "api_key_hash", apiKeyHash)
		return &walletInfo, nil
	}

	c.logger.Debug("Wallet cache miss, fetching from API", "api_key_hash", apiKeyHash)

	// Fetch wallet data from API
	walletEntries, err := c.fetchWallet(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wallet: %w", err)
	}

	// Get currency metadata for all currencies in wallet
	currencyIDs := make([]int, len(walletEntries))
	for i, entry := range walletEntries {
		currencyIDs[i] = entry.ID
	}

	currencies, err := c.GetCurrencies(ctx, currencyIDs)
	if err != nil {
		c.logger.Warn("Failed to get currency metadata", "error", err)
		// Continue without metadata
		currencies = make(map[int]Currency)
	}

	// Create wallet info
	walletInfo = WalletInfo{
		Entries:    walletEntries,
		Currencies: currencies,
		Total:      len(walletEntries),
		UpdatedAt:  time.Now(),
	}

	// Cache the result
	if err := c.cache.SetJSON(cacheKey, walletInfo, cache.WalletDataTTL); err != nil {
		c.logger.Warn("Failed to cache wallet data", "error", err)
	}

	return &walletInfo, nil
}

// GetCurrencies retrieves currency metadata
func (c *Client) GetCurrencies(ctx context.Context, ids []int) (map[int]Currency, error) {
	// If no specific IDs requested, get all currencies
	if len(ids) == 0 {
		return c.getAllCurrencies(ctx)
	}

	// Get specific currencies
	currencies := make(map[int]Currency)
	var missingIDs []int

	// Check cache for each currency
	for _, id := range ids {
		cacheKey := c.cache.GetCurrencyDetailKey(id)
		var currency Currency
		if c.cache.GetJSON(cacheKey, &currency) {
			currencies[id] = currency
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// Fetch missing currencies from API
	if len(missingIDs) > 0 {
		fetchedCurrencies, err := c.fetchCurrencies(ctx, missingIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch currencies: %w", err)
		}

		// Add fetched currencies to result and cache
		for _, currency := range fetchedCurrencies {
			currencies[currency.ID] = currency
			cacheKey := c.cache.GetCurrencyDetailKey(currency.ID)
			if err := c.cache.SetJSON(cacheKey, currency, cache.StaticDataTTL); err != nil {
				c.logger.Warn("Failed to cache currency", "id", currency.ID, "error", err)
			}
		}
	}

	return currencies, nil
}

// getAllCurrencies retrieves all available currencies
func (c *Client) getAllCurrencies(ctx context.Context) (map[int]Currency, error) {
	cacheKey := c.cache.GetCurrencyListKey()

	// Try cache first
	var currencies map[int]Currency
	if c.cache.GetJSON(cacheKey, &currencies) {
		c.logger.Debug("Currency list cache hit")
		return currencies, nil
	}

	c.logger.Debug("Currency list cache miss, fetching from API")

	// Fetch all currency IDs first
	currencyIDs, err := c.fetchCurrencyIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currency IDs: %w", err)
	}

	// Fetch all currency details
	currencyList, err := c.fetchCurrencies(ctx, currencyIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch currency details: %w", err)
	}

	// Convert to map
	currencies = make(map[int]Currency)
	for _, currency := range currencyList {
		currencies[currency.ID] = currency
	}

	// Cache the result
	if err := c.cache.SetJSON(cacheKey, currencies, cache.StaticDataTTL); err != nil {
		c.logger.Warn("Failed to cache currency list", "error", err)
	}

	return currencies, nil
}

// GetItems retrieves item metadata for the given IDs
func (c *Client) GetItems(ctx context.Context, ids []int) (map[int]Item, error) {
	items := make(map[int]Item)
	var missingIDs []int

	// Check cache for each item
	for _, id := range ids {
		cacheKey := c.cache.GetItemDetailKey(id)
		var item Item
		if c.cache.GetJSON(cacheKey, &item) {
			items[id] = item
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// Fetch missing items from API
	if len(missingIDs) > 0 {
		fetchedItems, err := c.fetchItems(ctx, missingIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch items: %w", err)
		}

		for _, item := range fetchedItems {
			items[item.ID] = item
			cacheKey := c.cache.GetItemDetailKey(item.ID)
			if err := c.cache.SetJSON(cacheKey, item, cache.ItemDataTTL); err != nil {
				c.logger.Warn("Failed to cache item", "id", item.ID, "error", err)
			}
		}
	}

	return items, nil
}

// fetchWallet makes the actual API call to get wallet data
func (c *Client) fetchWallet(ctx context.Context, apiKey string) ([]WalletEntry, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/account/wallet", http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var wallet []WalletEntry
	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

// fetchCurrencyIDs fetches all available currency IDs
func (c *Client) fetchCurrencyIDs(ctx context.Context) ([]int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/currencies", http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}

// fetchCurrencies fetches currency details for specific IDs
func (c *Client) fetchCurrencies(ctx context.Context, ids []int) ([]Currency, error) {
	// Convert IDs to comma-separated string
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	idsParam := strings.Join(idStrs, ",")

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/currencies?ids="+idsParam, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var currencies []Currency
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		return nil, err
	}

	return currencies, nil
}

// fetchItems fetches item details for specific IDs from /v2/items
func (c *Client) fetchItems(ctx context.Context, ids []int) ([]Item, error) {
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	idsParam := strings.Join(idStrs, ",")

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/items?ids="+idsParam, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	return items, nil
}

// GetPrices retrieves trading post prices for the given item IDs
func (c *Client) GetPrices(ctx context.Context, itemIDs []int) ([]PriceInfo, error) {
	var results []PriceInfo
	var missingIDs []int

	// Check cache for each item
	for _, id := range itemIDs {
		cacheKey := c.cache.GetTPPriceKey(id)
		var price PriceInfo
		if c.cache.GetJSON(cacheKey, &price) {
			results = append(results, price)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// Fetch missing prices from API
	if len(missingIDs) > 0 {
		fetched, err := c.fetchPrices(ctx, missingIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch prices: %w", err)
		}

		results = append(results, fetched...)

		// Cache each price
		for _, price := range fetched {
			cacheKey := c.cache.GetTPPriceKey(price.ID)
			if err := c.cache.SetJSON(cacheKey, price, cache.TPPriceTTL); err != nil {
				c.logger.Warn("Failed to cache price", "id", price.ID, "error", err)
			}
		}
	}

	// Enrich with item names
	allIDs := make([]int, len(results))
	for i, r := range results {
		allIDs[i] = r.ID
	}
	items, err := c.GetItems(ctx, allIDs)
	if err != nil {
		c.logger.Warn("Failed to get item metadata for prices", "error", err)
	} else {
		for i, r := range results {
			if item, ok := items[r.ID]; ok {
				results[i].ItemName = item.Name
			}
		}
	}

	// Add formatted prices
	for i, r := range results {
		results[i].BuyPrice = FormatCoins(r.Buys.UnitPrice)
		results[i].SellPrice = FormatCoins(r.Sells.UnitPrice)
	}

	return results, nil
}

// fetchPrices fetches trading post prices from /v2/commerce/prices
func (c *Client) fetchPrices(ctx context.Context, ids []int) ([]PriceInfo, error) {
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	idsParam := strings.Join(idStrs, ",")

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/commerce/prices?ids="+idsParam, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var prices []PriceInfo
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		return nil, err
	}

	return prices, nil
}

// GetListings retrieves trading post listings for the given item IDs
func (c *Client) GetListings(ctx context.Context, itemIDs []int) ([]ListingInfo, error) {
	var results []ListingInfo
	var missingIDs []int

	for _, id := range itemIDs {
		cacheKey := c.cache.GetTPListingKey(id)
		var listing ListingInfo
		if c.cache.GetJSON(cacheKey, &listing) {
			results = append(results, listing)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		fetched, err := c.fetchListings(ctx, missingIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch listings: %w", err)
		}

		results = append(results, fetched...)

		for _, listing := range fetched {
			cacheKey := c.cache.GetTPListingKey(listing.ID)
			if err := c.cache.SetJSON(cacheKey, listing, cache.TPListingTTL); err != nil {
				c.logger.Warn("Failed to cache listing", "id", listing.ID, "error", err)
			}
		}
	}

	// Enrich with item names
	allIDs := make([]int, len(results))
	for i, r := range results {
		allIDs[i] = r.ID
	}
	items, err := c.GetItems(ctx, allIDs)
	if err != nil {
		c.logger.Warn("Failed to get item metadata for listings", "error", err)
	} else {
		for i, r := range results {
			if item, ok := items[r.ID]; ok {
				results[i].ItemName = item.Name
			}
		}
	}

	return results, nil
}

// fetchListings fetches trading post listings from /v2/commerce/listings
func (c *Client) fetchListings(ctx context.Context, ids []int) ([]ListingInfo, error) {
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	idsParam := strings.Join(idStrs, ",")

	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/commerce/listings?ids="+idsParam, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var listings []ListingInfo
	if err := json.NewDecoder(resp.Body).Decode(&listings); err != nil {
		return nil, err
	}

	return listings, nil
}

// GetGemExchange retrieves gem exchange rates
func (c *Client) GetGemExchange(ctx context.Context, direction string, quantity int) (*ExchangeRate, error) {
	if direction != "coins" && direction != "gems" {
		return nil, fmt.Errorf("invalid direction %q: must be \"coins\" or \"gems\"", direction)
	}

	cacheKey := c.cache.GetTPExchangeKey(direction, quantity)

	var rate ExchangeRate
	if c.cache.GetJSON(cacheKey, &rate) {
		c.logger.Debug("TP exchange cache hit", "direction", direction, "quantity", quantity)
		return &rate, nil
	}

	c.logger.Debug("TP exchange cache miss, fetching from API", "direction", direction, "quantity", quantity)

	fetched, err := c.fetchGemExchange(ctx, direction, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gem exchange: %w", err)
	}

	fetched.Direction = direction
	fetched.FormattedRate = FormatCoins(fetched.CoinsPerGem) + " per gem"

	if err := c.cache.SetJSON(cacheKey, fetched, cache.TPExchangeTTL); err != nil {
		c.logger.Warn("Failed to cache exchange rate", "error", err)
	}

	return fetched, nil
}

// fetchGemExchange fetches gem exchange rates from /v2/commerce/exchange
func (c *Client) fetchGemExchange(ctx context.Context, direction string, quantity int) (*ExchangeRate, error) {
	url := fmt.Sprintf("%s/commerce/exchange/%s?quantity=%d", baseURL, direction, quantity)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var rate ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&rate); err != nil {
		return nil, err
	}

	return &rate, nil
}

// GetDelivery retrieves trading post delivery box for the given API key
func (c *Client) GetDelivery(ctx context.Context, apiKey string) (*DeliveryInfo, error) {
	hash := sha256.Sum256([]byte(apiKey))
	apiKeyHash := fmt.Sprintf("%x", hash[:8])

	cacheKey := c.cache.GetTPDeliveryKey(apiKeyHash)

	var delivery DeliveryInfo
	if c.cache.GetJSON(cacheKey, &delivery) {
		c.logger.Debug("TP delivery cache hit", "api_key_hash", apiKeyHash)
		return &delivery, nil
	}

	c.logger.Debug("TP delivery cache miss, fetching from API", "api_key_hash", apiKeyHash)

	fetched, err := c.fetchDelivery(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch delivery: %w", err)
	}

	delivery = *fetched
	delivery.FormattedCoins = FormatCoins(delivery.Coins)
	delivery.UpdatedAt = time.Now()

	// Enrich items with names
	if len(delivery.Items) > 0 {
		itemIDs := make([]int, len(delivery.Items))
		for i, item := range delivery.Items {
			itemIDs[i] = item.ID
		}
		items, err := c.GetItems(ctx, itemIDs)
		if err != nil {
			c.logger.Warn("Failed to get item metadata for delivery", "error", err)
		} else {
			for i, di := range delivery.Items {
				if item, ok := items[di.ID]; ok {
					delivery.Items[i].ItemName = item.Name
				}
			}
		}
	}

	if err := c.cache.SetJSON(cacheKey, delivery, cache.TPDeliveryTTL); err != nil {
		c.logger.Warn("Failed to cache delivery data", "error", err)
	}

	return &delivery, nil
}

// fetchDelivery fetches delivery box from /v2/commerce/delivery
func (c *Client) fetchDelivery(ctx context.Context, apiKey string) (*DeliveryInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/commerce/delivery", http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var delivery DeliveryInfo
	if err := json.NewDecoder(resp.Body).Decode(&delivery); err != nil {
		return nil, err
	}

	return &delivery, nil
}

// validTransactionTypes lists the allowed transaction type values
var validTransactionTypes = map[string]bool{
	"current/buys":  true,
	"current/sells": true,
	"history/buys":  true,
	"history/sells": true,
}

// GetTransactions retrieves trading post transactions for the given API key
func (c *Client) GetTransactions(ctx context.Context, apiKey string, txType string) (*TransactionList, error) {
	if !validTransactionTypes[txType] {
		return nil, fmt.Errorf("invalid transaction type %q: must be one of current/buys, current/sells, history/buys, history/sells", txType)
	}

	hash := sha256.Sum256([]byte(apiKey))
	apiKeyHash := fmt.Sprintf("%x", hash[:8])

	cacheKey := c.cache.GetTPTransactionKey(apiKeyHash, txType)

	var txList TransactionList
	if c.cache.GetJSON(cacheKey, &txList) {
		c.logger.Debug("TP transactions cache hit", "api_key_hash", apiKeyHash, "type", txType)
		return &txList, nil
	}

	c.logger.Debug("TP transactions cache miss, fetching from API", "api_key_hash", apiKeyHash, "type", txType)

	transactions, err := c.fetchTransactions(ctx, apiKey, txType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Enrich with item names and formatted prices
	if len(transactions) > 0 {
		itemIDs := make([]int, len(transactions))
		for i, tx := range transactions {
			itemIDs[i] = tx.ItemID
		}
		items, err := c.GetItems(ctx, itemIDs)
		if err != nil {
			c.logger.Warn("Failed to get item metadata for transactions", "error", err)
		} else {
			for i, tx := range transactions {
				if item, ok := items[tx.ItemID]; ok {
					transactions[i].ItemName = item.Name
				}
			}
		}

		for i, tx := range transactions {
			transactions[i].FormattedPrice = FormatCoins(tx.Price)
		}
	}

	txList = TransactionList{
		Type:         txType,
		Transactions: transactions,
		Total:        len(transactions),
		UpdatedAt:    time.Now(),
	}

	if err := c.cache.SetJSON(cacheKey, txList, cache.TPTransactionTTL); err != nil {
		c.logger.Warn("Failed to cache transactions", "error", err)
	}

	return &txList, nil
}

// fetchTransactions fetches transactions from /v2/commerce/transactions
func (c *Client) fetchTransactions(ctx context.Context, apiKey string, txType string) ([]Transaction, error) {
	url := fmt.Sprintf("%s/commerce/transactions/%s", baseURL, txType)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var transactions []Transaction
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}
