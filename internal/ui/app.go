package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmeiracorbal/tt/internal/api"
)

const pollInterval = 3 * time.Second

// minWidth / minHeight: below these the layout breaks down.
const minWidth = 40
const minHeight = 8

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

// panelHeight is rows available for the panel area: total minus header and footer.
func (m Model) panelHeight() int {
	h := m.height - 2
	if h < 4 {
		h = 4
	}
	return h
}

// resourcesVPSize: right panel content minus borders (2) and title row (1).
func (m Model) resourcesVPSize() (int, int) {
	w := m.rightPanelWidth() - 2
	if w < 10 {
		w = 10
	}
	h := m.panelHeight() - 3 // borders(2) + title(1)
	if h < 1 {
		h = 1
	}
	return w, h
}

// allLogsVPSize: full-width panel minus borders (2) and title row (1).
func (m Model) allLogsVPSize() (int, int) {
	w := m.width - 2
	if w < 10 {
		w = 10
	}
	h := m.panelHeight() - 3
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

	// Build spanID → resource map for banner detection.
	type resInfo struct {
		name string
		role string // ok / fail / working / pending
	}
	spanToRes := make(map[string]resInfo)
	for _, r := range m.resources {
		info := resInfo{name: r.Name, role: resourceRole(r)}
		if r.CurrentBuild.SpanID != "" {
			spanToRes[r.CurrentBuild.SpanID] = info
		}
		for _, b := range r.BuildHistory {
			if b.SpanID != "" {
				spanToRes[b.SpanID] = info
			}
		}
	}

	var sb strings.Builder
	lastRes := ""
	for _, seg := range m.allLogs {
		if info, ok := spanToRes[seg.SpanID]; ok && info.name != lastRes {
			sb.WriteString(renderBanner(info.name, info.role, m.width) + "\n")
			lastRes = info.name
		}
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

// renderBanner renders a resource section header for the All Logs view:
// "─── name  glyph ──────────────────────────"
func renderBanner(name, role string, width int) string {
	glyph, style := bannerIconAndStyle(role)
	label := fmt.Sprintf("─── %s  %s ", name, glyph)
	runeLen := len([]rune(label))
	fill := width - runeLen - 1
	if fill < 0 {
		fill = 0
	}
	return style.Render(label + strings.Repeat("─", fill))
}

func bannerIconAndStyle(role string) (string, lipgloss.Style) {
	switch role {
	case "ok":
		return "✓", bannerOkStyle
	case "fail":
		return "✗", bannerErrorStyle
	case "working":
		return "⟳", bannerBuildingStyle
	default:
		return "·", bannerDimStyle
	}
}

// ── Status helpers ────────────────────────────────────────────────────────────

// resourceRole returns one of: ok / fail / working / pending / unknown.
func resourceRole(r api.Resource) string {
	if !r.CurrentBuild.StartTime.IsZero() && r.CurrentBuild.FinishTime.IsZero() {
		return "working"
	}
	switch r.RuntimeStatus {
	case "ok":
		return "ok"
	case "error":
		return "fail"
	case "pending":
		return "pending"
	default:
		return "unknown"
	}
}

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
		return "·", dimStyle
	}
}

func timeAgo(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	default:
		return fmt.Sprintf("%dh", int(d.Hours()))
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
		return "\n  Loading tt...\n"
	}

	if m.width < minWidth || m.height < minHeight {
		return m.renderTooSmall()
	}

	header := m.renderHeader()
	footer := m.renderFooter()

	var body string
	switch m.tab {
	case tabResources:
		left := m.renderResourceList()
		right := m.renderLogPanel()
		body = lipgloss.JoinHorizontal(lipgloss.Top, left, right)
	case tabAllLogs:
		body = m.renderAllLogsPanel()
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m Model) renderTooSmall() string {
	msg := fmt.Sprintf(
		"\n\n  %s\n  Needs %dx%d  (current: %dx%d)\n\n  %s\n",
		errorStyle.Render("Terminal too small"),
		minWidth, minHeight,
		m.width, m.height,
		dimStyle.Render("q quit"),
	)
	return msg
}

