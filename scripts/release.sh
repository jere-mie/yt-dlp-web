#!/usr/bin/env bash

# Requires: GitHub CLI (gh), version.txt, and bin/ with built binaries

set -e

# Get the version number from version.txt
version=$(head -n 1 ./version.txt)
tag="v$version"

# Check if the release already exists
if gh release view "$tag" >/dev/null 2>&1; then
    echo "Release $tag already exists. Exiting."
    exit 1
fi

# Find all binaries in the bin directory
bin_files=(./bin/*)
if [ ${#bin_files[@]} -eq 0 ] || [ ! -e "${bin_files[0]}" ]; then
    echo "No binaries found in bin/. Run build.sh or build.ps1 first."
    exit 1
fi

# Create the release and upload binaries
echo "Creating GitHub release $tag..."
gh release create "$tag" "${bin_files[@]}" \
    --title "yt-dlp-web $version" \
    --notes "Automated release for version $version."

echo "Release $tag created and binaries uploaded."
