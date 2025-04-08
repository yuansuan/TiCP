# set script dir path to workdir
$scriptDirectory = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location -Path $scriptDirectory
$ysPath = [Environment]::GetEnvironmentVariable("YS_PATH", "Machine")

$logFile = $ysPath + "\logs\agent-daemon.log"

function LOG
{
    param(
        [string]$msg,
        [string]$level
    )

    "$( Get-Date ) $level $msg" | Out-File -Append $logFile
}

function Main
{
    $agentProcess = Get-Process -Name "agent" -ErrorAction SilentlyContinue

    if ($agentProcess -eq $null)
    {
        LOG -level INFO "starting agent.exe"
        C:\Windows\ys\agent.exe
    }
    else
    {
        LOG -level INFO "agent.exe already running"
    }
}

Main
