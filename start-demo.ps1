# Start all services for the Mini-Zanzibar Demo

Write-Host "🚀 Starting Mini-Zanzibar Demo Services..." -ForegroundColor Green

# Check if Mini-Zanzibar is running
Write-Host "1. Checking Mini-Zanzibar service..." -ForegroundColor Cyan
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -TimeoutSec 5
    Write-Host "✅ Mini-Zanzibar is already running" -ForegroundColor Green
} catch {
    Write-Host "❌ Mini-Zanzibar is not running. Please start it first:" -ForegroundColor Red
    Write-Host "   cd mini-zanzibar && go run cmd/server/main.go" -ForegroundColor Yellow
    exit 1
}

# Start Auth Service
Write-Host "`n2. Starting Auth Service..." -ForegroundColor Cyan
$authJob = Start-Job -ScriptBlock {
    Set-Location "c:\Users\mniko\source\repos\RBS\RBS-TEAM-10\auth-service"
    go run main.go
}

# Wait a bit for auth service to start
Start-Sleep -Seconds 3

# Check if Auth Service is running
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8081/health" -TimeoutSec 5
    Write-Host "✅ Auth Service is running on port 8081" -ForegroundColor Green
} catch {
    Write-Host "❌ Failed to start Auth Service" -ForegroundColor Red
    Stop-Job $authJob
    exit 1
}

# Start Web Client (simple HTTP server)
Write-Host "`n3. Starting Web Client..." -ForegroundColor Cyan

# Check if Python is available for simple HTTP server
try {
    python --version | Out-Null
    $webJob = Start-Job -ScriptBlock {
        Set-Location "c:\Users\mniko\source\repos\RBS\RBS-TEAM-10\web-client"
        python -m http.server 3000
    }
    Write-Host "✅ Web Client is running on port 3000 (Python HTTP server)" -ForegroundColor Green
} catch {
    try {
        # Try Node.js if Python is not available
        npx --version | Out-Null
        $webJob = Start-Job -ScriptBlock {
            Set-Location "c:\Users\mniko\source\repos\RBS\RBS-TEAM-10\web-client"
            npx http-server -p 3000 -c-1
        }
        Write-Host "✅ Web Client is running on port 3000 (Node.js HTTP server)" -ForegroundColor Green
    } catch {
        Write-Host "⚠️  Could not start HTTP server automatically." -ForegroundColor Yellow
        Write-Host "   Please serve the web-client folder manually on port 3000" -ForegroundColor Yellow
    }
}

Write-Host "`n🎉 Demo Environment Setup Complete!" -ForegroundColor Green
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

Write-Host "`n📋 Service Status:" -ForegroundColor White
Write-Host "   🔧 Mini-Zanzibar (Authorization):  http://localhost:8080" -ForegroundColor Cyan
Write-Host "   🔐 Auth Service (Authentication):   http://localhost:8081" -ForegroundColor Cyan  
Write-Host "   🌐 Web Client (Frontend):           http://localhost:3000" -ForegroundColor Cyan

Write-Host "`n👤 Demo Users:" -ForegroundColor White
Write-Host "   alice / alice123    (Admin - Full Access)" -ForegroundColor Green
Write-Host "   bob / bob123        (Editor - Edit Documents)" -ForegroundColor Yellow  
Write-Host "   charlie / charlie123 (Viewer - View Only)" -ForegroundColor Blue
Write-Host "   david / david123    (User - Limited Access)" -ForegroundColor Red

Write-Host "`n🔗 Architecture Flow:" -ForegroundColor White
Write-Host "   Web Client → Auth Service → Mini-Zanzibar" -ForegroundColor Gray

Write-Host "`n📖 Getting Started:" -ForegroundColor White
Write-Host "   1. Open http://localhost:3000 in your browser" -ForegroundColor Gray
Write-Host "   2. Login with any demo user (e.g., alice / alice123)" -ForegroundColor Gray
Write-Host "   3. Explore documents and authorization features" -ForegroundColor Gray

Write-Host "`n💡 Press Ctrl+C to stop all services" -ForegroundColor Yellow
Write-Host "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━" -ForegroundColor Gray

# Keep the script running and monitor services
try {
    while ($true) {
        Start-Sleep -Seconds 5
        
        # Check if services are still running
        try {
            Invoke-RestMethod -Uri "http://localhost:8081/health" -TimeoutSec 2 | Out-Null
        } catch {
            Write-Host "❌ Auth Service stopped unexpectedly" -ForegroundColor Red
            break
        }
        
        # Check job status
        if ($authJob.State -eq "Failed" -or $authJob.State -eq "Stopped") {
            Write-Host "❌ Auth Service job failed" -ForegroundColor Red
            break
        }
    }
} finally {
    Write-Host "`n🛑 Stopping services..." -ForegroundColor Yellow
    
    if ($authJob) {
        Stop-Job $authJob
        Remove-Job $authJob
    }
    
    if ($webJob) {
        Stop-Job $webJob  
        Remove-Job $webJob
    }
    
    Write-Host "✅ All services stopped" -ForegroundColor Green
}