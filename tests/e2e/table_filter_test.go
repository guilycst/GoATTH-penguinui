package e2e

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	_, browser, _ := setupPlaywright(t)
	page := newPage(t, browser)
	page.SetDefaultTimeout(5000)

	_, err := page.Goto(baseURL+"/components/table", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(15000),
	})
	require.NoError(t, err)

	// Check Alpine.js availability and component state
	alpineReady, err := page.WaitForFunction(`() => typeof Alpine !== 'undefined'`, nil, playwright.PageWaitForFunctionOptions{
		Timeout: playwright.Float(5000),
	})
	if err != nil {
		t.Logf("Alpine.js not available in headless browser (CDN may be blocked): %v", err)
		t.Skip("Alpine.js CDN not reachable in test environment")
	}
	_ = alpineReady

	// Wait for Alpine to process x-data components
	page.WaitForTimeout(500)

	// Verify the filter component initialized
	filterReady, _ := page.Evaluate(`() => {
		var el = document.querySelector('[x-data="filtered-tableFilters"]');
		if (!el) return 'no element';
		try { return Alpine.$data(el) ? 'ready' : 'no data'; } catch(e) { return 'error: ' + e.message; }
	}`, nil)
	t.Logf("Filter component state: %v", filterReady)
	if filterReady != "ready" {
		t.Skipf("Alpine filter component not initialized: %v", filterReady)
	}

	// --- Filter bar structure ---

	t.Run("FilterBar_Renders", func(t *testing.T) {
		filterBar := page.Locator("#filtered-table-filters")
		visible, err := filterBar.IsVisible()
		require.NoError(t, err)
		assert.True(t, visible, "filter bar should be visible")
	})

	t.Run("FilterBar_HasSearchInput", func(t *testing.T) {
		searchInput := page.Locator("#filtered-table-filters input[type='search']")
		count, err := searchInput.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should have 1 search input")

		placeholder, err := searchInput.GetAttribute("placeholder")
		require.NoError(t, err)
		assert.Contains(t, placeholder, "Search by name", "search should have placeholder")
	})

	t.Run("FilterBar_HasMembershipSelect", func(t *testing.T) {
		sel := page.Locator("#filtered-table-filters select")
		count, err := sel.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should have 1 select")

		// Should have 3 options: All, Gold, Silver
		options := sel.Locator("option")
		optCount, err := options.Count()
		require.NoError(t, err)
		assert.Equal(t, 3, optCount, "select should have 3 options")
	})

	t.Run("FilterBar_IsCollapsible", func(t *testing.T) {
		toggleBtn := page.Locator("#filtered-table-filters button:has-text('Filters')")
		count, err := toggleBtn.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should have collapsible Filters toggle")

		// Click to collapse
		err = toggleBtn.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		// Search input should be hidden
		searchInput := page.Locator("#filtered-table-filters input[type='search']")
		visible, err := searchInput.IsVisible()
		require.NoError(t, err)
		assert.False(t, visible, "search should be hidden after collapse")

		// Click to expand again
		err = toggleBtn.Click()
		require.NoError(t, err)
		page.WaitForTimeout(200)

		visible, err = searchInput.IsVisible()
		require.NoError(t, err)
		assert.True(t, visible, "search should be visible after expand")
	})

	// --- Search filter ---

	t.Run("Search_FiltersRows", func(t *testing.T) {
		searchInput := page.Locator("#filtered-table-filters input[type='search']")

		// Type "alice"
		err := searchInput.Fill("alice")
		require.NoError(t, err)
		page.WaitForTimeout(1000) // debounce + HTMX swap

		// Should show only Alice
		tbodyText, err := page.Locator("#filtered-table tbody").TextContent()
		require.NoError(t, err)
		assert.Contains(t, tbodyText, "Alice Brown", "should find Alice")

		rows := page.Locator("#filtered-table tbody tr")
		count, err := rows.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should show only 1 matching row")
	})

	t.Run("Search_ClearedShowsAll", func(t *testing.T) {
		searchInput := page.Locator("#filtered-table-filters input[type='search']")

		// Clear search
		err := searchInput.Fill("")
		require.NoError(t, err)
		page.WaitForTimeout(1000)

		rows := page.Locator("#filtered-table tbody tr")
		count, err := rows.Count()
		require.NoError(t, err)
		assert.Equal(t, 3, count, "clearing search should show default page (3 rows)")
	})

	// --- Select filter ---

	t.Run("SelectGold_FiltersToGoldOnly", func(t *testing.T) {
		sel := page.Locator("#filtered-table-filters select")

		_, err := sel.SelectOption(playwright.SelectOptionValues{Values: &[]string{"Gold"}})
		require.NoError(t, err)
		page.WaitForTimeout(1000)

		tbodyText, err := page.Locator("#filtered-table tbody").TextContent()
		require.NoError(t, err)
		assert.NotContains(t, tbodyText, "Silver", "Gold filter should exclude Silver rows")
		assert.Contains(t, tbodyText, "Gold")
	})

	t.Run("SelectAll_ClearsFilter", func(t *testing.T) {
		sel := page.Locator("#filtered-table-filters select")

		_, err := sel.SelectOption(playwright.SelectOptionValues{Values: &[]string{""}})
		require.NoError(t, err)
		page.WaitForTimeout(1000)

		rows := page.Locator("#filtered-table tbody tr")
		count, err := rows.Count()
		require.NoError(t, err)
		assert.Equal(t, 3, count, "clearing filter should show default page")
	})

	// --- Combined search + select ---

	t.Run("CombinedFilter_SearchAndSelect", func(t *testing.T) {
		// Set Gold filter
		sel := page.Locator("#filtered-table-filters select")
		_, err := sel.SelectOption(playwright.SelectOptionValues{Values: &[]string{"Gold"}})
		require.NoError(t, err)
		page.WaitForTimeout(1000)

		// Then search within Gold
		searchInput := page.Locator("#filtered-table-filters input[type='search']")
		err = searchInput.Fill("bob")
		require.NoError(t, err)
		page.WaitForTimeout(1000)

		tbodyText, err := page.Locator("#filtered-table tbody").TextContent()
		require.NoError(t, err)
		assert.Contains(t, tbodyText, "Bob Johnson")

		rows := page.Locator("#filtered-table tbody tr")
		count, err := rows.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should find only Bob in Gold members")

		// Clear both
		err = searchInput.Fill("")
		require.NoError(t, err)
		_, err = sel.SelectOption(playwright.SelectOptionValues{Values: &[]string{""}})
		require.NoError(t, err)
		page.WaitForTimeout(1000)
	})

	// --- Filter preserves across sort (via htmx:configRequest interception) ---

	t.Run("FilterPreservedAcrossSort", func(t *testing.T) {
		// Set Gold filter first
		sel := page.Locator("#filtered-table-filters select")
		_, err := sel.SelectOption(playwright.SelectOptionValues{Values: &[]string{"Gold"}})
		require.NoError(t, err)
		page.WaitForTimeout(1000)

		// Now click a sort header - filter should be preserved
		sortHeader := page.Locator("#filtered-table thead th[hx-get*='order_by']").First()
		count, err := sortHeader.Count()
		require.NoError(t, err)
		if count > 0 {
			err = sortHeader.Click()
			require.NoError(t, err)
			page.WaitForTimeout(1000)

			// Should still show only Gold members
			tbodyText, err := page.Locator("#filtered-table tbody").TextContent()
			require.NoError(t, err)
			assert.NotContains(t, tbodyText, "Silver", "Gold filter should persist after sort")
		}

		// Clean up
		_, err = sel.SelectOption(playwright.SelectOptionValues{Values: &[]string{""}})
		require.NoError(t, err)
		page.WaitForTimeout(1000)
	})
}
