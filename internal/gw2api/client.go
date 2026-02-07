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
	apiKey     string
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
func NewClient(cacheManager *cache.Manager, logger *log.Logger, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
		cache:  cacheManager,
		logger: logger,
		apiKey: apiKey,
	}
}

// APIKey returns the configured API key
func (c *Client) APIKey() string {
	return c.apiKey
}

// apiKeyHash returns a short hash of the API key for cache keys
func (c *Client) apiKeyHash() string {
	hash := sha256.Sum256([]byte(c.apiKey))
	return fmt.Sprintf("%x", hash[:8])
}

// GetWallet retrieves wallet information for the configured API key
func (c *Client) GetWallet(ctx context.Context) (*WalletInfo, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("GW2_API_KEY environment variable not configured")
	}

	apiKeyHash := c.apiKeyHash()
	cacheKey := c.cache.GetWalletKey(apiKeyHash)

	// Try to get from cache first
	var walletInfo WalletInfo
	if c.cache.GetJSON(cacheKey, &walletInfo) {
		c.logger.Debug("Wallet cache hit", "api_key_hash", apiKeyHash)
		return &walletInfo, nil
	}

	c.logger.Debug("Wallet cache miss, fetching from API", "api_key_hash", apiKeyHash)

	// Fetch wallet data from API
	walletEntries, err := c.fetchWallet(ctx, c.apiKey)
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

// GetDelivery retrieves trading post delivery box for the configured API key
func (c *Client) GetDelivery(ctx context.Context) (*DeliveryInfo, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("GW2_API_KEY environment variable not configured")
	}

	apiKeyHash := c.apiKeyHash()
	cacheKey := c.cache.GetTPDeliveryKey(apiKeyHash)

	var delivery DeliveryInfo
	if c.cache.GetJSON(cacheKey, &delivery) {
		c.logger.Debug("TP delivery cache hit", "api_key_hash", apiKeyHash)
		return &delivery, nil
	}

	c.logger.Debug("TP delivery cache miss, fetching from API", "api_key_hash", apiKeyHash)

	fetched, err := c.fetchDelivery(ctx, c.apiKey)
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

// GetTransactions retrieves trading post transactions for the configured API key
func (c *Client) GetTransactions(ctx context.Context, txType string) (*TransactionList, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("GW2_API_KEY environment variable not configured")
	}

	if !validTransactionTypes[txType] {
		return nil, fmt.Errorf("invalid transaction type %q: must be one of current/buys, current/sells, history/buys, history/sells", txType)
	}

	apiKeyHash := c.apiKeyHash()
	cacheKey := c.cache.GetTPTransactionKey(apiKeyHash, txType)

	var txList TransactionList
	if c.cache.GetJSON(cacheKey, &txList) {
		c.logger.Debug("TP transactions cache hit", "api_key_hash", apiKeyHash, "type", txType)
		return &txList, nil
	}

	c.logger.Debug("TP transactions cache miss, fetching from API", "api_key_hash", apiKeyHash, "type", txType)

	transactions, err := c.fetchTransactions(ctx, c.apiKey, txType)
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

// requireAPIKey checks that the API key is configured
func (c *Client) requireAPIKey() error {
	if c.apiKey == "" {
		return fmt.Errorf("GW2_API_KEY environment variable not configured")
	}
	return nil
}

// fetchAuthenticated performs an authenticated GET request and decodes JSON into dest
func (c *Client) fetchAuthenticated(ctx context.Context, path string, dest interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+path, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}

// fetchPublic performs an unauthenticated GET request and decodes JSON into dest
func (c *Client) fetchPublic(ctx context.Context, path string, dest interface{}) error {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+path, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			c.logger.Warn("Failed to close response body", "error", closeErr)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return fmt.Errorf("API request failed with status %d and failed to read body: %w", resp.StatusCode, readErr)
		}
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}

// fetchPublicRaw performs an unauthenticated GET request and returns raw JSON
func (c *Client) fetchPublicRaw(ctx context.Context, path string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+path, http.NoBody)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

