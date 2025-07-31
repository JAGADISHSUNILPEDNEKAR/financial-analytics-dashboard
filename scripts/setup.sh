#!/bin/bash
set -e

echo "Setting up Financial Analytics Platform..."

# Check prerequisites
command -v docker >/dev/null 2>&1 || { echo "Docker is required but not installed. Aborting." >&2; exit 1; }
command -v flutter >/dev/null 2>&1 || { echo "Flutter is required but not installed. Aborting." >&2; exit 1; }
command -v go >/dev/null 2>&1 || { echo "Go is required but not installed. Aborting." >&2; exit 1; }
command -v cargo >/dev/null 2>&1 || { echo "Rust/Cargo is required but not installed. Aborting." >&2; exit 1; }

# Create .env file if not exists
if [ ! -f .env ]; then
    echo "Creating .env file..."
    cat > .env << EOF
# Database
POSTGRES_PASSWORD=your_secure_password
MONGO_PASSWORD=your_secure_password
REDIS_PASSWORD=your_secure_password

# JWT
JWT_SECRET=your_jwt_secret_here

# Firebase
FIREBASE_PROJECT_ID=your_project_id
FIREBASE_API_KEY=your_api_key

# Grafana
GRAFANA_PASSWORD=admin

# API Keys (for market data)
ALPHA_VANTAGE_API_KEY=your_key
POLYGON_API_KEY=your_key
EOF
    echo ".env file created. Please update with your actual values."
fi

# Install Flutter dependencies
echo "Installing Flutter dependencies..."
cd frontend
flutter pub get
cd ..

# Install Go dependencies
echo "Installing Go dependencies..."
for service in api-gateway auth-service user-service dashboard-service; do
    if [ -d "backend/services/$service" ]; then
        echo "Installing dependencies for $service..."
        cd "backend/services/$service"
        go mod download
        cd ../../..
    fi
done

# Install Rust dependencies
echo "Installing Rust dependencies..."
cd backend/analytics-engine
cargo fetch
cd ../..

# Setup database
echo "Starting databases..."
docker-compose up -d postgres mongodb redis

# Wait for databases to be ready
echo "Waiting for databases to be ready..."
sleep 10

# Run migrations
echo "Running database migrations..."
docker-compose exec -T postgres psql -U postgres -d financial_analytics < database/migrations/001_initial_schema.sql

echo "Setup complete! You can now run './scripts/start-dev.sh' to start the development environment."