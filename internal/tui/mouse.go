package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/downloads"
	"github.com/programmersd21/mint/internal/platform"
)

const (
	mouseSidebarFirstRow = 2
	mouseHelpRows        = 3
)

func (m *model) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if !m.loaded {
		return nil
	}
	if msg.Action != tea.MouseActionPress {
		return nil
	}

	if m.pickerActive {
		return m.handlePickerMouse(msg)
	}
	if m.themeSwitcher.visible {
		return m.handleThemeSwitcherMouse(msg)
	}

	if page, ok := m.sidebarPageAt(msg); ok && msg.Button == tea.MouseButtonLeft {
		m.setPage(page)
		return nil
	}

	switch msg.Button {
	case tea.MouseButtonWheelUp:
		return m.mouseWheel(-1)
	case tea.MouseButtonWheelDown:
		return m.mouseWheel(1)
	case tea.MouseButtonLeft:
		return m.mouseClick(msg)
	case tea.MouseButtonRight:
		return m.mouseContextAction()
	}

	return nil
}

func (m *model) sidebarPageAt(msg tea.MouseMsg) (navPage, bool) {
	x0, _, ok := m.frameOrigin()
	if !ok || msg.X < x0 || msg.X >= x0+sidebarWidth() {
		return 0, false
	}
	pages := []navPage{navHome, navSearch, navDownloads, navCache, navSettings, navHelp}
	row := msg.Y - mouseSidebarFirstRow
	if row < 0 || row >= len(pages) {
		return 0, false
	}
	return pages[row], true
}

func (m *model) frameOrigin() (int, int, bool) {
	tw := totalWidth(m.width)
	if m.width <= 0 || tw <= 0 {
		return 0, 0, false
	}
	x0 := (m.width - tw) / 2
	if x0 < 0 {
		x0 = 0
	}
	return x0, 0, true
}

func (m *model) contentPoint(msg tea.MouseMsg) (int, int, bool) {
	x0, y0, ok := m.frameOrigin()
	if !ok {
		return 0, 0, false
	}
	cx := msg.X - x0 - sidebarWidth() - 1
	cy := msg.Y - y0
	if cx < 0 || cx >= contentWidth(m.width) || cy < 0 || cy >= contentHeight(m.height)+mouseHelpRows {
		return 0, 0, false
	}
	return cx, cy, true
}

func (m *model) mouseWheel(delta int) tea.Cmd {
	switch m.currentPage {
	case navHome:
		if delta < 0 && m.home.cursor > 0 {
			m.home.cursor--
		} else if delta > 0 && m.home.cursor < len(m.home.hits)-1 {
			m.home.cursor++
		}
	case navSearch:
		if m.search.textInput.Focused() {
			m.search.textInput.Blur()
		}
		if delta < 0 && m.search.cursor > 0 {
			m.search.cursor--
		} else if delta > 0 && m.search.cursor < len(m.search.results)-1 {
			m.search.cursor++
		}
	case navProject:
		if m.project.showBody {
			maxScroll := len(m.project.renderedLines) - bodyVisibleLines
			if maxScroll < 0 {
				maxScroll = 0
			}
			m.project.bodyScroll = clamp(m.project.bodyScroll+delta, 0, maxScroll)
		} else {
			m.project.urlCursor = clamp(m.project.urlCursor+delta, 0, len(m.project.urls)-1)
		}
	case navVersions:
		if m.version.page == versionDetail {
			if m.version.detail != nil {
				m.version.fileCursor = clamp(m.version.fileCursor+delta, 0, len(m.version.detail.Files)-1)
			}
		} else {
			m.version.selected = clamp(m.version.selected+delta, 0, len(m.version.versions)-1)
		}
	case navDownloads:
		items := m.downloadView.items()
		m.downloadView.cursor = clamp(m.downloadView.cursor+delta, 0, len(items)-1)
	case navCache:
		m.cacheView.cursor = clamp(m.cacheView.cursor+delta, 0, len(m.cacheView.recent)-1)
	case navSettings:
		m.settingsView.cursor = clamp(m.settingsView.cursor+delta, 0, 3)
	}
	return nil
}

