// Package server provides the MCP server implementation for Guild Wars 2 data access.
package server

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/AlyxPink/gw2-mcp/internal/cache"
	"github.com/AlyxPink/gw2-mcp/internal/gw2api"
	"github.com/AlyxPink/gw2-mcp/internal/wiki"

	"github.com/charmbracelet/log"
)

// MCPServer wraps the MCP server with GW2-specific functionality
type MCPServer struct {
	mcp    *mcp.Server
	logger *log.Logger
	cache  *cache.Manager
	gw2API *gw2api.Client
	wiki   *wiki.Client
}

// --- Argument structs for tools with parameters ---

type WikiSearchArgs struct {
	Query string `json:"query" jsonschema:"Search query for wiki content (e.g. 'Dragon Bash', 'currencies', 'wallet')"`
	Limit int    `json:"limit,omitempty" jsonschema:"Maximum number of results to return (default: 5)"`
}

type GetCurrenciesArgs struct {
	IDs []int `json:"ids,omitempty" jsonschema:"Specific currency IDs to fetch (optional, returns all if not specified)"`
}

type GetTPPricesArgs struct {
	ItemIDs []int `json:"item_ids" jsonschema:"Array of item IDs to get prices for (e.g. [19976, 19721] for Mystic Coin and Glob of Ectoplasm)"`
}

type GetTPListingsArgs struct {
	ItemIDs []int `json:"item_ids" jsonschema:"Array of item IDs to get listings for"`
}

type GetGemExchangeArgs struct {
	Direction string `json:"direction" jsonschema:"Exchange direction: 'coins' (coins to gems) or 'gems' (gems to coins)"`
	Quantity  int    `json:"quantity" jsonschema:"Amount to convert (coins in copper, or number of gems)"`
}

type GetTPTransactionsArgs struct {
	Type string `json:"type" jsonschema:"Transaction type: 'current/buys', 'current/sells', 'history/buys', or 'history/sells'"`
}

type GetCharactersArgs struct {
	Name string `json:"name,omitempty" jsonschema:"Character name to get details for (optional; omit to list all characters)"`
}

type GetAccountUnlocksArgs struct {
	Type string `json:"type" jsonschema:"Unlock type: skins, dyes, minis, titles, recipes, finishers, outfits, gliders, mailcarriers, novelties, emotes, mounts/skins, mounts/types, skiffs, jadebots"`
}

type GetAccountProgressArgs struct {
	Type string `json:"type" jsonschema:"Progress type: achievements, masteries, mastery/points, luck, legendaryarmory, progression"`
}

type GetAccountDailiesArgs struct {
	Type string `json:"type" jsonschema:"Daily type: dailycrafting, dungeons, raids, mapchests, worldbosses"`
}

type GetWizardsVaultObjectivesArgs struct {
	Type string `json:"type" jsonschema:"Objective type: daily, weekly, special"`
}

type GetItemsArgs struct {
	IDs []int `json:"ids" jsonschema:"Array of item IDs to look up"`
}

type GetSkinsArgs struct {
	IDs []int `json:"ids" jsonschema:"Array of skin IDs to look up"`
}

type GetRecipesArgs struct {
	IDs []int `json:"ids" jsonschema:"Array of recipe IDs to look up"`
}

type SearchRecipesArgs struct {
	Input  int `json:"input,omitempty" jsonschema:"Input item ID to search for recipes that use this item"`
	Output int `json:"output,omitempty" jsonschema:"Output item ID to search for recipes that produce this item"`
}

type GetAchievementsArgs struct {
	IDs []int `json:"ids" jsonschema:"Array of achievement IDs to look up"`
}

type GetGuildArgs struct {
	ID string `json:"id" jsonschema:"Guild ID (UUID)"`
}

type SearchGuildArgs struct {
	Name string `json:"name" jsonschema:"Guild name to search for"`
}

type GetGuildDetailsArgs struct {
	ID   string `json:"id" jsonschema:"Guild ID (UUID)"`
	Type string `json:"type" jsonschema:"Detail type: log, members, ranks, stash, storage, treasury, teams, upgrades"`
}

type GetColorsArgs struct {
	IDs []int `json:"ids" jsonschema:"Array of color IDs to look up"`
}

