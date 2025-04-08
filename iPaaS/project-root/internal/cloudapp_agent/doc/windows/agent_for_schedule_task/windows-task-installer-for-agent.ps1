# Get the path to the script to be started
$startScriptPath = $args[0]

if (-not$startScriptPath)
{
    $startScriptPath = Join-Path $PSScriptRoot "agent-bootstrap.ps1"
}
else
{
    # Check if the provided path is relative or absolute
    if ($startScriptPath -match "^\\")
    {
        $startScriptPath = $startScriptPath
    }
    else
    {
        $startScriptPath = Join-Path $PSScriptRoot $startScriptPath
    }
}

# Create a scheduled task to run the PowerShell script on logon
$action = New-ScheduledTaskAction -Execute 'powershell.exe' -Argument "-WindowStyle Hidden -File `"$startScriptPath`""
$trigger = New-ScheduledTaskTrigger -AtLogOn -User Administrator
$taskAction = Register-ScheduledTask -Action $action -Trigger $trigger -TaskName "ys-agent" -User "Administrator" -RunLevel Highest
