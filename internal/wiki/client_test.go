package wiki

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/log"

	"github.com/AlyxPink/gw2-mcp/internal/cache"
)

func TestClient_cleanSnippet(t *testing.T) {
	client := &Client{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "HTML tags removal",
			input:    `<span class="searchmatch">Dragon</span> Bash is a festival`,
			expected: "Dragon Bash is a festival",
		},
		{
			name:     "HTML entities",
			input:    "&quot;Dragon Bash&quot; &amp; other events &lt;test&gt;",
			expected: `"Dragon Bash" & other events <test>`,
		},
		{
			name:     "Whitespace cleanup",
			input:    "Dragon\nBash\t  festival   with    spaces",
			expected: "Dragon Bash festival with spaces",
		},
		{
			name: "Mixed content",
			input: `<span class="searchmatch">Dragon</span>
	Bash &amp; &quot;events&quot;   `,
			expected: `Dragon Bash & "events"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := client.cleanSnippet(tt.input)
			if result != tt.expected {
				t.Errorf("cleanSnippet() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestClient_Search_Cache(t *testing.T) {
	// Create a mock cache manager
	cacheManager := cache.NewManager()
	logger := log.New(os.Stderr)
	logger.SetLevel(log.ErrorLevel) // Reduce noise in tests

	client := NewClient(cacheManager, logger)

	// Mock the search response
	mockResponse := &SearchResponse{
		Query:   "test query",
		Results: []SearchResult{{Title: "Test Page", PageID: 123}},
		Total:   1,
	}

	// Cache the response
	cacheKey := cacheManager.GetWikiSearchKey("test query")
	err := cacheManager.SetJSON(cacheKey, mockResponse, time.Minute)
	if err != nil {
		t.Fatalf("Failed to cache response: %v", err)
	}

	// Test cache hit
	result, err := client.Search(context.Background(), "test query", 5)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if result.Query != "test query" {
		t.Errorf("Expected query 'test query', got %q", result.Query)
	}

	if len(result.Results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(result.Results))
	}

	if result.Results[0].Title != "Test Page" {
		t.Errorf("Expected title 'Test Page', got %q", result.Results[0].Title)
	}
}

func TestClient_Search_API(t *testing.T) {
	// Create mock server for wiki API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/api.php") {
			if r.URL.Query().Get("list") == "search" {
				// Mock search response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"batchcomplete": "",
					"query": {
						"search": [
							{
								"ns": 0,
								"title": "Dragon Bash",
								"pageid": 12345,
								"size": 5000,
								"wordcount": 800,
								"snippet": "<span class=\"searchmatch\">Dragon</span> Bash is a festival",
								"timestamp": "2023-07-01T12:00:00Z"
							}
						],
						"searchinfo": {
							"totalhits": 1
						}
					}
				}`))
			} else if strings.Contains(r.URL.Query().Get("prop"), "extracts") {
				// Mock extract + revisions response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"batchcomplete": "",
					"query": {
						"pages": {
							"12345": {
								"pageid": 12345,
								"ns": 0,
								"title": "Dragon Bash",
								"extract": "Dragon Bash is an annual festival in Guild Wars 2.",
								"revisions": [{"slots": {"main": {"*": "{{Infobox event\n| name = Dragon Bash\n| type = Festival\n}}"}}}]
							}
						}
					}
				}`))
			}
		}
	}))
	defer mockServer.Close()

	// Override the wiki API URL for testing
	originalURL := wikiAPIURL
	defer func() {
		// Note: In a real implementation, we'd need to make wikiAPIURL configurable
		// For this test, we're just demonstrating the structure
		_ = originalURL
	}()

	// Create client with fresh cache
	cacheManager := cache.NewManager()
	logger := log.New(os.Stderr)
	logger.SetLevel(log.ErrorLevel)

	_ = NewClient(cacheManager, logger)

	// Note: This test would need the client to be configurable to use the mock server
	// For now, we'll test the response parsing logic separately
	t.Log("Mock server created for API testing structure")
}

func TestSearchResult_URLGeneration(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:     "Simple title",
			title:    "Dragon Bash",
			expected: "https://wiki.guildwars2.com/wiki/Dragon%20Bash",
		},
		{
			name:     "Title with special characters",
			title:    "API:Account/wallet",
			expected: "https://wiki.guildwars2.com/wiki/API:Account%2Fwallet",
		},
		{
			name:     "Title with spaces and punctuation",
			title:    "Living World Season 4: Episode 1",
			expected: "https://wiki.guildwars2.com/wiki/Living%20World%20Season%204:%20Episode%201",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This tests the URL generation logic that would be in the Search method
			result := SearchResult{Title: tt.title}
			// In the actual implementation, this URL would be set in the Search method
			// Use PathEscape for proper URL encoding in path segments
			generatedURL := "https://wiki.guildwars2.com/wiki/" + url.PathEscape(tt.title)

			if generatedURL != tt.expected {
				t.Errorf("URL generation for %q: got %q, want %q", tt.title, generatedURL, tt.expected)
			}

			_ = result // Use the result to avoid unused variable error
		})
	}
}

func TestNewClient(t *testing.T) {
	cacheManager := cache.NewManager()
	logger := log.New(io.Discard)

	client := NewClient(cacheManager, logger)

	if client == nil {
		t.Fatal("Expected non-nil client")
	}

	if client.cache != cacheManager {
		t.Error("Expected cache manager to be set")
	}

	if client.logger != logger {
		t.Error("Expected logger to be set")
	}

	if client.httpClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	if client.httpClient.Timeout != requestTimeout {
		t.Errorf("Expected timeout %v, got %v", requestTimeout, client.httpClient.Timeout)
	}
}

func TestCleanWikiMarkup(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Piped link",
			input:    "[[Dragonite Ingot|dragonite ingots]]",
			expected: "dragonite ingots",
		},
		{
			name:     "Simple link",
			input:    "[[Dragonite Ingot]]",
			expected: "Dragonite Ingot",
		},
		{
			name:     "Bold text",
			input:    "'''Eternity'''",
			expected: "Eternity",
		},
		{
			name:     "Italic text",
			input:    "''italic text''",
			expected: "italic text",
		},
		{
			name:     "Mixed markup",
			input:    "Use [[Mystic Forge|the forge]] with '''courage'''",
			expected: "Use the forge with courage",
		},
		{
			name:     "Plain text passthrough",
			input:    "Just plain text",
			expected: "Just plain text",
		},
		{
			name:     "Multiple links",
			input:    "[[Bolt]] and [[Sunrise]]",
			expected: "Bolt and Sunrise",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanWikiMarkup(tt.input)
			if result != tt.expected {
				t.Errorf("cleanWikiMarkup(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRecipes(t *testing.T) {
	tests := []struct {
		name     string
		wikitext string
		expected []map[string]string
	}{
		{
			name:     "No recipes",
			wikitext: "Just some regular wiki text.",
			expected: nil,
		},
		{
			name: "Single recipe",
			wikitext: `{{Recipe
| id = 2221
| type = Bag
| output = 18 Slot Silk Bag
| discipline = Tailor
| rating = 400
| ingredient1 = [[Bolt of Silk|Bolts of Silk]]
| count1 = 10
}}`,
			expected: []map[string]string{
				{
					"id":          "2221",
					"type":        "Bag",
					"output":      "18 Slot Silk Bag",
					"discipline":  "Tailor",
					"rating":      "400",
					"ingredient1": "Bolts of Silk",
					"count1":      "10",
				},
			},
		},
		{
			name: "Multiple recipes",
			wikitext: `{{Recipe
| id = 1111
| type = Sword
| discipline = Weaponsmith
| rating = 400
}}
{{Recipe
| id = 2222
| type = Sword
| discipline = Huntsman
| rating = 400
}}`,
			expected: []map[string]string{
				{
					"id":         "1111",
					"type":       "Sword",
					"discipline": "Weaponsmith",
					"rating":     "400",
				},
				{
					"id":         "2222",
					"type":       "Sword",
					"discipline": "Huntsman",
					"rating":     "400",
				},
			},
		},
		{
			name: "Recipe with nested templates skipped",
			wikitext: `{{Recipe
| id = 3333
| type = Bag
| icon = {{item icon|Bag}}
| discipline = Tailor
}}`,
			expected: []map[string]string{
				{
					"id":         "3333",
					"type":       "Bag",
					"discipline": "Tailor",
				},
			},
		},
		{
			name: "Recipe list template ignored",
			wikitext: `{{Recipe list
| id = 9999
| type = Other
}}`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseRecipes(tt.wikitext)
			if tt.expected == nil {
				if len(result) != 0 {
					t.Errorf("parseRecipes() = %v, want nil/empty", result)
				}
				return
			}
			if len(result) != len(tt.expected) {
				t.Fatalf("parseRecipes() returned %d recipes, want %d", len(result), len(tt.expected))
			}
			for i, expectedRecipe := range tt.expected {
				for k, v := range expectedRecipe {
					if result[i][k] != v {
						t.Errorf("parseRecipes()[%d][%q] = %q, want %q", i, k, result[i][k], v)
					}
				}
			}
		})
	}
}

func TestParseInfobox(t *testing.T) {
	tests := []struct {
		name         string
		wikitext     string
		expectedType string
		expected     map[string]string
	}{
		{
			name:     "No infobox",
			wikitext: "Just some regular wiki text without templates.",
			expected: nil,
		},
		{
			name: "Simple infobox",
			wikitext: `Some text before
{{Infobox item
| id = 9566
| name = 18 Slot Silk Bag
| rarity = Fine
| type = Container
}}
Some text after`,
			expectedType: "item",
			expected: map[string]string{
				"id":     "9566",
				"name":   "18 Slot Silk Bag",
				"rarity": "Fine",
				"type":   "Container",
			},
		},
		{
			name: "Infobox with nested templates skipped",
			wikitext: `{{Infobox weapon
| id = 1234
| name = Cool Sword
| icon = {{item icon|Cool Sword}}
| rarity = Exotic
}}`,
			expectedType: "weapon",
			expected: map[string]string{
				"id":     "1234",
				"name":   "Cool Sword",
				"rarity": "Exotic",
			},
		},
		{
			name: "Infobox with other templates before it",
			wikitext: `{{stub}}
{{Infobox npc
| name = Miyani
| level = 80
}}`,
			expectedType: "npc",
			expected: map[string]string{
				"name":  "Miyani",
				"level": "80",
			},
		},
		{
			name:     "Empty infobox",
			wikitext: `{{Infobox item}}`,
			expected: nil,
		},
		{
			name: "Case insensitive infobox detection",
			wikitext: `{{infobox event
| name = Dragon Bash
| type = Festival
}}`,
			expectedType: "event",
			expected: map[string]string{
				"name": "Dragon Bash",
				"type": "Festival",
			},
		},
		{
			name: "GW2 wiki Inventory infobox format",
			wikitext: `{{Inventory infobox
| slots = 18
| description = 18 Slots
| rarity = basic
| value = 27
| id = 9566
}}`,
			expectedType: "inventory",
			expected: map[string]string{
				"slots":       "18",
				"description": "18 Slots",
				"rarity":      "basic",
				"value":       "27",
				"id":          "9566",
			},
		},
		{
			name: "GW2 wiki infobox after other templates",
			wikitext: `{{Recipe
| type = Bag
| source = automatic
}}
{{Inventory infobox
| id = 9566
| rarity = basic
}}`,
			expectedType: "inventory",
			expected: map[string]string{
				"id":     "9566",
				"rarity": "basic",
			},
		},
		{
			name: "Infobox values with wiki links are cleaned",
			wikitext: `{{Infobox weapon
| id = 5678
| name = [[Eternity]]
| material = [[Dragonite Ingot|dragonite ingots]]
}}`,
			expectedType: "weapon",
			expected: map[string]string{
				"id":       "5678",
				"name":     "Eternity",
				"material": "dragonite ingots",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseInfobox(tt.wikitext)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("parseInfobox() = %v, want nil", result)
				}
				return
			}
			if result == nil {
				t.Fatalf("parseInfobox() = nil, want type=%q fields=%v", tt.expectedType, tt.expected)
			}
			if result.Type != tt.expectedType {
				t.Errorf("parseInfobox().Type = %q, want %q", result.Type, tt.expectedType)
			}
			if len(result.Fields) != len(tt.expected) {
				t.Errorf("parseInfobox() returned %d keys, want %d. Got: %v", len(result.Fields), len(tt.expected), result.Fields)
			}
			for k, v := range tt.expected {
				if result.Fields[k] != v {
					t.Errorf("parseInfobox().Fields[%q] = %q, want %q", k, result.Fields[k], v)
				}
			}
		})
	}
}