// fetchAuthenticatedRaw performs an authenticated GET request and returns raw JSON
func (c *Client) fetchAuthenticatedRaw(ctx context.Context, path string) (json.RawMessage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+path, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(body), nil
}

// idsToParam converts a slice of ints to a comma-separated string for API queries
func idsToParam(ids []int) string {
	idStrs := make([]string, len(ids))
	for i, id := range ids {
		idStrs[i] = strconv.Itoa(id)
	}
	return strings.Join(idStrs, ",")
}

// --- Phase 2: Account Tools ---

// AccountInfo represents account data from /v2/account
type AccountInfo struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Age                 int       `json:"age"`
	World               int       `json:"world"`
	Guilds              []string  `json:"guilds"`
	GuildLeader         []string  `json:"guild_leader"`
	Created             string    `json:"created"`
	Access              []string  `json:"access"`
	Commander           bool      `json:"commander"`
	FractalLevel        int       `json:"fractal_level"`
	DailyAP             int       `json:"daily_ap"`
	MonthlyAP           int       `json:"monthly_ap"`
	WvWRank             int       `json:"wvw_rank"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// GetAccount retrieves account information
func (c *Client) GetAccount(ctx context.Context) (*AccountInfo, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	cacheKey := c.cache.GetAccountKey(c.apiKeyHash())
	var info AccountInfo
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	if err := c.fetchAuthenticated(ctx, "/account", &info); err != nil {
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}
	info.UpdatedAt = time.Now()

	if err := c.cache.SetJSON(cacheKey, info, cache.AccountDataTTL); err != nil {
		c.logger.Warn("Failed to cache account data", "error", err)
	}
	return &info, nil
}

// BankSlot represents a single bank vault slot
type BankSlot struct {
	ID        int    `json:"id"`
	Count     int    `json:"count"`
	Charges   int    `json:"charges,omitempty"`
	Skin      int    `json:"skin,omitempty"`
	Binding   string `json:"binding,omitempty"`
	BoundTo   string `json:"bound_to,omitempty"`
	ItemName  string `json:"item_name,omitempty"`
}

// BankInfo represents bank vault contents
type BankInfo struct {
	Slots     []*BankSlot `json:"slots"`
	UsedSlots int         `json:"used_slots"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// GetBank retrieves bank vault contents
func (c *Client) GetBank(ctx context.Context) (*BankInfo, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	cacheKey := c.cache.GetBankKey(c.apiKeyHash())
	var info BankInfo
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	var slots []*BankSlot
	if err := c.fetchAuthenticated(ctx, "/account/bank", &slots); err != nil {
		return nil, fmt.Errorf("failed to fetch bank: %w", err)
	}

	// Enrich with item names
	var itemIDs []int
	for _, slot := range slots {
		if slot != nil {
			itemIDs = append(itemIDs, slot.ID)
		}
	}
	if len(itemIDs) > 0 {
		items, err := c.GetItems(ctx, itemIDs)
		if err != nil {
			c.logger.Warn("Failed to get item metadata for bank", "error", err)
		} else {
			for _, slot := range slots {
				if slot != nil {
					if item, ok := items[slot.ID]; ok {
						slot.ItemName = item.Name
					}
				}
			}
		}
	}

	used := 0
	for _, slot := range slots {
		if slot != nil {
			used++
		}
	}

	info = BankInfo{Slots: slots, UsedSlots: used, UpdatedAt: time.Now()}
	if err := c.cache.SetJSON(cacheKey, info, cache.AccountDataTTL); err != nil {
		c.logger.Warn("Failed to cache bank data", "error", err)
	}
	return &info, nil
}

// MaterialSlot represents a material storage slot
type MaterialSlot struct {
	ID       int    `json:"id"`
	Category int    `json:"category"`
	Count    int    `json:"count"`
	ItemName string `json:"item_name,omitempty"`
}

