package v2

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBody_MultiSelect_RendersHiddenInputsPerSelected(t *testing.T) {
	resetRegistry()
	cfg := Config{
		ID:              "status",
		Name:            "status",
		Mode:            ModeMultiple,
		ToggleEndpoint:  "/toggle",
		OptionsEndpoint: "/options",
		ClearEndpoint:   "/clear",
		Source:          Source{Static: []Option{{Value: "creating", Label: "Creating"}, {Value: "running", Label: "Running"}}},
	}
	state := State{
		Options:  cfg.Source.Static,
		Selected: []string{"creating", "running"},
	}

	var buf bytes.Buffer
	require.NoError(t, Body(cfg, state).Render(context.Background(), &buf))
	html := buf.String()

	assert.Equal(t, 2, strings.Count(html, `type="hidden"`), "one hidden input per selected value")
	assert.Contains(t, html, `name="status"`)
	assert.Contains(t, html, `value="creating"`)
	assert.Contains(t, html, `value="running"`)
	assert.Contains(t, html, `data-combobox-body`)
	assert.Contains(t, html, `id="status-body"`)
}

func TestBody_MultiSelect_TriggerLabelReflectsSelection(t *testing.T) {
	cfg := Config{
		ID:          "status",
		Name:        "status",
		Mode:        ModeMultiple,
		Placeholder: "Select statuses",
		ToggleEndpoint: "/t", OptionsEndpoint: "/o", ClearEndpoint: "/c",
		Source: Source{Static: []Option{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}}},
	}

	t.Run("no selection shows placeholder", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, Body(cfg, State{Options: cfg.Source.Static}).Render(context.Background(), &buf))
		assert.Contains(t, buf.String(), "Select statuses")
	})

	t.Run("one selection shows label", func(t *testing.T) {
		var buf bytes.Buffer
		s := State{Options: cfg.Source.Static, Selected: []string{"a"}}
		require.NoError(t, Body(cfg, s).Render(context.Background(), &buf))
		assert.Contains(t, buf.String(), "A")
	})

	t.Run("multi-selection shows count", func(t *testing.T) {
		var buf bytes.Buffer
		s := State{Options: cfg.Source.Static, Selected: []string{"a", "b"}}
		require.NoError(t, Body(cfg, s).Render(context.Background(), &buf))
		assert.Contains(t, buf.String(), "2 selected")
	})
}

func TestBody_SingleSelect_TriggerLabelShowsSelectedOption(t *testing.T) {
	cfg := Config{
		ID:   "provider",
		Name: "provider",
		Mode: ModeSingle,
		ToggleEndpoint: "/t", OptionsEndpoint: "/o", ClearEndpoint: "/c",
		Source: Source{Static: []Option{{Value: "maas", Label: "MAAS"}, {Value: "eks", Label: "EKS"}}},
	}
	state := State{Options: cfg.Source.Static, Selected: []string{"maas"}}

	var buf bytes.Buffer
	require.NoError(t, Body(cfg, state).Render(context.Background(), &buf))
	html := buf.String()

	assert.Contains(t, html, "MAAS")
	assert.Equal(t, 1, strings.Count(html, `type="hidden"`), "single-select has one hidden input")
}

func TestOptionsList_AriaSelectedMatchesState(t *testing.T) {
	cfg := Config{
		ID: "t", Name: "t", Mode: ModeMultiple,
		ToggleEndpoint: "/t", OptionsEndpoint: "/o", ClearEndpoint: "/c",
		Source: Source{Static: []Option{{Value: "a"}, {Value: "b"}, {Value: "c"}}},
	}
	state := State{Options: cfg.Source.Static, Selected: []string{"a", "c"}}

	var buf bytes.Buffer
	require.NoError(t, OptionsList(cfg, state).Render(context.Background(), &buf))
	html := buf.String()

	assert.Contains(t, html, `data-value="a" aria-selected="true"`)
	assert.Contains(t, html, `data-value="b" aria-selected="false"`)
	assert.Contains(t, html, `data-value="c" aria-selected="true"`)
}

