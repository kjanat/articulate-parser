
#Requires -Version 5.1

<#
.SYNOPSIS
	Build articulate-parser for multiple platforms

.DESCRIPTION
	This script builds the articulate-parser application for multiple operating systems and architectures.
	It can build using either native Go (preferred) or WSL if Go is not available on Windows.

.PARAMETER Jobs
	Number of parallel build jobs to run (default: 4)

.PARAMETER Platforms
	Comma-separated list of platforms to build for (e.g., "windows,linux")
	Available: windows, linux, darwin, freebsd

.PARAMETER Architectures
	Comma-separated list of architectures to build for (e.g., "amd64,arm64")
	Available: amd64, arm64

.PARAMETER BuildDir
	Directory to place built binaries (default: 'build')

.PARAMETER EntryPoint
	Entry point Go file (default: 'main.go')
	Note: This script assumes the entry point is in the project root.

.PARAMETER UseWSL
	Force use of WSL even if Go is available on Windows

.PARAMETER Clean
	Clean build directory before building

.PARAMETER VerboseOutput
	Enable verbose output

.PARAMETER LdFlags
	Linker flags to pass to go build (default: "-s -w" for smaller binaries)
	Use empty string ("") to disable default ldflags

.PARAMETER SkipTests
	Skip running tests before building

.PARAMETER Version
	Version string to embed in binaries (auto-detected from git if not provided)

.PARAMETER ShowTargets
	Show available build targets and exit

.EXAMPLE
	.\build.ps1
	Build for all platforms and architectures

.EXAMPLE
	.\build.ps1 -Platforms "windows,linux" -Architectures "amd64"
	Build only for Windows and Linux on amd64

.EXAMPLE
	.\build.ps1 -Jobs 8 -VerboseOutput
	Build with 8 parallel jobs and verbose output

.EXAMPLE
	.\build.ps1 -Version "v1.2.3" -SkipTests
	Build with specific version and skip tests

.EXAMPLE
	.\build.ps1 -LdFlags "-X main.version=1.0.0"
	Build with custom ldflags (overrides default -s -w)

.EXAMPLE
	.\build.ps1 -LdFlags ""
	Build without any ldflags (disable defaults)

.EXAMPLE
	.\build.ps1 -ShowTargets
	Show all available build targets

.FUNCTIONALITY
	Build automation, Cross-platform compilation, Go builds, Multi-architecture, Parallel builds, Windows, Linux, macOS, FreeBSD, Release management, Go applications

.NOTES
	This script requires Go to be installed and available in the PATH.
	It also requires git if auto-detecting version from tags.
	If Go is not available on Windows, it will use WSL to perform the build.

	Ensure you have the necessary permissions to create directories and files in the specified BuildDir.

	For WSL builds, ensure you have a compatible Linux distribution installed and configured.

.OUTPUTS
	Outputs built binaries to the specified BuildDir.
	Displays build summary including successful and failed builds.
#>

[CmdletBinding()]
param(
	[Parameter(Mandatory = $false, Position = 0, HelpMessage = 'Number of parallel build jobs', ValueFromPipeline, ValueFromPipelineByPropertyName)]
	[int]$Jobs = 4,
	[Parameter(Mandatory = $false, Position = 1, HelpMessage = 'Comma-separated list of platforms to build for', ValueFromPipeline, ValueFromPipelineByPropertyName)]
	[string]$Platforms = 'windows,linux,darwin,freebsd',
	[Parameter(Mandatory = $false, Position = 2, HelpMessage = 'Comma-separated list of architectures to build for', ValueFromPipeline, ValueFromPipelineByPropertyName)]
	[string]$Architectures = 'amd64,arm64',
	[Parameter(Mandatory = $false, Position = 3, HelpMessage = 'Directory to place built binaries', ValueFromPipeline, ValueFromPipelineByPropertyName)]
	[string]$BuildDir = 'build',
	[Parameter(Mandatory = $false, Position = 4, HelpMessage = 'Entry point Go file', ValueFromPipeline, ValueFromPipelineByPropertyName, ValueFromRemainingArguments)]
	[string]$EntryPoint = 'main.go',
	[Parameter(Mandatory = $false, Position = 5, HelpMessage = 'Force use of WSL even if Go is available on Windows')]
	[switch]$UseWSL,
	[Parameter(Mandatory = $false, Position = 6, HelpMessage = 'Clean build directory before building')]
	[switch]$Clean,
	[Parameter(Mandatory = $false, Position = 7, HelpMessage = 'Enable verbose output')]
	[switch]$VerboseOutput,
	[Parameter(Mandatory = $false, Position = 8, HelpMessage = 'Linker flags to pass to go build')]
	[string]$LdFlags = '-s -w',
	[Parameter(Mandatory = $false, Position = 9, HelpMessage = 'Skip running tests before building')]
	[switch]$SkipTests,
	[Parameter(Mandatory = $false, Position = 10, HelpMessage = 'Version string to embed in binaries')]
	[string]$Version = '',
	[Parameter(Mandatory = $false, Position = 11, HelpMessage = 'Show available build targets and exit')]
	[switch]$ShowTargets
)

