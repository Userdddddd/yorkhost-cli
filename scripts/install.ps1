# Installe le CLI YorkHost sous Windows (dernière version, ou $env:YORKHOST_VERSION).
#
#   irm https://raw.githubusercontent.com/YorkHost-fr/yorkhost-cli/main/scripts/install.ps1 | iex
#
# Destination : $env:YORKHOST_INSTALL_DIR, sinon %LOCALAPPDATA%\Programs\yorkhost
# (ajouté au PATH de l'utilisateur). La somme SHA-256 est vérifiée.
$ErrorActionPreference = 'Stop'
$repo = 'YorkHost-fr/yorkhost-cli'

$arch = switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { 'x86_64' }
    'ARM64' { 'arm64' }
    default { throw "Architecture non prise en charge : $env:PROCESSOR_ARCHITECTURE" }
}

$version = $env:YORKHOST_VERSION
if (-not $version) {
    $version = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
}
$num = $version.TrimStart('v')
$archive = "yorkhost_${num}_windows_${arch}.zip"
$base = "https://github.com/$repo/releases/download/$version"

$tmp = Join-Path ([IO.Path]::GetTempPath()) ([Guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null
try {
    Write-Host "Téléchargement de yorkhost $version (windows/$arch)…"
    Invoke-WebRequest "$base/$archive" -OutFile "$tmp\$archive" -UseBasicParsing
    Invoke-WebRequest "$base/checksums.txt" -OutFile "$tmp\checksums.txt" -UseBasicParsing

    $line = Select-String -Path "$tmp\checksums.txt" -Pattern ([regex]::Escape($archive) + '$') | Select-Object -First 1
    if (-not $line) { throw "Somme de contrôle absente pour $archive" }
    $expected = ($line.Line -split '\s+')[0].ToLower()
    $actual = (Get-FileHash "$tmp\$archive" -Algorithm SHA256).Hash.ToLower()
    if ($expected -ne $actual) { throw "Somme de contrôle invalide (attendu $expected, obtenu $actual)" }

    Expand-Archive "$tmp\$archive" -DestinationPath "$tmp\x" -Force
    $dir = $env:YORKHOST_INSTALL_DIR
    if (-not $dir) { $dir = Join-Path $env:LOCALAPPDATA 'Programs\yorkhost' }
    New-Item -ItemType Directory -Path $dir -Force | Out-Null
    Copy-Item "$tmp\x\yorkhost.exe" (Join-Path $dir 'yorkhost.exe') -Force

    $userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
    if (($userPath -split ';') -notcontains $dir) {
        [Environment]::SetEnvironmentVariable('Path', "$userPath;$dir", 'User')
        Write-Host "$dir ajouté au PATH (ouvrez un nouveau terminal)."
    }
    Write-Host "Installé : $dir\yorkhost.exe"
    Write-Host 'Pour commencer : yorkhost login'
}
finally {
    Remove-Item -Recurse -Force $tmp
}
