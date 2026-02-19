// Package wiki provides functionality for interacting with the Guild Wars 2 wiki API.
package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/AlyxPink/gw2-mcp/internal/cache"
)

const (
	wikiBaseURL    = "https://wiki.guildwars2.com"
	wikiAPIURL     = wikiBaseURL + "/api.php"
	userAgent      = "github.com/AlyxPink/gw2-mcp"
	requestTimeout = 30 * time.Second
)

// Client handles wiki API requests
type Client struct {
	httpClient *http.Client
	cache      *cache.Manager
	logger     *log.Logger
}

// SearchResult represents a single search result from the wiki
type SearchResult struct {
	Title      string              `json:"title"`
	Snippet    string              `json:"snippet"`
	Timestamp  string              `json:"timestamp"`
	URL        string              `json:"url"`
	Extract    string              `json:"extract,omitempty"`
	Infobox    map[string]string   `json:"infobox,omitempty"`
	InfoboxType string             `json:"infobox_type,omitempty"`
	Recipes    []map[string]string `json:"recipes,omitempty"`
	PageID     int                 `json:"pageid"`
	Size       int                 `json:"size"`
	WordCount  int                 `json:"wordcount"`
}

// SearchResponse represents the complete search response
type SearchResponse struct {
	SearchedAt time.Time      `json:"searched_at"`
	Query      string         `json:"query"`
	Results    []SearchResult `json:"results"`
	Total      int            `json:"total"`
}

// APIResponse represents the MediaWiki API response structure
type APIResponse struct {
	BatchComplete string `json:"batchcomplete"`
	Query         struct {
		Search []struct {
			Title     string `json:"title"`
			Snippet   string `json:"snippet"`
			Timestamp string `json:"timestamp"`
			NS        int    `json:"ns"`
			PageID    int    `json:"pageid"`
			Size      int    `json:"size"`
			WordCount int    `json:"wordcount"`
		} `json:"search"`
		SearchInfo struct {
			TotalHits int `json:"totalhits"`
		} `json:"searchinfo"`
	} `json:"query"`
}

// PageContentResponse represents page content API response
type PageContentResponse struct {
	Query struct {
		Pages map[string]struct {
			Title     string `json:"title"`
			Extract   string `json:"extract"`
			Revisions []struct {
				Slots struct {
					Main struct {
						Content string `json:"*"`
					} `json:"main"`
				} `json:"slots"`
			} `json:"revisions"`
			PageID int `json:"pageid"`
			NS     int `json:"ns"`
		} `json:"pages"`
	} `json:"query"`
	BatchComplete string `json:"batchcomplete"`
}

// pageDetails holds the cached extract and infobox for a wiki page
type pageDetails struct {
	Extract     string              `json:"extract"`
	Infobox     map[string]string   `json:"infobox,omitempty"`
	InfoboxType string              `json:"infobox_type,omitempty"`
	Recipes     []map[string]string `json:"recipes,omitempty"`
}

// infoboxResult holds the parsed infobox type and field values
type infoboxResult struct {
	Type   string
	Fields map[string]string
}

// NewClient creates a new wiki client
func NewClient(cacheManager *cache.Manager, logger *log.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
		cache:  cacheManager,
		logger: logger,
	}
}

