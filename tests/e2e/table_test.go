package e2e

import (
	"net/http"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTable_DefaultTable tests the default table variant
func TestTable_DefaultTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/table", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Page_Loads_Successfully", func(t *testing.T) {
		title, err := page.Title()
		require.NoError(t, err)
		assert.Contains(t, title, "Table")
		t.Log("✓ Table page loads successfully")
	})

	t.Run("Default_Table_Renders", func(t *testing.T) {
		// Find tables on the page
		tables := page.Locator("table")
		count, err := tables.Count()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 1, "should have at least one table")

		// Check first table has correct structure
		firstTable := tables.First()
		classAttr, err := firstTable.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, classAttr, "w-full")
		assert.Contains(t, classAttr, "text-left")
		assert.Contains(t, classAttr, "text-sm")

		t.Log("✓ Default table renders with correct classes")
	})

	t.Run("Table_Has_Headers", func(t *testing.T) {
		// Check header cells exist
		headers := page.Locator("table").First().Locator("thead th")
		count, err := headers.Count()
		require.NoError(t, err)
		assert.Equal(t, 4, count, "default table should have 4 headers")

		// Verify header content
		firstHeader, err := headers.Nth(0).TextContent()
		require.NoError(t, err)
		assert.Equal(t, "CustomerID", firstHeader)

		t.Log("✓ Table headers render correctly")
	})

	t.Run("Table_Has_Rows", func(t *testing.T) {
		rows := page.Locator("table").First().Locator("tbody tr")
		count, err := rows.Count()
		require.NoError(t, err)
		assert.Equal(t, 3, count, "default table should have 3 rows")

		// Verify first row content
		firstCell, err := rows.Nth(0).Locator("td").First().TextContent()
		require.NoError(t, err)
		assert.Equal(t, "2335", firstCell)

		t.Log("✓ Table rows render correctly")
	})

	t.Run("Table_Container_Has_Border", func(t *testing.T) {
		container := page.Locator("table").First().Locator("xpath=..")
		classAttr, err := container.GetAttribute("class")
		require.NoError(t, err)

		// Walk up to the overflow container
		outerContainer := container.Locator("xpath=..")
		outerClass, err := outerContainer.GetAttribute("class")
		require.NoError(t, err)

		// One of the containers should have border and rounded classes
		combined := classAttr + " " + outerClass
		assert.Contains(t, combined, "border", "container should have border")
		assert.Contains(t, combined, "rounded-radius", "container should have rounded corners")

		t.Log("✓ Table container has correct styling")
	})
}

// TestTable_StripedVariant tests the striped table variant
func TestTable_StripedVariant(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/table", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Striped_Rows_Have_Even_Classes", func(t *testing.T) {
		// Find the striped table section
		stripedHeading := page.Locator("h4:has-text('Striped Table')")
		stripedSection := stripedHeading.Locator("xpath=..")

		// Get rows in the striped table
		rows := stripedSection.Locator("table tbody tr")
		count, err := rows.Count()
		require.NoError(t, err)
		assert.Equal(t, 6, count, "striped table should have 6 rows")

		// Check that rows have even: classes for striping
		firstRow := rows.First()
		classAttr, err := firstRow.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, classAttr, "even:bg-primary/5", "striped rows should have even background")

		t.Log("✓ Striped table rows have correct classes")
	})
}

// TestTable_WithCheckbox tests the checkbox table variant
func TestTable_WithCheckbox(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/table", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Checkbox_Table_Has_Checkboxes", func(t *testing.T) {
		// Find the checkbox table section
		checkboxHeading := page.Locator("h4:has-text('Table with Checkbox')")
		checkboxSection := checkboxHeading.Locator("xpath=..")

		// Check for checkboxes
		checkboxes := checkboxSection.Locator("input[type='checkbox']")
		count, err := checkboxes.Count()
		require.NoError(t, err)
		// Should have 1 header checkbox + 3 row checkboxes = 4
		assert.Equal(t, 4, count, "checkbox table should have 4 checkboxes (1 header + 3 rows)")

		t.Log("✓ Checkbox table has correct number of checkboxes")
	})

	t.Run("Check_All_Selects_All_Rows", func(t *testing.T) {
		// Find the checkbox table section
		checkboxHeading := page.Locator("h4:has-text('Table with Checkbox')")
		checkboxSection := checkboxHeading.Locator("xpath=..")

		// Click the header checkbox (check all)
		headerCheckbox := checkboxSection.Locator("thead input[type='checkbox']")
		err := headerCheckbox.Click()
		require.NoError(t, err)
		page.WaitForTimeout(300)

		// Verify all row checkboxes are checked
		rowCheckboxes := checkboxSection.Locator("tbody input[type='checkbox']")
		rowCount, err := rowCheckboxes.Count()
		require.NoError(t, err)

		for i := 0; i < rowCount; i++ {
			checked, err := rowCheckboxes.Nth(i).IsChecked()
			require.NoError(t, err)
			assert.True(t, checked, "row checkbox %d should be checked", i)
		}

		t.Log("✓ Check all selects all row checkboxes")
	})
}

// TestTable_WithAction tests the action button table variant
func TestTable_WithAction(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/table", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Action_Table_Has_Edit_Buttons", func(t *testing.T) {
		// Find the action table section
		actionHeading := page.Locator("h4:has-text('Table with Action')")
		actionSection := actionHeading.Locator("xpath=..")

		// Check for Edit buttons
		editButtons := actionSection.Locator("button:has-text('Edit')")
		count, err := editButtons.Count()
		require.NoError(t, err)
		assert.Equal(t, 3, count, "action table should have 3 Edit buttons")

		// Verify button styling
		firstButton := editButtons.First()
		classAttr, err := firstButton.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, classAttr, "text-primary", "Edit button should have primary color")

		t.Log("✓ Action table has correctly styled Edit buttons")
	})
}

// TestTable_UsersTable tests the users table with avatars and status badges
func TestTable_UsersTable(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page, err := browser.NewPage()
	require.NoError(t, err)

	_, err = page.Goto(baseURL+"/components/table", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Users_Table_Has_Avatars", func(t *testing.T) {
		// Find the users table section
		usersHeading := page.Locator("h4:has-text('Users Table')")
		usersSection := usersHeading.Locator("xpath=..")

		// Check for avatar images
		avatars := usersSection.Locator("img.rounded-full")
		count, err := avatars.Count()
		require.NoError(t, err)
		assert.Equal(t, 5, count, "users table should have 5 avatar images")

		t.Log("✓ Users table has avatar images")
	})

	t.Run("Users_Table_Has_Status_Badges", func(t *testing.T) {
		usersHeading := page.Locator("h4:has-text('Users Table')")
		usersSection := usersHeading.Locator("xpath=..")

		// Check for Active badges
		activeBadges := usersSection.Locator("span:has-text('Active')")
		activeCount, err := activeBadges.Count()
		require.NoError(t, err)
		assert.Equal(t, 4, activeCount, "should have 4 Active badges")

		// Check for Canceled badge
		canceledBadges := usersSection.Locator("span:has-text('Canceled')")
		canceledCount, err := canceledBadges.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, canceledCount, "should have 1 Canceled badge")

		t.Log("✓ Users table has correct status badges")
	})
}

// TestTable_AllVariantsRender tests that all table variants render without errors
func TestTable_AllVariantsRender(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	t.Run("Table_Page_Returns_200", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/components/table")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		t.Log("✓ Table page returns 200")
	})
}
