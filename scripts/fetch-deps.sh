#!/bin/sh
# Downloads the ffmpeg, yt-dlp and deno binaries that release builds embed
# (go build -tags bundled) into deps/<goos>_<goarch>/, gzip-compressed.
#
# Usage: scripts/fetch-deps.sh [goos_goarch ...]   (default: all release targets)
set -eu

FFMPEG_VERSION=b6.1.1
DENO_VERSION=v2.9.7
YTDLP_VERSION=${YTDLP_VERSION:-latest}

cd "$(dirname "$0")/.."
cache=deps/.download
mkdir -p "$cache"

# fetch URL -> path in the download cache (reused across targets)
fetch() {
	out="$cache/$(echo "$1" | tr '/:' '__')"
	if [ ! -s "$out" ]; then
		echo "  downloading $1" >&2
		curl -fsSL --retry 3 -o "$out.tmp" "$1"
		mv "$out.tmp" "$out"
	fi
	echo "$out"
}

ytdlp_url() {
	if [ "$YTDLP_VERSION" = latest ]; then
		echo "https://github.com/yt-dlp/yt-dlp/releases/latest/download/$1"
	else
		echo "https://github.com/yt-dlp/yt-dlp/releases/download/$YTDLP_VERSION/$1"
	fi
}

for target in ${*:-linux_amd64 linux_arm64 darwin_amd64 darwin_arm64 windows_amd64 windows_arm64}; do
	exe=
	case $target in
	linux_amd64)   ff=linux-x64;    yt=yt-dlp_linux;         dn=x86_64-unknown-linux-gnu ;;
	linux_arm64)   ff=linux-arm64;  yt=yt-dlp_linux_aarch64; dn=aarch64-unknown-linux-gnu ;;
	darwin_amd64)  ff=darwin-x64;   yt=yt-dlp_macos;         dn=x86_64-apple-darwin ;;
	darwin_arm64)  ff=darwin-arm64; yt=yt-dlp_macos;         dn=aarch64-apple-darwin ;;
	windows_amd64) ff=win32-x64;    yt=yt-dlp.exe;           dn=x86_64-pc-windows-msvc; exe=.exe ;;
	# No native ffmpeg build for Windows on ARM; Windows runs the x64 one emulated.
	windows_arm64) ff=win32-x64;    yt=yt-dlp_arm64.exe;     dn=aarch64-pc-windows-msvc; exe=.exe ;;
	*) echo "unknown target: $target" >&2; exit 1 ;;
	esac

	echo "$target"
	dir=deps/$target
	rm -rf "$dir"
	mkdir -p "$dir"
	cp "$(fetch "https://github.com/eugeneware/ffmpeg-static/releases/download/$FFMPEG_VERSION/ffmpeg-$ff.gz")" "$dir/ffmpeg$exe.gz"
	gzip -9c "$(fetch "$(ytdlp_url "$yt")")" > "$dir/yt-dlp$exe.gz"
	unzip -p "$(fetch "https://github.com/denoland/deno/releases/download/$DENO_VERSION/deno-$dn.zip")" "deno$exe" | gzip -9c > "$dir/deno$exe.gz"
done
