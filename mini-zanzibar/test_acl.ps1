# PowerShell script to test ACL creation
$baseUrl = "http://localhost:8080"

Write-Host "Testing Mini-Zanzibar ACL Creation..." -ForegroundColor Green

# Test 1: Create doc namespace first
Write-Host "`n1. Creating doc namespace..." -ForegroundColor Yellow
$namespaceData = @{
    namespace = "doc"
    relations = @{
        owner = @{}
        editor = @{
            union = @(
                @{ this = @{} },
                @{ computed_userset = @{ relation = "owner" } }
            )
        }
        viewer = @{
            union = @(
                @{ this = @{} },
                @{ computed_userset = @{ relation = "editor" } }
            )
        }
    }
} | ConvertTo-Json -Depth 10

try {
    $response = Invoke-RestMethod -Uri "$baseUrl/namespace" -Method POST -Body $namespaceData -ContentType "application/json"
    Write-Host "✅ Namespace created: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "❌ Namespace creation failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Response: $($_.Exception.Response.StatusCode)" -ForegroundColor Red
}

# Test 2: Create Alice as owner of document4
Write-Host "`n2. Creating Alice as owner of document4..." -ForegroundColor Yellow
$aclData = @{
    object = "doc:document4"
    relation = "owner"
    user = "user:alice"
} | ConvertTo-Json

try {
    $headers = @{
        "X-User-ID" = "user:alice"
        "Content-Type" = "application/json"
    }
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/acl" -Method POST -Body $aclData -Headers $headers
    Write-Host "✅ Alice ownership created: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "❌ Alice ownership creation failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody" -ForegroundColor Red
    }
}

# Test 3: Alice shares document4 with Bob as viewer
Write-Host "`n3. Alice sharing document4 with Bob as viewer..." -ForegroundColor Yellow
$shareData = @{
    object = "doc:document4"
    relation = "viewer"
    user = "user:bob"
} | ConvertTo-Json

try {
    $headers = @{
        "X-User-ID" = "user:alice"
        "Content-Type" = "application/json"
    }
    $response = Invoke-RestMethod -Uri "$baseUrl/api/v1/acl" -Method POST -Body $shareData -Headers $headers
    Write-Host "✅ Document shared successfully: $($response.message)" -ForegroundColor Green
} catch {
    Write-Host "❌ Document sharing failed: $($_.Exception.Message)" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host "Response body: $responseBody" -ForegroundColor Red
    }
}

# Test 4: Check if Bob has access
Write-Host "`n4. Checking if Bob has viewer access to document4..." -ForegroundColor Yellow
try {
    $checkUrl = "$baseUrl/api/v1/acl/check?object=doc:document4&relation=viewer&user=user:bob"
    $response = Invoke-RestMethod -Uri $checkUrl -Method GET
    Write-Host "✅ Bob's access check: Authorized = $($response.authorized)" -ForegroundColor Green
} catch {
    Write-Host "❌ Access check failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`nTest completed!" -ForegroundColor Green