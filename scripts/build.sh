#!/usr/bin/env bash
set -e

oses=(windows linux darwin)
archs=(amd64 386 arm64 arm)

mkdir -p bin

for os in "${oses[@]}"; do
  for arch in "${archs[@]}"; do
    # Skip unsupported OS/arch combinations
    if { [ "$os" = "darwin" ] && { [ "$arch" = "386" ] || [ "$arch" = "arm" ]; }; } || \
       { [ "$os" = "windows" ] && [ "$arch" = "arm" ]; }; then
      continue
    fi
    ext=""
    if [ "$os" = "windows" ]; then
      ext=".exe"
    fi
    outfile="bin/yt_dlp_web_${os}_${arch}${ext}"
    echo "Building for $os/$arch -> $outfile"
    GOOS="$os" GOARCH="$arch" go build -o "$outfile" .
  done
done
