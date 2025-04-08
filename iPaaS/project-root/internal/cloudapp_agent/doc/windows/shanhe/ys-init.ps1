$logFile = "C:\Windows\ys\logs\ys-init.log"

function LOG {
    param(
        [string]$msg,
        [string]$level
    )

    "$(Get-Date) $level $msg" | Out-File -Append $logFile
}

# 循环检测文件是否存在
while (-Not (Test-Path -Path "C:\Windows\ys\agent_env")) {
    LOG -level "WARN" "agent_env not exist"
    # 文件不存在时检查 init.ps1 文件是否存在
    if (Test-Path -Path "C:\etc\init.ps1") {
        LOG -level "INFO" "C:\etc\init.ps1 exist"
        # init.ps1 文件存在时执行
        & "C:\etc\init.ps1"
        break
    }
    LOG -level "ERROR" "C:\etc\init.ps1 not exist"

    Start-Sleep -Seconds 5
}