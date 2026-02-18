package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type DashboardService struct {
	db    *sql.DB
	redis *redis.Client
	kafka *kafka.Writer
}

type Dashboard struct {
	ID          string          `json:"id"`
	UserID      string          `json:"user_id"`
	Name        string          `json:"name"`
	Layout      json.RawMessage `json:"layout"`
	IsPublic    bool            `json:"is_public"`
	Widgets     []Widget        `json:"widgets"`
	Permissions []Permission    `json:"permissions"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Widget struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Config   json.RawMessage `json:"config"`
	Position json.RawMessage `json:"position"`
}

type Permission struct {
	UserID     string `json:"user_id"`
	Permission string `json:"permission"`
}

func main() {
	// Initialize database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Redis
	opt, _ := redis.ParseURL(os.Getenv("REDIS_URL"))
	redisClient := redis.NewClient(opt)

	// Initialize Kafka writer
	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   "dashboard-events",
	})

	service := &DashboardService{
		db:    db,
		redis: redisClient,
		kafka: kafkaWriter,
	}

	// Setup routes
	router := mux.NewRouter()

	// Dashboard routes
	router.HandleFunc("/dashboards", service.listDashboards).Methods("GET")
	router.HandleFunc("/dashboards", service.createDashboard).Methods("POST")
	router.HandleFunc("/dashboards/{id}", service.getDashboard).Methods("GET")
	router.HandleFunc("/dashboards/{id}", service.updateDashboard).Methods("PUT")
	router.HandleFunc("/dashboards/{id}", service.deleteDashboard).Methods("DELETE")

	// Widget routes
	router.HandleFunc("/dashboards/{id}/widgets", service.addWidget).Methods("POST")
	router.HandleFunc("/dashboards/{id}/widgets/{widgetId}", service.updateWidget).Methods("PUT")
	router.HandleFunc("/dashboards/{id}/widgets/{widgetId}", service.deleteWidget).Methods("DELETE")

	// Sharing routes
	router.HandleFunc("/dashboards/{id}/share", service.shareDashboard).Methods("POST")
	router.HandleFunc("/dashboards/{id}/permissions", service.getPermissions).Methods("GET")
	router.HandleFunc("/dashboards/{id}/permissions", service.updatePermissions).Methods("PUT")

	// Public dashboards
	router.HandleFunc("/public/dashboards", service.listPublicDashboards).Methods("GET")

	// Health check
	router.HandleFunc("/health", handleHealth).Methods("GET")

	log.Println("Dashboard service listening on :8084")
	log.Fatal(http.ListenAndServe(":8084", router))
}

func (s *DashboardService) listDashboards(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	rows, err := s.db.Query(`
        SELECT d.id, d.user_id, d.name, d.layout, d.is_public, d.created_at, d.updated_at
        FROM dashboards d
        LEFT JOIN dashboard_permissions dp ON d.id = dp.dashboard_id
        WHERE d.user_id = $1 OR dp.user_id = $1 OR d.is_public = true
        ORDER BY d.updated_at DESC
    `, userID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dashboards []Dashboard
	for rows.Next() {
		var d Dashboard
		err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.Layout, &d.IsPublic, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			continue
		}

		// Load widgets for each dashboard
		s.loadWidgets(&d)
		dashboards = append(dashboards, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboards)
}

func (s *DashboardService) createDashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var req struct {
		Name     string          `json:"name"`
		Layout   json.RawMessage `json:"layout"`
		IsPublic bool            `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var dashboardID string
	err := s.db.QueryRow(`
        INSERT INTO dashboards (user_id, name, layout, is_public)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, userID, req.Name, req.Layout, req.IsPublic).Scan(&dashboardID)

	if err != nil {
		http.Error(w, "Failed to create dashboard", http.StatusInternalServerError)
		return
	}

	// Publish event
	s.publishEvent("dashboard.created", map[string]interface{}{
		"dashboard_id": dashboardID,
		"user_id":      userID,
		"name":         req.Name,
	})

	dashboard := Dashboard{
		ID:        dashboardID,
		UserID:    userID,
		Name:      req.Name,
		Layout:    req.Layout,
		IsPublic:  req.IsPublic,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dashboard)
}

func (s *DashboardService) getDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Check cache first
	ctx := r.Context()
	cacheKey := "dashboard:" + dashboardID
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	var dashboard Dashboard
	err = s.db.QueryRow(`
        SELECT id, user_id, name, layout, is_public, created_at, updated_at
        FROM dashboards
        WHERE id = $1
    `, dashboardID).Scan(
		&dashboard.ID,
		&dashboard.UserID,
		&dashboard.Name,
		&dashboard.Layout,
		&dashboard.IsPublic,
		&dashboard.CreatedAt,
		&dashboard.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "Dashboard not found", http.StatusNotFound)
		return
	}

	// Check permissions
	if !dashboard.IsPublic && dashboard.UserID != userID {
		hasPermission := s.checkPermission(dashboardID, userID)
		if !hasPermission {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
	}

	// Load widgets
	s.loadWidgets(&dashboard)

	// Cache the result
	responseData, _ := json.Marshal(dashboard)
	s.redis.Set(ctx, cacheKey, responseData, 300*time.Second)

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseData)
}

func (s *DashboardService) updateDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Check ownership
	var ownerID string
	s.db.QueryRow("SELECT user_id FROM dashboards WHERE id = $1", dashboardID).Scan(&ownerID)
	if ownerID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req struct {
		Name     string          `json:"name"`
		Layout   json.RawMessage `json:"layout"`
		IsPublic bool            `json:"is_public"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
        UPDATE dashboards
        SET name = $1, layout = $2, is_public = $3, updated_at = NOW()
        WHERE id = $4
    `, req.Name, req.Layout, req.IsPublic, dashboardID)

	if err != nil {
		http.Error(w, "Failed to update dashboard", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "dashboard:"+dashboardID)

	// Publish event
	s.publishEvent("dashboard.updated", map[string]interface{}{
		"dashboard_id": dashboardID,
		"user_id":      userID,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dashboard updated successfully"})
}

func (s *DashboardService) deleteDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Check ownership
	var ownerID string
	s.db.QueryRow("SELECT user_id FROM dashboards WHERE id = $1", dashboardID).Scan(&ownerID)
	if ownerID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	_, err := s.db.Exec("DELETE FROM dashboards WHERE id = $1", dashboardID)
	if err != nil {
		http.Error(w, "Failed to delete dashboard", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "dashboard:"+dashboardID)

	// Publish event
	s.publishEvent("dashboard.deleted", map[string]interface{}{
		"dashboard_id": dashboardID,
		"user_id":      userID,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dashboard deleted successfully"})
}

func (s *DashboardService) addWidget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Check ownership
	var ownerID string
	s.db.QueryRow("SELECT user_id FROM dashboards WHERE id = $1", dashboardID).Scan(&ownerID)
	if ownerID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var widget Widget
	if err := json.NewDecoder(r.Body).Decode(&widget); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var widgetID string
	err := s.db.QueryRow(`
        INSERT INTO dashboard_widgets (dashboard_id, widget_type, config, position)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `, dashboardID, widget.Type, widget.Config, widget.Position).Scan(&widgetID)

	if err != nil {
		http.Error(w, "Failed to add widget", http.StatusInternalServerError)
		return
	}

	widget.ID = widgetID

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "dashboard:"+dashboardID)

	// Publish event
	s.publishEvent("widget.added", map[string]interface{}{
		"dashboard_id": dashboardID,
		"widget_id":    widgetID,
		"widget_type":  widget.Type,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(widget)
}

func (s *DashboardService) updateWidget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]
	widgetID := vars["widgetId"]
	userID := r.Header.Get("X-User-ID")

	// Check ownership
	var ownerID string
	s.db.QueryRow("SELECT user_id FROM dashboards WHERE id = $1", dashboardID).Scan(&ownerID)
	if ownerID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var widget Widget
	if err := json.NewDecoder(r.Body).Decode(&widget); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
        UPDATE dashboard_widgets
        SET config = $1, position = $2, updated_at = NOW()
        WHERE id = $3 AND dashboard_id = $4
    `, widget.Config, widget.Position, widgetID, dashboardID)

	if err != nil {
		http.Error(w, "Failed to update widget", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "dashboard:"+dashboardID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Widget updated successfully"})
}

func (s *DashboardService) deleteWidget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]
	widgetID := vars["widgetId"]
	userID := r.Header.Get("X-User-ID")

	// Check ownership
	var ownerID string
	s.db.QueryRow("SELECT user_id FROM dashboards WHERE id = $1", dashboardID).Scan(&ownerID)
	if ownerID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	_, err := s.db.Exec(`
        DELETE FROM dashboard_widgets
        WHERE id = $1 AND dashboard_id = $2
    `, widgetID, dashboardID)

	if err != nil {
		http.Error(w, "Failed to delete widget", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "dashboard:"+dashboardID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Widget deleted successfully"})
}

func (s *DashboardService) shareDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	// Check ownership
	var ownerID string
	s.db.QueryRow("SELECT user_id FROM dashboards WHERE id = $1", dashboardID).Scan(&ownerID)
	if ownerID != userID {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var req struct {
		UserIDs    []string `json:"user_ids"`
		Permission string   `json:"permission"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Add permissions
	for _, sharedUserID := range req.UserIDs {
		_, err := s.db.Exec(`
            INSERT INTO dashboard_permissions (dashboard_id, user_id, permission_type)
            VALUES ($1, $2, $3)
            ON CONFLICT (dashboard_id, user_id) 
            DO UPDATE SET permission_type = $3
        `, dashboardID, sharedUserID, req.Permission)

		if err != nil {
			log.Printf("Failed to share with user %s: %v", sharedUserID, err)
		}
	}

	// Publish event
	s.publishEvent("dashboard.shared", map[string]interface{}{
		"dashboard_id": dashboardID,
		"shared_with":  req.UserIDs,
		"permission":   req.Permission,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Dashboard shared successfully"})
}

func (s *DashboardService) loadWidgets(dashboard *Dashboard) {
	rows, err := s.db.Query(`
        SELECT id, widget_type, config, position
        FROM dashboard_widgets
        WHERE dashboard_id = $1
        ORDER BY created_at
    `, dashboard.ID)

	if err != nil {
		return
	}
	defer rows.Close()

	dashboard.Widgets = []Widget{}
	for rows.Next() {
		var w Widget
		if err := rows.Scan(&w.ID, &w.Type, &w.Config, &w.Position); err == nil {
			dashboard.Widgets = append(dashboard.Widgets, w)
		}
	}
}

func (s *DashboardService) checkPermission(dashboardID, userID string) bool {
	var exists bool
	s.db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM dashboard_permissions
            WHERE dashboard_id = $1 AND user_id = $2
        )
    `, dashboardID, userID).Scan(&exists)
	return exists
}

func (s *DashboardService) publishEvent(eventType string, data map[string]interface{}) {
	event := map[string]interface{}{
		"type":      eventType,
		"timestamp": time.Now().Unix(),
		"data":      data,
	}

	message, _ := json.Marshal(event)

	ctx := context.Background()
	err := s.kafka.WriteMessages(ctx, kafka.Message{
		Value: message,
	})

	if err != nil {
		log.Printf("Failed to publish event: %v", err)
	}
}

func (s *DashboardService) getPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dashboardID := vars["id"]

	rows, err := s.db.Query(`
        SELECT dp.user_id, dp.permission_type, u.email
        FROM dashboard_permissions dp
        JOIN users u ON u.id = dp.user_id
        WHERE dp.dashboard_id = $1
    `, dashboardID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var permissions []map[string]string
	for rows.Next() {
		var userID, permission, email string
		if err := rows.Scan(&userID, &permission, &email); err == nil {
			permissions = append(permissions, map[string]string{
				"user_id":    userID,
				"email":      email,
				"permission": permission,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
}

func (s *DashboardService) updatePermissions(_ http.ResponseWriter, _ *http.Request) {
	// Implementation similar to shareDashboard
}

func (s *DashboardService) listPublicDashboards(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(`
        SELECT d.id, d.user_id, d.name, d.layout, d.created_at, d.updated_at, u.email
        FROM dashboards d
        JOIN users u ON u.id = d.user_id
        WHERE d.is_public = true
        ORDER BY d.updated_at DESC
        LIMIT 50
    `)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var dashboards []map[string]interface{}
	for rows.Next() {
		var d Dashboard
		var userEmail string
		err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.Layout, &d.CreatedAt, &d.UpdatedAt, &userEmail)
		if err != nil {
			continue
		}

		dashboards = append(dashboards, map[string]interface{}{
			"id":         d.ID,
			"name":       d.Name,
			"owner":      userEmail,
			"created_at": d.CreatedAt,
			"updated_at": d.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dashboards)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
