package issue

import (
	"testing"

	"github.com/sahilm/fuzzy"
)

func TestFuzzySearch(t *testing.T) {
	// Test basic fuzzy matching
	subjects := []string{
		"Fix bug in login module",
		"Add new feature for user registration",
		"Update documentation for API",
		"Refactor authentication logic",
		"Bug fix for password reset",
	}

	// Test exact match
	matches := fuzzy.Find("login", subjects)
	if len(matches) == 0 {
		t.Error("Expected at least one match for 'login'")
	}
	if matches[0].Str != "Fix bug in login module" {
		t.Errorf("Expected first match to be 'Fix bug in login module', got '%s'", matches[0].Str)
	}

	// Test fuzzy match with typo
	matches = fuzzy.Find("lgn", subjects)
	if len(matches) == 0 {
		t.Error("Expected at least one match for 'lgn' (typo for login)")
	}

	// Test partial match
	matches = fuzzy.Find("bug", subjects)
	if len(matches) < 2 {
		t.Error("Expected at least two matches for 'bug'")
	}

	// Test fuzzy match with multiple words
	matches = fuzzy.Find("fxbg", subjects)
	if len(matches) == 0 {
		t.Error("Expected at least one match for 'fxbg' (fuzzy for 'fix bug')")
	}
}

func TestFuzzySearchCaseInsensitive(t *testing.T) {
	subjects := []string{
		"Fix Bug in Login Module",
		"Add New Feature",
	}

	// Test case insensitive matching
	matchesLower := fuzzy.Find("fix", subjects)
	matchesUpper := fuzzy.Find("FIX", subjects)
	matchesMixed := fuzzy.Find("FiX", subjects)

	if len(matchesLower) == 0 || len(matchesUpper) == 0 || len(matchesMixed) == 0 {
		t.Error("Fuzzy search should be case insensitive")
	}
}

func TestFuzzySearchNoMatch(t *testing.T) {
	subjects := []string{
		"Fix bug in login module",
		"Add new feature",
	}

	// Test no match scenario
	matches := fuzzy.Find("xyz123", subjects)
	if len(matches) != 0 {
		t.Error("Expected no matches for completely unrelated query")
	}
}
