# NATMap API 测试脚本

$BaseUrl = "http://localhost:8080"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "NATMap API 测试" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# 1. 健康检查
Write-Host "`n1. 测试健康检查接口..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/test" -Method GET
    Write-Host "   结果: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 2. 查询映射（应该返回 not found）
Write-Host "`n2. 测试查询映射接口（空数据库）..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/get?tenant_id=1&app_id=1" -Method GET
    Write-Host "   结果: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    $errorResponse = $_.ErrorDetails.Message | ConvertFrom-Json
    Write-Host "   结果: $($errorResponse | ConvertTo-Json)" -ForegroundColor Green
}

# 3. 创建租户
Write-Host "`n3. 创建租户..." -ForegroundColor Yellow
try {
    $body = @{ tenant_name = "测试租户" } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/admin?type=tenant" -Method POST -ContentType "application/json" -Body $body
    Write-Host "   结果: $($response | ConvertTo-Json)" -ForegroundColor Green
    $script:tenantId = $response.data.id
    Write-Host "   租户ID: $tenantId" -ForegroundColor Cyan
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 4. 创建应用
Write-Host "`n4. 创建应用..." -ForegroundColor Yellow
try {
    $body = @{ 
        tenant_id = [int]$tenantId
        app_name = "测试应用"
        description = "这是一个测试应用"
    } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/admin?type=app" -Method POST -ContentType "application/json" -Body $body
    Write-Host "   结果: $($response | ConvertTo-Json)" -ForegroundColor Green
    $script:appId = $response.data.id
    Write-Host "   应用ID: $appId" -ForegroundColor Cyan
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 5. 创建映射
Write-Host "`n5. 创建映射..." -ForegroundColor Yellow
try {
    $body = @{ 
        tenant_id = [int]$tenantId
        app_id = [int]$appId
        public_ip = "36.96.128.250"
        public_port = 10472
        local_ip = "192.168.1.100"
        local_port = 8080
        protocol = "tcp"
    } | ConvertTo-Json
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/admin?type=mapping" -Method POST -ContentType "application/json" -Body $body
    Write-Host "   结果: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 6. 查询映射（应该返回数据）
Write-Host "`n6. 测试查询映射接口（有数据）..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/get?tenant_id=$tenantId&app_id=$appId" -Method GET
    Write-Host "   结果: $($response | ConvertTo-Json)" -ForegroundColor Green
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 7. 再次查询（测试缓存）
Write-Host "`n7. 再次查询映射（测试缓存）..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/get?tenant_id=$tenantId&app_id=$appId" -Method GET
    Write-Host "   结果: $($response | ConvertTo-Json)" -ForegroundColor Green
    if ($response._cache -eq "HIT") {
        Write-Host "   ✅ 缓存命中！" -ForegroundColor Green
    }
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 8. 列出所有租户
Write-Host "`n8. 列出所有租户..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/admin?type=tenant" -Method GET
    Write-Host "   结果: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 9. 列出所有应用
Write-Host "`n9. 列出所有应用..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/admin?type=app" -Method GET
    Write-Host "   结果: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

# 10. 列出所有映射
Write-Host "`n10. 列出所有映射..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$BaseUrl/api/admin?type=mapping" -Method GET
    Write-Host "   结果: $($response | ConvertTo-Json -Depth 3)" -ForegroundColor Green
} catch {
    Write-Host "   错误: $_" -ForegroundColor Red
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "测试完成！" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
