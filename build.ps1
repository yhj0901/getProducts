# build.ps1

# Windows 64bit 빌드
Write-Host "Building for Windows 64bit..."
$env:GOARCH = "amd64"
$env:GOOS = "windows"
go build -tags windows -o getProduct_windows_amd64.exe
if ($?) {
    Write-Host "Windows 64bit build successful"
} else {
    Write-Host "Windows 64bit build failed"
    exit 1
}

# Windows 32bit 빌드
Write-Host "`nBuilding for Windows 32bit..."
$env:GOARCH = "386"
$env:GOOS = "windows"
go build -tags windows -o getProduct_windows_x86.exe
if ($?) {
    Write-Host "Windows 32bit build successful"
} else {
    Write-Host "Windows 32bit build failed"
    exit 1
}

# Linux 64bit 빌드
Write-Host "`nBuilding for Linux 64bit..."
$env:GOARCH = "amd64"
$env:GOOS = "linux"
go build -tags linux -o getProduct_linux_amd64
if ($?) {
    Write-Host "Linux 64bit build successful"
} else {
    Write-Host "Linux 64bit build failed"
    exit 1
}

Write-Host "`nAll builds completed successfully!"