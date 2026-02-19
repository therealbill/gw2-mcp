package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/AlyxPink/gw2-mcp/internal/gw2api"
	"github.com/AlyxPink/gw2-mcp/internal/wiki"
)

// textResult is a helper to build a text CallToolResult.
func textResult(text string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}, nil, nil
}

// errResult is a helper to build an error CallToolResult visible to the LLM.
func errResult(msg string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}, nil, nil
}

// jsonResult marshals v to indented JSON and returns it as a text result.
func jsonResult(v any) (*mcp.CallToolResult, any, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return errResult(fmt.Sprintf("Failed to format response: %v", err))
	}
	return textResult(string(data))
}

// handleWikiSearch handles wiki search requests
func (s *MCPServer) handleWikiSearch(ctx context.Context, _ *mcp.CallToolRequest, args WikiSearchArgs) (*mcp.CallToolResult, any, error) {
	if args.Query == "" {
		return errResult("query parameter is required")
	}

	limit := args.Limit
	if limit == 0 {
		limit = 5
	}

	s.logger.Debug("Wiki search request", "query", args.Query, "limit", limit)

	results, err := s.wiki.Search(ctx, args.Query, limit)
	if err != nil {
		return errResult(fmt.Sprintf("Wiki search failed: %v", err))
	}

	return jsonResult(results)
}

// handleGetWallet handles wallet information requests
func (s *MCPServer) handleGetWallet(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Wallet request")

	wallet, err := s.gw2API.GetWallet(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get wallet: %v", err))
	}

	return jsonResult(wallet)
}

// handleGetCurrencies handles currency information requests
func (s *MCPServer) handleGetCurrencies(ctx context.Context, _ *mcp.CallToolRequest, args GetCurrenciesArgs) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Currency request", "currency_ids", args.IDs)

	currencies, err := s.gw2API.GetCurrencies(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get currencies: %v", err))
	}

	return jsonResult(currencies)
}

// handleGetTPPrices handles trading post price requests
func (s *MCPServer) handleGetTPPrices(ctx context.Context, _ *mcp.CallToolRequest, args GetTPPricesArgs) (*mcp.CallToolResult, any, error) {
	if len(args.ItemIDs) == 0 {
		return errResult("item_ids parameter is required and must not be empty")
	}

	s.logger.Debug("TP prices request", "item_ids", args.ItemIDs)

	prices, err := s.gw2API.GetPrices(ctx, args.ItemIDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get trading post prices: %v", err))
	}

	return jsonResult(prices)
}

// handleGetTPListings handles trading post listing requests
func (s *MCPServer) handleGetTPListings(ctx context.Context, _ *mcp.CallToolRequest, args GetTPListingsArgs) (*mcp.CallToolResult, any, error) {
	if len(args.ItemIDs) == 0 {
		return errResult("item_ids parameter is required and must not be empty")
	}

	s.logger.Debug("TP listings request", "item_ids", args.ItemIDs)

	listings, err := s.gw2API.GetListings(ctx, args.ItemIDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get trading post listings: %v", err))
	}

	return jsonResult(listings)
}

// handleGetGemExchange handles gem exchange rate requests
func (s *MCPServer) handleGetGemExchange(ctx context.Context, _ *mcp.CallToolRequest, args GetGemExchangeArgs) (*mcp.CallToolResult, any, error) {
	if args.Direction == "" {
		return errResult("direction parameter is required")
	}
	if args.Quantity <= 0 {
		return errResult("quantity must be greater than 0")
	}

	s.logger.Debug("Gem exchange request", "direction", args.Direction, "quantity", args.Quantity)

	exchange, err := s.gw2API.GetGemExchange(ctx, args.Direction, args.Quantity)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get gem exchange rate: %v", err))
	}

	return jsonResult(exchange)
}

// handleGetTPDelivery handles trading post delivery box requests
func (s *MCPServer) handleGetTPDelivery(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("TP delivery request")

	delivery, err := s.gw2API.GetDelivery(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get trading post delivery: %v", err))
	}

	return jsonResult(delivery)
}

// handleGetTPTransactions handles trading post transaction requests
func (s *MCPServer) handleGetTPTransactions(ctx context.Context, _ *mcp.CallToolRequest, args GetTPTransactionsArgs) (*mcp.CallToolResult, any, error) {
	if args.Type == "" {
		return errResult("type parameter is required")
	}

	s.logger.Debug("TP transactions request", "type", args.Type)

	transactions, err := s.gw2API.GetTransactions(ctx, args.Type)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get trading post transactions: %v", err))
	}

	return jsonResult(transactions)
}

