package main

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// setupBundledTools extracts the ffmpeg, yt-dlp and deno binaries embedded in
// release builds into the user cache directory and puts that directory first
// on PATH, so they are found by exec.LookPath and by yt-dlp itself.
//
// A tool is only rewritten when the embedded copy changes, so a yt-dlp that
// updated itself with -update keeps its newer version.
func setupBundledTools() error {
	if bundledDir == "" {
		return nil
	}
	cache, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(cache, "ytmp3", "bin")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	entries, err := fs.ReadDir(bundledFS, bundledDir)
	if err != nil {
		return err
	}
	announced := false
	for _, e := range entries {
		data, err := fs.ReadFile(bundledFS, path.Join(bundledDir, e.Name()))
		if err != nil {
			return err
		}
		sum := sha256.Sum256(data)
		stamp := hex.EncodeToString(sum[:])
		dst := filepath.Join(dir, strings.TrimSuffix(e.Name(), ".gz"))
		if old, err := os.ReadFile(dst + ".sha256"); err == nil && string(old) == stamp {
			if _, err := os.Stat(dst); err == nil {
				continue
			}
		}
		if !announced {
			fmt.Println("Setting up bundled ffmpeg and yt-dlp (first run only)...")
			announced = true
		}
		if err := extractGzip(data, dst); err != nil {
			return fmt.Errorf("extracting %s: %w", filepath.Base(dst), err)
		}
		if err := os.WriteFile(dst+".sha256", []byte(stamp), 0o644); err != nil {
			return err
		}
	}
	return os.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

func extractGzip(data []byte, dst string) error {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	tmp := dst + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, zr); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	return os.Rename(tmp, dst)
}
