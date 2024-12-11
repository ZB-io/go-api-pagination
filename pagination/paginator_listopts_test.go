package pagination

import (
	"github.com/google/go-github/v65/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestlistOpts(t *testing.T) {
	// Define test cases using table-driven approach
	tests := []struct {
		name     string
		input    *PaginatorOpts
		expected *github.ListOptions
	}{
		{
			name:  "Scenario 1: Default Options When Input is Nil",
			input: nil,
			expected: &github.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		},
		{
			name: "Scenario 2: Default Options When ListOptions is Nil",
			input: &PaginatorOpts{
				ListOptions: nil,
			},
			expected: &github.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		},
		{
			name: "Scenario 3: Default PerPage When Zero Value Provided",
			input: &PaginatorOpts{
				ListOptions: &github.ListOptions{
					PerPage: 0,
					Page:    5,
				},
			},
			expected: &github.ListOptions{
				PerPage: 100,
				Page:    5,
			},
		},
		{
			name: "Scenario 4: Preserve Custom PerPage Value",
			input: &PaginatorOpts{
				ListOptions: &github.ListOptions{
					PerPage: 50,
					Page:    1,
				},
			},
			expected: &github.ListOptions{
				PerPage: 50,
				Page:    1,
			},
		},
		{
			name: "Scenario 5: Preserve Page Number",
			input: &PaginatorOpts{
				ListOptions: &github.ListOptions{
					PerPage: 100,
					Page:    10,
				},
			},
			expected: &github.ListOptions{
				PerPage: 100,
				Page:    10,
			},
		},
		{
			name: "Scenario 6: Handle Maximum PerPage Value",
			input: &PaginatorOpts{
				ListOptions: &github.ListOptions{
					PerPage: 1000,
					Page:    1,
				},
			},
			expected: &github.ListOptions{
				PerPage: 1000,
				Page:    1,
			},
		},
		{
			name: "Scenario 7: Maintain Original ListOptions Reference",
			input: &PaginatorOpts{
				ListOptions: &github.ListOptions{
					PerPage: 200,
					Page:    3,
				},
			},
			expected: &github.ListOptions{
				PerPage: 200,
				Page:    3,
			},
		},
	}

	// Execute test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log("Testing:", tt.name)

			// Act
			result := listOpts(tt.input)

			// Assert
			assert.NotNil(t, result, "Result should not be nil")
			assert.Equal(t, tt.expected.PerPage, result.PerPage, "PerPage values should match")
			assert.Equal(t, tt.expected.Page, result.Page, "Page values should match")

			// Additional check for reference maintenance when input is not nil
			if tt.input != nil && tt.input.ListOptions != nil {
				assert.Same(t, tt.input.ListOptions, result, "Should maintain the same reference when input is valid")
				t.Log("Reference check passed")
			}

			t.Log("Test passed successfully")
		})
	}
}

// TODO: Consider adding more edge cases if needed
// TODO: Consider adding tests for concurrent access if the function is used in concurrent contexts
// TODO: Consider adding performance benchmarks if pagination size impacts performance significantly