// MaterialStorage represents material storage contents
type MaterialStorage struct {
	Materials []MaterialSlot `json:"materials"`
	Total     int            `json:"total"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// GetMaterials retrieves material storage contents
func (c *Client) GetMaterials(ctx context.Context) (*MaterialStorage, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	cacheKey := c.cache.GetMaterialsKey(c.apiKeyHash())
	var info MaterialStorage
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	var materials []MaterialSlot
	if err := c.fetchAuthenticated(ctx, "/account/materials", &materials); err != nil {
		return nil, fmt.Errorf("failed to fetch materials: %w", err)
	}

	// Enrich with item names
	itemIDs := make([]int, len(materials))
	for i, m := range materials {
		itemIDs[i] = m.ID
	}
	if len(itemIDs) > 0 {
		items, err := c.GetItems(ctx, itemIDs)
		if err != nil {
			c.logger.Warn("Failed to get item metadata for materials", "error", err)
		} else {
			for i, m := range materials {
				if item, ok := items[m.ID]; ok {
					materials[i].ItemName = item.Name
				}
			}
		}
	}

	info = MaterialStorage{Materials: materials, Total: len(materials), UpdatedAt: time.Now()}
	if err := c.cache.SetJSON(cacheKey, info, cache.AccountDataTTL); err != nil {
		c.logger.Warn("Failed to cache materials data", "error", err)
	}
	return &info, nil
}

// SharedSlot represents a shared inventory slot
type SharedSlot struct {
	ID       int    `json:"id"`
	Count    int    `json:"count"`
	ItemName string `json:"item_name,omitempty"`
}

// InventoryInfo represents shared inventory contents
type InventoryInfo struct {
	Slots     []*SharedSlot `json:"slots"`
	UsedSlots int           `json:"used_slots"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// GetSharedInventory retrieves shared inventory slots
func (c *Client) GetSharedInventory(ctx context.Context) (*InventoryInfo, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	cacheKey := c.cache.GetSharedInventoryKey(c.apiKeyHash())
	var info InventoryInfo
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	var slots []*SharedSlot
	if err := c.fetchAuthenticated(ctx, "/account/inventory", &slots); err != nil {
		return nil, fmt.Errorf("failed to fetch shared inventory: %w", err)
	}

	// Enrich with item names
	var itemIDs []int
	for _, slot := range slots {
		if slot != nil {
			itemIDs = append(itemIDs, slot.ID)
		}
	}
	if len(itemIDs) > 0 {
		items, err := c.GetItems(ctx, itemIDs)
		if err != nil {
			c.logger.Warn("Failed to get item metadata for shared inventory", "error", err)
		} else {
			for _, slot := range slots {
				if slot != nil {
					if item, ok := items[slot.ID]; ok {
						slot.ItemName = item.Name
					}
				}
			}
		}
	}

	used := 0
	for _, slot := range slots {
		if slot != nil {
			used++
		}
	}

	info = InventoryInfo{Slots: slots, UsedSlots: used, UpdatedAt: time.Now()}
	if err := c.cache.SetJSON(cacheKey, info, cache.AccountDataTTL); err != nil {
		c.logger.Warn("Failed to cache shared inventory data", "error", err)
	}
	return &info, nil
}

// GetCharacters retrieves the list of character names
func (c *Client) GetCharacters(ctx context.Context) ([]string, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	cacheKey := c.cache.GetCharactersKey(c.apiKeyHash())
	var names []string
	if c.cache.GetJSON(cacheKey, &names) {
		return names, nil
	}

	if err := c.fetchAuthenticated(ctx, "/characters", &names); err != nil {
		return nil, fmt.Errorf("failed to fetch characters: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, names, cache.AccountDataTTL); err != nil {
		c.logger.Warn("Failed to cache characters", "error", err)
	}
	return names, nil
}

// CharacterInfo represents detailed character information
type CharacterInfo struct {
	Name        string          `json:"name"`
	Race        string          `json:"race"`
	Gender      string          `json:"gender"`
	Profession  string          `json:"profession"`
	Level       int             `json:"level"`
	Age         int             `json:"age"`
	Created     string          `json:"created"`
	Deaths      int             `json:"deaths"`
	Title       int             `json:"title,omitempty"`
	Guild       string          `json:"guild,omitempty"`
}

// GetCharacter retrieves detailed info for a specific character
func (c *Client) GetCharacter(ctx context.Context, name string) (*CharacterInfo, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	cacheKey := c.cache.GetCharacterKey(c.apiKeyHash(), name)
	var info CharacterInfo
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	if err := c.fetchAuthenticated(ctx, "/characters/"+name, &info); err != nil {
		return nil, fmt.Errorf("failed to fetch character %q: %w", name, err)
	}

	if err := c.cache.SetJSON(cacheKey, info, cache.AccountDataTTL); err != nil {
		c.logger.Warn("Failed to cache character data", "name", name, "error", err)
	}
	return &info, nil
}

// --- Phase 3: Account Unlocks ---

var validUnlockTypes = map[string]bool{
	"skins": true, "dyes": true, "minis": true, "titles": true,
	"recipes": true, "finishers": true, "outfits": true, "gliders": true,
	"mailcarriers": true, "novelties": true, "emotes": true,
	"mounts/skins": true, "mounts/types": true, "skiffs": true, "jadebots": true,
}

// GetAccountUnlocks retrieves unlocked IDs for the given unlock type
func (c *Client) GetAccountUnlocks(ctx context.Context, unlockType string) (json.RawMessage, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	if !validUnlockTypes[unlockType] {
		return nil, fmt.Errorf("invalid unlock type %q", unlockType)
	}

	cacheKey := c.cache.GetUnlocksKey(c.apiKeyHash(), unlockType)
	var cached json.RawMessage
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	data, err := c.fetchAuthenticatedRaw(ctx, "/account/"+unlockType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch unlocks for %s: %w", unlockType, err)
	}

	if err := c.cache.SetJSON(cacheKey, data, cache.UnlocksTTL); err != nil {
		c.logger.Warn("Failed to cache unlocks", "type", unlockType, "error", err)
	}
	return data, nil
}

// --- Phase 4: Account Progress ---

var validProgressTypes = map[string]bool{
	"achievements": true, "masteries": true, "mastery/points": true,
	"luck": true, "legendaryarmory": true, "progression": true,
}

// GetAccountProgress retrieves account progress data for the given type
func (c *Client) GetAccountProgress(ctx context.Context, progressType string) (json.RawMessage, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	if !validProgressTypes[progressType] {
		return nil, fmt.Errorf("invalid progress type %q", progressType)
	}

	cacheKey := c.cache.GetProgressKey(c.apiKeyHash(), progressType)
	var cached json.RawMessage
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	data, err := c.fetchAuthenticatedRaw(ctx, "/account/"+progressType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch progress for %s: %w", progressType, err)
	}

	if err := c.cache.SetJSON(cacheKey, data, cache.ProgressTTL); err != nil {
		c.logger.Warn("Failed to cache progress", "type", progressType, "error", err)
	}
	return data, nil
}

// --- Phase 5: Account Dailies ---

var validDailyTypes = map[string]bool{
	"dailycrafting": true, "dungeons": true, "raids": true,
	"mapchests": true, "worldbosses": true,
}

// GetAccountDailies retrieves completed daily IDs/names for the given type
func (c *Client) GetAccountDailies(ctx context.Context, dailyType string) (json.RawMessage, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	if !validDailyTypes[dailyType] {
		return nil, fmt.Errorf("invalid daily type %q", dailyType)
	}

	cacheKey := c.cache.GetDailiesKey(c.apiKeyHash(), dailyType)
	var cached json.RawMessage
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	data, err := c.fetchAuthenticatedRaw(ctx, "/account/"+dailyType)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch dailies for %s: %w", dailyType, err)
	}

	if err := c.cache.SetJSON(cacheKey, data, cache.DailiesTTL); err != nil {
		c.logger.Warn("Failed to cache dailies", "type", dailyType, "error", err)
	}
	return data, nil
}

// --- Phase 6: Wizard's Vault ---

// GetWizardsVault retrieves current wizard's vault season info
func (c *Client) GetWizardsVault(ctx context.Context) (json.RawMessage, error) {
	cacheKey := c.cache.GetWizardsVaultSeasonKey()
	var cached json.RawMessage
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	data, err := c.fetchPublicRaw(ctx, "/wizardsvault")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wizard's vault season: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, data, cache.WVSeasonTTL); err != nil {
		c.logger.Warn("Failed to cache wizard's vault season", "error", err)
	}
	return data, nil
}

