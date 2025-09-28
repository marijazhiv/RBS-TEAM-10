# Comprehensive test flow for Mini-Zanzibar system
# Tests all API endpoints with various scenarios

param(
    [string]$BaseUrl = "http://localhost:8080"
)

# Test counters
$TotalTests = 0
$PassedTests = 0
$FailedTests = 0

# Helper function to make HTTP requests
function Invoke-TestRequest {
    param(
        [string]$Method,
        [string]$Uri,
        [hashtable]$Body = $null,
        [string]$Description
    )
    
    Write-Host "Testing: $Description" -ForegroundColor Yellow
    $global:TotalTests++
    
    try {
        $params = @{
            Uri = $Uri
            Method = $Method
            ContentType = "application/json"
            UseBasicParsing = $true
        }
        
        if ($Body) {
            $params.Body = $Body | ConvertTo-Json -Depth 10
        }
        
        $response = Invoke-RestMethod @params
        
        Write-Host "PASS: $Description" -ForegroundColor Green
        $global:PassedTests++
        return @{ Success = $true; Response = $response }
    }
    catch {
        Write-Host "FAIL: $Description" -ForegroundColor Red
        Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
        $global:FailedTests++
        return @{ Success = $false; Error = $_.Exception.Message }
    }
}

# Helper function to verify response content
function Test-Response {
    param(
        [hashtable]$Result,
        [string]$TestName,
        [scriptblock]$Validation
    )
    
    $global:TotalTests++
    
    if ($Result.Success) {
        try {
            $isValid = & $Validation $Result.Response
            if ($isValid) {
                Write-Host "PASS: $TestName" -ForegroundColor Green
                $global:PassedTests++
            } else {
                Write-Host "FAIL: $TestName - Validation failed" -ForegroundColor Red
                $global:FailedTests++
            }
        }
        catch {
            Write-Host "FAIL: $TestName - Validation error: $($_.Exception.Message)" -ForegroundColor Red
            $global:FailedTests++
        }
    } else {
        Write-Host "FAIL: $TestName - Request failed" -ForegroundColor Red
        $global:FailedTests++
    }
}

Write-Host "Starting comprehensive Mini-Zanzibar test flow..." -ForegroundColor Cyan
Write-Host "Base URL: $BaseUrl" -ForegroundColor Cyan
Write-Host ""

# Test 1: Health Check
$healthResult = Invoke-TestRequest -Method GET -Uri "$BaseUrl/health" -Description "Health check"

# Test 2: Create namespace
$namespaceBody = @{
    name = "documents"
    config = @{
        relations = @{
            owner = @{}
            editor = @{
                union = @{
                    child = @(
                        @{ this = @{} },
                        @{ computed_userset = @{ relation = "owner" } }
                    )
                }
            }
            viewer = @{
                union = @{
                    child = @(
                        @{ this = @{} },
                        @{ computed_userset = @{ relation = "editor" } }
                    )
                }
            }
        }
    }
}

$createNamespaceResult = Invoke-TestRequest -Method POST -Uri "$BaseUrl/namespace" -Body $namespaceBody -Description "Create documents namespace"

# Test 3: Get namespace
$getNamespaceResult = Invoke-TestRequest -Method GET -Uri "$BaseUrl/namespace/documents" -Description "Get documents namespace"

# Test 4: Create ACL - Alice owns doc1
$aclBody1 = @{
    object = "documents:doc1"
    relation = "owner"
    user = "user:alice"
}

$createAclResult1 = Invoke-TestRequest -Method POST -Uri "$BaseUrl/acl" -Body $aclBody1 -Description "Create ACL: Alice owns doc1"

# Test 5: Create ACL - Bob edits doc1
$aclBody2 = @{
    object = "documents:doc1"
    relation = "editor"
    user = "user:bob"
}

$createAclResult2 = Invoke-TestRequest -Method POST -Uri "$BaseUrl/acl" -Body $aclBody2 -Description "Create ACL: Bob edits doc1"

