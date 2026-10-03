package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dankkomcg/tilt-tui/internal/api"
)

const pollInterval = 3 * time.Second

type activeTab int

const (
	tabResources activeTab = iota
	tabAllLogs
)

type Model struct {
	client      *api.Client
	addr        string
	resources   []api.Resource
	allLogs     []api.LogSegment
	cursor      int
	resourcesVP viewport.Model
	allLogsVP   viewport.Model
	width       int
	height      int
	ready       bool
	err         error
	loading     bool
	tab         activeTab
}

type viewFetchedMsg struct {
	view *api.View
	err  error
}

type tickMsg struct{}

func New(addr string) Model {
	if addr == "" {
		addr = api.DefaultAddr
	}
	return Model{
		client:  api.NewClient(addr),
		addr:    addr,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchView(m.client), scheduleTick())
}

func fetchView(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		view, err := client.FetchView()
		return viewFetchedMsg{view: view, err: err}
	}
}

func scheduleTick() tea.Cmd {
	return tea.Tick(pollInterval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		rW, rH := m.resourcesVPSize()
		aW, aH := m.allLogsVPSize()
		if !m.ready {
			m.resourcesVP = viewport.New(rW, rH)
			m.allLogsVP = viewport.New(aW, aH)
			m.ready = true
		} else {
			m.resourcesVP.Width = rW
			m.resourcesVP.Height = rH
			m.allLogsVP.Width = aW
			m.allLogsVP.Height = aH
		}
		m.resourcesVP.SetContent(m.resourceLogsContent())
		m.allLogsVP.SetContent(m.allLogsContent())

	case viewFetchedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.err = nil
			m.resources = msg.view.Resources
			m.allLogs = msg.view.LogList.Segments
			if m.ready {
				m.resourcesVP.SetContent(m.resourceLogsContent())
				m.allLogsVP.SetContent(m.allLogsContent())
			}
		}

	case tickMsg:
		cmds = append(cmds, fetchView(m.client), scheduleTick())

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "1":
			m.tab = tabResources
			return m, tea.Batch(cmds...)
		case "2":
			m.tab = tabAllLogs
			return m, tea.Batch(cmds...)
		case "tab":
			m.tab = (m.tab + 1) % 2
			return m, tea.Batch(cmds...)
		case "up", "k":
			if m.tab == tabResources && m.cursor > 0 {
				m.cursor--
				m.resourcesVP.GotoTop()
				m.resourcesVP.SetContent(m.resourceLogsContent())
			}
			return m, tea.Batch(cmds...)
		case "down", "j":
			if m.tab == tabResources && m.cursor < len(m.resources)-1 {
				m.cursor++
				m.resourcesVP.GotoTop()
				m.resourcesVP.SetContent(m.resourceLogsContent())
			}
			return m, tea.Batch(cmds...)
		case "r":
			m.loading = true
			cmds = append(cmds, fetchView(m.client))
			return m, tea.Batch(cmds...)
		}
	}

	var vpCmd tea.Cmd
	switch m.tab {
	case tabResources:
		m.resourcesVP, vpCmd = m.resourcesVP.Update(msg)
	case tabAllLogs:
		m.allLogsVP, vpCmd = m.allLogsVP.Update(msg)
	}
	if vpCmd != nil {
		cmds = append(cmds, vpCmd)
	}

	return m, tea.Batch(cmds...)
}

// ── Layout ────────────────────────────────────────────────────────────────────

func (m Model) leftPanelWidth() int {
	w := m.width / 3
	if w < 22 {
		w = 22
	}
	if w > 42 {
		w = 42
	}
	return w
}

func (m Model) rightPanelWidth() int {
	return m.width - m.leftPanelWidth()
}

// panelHeight is the rows available for panels: total minus header, 2 dot lines, menu.
func (m Model) panelHeight() int {
	h := m.height - 4
	if h < 4 {
		h = 4
	}
	return h
}

// resourcesVPSize: right panel minus padding(2), panel minus title bar(1).
func (m Model) resourcesVPSize() (int, int) {
	w := m.rightPanelWidth() - 2
	if w < 10 {
		w = 10
	}
	h := m.panelHeight() - 1
	if h < 1 {
		h = 1
	}
	return w, h
}

