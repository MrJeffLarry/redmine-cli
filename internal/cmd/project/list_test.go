package project

import (
	"testing"

	"github.com/sahilm/fuzzy"
)

func TestFuzzySearchProjects(t *testing.T) {
	// Test basic fuzzy matching for project names
	names := []string{
		"Website Redesign",
		"Mobile Application Development",
		"Backend API Refactoring",
		"Customer Portal",
		"Internal Tools",
	}

	// Test exact match
	matches := fuzzy.Find("website", names)
	if len(matches) == 0 {
		t.Error("Expected at least one match for 'website'")
	}
	if matches[0].Str != "Website Redesign" {
		t.Errorf("Expected first match to be 'Website Redesign', got '%s'", matches[0].Str)
	}

	// Test fuzzy match with typo
	matches = fuzzy.Find("mobil", names)
	if len(matches) == 0 {
		t.Error("Expected at least one match for 'mobil' (partial match for mobile)")
	}

	// Test partial match
	matches = fuzzy.Find("api", names)
	if len(matches) == 0 {
		t.Error("Expected at least one match for 'api'")
	}

	// Test fuzzy match with initials
	matches = fuzzy.Find("mad", names)
	if len(matches) == 0 {
		t.Error("Expected at least one match for 'mad' (initials for Mobile Application Development)")
	}
}

func TestFuzzySearchProjectsCaseInsensitive(t *testing.T) {
	names := []string{
		"Website Redesign",
		"Mobile App",
	}

	// Test case insensitive matching
	matchesLower := fuzzy.Find("website", names)
	matchesUpper := fuzzy.Find("WEBSITE", names)
	matchesMixed := fuzzy.Find("WebSite", names)

	if len(matchesLower) == 0 || len(matchesUpper) == 0 || len(matchesMixed) == 0 {
		t.Error("Fuzzy search should be case insensitive")
	}
}

func TestFuzzySearchProjectsNoMatch(t *testing.T) {
	names := []string{
		"Website Redesign",
		"Mobile App",
	}

	// Test no match scenario
	matches := fuzzy.Find("xyz999", names)
	if len(matches) != 0 {
		t.Error("Expected no matches for completely unrelated query")
	}
}