var validWVObjectiveTypes = map[string]bool{
	"daily": true, "weekly": true, "special": true,
}

// GetWizardsVaultObjectives retrieves wizard's vault objectives
func (c *Client) GetWizardsVaultObjectives(ctx context.Context, objType string) (json.RawMessage, error) {
	if !validWVObjectiveTypes[objType] {
		return nil, fmt.Errorf("invalid wizard's vault objective type %q: must be daily, weekly, or special", objType)
	}

	// Use authenticated endpoint if API key available
	if c.apiKey != "" {
		keyHash := c.apiKeyHash()
		cacheKey := c.cache.GetWizardsVaultObjectivesKey(keyHash, objType)
		var cached json.RawMessage
		if c.cache.GetJSON(cacheKey, &cached) {
			return cached, nil
		}

		data, err := c.fetchAuthenticatedRaw(ctx, "/account/wizardsvault/"+objType)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch wizard's vault objectives: %w", err)
		}

		if err := c.cache.SetJSON(cacheKey, data, cache.WVObjectivesAuthTTL); err != nil {
			c.logger.Warn("Failed to cache wizard's vault objectives", "error", err)
		}
		return data, nil
	}

	// Fall back to public endpoint
	cacheKey := c.cache.GetWizardsVaultObjectivesKey("public", objType)
	var cached json.RawMessage
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	// Public endpoint returns all objectives, not type-specific
	data, err := c.fetchPublicRaw(ctx, "/wizardsvault/objectives")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wizard's vault objectives: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, data, cache.WVObjectivesPublicTTL); err != nil {
		c.logger.Warn("Failed to cache wizard's vault objectives", "error", err)
	}
	return data, nil
}

