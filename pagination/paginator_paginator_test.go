package pagination

import (
	"context"
	"errors"
	"github.com/google/go-github/v65/github"
	"testing"
	"time"
)

// MockListFunc implements ListFunc interface for testing
type MockListFunc[T any] struct {
	items     [][]T
	errOnPage int
}

func (m *MockListFunc[T]) List(ctx context.Context, opts *github.ListOptions) ([]T, *github.Response, error) {
	if m.errOnPage > 0 && opts.Page == m.errOnPage {
		return nil, nil, errors.New("list function error")
	}

	if opts.Page > len(m.items) {
		return nil, &github.Response{NextPage: 0}, nil
	}

	resp := &github.Response{
		NextPage: opts.Page + 1,
	}
	if opts.Page == len(m.items) {
		resp.NextPage = 0
	}

	return m.items[opts.Page-1], resp, nil
}

// MockProcessFunc implements ProcessFunc interface for testing
type MockProcessFunc[T any] struct {
	processedItems []T
	errOnItem     int
}

func (m *MockProcessFunc[T]) Process(ctx context.Context, item T) error {
	if m.errOnItem > 0 && len(m.processedItems) == m.errOnItem-1 {
		return errors.New("process function error")
	}
	m.processedItems = append(m.processedItems, item)
	return nil
}

// MockRateLimitFunc implements RateLimitFunc interface for testing
type MockRateLimitFunc struct {
	allowCalls int
	callCount  int
}

func (m *MockRateLimitFunc) RateLimit(ctx context.Context, resp *github.Response) (bool, error) {
	m.callCount++
	return m.callCount <= m.allowCalls, nil
}

func TestPaginator(t *testing.T) {
	tests := []struct {
		name           string
		items          [][]int
		allowCalls     int
		errOnPage      int
		errOnItem      int
		expectedItems  int
		expectedError  bool
		cancelContext  bool
		description    string
	}{
		{
			name:          "Scenario 1: Successful Single Page Pagination",
			items:         [][]int{{1, 2, 3}},
			allowCalls:    1,
			expectedItems: 3,
			description:   "Tests basic functionality with single page",
		},
		{
			name:          "Scenario 2: Multi-Page Pagination Success",
			items:         [][]int{{1, 2}, {3, 4}, {5, 6}},
			allowCalls:    3,
			expectedItems: 6,
			description:   "Tests pagination across multiple pages",
		},
		{
			name:          "Scenario 3: Rate Limit Interruption",
			items:         [][]int{{1, 2}, {3, 4}, {5, 6}},
			allowCalls:    1,
			expectedItems: 2,
			description:   "Tests behavior when rate limit is reached",
		},
		{
			name:          "Scenario 4: List Function Error",
			items:         [][]int{{1, 2}, {3, 4}},
			errOnPage:     2,
			allowCalls:    2,
			expectedItems: 2,
			expectedError: true,
			description:   "Tests error handling in list function",
		},
		{
			name:          "Scenario 5: Process Function Error",
			items:         [][]int{{1, 2, 3}},
			errOnItem:     2,
			allowCalls:    1,
			expectedItems: 3,
			expectedError: true,
			description:   "Tests error handling in process function",
		},
		{
			name:           "Scenario 6: Context Cancellation",
			items:          [][]int{{1, 2}, {3, 4}},
			allowCalls:     2,
			cancelContext:  true,
			expectedError:  true,
			description:    "Tests context cancellation handling",
		},
		{
			name:          "Scenario 7: Empty Result Set",
			items:         [][]int{{}},
			allowCalls:    1,
			expectedItems: 0,
			description:   "Tests handling of empty result set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.description)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			listFunc := &MockListFunc[int]{
				items:     tt.items,
				errOnPage: tt.errOnPage,
			}
			processFunc := &MockProcessFunc[int]{
				errOnItem: tt.errOnItem,
			}
			rateLimitFunc := &MockRateLimitFunc{
				allowCalls: tt.allowCalls,
			}

			if tt.cancelContext {
				go func() {
					time.Sleep(100 * time.Millisecond)
					cancel()
				}()
			}

			opts := &PaginatorOpts{
				ListOptions: &github.ListOptions{
					PerPage: 100,
					Page:    1,
				},
			}

			items, err := Paginator(ctx, listFunc, processFunc, rateLimitFunc, opts)

			// Verify error conditions
			if tt.expectedError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify number of processed items
			if !tt.expectedError && len(processFunc.processedItems) != tt.expectedItems {
				t.Errorf("expected %d processed items, got %d", tt.expectedItems, len(processFunc.processedItems))
			}

			// Verify returned items
			if !tt.expectedError && len(items) != tt.expectedItems {
				t.Errorf("expected %d items, got %d", tt.expectedItems, len(items))
			}

			t.Logf("Test completed: %s", tt.name)
		})
	}
}
