package e2e

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

var (
	baseURL       = "http://localhost:8090"
	screenshotDir = "test-results/screenshots"
)

var serverCmd *exec.Cmd

// setupServer starts the GoATTH server for testing (singleton)
func setupServer(t *testing.T) func() {
	// Check if server is already running
	resp, err := http.Get(baseURL)
	if err == nil {
		resp.Body.Close()
		t.Logf("Server already running on %s", baseURL)
		return func() {} // No cleanup needed, server was already running
	}

	// Build server if not exists
	projectRoot, _ := filepath.Abs("../..")
	serverBin := filepath.Join(projectRoot, "bin", "server")

	if _, err := os.Stat(serverBin); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", "bin/server", "./cmd/server")
		buildCmd.Dir = projectRoot
		output, err := buildCmd.CombinedOutput()
		require.NoError(t, err, "failed to build server: %s", string(output))
	}

	// Start server
	serverCmd = exec.Command(serverBin, "-port", "8090")
	serverCmd.Dir = projectRoot
	serverCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := serverCmd.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}

	// Wait for server to be ready
	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 10*time.Second, 100*time.Millisecond, "server did not start")

	t.Logf("Server started on %s", baseURL)

	// Return cleanup function
	return func() {
		if serverCmd != nil && serverCmd.Process != nil {
			syscall.Kill(-serverCmd.Process.Pid, syscall.SIGTERM)
			serverCmd.Wait()
			serverCmd = nil
		}
	}
}

// setupPlaywright initializes Playwright and returns browser
func setupPlaywright(t *testing.T) (*playwright.Playwright, playwright.Browser, func()) {
	pw, err := playwright.Run()
	require.NoError(t, err, "failed to start playwright")

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	require.NoError(t, err, "failed to launch browser")

	cleanup := func() {
		browser.Close()
		pw.Stop()
	}

	return pw, browser, cleanup
}

// takeScreenshot captures a screenshot for debugging
func takeScreenshot(t *testing.T, page playwright.Page, name string) {
	os.MkdirAll(screenshotDir, 0755)
	path := filepath.Join(screenshotDir, fmt.Sprintf("%s-%d.png", name, time.Now().Unix()))
	_, err := page.Screenshot(playwright.PageScreenshotOptions{
		Path:     playwright.String(path),
		FullPage: playwright.Bool(true),
	})
	if err != nil {
		t.Logf("failed to take screenshot: %v", err)
	} else {
		t.Logf("Screenshot saved: %s", path)
	}
}

// normalizeHTMLSimple normalizes HTML for comparison (simple version)
func normalizeHTMLSimple(html string) string {
	// Remove extra whitespace
	html = strings.ReplaceAll(html, ">\n<", "><")
	html = strings.ReplaceAll(html, ">  <", "><")

	// Normalize spaces
	fields := strings.Fields(html)
	return strings.Join(fields, " ")
}
