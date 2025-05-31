package text

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Unit tests for search, expanded

func TestV2Search(t *testing.T) {
	t.Run("high frequency words", func(t *testing.T) {
		ts, err := NewSearcher("../files/Siddhartha.txt")
		assert.NoError(t, err)
		results := ts.Search("the", 0)
		if len(results) != 2220 { // VSC built-in checker shows 2221, but one of the "the" is one--the, which we could as a different word
			t.Fatalf("expected length %d, got %d", 2220, len(results))
		}
	})

	t.Run("double hypen", func(t *testing.T) {
		ts, err := NewSearcher("../files/Siddhartha.txt")
		assert.NoError(t, err)
		results := ts.Search("one--the", 10)
		assert.ElementsMatch(t,
			[]string{"Chandra Yenco, Isaac Jones Updated editions will replace the previous one--the old editions will be renamed. Creating the works from public"},
			results,
		)
	})

	t.Run("ignore special characters", func(t *testing.T) {
		ts, err := NewSearcher("../files/Siddhartha.txt")
		assert.NoError(t, err)
		results := ts.Search("***", 10)
		assert.ElementsMatch(t,
			[]string{},
			results,
		)
	})

	t.Run("ignore special characters in context", func(t *testing.T) {
		ts, err := NewSearcher("../files/Siddhartha.txt")
		assert.NoError(t, err)
		results := ts.Search("start", 10)
		for _, result := range results {
			if strings.Contains(result, "*") {
				t.Errorf("String should not contain the substring '%s'", "*")
			}
		}
	})
}