func (m *model) mouseClick(msg tea.MouseMsg) tea.Cmd {
	_, y, ok := m.contentPoint(msg)
	if !ok {
		return nil
	}

	if y >= contentHeight(m.height)-1 {
		return m.mouseFooterAction()
	}

	switch m.currentPage {
	case navHome:
		return m.homeMouseClick(y)
	case navSearch:
		return m.searchMouseClick(y)
	case navProject:
		return m.projectMouseClick(y)
	case navVersions:
		return m.versionMouseClick(y)
	case navDownloads:
		return m.downloadMouseClick(msg, y)
	case navCache:
		return m.cacheMouseClick(y)
	case navSettings:
		return m.settingsMouseClick(y)
	case navHelp:
		return nil
	}
	return nil
}

func (m *model) mouseFooterAction() tea.Cmd {
	switch m.currentPage {
	case navHome, navSearch, navProject:
		return m.startSelectedProjectDownload(false)
	case navVersions:
		return m.downloadSelectedVersionFile()
	case navDownloads:
		return m.downloadMouseAction()
	case navSettings:
		model, cmd := m.settingsView.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m.settingsView = model
		return cmd
	}
	return nil
}

func (m *model) mouseContextAction() tea.Cmd {
	switch m.currentPage {
	case navProject:
		if !m.project.showBody && len(m.project.urls) > 0 {
			openURL(m.project.urls[m.project.urlCursor].url)
		}
	case navHome, navSearch:
		return m.startSelectedProjectDownload(true)
	case navVersions:
		if m.version.page == versionDetail {
			openSelectedVersionFile(m.version)
		}
	case navDownloads:
		return m.downloadMouseAction()
	}
	return nil
}

func (m *model) homeMouseClick(y int) tea.Cmd {
	if len(m.home.hits) == 0 {
		return nil
	}
	idx := (y - 2) / 4
	if idx >= 0 && idx < len(m.home.hits) {
		if m.home.cursor == idx {
			m.home.selectedProjectID = m.home.hits[idx].ProjectID
			return m.updateHome(tea.KeyMsg{Type: tea.KeyEnter})
		}
		m.home.cursor = idx
	}
	return nil
}

func (m *model) searchMouseClick(y int) tea.Cmd {
	if y == 2 {
		m.search.textInput.Focus()
		return nil
	}
	if len(m.search.results) == 0 {
		return nil
	}
	idx := searchIndexAtY(m.search, y)
	if idx >= 0 && idx < len(m.search.results) {
		if m.search.cursor == idx {
			model, cmd := m.search.Update(tea.KeyMsg{Type: tea.KeyEnter})
			m.search = model
			if m.search.page == searchProject {
				m.setPage(navProject)
				return tea.Batch(cmd, m.project.load(m.search.selectedProjectID))
			}
			return cmd
		}
		m.search.cursor = idx
		m.search.textInput.Blur()
	}
	return nil
}

func searchIndexAtY(s searchModel, y int) int {
	itemsPerPage := (s.height - 9) / 3
	if itemsPerPage < 3 {
		itemsPerPage = 3
	}
	start := s.cursor - itemsPerPage/2
	if start < 0 {
		start = 0
	}
	end := start + itemsPerPage
	if end > len(s.results) {
		end = len(s.results)
		start = end - itemsPerPage
		if start < 0 {
			start = 0
		}
	}
	row := y - 7
	if row < 0 {
		return -1
	}
	return start + row/3
}

func (m *model) projectMouseClick(y int) tea.Cmd {
	if m.project.project == nil {
		return nil
	}
	if y >= contentHeight(m.height)-2 {
		m.project.showBody = !m.project.showBody
		m.project.bodyScroll = 0
		if m.project.showBody && m.project.renderedLines == nil {
			m.project.renderBody()
		}
		return nil
	}
	if !m.project.showBody {
		urlStart := 8
		if len(m.project.project.Loaders) > 0 {
			urlStart++
		}
		if len(m.project.project.GameVersions) > 0 {
			urlStart++
		}
		if len(m.project.project.Categories) > 0 {
			urlStart++
		}
		urlStart += 2
		idx := y - urlStart
		if idx >= 0 && idx < len(m.project.urls) {
			if m.project.urlCursor == idx {
				openURL(m.project.urls[idx].url)
			}
			m.project.urlCursor = idx
		}
	}
	return nil
}

