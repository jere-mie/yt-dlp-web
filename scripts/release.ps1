# Requires: GitHub CLI (gh), PowerShell 7+, version.txt, and bin/ with built binaries

# Get the version number from version.txt
$version = Get-Content ./version.txt | Select-Object -First 1
$tag = "v$version"

# Check if the release already exists
gh release view $tag 2>$null

if ($LASTEXITCODE -eq 0) {
    Write-Host "Release $tag already exists. Exiting."
    exit 1
}

# Find all binaries in the bin directory
$binFiles = Get-ChildItem ./bin -File

if ($binFiles.Count -eq 0) {
    Write-Host "No binaries found in bin/. Run build.ps1 first."
    exit 1
}

# Create the release and upload binaries
Write-Host "Creating GitHub release $tag..."
gh release create $tag $($binFiles | ForEach-Object { $_.FullName }) `
    --title "yt-dlp-web $version" `
    --notes "Automated release for version $version."

Write-Host "Release $tag created and binaries uploaded."