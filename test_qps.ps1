# 接口地址
$url = "http://localhost:8080/api/get?tenant_id=2&app_id=2"

# 并发线程数
$concurrency = 50

# 测试持续时间（秒）
$duration = 10

# 请求计数
$script:count = 0

# 开始时间
$start = Get-Date

# 创建线程任务
$jobs = @()

for ($i = 0; $i -lt $concurrency; $i++) {
    $jobs += Start-Job -ScriptBlock {
        param($url,$duration,$start)

        $localCount = 0

        while (((Get-Date) - $start).TotalSeconds -lt $duration) {
            try {
                Invoke-RestMethod -Uri $url -Method Get | Out-Null
                $localCount++
            } catch {}
        }

        return $localCount

    } -ArgumentList $url,$duration,$start
}

# 等待任务完成
$jobs | Wait-Job

# 统计结果
$total = 0
foreach ($job in $jobs) {
    $total += Receive-Job $job
}

$qps = [math]::Round($total / $duration,2)

Write-Host "总请求数: $total"
Write-Host "测试时间: $duration 秒"
Write-Host "平均 QPS: $qps"