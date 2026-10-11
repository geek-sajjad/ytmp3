#!/bin/sh
# Installs the latest ytmp3 release on macOS or Linux.
#
#   curl -fsSL https://raw.githubusercontent.com/geek-sajjad/ytmp3/main/install.sh | sh
#
# Set YTMP3_INSTALL_DIR to install somewhere other than ~/.local/bin.
set -eu

repo=geek-sajjad/ytmp3
bindir=${YTMP3_INSTALL_DIR:-$HOME/.local/bin}

case $(uname -s) in
Linux) os=linux ;;
Darwin) os=darwin ;;
*) echo "Unsupported OS: $(uname -s). On Windows, use install.ps1." >&2; exit 1 ;;
esac
case $(uname -m) in
x86_64 | amd64) arch=amd64 ;;
aarch64 | arm64) arch=arm64 ;;
*) echo "Unsupported CPU: $(uname -m)" >&2; exit 1 ;;
esac

# The releases/latest page redirects to .../tag/vX.Y.Z.
tag=$(curl -fsSLI -o /dev/null -w '%{url_effective}' "https://github.com/$repo/releases/latest")
tag=${tag##*/}
version=${tag#v}
archive="ytmp3_${version}_${os}_${arch}.tar.gz"
base="https://github.com/$repo/releases/download/$tag"

tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT

echo "Downloading ytmp3 $tag for $os/$arch..."
curl -fL --progress-bar -o "$tmp/$archive" "$base/$archive"
curl -fsSL -o "$tmp/checksums.txt" "$base/checksums.txt"

cd "$tmp"
grep " $archive\$" checksums.txt > expected.txt
if command -v sha256sum >/dev/null 2>&1; then
	sha256sum -c --quiet expected.txt
else
	shasum -a 256 -c --quiet expected.txt
fi

tar -xzf "$archive" ytmp3
mkdir -p "$bindir"
mv ytmp3 "$bindir/ytmp3"
chmod +x "$bindir/ytmp3"
echo "Installed ytmp3 to $bindir/ytmp3"

case ":$PATH:" in
*":$bindir:"*) ;;
*)
	case ${SHELL:-} in
	*/zsh) rc=$HOME/.zshrc ;;
	*/bash) if [ "$os" = darwin ]; then rc=$HOME/.bash_profile; else rc=$HOME/.bashrc; fi ;;
	*) rc=$HOME/.profile ;;
	esac
	echo "export PATH=\"$bindir:\$PATH\"" >> "$rc"
	echo "Added $bindir to your PATH in $rc."
	echo "Open a new terminal (or run: export PATH=\"$bindir:\$PATH\") to use ytmp3."
	;;
esac

echo
echo "Done! Try: ytmp3 \"https://www.youtube.com/watch?v=VIDEO_ID\""
