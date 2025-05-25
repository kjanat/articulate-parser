param(
	[switch]$AsCheckmark
)

# Get the list from 'go tool dist list'
$dists = & go tool dist list

# Parse into OS/ARCH pairs
$parsed = $dists | ForEach-Object {
	$split = $_ -split '/'
	[PSCustomObject]@{ OS = $split[0]; ARCH = $split[1] }
}

# Find all unique OSes and arches, sorted
$oses = $parsed | Select-Object -ExpandProperty OS -Unique | Sort-Object
$arches = $parsed | Select-Object -ExpandProperty ARCH -Unique | Sort-Object

# Group by OS, and build custom objects
$results = foreach ($os in $oses) {
	$props = @{}
	$props.OS = $os
	foreach ($arch in $arches) {
		$hasArch = $parsed | Where-Object { $_.OS -eq $os -and $_.ARCH -eq $arch }
		if ($hasArch) {
			if ($AsCheckmark) {
				$props[$arch] = '✅'
			} else {
				$props[$arch] = $true
			}
		} else {
			$props[$arch] = $false
		}
	}
	[PSCustomObject]$props
}

# Output
$results | Format-Table -AutoSize
