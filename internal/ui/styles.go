package ui

import "github.com/charmbracelet/lipgloss"

// Catppuccin Mocha palette — hex values used directly; lipgloss degrades to
// 256-color or 16-color automatically on terminals that need it.
const (
	clrGreen    = "#a6e3a1" // status.success
	clrRed      = "#f38ba8" // status.error
	clrYellow   = "#f9e2af" // status.warning
	clrBlue     = "#89b4fa" // accent.primary / building
	clrLavender = "#b4befe" // border.focus
	clrOverlay0 = "#6c7086" // border.default
	clrText     = "#cdd6f4" // fg.default
	clrSubtext0 = "#a6adc8" // fg.muted
	clrOverlay1 = "#7f849c" // fg.faint
	clrSelBg    = "#3b3d4f" // selection.bg
	clrMantle   = "#181825" // statusbar.bg
	clrSurface1 = "#45475a" // tab.active.bg
)

var (
	// Status glyphs
	okStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(clrGreen))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(clrRed))
	warnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color(clrYellow))
	buildingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(clrBlue))
	pendingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(clrSubtext0))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(clrOverlay1))

	// Header
	headerTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(clrBlue))

	headerAddrStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(clrSubtext0))

	headerFillStyle = lipgloss.NewStyle()

	// Panels — focused (left, list) and unfocused (right, logs)
	focusedPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(clrLavender))

	unfocusedPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(clrOverlay0))

	// Panel section title (first content row inside the panel)
	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(clrText))

	// Resource list items
	selectedItemStyle = lipgloss.NewStyle().
				Background(lipgloss.Color(clrSelBg)).
				Foreground(lipgloss.Color(clrText)).
				Bold(true)

	itemAgeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(clrSubtext0))

	// Bottom tab bar
	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(clrSurface1)).
			Foreground(lipgloss.Color(clrText)).
			Padding(0, 1)

	inactiveTabStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(clrOverlay1)).
				Background(lipgloss.Color(clrMantle)).
				Padding(0, 1)

	// Keybar (bottom hints)
	keybarFillStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(clrMantle))

	keyStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(clrMantle)).
			Foreground(lipgloss.Color(clrBlue))

	descStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(clrMantle)).
			Foreground(lipgloss.Color(clrOverlay1))

	// All Logs banner lines — colored by resource status
	bannerOkStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(clrGreen)).Bold(true)
	bannerErrorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(clrRed)).Bold(true)
	bannerBuildingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(clrBlue)).Bold(true)
	bannerDimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(clrSubtext0)).Bold(true)
)
