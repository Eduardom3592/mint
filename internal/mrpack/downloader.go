package mrpack

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const downloadUserAgent = "Mint/0.1.0 (terminal-modrinth-client)"

var downloadClient = &http.Client{Timeout: 30 * time.Second}

func DownloadModpack(url, destDir string) (string, error) {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("create modpacks dir: %w", err)
	}

	name := fmt.Sprintf("modpack_%d.mrpack", time.Now().UnixMilli())
	dest := filepath.Join(destDir, name)

	out, err := os.Create(dest)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	defer out.Close()

	resp, err := downloadURL(url)
	if err != nil {
		_ = os.Remove(dest)
		return "", fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(dest)
		return "", fmt.Errorf("http %d", resp.StatusCode)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		_ = os.Remove(dest)
		return "", fmt.Errorf("write: %w", err)
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(dest)
		return "", fmt.Errorf("close file: %w", err)
	}

	return dest, nil
}

func DownloadPackFile(url, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := downloadURL(url)
	if err != nil {
		_ = os.Remove(destPath)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(destPath)
		return fmt.Errorf("http %d", resp.StatusCode)
	}

	if _, err := io.Copy(out, resp.Body); err != nil {
		_ = os.Remove(destPath)
		return err
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(destPath)
		return err
	}

	return nil
}

func downloadURL(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", downloadUserAgent)
	return downloadClient.Do(req)
}
