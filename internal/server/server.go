// Package server provides the MCP server implementation for Guild Wars 2 data access.
package server

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"

	"github.com/AlyxPink/gw2-mcp/internal/cache"
	"github.com/AlyxPink/gw2-mcp/internal/gw2api"
	"github.com/AlyxPink/gw2-mcp/internal/wiki"

	"github.com/charmbracelet/log"
)

// MCPServer wraps the MCP server with GW2-specific functionality
type MCPServer struct {
	mcp    *mcpserver.MCPServer
	logger *log.Logger
	cache  *cache.Manager
	gw2API *gw2api.Client
	wiki   *wiki.Client
}

// NewMCPServer creates a new GW2 MCP server instance
func NewMCPServer(logger *log.Logger) (*MCPServer, error) {
	// Create cache manager
	cacheManager := cache.NewManager()

	// Create GW2 API client
	gw2Client := gw2api.NewClient(cacheManager, logger)

	// Create wiki client
	wikiClient := wiki.NewClient(cacheManager, logger)

	// Create MCP server
	mcpServer := mcpserver.NewMCPServer(
		"GW2 MCP Server",
		"1.0.0",
		mcpserver.WithToolCapabilities(true),
		mcpserver.WithResourceCapabilities(true, true),
		mcpserver.WithRecovery(),
	)

	gw2MCP := &MCPServer{
		mcp:    mcpServer,
		logger: logger,
		cache:  cacheManager,
		gw2API: gw2Client,
		wiki:   wikiClient,
	}

	// Register tools
	gw2MCP.registerTools()

	// Register resources
	gw2MCP.registerResources()

	return gw2MCP, nil
}

// Start starts the MCP server
func (s *MCPServer) Start(ctx context.Context) error {
	s.logger.Info("Starting MCP server on stdio")

	// Create a channel to capture ServeStdio errors
	errChan := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		errChan <- mcpserver.ServeStdio(s.mcp)
	}()

	// Wait for either context cancellation or server error
	select {
	case <-ctx.Done():
		s.logger.Info("Server shutdown requested")
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

// registerTools registers all available tools
func (s *MCPServer) registerTools() {
	// Wiki search tool
	wikiSearchTool := mcp.NewTool(
		"wiki_search",
		mcp.WithDescription("Search Guild Wars 2 wiki for information about game content"),
		mcp.WithString(
			"query",
			mcp.Required(),
			mcp.Description("Search query for wiki content (e.g., 'Dragon Bash', 'currencies', 'wallet')"),
		),
		mcp.WithNumber(
			"limit",
			mcp.Description("Maximum number of results to return (default: 5)"),
		),
	)

	s.mcp.AddTool(wikiSearchTool, s.handleWikiSearch)

	// Wallet info tool
	walletTool := mcp.NewTool(
		"get_wallet",
		mcp.WithDescription("Get user's wallet information including all currencies"),
		mcp.WithString(
			"api_key",
			mcp.Required(),
			mcp.Description("Guild Wars 2 API key with account scope"),
		),
	)

	s.mcp.AddTool(walletTool, s.handleGetWallet)

	// Currency info tool
	currencyTool := mcp.NewTool(
		"get_currencies",
		mcp.WithDescription("Get information about Guild Wars 2 currencies"),
		mcp.WithArray(
			"ids",
			mcp.Description("Specific currency IDs to fetch (optional, returns all if not specified)"),
		),
	)

	s.mcp.AddTool(currencyTool, s.handleGetCurrencies)

	// Trading Post prices tool
	tpPricesTool := mcp.NewTool(
		"get_tp_prices",
		mcp.WithDescription("Get Trading Post prices for items. Returns aggregated best buy/sell prices with item names and formatted coin values."),
		mcp.WithArray(
			"item_ids",
			mcp.Required(),
			mcp.Description("Array of item IDs to get prices for (e.g., [19976, 19721] for Mystic Coin and Glob of Ectoplasm)"),
		),
	)

	s.mcp.AddTool(tpPricesTool, s.handleGetTPPrices)

	// Trading Post listings tool
	tpListingsTool := mcp.NewTool(
		"get_tp_listings",
		mcp.WithDescription("Get Trading Post order book listings for items. Returns all buy/sell price tiers with quantities."),
		mcp.WithArray(
			"item_ids",
			mcp.Required(),
			mcp.Description("Array of item IDs to get listings for"),
		),
	)

	s.mcp.AddTool(tpListingsTool, s.handleGetTPListings)

	// Gem exchange tool
	gemExchangeTool := mcp.NewTool(
		"get_gem_exchange",
		mcp.WithDescription("Get gem exchange rates. Convert coins to gems or gems to coins."),
		mcp.WithString(
			"direction",
			mcp.Required(),
			mcp.Description("Exchange direction: 'coins' (coins to gems) or 'gems' (gems to coins)"),
		),
		mcp.WithNumber(
			"quantity",
			mcp.Required(),
			mcp.Description("Amount to convert (coins in copper, or number of gems)"),
		),
	)

	s.mcp.AddTool(gemExchangeTool, s.handleGetGemExchange)

	// Trading Post delivery tool
	tpDeliveryTool := mcp.NewTool(
		"get_tp_delivery",
		mcp.WithDescription("Get items and coins awaiting pickup from the Trading Post. Requires API key with account and tradingpost scopes."),
		mcp.WithString(
			"api_key",
			mcp.Required(),
			mcp.Description("Guild Wars 2 API key with account and tradingpost scopes"),
		),
	)

	s.mcp.AddTool(tpDeliveryTool, s.handleGetTPDelivery)

	// Trading Post transactions tool
	tpTransactionsTool := mcp.NewTool(
		"get_tp_transactions",
		mcp.WithDescription("Get Trading Post transaction history. View current orders or completed transactions from the past 90 days. Requires API key with account and tradingpost scopes."),
		mcp.WithString(
			"api_key",
			mcp.Required(),
			mcp.Description("Guild Wars 2 API key with account and tradingpost scopes"),
		),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Transaction type: 'current/buys', 'current/sells', 'history/buys', or 'history/sells'"),
		),
	)

	s.mcp.AddTool(tpTransactionsTool, s.handleGetTPTransactions)
}

// registerResources registers all available resources
func (s *MCPServer) registerResources() {
	// Currency list resource
	currencyListResource := mcp.NewResource(
		"gw2://currencies",
		"Guild Wars 2 Currencies",
		mcp.WithResourceDescription("Complete list of all Guild Wars 2 currencies with metadata"),
		mcp.WithMIMEType("application/json"),
	)

	s.mcp.AddResource(currencyListResource, s.handleCurrencyListResource)
}