// --- Account Handlers ---

// handleGetAccount handles account information requests
func (s *MCPServer) handleGetAccount(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Account request")

	account, err := s.gw2API.GetAccount(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get account: %v", err))
	}

	return jsonResult(account)
}

// handleGetBank handles bank vault requests
func (s *MCPServer) handleGetBank(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Bank request")

	bank, err := s.gw2API.GetBank(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get bank: %v", err))
	}

	return jsonResult(bank)
}

// handleGetMaterials handles material storage requests
func (s *MCPServer) handleGetMaterials(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Materials request")

	materials, err := s.gw2API.GetMaterials(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get materials: %v", err))
	}

	return jsonResult(materials)
}

// handleGetInventory handles shared inventory requests
func (s *MCPServer) handleGetInventory(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Shared inventory request")

	inventory, err := s.gw2API.GetSharedInventory(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get shared inventory: %v", err))
	}

	return jsonResult(inventory)
}

// handleGetCharacters handles character list/detail requests
func (s *MCPServer) handleGetCharacters(ctx context.Context, _ *mcp.CallToolRequest, args GetCharactersArgs) (*mcp.CallToolResult, any, error) {
	if args.Name != "" {
		s.logger.Debug("Character detail request", "name", args.Name)
		character, err := s.gw2API.GetCharacter(ctx, args.Name)
		if err != nil {
			return errResult(fmt.Sprintf("Failed to get character: %v", err))
		}
		return jsonResult(character)
	}

	s.logger.Debug("Characters list request")
	characters, err := s.gw2API.GetCharacters(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get characters: %v", err))
	}

	return jsonResult(characters)
}

// --- Account Unlocks, Progress, Dailies Handlers ---

// handleGetAccountUnlocks handles account unlock requests
func (s *MCPServer) handleGetAccountUnlocks(ctx context.Context, _ *mcp.CallToolRequest, args GetAccountUnlocksArgs) (*mcp.CallToolResult, any, error) {
	if args.Type == "" {
		return errResult("type parameter is required")
	}

	s.logger.Debug("Account unlocks request", "type", args.Type)

	unlocks, err := s.gw2API.GetAccountUnlocks(ctx, args.Type)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get account unlocks: %v", err))
	}

	return textResult(string(unlocks))
}

// handleGetAccountProgress handles account progress requests
func (s *MCPServer) handleGetAccountProgress(ctx context.Context, _ *mcp.CallToolRequest, args GetAccountProgressArgs) (*mcp.CallToolResult, any, error) {
	if args.Type == "" {
		return errResult("type parameter is required")
	}

	s.logger.Debug("Account progress request", "type", args.Type)

	progress, err := s.gw2API.GetAccountProgress(ctx, args.Type)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get account progress: %v", err))
	}

	return textResult(string(progress))
}

// handleGetAccountDailies handles account dailies requests
func (s *MCPServer) handleGetAccountDailies(ctx context.Context, _ *mcp.CallToolRequest, args GetAccountDailiesArgs) (*mcp.CallToolResult, any, error) {
	if args.Type == "" {
		return errResult("type parameter is required")
	}

	s.logger.Debug("Account dailies request", "type", args.Type)

	dailies, err := s.gw2API.GetAccountDailies(ctx, args.Type)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get account dailies: %v", err))
	}

	return textResult(string(dailies))
}

// --- Wizard's Vault Handlers ---

// handleGetWizardsVault handles wizard's vault season info requests
func (s *MCPServer) handleGetWizardsVault(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Wizard's vault season request")

	data, err := s.gw2API.GetWizardsVault(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get wizard's vault: %v", err))
	}

	return textResult(string(data))
}

// handleGetWizardsVaultObjectives handles wizard's vault objectives requests
func (s *MCPServer) handleGetWizardsVaultObjectives(ctx context.Context, _ *mcp.CallToolRequest, args GetWizardsVaultObjectivesArgs) (*mcp.CallToolResult, any, error) {
	if args.Type == "" {
		return errResult("type parameter is required")
	}

	s.logger.Debug("Wizard's vault objectives request", "type", args.Type)

	data, err := s.gw2API.GetWizardsVaultObjectives(ctx, args.Type)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get wizard's vault objectives: %v", err))
	}

	return textResult(string(data))
}

// handleGetWizardsVaultListings handles wizard's vault listings requests
func (s *MCPServer) handleGetWizardsVaultListings(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Wizard's vault listings request")

	data, err := s.gw2API.GetWizardsVaultListings(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get wizard's vault listings: %v", err))
	}

	return textResult(string(data))
}

// --- Game Data Handlers ---

