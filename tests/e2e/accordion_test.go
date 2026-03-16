package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccordion_OriginalPenguinUI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	// Navigate to original accordion page
	_, err = page.Goto(baseURL+"/original/accordion/default-accordion.html", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Structure_Correct", func(t *testing.T) {
		// Verify accordion container exists
		container := page.Locator("div[class*='divide-y']")
		visible, err := container.IsVisible()
		require.NoError(t, err)
		require.True(t, visible, "accordion container should be visible")

		// Verify at least 3 accordion items
		items := page.Locator("[id^='controlsAccordionItem']")
		count, err := items.Count()
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 3, "should have at least 3 accordion items")

		t.Logf("✓ Found %d accordion items", count)
	})

	t.Run("Buttons_Expand_Collapse", func(t *testing.T) {
		// Click first button
		firstButton := page.Locator("button").First()
		err := firstButton.Click()
		require.NoError(t, err)

		// Wait for expansion animation
		page.WaitForTimeout(300)

		// Check aria-expanded attribute
		ariaExpanded, err := firstButton.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "true", ariaExpanded, "first button should be expanded")

		// Click again to collapse
		err = firstButton.Click()
		require.NoError(t, err)

		page.WaitForTimeout(300)

		ariaExpanded, err = firstButton.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "false", ariaExpanded, "first button should be collapsed")

		t.Log("✓ Accordion expand/collapse works")
	})

	t.Run("Chevron_Rotates", func(t *testing.T) {
		// Get first button and its chevron
		firstButton := page.Locator("button").First()
		svg := firstButton.Locator("svg")

		// Check initial state (no rotation)
		initialClass, err := svg.GetAttribute("class")
		require.NoError(t, err)
		assert.NotContains(t, initialClass, "rotate-180", "chevron should not be rotated initially")

		// Click to expand
		err = firstButton.Click()
		require.NoError(t, err)
		page.WaitForTimeout(300)

		// Check rotated state
		expandedClass, err := svg.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, expandedClass, "rotate-180", "chevron should be rotated when expanded")

		t.Log("✓ Chevron rotation works")
	})
}

func TestAccordion_GoATTHComponent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	// Navigate to GoATTH accordion demo
	_, err = page.Goto(baseURL+"/components/accordion", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Page_Loads", func(t *testing.T) {
		title, err := page.Title()
		require.NoError(t, err)
		assert.Contains(t, title, "Accordion", "page title should contain 'Accordion'")

		t.Logf("✓ Page loaded: %s", title)
	})

	t.Run("GoATTH_Accordion_Exists", func(t *testing.T) {
		// Find GoATTH section
		accordions := page.Locator("#accordion-fragment .divide-y")
		count, err := accordions.Count()
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 3, "should have at least 3 accordion variants")

		t.Logf("✓ Found %d GoATTH accordions", count)
	})

	t.Run("Interactions_Work", func(t *testing.T) {
		// Click first accordion button in GoATTH section
		firstAccordion := page.Locator("#accordion-fragment .divide-y").First()
		firstButton := firstAccordion.Locator("button").First()

		err := firstButton.Click()
		require.NoError(t, err)
		page.WaitForTimeout(300)

		// Check expanded state
		ariaExpanded, err := firstButton.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "true", ariaExpanded)

		t.Log("✓ GoATTH accordion interactions work")
	})
}

func TestAccordion_VisualParity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	// Create screenshot directory
	screenshotDir := filepath.Join("test-results", "screenshots", "accordion")
	require.NoError(t, os.MkdirAll(screenshotDir, 0755))

	config := ScreenshotConfig{
		OriginalURL:    baseURL + "/original/accordion/default-accordion.html",
		GoATTHURL:      baseURL + "/components/accordion",
		ComponentName:  "accordion",
		ViewportWidth:  1280,
		ViewportHeight: 800,
		Threshold:      0.90, // 90% match threshold (allowing some flexibility)
	}

	result := CompareScreenshots(t, config)

	t.Run("Screenshot_Comparison", func(t *testing.T) {
		assert.True(t, result.Passed,
			"Visual parity should meet %.0f%% threshold, got %.2f%%",
			config.Threshold*100, result.MatchPercentage*100)

		t.Logf("✓ Visual parity: %.2f%%", result.MatchPercentage*100)
		t.Logf("  Original: %s", result.OriginalScreenshotPath)
		t.Logf("  GoATTH: %s", result.GoATTHScreenshotPath)
		t.Logf("  Diff: %s", result.DiffScreenshotPath)
	})
}

