package entities

import (
	"context"
	"testing"
)

func TestListResponse_HasMore(t *testing.T) {
	tests := []struct {
		name     string
		response ListResponse
		want     bool
	}{
		{
			name:     "has more pages",
			response: ListResponse{Next: "token123"},
			want:     true,
		},
		{
			name:     "no more pages",
			response: ListResponse{Next: ""},
			want:     false,
		},
		{
			name:     "empty response",
			response: ListResponse{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.response.HasMore(); got != tt.want {
				t.Errorf("HasMore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListAll(t *testing.T) {
	// Test that ListAll handles pagination correctly
	svc := &Service{
		doer: &paginationMockDoer{page: 0},
	}

	ctx := context.Background()
	opts := SearchOptions{Limit: 10}

	entities, err := svc.ListAll(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entities) != 2 {
		t.Errorf("expected 2 entities, got %d", len(entities))
	}
}

func TestListAllBlueprint(t *testing.T) {
	// Test that ListAllBlueprint handles pagination correctly
	svc := &Service{
		doer: &paginationMockDoer{page: 0},
	}

	ctx := context.Background()
	opts := SearchOptions{Limit: 10}

	entities, err := svc.ListAllBlueprint(ctx, "test-blueprint", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entities) != 2 {
		t.Errorf("expected 2 entities, got %d", len(entities))
	}
}

type paginationMockDoer struct {
	page int
}

func (m *paginationMockDoer) Do(ctx context.Context, method, path string, body any, out any) error {
	if out == nil {
		return nil
	}
	resp, ok := out.(*ListResponse)
	if !ok {
		return nil
	}

	// First page: return 1 entity with Next token
	if m.page == 0 {
		resp.Entities = []Entity{{Identifier: "entity1"}}
		resp.Next = "token123"
		resp.OK = true
		m.page++
		return nil
	}

	// Second page: return 1 entity with no Next token (end of pagination)
	if m.page == 1 {
		resp.Entities = []Entity{{Identifier: "entity2"}}
		resp.Next = ""
		resp.OK = true
		m.page++
		return nil
	}

	// Should not reach here
	return nil
}
