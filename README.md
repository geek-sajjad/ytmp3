# ytmp3

Download the audio of a YouTube / YouTube Music video as an MP3.

```sh
ytmp3 "https://music.youtube.com/watch?v=VIDEO_ID"
```

Files are named after the video title and tagged with title/artist metadata.
A link from a playlist (`&list=...`) downloads only that one video.

## Install

Release downloads come with everything included: ffmpeg, yt-dlp and deno
(the JavaScript engine yt-dlp needs for YouTube) are built into the one
file. You don't need to install anything else.

### One-line install (recommended)

**macOS / Linux**: open Terminal and run:

```sh
curl -fsSL https://raw.githubusercontent.com/geek-sajjad/ytmp3/main/install.sh | sh
```

**Windows**: open PowerShell (Start → type "PowerShell") and run:

```powershell
irm https://raw.githubusercontent.com/geek-sajjad/ytmp3/main/install.ps1 | iex
```

The installer picks the right download for your computer, checks its
checksum, and adds `ytmp3` to your PATH. Open a **new** terminal afterwards
and check it works:

```sh
ytmp3 -version
```

The first time you download a song, ytmp3 spends a few seconds unpacking its
bundled tools. It only does this once.

### Manual download

1. Go to the [latest release](https://github.com/geek-sajjad/ytmp3/releases/latest)
   and download the file for your computer:

   | Computer                                | File                       |
   |-----------------------------------------|----------------------------|
   | Windows (most PCs)                      | `ytmp3_*_windows_amd64.zip`  |
   | Windows on ARM (e.g. Snapdragon)        | `ytmp3_*_windows_arm64.zip`  |
   | Mac with Apple Silicon (M1 and newer)   | `ytmp3_*_darwin_arm64.tar.gz` |
   | Mac with Intel                          | `ytmp3_*_darwin_amd64.tar.gz` |
   | Linux (most PCs)                        | `ytmp3_*_linux_amd64.tar.gz`  |
   | Linux on ARM (e.g. Raspberry Pi 4/5)    | `ytmp3_*_linux_arm64.tar.gz`  |

2. Extract it (double-click on Windows and macOS, `tar -xzf` on Linux).
3. Run it from a terminal in that folder: `./ytmp3` (macOS/Linux) or
   `.\ytmp3.exe` (Windows). To run it from anywhere, move it into a folder on
   your PATH.

**Security warnings**: the binaries aren't code-signed.
- **Windows** SmartScreen may say "Windows protected your PC". Click
  **More info → Run anyway**.
- **macOS** may say the file "can't be opened" when you downloaded it in a
  browser. Clear the quarantine flag once: `xattr -d com.apple.quarantine ./ytmp3`.
  The one-line installer doesn't have this problem.

### Uninstall

Delete the `ytmp3` binary (`~/.local/bin/ytmp3`, or
`%LOCALAPPDATA%\Programs\ytmp3` on Windows) and its unpacked tools:

- Linux: `~/.cache/ytmp3`
- macOS: `~/Library/Caches/ytmp3`
- Windows: `%LOCALAPPDATA%\ytmp3`

### From source

With Go 1.26+ installed:

```sh
go install github.com/geek-sajjad/ytmp3@latest
```

A plain source build doesn't include the bundled tools, so it uses
`ffmpeg`, `yt-dlp` and `deno` (or `node`) from your PATH. Install them
yourself first, e.g. `brew install ffmpeg yt-dlp deno` or
`winget install Gyan.FFmpeg`, `winget install yt-dlp.yt-dlp` and
`winget install DenoLand.Deno`. To build the
self-contained version instead:

```sh
scripts/fetch-deps.sh linux_amd64     # your <os>_<arch>
go build -tags bundled -o ytmp3 .
```

## Usage

Interactive: paste links one at a time, empty line to quit:

```sh
ytmp3
```

Or pass one or more links directly:

```sh
ytmp3 -o ~/Music -b 320k "https://music.youtube.com/watch?v=VIDEO_ID"
```

| Flag       | Default | Description              |
|------------|---------|--------------------------|
| `-o`       | `.`     | Output directory         |
| `-b`       | `192k`  | MP3 bitrate              |
| `-update`  |         | Update yt-dlp and exit   |
| `-version` |         | Print version and exit   |

## Troubleshooting

- **Downloads suddenly fail** (`HTTP 403`, "Sign in to confirm", "unable to
  extract"): YouTube changed something and yt-dlp needs updating. Run
  `ytmp3 -update`, or install the latest ytmp3 release.
- **`ytmp3: command not found`**: open a new terminal so the PATH change
  takes effect.

## Releasing

Releases are built by [GoReleaser](https://goreleaser.com) in GitHub Actions
(`.github/workflows/release.yml`). To publish a new version, push a tag:

```sh
git tag v1.0.0
git push origin v1.0.0
```

The workflow runs `scripts/fetch-deps.sh` to download ffmpeg, yt-dlp (the
latest version) and deno, builds Windows, macOS and Linux binaries (amd64 +
arm64) with them embedded, and attaches the archives and checksums to a
GitHub Release. Releasing regularly keeps the bundled yt-dlp fresh.

## Bundled software

Release binaries include these third-party programs, unmodified, under
their own licenses:

- [FFmpeg](https://ffmpeg.org) 6.1.1, static builds from
  [eugeneware/ffmpeg-static](https://github.com/eugeneware/ffmpeg-static/releases/tag/b6.1.1)
  (GPL v3; source at [ffmpeg.org](https://ffmpeg.org/download.html#releases))
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) (Unlicense)
- [Deno](https://github.com/denoland/deno) (MIT)

Pinned versions are in `scripts/fetch-deps.sh`.
