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
func NewMCPServer(logger *log.Logger, apiKey string) (*MCPServer, error) {
	// Create cache manager
	cacheManager := cache.NewManager()

	// Create GW2 API client
	gw2Client := gw2api.NewClient(cacheManager, logger, apiKey)

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
		mcp.WithDescription("Get user's wallet information including all currencies. Requires GW2_API_KEY environment variable."),
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
		mcp.WithDescription("Get items and coins awaiting pickup from the Trading Post. Requires GW2_API_KEY with account and tradingpost scopes."),
	)

	s.mcp.AddTool(tpDeliveryTool, s.handleGetTPDelivery)

	// Trading Post transactions tool
	tpTransactionsTool := mcp.NewTool(
		"get_tp_transactions",
		mcp.WithDescription("Get Trading Post transaction history. View current orders or completed transactions from the past 90 days. Requires GW2_API_KEY with account and tradingpost scopes."),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Transaction type: 'current/buys', 'current/sells', 'history/buys', or 'history/sells'"),
		),
	)

	s.mcp.AddTool(tpTransactionsTool, s.handleGetTPTransactions)

	// --- Account Tools ---

	accountTool := mcp.NewTool(
		"get_account",
		mcp.WithDescription("Get account information including name, world, guilds, and access. Requires GW2_API_KEY."),
	)
	s.mcp.AddTool(accountTool, s.handleGetAccount)

	bankTool := mcp.NewTool(
		"get_bank",
		mcp.WithDescription("Get bank vault contents with item names. Requires GW2_API_KEY."),
	)
	s.mcp.AddTool(bankTool, s.handleGetBank)

	materialsTool := mcp.NewTool(
		"get_materials",
		mcp.WithDescription("Get material storage contents with item names. Requires GW2_API_KEY."),
	)
	s.mcp.AddTool(materialsTool, s.handleGetMaterials)

	inventoryTool := mcp.NewTool(
		"get_inventory",
		mcp.WithDescription("Get shared inventory slot contents with item names. Requires GW2_API_KEY."),
	)
	s.mcp.AddTool(inventoryTool, s.handleGetInventory)

	charactersTool := mcp.NewTool(
		"get_characters",
		mcp.WithDescription("Get list of character names, or detailed info for a specific character. Requires GW2_API_KEY."),
		mcp.WithString(
			"name",
			mcp.Description("Character name to get details for (optional; omit to list all characters)"),
		),
	)
	s.mcp.AddTool(charactersTool, s.handleGetCharacters)

	// --- Account Unlocks ---

	unlocksTool := mcp.NewTool(
		"get_account_unlocks",
		mcp.WithDescription("Get account unlocked IDs. Requires GW2_API_KEY."),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Unlock type: skins, dyes, minis, titles, recipes, finishers, outfits, gliders, mailcarriers, novelties, emotes, mounts/skins, mounts/types, skiffs, jadebots"),
		),
	)
	s.mcp.AddTool(unlocksTool, s.handleGetAccountUnlocks)

	// --- Account Progress ---

	progressTool := mcp.NewTool(
		"get_account_progress",
		mcp.WithDescription("Get account progress data. Requires GW2_API_KEY."),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Progress type: achievements, masteries, mastery/points, luck, legendaryarmory, progression"),
		),
	)
	s.mcp.AddTool(progressTool, s.handleGetAccountProgress)

	// --- Account Dailies ---

	dailiesTool := mcp.NewTool(
		"get_account_dailies",
		mcp.WithDescription("Get completed daily content IDs. Requires GW2_API_KEY."),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Daily type: dailycrafting, dungeons, raids, mapchests, worldbosses"),
		),
	)
	s.mcp.AddTool(dailiesTool, s.handleGetAccountDailies)

	// --- Wizard's Vault ---

	wvTool := mcp.NewTool(
		"get_wizards_vault",
		mcp.WithDescription("Get current Wizard's Vault season information."),
	)
	s.mcp.AddTool(wvTool, s.handleGetWizardsVault)

	wvObjectivesTool := mcp.NewTool(
		"get_wizards_vault_objectives",
		mcp.WithDescription("Get Wizard's Vault objectives. Uses authenticated endpoint if GW2_API_KEY is set, otherwise returns public objective list."),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Objective type: daily, weekly, special"),
		),
	)
	s.mcp.AddTool(wvObjectivesTool, s.handleGetWizardsVaultObjectives)

	wvListingsTool := mcp.NewTool(
		"get_wizards_vault_listings",
		mcp.WithDescription("Get Wizard's Vault reward listings. Uses authenticated endpoint if GW2_API_KEY is set."),
	)
	s.mcp.AddTool(wvListingsTool, s.handleGetWizardsVaultListings)

	// --- Game Data Lookups ---

	itemsTool := mcp.NewTool(
		"get_items",
		mcp.WithDescription("Get item metadata (name, type, rarity, level, icon) for given item IDs."),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of item IDs to look up"),
		),
	)
	s.mcp.AddTool(itemsTool, s.handleGetItems)

	skinsTool := mcp.NewTool(
		"get_skins",
		mcp.WithDescription("Get skin metadata (name, type, icon) for given skin IDs."),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of skin IDs to look up"),
		),
	)
	s.mcp.AddTool(skinsTool, s.handleGetSkins)

	recipesTool := mcp.NewTool(
		"get_recipes",
		mcp.WithDescription("Get recipe details (type, output, ingredients, disciplines) for given recipe IDs."),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of recipe IDs to look up"),
		),
	)
	s.mcp.AddTool(recipesTool, s.handleGetRecipes)

	searchRecipesTool := mcp.NewTool(
		"search_recipes",
		mcp.WithDescription("Search for recipe IDs by input or output item ID."),
		mcp.WithNumber(
			"input",
			mcp.Description("Input item ID to search for recipes that use this item"),
		),
		mcp.WithNumber(
			"output",
			mcp.Description("Output item ID to search for recipes that produce this item"),
		),
	)
	s.mcp.AddTool(searchRecipesTool, s.handleSearchRecipes)

	achievementsTool := mcp.NewTool(
		"get_achievements",
		mcp.WithDescription("Get achievement details (name, description, requirements) for given achievement IDs."),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of achievement IDs to look up"),
		),
	)
	s.mcp.AddTool(achievementsTool, s.handleGetAchievements)

	dailyAchievementsTool := mcp.NewTool(
		"get_daily_achievements",
		mcp.WithDescription("Get today's and tomorrow's daily achievements."),
	)
	s.mcp.AddTool(dailyAchievementsTool, s.handleGetDailyAchievements)

	// --- Guild Tools ---

	guildTool := mcp.NewTool(
		"get_guild",
		mcp.WithDescription("Get public guild information (name, tag, level)."),
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Description("Guild ID (UUID)"),
		),
	)
	s.mcp.AddTool(guildTool, s.handleGetGuild)

	searchGuildTool := mcp.NewTool(
		"search_guild",
		mcp.WithDescription("Search for a guild by name. Returns matching guild IDs."),
		mcp.WithString(
			"name",
			mcp.Required(),
			mcp.Description("Guild name to search for"),
		),
	)
	s.mcp.AddTool(searchGuildTool, s.handleSearchGuild)

	guildDetailsTool := mcp.NewTool(
		"get_guild_details",
		mcp.WithDescription("Get detailed guild data (log, members, ranks, stash, etc.). Requires GW2_API_KEY with guild leader permissions."),
		mcp.WithString(
			"id",
			mcp.Required(),
			mcp.Description("Guild ID (UUID)"),
		),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Detail type: log, members, ranks, stash, storage, treasury, teams, upgrades"),
		),
	)
	s.mcp.AddTool(guildDetailsTool, s.handleGetGuildDetails)

	// --- Game Metadata ---

	colorsTool := mcp.NewTool(
		"get_colors",
		mcp.WithDescription("Get dye color metadata (name) for given color IDs."),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of color IDs to look up"),
		),
	)
	s.mcp.AddTool(colorsTool, s.handleGetColors)

	minisTool := mcp.NewTool(
		"get_minis",
		mcp.WithDescription("Get miniature metadata (name, icon, item_id) for given mini IDs."),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of mini IDs to look up"),
		),
	)
	s.mcp.AddTool(minisTool, s.handleGetMinis)

	mountsInfoTool := mcp.NewTool(
		"get_mounts_info",
		mcp.WithDescription("Get mount skin or mount type metadata for given IDs."),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Mount info type: 'skins' or 'types'"),
		),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of mount skin or type IDs to look up"),
		),
	)
	s.mcp.AddTool(mountsInfoTool, s.handleGetMountsInfo)

	gameBuildTool := mcp.NewTool(
		"get_game_build",
		mcp.WithDescription("Get the current Guild Wars 2 game build number."),
	)
	s.mcp.AddTool(gameBuildTool, s.handleGetGameBuild)

	tokenInfoTool := mcp.NewTool(
		"get_token_info",
		mcp.WithDescription("Get API key information including name and permission scopes. Requires GW2_API_KEY."),
	)
	s.mcp.AddTool(tokenInfoTool, s.handleGetTokenInfo)

	dungeonsRaidsTool := mcp.NewTool(
		"get_dungeons_and_raids",
		mcp.WithDescription("Get dungeon or raid metadata (paths, wings, events) for given IDs."),
		mcp.WithString(
			"type",
			mcp.Required(),
			mcp.Description("Content type: 'dungeons' or 'raids'"),
		),
		mcp.WithArray(
			"ids",
			mcp.Required(),
			mcp.Description("Array of dungeon or raid IDs (e.g., 'ascalonian_catacombs', 'forsaken_thicket')"),
		),
	)
	s.mcp.AddTool(dungeonsRaidsTool, s.handleGetDungeonsAndRaids)
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