// allLogsVPSize: full width minus padding(2), panel minus title bar(1).
func (m Model) allLogsVPSize() (int, int) {
	w := m.width - 2
	if w < 10 {
		w = 10
	}
	h := m.panelHeight() - 1
	if h < 1 {
		h = 1
	}
	return w, h
}

// ── Content ───────────────────────────────────────────────────────────────────

func (m Model) resourceLogsContent() string {
	if len(m.resources) == 0 {
		return dimStyle.Render("No resources found")
	}
	return logsForResource(m.resources[m.cursor], m.allLogs)
}

func (m Model) allLogsContent() string {
	if len(m.allLogs) == 0 {
		return dimStyle.Render("No logs available")
	}
	var sb strings.Builder
	for _, seg := range m.allLogs {
		sb.WriteString(renderSegment(seg))
	}
	return sb.String()
}

func logsForResource(r api.Resource, allLogs []api.LogSegment) string {
	spanIDs := make(map[string]bool)
	for _, b := range r.BuildHistory {
		if b.SpanID != "" {
			spanIDs[b.SpanID] = true
		}
	}
	if r.CurrentBuild.SpanID != "" {
		spanIDs[r.CurrentBuild.SpanID] = true
	}

	var sb strings.Builder
	for _, seg := range allLogs {
		if len(spanIDs) > 0 && !spanIDs[seg.SpanID] {
			continue
		}
		sb.WriteString(renderSegment(seg))
	}
	if sb.Len() == 0 {
		return dimStyle.Render("No logs for this resource")
	}
	return sb.String()
}

func renderSegment(seg api.LogSegment) string {
	text := strings.TrimRight(seg.Text, "\n")
	var line string
	switch seg.Level {
	case "ERROR":
		line = errorStyle.Render(text)
	case "WARN":
		line = warnStyle.Render(text)
	default:
		line = text
	}
	return line + "\n"
}

// ── Status helpers ────────────────────────────────────────────────────────────

func statusIcon(r api.Resource) (string, lipgloss.Style) {
	if !r.CurrentBuild.StartTime.IsZero() && r.CurrentBuild.FinishTime.IsZero() {
		return "⟳", buildingStyle
	}
	switch r.RuntimeStatus {
	case "ok":
		return "✓", okStyle
	case "error":
		return "✗", errorStyle
	case "pending":
		return "●", pendingStyle
	default:
		return "○", dimStyle
	}
}

func timeAgo(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	default:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
}

type statusCounts struct{ ok, errored, building, pending int }

func countStatuses(resources []api.Resource) statusCounts {
	var c statusCounts
	for _, r := range resources {
		if !r.CurrentBuild.StartTime.IsZero() && r.CurrentBuild.FinishTime.IsZero() {
			c.building++
			continue
		}
		switch r.RuntimeStatus {
		case "ok":
			c.ok++
		case "error":
			c.errored++
		case "pending":
			c.pending++
		}
	}
	return c
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if !m.ready {
		return "\n  Loading tilt-tui...\n"
	}

	header := m.renderHeader()
	menu := m.renderMenu()
	dot := m.dotLine()

	var panels string
	switch m.tab {
	case tabResources:
		left := m.renderResourceList()
		right := m.renderLogPanel()
		panels = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	case tabAllLogs:
		panels = m.renderAllLogsPanel()
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, dot, panels, dot, menu)
}

func (m Model) dotLine() string {
	var sb strings.Builder
	for i := 0; i < m.width; i++ {
		if i%2 == 0 {
			sb.WriteRune('·')
		} else {
			sb.WriteRune(' ')
		}
	}
	return dotStyle.Render(sb.String())
}

