# ytmp3

Download the audio of a YouTube / YouTube Music video as an MP3.

```sh
ytmp3 "https://music.youtube.com/watch?v=VIDEO_ID"
```

Files are named after the video title and tagged with title/artist metadata.
A link from a playlist (`&list=...`) downloads only that one video.

## Install

ytmp3 is a single binary, but it needs two other tools at runtime:

- **ffmpeg** (required): converts the audio to MP3.
- **yt-dlp** (strongly recommended): fetches the audio. Without it ytmp3 falls
  back to a built-in Go client, which YouTube blocks (HTTP 403) for most music
  videos. Keep it current with `yt-dlp -U`: YouTube changes often break old
  versions.

Pick your platform below. Every binary is on the
[Releases page](https://github.com/geek-sajjad/ytmp3/releases/latest).

### Windows

1. Install the dependencies with [winget](https://learn.microsoft.com/windows/package-manager/winget/)
   (built into Windows 10/11), in PowerShell:

   ```powershell
   winget install Gyan.FFmpeg
   winget install yt-dlp.yt-dlp
   ```

2. Download `ytmp3_<version>_windows_amd64.zip` from the
   [Releases page](https://github.com/geek-sajjad/ytmp3/releases/latest)
   (`windows_arm64` for ARM devices such as Snapdragon laptops).
3. Extract `ytmp3.exe` to a folder, e.g. `C:\Tools\ytmp3`.
4. Add that folder to your PATH: Start → "Edit environment variables for your
   account" → **Path** → **Edit** → **New** → paste the folder path.
5. Open a **new** terminal and check it works:

   ```powershell
   ytmp3 -version
   ```

Windows SmartScreen may warn that the file is from an unknown publisher,
because the binary isn't code-signed. Click **More info → Run anyway**.

### macOS

1. Install the dependencies with [Homebrew](https://brew.sh):

   ```sh
   brew install ffmpeg yt-dlp
   ```

2. Download and install the binary. Use `darwin_arm64` for Apple Silicon
   (M1 and later) or `darwin_amd64` for Intel Macs:

   ```sh
   VERSION=1.0.0   # latest version from the Releases page, without the "v"
   ARCH=arm64      # or amd64
   curl -LO "https://github.com/geek-sajjad/ytmp3/releases/download/v${VERSION}/ytmp3_${VERSION}_darwin_${ARCH}.tar.gz"
   tar -xzf "ytmp3_${VERSION}_darwin_${ARCH}.tar.gz" ytmp3
   sudo mv ytmp3 /usr/local/bin/
   ```

3. The binary isn't signed by Apple, so Gatekeeper blocks it if you
   downloaded it in a browser. Remove the quarantine flag once:

   ```sh
   xattr -d com.apple.quarantine /usr/local/bin/ytmp3
   ```

4. Check it works: `ytmp3 -version`

### Linux

1. Install ffmpeg from your package manager:

   ```sh
   sudo apt install ffmpeg      # Debian / Ubuntu
   sudo dnf install ffmpeg      # Fedora (needs RPM Fusion)
   sudo pacman -S ffmpeg        # Arch
   ```

2. Install yt-dlp. Distro packages are usually outdated, so use the
   official release:

   ```sh
   mkdir -p ~/.local/bin
   curl -L -o ~/.local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux
   chmod +x ~/.local/bin/yt-dlp
   ```

3. Download and install ytmp3 (`linux_amd64` or `linux_arm64`, e.g. for a
   Raspberry Pi):

   ```sh
   VERSION=1.0.0   # latest version from the Releases page, without the "v"
   ARCH=amd64      # or arm64
   curl -LO "https://github.com/geek-sajjad/ytmp3/releases/download/v${VERSION}/ytmp3_${VERSION}_linux_${ARCH}.tar.gz"
   tar -xzf "ytmp3_${VERSION}_linux_${ARCH}.tar.gz" ytmp3
   mv ytmp3 ~/.local/bin/
   ```

   Make sure `~/.local/bin` is in your PATH.

4. Check it works: `ytmp3 -version`

### From source (any OS)

With Go 1.26+ installed (an older Go fetches the right toolchain
automatically):

```sh
go install github.com/geek-sajjad/ytmp3@latest
```

Or clone and build:

```sh
git clone https://github.com/geek-sajjad/ytmp3.git
cd ytmp3
go build -o ytmp3 .
```

You still need ffmpeg and yt-dlp as described above.

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
| `-version` |         | Print version and exit   |

## Troubleshooting

- **`HTTP 403` / "YouTube blocked the download"**: install yt-dlp, or update
  it with `yt-dlp -U`.
- **`ffmpeg not found in PATH`**: install ffmpeg and open a new terminal.
- yt-dlp needs a JavaScript runtime for some videos. ytmp3 uses `deno` if
  installed, otherwise `node`.

## Releasing

Releases are built by [GoReleaser](https://goreleaser.com) in GitHub Actions
(`.github/workflows/release.yml`). To publish a new version, push a tag:

```sh
git tag v1.0.0
git push origin v1.0.0
```

The workflow builds Windows, macOS and Linux binaries (amd64 + arm64),
packages them with checksums and attaches them to a GitHub Release.
