package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	db    *sql.DB
	redis *redis.Client
}

type UserProfile struct {
	ID          string                 `json:"id"`
	Email       string                 `json:"email"`
	DisplayName string                 `json:"display_name"`
	AvatarURL   string                 `json:"avatar_url"`
	Preferences map[string]interface{} `json:"preferences"`
	CreatedAt   string                 `json:"created_at"`
}

func main() {
	// Initialize database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Redis
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}
	redisClient := redis.NewClient(opt)

	service := &UserService{
		db:    db,
		redis: redisClient,
	}

	// Setup routes
	router := mux.NewRouter()

	// User profile routes
	router.HandleFunc("/users/{id}", service.getUser).Methods("GET")
	router.HandleFunc("/users/{id}", service.updateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", service.deleteUser).Methods("DELETE")

	// Preferences routes
	router.HandleFunc("/users/{id}/preferences", service.getPreferences).Methods("GET")
	router.HandleFunc("/users/{id}/preferences", service.updatePreferences).Methods("PUT")

	// Settings routes
	router.HandleFunc("/users/{id}/settings", service.getSettings).Methods("GET")
	router.HandleFunc("/users/{id}/settings", service.updateSettings).Methods("PUT")

	// Avatar upload
	router.HandleFunc("/users/{id}/avatar", service.uploadAvatar).Methods("POST")

	// Health check
	router.HandleFunc("/health", handleHealth).Methods("GET")

	log.Println("User service listening on :8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}

func (s *UserService) getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Check cache first
	ctx := r.Context()
	cacheKey := "user:" + userID
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(cached)); err != nil {
			log.Println("Failed to write cached response:", err)
		}
		return
	}

	// Get from database
	var profile UserProfile
	var prefsJSON []byte

	err = s.db.QueryRow(`
        SELECT u.id, u.email, u.display_name, u.avatar_url, u.created_at,
               p.settings
        FROM users u
        LEFT JOIN user_preferences p ON p.user_id = u.id
        WHERE u.id = $1
    `, userID).Scan(
		&profile.ID,
		&profile.Email,
		&profile.DisplayName,
		&profile.AvatarURL,
		&profile.CreatedAt,
		&prefsJSON,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Parse preferences
	if prefsJSON != nil {
		if err := json.Unmarshal(prefsJSON, &profile.Preferences); err != nil {
			log.Println("Failed to unmarshal preferences:", err)
		}
	}

	// Cache the result
	responseData, err := json.Marshal(profile)
	if err != nil {
		log.Println("Failed to marshal profile:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	s.redis.Set(ctx, cacheKey, responseData, 300*time.Second)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseData); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *UserService) updateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var update struct {
		DisplayName string `json:"display_name"`
		AvatarURL   string `json:"avatar_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Update database
	_, err := s.db.Exec(`
        UPDATE users 
        SET display_name = $1, avatar_url = $2, updated_at = NOW()
        WHERE id = $3
    `, update.DisplayName, update.AvatarURL, userID)

	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "user:"+userID)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"}); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *UserService) getPreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var prefs struct {
		Theme                string          `json:"theme"`
		Timezone             string          `json:"timezone"`
		NotificationsEnabled bool            `json:"notifications_enabled"`
		DefaultDashboardID   *string         `json:"default_dashboard_id"`
		Settings             json.RawMessage `json:"settings"`
	}

	err := s.db.QueryRow(`
        SELECT theme, timezone, notifications_enabled, default_dashboard_id, settings
        FROM user_preferences
        WHERE user_id = $1
    `, userID).Scan(
		&prefs.Theme,
		&prefs.Timezone,
		&prefs.NotificationsEnabled,
		&prefs.DefaultDashboardID,
		&prefs.Settings,
	)

	if err == sql.ErrNoRows {
		// Create default preferences
		_, err = s.db.Exec(`
            INSERT INTO user_preferences (user_id) VALUES ($1)
        `, userID)
		if err != nil {
			http.Error(w, "Failed to create preferences", http.StatusInternalServerError)
			return
		}

		// Return defaults
		prefs.Theme = "light"
		prefs.Timezone = "UTC"
		prefs.NotificationsEnabled = true
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(prefs); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *UserService) updatePreferences(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var prefs struct {
		Theme                string          `json:"theme"`
		Timezone             string          `json:"timezone"`
		NotificationsEnabled bool            `json:"notifications_enabled"`
		DefaultDashboardID   *string         `json:"default_dashboard_id"`
		Settings             json.RawMessage `json:"settings"`
	}

	if err := json.NewDecoder(r.Body).Decode(&prefs); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
        UPDATE user_preferences
        SET theme = $1, timezone = $2, notifications_enabled = $3, 
            default_dashboard_id = $4, settings = $5, updated_at = NOW()
        WHERE user_id = $6
    `, prefs.Theme, prefs.Timezone, prefs.NotificationsEnabled,
		prefs.DefaultDashboardID, prefs.Settings, userID)

	if err != nil {
		http.Error(w, "Failed to update preferences", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "user:"+userID)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Preferences updated successfully"}); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *UserService) deleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Soft delete
	_, err := s.db.Exec(`
        UPDATE users 
        SET deleted_at = NOW()
        WHERE id = $1
    `, userID)

	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "user:"+userID)

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"}); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *UserService) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Missing avatar file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// In production, upload to S3 or similar
	// For now, we'll just save the filename
	avatarURL := "/avatars/" + userID + "_" + header.Filename

	_, err = s.db.Exec(`
        UPDATE users 
        SET avatar_url = $1, updated_at = NOW()
        WHERE id = $2
    `, avatarURL, userID)

	if err != nil {
		http.Error(w, "Failed to update avatar", http.StatusInternalServerError)
		return
	}

	// Invalidate cache
	ctx := r.Context()
	s.redis.Del(ctx, "user:"+userID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"avatar_url": avatarURL,
		"message":    "Avatar uploaded successfully",
	}); err != nil {
		log.Println("Failed to write response:", err)
	}
}

func (s *UserService) getSettings(_ http.ResponseWriter, _ *http.Request) {
	// Implementation similar to getPreferences but for app-specific settings
}

func (s *UserService) updateSettings(_ http.ResponseWriter, _ *http.Request) {
	// Implementation similar to updatePreferences but for app-specific settings
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "healthy"}); err != nil {
		log.Println("Failed to write response:", err)
	}
}
