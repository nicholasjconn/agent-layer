param(
    [string]$Version = "latest",
    [string]$InstallDir = ""
)

Write-Warning "Windows support is not tested and is not guaranteed to work."

if ([string]::IsNullOrWhiteSpace($InstallDir)) {
    if ($env:LOCALAPPDATA) {
        $InstallDir = Join-Path $env:LOCALAPPDATA "Programs\AgentLayer\bin"
    } elseif ($env:USERPROFILE) {
        $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
    } else {
        throw "Unable to determine a user install directory. Set -InstallDir explicitly."
    }
}

$baseUrl = "https://github.com/conn-castle/agent-layer/releases"
$asset = "al-windows-amd64.exe"

$tag = "latest"
if ($Version -ne "latest") {
    $normalized = $Version.TrimStart("v")
    if ($normalized -notmatch "^\d+\.\d+\.\d+$") {
        throw "Invalid version: $Version (expected vX.Y.Z or X.Y.Z)"
    }
    $tag = "v$normalized"
}

$downloadUrl = if ($tag -eq "latest") { "$baseUrl/latest/download/$asset" } else { "$baseUrl/download/$tag/$asset" }
$checksumsUrl = if ($tag -eq "latest") { "$baseUrl/latest/download/checksums.txt" } else { "$baseUrl/download/$tag/checksums.txt" }

$tempAsset = Join-Path $env:TEMP ([System.IO.Path]::GetRandomFileName() + ".exe")
$tempChecksums = Join-Path $env:TEMP ([System.IO.Path]::GetRandomFileName() + ".txt")

try {
    Invoke-WebRequest -Uri $downloadUrl -OutFile $tempAsset
    Invoke-WebRequest -Uri $checksumsUrl -OutFile $tempChecksums

    $expected = $null
    foreach ($line in Get-Content $tempChecksums) {
        $trimmed = $line.Trim()
        if ([string]::IsNullOrWhiteSpace($trimmed)) {
            continue
        }
        $parts = $trimmed -split "\s+"
        if ($parts.Length -lt 2) {
            continue
        }
        $path = $parts[1] -replace "^\*", ""
        $path = $path -replace "^\./", ""
        if ($path -eq $asset) {
            $expected = $parts[0]
            break
        }
    }

    if (-not $expected) {
        throw "Checksum for $asset not found in $checksumsUrl."
    }

    $actual = (Get-FileHash -Algorithm SHA256 -Path $tempAsset).Hash.ToLowerInvariant()
    if ($actual -ne $expected.ToLowerInvariant()) {
        throw "Checksum mismatch for $asset."
    }

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $dest = Join-Path $InstallDir "al.exe"
    Move-Item -Force -Path $tempAsset -Destination $dest

    Write-Host "Installed al ($tag) to $dest"
    if ($env:PATH -notmatch [regex]::Escape($InstallDir)) {
        Write-Host "Add $InstallDir to your PATH to run al from any shell."
    }
} finally {
    if (Test-Path $tempAsset) {
        Remove-Item -Force $tempAsset
    }
    if (Test-Path $tempChecksums) {
        Remove-Item -Force $tempChecksums
    }
}
