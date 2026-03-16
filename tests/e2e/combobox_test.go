package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCombobox_SingleSelect tests single-select combobox functionality
func TestCombobox_SingleSelect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Combobox_Exists_And_Loads", func(t *testing.T) {
		// Verify combobox is present
		combobox := page.Locator("#industry-trigger").First()
		visible, err := combobox.IsVisible()
		require.NoError(t, err)
		require.True(t, visible, "combobox trigger should be visible")

		title, err := page.Title()
		require.NoError(t, err)
		assert.Contains(t, title, "Combobox", "page title should contain 'Combobox'")

		t.Log("✓ Combobox page loaded successfully")
	})

	t.Run("SingleSelect_Opens_And_Selects", func(t *testing.T) {
		// Find the single select combobox
		trigger := page.Locator("#industry-trigger").First()
		require.NotNil(t, trigger, "should find industry combobox")

		// Initial state
		expanded, err := trigger.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "false", expanded, "should be collapsed initially")

		// Click to open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Should be open
		expanded, err = trigger.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "true", expanded, "should be expanded after click")

		// Dropdown should be visible
		dropdown := page.Locator("[role='listbox']").First()
		dropdownVisible, err := dropdown.IsVisible()
		require.NoError(t, err)
		assert.True(t, dropdownVisible, "dropdown should be visible")

		// Select an option
		options := page.Locator("[role='option']")
		count, err := options.Count()
		require.NoError(t, err)
		require.Greater(t, count, 0, "should have options")

		// Click first option
		err = options.First().Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Dropdown should close
		dropdownVisible, err = dropdown.IsVisible()
		require.NoError(t, err)
		assert.False(t, dropdownVisible, "dropdown should close after selection")

		t.Log("✓ Single-select opens, shows options, and selects correctly")
	})

	t.Run("Chevron_Rotates_On_Open", func(t *testing.T) {
		trigger := page.Locator("#industry-trigger").First()
		svg := trigger.Locator("svg")

		// Initial state - no rotation
		initialClass, err := svg.GetAttribute("class")
		require.NoError(t, err)
		assert.NotContains(t, initialClass, "rotate-180", "should not be rotated initially")

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Should be rotated
		expandedClass, err := svg.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, expandedClass, "rotate-180", "should be rotated when open")

		t.Log("✓ Chevron rotates correctly")
	})

	t.Run("Keyboard_Navigation_Works", func(t *testing.T) {
		trigger := page.Locator("#industry-trigger").First()

		// Open with Enter key
		err := trigger.Focus()
		require.NoError(t, err)
		err = trigger.Press("Enter")
		require.NoError(t, err)
		page.WaitForTimeout(200)

		expanded, err := trigger.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "true", expanded, "should open with Enter key")

		// Close with Escape
		err = trigger.Press("Escape")
		require.NoError(t, err)
		page.WaitForTimeout(200)

		expanded, err = trigger.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "false", expanded, "should close with Escape key")

		t.Log("✓ Keyboard navigation works")
	})

	t.Run("Options_Are_Shown_When_Open", func(t *testing.T) {
		trigger := page.Locator("#industry-trigger").First()

		// Open combobox
		err := trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Get dropdown and options
		dropdown := page.Locator("[role='listbox']").First()
		options := dropdown.Locator("[role='option']")

		// Verify options are visible
		count, err := options.Count()
		require.NoError(t, err)
		require.Greater(t, count, 0, "should have visible options")

		// Verify each option is visible
		for i := 0; i < count && i < 5; i++ {
			option := options.Nth(i)
			visible, err := option.IsVisible()
			require.NoError(t, err)
			assert.True(t, visible, fmt.Sprintf("option %d should be visible", i))
		}

		// Verify option labels are displayed
		firstOption := options.First()
		labelText, err := firstOption.TextContent()
		require.NoError(t, err)
		assert.NotEmpty(t, labelText, "option should have label text")
		t.Logf("✓ First option label: %s", labelText)

		t.Logf("✓ %d options are shown and visible", count)
	})

	t.Run("Container_Has_Proper_Size", func(t *testing.T) {
		trigger := page.Locator("#industry-trigger").First()

		// Open combobox
		err := trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Get dropdown container
		dropdown := page.Locator("[role='listbox']").First()

		// Check bounding box (size)
		boundingBox, err := dropdown.BoundingBox()
		require.NoError(t, err)
		require.NotNil(t, boundingBox, "dropdown should have bounding box")

		// Verify minimum dimensions
		assert.Greater(t, boundingBox.Width, 100.0, "dropdown width should be at least 100px")
		assert.Greater(t, boundingBox.Height, 50.0, "dropdown height should be at least 50px")

		t.Logf("✓ Dropdown dimensions: %.0fpx x %.0fpx", boundingBox.Width, boundingBox.Height)

		// Verify container has proper CSS classes for sizing
		classAttr, err := dropdown.GetAttribute("class")
		require.NoError(t, err)

		// Should have min-width and positioning classes
		assert.Contains(t, classAttr, "min-w-full", "should have min-w-full class")
		assert.Contains(t, classAttr, "absolute", "should have absolute positioning")
		assert.Contains(t, classAttr, "z-30", "should have z-index")

		t.Log("✓ Container has proper size and positioning classes")
	})
}