# Set error action preference
$ErrorActionPreference = 'Stop'

# Get script directory and project root
$ScriptDir = $PSScriptRoot
$ProjectRoot = Split-Path $ScriptDir -Parent
$BuildDir = Join-Path $ProjectRoot $BuildDir

# Ensure we're in the project root
Push-Location $ProjectRoot

try {
	# Show targets and exit if requested
	if ($ShowTargets) {
		Write-Host 'Available build targets:' -ForegroundColor Cyan
		
		# Get available platforms and architectures from Go toolchain
		try {
			$GoTargets = @(go tool dist list 2>$null)
			if ($LASTEXITCODE -ne 0 -or $GoTargets.Count -eq 0) {
				throw 'Failed to get target list from Go toolchain'
			}
		} catch {
			Write-Host '⚠️ Could not retrieve targets from Go. Using default targets.' -ForegroundColor Yellow
			$PlatformList = $Platforms.Split(',') | ForEach-Object { $_.Trim() }
			$ArchList = $Architectures.Split(',') | ForEach-Object { $_.Trim() }
			
			foreach ($platform in $PlatformList) {
				foreach ($arch in $ArchList) {
					$BinaryName = "articulate-parser-$platform-$arch"
					if ($platform -eq 'windows') { $BinaryName += '.exe' }
					Write-Host "  $platform/$arch -> $BinaryName" -ForegroundColor Gray
				}
			}
			return
		}

		# Filter targets from go tool dist list
		$SelectedTargets = @()
		$PlatformList = $Platforms.Split(',') | ForEach-Object { $_.Trim() }
		$ArchList = $Architectures.Split(',') | ForEach-Object { $_.Trim() }
		
		foreach ($target in $GoTargets) {
			$parts = $target.Split('/')
			$platform = $parts[0]
			$arch = $parts[1]
			
			if ($PlatformList -contains $platform -and $ArchList -contains $arch) {
				$SelectedTargets += @{
					Platform = $platform
					Arch     = $arch
					Original = $target
				}
			}
		}
		
		# Display filtered targets
		foreach ($target in $SelectedTargets) {
			$BinaryName = "articulate-parser-$($target.Platform)-$($target.Arch)"
			if ($target.Platform -eq 'windows') { $BinaryName += '.exe' }
			Write-Host "  $($target.Original) -> $BinaryName" -ForegroundColor Gray
		}
		
		# Show all available targets if verbose
		if ($VerboseOutput) {
			Write-Host "`nAll Go targets available on this system:" -ForegroundColor Cyan
			foreach ($target in $GoTargets) {
				Write-Host "  $target" -ForegroundColor DarkGray
			}
		}
		return
	}

	# Validate required files exist
	$RequiredFiles = @('go.mod', 'main.go')
	foreach ($file in $RequiredFiles) {
		if (-not (Test-Path $file)) {
			Write-Error "Required file not found: $file. Make sure you're in the project root."
			exit 1
		}
	}

	# Auto-detect version from git if not provided
	if (-not $Version) {
		try {
			$gitTag = git describe --tags --always --dirty 2>$null
			if ($LASTEXITCODE -eq 0 -and $gitTag) {
				$Version = $gitTag.Trim()
				if ($VerboseOutput) { Write-Host "✓ Auto-detected version: $Version" -ForegroundColor Green }
			} else {
				$Version = 'dev'
				if ($VerboseOutput) { Write-Host "⚠ Using default version: $Version" -ForegroundColor Yellow }
			}
		} catch {
			$Version = 'dev'
			if ($VerboseOutput) { Write-Host "⚠ Git not available, using default version: $Version" -ForegroundColor Yellow }
		}
	}

	# Get build timestamp
	$BuildTime = Get-Date -Format 'yyyy-MM-ddTHH:mm:ssZ'

	# Get commit hash if available
	$CommitHash = 'unknown'
	try {
		$gitCommit = git rev-parse --short HEAD 2>$null
		if ($LASTEXITCODE -eq 0 -and $gitCommit) {
			$CommitHash = $gitCommit.Trim()
		}
	} catch {
		# Git not available or not in a git repo
	}

	# Prepare enhanced ldflags with version info
	$VersionLdFlags = @(
		"-X main.Version=$Version",
		"-X main.BuildTime=$BuildTime",
		"-X main.CommitHash=$CommitHash"
	)

	# Combine base ldflags with version ldflags
	$AllLdFlags = @()
	if ($LdFlags) {
		# Remove quotes if present and split by space
		$BaseLdFlags = $LdFlags.Trim('"', "'").Split(' ', [StringSplitOptions]::RemoveEmptyEntries)
		$AllLdFlags += $BaseLdFlags
	}
	$AllLdFlags += $VersionLdFlags

	$EnhancedLdFlags = $AllLdFlags -join ' '

	if ($VerboseOutput) {
		Write-Host "🔍 Enhanced ldflags: '$EnhancedLdFlags'" -ForegroundColor Magenta
	}
	# Validate Go installation
	$GoAvailable = $false
	try {
		$goVersion = go version 2>$null
		if ($LASTEXITCODE -eq 0) {
			$GoAvailable = $true
			if ($VerboseOutput) { Write-Host "✓ Go is available: $goVersion" -ForegroundColor Green }
		}
	} catch {
		# Go not available
	}

	# Check if we should use WSL
	$UseWSLBuild = $UseWSL -or (-not $GoAvailable)

	if ($UseWSLBuild) {
		# Check WSL availability
		try {
			wsl.exe --status >$null 2>&1
			if ($LASTEXITCODE -ne 0) {
				throw 'WSL is not available'
			}
		} catch {
			Write-Error 'Neither Go nor WSL is available. Please install Go or WSL to build the project.'
			exit 1
		}

		Write-Host '🔄 Using WSL for build...' -ForegroundColor Yellow

		# Build script path
		$bashScript = Join-Path $ScriptDir 'build.sh'
		if (-not (Test-Path $bashScript)) {
			Write-Error "Build script not found at $bashScript"
			exit 1
		}

		# Prepare arguments for bash script
		$bashArgs = @()
		if ($Jobs -ne 4) {
			$bashArgs += '-j', $Jobs
		}
		if ($EnhancedLdFlags) {
			$bashArgs += '-ldflags', $EnhancedLdFlags
		}
		# Pass build directory and entry point
		$bashArgs += '-o', $BuildDir
		$bashArgs += '-e', $EntryPoint

		# Execute WSL build
		wsl.exe bash "$bashScript" @bashArgs
		if ($LASTEXITCODE -ne 0) {
			Write-Error "WSL build script failed with exit code $LASTEXITCODE"
			exit $LASTEXITCODE
		}
		return
	}

	# Run tests before building (unless skipped)
	if (-not $SkipTests) {
		Write-Host '🧪 Running tests...' -ForegroundColor Cyan
		$TestResult = go test -v ./... 2>&1
		if ($LASTEXITCODE -ne 0) {
			Write-Host '❌ Tests failed:' -ForegroundColor Red
			Write-Host $TestResult -ForegroundColor Red
			Write-Error 'Tests failed. Use -SkipTests to build anyway.'
			exit 1
		}
		Write-Host '✅ All tests passed' -ForegroundColor Green
	}

	# Native PowerShell build
	Write-Host '🔨 Building articulate-parser natively...' -ForegroundColor Cyan

	# Clean build directory if requested
	if ($Clean -and (Test-Path $BuildDir)) {
		Write-Host '🧹 Cleaning build directory...' -ForegroundColor Yellow
		Remove-Item $BuildDir -Recurse -Force
	}

	# Create build directory
	if (-not (Test-Path $BuildDir)) {
		New-Item -ItemType Directory -Path $BuildDir | Out-Null
	}

	# Parse platforms and architectures
	$PlatformList = $Platforms.Split(',') | ForEach-Object { $_.Trim() }
	$ArchList = $Architectures.Split(',') | ForEach-Object { $_.Trim() }

	# Validate platforms and architectures
	$ValidPlatforms = @('windows', 'linux', 'darwin', 'freebsd')
	$ValidArchs = @('amd64', 'arm64')

	foreach ($platform in $PlatformList) {
		if ($platform -notin $ValidPlatforms) {
			Write-Error "Invalid platform: $platform. Valid platforms: $($ValidPlatforms -join ', ')"
			exit 1
		}
	}

	foreach ($arch in $ArchList) {
		if ($arch -notin $ValidArchs) {
			Write-Error "Invalid architecture: $arch. Valid architectures: $($ValidArchs -join ', ')"
			exit 1
		}
	}

	# Generate build targets
	$Targets = @()
	foreach ($platform in $PlatformList) {
		foreach ($arch in $ArchList) {
			$BinaryName = "articulate-parser-$platform-$arch"
			if ($platform -eq 'windows') {
				$BinaryName += '.exe'
			}
			$Targets += @{
				Platform = $platform
				Arch     = $arch
				Binary   = $BinaryName
				Path     = Join-Path $BuildDir $BinaryName
			}
		}
	}

	Write-Host "📋 Building $($Targets.Count) targets with $Jobs parallel jobs" -ForegroundColor Cyan

	# Display targets
	if ($VerboseOutput) {
		foreach ($target in $Targets) {
			Write-Host "  - $($target.Platform)/$($target.Arch) -> $($target.Binary)" -ForegroundColor Gray
		}
	}

	# Build function
	$BuildTarget = {
		param($Target, $EnhancedLdFlags, $VerboseOutput, $BuildDir, $EntryPoint, $ProjectRoot)

		$env:GOOS = $Target.Platform
		$env:GOARCH = $Target.Arch
		$env:CGO_ENABLED = '0'

		# Construct build arguments
		$BuildArgs = @('build')
		if ($EnhancedLdFlags) {
			$BuildArgs += '-ldflags'
			$BuildArgs += "`"$EnhancedLdFlags`""
		}
		$BuildArgs += '-o'
		$BuildArgs += $Target.Path
		
		# If using custom entry point that's not main.go
		# we need to use the file explicitly to avoid duplicate declarations
		$EntryPointPath = Join-Path $ProjectRoot $EntryPoint
		$EntryPointFile = Split-Path $EntryPointPath -Leaf
		$IsCustomEntryPoint = ($EntryPointFile -ne 'main.go')
		
		if ($IsCustomEntryPoint) {
			# When using custom entry point, compile only that file
			$BuildArgs += $EntryPointPath
		} else {
			# For standard main.go, let Go find and compile all package files
			$PackagePath = Split-Path $EntryPointPath -Parent
			$BuildArgs += $PackagePath
		}
		
		# For verbose output, show the command that will be executed
		if ($VerboseOutput) {
			Write-Host "Command: go $($BuildArgs -join ' ')" -ForegroundColor DarkCyan
		}

		$LogFile = "$($Target.Path).log"

		try {
			if ($VerboseOutput) {
				Write-Host "🔨 Building $($Target.Binary)..." -ForegroundColor Yellow
			}

			$Process = Start-Process -FilePath 'go' -ArgumentList $BuildArgs -Wait -PassThru -NoNewWindow -RedirectStandardError $LogFile

			if ($Process.ExitCode -eq 0) {
				# Remove log file on success
				if (Test-Path $LogFile) {
					Remove-Item $LogFile -Force
				}
				return @{ Success = $true; Target = $Target.Binary }
			} else {
				return @{ Success = $false; Target = $Target.Binary; LogFile = $LogFile }
			}
		} catch {
			return @{ Success = $false; Target = $Target.Binary; Error = $_.Exception.Message }
		}
	}

	# Execute builds with throttling
	$RunspacePool = [runspacefactory]::CreateRunspacePool(1, $Jobs)
	$RunspacePool.Open()

	$BuildJobs = @()
	foreach ($target in $Targets) {
		$PowerShell = [powershell]::Create()
		$PowerShell.RunspacePool = $RunspacePool
		$PowerShell.AddScript($BuildTarget).AddParameters(@{
				Target          = $target
				EnhancedLdFlags = $EnhancedLdFlags
				VerboseOutput   = $VerboseOutput
				BuildDir        = $BuildDir
				EntryPoint      = $EntryPoint
				ProjectRoot     = $ProjectRoot
			}) | Out-Null

		$BuildJobs += @{
			PowerShell  = $PowerShell
			AsyncResult = $PowerShell.BeginInvoke()
			Target      = $target.Binary
		}
	}

	# Wait for results and display progress
	$Completed = 0
	$Successful = 0
	$Failed = 0

	Write-Host ''
	while ($Completed -lt $BuildJobs.Count) {
		foreach ($job in $BuildJobs | Where-Object { $_.AsyncResult.IsCompleted -and -not $_.Processed }) {
			$job.Processed = $true
			$Result = $job.PowerShell.EndInvoke($job.AsyncResult)
			$job.PowerShell.Dispose()

			$Completed++
			if ($Result.Success) {
				$Successful++
				Write-Host "✅ $($Result.Target)" -ForegroundColor Green
			} else {
				$Failed++
				if ($Result.LogFile) {
					Write-Host "❌ $($Result.Target) (see $($Result.LogFile))" -ForegroundColor Red
				} else {
					Write-Host "❌ $($Result.Target): $($Result.Error)" -ForegroundColor Red
				}
			}
		}
		Start-Sleep -Milliseconds 100
	}

	$RunspacePool.Close()
	$RunspacePool.Dispose()

	# Summary
	Write-Host ''
	Write-Host '📊 Build Summary:' -ForegroundColor Cyan
	Write-Host "  🏷️  Version:    $Version" -ForegroundColor Gray
	Write-Host "  🔨 Commit:     $CommitHash" -ForegroundColor Gray
	Write-Host "  ⏰ Build Time: $BuildTime" -ForegroundColor Gray
	Write-Host "  ✅ Successful: $Successful" -ForegroundColor Green
	Write-Host "  ❌ Failed:     $Failed" -ForegroundColor Red
	Write-Host "  📁 Output:     $BuildDir" -ForegroundColor Yellow

	if ($Successful -gt 0) {
		Write-Host ''
		Write-Host '📦 Built binaries:' -ForegroundColor Cyan
		Get-ChildItem $BuildDir -File | Where-Object { $_.Name -notlike '*.log' } | Sort-Object Name | ForEach-Object {
			$Size = [math]::Round($_.Length / 1MB, 2)
			$LastWrite = $_.LastWriteTime.ToString('HH:mm:ss')
			Write-Host "  $($_.Name.PadRight(35)) $($Size.ToString().PadLeft(6)) MB  ($LastWrite)" -ForegroundColor Gray
		}

		# Calculate total size
		$TotalSize = (Get-ChildItem $BuildDir -File | Where-Object { $_.Name -notlike '*.log' } | Measure-Object -Property Length -Sum).Sum
		$TotalSizeMB = [math]::Round($TotalSize / 1MB, 2)
		Write-Host "  $('Total:'.PadRight(35)) $($TotalSizeMB.ToString().PadLeft(6)) MB" -ForegroundColor Cyan
	}

	if ($Failed -gt 0) {
		exit 1
	}
} finally {
	Pop-Location
}
