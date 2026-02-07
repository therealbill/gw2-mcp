package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// handleWikiSearch handles wiki search requests
func (s *MCPServer) handleWikiSearch(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid query parameter: %v", err)), nil
	}

	// Get limit parameter (optional)
	const defaultLimit = 5
	limit := request.GetInt("limit", defaultLimit)

	s.logger.Debug("Wiki search request", "query", query, "limit", limit)

	// Perform wiki search
	results, err := s.wiki.Search(ctx, query, limit)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Wiki search failed: %v", err)), nil
	}

	// Format results as JSON
	resultJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(resultJSON)), nil
}

// handleGetWallet handles wallet information requests
func (s *MCPServer) handleGetWallet(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Wallet request")

	// Get wallet information
	wallet, err := s.gw2API.GetWallet(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get wallet: %v", err)), nil
	}

	// Format wallet as JSON
	walletJSON, err := json.MarshalIndent(wallet, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format wallet: %v", err)), nil
	}

	return mcp.NewToolResultText(string(walletJSON)), nil
}

// handleGetCurrencies handles currency information requests
func (s *MCPServer) handleGetCurrencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse optional currency IDs
	currencyIDs := request.GetIntSlice("ids", nil)

	s.logger.Debug("Currency request", "currency_ids", currencyIDs)

	// Get currency information
	currencies, err := s.gw2API.GetCurrencies(ctx, currencyIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get currencies: %v", err)), nil
	}

	// Format currencies as JSON
	currenciesJSON, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format currencies: %v", err)), nil
	}

	return mcp.NewToolResultText(string(currenciesJSON)), nil
}

// handleGetTPPrices handles trading post price requests
func (s *MCPServer) handleGetTPPrices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	itemIDs := request.GetIntSlice("item_ids", nil)
	if len(itemIDs) == 0 {
		return mcp.NewToolResultError("item_ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("TP prices request", "item_ids", itemIDs)

	prices, err := s.gw2API.GetPrices(ctx, itemIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get trading post prices: %v", err)), nil
	}

	pricesJSON, err := json.MarshalIndent(prices, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format prices: %v", err)), nil
	}

	return mcp.NewToolResultText(string(pricesJSON)), nil
}

// handleGetTPListings handles trading post listing requests
func (s *MCPServer) handleGetTPListings(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	itemIDs := request.GetIntSlice("item_ids", nil)
	if len(itemIDs) == 0 {
		return mcp.NewToolResultError("item_ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("TP listings request", "item_ids", itemIDs)

	listings, err := s.gw2API.GetListings(ctx, itemIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get trading post listings: %v", err)), nil
	}

	listingsJSON, err := json.MarshalIndent(listings, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format listings: %v", err)), nil
	}

	return mcp.NewToolResultText(string(listingsJSON)), nil
}

// handleGetGemExchange handles gem exchange rate requests
func (s *MCPServer) handleGetGemExchange(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	direction, err := request.RequireString("direction")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid direction parameter: %v", err)), nil
	}

	quantity := request.GetInt("quantity", 0)
	if quantity <= 0 {
		return mcp.NewToolResultError("quantity must be greater than 0"), nil
	}

	s.logger.Debug("Gem exchange request", "direction", direction, "quantity", quantity)

	exchange, err := s.gw2API.GetGemExchange(ctx, direction, quantity)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get gem exchange rate: %v", err)), nil
	}

	exchangeJSON, err := json.MarshalIndent(exchange, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format exchange rate: %v", err)), nil
	}

	return mcp.NewToolResultText(string(exchangeJSON)), nil
}

// handleGetTPDelivery handles trading post delivery box requests
func (s *MCPServer) handleGetTPDelivery(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("TP delivery request")

	delivery, err := s.gw2API.GetDelivery(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get trading post delivery: %v", err)), nil
	}

	deliveryJSON, err := json.MarshalIndent(delivery, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format delivery: %v", err)), nil
	}

	return mcp.NewToolResultText(string(deliveryJSON)), nil
}

// handleGetTPTransactions handles trading post transaction requests
func (s *MCPServer) handleGetTPTransactions(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	txType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	s.logger.Debug("TP transactions request", "type", txType)

	transactions, err := s.gw2API.GetTransactions(ctx, txType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get trading post transactions: %v", err)), nil
	}

	transactionsJSON, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format transactions: %v", err)), nil
	}

	return mcp.NewToolResultText(string(transactionsJSON)), nil
}

// --- Phase 2: Account Handlers ---