// handleGetItems handles item lookup requests
func (s *MCPServer) handleGetItems(ctx context.Context, _ *mcp.CallToolRequest, args GetItemsArgs) (*mcp.CallToolResult, any, error) {
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Items request", "ids", args.IDs)

	items, err := s.gw2API.GetItems(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get items: %v", err))
	}

	return jsonResult(items)
}

// handleGetSkins handles skin lookup requests
func (s *MCPServer) handleGetSkins(ctx context.Context, _ *mcp.CallToolRequest, args GetSkinsArgs) (*mcp.CallToolResult, any, error) {
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Skins request", "ids", args.IDs)

	skins, err := s.gw2API.GetSkins(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get skins: %v", err))
	}

	return jsonResult(skins)
}

// handleGetRecipes handles recipe lookup requests
func (s *MCPServer) handleGetRecipes(ctx context.Context, _ *mcp.CallToolRequest, args GetRecipesArgs) (*mcp.CallToolResult, any, error) {
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Recipes request", "ids", args.IDs)

	recipes, err := s.gw2API.GetRecipes(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get recipes: %v", err))
	}

	return jsonResult(recipes)
}

// handleSearchRecipes handles recipe search requests
func (s *MCPServer) handleSearchRecipes(ctx context.Context, _ *mcp.CallToolRequest, args SearchRecipesArgs) (*mcp.CallToolResult, any, error) {
	if args.Input == 0 && args.Output == 0 {
		return errResult("either input or output item ID must be provided")
	}

	s.logger.Debug("Recipe search request", "input", args.Input, "output", args.Output)

	ids, err := s.gw2API.SearchRecipes(ctx, args.Input, args.Output)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to search recipes: %v", err))
	}

	return jsonResult(ids)
}

// handleGetAchievements handles achievement lookup requests
func (s *MCPServer) handleGetAchievements(ctx context.Context, _ *mcp.CallToolRequest, args GetAchievementsArgs) (*mcp.CallToolResult, any, error) {
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Achievements request", "ids", args.IDs)

	achievements, err := s.gw2API.GetAchievements(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get achievements: %v", err))
	}

	return jsonResult(achievements)
}

// handleGetDailyAchievements handles daily achievement requests
func (s *MCPServer) handleGetDailyAchievements(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Daily achievements request")

	dailies, err := s.gw2API.GetDailyAchievements(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get daily achievements: %v", err))
	}

	return jsonResult(dailies)
}

// --- Guild Handlers ---

// handleGetGuild handles public guild info requests
func (s *MCPServer) handleGetGuild(ctx context.Context, _ *mcp.CallToolRequest, args GetGuildArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id parameter is required")
	}

	s.logger.Debug("Guild info request", "id", args.ID)

	guild, err := s.gw2API.GetGuild(ctx, args.ID)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get guild: %v", err))
	}

	return jsonResult(guild)
}

// handleSearchGuild handles guild search requests
func (s *MCPServer) handleSearchGuild(ctx context.Context, _ *mcp.CallToolRequest, args SearchGuildArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name parameter is required")
	}

	s.logger.Debug("Guild search request", "name", args.Name)

	ids, err := s.gw2API.SearchGuild(ctx, args.Name)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to search guild: %v", err))
	}

	return jsonResult(ids)
}

// handleGetGuildDetails handles authenticated guild detail requests
func (s *MCPServer) handleGetGuildDetails(ctx context.Context, _ *mcp.CallToolRequest, args GetGuildDetailsArgs) (*mcp.CallToolResult, any, error) {
	if args.ID == "" {
		return errResult("id parameter is required")
	}
	if args.Type == "" {
		return errResult("type parameter is required")
	}

	s.logger.Debug("Guild detail request", "id", args.ID, "type", args.Type)

	details, err := s.gw2API.GetGuildDetails(ctx, args.ID, args.Type)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get guild details: %v", err))
	}

	return textResult(string(details))
}

// --- Game Metadata Handlers ---

// handleGetColors handles color lookup requests
func (s *MCPServer) handleGetColors(ctx context.Context, _ *mcp.CallToolRequest, args GetColorsArgs) (*mcp.CallToolResult, any, error) {
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Colors request", "ids", args.IDs)

	colors, err := s.gw2API.GetColors(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get colors: %v", err))
	}

	return jsonResult(colors)
}

// handleGetMinis handles mini lookup requests
func (s *MCPServer) handleGetMinis(ctx context.Context, _ *mcp.CallToolRequest, args GetMinisArgs) (*mcp.CallToolResult, any, error) {
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Minis request", "ids", args.IDs)

	minis, err := s.gw2API.GetMinis(ctx, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get minis: %v", err))
	}

	return jsonResult(minis)
}

