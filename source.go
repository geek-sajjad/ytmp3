package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"

	"github.com/kkdai/youtube/v2"
)

// track is an opened audio stream plus the metadata needed to save it.
type track struct {
	Title  string
	Artist string
	Length string
	Size   int64 // 0 if unknown
	Body   io.ReadCloser
}

// source opens the best available audio stream for a video URL.
type source interface {
	Name() string
	Open(ctx context.Context, url string) (*track, error)
}

// newSource prefers yt-dlp, which keeps up with YouTube's download
// protections, and falls back to the pure-Go client when it isn't installed.
func newSource() source {
	if path, err := exec.LookPath("yt-dlp"); err == nil {
		return &ytdlpSource{bin: path}
	}
	return &goSource{client: &youtube.Client{}}
}

// ytdlpSource fetches audio with the yt-dlp command-line tool.
type ytdlpSource struct{ bin string }

func (s *ytdlpSource) Name() string { return "yt-dlp" }

func (s *ytdlpSource) args(extra ...string) []string {
	args := []string{"--no-playlist", "--no-warnings", "-f", "bestaudio/best"}
	// yt-dlp only enables deno by default; use node if that is what's installed.
	if _, err := exec.LookPath("deno"); err != nil {
		if _, err := exec.LookPath("node"); err == nil {
			args = append(args, "--js-runtimes", "node")
		}
	}
	return append(args, extra...)
}

func (s *ytdlpSource) Open(ctx context.Context, url string) (*track, error) {
	var stderr strings.Builder
	info := exec.CommandContext(ctx, s.bin, s.args("--dump-json", "--", url)...)
	info.Stderr = &stderr
	out, err := info.Output()
	if err != nil {
		return nil, fmt.Errorf("yt-dlp: %v: %s", err, strings.TrimSpace(stderr.String()))
	}
	var meta struct {
		Title          string  `json:"title"`
		Track          string  `json:"track"`
		Artist         string  `json:"artist"`
		Uploader       string  `json:"uploader"`
		Duration       float64 `json:"duration"`
		Filesize       int64   `json:"filesize"`
		FilesizeApprox float64 `json:"filesize_approx"`
	}
	if err := json.Unmarshal(out, &meta); err != nil {
		return nil, fmt.Errorf("parsing yt-dlp output: %w", err)
	}

	t := &track{
		Title:  firstNonEmpty(meta.Track, meta.Title),
		Artist: firstNonEmpty(meta.Artist, meta.Uploader),
		Length: fmt.Sprintf("%dm%02ds", int(meta.Duration)/60, int(meta.Duration)%60),
		Size:   meta.Filesize,
	}
	if t.Size == 0 {
		t.Size = int64(meta.FilesizeApprox)
	}

	stream := exec.CommandContext(ctx, s.bin, s.args("--quiet", "-o", "-", "--", url)...)
	var streamErr strings.Builder
	stream.Stderr = &streamErr
	stdout, err := stream.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := stream.Start(); err != nil {
		return nil, fmt.Errorf("starting yt-dlp: %w", err)
	}
	t.Body = &cmdReader{r: stdout, cmd: stream, stderr: &streamErr}
	return t, nil
}

// cmdReader reads a command's stdout and surfaces its exit error at EOF.
type cmdReader struct {
	r      io.Reader
	cmd    *exec.Cmd
	stderr *strings.Builder
	done   bool
}

func (c *cmdReader) Read(b []byte) (int, error) {
	n, err := c.r.Read(b)
	if err == io.EOF && !c.done {
		c.done = true
		if werr := c.cmd.Wait(); werr != nil {
			return n, fmt.Errorf("yt-dlp: %v: %s", werr, strings.TrimSpace(c.stderr.String()))
		}
	}
	return n, err
}

func (c *cmdReader) Close() error {
	if !c.done {
		c.done = true
		c.cmd.Process.Kill()
		c.cmd.Wait()
	}
	return nil
}

// goSource fetches audio with the pure-Go kkdai/youtube client. It needs no
// external tools but YouTube blocks it for many music videos (HTTP 403).
type goSource struct{ client *youtube.Client }

func (s *goSource) Name() string { return "built-in client" }

func (s *goSource) Open(ctx context.Context, url string) (*track, error) {
	video, err := s.client.GetVideoContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("fetching video info: %w", err)
	}
	format, err := bestAudioFormat(video.Formats)
	if err != nil {
		return nil, err
	}
	body, size, err := s.client.GetStreamContext(ctx, video, format)
	if err != nil {
		return nil, fmt.Errorf("opening audio stream: %w", err)
	}
	return &track{
		Title:  video.Title,
		Artist: video.Author,
		Length: video.Duration.String(),
		Size:   size,
		Body:   &hintReader{body},
	}, nil
}

// hintReader adds advice to the 403 error YouTube returns for protected videos.
type hintReader struct{ io.ReadCloser }

func (h *hintReader) Read(b []byte) (int, error) {
	n, err := h.ReadCloser.Read(b)
	if err != nil && strings.Contains(err.Error(), "403") {
		err = fmt.Errorf("%w (YouTube blocked the download; install yt-dlp and try again)", err)
	}
	return n, err
}

// bestAudioFormat picks the highest-bitrate audio-only stream, falling back
// to any format that carries audio.
func bestAudioFormat(formats youtube.FormatList) (*youtube.Format, error) {
	candidates := formats.Type("audio/").WithAudioChannels()
	if len(candidates) == 0 {
		candidates = formats.WithAudioChannels()
	}
	if len(candidates) == 0 {
		return nil, errors.New("no audio stream available for this video")
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Bitrate > candidates[j].Bitrate
	})
	return &candidates[0], nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