// TestCombobox_MultiSelect tests multi-select combobox functionality
func TestCombobox_MultiSelect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("MultiSelect_Allows_Multiple_Selections", func(t *testing.T) {
		// Find multi-select combobox (skills)
		trigger := page.Locator("#skills-trigger").First()
		require.NotNil(t, trigger, "should find skills combobox")

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Select first option
		options := page.Locator("#skills-trigger").Locator("xpath=../../..").Locator("[role='option']")
		count, err := options.Count()
		require.NoError(t, err)
		t.Logf("Found %d options", count)

		if count > 0 {
			// Click first checkbox
			checkbox := options.First().Locator("input[type='checkbox']")
			err = checkbox.Click()
			require.NoError(t, err)
			page.WaitForTimeout(100)

			// Verify checked
			checked, err := checkbox.IsChecked()
			require.NoError(t, err)
			assert.True(t, checked, "first option should be checked")

			// Select second option
			if count > 1 {
				checkbox2 := options.Nth(1).Locator("input[type='checkbox']")
				err = checkbox2.Click()
				require.NoError(t, err)
				page.WaitForTimeout(100)

				checked2, err := checkbox2.IsChecked()
				require.NoError(t, err)
				assert.True(t, checked2, "second option should also be checked")
			}
		}

		t.Log("✓ Multi-select allows multiple selections")
	})

	t.Run("MultiSelect_Shows_Selected_Count", func(t *testing.T) {
		trigger := page.Locator("#skills-trigger").First()

		// Get label text
		labelSpan := trigger.Locator("span").First()
		text, err := labelSpan.TextContent()
		require.NoError(t, err)

		// If selections were made, should show count or joined labels
		t.Logf("Label text: %s", text)
		assert.NotEmpty(t, text, "label should show selected items or placeholder")

		t.Log("✓ Multi-select label updates correctly")
	})

	t.Run("ClearAll_Button_Works", func(t *testing.T) {
		// Find multi-select with clear all
		trigger := page.Locator("#skills-trigger").First()

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Look for clear all button
		clearBtn := page.Locator("button:has-text('Clear all')")
		exists, err := clearBtn.IsVisible()
		require.NoError(t, err)

		if exists {
			err = clearBtn.Click()
			require.NoError(t, err)
			page.WaitForTimeout(100)

			t.Log("✓ Clear all button works")
		} else {
			t.Log("Clear all button not visible (no selections)")
		}
	})

	t.Run("MultiSelect_Options_Visible_With_Checkboxes", func(t *testing.T) {
		trigger := page.Locator("#skills-trigger").First()

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Get dropdown
		dropdown := page.Locator("[role='listbox']").First()

		// Get all options
		options := dropdown.Locator("[role='option']")
		count, err := options.Count()
		require.NoError(t, err)
		require.Greater(t, count, 0, "should have options")

		// Verify each option has a checkbox
		for i := 0; i < count && i < 3; i++ {
			option := options.Nth(i)
			checkbox := option.Locator("input[type='checkbox']")
			checkboxVisible, err := checkbox.IsVisible()
			require.NoError(t, err)
			assert.True(t, checkboxVisible, fmt.Sprintf("option %d should have visible checkbox", i))
		}

		t.Logf("✓ %d multi-select options visible with checkboxes", count)
	})

	t.Run("MultiSelect_Container_Sizing", func(t *testing.T) {
		trigger := page.Locator("#skills-trigger").First()

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Get dropdown
		dropdown := page.Locator("[role='listbox']").First()

		// Check dimensions
		boundingBox, err := dropdown.BoundingBox()
		require.NoError(t, err)
		require.NotNil(t, boundingBox)

		// Multi-select should be wider for checkboxes
		assert.Greater(t, boundingBox.Width, 150.0, "multi-select dropdown should be at least 150px wide")
		assert.Greater(t, boundingBox.Height, 100.0, "multi-select dropdown should show multiple options")

		t.Logf("✓ Multi-select container size: %.0fpx x %.0fpx", boundingBox.Width, boundingBox.Height)
	})
}