// handleGetMountsInfo handles mount info requests
func (s *MCPServer) handleGetMountsInfo(ctx context.Context, _ *mcp.CallToolRequest, args GetMountsInfoArgs) (*mcp.CallToolResult, any, error) {
	if args.Type == "" {
		return errResult("type parameter is required")
	}
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Mounts info request", "type", args.Type, "ids", args.IDs)

	data, err := s.gw2API.GetMountsInfo(ctx, args.Type, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get mount info: %v", err))
	}

	return textResult(string(data))
}

// handleGetGameBuild handles game build number requests
func (s *MCPServer) handleGetGameBuild(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Game build request")

	build, err := s.gw2API.GetGameBuild(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get game build: %v", err))
	}

	return jsonResult(build)
}

// handleGetTokenInfo handles API token info requests
func (s *MCPServer) handleGetTokenInfo(ctx context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	s.logger.Debug("Token info request")

	info, err := s.gw2API.GetTokenInfo(ctx)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get token info: %v", err))
	}

	return jsonResult(info)
}

// handleGetDungeonsAndRaids handles dungeon/raid metadata requests
func (s *MCPServer) handleGetDungeonsAndRaids(ctx context.Context, _ *mcp.CallToolRequest, args GetDungeonsAndRaidsArgs) (*mcp.CallToolResult, any, error) {
	if args.Type == "" {
		return errResult("type parameter is required")
	}
	if len(args.IDs) == 0 {
		return errResult("ids parameter is required and must not be empty")
	}

	s.logger.Debug("Dungeons/raids request", "type", args.Type, "ids", args.IDs)

	data, err := s.gw2API.GetDungeonsAndRaids(ctx, args.Type, args.IDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get %s: %v", args.Type, err))
	}

	return textResult(string(data))
}

// --- Composite Tool Handlers ---

// EnrichedRecipe wraps a Recipe with resolved item names
type EnrichedRecipe struct {
	gw2api.Recipe
	OutputItemName  string         `json:"output_item_name"`
	IngredientNames map[int]string `json:"ingredient_names"`
}

// ItemRecipeResult is the response for get_item_recipe_by_name
type ItemRecipeResult struct {
	ItemName string           `json:"item_name"`
	ItemID   int              `json:"item_id,omitempty"`
	Recipes  []EnrichedRecipe `json:"recipes"`
}

// extractItemIDFromWikiResult extracts and parses the item ID from a wiki search result's infobox
func extractItemIDFromWikiResult(result wiki.SearchResult) (int, error) {
	idStr, ok := result.Infobox["id"]
	if !ok {
		return 0, fmt.Errorf("no item ID found in wiki result for %q", result.Title)
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid item ID %q in wiki result for %q: %w", idStr, result.Title, err)
	}
	return id, nil
}