# Test 6: Create ACL - Charlie views doc1
$aclBody3 = @{
    object = "documents:doc1"
    relation = "viewer"
    user = "user:charlie"
}

$createAclResult3 = Invoke-TestRequest -Method POST -Uri "$BaseUrl/acl" -Body $aclBody3 -Description "Create ACL: Charlie views doc1"

# Test 7: Check Alice can own doc1
$checkResult1 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=documents:doc1&relation=owner&user=user:alice" -Description "Check Alice owns doc1"

# Test 8: Check Bob can edit doc1
$checkResult2 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=documents:doc1&relation=editor&user=user:bob" -Description "Check Bob edits doc1"

# Test 9: Check Charlie can view doc1
$checkResult3 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=documents:doc1&relation=viewer&user=user:charlie" -Description "Check Charlie views doc1"

# Test 10: Check unauthorized access (Dave should not have access)
$checkResult4 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=documents:doc1&relation=viewer&user=user:dave" -Description "Check Dave cannot view doc1 (should fail)"

# Test 11: List ACLs for doc1
$listResult1 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl?object=documents:doc1" -Description "List ACLs for doc1"

# Test 12: List ACLs for Alice
$listResult2 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl?user=user:alice" -Description "List ACLs for Alice"

# Test 13: Create another document with different permissions
$aclBody4 = @{
    object = "documents:doc2"
    relation = "owner"
    user = "user:bob"
}

$createAclResult4 = Invoke-TestRequest -Method POST -Uri "$BaseUrl/acl" -Body $aclBody4 -Description "Create ACL: Bob owns doc2"

# Test 14: Check Bob can own doc2
$checkResult5 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=documents:doc2&relation=owner&user=user:bob" -Description "Check Bob owns doc2"

# Test 15: Check Alice cannot access doc2
$checkResult6 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=documents:doc2&relation=viewer&user=user:alice" -Description "Check Alice cannot view doc2 (should fail)"

# Test 16: Delete ACL
$deleteResult1 = Invoke-TestRequest -Method DELETE -Uri "$BaseUrl/acl" -Body $aclBody3 -Description "Delete Charlie's view access to doc1"

# Test 17: Verify deletion worked
$checkResult7 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=documents:doc1&relation=viewer&user=user:charlie" -Description "Verify Charlie cannot view doc1 after deletion"

# Test 18: Create a different namespace
$namespaceBody2 = @{
    name = "files"
    config = @{
        relations = @{
            owner = @{}
            reader = @{
                union = @{
                    child = @(
                        @{ this = @{} },
                        @{ computed_userset = @{ relation = "owner" } }
                    )
                }
            }
        }
    }
}

$createNamespaceResult2 = Invoke-TestRequest -Method POST -Uri "$BaseUrl/namespace" -Body $namespaceBody2 -Description "Create files namespace"

# Test 19: Create ACL in different namespace
$aclBody5 = @{
    object = "files:config.json"
    relation = "owner"
    user = "user:admin"
}

$createAclResult5 = Invoke-TestRequest -Method POST -Uri "$BaseUrl/acl" -Body $aclBody5 -Description "Create ACL: Admin owns config.json"

# Test 20: Check admin can access file
$checkResult8 = Invoke-TestRequest -Method GET -Uri "$BaseUrl/acl/check?object=files:config.json&relation=owner&user=user:admin" -Description "Check Admin owns config.json"

Write-Host ""

# Test Summary
Write-Host "=== Test Summary ===" -ForegroundColor Yellow
Write-Host "Total Tests: $TotalTests"
Write-Host "Passed: $PassedTests" -ForegroundColor Green
Write-Host "Failed: $FailedTests" -ForegroundColor Red

if ($FailedTests -eq 0) {
    Write-Host "SUCCESS: All tests passed!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "FAILED: Some tests failed. Check the output above." -ForegroundColor Red
    exit 1
}