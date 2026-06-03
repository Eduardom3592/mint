package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/programmersd21/mint/internal/api"
	"github.com/programmersd21/mint/internal/cache"
	"github.com/programmersd21/mint/internal/downloads"
	"github.com/programmersd21/mint/internal/tui"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildDate = "unknown"
)

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	dataDir := flag.String("data-dir", defaultDataDir(), "Data directory for cache and settings")
	downloadDir := flag.String("download-dir", defaultDownloadDir(), "Download directory")
	apiKey := flag.String("api-key", "", "Modrinth API key")
	flag.Parse()

	if *showVersion {
		fmt.Printf("mint %s (commit: %s, built: %s)\n", version, commit, buildDate)
		return
	}

	c, err := cache.Open(cache.Config{DataDir: *dataDir})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening cache: %s\n", err)
		os.Exit(1)
	}
	defer c.Close()

	client := api.NewClient()
	if *apiKey != "" {
		client = api.NewClient(api.WithAPIKey(*apiKey))
	}

	dlDir := *downloadDir
	if dlDir == "" {
		dlDir = *dataDir + "/downloads"
	}
	dlmgr := downloads.NewManager(dlDir, 3, c)
	defer dlmgr.Close()

	p := tea.NewProgram(
		tui.New(client, c, dlmgr),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		p.Quit()
	}()

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func defaultDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".mint"
	}
	return home + "/.local/share/mint"
}

func defaultDownloadDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "downloads"
	}
	return home + "/Downloads/mint"
}
