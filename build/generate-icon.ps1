# Windows 图标生成脚本（PowerShell）
# 使用 .NET 将 SVG 转换为 PNG，然后生成应用图标资源
# 需要: Windows 10+ with .NET Framework

param(
    [string]$SvgPath = "frontend/public/favicon.svg",
    [string]$OutputDir = "build"
)

$ErrorActionPreference = "Stop"

Write-Host "=== General Agent 图标生成 ==="

# 输出路径
$png256 = Join-Path $OutputDir "appicon.png"
$icoPath = Join-Path $OutputDir "appicon.ico"

# 确保输出目录存在
if (-not (Test-Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

# 1. SVG → PNG (256x256)
# .NET 不原生支持 SVG，使用备用方案
Write-Host "[1/3] 正在生成 256x256 PNG..."
try {
    # 使用 MSHTML (IE engine) 渲染 SVG
    Add-Type -AssemblyName System.Drawing
    Add-Type -AssemblyName System.Windows.Forms

    $svgContent = Get-Content $SvgPath -Raw -Encoding UTF8
    $html = @"
<!DOCTYPE html><html><head><meta charset='utf-8'></head>
<body style='margin:0'>$svgContent</body></html>
"@

    $wb = New-Object System.Windows.Forms.WebBrowser
    $wb.ScrollBarsEnabled = $false
    $wb.ScriptErrorsSuppressed = $true
    $wb.Size = New-Object System.Drawing.Size(256, 256)

    $tempHtml = [System.IO.Path]::GetTempFileName() + ".html"
    $html | Out-File $tempHtml -Encoding UTF8
    $wb.Navigate("file:///$tempHtml")

    # 等待渲染
    while ($wb.ReadyState -ne 'Complete') {
        [System.Windows.Forms.Application]::DoEvents()
    }
    Start-Sleep -Seconds 2

    $bitmap = New-Object System.Drawing.Bitmap 256, 256
    $wb.DrawToBitmap($bitmap, (New-Object System.Drawing.Rectangle 0, 0, 256, 256))
    $bitmap.Save($png256, [System.Drawing.Imaging.ImageFormat]::Png)
    $wb.Dispose()
    Remove-Item $tempHtml
    Write-Host "  ✓ appicon.png 已生成"
} catch {
    Write-Host "  ⚠ SVG 自动转换失败: $_"
    Write-Host "  请手动将 favicon.svg 转换为 256x256 PNG 并保存为 build/appicon.png"
    Write-Host "  推荐在线工具: https://svgtopng.com/"
    exit 1
}

# 2. 生成 Windows .ico（包含多尺寸）
Write-Host "[2/3] 正在生成 .ico 文件..."
try {
    $bitmap = [System.Drawing.Image]::FromFile($png256)
    $icon = [System.Drawing.Icon]::FromHandle((New-Object System.Drawing.Bitmap 32, 32).GetHicon())

    # 使用 System.Drawing 生成包含常用尺寸的 ico
    $sizes = @(16, 24, 32, 48, 64, 128, 256)
    $ms = New-Object System.IO.MemoryStream
    $writer = New-Object System.IO.BinaryWriter($ms)

    # ICO header
    $writer.Write([Int16]0)      # Reserved
    $writer.Write([Int16]1)      # ICO type
    $writer.Write([Int16]$sizes.Length)  # Number of images

    $images = @()
    $dataStart = 6 + 16 * $sizes.Length

    foreach ($size in $sizes) {
        $resized = New-Object System.Drawing.Bitmap($bitmap, $size, $size)
        $pngStream = New-Object System.IO.MemoryStream
        $resized.Save($pngStream, [System.Drawing.Imaging.ImageFormat]::Png)
        $pngBytes = $pngStream.ToArray()
        $images += @{ Size = $size; Data = $pngBytes; Offset = $dataStart }
        $dataStart += $pngBytes.Length
        $resized.Dispose()
        $pngStream.Dispose()
    }

    # Write directory entries
    foreach ($img in $images) {
        $s = if ($img.Size -ge 256) { 0 } else { $img.Size }
        $writer.Write([Byte]$s)        # Width
        $writer.Write([Byte]$s)        # Height
        $writer.Write([Byte]0)         # Colors
        $writer.Write([Byte]0)         # Reserved
        $writer.Write([Int16]1)        # Planes
        $writer.Write([Int16]32)       # Bits per pixel
        $writer.Write([Int32]$img.Data.Length)
        $writer.Write([Int32]$img.Offset)
    }

    # Write image data
    foreach ($img in $images) {
        $writer.Write($img.Data)
    }

    $writer.Flush()
    [System.IO.File]::WriteAllBytes($icoPath, $ms.ToArray())
    $ms.Dispose()
    $writer.Dispose()
    $bitmap.Dispose()

    Write-Host "  ✓ appicon.ico 已生成"
} catch {
    Write-Host "  ⚠ ICO 生成失败: $_"
    exit 1
}

# 3. 提示下一步
Write-Host "[3/3] 图标文件就绪"
Write-Host ""
Write-Host "下一步: 在项目根目录执行以下命令完成打包:"
Write-Host "  wails3 task generate:syso ARCH=amd64"
Write-Host "  wails3 task package:windows"
Write-Host ""
Write-Host "或手动构建:"
Write-Host "  go build -o bin/GeneralAgent.exe ./cmd/wails/"