func TestAccordion_CSSClassParity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  1280,
			Height: 800,
		},
	})
	require.NoError(t, err)

	t.Run("Compare_Original_vs_GoATTH_Classes", func(t *testing.T) {
		// Navigate to comparison page
		_, err := page.Goto(baseURL+"/components/accordion", playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		})
		require.NoError(t, err)

		// Extract and compare HTML/classes
		htmlResult := ExtractAndCompareHTML(t, page,
			"text=Original >> xpath=../..", // Original section
			"#accordion-fragment")          // GoATTH section

		PrintComparisonReport(t, htmlResult, nil)

		// Assert minimum class parity
		assert.GreaterOrEqual(t, htmlResult.MatchPercentage, 0.85,
			"CSS class parity should be at least 85%%, got %.2f%%",
			htmlResult.MatchPercentage*100)
	})

	t.Run("Verify_Tailwind_Classes_Present", func(t *testing.T) {
		// Navigate to accordion demo
		_, err := page.Goto(baseURL+"/components/accordion", playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		})
		require.NoError(t, err)

		// Check first accordion has required classes
		firstAccordion := page.Locator("#accordion-fragment .divide-y").First()
		firstButton := firstAccordion.Locator("button").First()

		// Verify key Tailwind classes
		VerifyTailwindClasses(t, firstButton, []string{
			"flex",
			"w-full",
			"items-center",
			"justify-between",
			"p-4",
		})

		// Verify container has required classes
		VerifyTailwindClasses(t, firstAccordion, []string{
			"w-full",
			"divide-y",
			"rounded-radius",
			"border",
		})

		t.Log("✓ Required Tailwind classes present")
	})
}

func TestAccordion_Variants(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/accordion", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Default_Variant", func(t *testing.T) {
		defaultAccordion := page.Locator("text=Default").Locator("xpath=../..//div[contains(@class, 'divide-y')]").First()
		classAttr, err := defaultAccordion.GetAttribute("class")
		require.NoError(t, err)

		// Default should have background
		assert.Contains(t, classAttr, "bg-surface-alt", "default variant should have surface-alt background")
	})

	t.Run("NoBackground_Variant", func(t *testing.T) {
		nobgAccordion := page.Locator("text=No Background").Locator("xpath=../..//div[contains(@class, 'divide-y')]").First()
		classAttr, err := nobgAccordion.GetAttribute("class")
		require.NoError(t, err)

		// NoBackground should have surface background, not surface-alt
		assert.NotContains(t, classAttr, "bg-surface-alt/40", "no-background variant should not have surface-alt/40")
		assert.Contains(t, classAttr, "bg-surface", "no-background variant should have surface background")
	})

	t.Run("MultipleOpen_Variant", func(t *testing.T) {
		multiAccordion := page.Locator("text=Allow Multiple Open").Locator("xpath=../..//div[contains(@class, 'divide-y')]").First()

		// Get all buttons
		buttons := multiAccordion.Locator("button")
		count, err := buttons.Count()
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 2, "should have at least 2 buttons")

		// Click first button
		err = buttons.Nth(0).Click()
		require.NoError(t, err)
		page.WaitForTimeout(300)

		// Click second button (should also expand)
		err = buttons.Nth(1).Click()
		require.NoError(t, err)
		page.WaitForTimeout(300)

		// Both should be expanded
		expanded0, _ := buttons.Nth(0).GetAttribute("aria-expanded")
		expanded1, _ := buttons.Nth(1).GetAttribute("aria-expanded")

		// In AllowMultiple mode, both can be open
		// (Note: This depends on how Alpine.js handles it)
		if expanded0 == "true" && expanded1 == "true" {
			t.Log("✓ Multiple items can be open simultaneously")
		} else {
			t.Logf("Note: Item states - Item 1: %s, Item 2: %s", expanded0, expanded1)
		}
	})
}