// GetWizardsVaultListings retrieves wizard's vault reward listings
func (c *Client) GetWizardsVaultListings(ctx context.Context) (json.RawMessage, error) {
	if c.apiKey != "" {
		keyHash := c.apiKeyHash()
		cacheKey := c.cache.GetWizardsVaultListingsKey(keyHash)
		var cached json.RawMessage
		if c.cache.GetJSON(cacheKey, &cached) {
			return cached, nil
		}

		data, err := c.fetchAuthenticatedRaw(ctx, "/account/wizardsvault/listings")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch wizard's vault listings: %w", err)
		}

		if err := c.cache.SetJSON(cacheKey, data, cache.WVListingsTTL); err != nil {
			c.logger.Warn("Failed to cache wizard's vault listings", "error", err)
		}
		return data, nil
	}

	cacheKey := c.cache.GetWizardsVaultListingsKey("public")
	var cached json.RawMessage
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	data, err := c.fetchPublicRaw(ctx, "/wizardsvault/listings")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wizard's vault listings: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, data, cache.WVListingsTTL); err != nil {
		c.logger.Warn("Failed to cache wizard's vault listings", "error", err)
	}
	return data, nil
}

// --- Phase 7: Game Data Lookups ---

// Skin represents skin metadata from /v2/skins
type Skin struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Icon  string `json:"icon"`
}

