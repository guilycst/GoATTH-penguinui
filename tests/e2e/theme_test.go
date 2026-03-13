package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTheme_Colors_VerifyComputedValues tests that buttons have the correct computed CSS colors
func TestTheme_Colors_VerifyComputedValues(t *testing.T) {
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

	// Navigate to button demo
	_, err = page.Goto(baseURL+"/components/button", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("PrimaryButton_LightMode", func(t *testing.T) {
		button := page.Locator("button:has-text('Primary')").First()

		// Get computed background color
		bgColor, err := button.Evaluate("el => window.getComputedStyle(el).backgroundColor", nil)
		require.NoError(t, err)

		// Primary should be black: rgb(0, 0, 0)
		assert.Equal(t, "rgb(0, 0, 0)", bgColor, "Primary button background should be black")

		// Get computed text color
		textColor, err := button.Evaluate("el => window.getComputedStyle(el).color", nil)
		require.NoError(t, err)

		// Text should be white: rgb(255, 255, 255)
		assert.Equal(t, "rgb(255, 255, 255)", textColor, "Primary button text should be white")

		t.Log("✓ Primary button colors correct in light mode")
	})

	t.Run("SuccessButton_LightMode", func(t *testing.T) {
		button := page.Locator("button:has-text('Success')").First()

		// Get computed background color
		bgColor, err := button.Evaluate("el => window.getComputedStyle(el).backgroundColor", nil)
		require.NoError(t, err)

		// Success should be green-300: rgb(134, 239, 172)
		assert.Equal(t, "rgb(134, 239, 172)", bgColor, "Success button background should be green-300")

		// Get computed text color
		textColor, err := button.Evaluate("el => window.getComputedStyle(el).color", nil)
		require.NoError(t, err)

		// Text should be slate-900: rgb(15, 23, 42)
		assert.Equal(t, "rgb(15, 23, 42)", textColor, "Success button text should be slate-900")

		t.Log("✓ Success button colors correct in light mode")
	})

	t.Run("WarningButton_LightMode", func(t *testing.T) {
		button := page.Locator("button:has-text('Warning')").First()

		// Get computed background color
		bgColor, err := button.Evaluate("el => window.getComputedStyle(el).backgroundColor", nil)
		require.NoError(t, err)

		// Warning should be amber-300: rgb(252, 211, 77)
		assert.Equal(t, "rgb(252, 211, 77)", bgColor, "Warning button background should be amber-300")

		// Get computed text color
		textColor, err := button.Evaluate("el => window.getComputedStyle(el).color", nil)
		require.NoError(t, err)

		// Text should be amber-900: rgb(120, 53, 15)
		assert.Equal(t, "rgb(120, 53, 15)", textColor, "Warning button text should be amber-900")

		t.Log("✓ Warning button colors correct in light mode")
	})
}

// TestTheme_Classes_Presence verifies Tailwind utility classes are applied
func TestTheme_Classes_Presence(t *testing.T) {
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

	// Navigate to button demo
	_, err = page.Goto(baseURL+"/components/button", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("GoTTHA_Buttons_Have_Correct_Classes", func(t *testing.T) {
		// Get all GoTTHA buttons
		buttons := page.Locator(".bg-primary, .bg-secondary, .bg-info, .bg-danger, .bg-warning, .bg-success")
		count, err := buttons.Count()
		require.NoError(t, err)
		require.GreaterOrEqual(t, count, 8, "Should have at least 8 buttons")

		// Check that buttons have the expected Tailwind classes
		for i := 0; i < count && i < 8; i++ {
			button := buttons.Nth(i)

			// Check for rounded-2xl class (border-radius: 1rem)
			hasRounded, err := button.Evaluate("el => el.classList.contains('rounded-2xl')", nil)
			require.NoError(t, err)
			assert.True(t, hasRounded.(bool), fmt.Sprintf("Button %d should have rounded-2xl class", i))

			// Check for font-medium class
			hasFontMedium, err := button.Evaluate("el => el.classList.contains('font-medium')", nil)
			require.NoError(t, err)
			assert.True(t, hasFontMedium.(bool), fmt.Sprintf("Button %d should have font-medium class", i))
		}

		t.Logf("✓ All %d buttons have correct Tailwind classes", count)
	})
}

// TestTheme_DarkMode_Toggle verifies dark mode toggle works
func TestTheme_DarkMode_Toggle(t *testing.T) {
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

	// Navigate to button demo
	_, err = page.Goto(baseURL+"/components/button", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("DarkMode_Toggle_Adds_Class", func(t *testing.T) {
		// Initially should not have dark class
		hasDarkClass, err := page.Evaluate("() => document.documentElement.classList.contains('dark')", nil)
		require.NoError(t, err)
		assert.False(t, hasDarkClass.(bool), "Should not have dark class initially")

		// Click dark mode toggle button (moon/sun icon)
		toggleBtn := page.Locator("header button").Nth(1)
		err = toggleBtn.Click()
		require.NoError(t, err)

		// Wait a moment for Alpine.js to update
		page.WaitForTimeout(100)

		// Now should have dark class
		hasDarkClass, err = page.Evaluate("() => document.documentElement.classList.contains('dark')", nil)
		require.NoError(t, err)
		assert.True(t, hasDarkClass.(bool), "Should have dark class after toggle")

		t.Log("✓ Dark mode toggle works correctly")
	})

	t.Run("DarkMode_Persists_In_LocalStorage", func(t *testing.T) {
		// Check localStorage
		darkModeValue, err := page.Evaluate("() => localStorage.getItem('darkMode')", nil)
		require.NoError(t, err)

		// Should be 'true' after toggle
		assert.Equal(t, "true", darkModeValue, "darkMode should be persisted in localStorage")

		t.Log("✓ Dark mode persists in localStorage")
	})
}

// TestTheme_Visual_Comparison takes screenshots for visual regression
func TestTheme_Visual_Comparison(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	// Setup server
	cleanupServer := setupServer(t)
	defer cleanupServer()

	// Setup Playwright
	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	// Create page with specific viewport for consistent screenshots
	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  1280,
			Height: 800,
		},
	})
	require.NoError(t, err)

	// Navigate to button demo
	_, err = page.Goto(baseURL+"/components/button", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Screenshot_Original_Section", func(t *testing.T) {
		// Take screenshot of Original section
		originalSection := page.Locator("text=Original").First()
		screenshotPath := fmt.Sprintf("test-results/screenshots/original-section-%d.png", time.Now().Unix())

		_, err := originalSection.Screenshot(playwright.LocatorScreenshotOptions{
			Path: playwright.String(screenshotPath),
			Type: playwright.ScreenshotTypePng,
		})
		require.NoError(t, err)

		t.Logf("✓ Screenshot saved: %s", screenshotPath)
	})

	t.Run("Screenshot_GoTTHA_Section", func(t *testing.T) {
		// Take screenshot of GoTTHA section
		gotthaSection := page.Locator("text=GoTTHA").First()
		screenshotPath := fmt.Sprintf("test-results/screenshots/gottha-section-%d.png", time.Now().Unix())

		_, err := gotthaSection.Screenshot(playwright.LocatorScreenshotOptions{
			Path: playwright.String(screenshotPath),
			Type: playwright.ScreenshotTypePng,
		})
		require.NoError(t, err)

		t.Logf("✓ Screenshot saved: %s", screenshotPath)
	})
}