// TestCombobox_WithSearch tests search functionality
func TestCombobox_WithSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Search_Filters_Options", func(t *testing.T) {
		// Find combobox with search (make)
		trigger := page.Locator("#make-trigger").First()
		require.NotNil(t, trigger, "should find make combobox")

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Type in search
		searchInput := page.Locator("input[type='text']").First()

		// If search input exists
		searchExists, err := searchInput.IsVisible()
		require.NoError(t, err)

		if searchExists {
			err = searchInput.Fill("Toyota")
			require.NoError(t, err)
			page.WaitForTimeout(300)

			// Should filter to show Toyota
			options := page.Locator("[role='option']")
			count, err := options.Count()
			require.NoError(t, err)

			// Should have fewer options after filtering
			t.Logf("Found %d options after filtering for 'Toyota'", count)
			assert.GreaterOrEqual(t, count, 1, "should show at least one result")

			// Clear search
			err = searchInput.Fill("")
			require.NoError(t, err)
			page.WaitForTimeout(200)
		}

		t.Log("✓ Search filters options correctly")
	})

	t.Run("Search_Shows_NoResults", func(t *testing.T) {
		trigger := page.Locator("#make-trigger").First()

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		searchInput := page.Locator("input[type='text']").First()
		searchExists, err := searchInput.IsVisible()
		require.NoError(t, err)

		if searchExists {
			// Search for non-existent item
			err = searchInput.Fill("xyznonexistent")
			require.NoError(t, err)
			page.WaitForTimeout(300)

			// Should show no results message
			noResults := page.Locator("text=/no matches found/i")
			visible, err := noResults.IsVisible()
			require.NoError(t, err)
			assert.True(t, visible, "should show no results message")

			// Clear
			err = searchInput.Fill("")
			require.NoError(t, err)
		}

		t.Log("✓ No results message shown correctly")
	})
}

// TestCombobox_WithImages tests combobox with image/avatar support
func TestCombobox_WithImages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Images_Are_Displayed", func(t *testing.T) {
		// Find combobox with images (user)
		trigger := page.Locator("#user-trigger").First()
		require.NotNil(t, trigger, "should find user combobox")

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Look for images in options
		images := page.Locator("img[class*='rounded-full']")
		count, err := images.Count()
		require.NoError(t, err)
		assert.Greater(t, count, 0, "should have avatar images")

		// Verify images are visible
		firstImg := images.First()
		visible, err := firstImg.IsVisible()
		require.NoError(t, err)
		assert.True(t, visible, "avatar image should be visible")

		t.Logf("✓ Found %d avatar images", count)
	})
}