func TestOptionsList_LiHasHXAttributesForToggle(t *testing.T) {
	cfg := Config{
		ID: "status", Name: "status", Mode: ModeMultiple,
		ToggleEndpoint: "/ui/combobox/status/toggle",
		OptionsEndpoint: "/ui/combobox/status/options",
		ClearEndpoint: "/ui/combobox/status/clear",
		Source: Source{Static: []Option{{Value: "running", Label: "Running"}}},
		DependsOn: []string{"provider"},
	}
	state := State{Options: cfg.Source.Static}

	var buf bytes.Buffer
	require.NoError(t, OptionsList(cfg, state).Render(context.Background(), &buf))
	html := buf.String()

	assert.Contains(t, html, `hx-post="/ui/combobox/status/toggle"`)
	assert.Contains(t, html, `hx-target="closest [data-combobox-body]"`)
	assert.Contains(t, html, `hx-swap="outerHTML"`)
	assert.Contains(t, html, `hx-vals`)
	assert.Contains(t, html, `&#34;value&#34;:&#34;running&#34;`)
	assert.Contains(t, html, `hx-include="closest [data-combobox] input[type=hidden],[name=&#39;provider&#39;]"`)
}

func TestBody_SearchInput_RenderedWhenEnabled(t *testing.T) {
	cfg := Config{
		ID: "users", Name: "users", Mode: ModeMultiple,
		ToggleEndpoint: "/t", OptionsEndpoint: "/ui/combobox/users/options", ClearEndpoint: "/c",
		Source: Source{Static: []Option{{Value: "a"}}},
		EnableSearch: true,
	}

	t.Run("renders when enabled", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, Body(cfg, State{Options: cfg.Source.Static, Search: "foo"}).Render(context.Background(), &buf))
		html := buf.String()
		assert.Contains(t, html, `data-combobox-search`)
		assert.Contains(t, html, `value="foo"`)
		assert.Contains(t, html, `hx-get="/ui/combobox/users/options"`)
		assert.Contains(t, html, `hx-trigger="input changed delay:200ms"`)
		assert.Contains(t, html, `hx-target="#users-options"`)
		assert.Contains(t, html, `hx-swap="outerHTML"`)
	})

	t.Run("absent when disabled", func(t *testing.T) {
		cfg2 := cfg
		cfg2.EnableSearch = false
		var buf bytes.Buffer
		require.NoError(t, Body(cfg2, State{Options: cfg2.Source.Static}).Render(context.Background(), &buf))
		assert.NotContains(t, buf.String(), `data-combobox-search`)
	})
}

func TestBody_ClearAllButton(t *testing.T) {
	cfg := Config{
		ID: "t", Name: "t", Mode: ModeMultiple,
		ToggleEndpoint: "/t", OptionsEndpoint: "/o", ClearEndpoint: "/ui/combobox/t/clear",
		Source: Source{Static: []Option{{Value: "a"}}},
		EnableClearAll: true,
	}

	t.Run("renders when enabled and selection non-empty", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, Body(cfg, State{Options: cfg.Source.Static, Selected: []string{"a"}}).Render(context.Background(), &buf))
		html := buf.String()
		assert.Contains(t, html, `data-combobox-clear-all`)
		assert.Contains(t, html, `hx-post="/ui/combobox/t/clear"`)
		assert.Contains(t, html, `hx-target="closest [data-combobox-body]"`)
	})

	t.Run("absent when selection empty", func(t *testing.T) {
		var buf bytes.Buffer
		require.NoError(t, Body(cfg, State{Options: cfg.Source.Static}).Render(context.Background(), &buf))
		assert.NotContains(t, buf.String(), `data-combobox-clear-all`)
	})

	t.Run("absent when disabled", func(t *testing.T) {
		cfg2 := cfg
		cfg2.EnableClearAll = false
		var buf bytes.Buffer
		require.NoError(t, Body(cfg2, State{Options: cfg2.Source.Static, Selected: []string{"a"}}).Render(context.Background(), &buf))
		assert.NotContains(t, buf.String(), `data-combobox-clear-all`)
	})
}

func TestCombobox_OuterShell_HasAlpineDataAndContainsBody(t *testing.T) {
	cfg := Config{
		ID: "status", Name: "status", Mode: ModeMultiple,
		ToggleEndpoint: "/t", OptionsEndpoint: "/o", ClearEndpoint: "/c",
		Source: Source{Static: []Option{{Value: "a", Label: "A"}}},
	}
	state := State{Options: cfg.Source.Static}

	var buf bytes.Buffer
	require.NoError(t, Combobox(cfg, state).Render(context.Background(), &buf))
	html := buf.String()

	assert.Contains(t, html, `data-combobox`)
	assert.Contains(t, html, `x-data="{isOpen:false, openedWithKeyboard:false, focusIndex:-1}"`)
	assert.Contains(t, html, `id="status-body"`, "body partial embedded inside shell")
	assert.NotContains(t, html, `allOptions`, "no data state in Alpine")
	assert.NotContains(t, html, `selectedValues`, "no data state in Alpine")
	assert.NotContains(t, html, `filteredOptions`, "no data state in Alpine")
}