// handleGetAccount handles account information requests
func (s *MCPServer) handleGetAccount(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Account request")

	account, err := s.gw2API.GetAccount(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get account: %v", err)), nil
	}

	data, err := json.MarshalIndent(account, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format account: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetBank handles bank vault requests
func (s *MCPServer) handleGetBank(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Bank request")

	bank, err := s.gw2API.GetBank(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get bank: %v", err)), nil
	}

	data, err := json.MarshalIndent(bank, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format bank: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetMaterials handles material storage requests
func (s *MCPServer) handleGetMaterials(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Materials request")

	materials, err := s.gw2API.GetMaterials(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get materials: %v", err)), nil
	}

	data, err := json.MarshalIndent(materials, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format materials: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetInventory handles shared inventory requests
func (s *MCPServer) handleGetInventory(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Shared inventory request")

	inventory, err := s.gw2API.GetSharedInventory(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get shared inventory: %v", err)), nil
	}

	data, err := json.MarshalIndent(inventory, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format shared inventory: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetCharacters handles character list/detail requests
func (s *MCPServer) handleGetCharacters(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")

	if name != "" {
		s.logger.Debug("Character detail request", "name", name)
		character, err := s.gw2API.GetCharacter(ctx, name)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get character: %v", err)), nil
		}
		data, err := json.MarshalIndent(character, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format character: %v", err)), nil
		}
		return mcp.NewToolResultText(string(data)), nil
	}

	s.logger.Debug("Characters list request")
	characters, err := s.gw2API.GetCharacters(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get characters: %v", err)), nil
	}

	data, err := json.MarshalIndent(characters, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format characters: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// --- Phase 3-5: Account Unlocks, Progress, Dailies Handlers ---

// handleGetAccountUnlocks handles account unlock requests
func (s *MCPServer) handleGetAccountUnlocks(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	unlockType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	s.logger.Debug("Account unlocks request", "type", unlockType)

	unlocks, err := s.gw2API.GetAccountUnlocks(ctx, unlockType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get account unlocks: %v", err)), nil
	}

	return mcp.NewToolResultText(string(unlocks)), nil
}

// handleGetAccountProgress handles account progress requests
func (s *MCPServer) handleGetAccountProgress(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	progressType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	s.logger.Debug("Account progress request", "type", progressType)

	progress, err := s.gw2API.GetAccountProgress(ctx, progressType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get account progress: %v", err)), nil
	}

	return mcp.NewToolResultText(string(progress)), nil
}

// handleGetAccountDailies handles account dailies requests
func (s *MCPServer) handleGetAccountDailies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dailyType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	s.logger.Debug("Account dailies request", "type", dailyType)

	dailies, err := s.gw2API.GetAccountDailies(ctx, dailyType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get account dailies: %v", err)), nil
	}

	return mcp.NewToolResultText(string(dailies)), nil
}

// --- Phase 6: Wizard's Vault Handlers ---

// handleGetWizardsVault handles wizard's vault season info requests
func (s *MCPServer) handleGetWizardsVault(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Wizard's vault season request")

	data, err := s.gw2API.GetWizardsVault(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get wizard's vault: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// handleGetWizardsVaultObjectives handles wizard's vault objectives requests
func (s *MCPServer) handleGetWizardsVaultObjectives(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	objType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	s.logger.Debug("Wizard's vault objectives request", "type", objType)

	data, err := s.gw2API.GetWizardsVaultObjectives(ctx, objType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get wizard's vault objectives: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// handleGetWizardsVaultListings handles wizard's vault listings requests
func (s *MCPServer) handleGetWizardsVaultListings(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Wizard's vault listings request")

	data, err := s.gw2API.GetWizardsVaultListings(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get wizard's vault listings: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// --- Phase 7: Game Data Handlers ---

// handleGetItems handles item lookup requests
func (s *MCPServer) handleGetItems(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	itemIDs := request.GetIntSlice("ids", nil)
	if len(itemIDs) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Items request", "ids", itemIDs)

	items, err := s.gw2API.GetItems(ctx, itemIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get items: %v", err)), nil
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format items: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetSkins handles skin lookup requests
func (s *MCPServer) handleGetSkins(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	skinIDs := request.GetIntSlice("ids", nil)
	if len(skinIDs) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Skins request", "ids", skinIDs)

	skins, err := s.gw2API.GetSkins(ctx, skinIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get skins: %v", err)), nil
	}

	data, err := json.MarshalIndent(skins, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format skins: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetRecipes handles recipe lookup requests
func (s *MCPServer) handleGetRecipes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	recipeIDs := request.GetIntSlice("ids", nil)
	if len(recipeIDs) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Recipes request", "ids", recipeIDs)

	recipes, err := s.gw2API.GetRecipes(ctx, recipeIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get recipes: %v", err)), nil
	}

	data, err := json.MarshalIndent(recipes, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format recipes: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleSearchRecipes handles recipe search requests
func (s *MCPServer) handleSearchRecipes(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	input := request.GetInt("input", 0)
	output := request.GetInt("output", 0)

	if input == 0 && output == 0 {
		return mcp.NewToolResultError("either input or output item ID must be provided"), nil
	}

	s.logger.Debug("Recipe search request", "input", input, "output", output)

	ids, err := s.gw2API.SearchRecipes(ctx, input, output)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search recipes: %v", err)), nil
	}

	data, err := json.MarshalIndent(ids, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format recipe search results: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetAchievements handles achievement lookup requests
func (s *MCPServer) handleGetAchievements(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	achIDs := request.GetIntSlice("ids", nil)
	if len(achIDs) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Achievements request", "ids", achIDs)

	achievements, err := s.gw2API.GetAchievements(ctx, achIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get achievements: %v", err)), nil
	}

	data, err := json.MarshalIndent(achievements, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format achievements: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetDailyAchievements handles daily achievement requests
func (s *MCPServer) handleGetDailyAchievements(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Daily achievements request")

	dailies, err := s.gw2API.GetDailyAchievements(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get daily achievements: %v", err)), nil
	}

	data, err := json.MarshalIndent(dailies, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format daily achievements: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// --- Phase 8: Guild Handlers ---

// handleGetGuild handles public guild info requests
func (s *MCPServer) handleGetGuild(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	guildID, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid id parameter: %v", err)), nil
	}

	s.logger.Debug("Guild info request", "id", guildID)

	guild, err := s.gw2API.GetGuild(ctx, guildID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get guild: %v", err)), nil
	}

	data, err := json.MarshalIndent(guild, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format guild: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleSearchGuild handles guild search requests
func (s *MCPServer) handleSearchGuild(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid name parameter: %v", err)), nil
	}

	s.logger.Debug("Guild search request", "name", name)

	ids, err := s.gw2API.SearchGuild(ctx, name)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to search guild: %v", err)), nil
	}

	data, err := json.MarshalIndent(ids, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format guild search results: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetGuildDetails handles authenticated guild detail requests
func (s *MCPServer) handleGetGuildDetails(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	guildID, err := request.RequireString("id")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid id parameter: %v", err)), nil
	}

	detailType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	s.logger.Debug("Guild detail request", "id", guildID, "type", detailType)

	details, err := s.gw2API.GetGuildDetails(ctx, guildID, detailType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get guild details: %v", err)), nil
	}

	return mcp.NewToolResultText(string(details)), nil
}

// --- Phase 9: Game Metadata Handlers ---

// handleGetColors handles color lookup requests
func (s *MCPServer) handleGetColors(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	colorIDs := request.GetIntSlice("ids", nil)
	if len(colorIDs) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Colors request", "ids", colorIDs)

	colors, err := s.gw2API.GetColors(ctx, colorIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get colors: %v", err)), nil
	}

	data, err := json.MarshalIndent(colors, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format colors: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetMinis handles mini lookup requests
func (s *MCPServer) handleGetMinis(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	miniIDs := request.GetIntSlice("ids", nil)
	if len(miniIDs) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Minis request", "ids", miniIDs)

	minis, err := s.gw2API.GetMinis(ctx, miniIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get minis: %v", err)), nil
	}

	data, err := json.MarshalIndent(minis, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format minis: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetMountsInfo handles mount info requests
func (s *MCPServer) handleGetMountsInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	mountType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	mountIDs := request.GetIntSlice("ids", nil)
	if len(mountIDs) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Mounts info request", "type", mountType, "ids", mountIDs)

	data, err := s.gw2API.GetMountsInfo(ctx, mountType, mountIDs)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get mount info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// handleGetGameBuild handles game build number requests
func (s *MCPServer) handleGetGameBuild(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Game build request")

	build, err := s.gw2API.GetGameBuild(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get game build: %v", err)), nil
	}

	data, err := json.MarshalIndent(build, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format game build: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetTokenInfo handles API token info requests
func (s *MCPServer) handleGetTokenInfo(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	s.logger.Debug("Token info request")

	info, err := s.gw2API.GetTokenInfo(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get token info: %v", err)), nil
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format token info: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

// handleGetDungeonsAndRaids handles dungeon/raid metadata requests
func (s *MCPServer) handleGetDungeonsAndRaids(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	contentType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	idsRaw := request.GetStringSlice("ids", nil)
	if len(idsRaw) == 0 {
		return mcp.NewToolResultError("ids parameter is required and must not be empty"), nil
	}

	s.logger.Debug("Dungeons/raids request", "type", contentType, "ids", idsRaw)

	data, err := s.gw2API.GetDungeonsAndRaids(ctx, contentType, idsRaw)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get %s: %v", contentType, err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// handleCurrencyListResource handles the currency list resource
func (s *MCPServer) handleCurrencyListResource(ctx context.Context,
	_ mcp.ReadResourceRequest,
) ([]mcp.ResourceContents, error) {
	s.logger.Debug("Currency list resource request")

	// Get all currencies
	currencies, err := s.gw2API.GetCurrencies(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get currencies: %w", err)
	}

	// Format currencies as JSON
	currenciesJSON, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format currencies: %w", err)
	}

	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "gw2://currencies",
			MIMEType: "application/json",
			Text:     string(currenciesJSON),
		},
	}, nil
}