// TestTheme_Visual_Parity_99_99 tests for 99.99% visual parity between Original and GoTTHA
func TestTheme_Visual_Parity_99_99(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping E2E test in short mode")
	}

	// Setup server
	cleanupServer := setupServer(t)
	defer cleanupServer()

	// Setup Playwright
	_, browser, cleanupPW := setupPlaywright(t)
	defer cleanupPW()

	// Create page with specific viewport for consistent screenshots
	page, err := browser.NewPage(playwright.BrowserNewPageOptions{
		Viewport: &playwright.Size{
			Width:  1280,
			Height: 800,
		},
	})
	require.NoError(t, err)

	// Navigate to button demo
	_, err = page.Goto(baseURL+"/components/button", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	require.NoError(t, err)

	t.Run("Visual_Parity_99.99_Percent", func(t *testing.T) {
		// Take screenshot of button preview areas only (not the full sections)
		// Get the first button grid from Original section
		originalGrid := page.Locator("#buttonDefault").Locator("..").Locator(".grid").First()
		gotthaGrid := page.Locator("#gottha-button-preview").Locator(".grid").First()

		// Screenshot both grids
		originalScreenshot, err := originalGrid.Screenshot(playwright.LocatorScreenshotOptions{
			Type: playwright.ScreenshotTypePng,
		})
		require.NoError(t, err)

		gotthaScreenshot, err := gotthaGrid.Screenshot(playwright.LocatorScreenshotOptions{
			Type: playwright.ScreenshotTypePng,
		})
		require.NoError(t, err)

		// Compare screenshots using Playwright's built-in comparison
		similarity, err := page.Evaluate(`([original, gottha]) => {
			return new Promise((resolve) => {
				const img1 = new Image();
				const img2 = new Image();
				
				img1.onload = () => {
					img2.onload = () => {
						const canvas1 = document.createElement('canvas');
						const canvas2 = document.createElement('canvas');
						canvas1.width = img1.width;
						canvas1.height = img1.height;
						canvas2.width = img2.width;
						canvas2.height = img2.height;
						
						const ctx1 = canvas1.getContext('2d');
						const ctx2 = canvas2.getContext('2d');
						ctx1.drawImage(img1, 0, 0);
						ctx2.drawImage(img2, 0, 0);
						
						const data1 = ctx1.getImageData(0, 0, canvas1.width, canvas1.height).data;
						const data2 = ctx2.getImageData(0, 0, canvas2.width, canvas2.height).data;
						
						let matchingPixels = 0;
						let totalPixels = data1.length / 4;
						
						for (let i = 0; i < data1.length; i += 4) {
							const r1 = data1[i], g1 = data1[i+1], b1 = data1[i+2], a1 = data1[i+3];
							const r2 = data2[i], g2 = data2[i+1], b2 = data2[i+2], a2 = data2[i+3];
							
							// Calculate color distance
							const distance = Math.sqrt(
								Math.pow(r1 - r2, 2) +
								Math.pow(g1 - g2, 2) +
								Math.pow(b1 - b2, 2) +
								Math.pow(a1 - a2, 2)
							);
							
							// Very strict tolerance for 99.99% parity (only allow minimal anti-aliasing differences)
							if (distance < 5) {
								matchingPixels++;
							}
						}
						
						resolve(matchingPixels / totalPixels);
					};
					img2.src = gottha;
				};
				img1.src = original;
			});
		}`, []interface{}{originalScreenshot, gotthaScreenshot})

		require.NoError(t, err)
		similarityPct := similarity.(float64) * 100

		// Assert 99.99% or higher similarity
		assert.GreaterOrEqual(t, similarityPct, 99.99,
			"Visual parity should be at least 99.99%%, but got %.2f%%", similarityPct)

		t.Logf("✓ Visual parity: %.2f%% (threshold: 99.99%%)", similarityPct)
	})
}
