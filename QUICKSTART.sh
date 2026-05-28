#!/bin/bash

# Vinhomes Property Management - Quick Start Script
# This script helps set up and run the frontend

echo "🏢 Vinhomes Property Management System - Quick Start"
echo "=================================================="
echo ""

# Check if backend is running
echo "📡 Checking backend connection..."
if curl -s http://localhost:8080/health | grep -q "healthy"; then
  echo "✅ Backend is running on http://localhost:8080"
else
  echo "⚠️  Backend doesn't appear to be running"
  echo "   Please start the backend first:"
  echo "   cd backend && go run ./cmd/main.go"
  echo ""
fi

# Check if we have a server running
echo ""
echo "🌐 Frontend Server Options:"
echo ""
echo "Option 1: Python 3 (recommended)"
echo "  python3 -m http.server 8000"
echo ""
echo "Option 2: Node.js / npm"
echo "  npx http-server -p 8000"
echo ""
echo "Option 3: PHP"
echo "  php -S localhost:8000"
echo ""
echo "After starting the server, open:"
echo "  http://localhost:8000/pages/login.html"
echo ""
echo "Demo Credentials:"
echo "  Admin:    admin@vinhomes.com / admin123"
echo "  Manager:  manager@vinhomes.com / manager123"
echo "  Staff:    staff@vinhomes.com / staff123"
echo "  Resident: resident@vinhomes.com / resident123"
echo ""
echo "=================================================="
