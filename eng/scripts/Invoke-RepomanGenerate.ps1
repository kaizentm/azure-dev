param(
    [string] $TargetBranchName,
    [string] $ResultsFileLocation = "$([System.IO.Path]::GetTempPath())/repoman.md",
    [string] $RunnerTemp = [System.IO.Path]::GetTempPath(),
    [switch] $WhatIf
)

$projectsJson = repoman list --format json | Out-String
$projects = ConvertFrom-Json $projectsJson

foreach ($project in $projects) {
    $additionalParameters = '--update'
    if ($WhatIf) {
        $additionalParameters = ''
    }

    $projectPath = $project.projectPath
    $templatePath = $projectPath.templatePath

    Write-Host @"
repoman generate `
    -s $projectPath `
    -o $RunnerTemp `
    -t $templatePath `
    --branch "$TargetBranchName" `
    --https `
    --fail-on-update-error `
    --resultsFile $ResultsFileLocation `
    $additionalParameters
"@

    repoman generate `
        -s $projectPath `
        -o $RunnerTemp `
        -t $templatePath `
        --branch `"$TargetBranchName`" `
        --https `
        --fail-on-update-error `
        --resultsFile $ResultsFileLocation `
        $additionalParameters

    if ($LASTEXITCODE) {
        Write-Error "Error running repoman generate. Exit code: $LASTEXITCODE"
        exit $LASTEXITCODE
    }
}