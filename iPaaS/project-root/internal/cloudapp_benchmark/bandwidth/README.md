# 性能测试-带宽定量测试工具

## 准备工作

使用 `make win` 编译出windows下运行工具`benchmark-bandwidth.exe`

将`benchmark-bandwidth.exe` 与 `benchmark.mp4`(下载地址：http://wiki.yuansuan.com/pages/viewpage.action?pageId=26985064 带宽基准测试视频) 放置与同级目录


### 预热

双击benchmark.mp4文件，并使其最大化播放，关闭mp4（此处为了让播放器记住最大化播放该测试视频）

## 运行

于`benchmark-bandwidth.exe`目录打开`cmd`，并且`.\benchmark-bandwidth.exe` 运行测试工具，预期测试工具会自动打开播放benchmark.mp4文件，并且播放结束后产生测试报告。

## 测试方案

参考跳水评分规则，建议运行5次工具，得到5次报告，去掉数值最大/小的2份，其余3份取平均数得到本次性能测试数据。

## 测试报告

一份json文件如下：
```json
{
    "upload_average_bandwidth": "0.390508 MB/s",   // 上行平均带宽
    "upload_peek_bandwidth": "0.829853 MB/s",      // 上行最大带宽
    "download_average_bandwidth": "0.002791 MB/s", // 下行平均带宽
    "download_peek_bandwidth": "0.007519 MB/s"     // 下行最大带宽
}
```

## 自定义参数
```shell
PS C:\Windows\ys\benchmark> .\benchmark-bandwidth.exe -h
bandwidth benchmark tools

Usage:
  benchmark-bandwidth [flags]

Flags:
      --check-interval duration   net traffic check interval e.g. [1s | 1m] (default 1s)
      --duration duration         net traffic collect duration e.g. [60s | 10m] (default 1m0s)
  -h, --help                      help for benchmark-bandwidth
      --net-name string           catch bandwidth net name
      --video-path string         video path for benchmark (default "benchmark.mp4")
PS C:\Windows\ys\benchmark>
```
check-interval 网络流量检测间隔，默认1s

duration 网络流量监测时长，默认1分钟

net-name 流量监控的网卡

video-path 测试视频的路径