<#
1. 采用单行中文注释并且换行符为Unix(LF)时，某些中文字符可能会导致windows上的powershell脚本无法正确识别换行符

2. 推测跟换行符以及powershell的编码有关. 例如 下面的代码将导致变量hostsPath无法识别
#立
$hostsPath = "aaaa"

3. 所以保险点，本脚本里的注释尽量使用英文或者多行注释方式
#>

<#
多会话时，非Administrator用户无法通过配置权限获得C:\Windows\下的读写权限，故将挂载日志存放至用户目录
#>
$userProfile = [System.Environment]::GetFolderPath('UserProfile')
$logFolderPath = Join-Path -Path $userProfile -ChildPath "logs"
# logFile 路径为C:\Users\<username>\logs\start-share.log
$logFile = Join-Path -Path $logFolderPath -ChildPath "start-share.log"

$hostsPath = "C:\Windows\System32\drivers\etc\hosts"
<#
确保logs文件夹存在，如果不存在则创建
#>
if (-not (Test-Path -Path $logFolderPath)) {
    New-Item -Path $logFolderPath -ItemType Directory
}

<#
因为windows不允许一个用户使用多个凭证和同一个smb server建立多个连接，所以为每个共享点，创建一个host映射来规避这个限制。
#>
function AddHost {
    param(
        [string]$shareHostIp,
        [string]$shareServerHost
    )
    try {
        $newHostsEntry = "`n$shareHostIp    $shareServerHost"
        $hostsContent = Get-Content -Path $hostsPath -Raw
        if ($hostsContent -notmatch $newHostsEntry){
            Add-Content -Path $hostsPath -Value $newHostsEntry
        }
    }
    catch {
        LOG -level "ERROR" -msg "Error Add Host: $newHostsEntry, $_"
    }
}

function MountShare {
    param(
        [string]$mountSrc,
        [string]$mountPoint,
        [string]$mountHost,
        [string]$mountUser,
        [string]$mountPass
    )
    $shareServerHost = $mountSrc

    AddHost -shareHostIp $mountHost -shareServerHost $shareServerHost

    $command = "net use '$mountPoint' '\\$shareServerHost\$mountSrc' /USER:'$mountUser' '$mountPass' 2>&1"

    $retryCount = 0
    $maxRetries = 20
    while ($retryCount -lt $maxRetries) {
        try {
            LOG -level "INFO" -msg "$command"
            $output = Invoke-Expression $command
            if ($LASTEXITCODE -eq 0) {
                LOG -level "INFO" -msg "mount from $mountSrc to $mountPoint success!"
                break
            } else {
                LOG -level "ERROR" -msg "mount from $mountSrc to $mountPoint failed!, error output in below, count: $retryCount"
                LOG -level "ERROR" -msg $output
            }
        }
        catch {
            LOG -level "ERROR" -msg "Error mounting share: $_"
        }
        $retryCount++
        if ($retryCount -lt $maxRetries) {
            Start-Sleep -Seconds 5
        }
    }
    if ($retryCount -ge $maxRetries) {
        LOG -level "ERROR" -msg "All $maxRetries attempts to mount from $mountSrc to $mountPoint failed."
    }
}

function LOG {
    param(
        [string]$msg,
        [string]$level
    )

    "$(Get-Date) $level $msg" | Out-File -Append $logFile
}

function initEnvs {
    param(
        [string]$file
    )

    $contents = Get-Content $file

    foreach ($line in $contents) {
        if ($line -match '^(.*?)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()

            if ([string]::IsNullOrEmpty($value)) {
                $value = ""
            }

            if ($key -eq "SHARE_SERVER") {
                LOG -level "INFO" -msg "get SHARE_SERVER from envfile, value = $value"
                $mountHost = $value
            }

            if ($key -eq "SHARE_USERNAME") {
                LOG -level "INFO" -msg "get SHARE_USERNAME from envfile, value = $value"
                $mountUser = $value
            }

            if ($key -eq "SHARE_PASSWORD") {
                LOG -level "INFO" -msg "get SHARE_PASSWORD from envfile, value = $value"
                $mountPass = $value
            }

            if ($key -eq "SHARE_MOUNT_PATHS") {
                LOG -level "INFO" -msg "get SHARE_MOUNT_PATHS from envfile, value = $value"
                $mountSubPaths = $value
            }

            if ($key -eq "SHARE_REGISTER_ADDRESS") {
                LOG -level "INFO" -msg "get SHARE_REGISTER_ADDRESS from envfile, value = $value"
                $shareRegisterAddress = $value
            }
        }
    }

    $mountHost
    $mountUser
    $mountPass
    $mountSubPaths
    $shareRegisterAddress
}

function Main {
    Set-Location -PassThru ([Environment]::GetEnvironmentVariable("YS_PATH", "Machine"))
    $envfile = [Environment]::GetEnvironmentVariable("YS_PATH", "Machine") + "\agent_env"

    LOG -level "INFO" -msg "starting to loop check env_file exist or not"
    while (!(Test-Path $envfile)) {
        LOG -level "WARN" -msg "env_file not exist"
        Start-Sleep -Seconds 1
    }
    LOG -level "INFO" -msg "get SHARE_SERVER env done, starting to mount shares"

    $mountHost, $mountUser, $mountPass, $mountSubPaths, $shareRegisterAddress = initEnvs -file $envfile

    Write-Host "start to mount user volume, please do not close it..."

    if (![string]::IsNullOrEmpty($mountSubPaths)) {
        Add-Content -Path $hostsPath -Value "`r`n"

        $fields = $mountSubPaths.Split(',')
        $shareUsername = $mountUser.Split(',')
        $sharePassword = $mountPass.Split(',')
        foreach ($index in 0..($fields.Length-1)) {
            $keyValue = $fields[$index].Split('=')
            $mountSrc = $shareUsername[$index]
            $mountPassword = $sharePassword[$index]
            $mountPoint = $keyValue[1]

            LOG -level "INFO" -msg "starting to mount from $mountSrc to $mountPoint"

            $fulfilledMountSrc = "$mountSrc"
            LOG -level "INFO" -msg "fullfilled mount src is: $fulfilledMountSrc"

            MountShare -mountPoint $mountPoint -mountSrc $fulfilledMountSrc -mountHost $mountHost -mountUser $mountSrc -mountPass $mountPassword

            $volumeName = "Data"
            if ($mountSrc -match "\\") {
                $substrings = $mountSrc -split "\\"
                $volumeName = $substrings[-1]
            } else {
                $volumeName = $mountSrc
            }

            $app = New-Object -ComObject shell.application
            $app.NameSpace("$mountPoint").self.name = "$volumeName"
            LOG -level "INFO" -msg "label $mountPoint to $volumeName"
        }
    } else {
        LOG -level "WARN" -msg "Empty mount paths"
    }
}

Main