// Search performs a search on the Guild Wars 2 wiki
func (c *Client) Search(ctx context.Context, query string, limit int) (*SearchResponse, error) {
	// Normalize query for caching
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	cacheKey := c.cache.GetWikiSearchKey(normalizedQuery)

	// Try cache first
	var searchResponse SearchResponse
	if c.cache.GetJSON(cacheKey, &searchResponse) {
		c.logger.Debug("Wiki search cache hit", "query", query)
		return &searchResponse, nil
	}

	c.logger.Debug("Wiki search cache miss, fetching from API", "query", query)

	// Perform search
	searchResults, err := c.performSearch(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Enhance results with page extracts and infobox data
	for i := range searchResults {
		details, err := c.getPageDetails(ctx, searchResults[i].Title)
		if err != nil {
			c.logger.Warn("Failed to get page details", "title", searchResults[i].Title, "error", err)
		} else {
			searchResults[i].Extract = details.Extract
			searchResults[i].Infobox = details.Infobox
			searchResults[i].InfoboxType = details.InfoboxType
			searchResults[i].Recipes = details.Recipes
		}
		searchResults[i].URL = fmt.Sprintf("%s/wiki/%s", wikiBaseURL, url.QueryEscape(searchResults[i].Title))
	}

	// Create response
	searchResponse = SearchResponse{
		Query:      query,
		Results:    searchResults,
		Total:      len(searchResults),
		SearchedAt: time.Now(),
	}

	// Cache the result
	if err := c.cache.SetJSON(cacheKey, searchResponse, cache.WikiDataTTL); err != nil {
		c.logger.Warn("Failed to cache search results", "error", err)
	}

	return &searchResponse, nil
}

// performSearch makes the actual search API call
func (c *Client) performSearch(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	// Build search URL
	params := url.Values{
		"action":   {"query"},
		"format":   {"json"},
		"list":     {"search"},
		"srsearch": {query},
		"srlimit":  {fmt.Sprintf("%d", limit)},
		"srprop":   {"size|wordcount|timestamp|snippet"},
	}

	searchURL := fmt.Sprintf("%s?%s", wikiAPIURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, http.NoBody)
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
			return nil, fmt.Errorf("wiki API request failed with status %d and failed to read body: %w",
				resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("wiki API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	// Convert to our format
	results := make([]SearchResult, len(apiResponse.Query.Search))
	for i, item := range apiResponse.Query.Search {
		results[i] = SearchResult{
			Title:     item.Title,
			PageID:    item.PageID,
			Size:      item.Size,
			WordCount: item.WordCount,
			Snippet:   c.cleanSnippet(item.Snippet),
			Timestamp: item.Timestamp,
		}
	}

	return results, nil
}

// getPageDetails retrieves the extract and infobox data for a wiki page
func (c *Client) getPageDetails(ctx context.Context, title string) (*pageDetails, error) {
	cacheKey := c.cache.GetWikiPageKey(title)

	// Try cache first
	var cached pageDetails
	if c.cache.GetJSON(cacheKey, &cached) {
		return &cached, nil
	}

	// Build URL requesting both extracts and wikitext revisions
	params := url.Values{
		"action":          {"query"},
		"format":          {"json"},
		"prop":            {"extracts|revisions"},
		"titles":          {title},
		"exintro":         {"true"},
		"explaintext":     {"true"},
		"exsectionformat": {"plain"},
		"exchars":         {"500"},
		"rvprop":          {"content"},
		"rvslots":         {"main"},
	}

	detailsURL := fmt.Sprintf("%s?%s", wikiAPIURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", detailsURL, http.NoBody)
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
		return nil, fmt.Errorf("page details API request failed with status %d", resp.StatusCode)
	}

	var contentResponse PageContentResponse
	if err := json.NewDecoder(resp.Body).Decode(&contentResponse); err != nil {
		return nil, fmt.Errorf("failed to decode page details response: %w", err)
	}

	// Extract the content and parse infobox from wikitext
	details := &pageDetails{}
	for _, page := range contentResponse.Query.Pages {
		details.Extract = page.Extract
		if len(page.Revisions) > 0 {
			wikitext := page.Revisions[0].Slots.Main.Content
			if ib := parseInfobox(wikitext); ib != nil {
				details.Infobox = ib.Fields
				details.InfoboxType = ib.Type
			}
			details.Recipes = parseRecipes(wikitext)
		}
		break // Take the first (and should be only) page
	}

	// Cache the details
	if err := c.cache.SetJSON(cacheKey, details, cache.WikiDataTTL); err != nil {
		c.logger.Warn("Failed to cache page details", "error", err)
	}

	return details, nil
}

// Compiled regexes for cleanWikiMarkup
var (
	rePipedLink  = regexp.MustCompile(`\[\[[^\]|]+\|([^\]]+)\]\]`)
	reSimpleLink = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	reBold       = regexp.MustCompile(`'''(.+?)'''`)
	reItalic     = regexp.MustCompile(`''(.+?)''`)
)

// cleanWikiMarkup strips common wiki markup from a string, leaving plain text.
func cleanWikiMarkup(s string) string {
	s = rePipedLink.ReplaceAllString(s, "$1")
	s = reSimpleLink.ReplaceAllString(s, "$1")
	s = reBold.ReplaceAllString(s, "$1")
	s = reItalic.ReplaceAllString(s, "$1")
	return strings.TrimSpace(s)
}

// parseInfobox extracts key-value pairs from the first infobox template in wikitext.
// It handles nested {{...}} blocks and returns nil if no infobox is found.
func parseInfobox(wikitext string) *infoboxResult {
	// Find the start of an infobox template (case-insensitive search for "infobox")
	lower := strings.ToLower(wikitext)
	idx := strings.Index(lower, "{{")
	var rawTemplateName string
	for idx >= 0 {
		rest := lower[idx+2:]
		trimmed := strings.TrimSpace(rest)
		// Extract just the template name (before first | or newline)
		templateName := trimmed
		if pipeIdx := strings.IndexAny(templateName, "|\n"); pipeIdx >= 0 {
			templateName = templateName[:pipeIdx]
		}
		templateName = strings.TrimSpace(templateName)
		if strings.Contains(templateName, "infobox") {
			rawTemplateName = templateName
			break
		}
		// Look for next {{
		next := strings.Index(lower[idx+2:], "{{")
		if next < 0 {
			return nil
		}
		idx = idx + 2 + next
	}
	if idx < 0 {
		return nil
	}

	// Derive the infobox type from the template name
	// e.g. "inventory infobox" → "inventory", "infobox item" → "item"
	infoboxType := strings.ReplaceAll(rawTemplateName, "infobox", "")
	infoboxType = strings.TrimSpace(infoboxType)

	// Find the matching closing }} while tracking nesting depth
	depth := 1
	pos := idx + 2
	for pos < len(wikitext)-1 && depth > 0 {
		if wikitext[pos] == '{' && wikitext[pos+1] == '{' {
			depth++
			pos += 2
		} else if wikitext[pos] == '}' && wikitext[pos+1] == '}' {
			depth--
			if depth == 0 {
				break
			}
			pos += 2
		} else {
			pos++
		}
	}

	if depth != 0 {
		return nil
	}

	block := wikitext[idx+2 : pos]

	// Extract | key = value pairs, skipping nested templates
	result := make(map[string]string)
	lines := strings.Split(block, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		line = line[1:] // strip leading |
		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eqIdx])
		value := strings.TrimSpace(line[eqIdx+1:])
		if key == "" {
			continue
		}
		// Skip values that are themselves template calls
		if strings.HasPrefix(value, "{{") {
			continue
		}
		result[key] = cleanWikiMarkup(value)
	}

	if len(result) == 0 {
		return nil
	}
	return &infoboxResult{Type: infoboxType, Fields: result}
}

// parseRecipes extracts all {{Recipe}} templates from wikitext.
// Each recipe is returned as a map of key-value pairs with cleaned markup.
func parseRecipes(wikitext string) []map[string]string {
	var recipes []map[string]string
	lower := strings.ToLower(wikitext)
	searchFrom := 0

	for searchFrom < len(wikitext) {
		// Find next {{Recipe (case-insensitive)
		idx := strings.Index(lower[searchFrom:], "{{recipe")
		if idx < 0 {
			break
		}
		idx += searchFrom

		// Verify it's actually "recipe" and not e.g. "recipe list" by checking the template name
		rest := lower[idx+2:]
		trimmed := strings.TrimSpace(rest)
		templateName := trimmed
		if pipeIdx := strings.IndexAny(templateName, "|\n"); pipeIdx >= 0 {
			templateName = templateName[:pipeIdx]
		}
		templateName = strings.TrimSpace(templateName)

		// Skip if it's not exactly "recipe" (e.g. "recipe list", "recipe box")
		if templateName != "recipe" {
			searchFrom = idx + 2
			continue
		}

		// Find matching }}
		depth := 1
		pos := idx + 2
		for pos < len(wikitext)-1 && depth > 0 {
			if wikitext[pos] == '{' && wikitext[pos+1] == '{' {
				depth++
				pos += 2
			} else if wikitext[pos] == '}' && wikitext[pos+1] == '}' {
				depth--
				if depth == 0 {
					break
				}
				pos += 2
			} else {
				pos++
			}
		}

		if depth != 0 {
			break
		}

		block := wikitext[idx+2 : pos]
		searchFrom = pos + 2

		// Extract | key = value pairs
		fields := make(map[string]string)
		lines := strings.Split(block, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "|") {
				continue
			}
			line = line[1:]
			eqIdx := strings.Index(line, "=")
			if eqIdx < 0 {
				continue
			}
			key := strings.TrimSpace(line[:eqIdx])
			value := strings.TrimSpace(line[eqIdx+1:])
			if key == "" {
				continue
			}
			if strings.HasPrefix(value, "{{") {
				continue
			}
			fields[key] = cleanWikiMarkup(value)
		}

		if len(fields) > 0 {
			recipes = append(recipes, fields)
		}
	}

	return recipes
}

// cleanSnippet removes HTML tags and cleans up the snippet text
func (c *Client) cleanSnippet(snippet string) string {
	// Remove HTML tags
	cleaned := strings.ReplaceAll(snippet, "<span class=\"searchmatch\">", "")
	cleaned = strings.ReplaceAll(cleaned, "</span>", "")
	cleaned = strings.ReplaceAll(cleaned, "&quot;", "\"")
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&lt;", "<")
	cleaned = strings.ReplaceAll(cleaned, "&gt;", ">")

	// Clean up whitespace
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.ReplaceAll(cleaned, "\n", " ")
	cleaned = strings.ReplaceAll(cleaned, "\t", " ")

	// Remove multiple spaces
	for strings.Contains(cleaned, "  ") {
		cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	}

	return cleaned
}
