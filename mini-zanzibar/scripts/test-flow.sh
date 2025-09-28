#!/bin/bash

# Mini-Zanzibar Complete Flow Test Script
# This script tests the entire authorization flow

BASE_URL="http://localhost:8080"
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print test results
print_test_result() {
    local test_name="$1"
    local expected="$2"
    local actual="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [[ "$actual" == *"$expected"* ]]; then
        echo -e "${GREEN}✓ PASS${NC}: $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ FAIL${NC}: $test_name"
        echo -e "  Expected: $expected"
        echo -e "  Actual: $actual"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# Function to make HTTP requests and handle errors
make_request() {
    local method="$1"
    local url="$2"
    local data="$3"
    local description="$4"
    
    echo -e "${BLUE}Testing:${NC} $description"
    
    if [[ -n "$data" ]]; then
        response=$(curl -s -X "$method" "$url" -H "Content-Type: application/json" -d "$data" 2>/dev/null)
    else
        response=$(curl -s -X "$method" "$url" 2>/dev/null)
    fi
    
    if [[ $? -ne 0 ]]; then
        echo -e "${RED}✗ FAIL${NC}: Failed to connect to $url"
        TOTAL_TESTS=$((TOTAL_TESTS + 1))
        FAILED_TESTS=$((FAILED_TESTS + 1))
        return 1
    fi
    
    echo "$response"
}

echo -e "${YELLOW}=== Mini-Zanzibar Complete Flow Test ===${NC}"
echo ""

# Test 1: Health Check
echo -e "${BLUE}=== Step 1: Health Check ===${NC}"
response=$(make_request "GET" "$BASE_URL/health" "" "Health check endpoint")
print_test_result "Health check" '"status":"healthy"' "$response"
echo ""

# Test 2: Create Namespace
echo -e "${BLUE}=== Step 2: Create Namespace Configuration ===${NC}"
namespace_data='{
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
}'

response=$(make_request "POST" "$BASE_URL/namespace" "$namespace_data" "Create doc namespace")
print_test_result "Namespace creation" '"message":"Namespace created successfully"' "$response"
echo ""

# Test 3: Verify Namespace
echo -e "${BLUE}=== Step 3: Verify Namespace Configuration ===${NC}"
response=$(make_request "GET" "$BASE_URL/namespace/doc" "" "Get doc namespace")
print_test_result "Namespace retrieval" '"namespace":"doc"' "$response"
print_test_result "Namespace version" '"version":1' "$response"
echo ""

# Test 4: Create ACL Tuples
echo -e "${BLUE}=== Step 4: Create ACL Tuples ===${NC}"

# Alice as owner
alice_owner='{"object": "doc:readme", "relation": "owner", "user": "user:alice"}'
response=$(make_request "POST" "$BASE_URL/acl" "$alice_owner" "Alice as owner of doc:readme")
print_test_result "Alice owner ACL" '"message":"ACL created successfully"' "$response"

# Bob as editor
bob_editor='{"object": "doc:readme", "relation": "editor", "user": "user:bob"}'
response=$(make_request "POST" "$BASE_URL/acl" "$bob_editor" "Bob as editor of doc:readme")
print_test_result "Bob editor ACL" '"message":"ACL created successfully"' "$response"

# Charlie as viewer
charlie_viewer='{"object": "doc:readme", "relation": "viewer", "user": "user:charlie"}'
response=$(make_request "POST" "$BASE_URL/acl" "$charlie_viewer" "Charlie as viewer of doc:readme")
print_test_result "Charlie viewer ACL" '"message":"ACL created successfully"' "$response"
echo ""

# Test 5: Authorization Checks
echo -e "${BLUE}=== Step 5: Authorization Checks ===${NC}"

# Direct access checks
response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=owner&user=user:alice" "" "Alice owner access")
print_test_result "Alice has owner access" '"authorized":true' "$response"

response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=editor&user=user:bob" "" "Bob editor access")
print_test_result "Bob has editor access" '"authorized":true' "$response"

response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=viewer&user=user:charlie" "" "Charlie viewer access")
print_test_result "Charlie has viewer access" '"authorized":true' "$response"

# Computed userset checks (TODO: These will fail until computed usersets are implemented)
response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=viewer&user=user:alice" "" "Alice viewer access (computed)")
print_test_result "Alice has viewer access (computed)" '"authorized":true' "$response"

response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=viewer&user=user:bob" "" "Bob viewer access (computed)")
print_test_result "Bob has viewer access (computed)" '"authorized":true' "$response"

# Negative checks
response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=owner&user=user:bob" "" "Bob owner access (should fail)")
print_test_result "Bob does not have owner access" '"authorized":false' "$response"

response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=editor&user=user:charlie" "" "Charlie editor access (should fail)")
print_test_result "Charlie does not have editor access" '"authorized":false' "$response"

response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=viewer&user=user:unknown" "" "Unknown user access (should fail)")
print_test_result "Unknown user has no access" '"authorized":false' "$response"
echo ""

# Test 6: List ACLs
echo -e "${BLUE}=== Step 6: List ACL Operations ===${NC}"

response=$(make_request "GET" "$BASE_URL/acl/object/doc:readme" "" "List ACLs for doc:readme")
print_test_result "List ACLs by object" '"tuples"' "$response"

response=$(make_request "GET" "$BASE_URL/acl/user/user:alice" "" "List ACLs for user:alice")
print_test_result "List ACLs by user" '"tuples"' "$response"
echo ""

# Test 7: ACL Deletion
echo -e "${BLUE}=== Step 7: ACL Deletion ===${NC}"

delete_data='{"object": "doc:readme", "relation": "editor", "user": "user:bob"}'
response=$(make_request "DELETE" "$BASE_URL/acl" "$delete_data" "Delete Bob's editor access")
print_test_result "ACL deletion" '"message":"ACL deleted successfully"' "$response"

# Verify deletion
response=$(make_request "GET" "$BASE_URL/acl/check?object=doc:readme&relation=editor&user=user:bob" "" "Bob editor access after deletion")
print_test_result "Bob no longer has editor access" '"authorized":false' "$response"
echo ""

# Test 8: Namespace Listing
echo -e "${BLUE}=== Step 8: Namespace Management ===${NC}"

response=$(make_request "GET" "$BASE_URL/namespaces" "" "List all namespaces")
print_test_result "List namespaces" '"namespaces"' "$response"
echo ""

# Test 9: Error Handling
echo -e "${BLUE}=== Step 9: Error Handling ===${NC}"

# Invalid ACL creation
invalid_acl='{"object": "doc:readme"}'
response=$(make_request "POST" "$BASE_URL/acl" "$invalid_acl" "Invalid ACL creation (missing fields)")
print_test_result "Invalid ACL rejection" '"error"' "$response"

# Non-existent namespace
response=$(make_request "GET" "$BASE_URL/namespace/nonexistent" "" "Non-existent namespace")
print_test_result "Non-existent namespace error" '"error"' "$response"
echo ""

# Test Summary
echo -e "${YELLOW}=== Test Summary ===${NC}"
echo -e "Total Tests: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"

if [[ $FAILED_TESTS -eq 0 ]]; then
    echo -e "${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed. Check the output above.${NC}"
    exit 1
fi