func (m *model) versionMouseClick(y int) tea.Cmd {
	if m.version.page == versionDetail {
		idx := y - versionFileStartY(m.version)
		if m.version.detail != nil && idx >= 0 && idx < len(m.version.detail.Files) {
			m.version.fileCursor = idx
			return nil
		}
		return nil
	}
	if len(m.version.versions) == 0 {
		return nil
	}
	idx := versionIndexAtY(m.version, y)
	if idx >= 0 && idx < len(m.version.versions) {
		if m.version.selected == idx {
			model, cmd := m.version.Update(tea.KeyMsg{Type: tea.KeyEnter})
			m.version = model
			return cmd
		}
		m.version.selected = idx
	}
	return nil
}

func versionIndexAtY(v versionModel, y int) int {
	height := 20
	start := v.selected - height/2
	if start < 0 {
		start = 0
	}
	end := start + height
	if end > len(v.versions) {
		end = len(v.versions)
		start = end - height
		if start < 0 {
			start = 0
		}
	}
	row := y - 2
	if row < 0 {
		return -1
	}
	return start + row/2
}

func versionFileStartY(v versionModel) int {
	if v.detail == nil {
		return 0
	}
	y := 10
	if len(v.detail.Dependencies) > 0 {
		y += len(v.detail.Dependencies) + 2
	}
	return y
}

func (m *model) downloadMouseClick(msg tea.MouseMsg, y int) tea.Cmd {
	if y == 2 {
		tab := downloadTabAtX(m, msg.X)
		if tab >= 0 {
			m.downloadView.tab = tab
			m.downloadView.cursor = 0
			m.downloadView.scrollTop = 0
		}
		return nil
	}
	items := m.downloadView.items()
	idx := y - 4
	if idx >= 0 && idx < len(items) {
		if m.downloadView.cursor == idx {
			return m.downloadMouseAction()
		}
		m.downloadView.cursor = idx
	}
	return nil
}

func downloadTabAtX(m *model, screenX int) downloadTab {
	// The rendered tab labels occupy one line in fixed order; use coarse ranges
	// so terminal styling width differences do not make tabs hard to click.
	x0, _, ok := m.frameOrigin()
	if !ok {
		return -1
	}
	x := screenX - x0 - sidebarWidth() - 1
	switch {
	case x >= 0 && x < 10:
		return tabActive
	case x >= 10 && x < 22:
		return tabHistory
	case x >= 22 && x < 43:
		return tabInstalledModpacks
	case x >= 43 && x < 55:
		return tabFailed
	default:
		return -1
	}
}

func (m *model) downloadMouseAction() tea.Cmd {
	items := m.downloadView.items()
	if m.downloadView.cursor < 0 || m.downloadView.cursor >= len(items) {
		return nil
	}
	item := items[m.downloadView.cursor]
	switch item.Status {
	case downloads.StatusFailed, downloads.StatusCancelled:
		m.dlmgr.Retry(item.ID)
	case downloads.StatusQueued, downloads.StatusPreparing, downloads.StatusDownloading:
		m.dlmgr.Cancel(item.ID)
	case downloads.StatusCompleted:
		if p := item.DestPath(); p != "" {
			platform.OpenFileInDir(p)
		}
	}
	return nil
}

func (m *model) cacheMouseClick(y int) tea.Cmd {
	if y == 6 {
		m.cacheView.askClear = true
		m.cacheView.clearMsg = "clear all cached data? (enter=yes, any key=no)"
		return nil
	}
	idx := y - 9
	if idx >= 0 && idx < len(m.cacheView.recent) {
		m.cacheView.cursor = idx
		id := m.cacheView.recent[idx].EntityID
		return func() tea.Msg {
			return openProjectMsg{id: id}
		}
	}
	return nil
}

func (m *model) settingsMouseClick(y int) tea.Cmd {
	if y == 6 {
		m.settingsView.cursor = 3
		m.settingsView.askReset = true
		m.settingsView.resetMsg = "reset all settings to defaults? (enter=yes, any key=no)"
		return nil
	}
	idx := y - 2
	if idx >= 0 && idx < 3 {
		m.settingsView.cursor = idx
		model, cmd := m.settingsView.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m.settingsView = model
		return cmd
	}
	return nil
}

