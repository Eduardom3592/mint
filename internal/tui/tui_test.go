package tui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/cache"
	"github.com/programmersd21/mint/internal/models"
)

func TestPrevPageToSearchResetsSearchList(t *testing.T) {
	m := &model{
		currentPage: navDownloads,
		search: searchModel{
			page:              searchProject,
			selectedProjectID: "project-id",
		},
	}

	m.prevPage()

	if m.currentPage != navSearch {
		t.Fatalf("expected current page search, got %v", m.currentPage)
	}
	if m.search.page != searchList {
		t.Fatalf("expected search list page, got %v", m.search.page)
	}
	if m.search.selectedProjectID != "" {
		t.Fatalf("expected selected project to be cleared, got %q", m.search.selectedProjectID)
	}
}

func TestSetPageSearchResetsSearchList(t *testing.T) {
	m := &model{
		currentPage: navDownloads,
		search: searchModel{
			page:              searchProject,
			selectedProjectID: "project-id",
		},
	}

	m.setPage(navSearch)

	if m.search.page != searchList {
		t.Fatalf("expected search list page, got %v", m.search.page)
	}
	if m.search.selectedProjectID != "" {
		t.Fatalf("expected selected project to be cleared, got %q", m.search.selectedProjectID)
	}
}

func TestMouseSidebarNavigatesToSearch(t *testing.T) {
	m := &model{
		loaded:      true,
		width:       100,
		height:      40,
		currentPage: navDownloads,
		search: searchModel{
			page:              searchProject,
			selectedProjectID: "project-id",
		},
	}

	cmd := m.handleMouse(tea.MouseMsg{
		X:      0,
		Y:      3,
		Button: tea.MouseButtonLeft,
		Action: tea.MouseActionPress,
	})

	if cmd != nil {
		t.Fatal("expected no command")
	}
	if m.currentPage != navSearch {
		t.Fatalf("expected current page search, got %v", m.currentPage)
	}
	if m.search.page != searchList {
		t.Fatalf("expected search list page, got %v", m.search.page)
	}
}

func TestMouseWheelMovesHomeSelection(t *testing.T) {
	m := &model{
		loaded:      true,
		width:       100,
		height:      40,
		currentPage: navHome,
		home: homeModel{
			hits: make([]models.SearchHit, 3),
		},
	}

	m.handleMouse(tea.MouseMsg{
		X:      30,
		Y:      10,
		Button: tea.MouseButtonWheelDown,
		Action: tea.MouseActionPress,
	})

	if m.home.cursor != 1 {
		t.Fatalf("expected home cursor 1, got %d", m.home.cursor)
	}
}

func TestMouseClickCacheRecentItem(t *testing.T) {
	c, err := cache.Open(cache.Config{DataDir: t.TempDir()})
	if err != nil {
		t.Fatalf("open cache: %v", err)
	}
	defer c.Close()

	m := &model{
		cache:       c,
		loaded:      true,
		width:       100,
		height:      40,
		currentPage: navCache,
		project: projectModel{
			cache: c,
		},
		cacheView: cacheViewModel{
			loaded: true,
			recent: []recentItem{
				{EntityID: "abc", Title: "project a", ViewedAt: time.Now()},
			},
		},
	}

	cmd := m.cacheMouseClick(9)
	if cmd == nil {
		t.Fatal("expected command to open project")
	}
	msg := cmd()
	updated, _ := m.Update(msg)
	m = updated.(*model)
	if m.currentPage != navProject {
		t.Fatalf("expected project page, got %v", m.currentPage)
	}
}

func TestMouseClickSettingsRows(t *testing.T) {
	m := &model{
		loaded:      true,
		width:       100,
		height:      40,
		currentPage: navSettings,
		settingsView: settingsViewModel{
			cursor: 0,
			focus:  -1,
			inputs: []textinput.Model{textinput.New(), textinput.New(), textinput.New()},
		},
	}

	m.settingsMouseClick(4)
	if m.settingsView.cursor != 2 {
		t.Fatalf("expected settings cursor 2, got %d", m.settingsView.cursor)
	}

	m.settingsMouseClick(6)
	if !m.settingsView.askReset {
		t.Fatal("expected reset prompt to open")
	}
}
