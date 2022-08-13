param(
    [string] $TargetBranchName,
    [string] $RunnerTemp = [System.IO.Path]::GetTempPath()
)

$projectsJson = repoman list --format json | Out-String
$projects = ConvertFrom-Json $projectsJson

foreach ($project in $projects) {
    $projectPath = $project.projectPath
    $templatePath = $project.templatePath.Replace($projectPath, "")
    Write-Host @"

repoman clean `
    -s $projectPath `
    -o $RunnerTemp `
    -t $templatePath `
    --branch "$TargetBranchName" `
    --fail-on-clean-error `
    --https

"@

    repoman clean `
        -s $projectPath `
        -o $RunnerTemp `
        -t $templatePath `
        --branch $TargetBranchName `
        --fail-on-clean-error  `
        --https

    if ($LASTEXITCODE) {
        Write-Error "Error running repoman clean. Exit code: $LASTEXITCODE"
        exit $LASTEXITCODE
    }
}