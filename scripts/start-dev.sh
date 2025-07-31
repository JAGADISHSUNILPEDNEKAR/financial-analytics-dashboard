#!/bin/bash
set -e

echo "Starting Financial Analytics Platform in development mode..."

# Start infrastructure services
echo "Starting infrastructure services..."
docker-compose up -d postgres mongodb redis kafka zookeeper

# Wait for services to be ready
echo "Waiting for services to be ready..."
sleep 15

# Start backend services in background
echo "Starting backend services..."

# API Gateway
cd backend/api-gateway
go run cmd/main.go &
API_GATEWAY_PID=$!
cd ../..

# Analytics Engine
cd backend/analytics-engine
cargo run &
ANALYTICS_PID=$!
cd ../..

echo "Backend services started:"
echo "  API Gateway PID: $API_GATEWAY_PID"
echo "  Analytics Engine PID: $ANALYTICS_PID"

# Create stop script
cat > stop-dev.sh << EOF
#!/bin/bash
echo "Stopping development services..."
kill $API_GATEWAY_PID 2>/dev/null || true
kill $ANALYTICS_PID 2>/dev/null || true
docker-compose down
echo "Services stopped."
EOF
chmod +x stop-dev.sh

echo ""
echo "Development environment is running!"
echo "API Gateway: http://localhost:8080"
echo "Analytics Engine: http://localhost:8081"
echo ""
echo "To start the Flutter app, run:"
echo "  cd frontend && flutter run"
echo ""
echo "To stop all services, run:"
echo "  ./stop-dev.sh"