func (m Model) renderHeader() string {
	counts := countStatuses(m.resources)

	var badges []string
	if counts.ok > 0 {
		badges = append(badges, okStyle.Render(fmt.Sprintf("✓ %d", counts.ok)))
	}
	if counts.errored > 0 {
		badges = append(badges, errorStyle.Render(fmt.Sprintf("✗ %d", counts.errored)))
	}
	if counts.building > 0 {
		badges = append(badges, buildingStyle.Render(fmt.Sprintf("⟳ %d", counts.building)))
	}
	if counts.pending > 0 {
		badges = append(badges, pendingStyle.Render(fmt.Sprintf("● %d", counts.pending)))
	}
	if m.loading {
		badges = append(badges, dimStyle.Render("syncing…"))
	}
	if m.err != nil {
		badges = append(badges, errorStyle.Render("unreachable"))
	}

	left := headerStyle.Render("⚡ tilt-tui  " + dimStyle.Render(m.addr))
	right := headerRightStyle.Render(strings.Join(badges, "  "))
	fill := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if fill < 0 {
		fill = 0
	}
	return left + headerFillStyle.Render(strings.Repeat(" ", fill)) + right
}

func (m Model) renderMenu() string {
	type tabDef struct {
		id    activeTab
		label string
		key   string
	}
	tabs := []tabDef{
		{tabResources, "Resources", "1"},
		{tabAllLogs, "All Logs", "2"},
	}

	var parts []string
	for _, t := range tabs {
		label := fmt.Sprintf(" %s %s ", t.key, t.label)
		if m.tab == t.id {
			parts = append(parts, activeTabStyle.Render(label))
		} else {
			parts = append(parts, inactiveTabStyle.Render(label))
		}
	}

	tabBar := strings.Join(parts, "")
	shortcuts := menuShortcutsStyle.Render("↑↓/jk Navigate · Tab Switch · PgUp/PgDn Scroll · r Reload · q Quit")
	fill := m.width - lipgloss.Width(tabBar) - lipgloss.Width(shortcuts)
	if fill < 0 {
		fill = 0
	}
	return tabBar + menuFillStyle.Render(strings.Repeat(" ", fill)) + shortcuts
}

func (m Model) renderResourceList() string {
	var sb strings.Builder
	sb.WriteString(sectionTitleStyle.Render("Resources") + "\n")

	itemW := m.leftPanelWidth() - 2 // account for border

	for i, r := range m.resources {
		icon, style := statusIcon(r)

		name := r.Name
		if r.HasPendingChanges {
			name += " *"
		}
		line1 := fmt.Sprintf(" %s %s", icon, name)

		// Second line: status + last build time
		lastBuild := r.CurrentBuild.FinishTime
		if lastBuild.IsZero() && len(r.BuildHistory) > 0 {
			lastBuild = r.BuildHistory[0].FinishTime
		}
		status := r.RuntimeStatus
		if !r.CurrentBuild.StartTime.IsZero() && r.CurrentBuild.FinishTime.IsZero() {
			status = "building"
		}
		line2 := fmt.Sprintf("   %s · %s", status, timeAgo(lastBuild))

		if i == m.cursor {
			sb.WriteString(selectedItemStyle.Width(itemW).Render(line1) + "\n")
			sb.WriteString(selectedSubStyle.Width(itemW).Render(line2) + "\n")
		} else {
			sb.WriteString(style.Render(line1) + "\n")
			sb.WriteString(itemSubStyle.Render(line2) + "\n")
		}
	}

	if m.err != nil {
		sb.WriteString("\n" + errorStyle.Render(" "+m.err.Error()))
	}

	return listPanelStyle.
		Width(m.leftPanelWidth()).
		Height(m.panelHeight()).
		Render(sb.String())
}

func (m Model) renderLogPanel() string {
	name := "—"
	if len(m.resources) > 0 {
		name = m.resources[m.cursor].Name
	}
	title := sectionTitleStyle.Render("Logs: " + name)
	content := lipgloss.JoinVertical(lipgloss.Left, title, m.resourcesVP.View())

	return logPanelStyle.
		Width(m.rightPanelWidth()).
		Height(m.panelHeight()).
		Render(content)
}

func (m Model) renderAllLogsPanel() string {
	title := sectionTitleStyle.Render("All Logs")
	content := lipgloss.JoinVertical(lipgloss.Left, title, m.allLogsVP.View())

	return logPanelStyle.
		Width(m.width).
		Height(m.panelHeight()).
		Render(content)
}
