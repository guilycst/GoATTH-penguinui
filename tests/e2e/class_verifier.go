package e2e

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ElementClasses holds class information for an element
type ElementClasses struct {
	Selector string
	TagName  string
	Classes  []string
	Computed map[string]string
}

// ClassComparisonResult holds the result of comparing classes between two elements
type ClassComparisonResult struct {
	OriginalSelector    string
	GoATTHSelector      string
	MissingClasses      []string
	ExtraClasses        []string
	MatchingClasses     []string
	ComputedDifferences map[string]struct {
		Original string
		GoATTH   string
	}
	MatchPercentage float64
}

// ExtractClassesFromPage extracts all classes from elements matching a selector
func ExtractClassesFromPage(t *testing.T, page playwright.Page, selector string) []ElementClasses {
	t.Helper()

	locators := page.Locator(selector)
	count, err := locators.Count()
	require.NoError(t, err)

	var results []ElementClasses

	for i := 0; i < count; i++ {
		locator := locators.Nth(i)

		// Get tag name
		tagName, err := locator.Evaluate("el => el.tagName.toLowerCase()", nil)
		require.NoError(t, err)

		// Get class attribute
		classAttr, err := locator.GetAttribute("class")
		require.NoError(t, err)

		// Parse classes
		classes := parseClassAttribute(classAttr)

		// Get computed styles
		computed := extractComputedStyles(t, locator)

		results = append(results, ElementClasses{
			Selector: fmt.Sprintf("%s:nth-of-type(%d)", selector, i+1),
			TagName:  tagName.(string),
			Classes:  classes,
			Computed: computed,
		})
	}

	return results
}

// CompareElementClasses compares classes between original and GoATTH implementations
func CompareElementClasses(t *testing.T, originalClasses, goatthClasses []ElementClasses) []ClassComparisonResult {
	t.Helper()

	var results []ClassComparisonResult

	minLen := len(originalClasses)
	if len(goatthClasses) < minLen {
		minLen = len(goatthClasses)
	}

	for i := 0; i < minLen; i++ {
		orig := originalClasses[i]
		goatth := goatthClasses[i]

		result := compareSingleElement(t, orig, goatth)
		results = append(results, result)
	}

	if len(originalClasses) != len(goatthClasses) {
		t.Logf("Warning: Different number of elements - Original: %d, GoATTH: %d",
			len(originalClasses), len(goatthClasses))
	}

	return results
}

// compareSingleElement compares classes for a single element pair
func compareSingleElement(t *testing.T, original, goatth ElementClasses) ClassComparisonResult {
	result := ClassComparisonResult{
		OriginalSelector:    original.Selector,
		GoATTHSelector:      goatth.Selector,
		ComputedDifferences: make(map[string]struct{ Original, GoATTH string }),
	}

	origMap := makeSet(original.Classes)
	goatthMap := makeSet(goatth.Classes)

	// Find missing classes (in original but not in GoATTH)
	for class := range origMap {
		if !goatthMap[class] {
			result.MissingClasses = append(result.MissingClasses, class)
		}
	}

	// Find extra classes (in GoATTH but not in original)
	for class := range goatthMap {
		if !origMap[class] {
			result.ExtraClasses = append(result.ExtraClasses, class)
		}
	}

	// Find matching classes
	for class := range origMap {
		if goatthMap[class] {
			result.MatchingClasses = append(result.MatchingClasses, class)
		}
	}

	// Calculate match percentage
	totalClasses := len(origMap) + len(result.ExtraClasses)
	if totalClasses > 0 {
		result.MatchPercentage = float64(len(result.MatchingClasses)) / float64(totalClasses)
	} else {
		result.MatchPercentage = 1.0
	}

	// Compare computed styles
	for key, origValue := range original.Computed {
		if goatthValue, exists := goatth.Computed[key]; exists {
			if origValue != goatthValue {
				result.ComputedDifferences[key] = struct {
					Original string
					GoATTH   string
				}{
					Original: origValue,
					GoATTH:   goatthValue,
				}
			}
		}
	}

	return result
}

// parseClassAttribute parses a class attribute string into individual classes
func parseClassAttribute(classAttr string) []string {
	if classAttr == "" {
		return nil
	}

	// Split by whitespace and filter empty strings
	parts := strings.Fields(classAttr)

	// Remove duplicates while preserving order
	seen := make(map[string]bool)
	var unique []string
	for _, part := range parts {
		if !seen[part] {
			seen[part] = true
			unique = append(unique, part)
		}
	}

	return unique
}

// extractComputedStyles extracts key computed CSS properties
func extractComputedStyles(t *testing.T, locator playwright.Locator) map[string]string {
	t.Helper()

	properties := []string{
		"backgroundColor",
		"color",
		"borderColor",
		"borderRadius",
		"padding",
		"margin",
		"fontSize",
		"fontWeight",
		"display",
	}

	result := make(map[string]string)

	for _, prop := range properties {
		value, err := locator.Evaluate(
			fmt.Sprintf("el => window.getComputedStyle(el).%s", prop),
			nil,
		)
		if err == nil && value != nil {
			result[prop] = fmt.Sprintf("%v", value)
		}
	}

	return result
}

// makeSet creates a set (map) from a slice
func makeSet(items []string) map[string]bool {
	set := make(map[string]bool)
	for _, item := range items {
		set[item] = true
	}
	return set
}

