# Build stage
FROM rust:1.72 as builder

WORKDIR /app

# Copy manifests
COPY Cargo.toml Cargo.lock ./

# Build dependencies (for better caching)
RUN mkdir src && \
    echo "fn main() {}" > src/main.rs && \
    touch src/lib.rs && \
    cargo build --release && \
    rm -rf src

# Copy source code
COPY src ./src

# Build application
RUN cargo build --release

# Build WASM module
RUN rustup target add wasm32-unknown-unknown
RUN cargo build --target wasm32-unknown-unknown --release

# Final stage
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/target/release/analytics-service /app/
COPY --from=builder /app/target/wasm32-unknown-unknown/release/*.wasm /app/wasm/

EXPOSE 8081

CMD ["./analytics-service"]