func (m *model) handlePickerMouse(msg tea.MouseMsg) tea.Cmd {
	if msg.Button == tea.MouseButtonWheelUp {
		model, cmd := m.picker.Update(tea.KeyMsg{Type: tea.KeyUp})
		m.picker = model
		return cmd
	}
	if msg.Button == tea.MouseButtonWheelDown {
		model, cmd := m.picker.Update(tea.KeyMsg{Type: tea.KeyDown})
		m.picker = model
		return cmd
	}
	if msg.Button != tea.MouseButtonLeft {
		return nil
	}
	_, y, ok := m.contentPoint(msg)
	if !ok {
		return nil
	}
	if y >= contentHeight(m.height)-1 {
		model, cmd := m.picker.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m.picker = model
		if m.picker.allVersions == nil {
			m.pickerActive = false
			m.setPage(navDownloads)
		}
		return cmd
	}
	if m.picker.selected {
		idx := y - 13
		if m.picker.cursor >= 0 && m.picker.cursor < len(m.picker.filtered) && idx >= 0 && idx < len(m.picker.filtered[m.picker.cursor].Files) {
			m.picker.fileCursor = idx
		}
		return nil
	}
	if m.picker.showFilters && y >= 3 && y <= 5 {
		m.picker.filterCursor = y - 3
		model, cmd := m.picker.Update(tea.KeyMsg{Type: tea.KeyRight})
		m.picker = model
		return cmd
	}
	idx := y - 5
	if idx >= 0 && idx < len(m.picker.filtered) {
		if m.picker.cursor == idx {
			model, cmd := m.picker.Update(tea.KeyMsg{Type: tea.KeyEnter})
			m.picker = model
			return cmd
		}
		m.picker.cursor = idx
	}
	return nil
}

func (m *model) handleThemeSwitcherMouse(msg tea.MouseMsg) tea.Cmd {
	if msg.Button == tea.MouseButtonWheelUp && m.themeSwitcher.cursor > 0 {
		m.themeSwitcher.cursor--
		return nil
	}
	if msg.Button == tea.MouseButtonWheelDown && m.themeSwitcher.cursor < len(m.themeSwitcher.themes)-1 {
		m.themeSwitcher.cursor++
		return nil
	}
	if msg.Button != tea.MouseButtonLeft {
		return nil
	}
	idx := msg.Y - (m.height-len(m.themeSwitcher.themes))/2 - 2
	if idx >= 0 && idx < len(m.themeSwitcher.themes) {
		m.themeSwitcher.cursor = idx
		model, cmd := m.themeSwitcher.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m.themeSwitcher = model
		return cmd
	}
	return nil
}

func (m *model) startSelectedProjectDownload(quick bool) tea.Cmd {
	id, title := m.selectedProject()
	if id == "" {
		return nil
	}
	m.pickerActive = true
	return m.picker.loadVersions(id, title, quick)
}

func (m *model) downloadSelectedVersionFile() tea.Cmd {
	if m.version.page != versionDetail || m.version.detail == nil || m.version.fileCursor < 0 || m.version.fileCursor >= len(m.version.detail.Files) {
		return nil
	}
	f := m.version.detail.Files[m.version.fileCursor]
	var hash *downloads.HashInfo
	if h, ok := f.Hashes["sha1"]; ok {
		hash = &downloads.HashInfo{Type: downloads.HashSHA1, Value: h}
	} else if h, ok := f.Hashes["sha512"]; ok {
		hash = &downloads.HashInfo{Type: downloads.HashSHA512, Value: h}
	}
	m.dlmgr.Enqueue(m.version.projectID, "", m.version.detail.ID, m.version.detail.VersionNumber, f.URL, f.Filename, f.Size, hash)
	return nil
}

func openSelectedVersionFile(v versionModel) {
	if v.page == versionDetail && v.detail != nil && len(v.detail.Files) > 0 && v.fileCursor >= 0 && v.fileCursor < len(v.detail.Files) {
		openURL(v.detail.Files[v.fileCursor].URL)
	}
}

func clamp(v, min, max int) int {
	if max < min {
		return min
	}
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