// extractRecipeIDsFromWikiResult extracts and parses recipe IDs from a wiki search result
func extractRecipeIDsFromWikiResult(result wiki.SearchResult) ([]int, error) {
	var ids []int
	for _, recipe := range result.Recipes {
		idStr, ok := recipe["id"]
		if !ok {
			continue
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// handleGetItemByName handles item lookup by name via wiki search
func (s *MCPServer) handleGetItemByName(ctx context.Context, _ *mcp.CallToolRequest, args GetItemByNameArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name parameter is required")
	}

	s.logger.Debug("Item by name request", "name", args.Name)

	// Search wiki for the item
	results, err := s.wiki.Search(ctx, args.Name, 1)
	if err != nil {
		return errResult(fmt.Sprintf("Wiki search failed: %v", err))
	}
	if len(results.Results) == 0 {
		return errResult(fmt.Sprintf("No wiki results found for %q", args.Name))
	}

	// Extract item ID from infobox
	id, err := extractItemIDFromWikiResult(results.Results[0])
	if err != nil {
		return errResult(err.Error())
	}

	// Fetch full item details
	items, err := s.gw2API.GetItems(ctx, []int{id})
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get item details: %v", err))
	}

	item, ok := items[id]
	if !ok {
		return errResult(fmt.Sprintf("Item ID %d not found in API response", id))
	}

	return jsonResult(item)
}

// handleGetItemRecipeByName handles recipe lookup by item name via wiki search
func (s *MCPServer) handleGetItemRecipeByName(ctx context.Context, _ *mcp.CallToolRequest, args GetItemRecipeByNameArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name parameter is required")
	}

	s.logger.Debug("Item recipe by name request", "name", args.Name)

	// Search wiki for the item
	results, err := s.wiki.Search(ctx, args.Name, 1)
	if err != nil {
		return errResult(fmt.Sprintf("Wiki search failed: %v", err))
	}
	if len(results.Results) == 0 {
		return errResult(fmt.Sprintf("No wiki results found for %q", args.Name))
	}

	wikiResult := results.Results[0]

	// Try to get recipe IDs from wiki result
	recipeIDs, _ := extractRecipeIDsFromWikiResult(wikiResult)

	// If no wiki recipes, try fallback via API using item ID
	itemID := 0
	if len(recipeIDs) == 0 {
		id, err := extractItemIDFromWikiResult(wikiResult)
		if err != nil {
			return errResult(fmt.Sprintf("No recipes found in wiki and %v", err))
		}
		itemID = id

		apiRecipeIDs, err := s.gw2API.SearchRecipes(ctx, 0, itemID)
		if err != nil {
			return errResult(fmt.Sprintf("Failed to search recipes by output item: %v", err))
		}
		if len(apiRecipeIDs) == 0 {
			return errResult(fmt.Sprintf("No recipes found for %q (item ID %d)", args.Name, itemID))
		}
		recipeIDs = apiRecipeIDs
	} else if id, err := extractItemIDFromWikiResult(wikiResult); err == nil {
		itemID = id
	}

	// Fetch full recipe details
	recipes, err := s.gw2API.GetRecipes(ctx, recipeIDs)
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get recipe details: %v", err))
	}

	// Collect all item IDs for name resolution
	itemIDSet := make(map[int]bool)
	for _, r := range recipes {
		itemIDSet[r.OutputItemID] = true
		for _, ing := range r.Ingredients {
			itemIDSet[ing.ItemID] = true
		}
	}
	allItemIDs := make([]int, 0, len(itemIDSet))
	for id := range itemIDSet {
		allItemIDs = append(allItemIDs, id)
	}

	// Resolve item names
	itemMap, err := s.gw2API.GetItems(ctx, allItemIDs)
	if err != nil {
		s.logger.Warn("Failed to resolve item names for recipes", "error", err)
		itemMap = make(map[int]gw2api.Item)
	}

	// Build enriched recipes
	enriched := make([]EnrichedRecipe, len(recipes))
	for i, r := range recipes {
		ingNames := make(map[int]string)
		for _, ing := range r.Ingredients {
			if item, ok := itemMap[ing.ItemID]; ok {
				ingNames[ing.ItemID] = item.Name
			}
		}
		outputName := ""
		if item, ok := itemMap[r.OutputItemID]; ok {
			outputName = item.Name
		}
		enriched[i] = EnrichedRecipe{
			Recipe:          r,
			OutputItemName:  outputName,
			IngredientNames: ingNames,
		}
	}

	result := ItemRecipeResult{
		ItemName: wikiResult.Title,
		ItemID:   itemID,
		Recipes:  enriched,
	}

	return jsonResult(result)
}

// handleGetTPPriceByName handles trading post price lookup by item name via wiki search
func (s *MCPServer) handleGetTPPriceByName(ctx context.Context, _ *mcp.CallToolRequest, args GetTPPriceByNameArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return errResult("name parameter is required")
	}

	s.logger.Debug("TP price by name request", "name", args.Name)

	// Search wiki for the item
	results, err := s.wiki.Search(ctx, args.Name, 1)
	if err != nil {
		return errResult(fmt.Sprintf("Wiki search failed: %v", err))
	}
	if len(results.Results) == 0 {
		return errResult(fmt.Sprintf("No wiki results found for %q", args.Name))
	}

	// Extract item ID from infobox
	id, err := extractItemIDFromWikiResult(results.Results[0])
	if err != nil {
		return errResult(err.Error())
	}

	// Fetch trading post prices
	prices, err := s.gw2API.GetPrices(ctx, []int{id})
	if err != nil {
		return errResult(fmt.Sprintf("Failed to get trading post prices: %v", err))
	}
	if len(prices) == 0 {
		return errResult(fmt.Sprintf("No trading post data found for item ID %d", id))
	}

	return jsonResult(prices[0])
}

// handleCurrencyListResource handles the currency list resource
func (s *MCPServer) handleCurrencyListResource(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	s.logger.Debug("Currency list resource request")

	currencies, err := s.gw2API.GetCurrencies(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get currencies: %w", err)
	}

	currenciesJSON, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to format currencies: %w", err)
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "gw2://currencies",
				MIMEType: "application/json",
				Text:     string(currenciesJSON),
			},
		},
	}, nil
}
