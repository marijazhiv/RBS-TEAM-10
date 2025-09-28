# Simple Mini-Zanzibar Test Script
# This is a simplified version to avoid PowerShell parsing issues

$BaseURL = "http://localhost:8080"

Write-Host "=== Mini-Zanzibar Simple Test ===" -ForegroundColor Yellow
Write-Host ""

# Test 1: Health Check
Write-Host "Testing Health Check..." -ForegroundColor Blue
try {
    $response = Invoke-RestMethod -Uri "$BaseURL/health" -Method GET
    Write-Host "PASS: Health check successful" -ForegroundColor Green
    Write-Host "Response: $($response | ConvertTo-Json -Compress)" -ForegroundColor Gray
} catch {
    Write-Host "FAIL: Health check failed - $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 2: Create Namespace
Write-Host "Testing Namespace Creation..." -ForegroundColor Blue
$namespaceData = @"
{
  "namespace": "doc",
  "relations": {
    "owner": {},
    "editor": {
      "union": [
        {"this": {}},
        {"computed_userset": {"relation": "owner"}}
      ]
    },
    "viewer": {
      "union": [
        {"this": {}},
        {"computed_userset": {"relation": "editor"}}
      ]
    }
  }
}
"@

try {
    $headers = @{ "Content-Type" = "application/json" }
    $response = Invoke-RestMethod -Uri "$BaseURL/namespace" -Method POST -Body $namespaceData -Headers $headers
    Write-Host "PASS: Namespace created" -ForegroundColor Green
    Write-Host "Response: $($response | ConvertTo-Json -Compress)" -ForegroundColor Gray
} catch {
    Write-Host "FAIL: Namespace creation failed - $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 3: Create ACL
Write-Host "Testing ACL Creation..." -ForegroundColor Blue
$aclData = @"
{
  "object": "doc:readme",
  "relation": "owner",
  "user": "user:alice"
}
"@

try {
    $headers = @{ "Content-Type" = "application/json" }
    $response = Invoke-RestMethod -Uri "$BaseURL/acl" -Method POST -Body $aclData -Headers $headers
    Write-Host "PASS: ACL created" -ForegroundColor Green
    Write-Host "Response: $($response | ConvertTo-Json -Compress)" -ForegroundColor Gray
} catch {
    Write-Host "FAIL: ACL creation failed - $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# Test 4: Check Authorization
Write-Host "Testing Authorization Check..." -ForegroundColor Blue
$checkUrl = "$BaseURL/acl/check?object=doc:readme&relation=owner&user=user:alice"
try {
    $response = Invoke-RestMethod -Uri $checkUrl -Method GET
    Write-Host "PASS: Authorization check completed" -ForegroundColor Green
    Write-Host "Response: $($response | ConvertTo-Json -Compress)" -ForegroundColor Gray
    
    if ($response.authorized -eq $true) {
        Write-Host "SUCCESS: Alice has owner access as expected" -ForegroundColor Green
    } else {
        Write-Host "UNEXPECTED: Alice should have owner access" -ForegroundColor Yellow
    }
} catch {
    Write-Host "FAIL: Authorization check failed - $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

Write-Host "=== Test Complete ===" -ForegroundColor Yellow
Write-Host ""
Write-Host "Next Steps:" -ForegroundColor Cyan
Write-Host "1. Try more ACL operations" -ForegroundColor White
Write-Host "2. Test computed usersets (currently not implemented)" -ForegroundColor White
Write-Host "3. Check the server logs for any issues" -ForegroundColor White