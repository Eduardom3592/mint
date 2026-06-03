package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/programmersd21/mint/internal/api"
	"github.com/programmersd21/mint/internal/cache"
	"github.com/programmersd21/mint/internal/downloads"
	"github.com/programmersd21/mint/internal/mrpack"
	"github.com/programmersd21/mint/internal/tui/theme"
)

type tickMsg int

type openProjectMsg struct {
	id string
}

type mrpackInspectMsg struct {
	index *mrpack.PackIndex
	meta  *mrpack.PackMetadata
	err   error
}

type mrpackInstallMsg struct {
	name       string
	installDir string
	downloadID int
	err        error
}

type downloadCompleteMsg struct {
	item *downloads.Item
}

type model struct {
	client *api.Client
	cache  *cache.Cache
	dlmgr  *downloads.Manager

	themeManager *theme.Manager

	width  int
	height int

	currentPage   navPage
	sidebarCursor int

	home          homeModel
	search        searchModel
	project       projectModel
	version       versionModel
	dependency    dependencyModel
	downloadView  downloadViewModel
	cacheView     cacheViewModel
	settingsView  settingsViewModel
	helpView      helpModel
	themeSwitcher themeSwitcherModel
	picker        versionPickerModel
	pickerActive  bool

	styles Styles

	frame        int
	loaded       bool
	completionCh chan *downloads.Item
}

func frameTick() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(1)
	})
}

func inspectMRPack(dlDir, url, filename string, size int64, hashes map[string]string) tea.Cmd {
	return func() tea.Msg {
		dest := filepath.Join(dlDir, filename)
		if _, err := os.Stat(dest); err != nil {
			return mrpackInspectMsg{err: fmt.Errorf("file not downloaded yet: download with d first")}
		}
		index, meta, err := mrpack.Parse(dest)
		return mrpackInspectMsg{index: index, meta: meta, err: err}
	}
}

func installMRPack(path, installDir string, downloadID int) tea.Cmd {
	return func() tea.Msg {
		name := filepath.Base(path)
		installPath := filepath.Join(installDir, "installed", name[:len(name)-len(".mrpack")])
		if err := os.MkdirAll(installPath, 0755); err != nil {
			return mrpackInstallMsg{name: name, installDir: installDir, downloadID: downloadID, err: fmt.Errorf("create install dir: %w", err)}
		}
		_, _, err := mrpack.Parse(path)
		if err != nil {
			return mrpackInstallMsg{name: name, installDir: installDir, downloadID: downloadID, err: fmt.Errorf("parse mrpack: %w", err)}
		}
		if err := mrpack.Extract(path, installPath, nil); err != nil {
			return mrpackInstallMsg{name: name, installDir: installDir, downloadID: downloadID, err: fmt.Errorf("extract: %w", err)}
		}
		if err := mrpack.ApplyAllOverrides(installPath, installPath); err != nil {
			return mrpackInstallMsg{name: name, installDir: installDir, downloadID: downloadID, err: fmt.Errorf("overrides: %w", err)}
		}
		return mrpackInstallMsg{name: name, installDir: installDir, downloadID: downloadID, err: nil}
	}
}

func New(client *api.Client, c *cache.Cache, dlmgr *downloads.Manager) *model {
	tm := theme.NewManager(c.DB())
	setGlobalTheme(tm.Current())

	ch := make(chan *downloads.Item, 10)
	dlmgr.OnComplete = func(item *downloads.Item) {
		select {
		case ch <- item:
		default:
		}
	}

	m := &model{
		client:        client,
		cache:         c,
		dlmgr:         dlmgr,
		themeManager:  tm,
		currentPage:   navHome,
		styles:        glob,
		home:          newHomeModel(client),
		search:        newSearchModel(client, c),
		project:       newProjectModel(client, c),
		version:       newVersionModel(client, c),
		dependency:    newDependencyModel(c),
		downloadView:  newDownloadViewModel(dlmgr),
		cacheView:     newCacheViewModel(c),
		settingsView:  newSettingsViewModel(c),
		helpView:      newHelpModel(),
		themeSwitcher: newThemeSwitcherModel(tm),
		picker:        newVersionPickerModel(client, dlmgr),
		completionCh:  ch,
	}

	m.syncStyles()

	return m
}

func (m *model) syncStyles() {
	s := NewStyles(m.themeManager.Current())
	m.styles = s
	m.home.setStyles(s)
	m.search.setStyles(s)
	m.project.setStyles(s)
	m.version.setStyles(s)
	m.dependency.setStyles(s)
	m.downloadView.setStyles(s)
	m.cacheView.setStyles(s)
	m.settingsView.setStyles(s)
	m.helpView.setStyles(s)
	m.themeSwitcher.setStyles(s)
	m.picker.setStyles(s)
}