type GetMinisArgs struct {
	IDs []int `json:"ids" jsonschema:"Array of mini IDs to look up"`
}

type GetMountsInfoArgs struct {
	Type string `json:"type" jsonschema:"Mount info type: 'skins' or 'types'"`
	IDs  []int  `json:"ids" jsonschema:"Array of mount skin or type IDs to look up"`
}

type GetDungeonsAndRaidsArgs struct {
	Type string   `json:"type" jsonschema:"Content type: 'dungeons' or 'raids'"`
	IDs  []string `json:"ids" jsonschema:"Array of dungeon or raid IDs (e.g. 'ascalonian_catacombs', 'forsaken_thicket')"`
}

type GetItemByNameArgs struct {
	Name string `json:"name" jsonschema:"Item name to search for (e.g. 'Mystic Coin', 'Dusk')"`
}

type GetItemRecipeByNameArgs struct {
	Name string `json:"name" jsonschema:"Item name to find recipes for (e.g. '18 Slot Silk Bag', 'Dawn')"`
}

type GetTPPriceByNameArgs struct {
	Name string `json:"name" jsonschema:"Item name to get trading post prices for (e.g. 'Glob of Ectoplasm', 'Mystic Coin')"`
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
	mcpServer := mcp.NewServer(
		&mcp.Implementation{
			Name:    "GW2 MCP Server",
			Version: "1.0.0",
		},
		nil,
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
	return s.mcp.Run(ctx, &mcp.StdioTransport{})
}

// registerTools registers all available tools
func (s *MCPServer) registerTools() {
	// Wiki search tool
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "wiki_search",
		Description: "Search Guild Wars 2 wiki for information about game content",
	}, s.handleWikiSearch)

	// Wallet info tool (no params)
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_wallet",
		Description: "Get user's wallet information including all currencies. Requires GW2_API_KEY environment variable.",
	}, s.handleGetWallet)

	// Currency info tool
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_currencies",
		Description: "Get information about Guild Wars 2 currencies",
	}, s.handleGetCurrencies)

	// Trading Post prices tool
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_tp_prices",
		Description: "Get Trading Post prices for items. Returns aggregated best buy/sell prices with item names and formatted coin values.",
	}, s.handleGetTPPrices)

	// Trading Post listings tool
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_tp_listings",
		Description: "Get Trading Post order book listings for items. Returns all buy/sell price tiers with quantities.",
	}, s.handleGetTPListings)

	// Gem exchange tool
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_gem_exchange",
		Description: "Get gem exchange rates. Convert coins to gems or gems to coins.",
	}, s.handleGetGemExchange)

	// Trading Post delivery tool (no params)
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_tp_delivery",
		Description: "Get items and coins awaiting pickup from the Trading Post. Requires GW2_API_KEY with account and tradingpost scopes.",
	}, s.handleGetTPDelivery)

	// Trading Post transactions tool
	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_tp_transactions",
		Description: "Get Trading Post transaction history. View current orders or completed transactions from the past 90 days. Requires GW2_API_KEY with account and tradingpost scopes.",
	}, s.handleGetTPTransactions)

	// --- Account Tools ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_account",
		Description: "Get account information including name, world, guilds, and access. Requires GW2_API_KEY.",
	}, s.handleGetAccount)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_bank",
		Description: "Get bank vault contents with item names. Requires GW2_API_KEY.",
	}, s.handleGetBank)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_materials",
		Description: "Get material storage contents with item names. Requires GW2_API_KEY.",
	}, s.handleGetMaterials)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_inventory",
		Description: "Get shared inventory slot contents with item names. Requires GW2_API_KEY.",
	}, s.handleGetInventory)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_characters",
		Description: "Get list of character names, or detailed info for a specific character including crafting disciplines, equipment, skills, specializations, and build tabs. Requires GW2_API_KEY.",
	}, s.handleGetCharacters)

	// --- Account Unlocks ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_account_unlocks",
		Description: "Get account unlocked IDs. Requires GW2_API_KEY.",
	}, s.handleGetAccountUnlocks)

	// --- Account Progress ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_account_progress",
		Description: "Get account progress data. Requires GW2_API_KEY.",
	}, s.handleGetAccountProgress)

	// --- Account Dailies ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_account_dailies",
		Description: "Get completed daily content IDs. Requires GW2_API_KEY.",
	}, s.handleGetAccountDailies)

	// --- Wizard's Vault ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_wizards_vault",
		Description: "Get current Wizard's Vault season information.",
	}, s.handleGetWizardsVault)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_wizards_vault_objectives",
		Description: "Get Wizard's Vault objectives. Uses authenticated endpoint if GW2_API_KEY is set, otherwise returns public objective list.",
	}, s.handleGetWizardsVaultObjectives)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_wizards_vault_listings",
		Description: "Get Wizard's Vault reward listings. Uses authenticated endpoint if GW2_API_KEY is set.",
	}, s.handleGetWizardsVaultListings)

	// --- Game Data Lookups ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_items",
		Description: "Get item metadata (name, type, rarity, level, icon, description, vendor value, flags, game types, restrictions, and type-specific details) for given item IDs.",
	}, s.handleGetItems)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_skins",
		Description: "Get skin metadata (name, type, icon, rarity, description, flags, restrictions, and type-specific details) for given skin IDs.",
	}, s.handleGetSkins)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_recipes",
		Description: "Get recipe details (type, output, ingredients, disciplines, crafting time, flags, guild ingredients, chat link) for given recipe IDs.",
	}, s.handleGetRecipes)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "search_recipes",
		Description: "Search for recipe IDs by input or output item ID.",
	}, s.handleSearchRecipes)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_achievements",
		Description: "Get achievement details (name, description, requirements, tiers, prerequisites, rewards, bits, icon) for given achievement IDs.",
	}, s.handleGetAchievements)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_daily_achievements",
		Description: "Get today's and tomorrow's daily achievements.",
	}, s.handleGetDailyAchievements)

	// --- Guild Tools ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_guild",
		Description: "Get public guild information (name, tag, level).",
	}, s.handleGetGuild)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "search_guild",
		Description: "Search for a guild by name. Returns matching guild IDs.",
	}, s.handleSearchGuild)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_guild_details",
		Description: "Get detailed guild data (log, members, ranks, stash, etc.). Requires GW2_API_KEY with guild leader permissions.",
	}, s.handleGetGuildDetails)

	// --- Game Metadata ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_colors",
		Description: "Get dye color metadata (name, base RGB, cloth/leather/metal/fur material adjustments) for given color IDs.",
	}, s.handleGetColors)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_minis",
		Description: "Get miniature metadata (name, icon, item_id) for given mini IDs.",
	}, s.handleGetMinis)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_mounts_info",
		Description: "Get mount skin or mount type metadata for given IDs.",
	}, s.handleGetMountsInfo)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_game_build",
		Description: "Get the current Guild Wars 2 game build number.",
	}, s.handleGetGameBuild)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_token_info",
		Description: "Get API key information including name and permission scopes. Requires GW2_API_KEY.",
	}, s.handleGetTokenInfo)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_dungeons_and_raids",
		Description: "Get dungeon or raid metadata (paths, wings, events) for given IDs.",
	}, s.handleGetDungeonsAndRaids)

	// --- Composite Tools ---

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_item_by_name",
		Description: "Look up a GW2 item by name. Searches the wiki to find the item ID, then returns full item details from the API.",
	}, s.handleGetItemByName)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_item_recipe_by_name",
		Description: "Find crafting recipes for a GW2 item by name. Searches the wiki for recipe data, then returns full recipe details with resolved ingredient names.",
	}, s.handleGetItemRecipeByName)

	mcp.AddTool(s.mcp, &mcp.Tool{
		Name:        "get_tp_price_by_name",
		Description: "Get Trading Post prices for an item by name. Searches the wiki to find the item ID, then returns current buy/sell prices.",
	}, s.handleGetTPPriceByName)
}

// registerResources registers all available resources
func (s *MCPServer) registerResources() {
	s.mcp.AddResource(&mcp.Resource{
		URI:         "gw2://currencies",
		Name:        "Guild Wars 2 Currencies",
		Description: "Complete list of all Guild Wars 2 currencies with metadata",
		MIMEType:    "application/json",
	}, s.handleCurrencyListResource)
}
