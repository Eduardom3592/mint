package platform

import (
	"os/exec"
	"runtime"
)

func OpenFolder(dir string) {
	cmd := openFolderCommand(dir)
	if cmd != nil {
		_ = cmd.Start()
	}
}

func OpenFileInDir(p string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", "/select,", p)
	case "darwin":
		cmd = exec.Command("open", "-R", p)
	default:
		cmd = exec.Command("xdg-open", p)
	}
	if cmd != nil {
		_ = cmd.Start()
	}
}

func OpenURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if cmd != nil {
		_ = cmd.Start()
	}
}

func openFolderCommand(dir string) *exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("explorer", dir)
	case "darwin":
		return exec.Command("open", dir)
	default:
		return exec.Command("xdg-open", dir)
	}
}
