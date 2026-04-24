package drawer

// Side is the edge the drawer slides from.
type Side string

const (
	SideRight Side = "right"
	SideLeft  Side = "left"
)

// Width is a preset panel width.
type Width string

const (
	WidthSM   Width = "sm"   // 320px
	WidthMD   Width = "md"   // 420px (default)
	WidthLG   Width = "lg"   // 560px
	WidthXL   Width = "xl"   // 720px
	WidthFull Width = "full" // 100vw on mobile, capped on desktop
)

// Config holds the drawer configuration. The drawer is opened/closed via Alpine
// state whose name is derived from ID ("{ID}IsOpen"). To open from outside,
// dispatch a window event: `drawer:open` with `{detail: {id: "<ID>"}}`. The
// drawer listens and matches on ID.
//
// For HTMX-driven flows, return an HX-Trigger header: `{"drawer:open": {"id": "<ID>"}}`.
type Config struct {
	// ID uniquely identifies the drawer. Required. Used for the Alpine state
	// var name (`{ID}IsOpen`) and for the aria-labelledby target (`{ID}Title`).
	ID string

	// Title is the drawer heading. Required for accessibility.
	Title string

	// Side the drawer slides in from. Default: SideRight.
	Side Side

	// Width preset. Default: WidthMD.
	Width Width

	// BodyID is the id attribute of the inner content container. Exposed so
	// HTMX targets can swap content directly: hx-target="#{BodyID}".
	// Default: "{ID}-body".
	BodyID string

	// Persistent disables click-backdrop and Esc-to-close. Default: false.
	Persistent bool

	// Class allows extra CSS classes on the panel (not the overlay).
	Class string
}

// stateVar returns the Alpine state-variable name for this drawer's open bit.
func (cfg Config) StateVar() string {
	return cfg.ID + "IsOpen"
}

// GetBodyID returns the resolved body slot id.
func (cfg Config) GetBodyID() string {
	if cfg.BodyID != "" {
		return cfg.BodyID
	}
	return cfg.ID + "-body"
}

// TitleID returns the id used on the drawer's <h2> for aria-labelledby.
func (cfg Config) TitleID() string {
	return cfg.ID + "Title"
}

// OverlayClasses returns classes for the backdrop overlay.
func (cfg Config) OverlayClasses() string {
	return "fixed inset-0 z-40 bg-black/40 dark:bg-black/60"
}

// PanelClasses returns classes for the sliding panel.
func (cfg Config) PanelClasses() string {
	base := "fixed z-50 top-0 bottom-0 flex flex-col bg-surface dark:bg-surface-dark border-outline dark:border-outline-dark shadow-xl"

	switch cfg.Side {
	case SideLeft:
		base += " left-0 border-r"
	default:
		base += " right-0 border-l"
	}

	switch cfg.Width {
	case WidthSM:
		base += " w-full max-w-[320px]"
	case WidthLG:
		base += " w-full max-w-[560px]"
	case WidthXL:
		base += " w-full max-w-[720px]"
	case WidthFull:
		base += " w-full max-w-full md:max-w-[90vw]"
	default:
		base += " w-full max-w-[420px]"
	}

	if cfg.Class != "" {
		base += " " + cfg.Class
	}
	return base
}

// EnterStart returns the Alpine transition enter-start classes.
func (cfg Config) EnterStart() string {
	if cfg.Side == SideLeft {
		return "-translate-x-full"
	}
	return "translate-x-full"
}

// EnterEnd returns the Alpine transition enter-end classes.
func (cfg Config) EnterEnd() string {
	return "translate-x-0"
}
