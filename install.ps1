# Installs the latest ytmp3 release on Windows.
#
#   irm https://raw.githubusercontent.com/geek-sajjad/ytmp3/main/install.ps1 | iex
#
# Installs to %LOCALAPPDATA%\Programs\ytmp3 and adds it to your user PATH.
$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'  # Invoke-WebRequest is very slow with the progress bar on

$repo = 'geek-sajjad/ytmp3'
$installDir = Join-Path $env:LOCALAPPDATA 'Programs\ytmp3'

$arch = if ($env:PROCESSOR_ARCHITECTURE -eq 'ARM64') { 'arm64' } else { 'amd64' }

$release = Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest"
$tag = $release.tag_name
$version = $tag.TrimStart('v')
$archive = "ytmp3_${version}_windows_${arch}.zip"
$base = "https://github.com/$repo/releases/download/$tag"

$tmp = Join-Path ([IO.Path]::GetTempPath()) ([Guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null
try {
    Write-Host "Downloading ytmp3 $tag for windows/$arch (about 100 MB)..."
    $zip = Join-Path $tmp $archive
    Invoke-WebRequest "$base/$archive" -OutFile $zip
    $sums = (Invoke-WebRequest "$base/checksums.txt" -UseBasicParsing).Content

    $line = ($sums -split "`n") | Where-Object { $_ -match " $([regex]::Escape($archive))$" }
    $expected = ($line -split '\s+')[0]
    $actual = (Get-FileHash $zip -Algorithm SHA256).Hash
    if ($actual -ne $expected) { throw "Checksum mismatch for $archive" }

    Expand-Archive $zip -DestinationPath $tmp -Force
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Copy-Item (Join-Path $tmp 'ytmp3.exe') $installDir -Force
} finally {
    Remove-Item $tmp -Recurse -Force -ErrorAction SilentlyContinue
}
Write-Host "Installed ytmp3 to $installDir\ytmp3.exe"

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (-not (($userPath -split ';') -contains $installDir)) {
    $newPath = if ($userPath) { "$userPath;$installDir" } else { $installDir }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    $env:Path = "$env:Path;$installDir"
    Write-Host "Added $installDir to your PATH. Open a new terminal to use ytmp3 everywhere."
}

Write-Host ''
Write-Host 'Done! Try: ytmp3 "https://www.youtube.com/watch?v=VIDEO_ID"'
