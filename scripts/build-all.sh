#!/bin/bash
set -e

VERSION=${1:-latest}
REGISTRY=${2:-financial-analytics}

echo "Building all services with version: $VERSION"

# Build Go services
for service in api-gateway auth-service user-service dashboard-service; do
    if [ -d "backend/services/$service" ]; then
        echo "Building $service..."
        docker build -t $REGISTRY/$service:$VERSION backend/services/$service
    fi
done

# Build Analytics Engine
echo "Building analytics-engine..."
docker build -t $REGISTRY/analytics-engine:$VERSION backend/analytics-engine

# Build ML services
for service in anomaly-detection prediction-service; do
    if [ -d "backend/ml-services/$service" ]; then
        echo "Building $service..."
        docker build -t $REGISTRY/$service:$VERSION backend/ml-services/$service
    fi
done

# Build Flutter for all platforms
echo "Building Flutter applications..."
cd frontend

# Web
echo "Building Flutter web..."
flutter build web --release

# Desktop builds (if on appropriate platform)
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    echo "Building Flutter Linux..."
    flutter build linux --release
elif [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Building Flutter macOS..."
    flutter build macos --release
elif [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
    echo "Building Flutter Windows..."
    flutter build windows --release
fi

cd ..

echo "Build complete!"