func waitForDownloadComplete(ch chan *downloads.Item) tea.Cmd {
	return func() tea.Msg {
		item := <-ch
		return downloadCompleteMsg{item: item}
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(
		frameTick(),
		waitForDownloadComplete(m.completionCh),
		m.home.Init(),
		m.search.Init(),
		m.project.Init(),
		m.version.Init(),
		m.dependency.Init(),
		m.downloadView.Init(),
		m.cacheView.Init(),
		m.settingsView.Init(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case themeChangedMsg:
		m.syncStyles()
		return m, nil

	case tickMsg:
		m.frame++
		m.home.frame = m.frame
		m.search.frame = m.frame
		m.project.frame = m.frame
		m.version.frame = m.frame
		return m, frameTick()

	case mrpackInspectMsg:
		if msg.err != nil {
			m.version.err = msg.err
		} else if msg.index != nil {
			m.version.mrpackInfo = msg.index
		}
		return m, nil

	case downloadCompleteMsg:
		item := msg.item
		if filepath.Ext(item.Filename) == ".mrpack" && item.Status == downloads.StatusCompleted {
			return m, installMRPack(item.DestPath(), m.dlmgr.Dir(), item.ID)
		}
		return m, waitForDownloadComplete(m.completionCh)

	case mrpackInstallMsg:
		if msg.err == nil && msg.downloadID > 0 {
			if p := m.dlmgr.Persist(); p != nil {
				_ = p.MarkInstalled(int64(msg.downloadID))
			}
		}
		return m, waitForDownloadComplete(m.completionCh)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.loaded = true
		m.home.width = contentWidth(m.width)
		m.search.height = contentHeight(m.height)
		m.search.width = contentWidth(m.width)
		m.search.textInput.Width = m.search.width - 4
		if m.search.textInput.Width < 5 {
			m.search.textInput.Width = 5
		}
		m.project.width = contentWidth(m.width)
		if m.project.showBody && m.project.renderedLines != nil && m.project.width != m.project.renderWidth {
			m.project.renderBody()
		}
		m.picker.width = contentWidth(m.width)

	case homeLoadedMsg:
		model, _ := m.home.Update(msg)
		m.home = model
		return m, nil

	case homeErrorMsg:
		model, _ := m.home.Update(msg)
		m.home = model
		return m, nil

	case cacheLoadedMsg:
		model, cmd := m.cacheView.Update(msg)
		m.cacheView = model
		return m, cmd

	case openProjectMsg:
		m.setPage(navProject)
		return m, m.project.load(msg.id)

	case searchResultsMsg:
		model, _ := m.search.Update(msg)
		m.search = model
		return m, nil

	case searchErrorMsg:
		model, _ := m.search.Update(msg)
		m.search = model
		return m, nil

	case projectLoadedMsg:
		model, _ := m.project.Update(msg)
		m.project = model
		return m, nil

	case projectErrorMsg:
		model, _ := m.project.Update(msg)
		m.project = model
		return m, nil

	case versionsLoadedMsg:
		model, _ := m.version.Update(msg)
		m.version = model
		return m, nil

	case versionErrorMsg:
		model, _ := m.version.Update(msg)
		m.version = model
		return m, nil

	case pickerVersionsMsg:
		model, cmd := m.picker.Update(msg)
		m.picker = model
		if !m.pickerActive {
			return m, nil
		}
		if m.picker.allVersions == nil {
			m.pickerActive = false
			m.setPage(navDownloads)
		}
		return m, cmd

	case pickerErrorMsg:
		model, cmd := m.picker.Update(msg)
		m.picker = model
		return m, cmd

	case tea.KeyMsg:
		if m.pickerActive {
			model, cmd := m.picker.Update(msg)
			m.picker = model
			if m.picker.allVersions == nil {
				m.pickerActive = false
				m.setPage(navDownloads)
			}
			return m, cmd
		}

		if m.themeSwitcher.visible {
			model, cmd := m.themeSwitcher.Update(msg)
			m.themeSwitcher = model
			return m, cmd
		}

		if (key.Matches(msg, keys.Tab) || key.Matches(msg, keys.ShiftTab)) && m.currentPage == navProject {
			break
		}

		switch {
		case key.Matches(msg, keys.Quit) && !m.inputFocused():
			return m, tea.Quit
		case key.Matches(msg, keys.Help) && !m.inputFocused():
			m.setPage(navHelp)
			return m, nil
		case key.Matches(msg, keys.Theme) && !m.inputFocused():
			m.themeSwitcher.visible = true
			m.themeSwitcher.cursor = 0
			current := m.themeManager.Current()
			for i, t := range m.themeSwitcher.themes {
				if t.Name == current.Name {
					m.themeSwitcher.cursor = i
					break
				}
			}
			return m, nil
		case key.Matches(msg, keys.Download) && m.currentPage == navVersions && m.version.detail != nil && !m.inputFocused():
			if m.version.fileCursor >= 0 && m.version.fileCursor < len(m.version.detail.Files) {
				f := m.version.detail.Files[m.version.fileCursor]
				var hash *downloads.HashInfo
				if h, ok := f.Hashes["sha1"]; ok {
					hash = &downloads.HashInfo{Type: downloads.HashSHA1, Value: h}
				} else if h, ok := f.Hashes["sha512"]; ok {
					hash = &downloads.HashInfo{Type: downloads.HashSHA512, Value: h}
				}
				m.dlmgr.Enqueue(m.version.projectID, "", m.version.detail.ID, m.version.detail.VersionNumber, f.URL, f.Filename, f.Size, hash)
			}
			return m, nil

		case key.Matches(msg, keys.Download) && m.currentPage.supportsDownload() && !m.inputFocused():
			id, title := m.selectedProject()
			if id != "" {
				m.pickerActive = true
				return m, m.picker.loadVersions(id, title, false)
			}
			return m, nil

		case key.Matches(msg, keys.QuickDL) && m.currentPage.supportsDownload() && !m.inputFocused():
			id, title := m.selectedProject()
			if id != "" {
				m.pickerActive = true
				return m, m.picker.loadVersions(id, title, true)
			}
			return m, nil

		case key.Matches(msg, keys.Inspect) && m.currentPage == navVersions && m.version.detail != nil && !m.inputFocused():
			if m.version.fileCursor >= 0 && m.version.fileCursor < len(m.version.detail.Files) {
				f := m.version.detail.Files[m.version.fileCursor]
				if filepath.Ext(f.Filename) == ".mrpack" {
					return m, inspectMRPack(m.dlmgr.Dir(), f.URL, f.Filename, f.Size, f.Hashes)
				}
			}
			return m, nil

		case key.Matches(msg, keys.Install) && m.currentPage == navVersions && m.version.detail != nil && !m.inputFocused():
			if m.version.fileCursor >= 0 && m.version.fileCursor < len(m.version.detail.Files) {
				f := m.version.detail.Files[m.version.fileCursor]
				if filepath.Ext(f.Filename) == ".mrpack" {
					var hash *downloads.HashInfo
					if h, ok := f.Hashes["sha1"]; ok {
						hash = &downloads.HashInfo{Type: downloads.HashSHA1, Value: h}
					} else if h, ok := f.Hashes["sha512"]; ok {
						hash = &downloads.HashInfo{Type: downloads.HashSHA512, Value: h}
					}
					m.dlmgr.Enqueue(m.version.projectID, "", m.version.detail.ID, m.version.detail.VersionNumber, f.URL, f.Filename, f.Size, hash)
				}
			}
			return m, nil

		case key.Matches(msg, keys.Tab):
			m.nextPage()
			return m, nil
		case key.Matches(msg, keys.ShiftTab):
			m.prevPage()
			return m, nil
		}

	case tea.MouseMsg:
		return m, m.handleMouse(msg)
	}

	var cmd tea.Cmd
	switch m.currentPage {
	case navHome:
		cmd = m.updateHome(msg)
	case navSearch:
		cmd = m.updateSearch(msg)
	case navProject:
		cmd = m.updateProject(msg)
	case navVersions:
		cmd = m.updateVersion(msg)
	case navDependencies:
		cmd = m.updateDependency(msg)
	case navDownloads:
		cmd = m.updateDownloadView(msg)
	case navCache:
		cmd = m.updateCacheView(msg)
	case navSettings:
		cmd = m.updateSettingsView(msg)
	case navHelp:
		cmd = m.updateHelp(msg)

	}

	return m, cmd
}

func dots(frame int) string {
	n := (frame / 3) % 4
	return strings.Repeat(".", n)
}

func totalWidth(mw int) int {
	if mw > 120 {
		return 120
	}
	return mw
}

func sidebarWidth() int {
	return 20
}

func contentWidth(mw int) int {
	tw := totalWidth(mw)
	w := tw - sidebarWidth() - 1
	if w < 12 {
		w = 12
	}
	return w
}

func contentHeight(mh int) int {
	h := mh - 4
	if h < 5 {
		h = 5
	}
	return h
}

func hyperlink(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

func openURL(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	_ = exec.Command(cmd, args...).Start()
}

func (m *model) inputFocused() bool {
	if m.currentPage == navSearch && m.search.textInput.Focused() {
		return true
	}
	if m.currentPage == navSettings && m.settingsView.focus >= 0 {
		return true
	}
	return false
}

func (m *model) selectedProject() (id, title string) {
	switch m.currentPage {
	case navHome:
		if len(m.home.hits) > 0 && m.home.cursor >= 0 && m.home.cursor < len(m.home.hits) {
			h := m.home.hits[m.home.cursor]
			return h.ProjectID, h.Title
		}
	case navSearch:
		if len(m.search.results) > 0 && m.search.cursor >= 0 && m.search.cursor < len(m.search.results) {
			h := m.search.results[m.search.cursor]
			return h.ProjectID, h.Title
		}
	case navProject:
		if m.project.project != nil {
			return m.project.project.ID, m.project.project.Title
		}
	}
	return "", ""
}

func scrollLine(pos, total, visible int) string {
	if total <= visible || total <= 0 {
		return ""
	}
	return fmt.Sprintf("-- %d/%d --", pos+1, total)
}

func (m *model) View() string {
	if !m.loaded {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#8B8FA3")).Render(
			fmt.Sprintf("mint warming up%s", dots(m.frame)),
		)
	}

	s := m.styles

	sidebar := m.renderSidebar()
	content := m.renderContent()
	statusBar := m.renderStatusBar()

	sw := sidebarWidth()
	cw := contentWidth(m.width)
	ch := contentHeight(m.height)
	tw := totalWidth(m.width)

	contentClipped := lipgloss.NewStyle().
		MaxHeight(ch).
		Render(content)

	main := lipgloss.JoinHorizontal(
		lipgloss.Top,
		s.Sidebar.Width(sw).Render(sidebar),
		s.Content.Width(cw).Align(lipgloss.Center).Render(contentClipped),
	)

	statusBarRendered := s.StatusBar.Width(tw).Render(statusBar)

	inner := s.App.Width(tw).Render(lipgloss.JoinVertical(lipgloss.Top, main, statusBarRendered))

	view := lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Top,
		inner,
	)

	if m.themeSwitcher.visible {
		overlay := m.themeSwitcher.View()
		if overlay != "" {
			view = lipgloss.Place(
				m.width, m.height,
				lipgloss.Center, lipgloss.Center,
				overlay,
			)
		}
	}

	if m.pickerActive {
		overlay := m.picker.View()
		if overlay != "" {
			pickerContent := lipgloss.NewStyle().
				MaxHeight(m.height - 4).
				Width(contentWidth(m.width)).
				Render(overlay)
			view = lipgloss.Place(
				m.width, m.height,
				lipgloss.Center, lipgloss.Center,
				pickerContent,
			)
		}
	}

	return view
}

func (m *model) nextPage() {
	pages := []navPage{navHome, navSearch, navDownloads, navCache, navSettings, navHelp}
	for i, p := range pages {
		if p == m.currentPage {
			m.setPage(pages[(i+1)%len(pages)])
			m.sidebarCursor = (i + 1) % len(pages)
			return
		}
	}
}

func (m *model) prevPage() {
	pages := []navPage{navHome, navSearch, navDownloads, navCache, navSettings, navHelp}
	for i, p := range pages {
		if p == m.currentPage {
			m.setPage(pages[(i-1+len(pages))%len(pages)])
			m.sidebarCursor = (i - 1 + len(pages)) % len(pages)
			return
		}
	}
}

func (m *model) setPage(page navPage) {
	m.currentPage = page
	if page == navSearch {
		m.search.page = searchList
		m.search.selectedProjectID = ""
	}
}

func (m *model) renderSidebar() string {
	var sb strings.Builder
	s := m.styles

	sb.WriteString(s.Title.Render("mint"))
	sb.WriteString("\n\n")

	items := []struct {
		label string
		page  navPage
	}{
		{"home", navHome},
		{"search", navSearch},
		{"downloads", navDownloads},

		{"cache", navCache},
		{"settings", navSettings},
		{"help", navHelp},
	}

	for _, item := range items {
		if item.page == m.currentPage {
			sb.WriteString(s.ActiveItem.Render(fmt.Sprintf(" > %s", item.label)))
		} else {
			sb.WriteString(s.Inactive.Render(fmt.Sprintf("   %s", item.label)))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	if dlmgr := m.dlmgr; dlmgr != nil {
		active := dlmgr.ActiveCount()
		if active > 0 {
			sb.WriteString(s.Warning.Render(fmt.Sprintf("↓ %d active", active)))
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(s.Dimmed.Render("tab nav"))
	sb.WriteString("\n")
	sb.WriteString(s.Dimmed.Render("t theme"))

	return sb.String()
}

func (m *model) renderContent() string {
	var content string
	switch m.currentPage {
	case navHome:
		content = m.home.View()
	case navSearch:
		content = m.search.View()
	case navProject:
		content = m.project.View()
	case navVersions:
		content = m.version.View()
	case navDependencies:
		content = m.dependency.View()
	case navDownloads:
		content = m.downloadView.View()
	case navCache:
		content = m.cacheView.View()
	case navSettings:
		content = m.settingsView.View()
	case navHelp:
		content = m.helpView.View()

	default:
		content = "unknown page"
	}
	return content
}

func (m *model) renderStatusBar() string {
	var parts []string
	s := m.styles

	cacheProjects, cacheVersions, _, _ := m.cache.Stats()

	parts = append(parts, s.Info.Render(fmt.Sprintf("%d", cacheProjects)))
	parts = append(parts, s.Info.Render(fmt.Sprintf("%d", cacheVersions)))

	if active := m.dlmgr.ActiveCount(); active > 0 {
		parts = append(parts, s.Info.Render(fmt.Sprintf("↓ %d", active)))
	}

	parts = append(parts, s.Info.Render(m.themeManager.Current().Name))

	return s.StatusBar.Width(m.width).Render(
		strings.Join(parts, "  "),
	)
}

func (m *model) updateHome(msg tea.Msg) tea.Cmd {
	model, cmd := m.home.Update(msg)
	m.home = model
	if m.home.selectedProjectID != "" {
		id := m.home.selectedProjectID
		m.home.selectedProjectID = ""
		m.setPage(navProject)
		return tea.Batch(cmd, m.project.load(id))
	}
	return cmd
}

func (m *model) updateSearch(msg tea.Msg) tea.Cmd {
	oldPage := m.search.page
	model, cmd := m.search.Update(msg)
	m.search = model
	if m.search.page != oldPage {
		switch m.search.page {
		case searchProject:
			m.setPage(navProject)
			return tea.Batch(cmd, m.project.load(m.search.selectedProjectID))
		case searchVersion:
			m.setPage(navVersions)
		}
	}
	return cmd
}

func (m *model) updateProject(msg tea.Msg) tea.Cmd {
	oldPage := m.project.page
	model, cmd := m.project.Update(msg)
	m.project = model
	if m.project.page != oldPage {
		switch m.project.page {
		case projectSearch:
			m.setPage(navSearch)
		case projectVersion:
			m.setPage(navVersions)
			if m.project.project != nil {
				return tea.Batch(cmd, m.version.load(m.project.project.ID))
			}
		case projectDependency:
			m.setPage(navDependencies)
		case projectDownload:
			m.setPage(navDownloads)
		}
	}
	return cmd
}

func (m *model) updateVersion(msg tea.Msg) tea.Cmd {
	oldPage := m.version.page
	model, cmd := m.version.Update(msg)
	m.version = model
	if m.version.page != oldPage {
		switch m.version.page {
		case versionProject:
			m.setPage(navProject)
		case versionDependency:
			m.setPage(navDependencies)
		case versionSearch:
			m.setPage(navSearch)
		case versionDownload:
			m.setPage(navDownloads)
		}
	}
	return cmd
}

func (m *model) updateDependency(msg tea.Msg) tea.Cmd {
	model, cmd := m.dependency.Update(msg)
	m.dependency = model
	return cmd
}

func (m *model) updateDownloadView(msg tea.Msg) tea.Cmd {
	model, cmd := m.downloadView.Update(msg)
	m.downloadView = model
	return cmd
}

func (m *model) updateCacheView(msg tea.Msg) tea.Cmd {
	model, cmd := m.cacheView.Update(msg)
	m.cacheView = model
	return cmd
}

func (m *model) updateSettingsView(msg tea.Msg) tea.Cmd {
	model, cmd := m.settingsView.Update(msg)
	m.settingsView = model
	return cmd
}

func (m *model) updateHelp(msg tea.Msg) tea.Cmd {
	model, cmd := m.helpView.Update(msg)
	m.helpView = model
	return cmd
}