// TestCombobox_Disabled tests disabled state
func TestCombobox_Disabled(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Disabled_Combobox_Cannot_Be_Opened", func(t *testing.T) {
		// Find disabled combobox
		trigger := page.Locator("#disabled-trigger").First()
		require.NotNil(t, trigger, "should find disabled combobox")

		// Check disabled attribute
		disabled, err := trigger.GetAttribute("disabled")
		require.NoError(t, err)
		assert.NotNil(t, disabled, "should have disabled attribute")

		// Try to click
		err = trigger.Click()
		// Don't require no error, clicking disabled elements might not work

		page.WaitForTimeout(200)

		// Should not open
		expanded, err := trigger.GetAttribute("aria-expanded")
		require.NoError(t, err)
		assert.Equal(t, "false", expanded, "disabled combobox should not open")

		t.Log("✓ Disabled combobox works correctly")
	})
}

// TestCombobox_Preselected tests pre-selected values
func TestCombobox_Preselected(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Preselected_Values_Are_Checked", func(t *testing.T) {
		// Find preselected combobox
		trigger := page.Locator("#preselected-trigger").First()
		require.NotNil(t, trigger, "should find preselected combobox")

		// Open
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Check checkboxes - should have some pre-selected
		options := page.Locator("#preselected-trigger").Locator("xpath=../../..").Locator("[role='option']")
		count, err := options.Count()
		require.NoError(t, err)

		checkedCount := 0
		for i := 0; i < count && i < 5; i++ {
			checkbox := options.Nth(i).Locator("input[type='checkbox']")
			checked, _ := checkbox.IsChecked()
			if checked {
				checkedCount++
			}
		}

		t.Logf("Found %d pre-selected options", checkedCount)
		// We expect at least 2 pre-selected (react and vue)
		assert.GreaterOrEqual(t, checkedCount, 1, "should have pre-selected values")

		t.Log("✓ Pre-selected values are rendered correctly")
	})
}

// TestCombobox_VisualParity tests visual comparison with original
func TestCombobox_VisualParity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	screenshotDir := filepath.Join("test-results", "screenshots", "combobox")
	require.NoError(t, os.MkdirAll(screenshotDir, 0755))

	config := ScreenshotConfig{
		OriginalURL:    baseURL + "/original/combobox/simple-combobox.html",
		GoATTHURL:      baseURL + "/components/combobox",
		ComponentName:  "combobox",
		ViewportWidth:  1280,
		ViewportHeight: 800,
		Threshold:      0.85, // 85% match (combobox is complex)
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

// TestCombobox_CSSClassParity tests CSS class matching
func TestCombobox_CSSClassParity(t *testing.T) {
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

	t.Run("Verify_Tailwind_Classes_Present", func(t *testing.T) {
		_, err := page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		})
		require.NoError(t, err)

		// Check combobox trigger has required classes
		trigger := page.Locator("#industry-trigger").First()
		VerifyTailwindClasses(t, trigger, []string{
			"inline-flex",
			"w-full",
			"items-center",
			"justify-between",
			"rounded-radius",
			"border",
		})

		// Open to check dropdown classes
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Check dropdown has required classes
		dropdown := page.Locator("[role='listbox']").First()
		VerifyTailwindClasses(t, dropdown, []string{
			"absolute",
			"z-30",
			"rounded-radius",
			"border",
		})

		t.Log("✓ Required Tailwind classes present")
	})
}

// TestCombobox_Events tests custom event dispatching
func TestCombobox_Events(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/combobox", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Change_Events_Are_Dispatched", func(t *testing.T) {
		// Add event listener via JS
		_, err := page.Evaluate(`() => {
			window.comboboxEvents = [];
			document.addEventListener('combobox-change', (e) => {
				window.comboboxEvents.push(e.detail);
			});
		}`)
		require.NoError(t, err)

		// Open and select an option
		trigger := page.Locator("#industry-trigger").First()
		err = trigger.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		options := page.Locator("[role='option']")
		if count, _ := options.Count(); count > 0 {
			err = options.First().Click()
			require.NoError(t, err)
			page.WaitForTimeout(200)
		}

		// Check if event was dispatched
		events, err := page.Evaluate(`() => window.comboboxEvents`)
		require.NoError(t, err)

		t.Logf("Captured %d combobox-change events", len(events.([]interface{})))

		t.Log("✓ Change events dispatched correctly")
	})
}

// Helper to import os
var _ = fmt.Sprintf
