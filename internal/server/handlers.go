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
func (s *MCPServer) handleGetWallet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := request.RequireString("api_key")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid api_key parameter: %v", err)), nil
	}

	s.logger.Debug("Wallet request", "api_key_length", len(apiKey))

	// Get wallet information
	wallet, err := s.gw2API.GetWallet(ctx, apiKey)
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
func (s *MCPServer) handleGetTPDelivery(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, err := request.RequireString("api_key")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid api_key parameter: %v", err)), nil
	}

	s.logger.Debug("TP delivery request", "api_key_length", len(apiKey))

	delivery, err := s.gw2API.GetDelivery(ctx, apiKey)
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
	apiKey, err := request.RequireString("api_key")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid api_key parameter: %v", err)), nil
	}

	txType, err := request.RequireString("type")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Invalid type parameter: %v", err)), nil
	}

	s.logger.Debug("TP transactions request", "api_key_length", len(apiKey), "type", txType)

	transactions, err := s.gw2API.GetTransactions(ctx, apiKey, txType)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to get trading post transactions: %v", err)), nil
	}

	transactionsJSON, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to format transactions: %v", err)), nil
	}

	return mcp.NewToolResultText(string(transactionsJSON)), nil
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
