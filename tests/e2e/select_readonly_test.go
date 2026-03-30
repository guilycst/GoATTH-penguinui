package e2e

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelect_ReadonlyDisabledWithHiddenInput(t *testing.T) {
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

	// The visible select should be disabled
	sel := page.Locator("#os-readonly")
	require.NoError(t, sel.WaitFor())

	disabled, err := sel.IsDisabled()
	require.NoError(t, err)
	assert.True(t, disabled, "readonly select should be rendered as disabled")

	// Should show the selected value "Windows"
	val, err := sel.InputValue()
	require.NoError(t, err)
	assert.Equal(t, "windows", val, "readonly select should show the selected value")

	// There should be a hidden input with the same name and value for form submission
	hidden := page.Locator("form#readonlySelectForm input[type='hidden'][name='os-readonly']")
	count, err := hidden.Count()
	require.NoError(t, err)
	assert.Equal(t, 1, count, "should have a hidden input for form submission")

	hiddenVal, err := hidden.GetAttribute("value")
	require.NoError(t, err)
	assert.Equal(t, "windows", hiddenVal, "hidden input should have the selected value")
}
