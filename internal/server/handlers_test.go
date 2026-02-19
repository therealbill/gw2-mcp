package server

import (
	"testing"

	"github.com/AlyxPink/gw2-mcp/internal/wiki"
)

func TestExtractItemIDFromWikiResult(t *testing.T) {
	tests := []struct {
		name    string
		result  wiki.SearchResult
		wantID  int
		wantErr bool
	}{
		{
			name: "valid ID",
			result: wiki.SearchResult{
				Title:   "Mystic Coin",
				Infobox: map[string]string{"id": "19976"},
			},
			wantID:  19976,
			wantErr: false,
		},
		{
			name: "missing infobox",
			result: wiki.SearchResult{
				Title: "Some Page",
			},
			wantErr: true,
		},
		{
			name: "missing id field",
			result: wiki.SearchResult{
				Title:   "Some Page",
				Infobox: map[string]string{"name": "Foo"},
			},
			wantErr: true,
		},
		{
			name: "non-numeric id",
			result: wiki.SearchResult{
				Title:   "Some Page",
				Infobox: map[string]string{"id": "abc"},
			},
			wantErr: true,
		},
		{
			name: "empty id",
			result: wiki.SearchResult{
				Title:   "Some Page",
				Infobox: map[string]string{"id": ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := extractItemIDFromWikiResult(tt.result)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if id != tt.wantID {
				t.Errorf("got ID %d, want %d", id, tt.wantID)
			}
		})
	}
}

func TestExtractRecipeIDsFromWikiResult(t *testing.T) {
	tests := []struct {
		name    string
		result  wiki.SearchResult
		wantIDs []int
	}{
		{
			name: "single recipe",
			result: wiki.SearchResult{
				Title: "18 Slot Silk Bag",
				Recipes: []map[string]string{
					{"id": "7800", "output": "18 Slot Silk Bag"},
				},
			},
			wantIDs: []int{7800},
		},
		{
			name: "multiple recipes",
			result: wiki.SearchResult{
				Title: "Some Item",
				Recipes: []map[string]string{
					{"id": "100"},
					{"id": "200"},
					{"id": "300"},
				},
			},
			wantIDs: []int{100, 200, 300},
		},
		{
			name: "no recipes",
			result: wiki.SearchResult{
				Title: "Some Page",
			},
			wantIDs: nil,
		},
		{
			name: "recipe without id",
			result: wiki.SearchResult{
				Title: "Some Item",
				Recipes: []map[string]string{
					{"output": "Foo"},
				},
			},
			wantIDs: nil,
		},
		{
			name: "mixed valid and invalid ids",
			result: wiki.SearchResult{
				Title: "Some Item",
				Recipes: []map[string]string{
					{"id": "100"},
					{"id": "abc"},
					{"id": "300"},
				},
			},
			wantIDs: []int{100, 300},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, err := extractRecipeIDsFromWikiResult(tt.result)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if len(ids) != len(tt.wantIDs) {
				t.Errorf("got %d IDs, want %d", len(ids), len(tt.wantIDs))
				return
			}
			for i, id := range ids {
				if id != tt.wantIDs[i] {
					t.Errorf("ID[%d] = %d, want %d", i, id, tt.wantIDs[i])
				}
			}
		})
	}
}
