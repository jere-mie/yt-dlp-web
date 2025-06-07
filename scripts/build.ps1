# Define the target operating systems and architectures
$oses = @('windows', 'linux', 'darwin')
$archs = @('amd64', '386', 'arm64', 'arm')

# Create the output directory if it doesn't exist
if (!(Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

# Loop through each combination of OS and architecture
foreach ($os in $oses) {
    foreach ($arch in $archs) {
        # Skip unsupported OS/arch combinations
        if ($os -eq 'darwin' -and ($arch -eq '386' -or $arch -eq 'arm')) {
            continue
        }
        # Set the file extension for Windows executables
        $ext = if ($os -eq 'windows') { '.exe' } else { '' }
        # Define the output file path
        $outfile = "bin/yt_dlp_web_${os}_${arch}$ext"
        Write-Host "Building for $os/$arch -> $outfile"
        # Set the environment variables for cross-compilation
        $env:GOOS = $os
        $env:GOARCH = $arch
        # Run the Go build command
        go build -o $outfile .
    }
}