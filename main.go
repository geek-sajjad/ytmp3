// Command ytmp3 downloads the audio track of a YouTube (or YouTube Music)
// video and converts it to MP3 using ffmpeg.
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func main() {
	outDir := flag.String("o", ".", "directory to save MP3 files in")
	bitrate := flag.String("b", "192k", "MP3 bitrate (e.g. 128k, 192k, 320k)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [URL ...]\n\n", filepath.Base(os.Args[0]))
		fmt.Fprintln(os.Stderr, "Downloads YouTube / YouTube Music audio as MP3.")
		fmt.Fprintln(os.Stderr, "With no URLs, prompts for links interactively.\n\nFlags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		fatal("ffmpeg not found in PATH; install it first (e.g. sudo apt install ffmpeg)")
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fatal("cannot create output directory: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	d := &downloader{src: newSource(), outDir: *outDir, bitrate: *bitrate}
	fmt.Printf("Using %s to fetch audio.\n", d.src.Name())

	if flag.NArg() > 0 {
		failed := 0
		for _, url := range flag.Args() {
			if !d.run(ctx, url) {
				failed++
			}
		}
		if failed > 0 {
			os.Exit(1)
		}
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("\nEnter YouTube link (empty to quit): ")
		if !scanner.Scan() {
			fmt.Println()
			return
		}
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			return
		}
		d.run(ctx, url)
		if ctx.Err() != nil {
			return
		}
	}
}

type downloader struct {
	src     source
	outDir  string
	bitrate string
}

// run downloads one URL and reports the outcome; it returns false on failure.
func (d *downloader) run(ctx context.Context, url string) bool {
	path, err := d.download(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
		return false
	}
	fmt.Printf("\nSaved: %s\n", path)
	return true
}

func (d *downloader) download(ctx context.Context, url string) (string, error) {
	fmt.Println("Fetching video info...")
	t, err := d.src.Open(ctx, url)
	if err != nil {
		return "", err
	}
	defer t.Body.Close()
	fmt.Printf("Title:  %s\nArtist: %s\nLength: %s\n", t.Title, t.Artist, t.Length)

	finalPath := filepath.Join(d.outDir, sanitizeFilename(t.Title)+".mp3")
	tmpPath := finalPath + ".part"

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-hide_banner", "-loglevel", "error", "-y",
		"-i", "pipe:0",
		"-vn",
		"-codec:a", "libmp3lame",
		"-b:a", d.bitrate,
		"-metadata", "title="+t.Title,
		"-metadata", "artist="+t.Artist,
		"-f", "mp3",
		tmpPath,
	)
	in := &progressReader{r: t.Body, total: t.Size}
	cmd.Stdin = in
	var stderr strings.Builder
	cmd.Stderr = &stderr

	err = cmd.Run()
	// A failed download is the root cause of any ffmpeg error that follows it.
	if in.err != nil && in.err != io.EOF {
		err = in.err
	} else if err != nil {
		err = fmt.Errorf("ffmpeg: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	if err != nil {
		os.Remove(tmpPath)
		if ctx.Err() != nil {
			return "", errors.New("cancelled")
		}
		return "", err
	}
	if err := os.Rename(tmpPath, finalPath); err != nil {
		return "", err
	}
	return finalPath, nil
}

// sanitizeFilename strips characters that are invalid in file names.
func sanitizeFilename(name string) string {
	name = strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '_'
		}
		if r < 32 {
			return -1
		}
		return r
	}, name)
	name = strings.Trim(strings.TrimSpace(name), ".")
	if name == "" {
		name = "audio"
	}
	if len(name) > 200 {
		name = name[:200]
	}
	return name
}

// progressReader prints download progress as data is read.
type progressReader struct {
	r     io.Reader
	total int64
	read  int64
	last  time.Time
	err   error // first read error, including io.EOF
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.read += int64(n)
	if err != nil && p.err == nil {
		p.err = err
	}
	if err != nil || time.Since(p.last) > 200*time.Millisecond {
		p.last = time.Now()
		if p.total > 0 {
			fmt.Printf("\rDownloading: %5.1f%% (%.1f / %.1f MB)",
				float64(p.read)*100/float64(p.total), mb(p.read), mb(p.total))
		} else {
			fmt.Printf("\rDownloading: %.1f MB", mb(p.read))
		}
	}
	return n, err
}

func mb(n int64) float64 { return float64(n) / (1 << 20) }

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
	os.Exit(1)
}