// AssertClassParity asserts that class parity meets the threshold
func AssertClassParity(t *testing.T, results []ClassComparisonResult, threshold float64) {
	t.Helper()

	allPassed := true

	for i, result := range results {
		t.Run(fmt.Sprintf("Element_%d", i), func(t *testing.T) {
			if result.MatchPercentage < threshold {
				allPassed = false
				t.Errorf("Class match %.2f%% below threshold %.2f%% for %s",
					result.MatchPercentage*100, threshold*100, result.GoATTHSelector)

				if len(result.MissingClasses) > 0 {
					t.Errorf("Missing classes: %v", result.MissingClasses)
				}
				if len(result.ExtraClasses) > 0 {
					t.Errorf("Extra classes: %v", result.ExtraClasses)
				}
			} else {
				t.Logf("✓ Element %d: %.2f%% class match", i+1, result.MatchPercentage*100)
			}
		})
	}

	if allPassed {
		t.Logf("✓ All elements meet class parity threshold of %.2f%%", threshold*100)
	}
}

// VerifyTailwindClasses checks that specific Tailwind classes are present
func VerifyTailwindClasses(t *testing.T, locator playwright.Locator, expectedClasses []string) {
	t.Helper()

	classAttr, err := locator.GetAttribute("class")
	require.NoError(t, err)

	actualClasses := parseClassAttribute(classAttr)
	actualSet := makeSet(actualClasses)

	for _, expectedClass := range expectedClasses {
		assert.True(t, actualSet[expectedClass],
			"Expected class '%s' not found. Actual classes: %v",
			expectedClass, actualClasses)
	}
}

// ExtractAndCompareHTML extracts HTML from both implementations and compares
func ExtractAndCompareHTML(t *testing.T, page playwright.Page, originalSelector, goatthSelector string) HTMLComparisonResult {
	t.Helper()

	result := HTMLComparisonResult{}

	// Get original HTML
	originalLocator := page.Locator(originalSelector)
	originalHTML, err := originalLocator.InnerHTML()
	require.NoError(t, err)

	// Get GoATTH HTML
	goatthLocator := page.Locator(goatthSelector)
	goatthHTML, err := goatthLocator.InnerHTML()
	require.NoError(t, err)

	result.OriginalHTML = normalizeHTML(originalHTML)
	result.GoATTHHTML = normalizeHTML(goatthHTML)

	// Extract class attributes from both
	result.OriginalClasses = extractAllClasses(result.OriginalHTML)
	result.GoATTHClasses = extractAllClasses(result.GoATTHHTML)

	// Compare
	origSet := makeSet(result.OriginalClasses)
	goatthSet := makeSet(result.GoATTHClasses)

	for _, class := range result.OriginalClasses {
		if goatthSet[class] {
			result.MatchingClasses = append(result.MatchingClasses, class)
		} else {
			result.MissingClasses = append(result.MissingClasses, class)
		}
	}

	for _, class := range result.GoATTHClasses {
		if !origSet[class] {
			result.ExtraClasses = append(result.ExtraClasses, class)
		}
	}

	// Calculate match
	allClasses := make(map[string]bool)
	for _, c := range result.OriginalClasses {
		allClasses[c] = true
	}
	for _, c := range result.GoATTHClasses {
		allClasses[c] = true
	}

	if len(allClasses) > 0 {
		result.MatchPercentage = float64(len(result.MatchingClasses)) / float64(len(allClasses))
	} else {
		result.MatchPercentage = 1.0
	}

	return result
}

// HTMLComparisonResult holds HTML comparison data
type HTMLComparisonResult struct {
	OriginalHTML    string
	GoATTHHTML      string
	OriginalClasses []string
	GoATTHClasses   []string
	MatchingClasses []string
	MissingClasses  []string
	ExtraClasses    []string
	MatchPercentage float64
}

// extractAllClasses extracts all class attributes from HTML
func extractAllClasses(html string) []string {
	// Regex to find class="..." attributes
	re := regexp.MustCompile(`class="([^"]*)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	var allClasses []string
	for _, match := range matches {
		if len(match) > 1 {
			classes := parseClassAttribute(match[1])
			allClasses = append(allClasses, classes...)
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, c := range allClasses {
		if !seen[c] {
			seen[c] = true
			unique = append(unique, c)
		}
	}

	sort.Strings(unique)
	return unique
}

// PrintComparisonReport prints a detailed comparison report
func PrintComparisonReport(t *testing.T, htmlResult HTMLComparisonResult, classResults []ClassComparisonResult) {
	t.Helper()

	t.Log("\n=== Visual Parity Report ===\n")

	t.Logf("Overall Class Match: %.2f%%", htmlResult.MatchPercentage*100)

	if len(htmlResult.MissingClasses) > 0 {
		t.Logf("\nMissing Classes (%d):", len(htmlResult.MissingClasses))
		for _, c := range htmlResult.MissingClasses {
			t.Logf("  - %s", c)
		}
	}

	if len(htmlResult.ExtraClasses) > 0 {
		t.Logf("\nExtra Classes (%d):", len(htmlResult.ExtraClasses))
		for _, c := range htmlResult.ExtraClasses {
			t.Logf("  + %s", c)
		}
	}

	t.Logf("\nElement-wise Breakdown (%d elements):", len(classResults))
	for i, r := range classResults {
		status := "✓"
		if r.MatchPercentage < 0.95 {
			status = "✗"
		}
		t.Logf("  %s Element %d: %.1f%% match", status, i+1, r.MatchPercentage*100)
	}
}

// normalizeHTML removes whitespace variations for comparison
func normalizeHTML(html string) string {
	// Remove extra whitespace between tags
	re := regexp.MustCompile(`>\s+<`)
	html = re.ReplaceAllString(html, "><")

	// Normalize spaces
	re = regexp.MustCompile(`\s+`)
	html = re.ReplaceAllString(html, " ")

	return strings.TrimSpace(html)
}
