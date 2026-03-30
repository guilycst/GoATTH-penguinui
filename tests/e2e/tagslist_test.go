package e2e

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagsList_AddAndRemoveTags(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}
	cleanupServer := setupServer(t)
	defer cleanupServer()
	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page := newPage(t, browser)

	_, err := page.Goto(baseURL+"/components/tags-list", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)

	// Wait for Alpine to initialize
	_, err = page.WaitForFunction("() => typeof Alpine !== 'undefined'", nil)
	require.NoError(t, err)

	container := page.Locator("#tagsDemo")
	require.NoError(t, container.WaitFor())

	// Should have 3 initial tags (prod, critical, gpu)
	inputs := container.Locator("input[type='text']")
	count, err := inputs.Count()
	require.NoError(t, err)
	assert.Equal(t, 3, count, "should have 3 initial tags")

	// Add a tag
	addBtn := container.Locator("[data-add-tag]")
	err = addBtn.Click()
	require.NoError(t, err)

	count, err = inputs.Count()
	require.NoError(t, err)
	assert.Equal(t, 4, count, "should have 4 tags after adding")

	// Remove first tag
	removeBtn := container.Locator("button[aria-label='Remove tag']").First()
	err = removeBtn.Click()
	require.NoError(t, err)

	count, err = inputs.Count()
	require.NoError(t, err)
	assert.Equal(t, 3, count, "should have 3 tags after removing one")
}

func TestTagsList_DisabledState(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}
	cleanupServer := setupServer(t)
	defer cleanupServer()
	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	page := newPage(t, browser)

	_, err := page.Goto(baseURL+"/components/tags-list", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	require.NoError(t, err)

	_, err = page.WaitForFunction("() => typeof Alpine !== 'undefined'", nil)
	require.NoError(t, err)

	container := page.Locator("#disabledTagsDemo")
	require.NoError(t, container.WaitFor())

	// Should have inputs but no add button and no remove buttons
	addBtn := container.Locator("[data-add-tag]")
	addCount, err := addBtn.Count()
	require.NoError(t, err)
	assert.Equal(t, 0, addCount, "disabled tags list should not have add button")

	removeBtn := container.Locator("button[aria-label='Remove tag']")
	removeCount, err := removeBtn.Count()
	require.NoError(t, err)
	assert.Equal(t, 0, removeCount, "disabled tags list should not have remove buttons")
}
