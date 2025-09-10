#!/bin/bash

# Database Testing Comparison Demo
# This script demonstrates the performance difference between containerized and existing databases

echo "🚀 Database Testing Comparison Demo"
echo "===================================="
echo ""

echo "📦 Testing with Containerized Database (testcontainers)..."
echo "This will create a fresh PostgreSQL container for each test"
echo ""
time go test -count=1 -v ./... 2>/dev/null | grep -E "(Starting|completed in|PASS|FAIL)"
echo ""

echo "🔄 Now let's test with a Pre-existing Database..."
echo "Using remote PostgreSQL instance"
echo ""
echo "📊 Testing with Pre-existing Database..."
echo "This will reuse the same PostgreSQL instance"
echo ""
time go test -count=1 -v -existing-db ./... 2>/dev/null | grep -E "(Starting|completed in|PASS|FAIL)"

echo ""
echo "✅ Demo completed!"
echo ""
echo "Key Observations:"
echo "- Containerized tests take longer due to container startup time"
echo "- Pre-existing database tests are much faster"
echo "- Containerized tests provide complete isolation"
echo "- Pre-existing TimescaleDB is shared and always available"
