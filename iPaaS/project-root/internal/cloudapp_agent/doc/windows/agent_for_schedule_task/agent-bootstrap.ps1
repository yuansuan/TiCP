# the reason why add bootstrap.ps1 to call daemon.ps1 :
# hide the agent running shell window

# set script dir path to workdir
$scriptDirectory = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location -Path $scriptDirectory

Start-Process powershell -ArgumentList '-ExecutionPolicy Unrestricted -NoProfile -Sta -File agent-daemon.ps1' -WindowStyle Hidden
