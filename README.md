# ytmp3

Download the audio of a YouTube / YouTube Music video as an MP3.

## Requirements

- Go 1.26+ (fetched automatically by the Go toolchain if you have an older Go)
- `ffmpeg` in your PATH
- `yt-dlp` in your PATH (strongly recommended). Without it, a built-in Go
  client is used, but YouTube blocks it (HTTP 403) for most music videos.
  The apt package is usually outdated; install the latest release instead:

  ```sh
  curl -L -o ~/.local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux
  chmod +x ~/.local/bin/yt-dlp
  ```

  `yt-dlp -U` updates it when YouTube changes break downloads.

A link from a playlist (`&list=...`) downloads only that one video.

## Build

```sh
go build -o ytmp3 .
```

## Usage

Interactive: paste links one at a time, empty line to quit:

```sh
./ytmp3
```

Or pass one or more links directly:

```sh
./ytmp3 -o ~/Music -b 320k "https://music.youtube.com/watch?v=VIDEO_ID"
```

| Flag | Default | Description              |
|------|---------|--------------------------|
| `-o` | `.`     | Output directory         |
| `-b` | `192k`  | MP3 bitrate              |

Files are named after the video title and tagged with title/artist metadata.