// GetSkins retrieves skin metadata for the given IDs
func (c *Client) GetSkins(ctx context.Context, ids []int) ([]Skin, error) {
	var results []Skin
	var missingIDs []int

	for _, id := range ids {
		cacheKey := c.cache.GetSkinDetailKey(id)
		var skin Skin
		if c.cache.GetJSON(cacheKey, &skin) {
			results = append(results, skin)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		var fetched []Skin
		if err := c.fetchPublic(ctx, "/skins?ids="+idsToParam(missingIDs), &fetched); err != nil {
			return nil, fmt.Errorf("failed to fetch skins: %w", err)
		}
		for _, skin := range fetched {
			results = append(results, skin)
			if err := c.cache.SetJSON(c.cache.GetSkinDetailKey(skin.ID), skin, cache.ItemDataTTL); err != nil {
				c.logger.Warn("Failed to cache skin", "id", skin.ID, "error", err)
			}
		}
	}

	return results, nil
}

// Recipe represents recipe data from /v2/recipes
type Recipe struct {
	ID              int              `json:"id"`
	Type            string           `json:"type"`
	OutputItemID    int              `json:"output_item_id"`
	OutputItemCount int              `json:"output_item_count"`
	Disciplines     []string         `json:"disciplines"`
	MinRating       int              `json:"min_rating"`
	Ingredients     []RecipeIngredient `json:"ingredients"`
}

// RecipeIngredient represents an ingredient in a recipe
type RecipeIngredient struct {
	ItemID int `json:"item_id"`
	Count  int `json:"count"`
}

// GetRecipes retrieves recipe data for the given IDs
func (c *Client) GetRecipes(ctx context.Context, ids []int) ([]Recipe, error) {
	var results []Recipe
	var missingIDs []int

	for _, id := range ids {
		cacheKey := c.cache.GetRecipeDetailKey(id)
		var recipe Recipe
		if c.cache.GetJSON(cacheKey, &recipe) {
			results = append(results, recipe)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		var fetched []Recipe
		if err := c.fetchPublic(ctx, "/recipes?ids="+idsToParam(missingIDs), &fetched); err != nil {
			return nil, fmt.Errorf("failed to fetch recipes: %w", err)
		}
		for _, recipe := range fetched {
			results = append(results, recipe)
			if err := c.cache.SetJSON(c.cache.GetRecipeDetailKey(recipe.ID), recipe, cache.RecipeDataTTL); err != nil {
				c.logger.Warn("Failed to cache recipe", "id", recipe.ID, "error", err)
			}
		}
	}

	return results, nil
}

// SearchRecipes searches for recipes by input or output item ID
func (c *Client) SearchRecipes(ctx context.Context, input, output int) ([]int, error) {
	var direction string
	var itemID int
	if input > 0 {
		direction = "input"
		itemID = input
	} else if output > 0 {
		direction = "output"
		itemID = output
	} else {
		return nil, fmt.Errorf("either input or output item ID must be provided")
	}

	cacheKey := c.cache.GetRecipeSearchKey(direction, itemID)
	var cached []int
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	path := fmt.Sprintf("/recipes/search?%s=%d", direction, itemID)
	var ids []int
	if err := c.fetchPublic(ctx, path, &ids); err != nil {
		return nil, fmt.Errorf("failed to search recipes: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, ids, cache.RecipeDataTTL); err != nil {
		c.logger.Warn("Failed to cache recipe search", "error", err)
	}
	return ids, nil
}

// Achievement represents achievement data from /v2/achievements
type Achievement struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Requirement string   `json:"requirement"`
	Type        string   `json:"type"`
	Flags       []string `json:"flags"`
}

// GetAchievements retrieves achievement data for the given IDs
func (c *Client) GetAchievements(ctx context.Context, ids []int) ([]Achievement, error) {
	var results []Achievement
	var missingIDs []int

	for _, id := range ids {
		cacheKey := c.cache.GetAchievementKey(id)
		var ach Achievement
		if c.cache.GetJSON(cacheKey, &ach) {
			results = append(results, ach)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		var fetched []Achievement
		if err := c.fetchPublic(ctx, "/achievements?ids="+idsToParam(missingIDs), &fetched); err != nil {
			return nil, fmt.Errorf("failed to fetch achievements: %w", err)
		}
		for _, ach := range fetched {
			results = append(results, ach)
			if err := c.cache.SetJSON(c.cache.GetAchievementKey(ach.ID), ach, cache.AchievementDataTTL); err != nil {
				c.logger.Warn("Failed to cache achievement", "id", ach.ID, "error", err)
			}
		}
	}

	return results, nil
}

// DailyAchievements represents daily achievement categories
type DailyAchievements struct {
	Today    json.RawMessage `json:"today"`
	Tomorrow json.RawMessage `json:"tomorrow"`
}

// GetDailyAchievements retrieves today's and tomorrow's daily achievements
func (c *Client) GetDailyAchievements(ctx context.Context) (*DailyAchievements, error) {
	cacheKey := c.cache.GetDailyAchievementKey()
	var cached DailyAchievements
	if c.cache.GetJSON(cacheKey, &cached) {
		return &cached, nil
	}

	today, err := c.fetchPublicRaw(ctx, "/achievements/daily")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch daily achievements: %w", err)
	}

	tomorrow, err := c.fetchPublicRaw(ctx, "/achievements/daily/tomorrow")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tomorrow's daily achievements: %w", err)
	}

	result := &DailyAchievements{Today: today, Tomorrow: tomorrow}
	if err := c.cache.SetJSON(cacheKey, result, cache.DailyAchievementTTL); err != nil {
		c.logger.Warn("Failed to cache daily achievements", "error", err)
	}
	return result, nil
}

// --- Phase 8: Guild Tools ---

// GuildInfo represents public guild information
type GuildInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Tag      string `json:"tag"`
	Level    int    `json:"level"`
	MOTD     string `json:"motd,omitempty"`
	Aetherium int   `json:"aetherium,omitempty"`
	Favor    int    `json:"favor,omitempty"`
}

