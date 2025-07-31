use axum::{
    extract::{Path, State},
    http::StatusCode,
    response::Json,
    routing::{get, post},
    Router,
};
use serde::{Deserialize, Serialize};
use std::sync::Arc;
use tokio::net::TcpListener;
use tracing::{error, info};

mod analytics;
mod indicators;
mod streaming;

use analytics::AnalyticsEngine;
use streaming::StreamProcessor;

#[derive(Clone)]
struct AppState {
    analytics_engine: Arc<AnalyticsEngine>,
    stream_processor: Arc<StreamProcessor>,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    tracing_subscriber::fmt::init();

    // Initialize components
    let analytics_engine = Arc::new(AnalyticsEngine::new().await?);
    let stream_processor = Arc::new(StreamProcessor::new().await?);
    
    // Start stream processing
    let processor = stream_processor.clone();
    tokio::spawn(async move {
        processor.start().await;
    });

    let app_state = AppState {
        analytics_engine,
        stream_processor,
    };

    // Build router
    let app = Router::new()
        .route("/health", get(health_check))
        .route("/indicators/:symbol", get(get_indicators))
        .route("/calculate", post(calculate_custom))
        .route("/historical/:symbol", get(get_historical))
        .with_state(app_state);

    // Start server
    let listener = TcpListener::bind("0.0.0.0:8081").await?;
    info!("Analytics engine listening on {}", listener.local_addr()?);
    
    axum::serve(listener, app).await?;
    Ok(())
}

async fn health_check() -> StatusCode {
    StatusCode::OK
}

#[derive(Deserialize)]
struct IndicatorQuery {
    indicators: Vec<String>,
    period: Option<usize>,
}

#[derive(Serialize)]
struct IndicatorResponse {
    symbol: String,
    indicators: serde_json::Value,
    timestamp: i64,
}

async fn get_indicators(
    Path(symbol): Path<String>,
    State(state): State<AppState>,
    Json(query): Json<IndicatorQuery>,
) -> Result<Json<IndicatorResponse>, StatusCode> {
    match state.analytics_engine.calculate_indicators(&symbol, &query.indicators, query.period).await {
        Ok(result) => Ok(Json(IndicatorResponse {
            symbol,
            indicators: result,
            timestamp: chrono::Utc::now().timestamp(),
        })),
        Err(e) => {
            error!("Failed to calculate indicators: {}", e);
            Err(StatusCode::INTERNAL_SERVER_ERROR)
        }
    }
}

#[derive(Deserialize)]
struct CalculationRequest {
    symbol: String,
    data: Vec<f64>,
    calculation_type: String,
    params: serde_json::Value,
}

#[derive(Serialize)]
struct CalculationResponse {
    result: serde_json::Value,
    execution_time_ms: u64,
}

async fn calculate_custom(
    State(state): State<AppState>,
    Json(request): Json<CalculationRequest>,
) -> Result<Json<CalculationResponse>, StatusCode> {
    let start = std::time::Instant::now();
    
    match state.analytics_engine.calculate_custom(&request).await {
        Ok(result) => Ok(Json(CalculationResponse {
            result,
            execution_time_ms: start.elapsed().as_millis() as u64,
        })),
        Err(e) => {
            error!("Calculation failed: {}", e);
            Err(StatusCode::INTERNAL_SERVER_ERROR)
        }
    }
}

async fn get_historical(
    Path(symbol): Path<String>,
    State(state): State<AppState>,
) -> Result<Json<serde_json::Value>, StatusCode> {
    match state.analytics_engine.get_historical_data(&symbol).await {
        Ok(data) => Ok(Json(data)),
        Err(e) => {
            error!("Failed to get historical data: {}", e);
            Err(StatusCode::INTERNAL_SERVER_ERROR)
        }
    }
}