func (m Model) renderHeader() string {
	counts := countStatuses(m.resources)

	left := headerTitleStyle.Render(" tt ") + headerAddrStyle.Render(" "+m.addr)

	var badges []string
	if counts.ok > 0 {
		badges = append(badges, okStyle.Render(fmt.Sprintf(" ✓ %d ", counts.ok)))
	}
	if counts.errored > 0 {
		badges = append(badges, errorStyle.Render(fmt.Sprintf(" ✗ %d ", counts.errored)))
	}
	if counts.building > 0 {
		badges = append(badges, buildingStyle.Render(fmt.Sprintf(" ⟳ %d ", counts.building)))
	}
	if counts.pending > 0 {
		badges = append(badges, pendingStyle.Render(fmt.Sprintf(" ● %d ", counts.pending)))
	}
	if m.loading {
		badges = append(badges, dimStyle.Render(" syncing… "))
	}
	if m.err != nil {
		badges = append(badges, errorStyle.Render(" unreachable "))
	}
	right := strings.Join(badges, "")

	fill := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if fill < 0 {
		fill = 0
	}
	return left + headerFillStyle.Render(strings.Repeat(" ", fill)) + right
}

// renderFooter renders the bottom tab bar + key hints.
func (m Model) renderFooter() string {
	var tab1, tab2 string
	if m.tab == tabResources {
		tab1 = activeTabStyle.Render("1 Resources")
		tab2 = inactiveTabStyle.Render("2 All Logs")
	} else {
		tab1 = inactiveTabStyle.Render("1 Resources")
		tab2 = activeTabStyle.Render("2 All Logs")
	}
	tabs := tab1 + tab2

	var hints string
	if m.tab == tabResources {
		hints = renderHint("↑↓", "select") + "  " +
			renderHint("PgUp/Dn", "scroll") + "  " +
			renderHint("r", "reload")
	} else {
		hints = renderHint("PgUp/Dn", "scroll") + "  " +
			renderHint("r", "reload")
	}
	rightHints := renderHint("?", "help") + "  " + renderHint("q", "quit")

	tabsW := lipgloss.Width(tabs)
	hintsW := lipgloss.Width(hints)
	rightW := lipgloss.Width(rightHints)

	fill := m.width - tabsW - hintsW - rightW - 2
	if fill < 0 {
		fill = 0
	}
	return tabs + " " + hints + keybarFillStyle.Render(strings.Repeat(" ", fill)) + " " + rightHints
}

func renderHint(key, verb string) string {
	return keyStyle.Render(key) + " " + descStyle.Render(verb)
}

func (m Model) renderResourceList() string {
	var sb strings.Builder

	innerW := m.leftPanelWidth() - 2 // subtract left+right border
	// layout per row: cursor(1) + sp(1) + glyph(1) + sp(1) + name(flex) + sp(1) + age(4) = 9 fixed
	ageW := 4
	nameW := innerW - 9
	if nameW < 1 {
		nameW = 1
	}

	for i, r := range m.resources {
		icon, iconStyle := statusIcon(r)

		name := r.Name
		if r.HasPendingChanges {
			name += " *"
		}
		if len([]rune(name)) > nameW {
			runes := []rune(name)
			name = string(runes[:nameW-1]) + "⋯"
		}

		lastBuild := r.CurrentBuild.FinishTime
		if lastBuild.IsZero() && len(r.BuildHistory) > 0 {
			lastBuild = r.BuildHistory[0].FinishTime
		}
		age := timeAgo(lastBuild)

		if i == m.cursor {
			line := fmt.Sprintf("❯ %s %-*s %*s", icon, nameW, name, ageW, age)
			sb.WriteString(selectedItemStyle.Width(innerW).Render(line) + "\n")
		} else {
			line := "  " + iconStyle.Render(icon) +
				fmt.Sprintf(" %-*s ", nameW, name) +
				itemAgeStyle.Render(fmt.Sprintf("%*s", ageW, age))
			sb.WriteString(line + "\n")
		}
	}

	if m.err != nil {
		sb.WriteString("\n" + errorStyle.Render("  "+m.err.Error()))
	}

	return focusedPanelStyle.
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

	return unfocusedPanelStyle.
		Width(m.rightPanelWidth()).
		Height(m.panelHeight()).
		Render(content)
}

func (m Model) renderAllLogsPanel() string {
	title := sectionTitleStyle.Render("All Logs")
	content := lipgloss.JoinVertical(lipgloss.Left, title, m.allLogsVP.View())

	return unfocusedPanelStyle.
		Width(m.width).
		Height(m.panelHeight()).
		Render(content)
}