// GetGuild retrieves public guild info
func (c *Client) GetGuild(ctx context.Context, guildID string) (*GuildInfo, error) {
	cacheKey := c.cache.GetGuildInfoKey(guildID)
	var info GuildInfo
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	if err := c.fetchPublic(ctx, "/guild/"+guildID, &info); err != nil {
		return nil, fmt.Errorf("failed to fetch guild: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, info, cache.GuildInfoTTL); err != nil {
		c.logger.Warn("Failed to cache guild info", "id", guildID, "error", err)
	}
	return &info, nil
}

// SearchGuild finds a guild by name
func (c *Client) SearchGuild(ctx context.Context, name string) ([]string, error) {
	cacheKey := c.cache.GetGuildSearchKey(name)
	var cached []string
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	var ids []string
	if err := c.fetchPublic(ctx, "/guild/search?name="+name, &ids); err != nil {
		return nil, fmt.Errorf("failed to search guild: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, ids, cache.GuildSearchTTL); err != nil {
		c.logger.Warn("Failed to cache guild search", "error", err)
	}
	return ids, nil
}

var validGuildDetailTypes = map[string]bool{
	"log": true, "members": true, "ranks": true, "stash": true,
	"storage": true, "treasury": true, "teams": true, "upgrades": true,
}

// GetGuildDetails retrieves authenticated guild detail data
func (c *Client) GetGuildDetails(ctx context.Context, guildID, detailType string) (json.RawMessage, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	if !validGuildDetailTypes[detailType] {
		return nil, fmt.Errorf("invalid guild detail type %q", detailType)
	}

	cacheKey := c.cache.GetGuildDetailKey(guildID, detailType)
	var cached json.RawMessage
	if c.cache.GetJSON(cacheKey, &cached) {
		return cached, nil
	}

	path := fmt.Sprintf("/guild/%s/%s", guildID, detailType)
	data, err := c.fetchAuthenticatedRaw(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch guild %s: %w", detailType, err)
	}

	if err := c.cache.SetJSON(cacheKey, data, cache.GuildDetailTTL); err != nil {
		c.logger.Warn("Failed to cache guild detail", "type", detailType, "error", err)
	}
	return data, nil
}

// --- Phase 9: Game Metadata ---

// Color represents a dye color from /v2/colors
type Color struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GetColors retrieves color metadata for the given IDs
func (c *Client) GetColors(ctx context.Context, ids []int) ([]Color, error) {
	var results []Color
	var missingIDs []int

	for _, id := range ids {
		cacheKey := c.cache.GetColorDetailKey(id)
		var color Color
		if c.cache.GetJSON(cacheKey, &color) {
			results = append(results, color)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		var fetched []Color
		if err := c.fetchPublic(ctx, "/colors?ids="+idsToParam(missingIDs), &fetched); err != nil {
			return nil, fmt.Errorf("failed to fetch colors: %w", err)
		}
		for _, color := range fetched {
			results = append(results, color)
			if err := c.cache.SetJSON(c.cache.GetColorDetailKey(color.ID), color, cache.ColorDataTTL); err != nil {
				c.logger.Warn("Failed to cache color", "id", color.ID, "error", err)
			}
		}
	}

	return results, nil
}

// Mini represents a miniature from /v2/minis
type Mini struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	ItemID int    `json:"item_id"`
}

// GetMinis retrieves mini metadata for the given IDs
func (c *Client) GetMinis(ctx context.Context, ids []int) ([]Mini, error) {
	var results []Mini
	var missingIDs []int

	for _, id := range ids {
		cacheKey := c.cache.GetMiniDetailKey(id)
		var mini Mini
		if c.cache.GetJSON(cacheKey, &mini) {
			results = append(results, mini)
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	if len(missingIDs) > 0 {
		var fetched []Mini
		if err := c.fetchPublic(ctx, "/minis?ids="+idsToParam(missingIDs), &fetched); err != nil {
			return nil, fmt.Errorf("failed to fetch minis: %w", err)
		}
		for _, mini := range fetched {
			results = append(results, mini)
			if err := c.cache.SetJSON(c.cache.GetMiniDetailKey(mini.ID), mini, cache.MiniDataTTL); err != nil {
				c.logger.Warn("Failed to cache mini", "id", mini.ID, "error", err)
			}
		}
	}

	return results, nil
}

// GetMountsInfo retrieves mount info (skins or types) for the given IDs
func (c *Client) GetMountsInfo(ctx context.Context, mountType string, ids []int) (json.RawMessage, error) {
	if mountType != "skins" && mountType != "types" {
		return nil, fmt.Errorf("invalid mount type %q: must be 'skins' or 'types'", mountType)
	}

	path := fmt.Sprintf("/mounts/%s?ids=%s", mountType, idsToParam(ids))
	data, err := c.fetchPublicRaw(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mount %s: %w", mountType, err)
	}

	return data, nil
}

// BuildInfo represents the game build number
type BuildInfo struct {
	ID int `json:"id"`
}

// GetGameBuild retrieves the current game build number
func (c *Client) GetGameBuild(ctx context.Context) (*BuildInfo, error) {
	cacheKey := c.cache.GetGameBuildKey()
	var info BuildInfo
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	if err := c.fetchPublic(ctx, "/build", &info); err != nil {
		return nil, fmt.Errorf("failed to fetch game build: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, info, cache.GameBuildTTL); err != nil {
		c.logger.Warn("Failed to cache game build", "error", err)
	}
	return &info, nil
}

// TokenInfo represents API token info from /v2/tokeninfo
type TokenInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

// GetTokenInfo retrieves API token info
func (c *Client) GetTokenInfo(ctx context.Context) (*TokenInfo, error) {
	if err := c.requireAPIKey(); err != nil {
		return nil, err
	}

	cacheKey := c.cache.GetTokenInfoKey(c.apiKeyHash())
	var info TokenInfo
	if c.cache.GetJSON(cacheKey, &info) {
		return &info, nil
	}

	if err := c.fetchAuthenticated(ctx, "/tokeninfo", &info); err != nil {
		return nil, fmt.Errorf("failed to fetch token info: %w", err)
	}

	if err := c.cache.SetJSON(cacheKey, info, cache.TokenInfoTTL); err != nil {
		c.logger.Warn("Failed to cache token info", "error", err)
	}
	return &info, nil
}

// GetDungeonsAndRaids retrieves dungeon or raid metadata for the given IDs
func (c *Client) GetDungeonsAndRaids(ctx context.Context, contentType string, ids []string) (json.RawMessage, error) {
	if contentType != "dungeons" && contentType != "raids" {
		return nil, fmt.Errorf("invalid type %q: must be 'dungeons' or 'raids'", contentType)
	}

	idsParam := strings.Join(ids, ",")
	path := fmt.Sprintf("/%s?ids=%s", contentType, idsParam)
	data, err := c.fetchPublicRaw(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", contentType, err)
	}

	return data, nil
}
