# Mini-Zanzibar Comprehensive Test for Windows - Legacy docker-compose Version

param(
    [string]$ProjectRoot = ".."  # Points to mini-zanzibar directory from scripts folder
)

Write-Host "=== Mini-Zanzibar Comprehensive Test (docker-compose) ===" -ForegroundColor Green

# Set the path to docker-compose.yml
$ComposeFile = "$ProjectRoot\deployments\docker\docker-compose.yml"

# Function to run docker-compose with proper path (legacy syntax)
function Invoke-DockerCompose($Command) {
    if (Test-Path $ComposeFile) {
        Write-Host "  Running: docker-compose -f `"$ComposeFile`" $Command" -ForegroundColor Gray
        docker-compose -f "$ComposeFile" $Command
    } else {
        Write-Host "❌ docker-compose.yml not found at: $ComposeFile" -ForegroundColor Red
        return $null
    }
}

# 0. Check if we're in the right location
Write-Host "0. Checking project structure..." -ForegroundColor Yellow
if (-not (Test-Path $ComposeFile)) {
    Write-Host "❌ docker-compose.yml not found at: $ComposeFile" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Project structure looks good" -ForegroundColor Green

# 1. Start services if not running
Write-Host "`n1. Starting services..." -ForegroundColor Yellow
Invoke-DockerCompose "up -d"
Start-Sleep -Seconds 10

# 2. Test health endpoint
Write-Host "`n2. Testing health endpoint..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method Get -TimeoutSec 10
    Write-Host "✅ Health check passed" -ForegroundColor Green
    $response | ConvertTo-Json
} catch {
    Write-Host "❌ Health check failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Checking application logs..." -ForegroundColor Yellow
    Invoke-DockerCompose "logs app"
    exit 1
}

# 3. Test Redis connection
Write-Host "`n3. Testing Redis connection..." -ForegroundColor Yellow
try {
    $result = Invoke-DockerCompose "exec redis redis-cli ping"
    if ($result -eq "PONG") {
        Write-Host "✅ Redis is responding: $result" -ForegroundColor Green
    } else {
        Write-Host "❌ Redis not responding correctly: $result" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Redis test failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. Check if Consul is running (required for namespaces)
Write-Host "`n4. Checking Consul availability..." -ForegroundColor Yellow
try {
    $consulResponse = Invoke-WebRequest -Uri "http://localhost:8500/v1/status/leader" -TimeoutSec 5 -UseBasicParsing
    Write-Host "✅ Consul is running" -ForegroundColor Green
} catch {
    Write-Host "❌ Consul is not running on localhost:8500" -ForegroundColor Red
    Write-Host "Please start Consul or update CONSUL_ADDRESS in your configuration" -ForegroundColor Yellow
    Write-Host "You can run: consul agent -dev" -ForegroundColor Yellow
}

# 5. Create test namespace
Write-Host "`n5. Creating test namespace..." -ForegroundColor Yellow
$namespaceData = @{
    namespace = "document"
    relations = @{
        viewer = @{
            union = @(
                @{ this = @{} }
            )
        }
        editor = @{
            union = @(
                @{ this = @{} },
                @{ computed_userset = @{ relation = "owner" } }
            )
        }
        owner = @{
            union = @(
                @{ this = @{} }
            )
        }
    }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/namespace" -Method Post -Body $namespaceData -ContentType "application/json"
    Write-Host "✅ Namespace created successfully" -ForegroundColor Green
    $response | ConvertTo-Json
} catch {
    Write-Host "❌ Namespace creation failed: $($_.Exception.Message)" -ForegroundColor Red
    # If Consul isn't running, we can't create namespaces
    if ($_.Exception.Message -like "*Consul*" -or $_.Exception.Message -like "*connection*") {
        Write-Host "This is expected because Consul is not running" -ForegroundColor Yellow
        Write-Host "Skipping remaining ACL tests..." -ForegroundColor Yellow
        Write-Host "`n=== Test Complete (Partial - Consul Required) ===" -ForegroundColor Green
        exit 0
    }
}

# 6. Wait for namespace to be processed
Start-Sleep -Seconds 2

# 7. Create test ACL
Write-Host "`n6. Creating test ACL..." -ForegroundColor Yellow
$aclData = @{
    object = "document:1"
    relation = "viewer"
    user = "user:alice"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/acl" -Method Post -Body $aclData -ContentType "application/json"
    Write-Host "✅ ACL created successfully" -ForegroundColor Green
    $response | ConvertTo-Json
    $successfulRelation = "viewer"
} catch {
    Write-Host "❌ ACL creation failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "This might be because:" -ForegroundColor Yellow
    Write-Host "  - Consul is not running" -ForegroundColor Yellow
    Write-Host "  - Namespace wasn't properly created" -ForegroundColor Yellow
    Write-Host "  - Relation doesn't exist in the namespace" -ForegroundColor Yellow
    $successfulRelation = "viewer"  # fallback for auth checks
}

# 8. Test authorization (first time - should cache)
Write-Host "`n7. Testing authorization (first request - cache miss)..." -ForegroundColor Yellow
$stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/acl/check?object=document:1&relation=$successfulRelation&user=user:alice" -Method Get
    $stopwatch.Stop()
    Write-Host "Time: $($stopwatch.ElapsedMilliseconds)ms" -ForegroundColor Cyan
    if ($response.authorized -eq $true) {
        Write-Host "✅ Authorization check passed - user is authorized" -ForegroundColor Green
    } else {
        Write-Host "⚠️  Authorization check returned false - user is NOT authorized" -ForegroundColor Yellow
        Write-Host "This is expected if ACL creation failed" -ForegroundColor Gray
    }
    $response | ConvertTo-Json
} catch {
    Write-Host "❌ Authorization check failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 9. Test authorization (second time - should use cache)
Write-Host "`n8. Testing authorization (second request - cache hit)..." -ForegroundColor Yellow
$stopwatch = [System.Diagnostics.Stopwatch]::StartNew()
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/acl/check?object=document:1&relation=$successfulRelation&user=user:alice" -Method Get
    $stopwatch.Stop()
    Write-Host "Time: $($stopwatch.ElapsedMilliseconds)ms" -ForegroundColor Cyan
    Write-Host "Note: Second request should be faster if caching works" -ForegroundColor Gray
    $response | ConvertTo-Json
} catch {
    Write-Host "❌ Authorization check failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 10. List operations
Write-Host "`n9. Testing list operations..." -ForegroundColor Yellow
Write-Host "ACLs by object (document:1):" -ForegroundColor Cyan
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/acl/object/document:1" -Method Get
    if ($response.tuples -and $response.tuples.Count -gt 0) {
        Write-Host "✅ Found $($response.tuples.Count) ACL tuples" -ForegroundColor Green
        $response.tuples | ForEach-Object { 
            Write-Host "   - $($_.object)#$($_.relation)@$($_.user)" -ForegroundColor White
        }
    } else {
        Write-Host "⚠️  No ACL tuples found for object" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ List by object failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nACLs by user (user:alice):" -ForegroundColor Cyan
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/acl/user/user:alice" -Method Get
    if ($response.tuples -and $response.tuples.Count -gt 0) {
        Write-Host "✅ Found $($response.tuples.Count) ACL tuples" -ForegroundColor Green
        $response.tuples | ForEach-Object { 
            Write-Host "   - $($_.object)#$($_.relation)@$($_.user)" -ForegroundColor White
        }
    } else {
        Write-Host "⚠️  No ACL tuples found for user" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ List by user failed: $($_.Exception.Message)" -ForegroundColor Red
}

# 11. Check Redis keys
Write-Host "`n10. Checking Redis cache keys..." -ForegroundColor Yellow
try {
    $keys = Invoke-DockerCompose "exec redis redis-cli keys '*'"
    if ($keys -and $keys -ne "") {
        Write-Host "✅ Redis keys found:" -ForegroundColor Green
        $keys.Split("`n") | ForEach-Object { 
            if ($_ -ne "" -and $_ -notlike "*Usage:*" -and $_ -notlike "*docker*") {
                Write-Host "   - $_" -ForegroundColor White
            }
        }
    } else {
        Write-Host "⚠️  No Redis keys found" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ Redis keys check failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== Test Complete ===" -ForegroundColor Green