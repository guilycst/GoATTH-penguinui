package e2e

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSelect_DefaultRendering tests the default select component renders correctly
func TestSelect_DefaultRendering(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page := newPage(t, browser)

	_, err := page.Goto(baseURL+"/components/select", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)

	t.Run("Default_Select_Has_Label_And_Options", func(t *testing.T) {
		// Check label exists
		label := page.Locator("label[for='os']")
		labelText, err := label.TextContent()
		require.NoError(t, err)
		assert.Contains(t, labelText, "Operating System")

		// Check select element exists with options
		selectEl := page.Locator("select#os")
		count, err := selectEl.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, count, "should have exactly one default select")

		// Check it has options (placeholder + 3 options)
		options := selectEl.Locator("option")
		optCount, err := options.Count()
		require.NoError(t, err)
		assert.Equal(t, 4, optCount, "should have placeholder + 3 options")

		t.Log("Default select renders with label and options")
	})

	t.Run("Select_Has_Chevron_Icon", func(t *testing.T) {
		// The chevron SVG should be present as a sibling of the select
		wrapper := page.Locator("select#os").Locator("xpath=..")
		svg := wrapper.Locator("svg")
		svgCount, err := svg.Count()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, svgCount, 1, "should have chevron icon")

		t.Log("Select has chevron icon")
	})

	t.Run("Select_Has_Appearance_None", func(t *testing.T) {
		selectEl := page.Locator("select#os")
		class, err := selectEl.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, class, "appearance-none", "select should have appearance-none class")

		t.Log("Select has appearance-none for custom styling")
	})
}

// TestSelect_ValidationStates tests error and success states
func TestSelect_ValidationStates(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page := newPage(t, browser)

	_, err := page.Goto(baseURL+"/components/select", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)

	t.Run("Error_State_Has_Danger_Border", func(t *testing.T) {
		selectEl := page.Locator("select#os-error")
		class, err := selectEl.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, class, "border-danger", "error select should have border-danger")

		t.Log("Error state has danger border")
	})

	t.Run("Error_State_Has_Helper_Text", func(t *testing.T) {
		helperText := page.Locator("select#os-error").Locator("xpath=../small")
		text, err := helperText.TextContent()
		require.NoError(t, err)
		assert.Contains(t, text, "Error: Please select an operating system")

		t.Log("Error state shows helper text")
	})

	t.Run("Error_State_Label_Has_Icon", func(t *testing.T) {
		label := page.Locator("label[for='os-error']")
		svg := label.Locator("svg")
		count, err := svg.Count()
		require.NoError(t, err)
		assert.Equal(t, 1, count, "error label should have an icon")

		labelClass, err := label.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, labelClass, "text-danger", "error label should have text-danger")

		t.Log("Error state label has icon and danger color")
	})

	t.Run("Success_State_Has_Success_Border", func(t *testing.T) {
		selectEl := page.Locator("select#os-success")
		class, err := selectEl.GetAttribute("class")
		require.NoError(t, err)
		assert.Contains(t, class, "border-success", "success select should have border-success")

		t.Log("Success state has success border")
	})

	t.Run("Success_State_Has_Preselected_Value", func(t *testing.T) {
		selectEl := page.Locator("select#os-success")
		value, err := selectEl.InputValue()
		require.NoError(t, err)
		assert.Equal(t, "mac", value, "success select should have Mac preselected")

		t.Log("Success state has preselected value")
	})
}

// TestSelect_DisabledState tests the disabled select
func TestSelect_DisabledState(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page := newPage(t, browser)

	_, err := page.Goto(baseURL+"/components/select", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)

	t.Run("Disabled_Select_Has_Disabled_Attribute", func(t *testing.T) {
		selectEl := page.Locator("select#os-disabled")
		isDisabled, err := selectEl.IsDisabled()
		require.NoError(t, err)
		assert.True(t, isDisabled, "disabled select should be disabled")

		t.Log("Disabled select has disabled attribute")
	})
}

// TestSelect_CountrySelect tests the country select with many options
func TestSelect_CountrySelect(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	cleanupServer := setupServer(t)
	defer cleanupServer()

	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page := newPage(t, browser)

	_, err := page.Goto(baseURL+"/components/select", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)

	t.Run("Country_Select_Has_Many_Options", func(t *testing.T) {
		selectEl := page.Locator("select#country")
		options := selectEl.Locator("option")
		count, err := options.Count()
		require.NoError(t, err)
		assert.Greater(t, count, 50, "country select should have many options")

		t.Log("Country select has many options")
	})

	t.Run("Country_Select_Has_Autocomplete", func(t *testing.T) {
		selectEl := page.Locator("select#country")
		autocomplete, err := selectEl.GetAttribute("autocomplete")
		require.NoError(t, err)
		assert.Equal(t, "country", autocomplete, "country select should have autocomplete=country")

		t.Log("Country select has autocomplete attribute")
	})
}
