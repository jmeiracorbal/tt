package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Header
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 1)

	headerFillStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("62"))

	headerRightStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("255")).
				Padding(0, 1)

	// Bottom menu
	menuFillStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("237"))

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("255")).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("237")).
				Foreground(lipgloss.Color("243")).
				Padding(0, 2)

	menuShortcutsStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("237")).
				Foreground(lipgloss.Color("243")).
				Padding(0, 1)

	// Panels
	listPanelStyle = lipgloss.NewStyle().
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("238"))

	logPanelStyle = lipgloss.NewStyle().
			Padding(0, 1)

	sectionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("62")).
				Padding(0, 1)

	// Resource list items
	selectedItemStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("62")).
				Foreground(lipgloss.Color("255")).
				Bold(true)

	selectedSubStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("57")).
				Foreground(lipgloss.Color("189")).
				Italic(true)

	itemSubStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true)

	// Status colors
	okStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	errorStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	pendingStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	buildingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("33"))
	warnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	dimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
)
