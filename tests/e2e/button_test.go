package e2e

import (
	"strings"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func TestButton_OriginalPenguinUI(t *testing.T) {
	// Skip if Playwright not available
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	// Setup server
	cleanupServer := setupServer(t)
	defer cleanupServer()

	// Setup Playwright
	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	// Create page
	page, err := browser.NewPage()
	require.NoError(t, err)

	// Navigate to original PenguinUI button page
	_, err = page.Goto(baseURL+"/original/buttons/default-button.html", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	// Test: All button variants are present
	t.Run("AllVariantsPresent", func(t *testing.T) {
		buttons := page.Locator("button")
		count, err := buttons.Count()
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 8, "expected at least 8 buttons")

		// Verify each button is visible
		for i := 0; i < count; i++ {
			button := buttons.Nth(i)
			visible, err := button.IsVisible()
			require.NoError(t, err)
			require.True(t, visible, "button %d should be visible", i)
		}

		t.Logf("✓ Found %d buttons", count)
	})

	// Test: Primary button has correct styling
	t.Run("PrimaryButtonStyling", func(t *testing.T) {
		// Get first button (should be Primary)
		button := page.Locator("button").First()

		// Check text content
		text, err := button.TextContent()
		require.NoError(t, err)
		require.Contains(t, strings.ToLower(text), "primary", "first button should be Primary")

		// Get class attribute
		classAttr, err := button.GetAttribute("class")
		require.NoError(t, err)

		// Verify key classes are present
		require.Contains(t, classAttr, "bg-primary", "should have bg-primary class")
		require.Contains(t, classAttr, "text-on-primary", "should have text-on-primary class")
		require.Contains(t, classAttr, "rounded-radius", "should have rounded-radius class")

		t.Logf("✓ Primary button has correct styling")
	})

	// Test: Button is clickable
	t.Run("ButtonIsClickable", func(t *testing.T) {
		button := page.Locator("button").First()

		// Check if enabled
		disabled, err := button.IsDisabled()
		require.NoError(t, err)
		require.False(t, disabled, "button should be enabled")

		// Try clicking
		err = button.Click()
		require.NoError(t, err)

		t.Logf("✓ Button is clickable")
	})
}

func TestButton_GoATTHComponent(t *testing.T) {
	// Skip if Playwright not available
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	// Setup server
	cleanupServer := setupServer(t)
	defer cleanupServer()

	// Setup Playwright
	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	// Create page
	page, err := browser.NewPage()
	require.NoError(t, err)

	// Navigate to GoATTH button demo
	_, err = page.Goto(baseURL+"/gottha/button", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	// Take screenshot for reference
	takeScreenshot(t, page, "gottha-button-demo")

	// Test: Page loads successfully
	t.Run("PageLoads", func(t *testing.T) {
		title, err := page.Title()
		require.NoError(t, err)
		require.Contains(t, title, "Button Component", "page title should contain 'Button Component'")

		t.Logf("✓ Page loaded with title: %s", title)
	})

	// Test: All button variants render
	t.Run("AllVariantsRender", func(t *testing.T) {
		buttons := page.Locator("button")
		count, err := buttons.Count()
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 8, "expected at least 8 buttons")

		t.Logf("✓ Found %d GoATTH buttons", count)
	})

	// Test: HTMX button works
	t.Run("HTMXButtonWorks", func(t *testing.T) {
		// Find HTMX button by its text
		htmxButton := page.Locator("button:has-text('Load Content')")

		// Check hx-get attribute
		hxGet, err := htmxButton.GetAttribute("hx-get")
		require.NoError(t, err)
		require.Equal(t, "/api/hello", hxGet, "should have hx-get='/api/hello'")

		t.Logf("✓ HTMX button configured correctly")
	})

	// Test: Alpine.js button works
	t.Run("AlpineButtonWorks", func(t *testing.T) {
		// Find Alpine button
		alpineButton := page.Locator("button:has-text('Increment Counter')")

		// Check x-on:click attribute
		onClick, err := alpineButton.GetAttribute("x-on:click")
		require.NoError(t, err)
		require.Equal(t, "count++", onClick, "should have x-on:click='count++'")

		t.Logf("✓ Alpine.js button configured correctly")
	})
}