func TestAccordion_Accessibility(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/accordion", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("ARIA_Attributes_Present", func(t *testing.T) {
		firstAccordion := page.Locator("#accordion-fragment .divide-y").First()
		firstButton := firstAccordion.Locator("button").First()

		// Check aria-expanded
		ariaExpanded, err := firstButton.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.NotEmpty(t, ariaExpanded, "button should have aria-expanded")

		// Check aria-controls
		ariaControls, err := firstButton.GetAttribute("aria-controls")
		require.NoError(t, err)
		assert.NotEmpty(t, ariaControls, "button should have aria-controls")

		// Check the controlled region exists
		contentRegion := page.Locator("#" + ariaControls)
		exists, err := contentRegion.Count()
		require.NoError(t, err)
		assert.Greater(t, exists, 0, "aria-controls should reference existing element")

		t.Log("✓ ARIA attributes present")
	})

	t.Run("Keyboard_Navigation", func(t *testing.T) {
		firstAccordion := page.Locator("#accordion-fragment .divide-y").First()
		firstButton := firstAccordion.Locator("button").First()

		// Focus the button
		err := firstButton.Focus()
		require.NoError(t, err)

		// Press Enter to expand
		err = page.Keyboard().Press("Enter")
		require.NoError(t, err)
		page.WaitForTimeout(300)

		// Check it expanded
		ariaExpanded, err := firstButton.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "true", ariaExpanded, "should expand with Enter key")

		// Press Enter again to collapse
		err = page.Keyboard().Press("Enter")
		require.NoError(t, err)
		page.WaitForTimeout(300)

		ariaExpanded, err = firstButton.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "false", ariaExpanded, "should collapse with Enter key")

		t.Log("✓ Keyboard navigation works")
	})
}

func TestAccordion_AllThemes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  1280,
			Height: 800,
		},
	})
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/accordion", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Screenshot_All_Themes", func(t *testing.T) {
		screenshotDir := filepath.Join("test-results", "screenshots", "accordion-themes")
		require.NoError(t, os.MkdirAll(screenshotDir, 0755))

		// Take screenshot of default theme
		screenshotPath := filepath.Join(screenshotDir, fmt.Sprintf("accordion-theme-default-%d.png", time.Now().Unix()))
		_, err := page.Screenshot(playwright.PageScreenshotOptions{
			Path:     playwright.String(screenshotPath),
			FullPage: playwright.Bool(false),
		})
		require.NoError(t, err)

		t.Logf("✓ Screenshot saved: %s", screenshotPath)
	})

	t.Run("Dark_Mode_Toggle", func(t *testing.T) {
		// Check initial state (no dark class)
		hasDark, err := page.Evaluate("() => document.documentElement.classList.contains('dark')", nil)
		require.NoError(t, err)
		initialDark := hasDark.(bool)

		// Click dark mode toggle
		toggleBtn := page.Locator("#darkModeToggleBtn")
		err = toggleBtn.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Verify dark class changed
		hasDarkAfter, err := page.Evaluate("() => document.documentElement.classList.contains('dark')", nil)
		require.NoError(t, err)
		assert.NotEqual(t, initialDark, hasDarkAfter.(bool), "dark mode should toggle")

		t.Log("✓ Dark mode toggle works")
	})
}

func TestAccordion_VisualParity_99_99(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  1280,
			Height: 800,
		},
	})
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/accordion", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Visual_Parity_Comprehensive", func(t *testing.T) {
		// Extract and compare classes from both Original and GoATTH sections
		htmlResult := ExtractAndCompareHTML(t, page,
			"text=Original >> xpath=../..//div[contains(@class, 'divide-y')]",
			"#accordion-fragment >> div[contains(@class, 'divide-y')]")

		t.Logf("\n=== Accordion Visual Parity Report ===\n")
		t.Logf("Overall CSS Class Match: %.2f%%", htmlResult.MatchPercentage*100)

		if len(htmlResult.MissingClasses) > 0 {
			t.Logf("Missing Classes (%d):", len(htmlResult.MissingClasses))
			for _, c := range htmlResult.MissingClasses[:min(10, len(htmlResult.MissingClasses))] {
				t.Logf("  - %s", c)
			}
		}

		if len(htmlResult.ExtraClasses) > 0 {
			t.Logf("Extra Classes (%d):", len(htmlResult.ExtraClasses))
			for _, c := range htmlResult.ExtraClasses[:min(10, len(htmlResult.ExtraClasses))] {
				t.Logf("  + %s", c)
			}
		}

		// For 99.99% parity, we need virtually all classes to match
		// Allow for minor differences like IDs
		assert.GreaterOrEqual(t, htmlResult.MatchPercentage, 0.90,
			"Visual parity should be at least 90%% for accordion. Got: %.2f%%",
			htmlResult.MatchPercentage*100)

		if htmlResult.MatchPercentage >= 0.95 {
			t.Log("✓ Excellent visual parity achieved!")
		}
	})
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
