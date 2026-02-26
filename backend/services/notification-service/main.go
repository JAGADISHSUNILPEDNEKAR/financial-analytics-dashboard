package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type NotificationService struct {
	db    *sql.DB
	fcm   *messaging.Client
	kafka *kafka.Reader
}

type Notification struct {
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"`
	Type      string                 `json:"type"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Data      map[string]interface{} `json:"data"`
	Read      bool                   `json:"read"`
	CreatedAt string                 `json:"created_at"`
}

func main() {
	// Initialize database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize Firebase
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatal("Failed to initialize Firebase:", err)
	}

	fcmClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatal("Failed to get FCM client:", err)
	}

	// Initialize Kafka reader
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{os.Getenv("KAFKA_BROKERS")},
		Topic:   "notifications",
		GroupID: "notification-service",
	})

	service := &NotificationService{
		db:    db,
		fcm:   fcmClient,
		kafka: kafkaReader,
	}

	// Start Kafka consumer
	go service.consumeNotifications()

	// Setup HTTP routes
	router := mux.NewRouter()

	// Notification routes
	router.HandleFunc("/notifications", service.getNotifications).Methods("GET")
	router.HandleFunc("/notifications/{id}/read", service.markAsRead).Methods("PUT")
	router.HandleFunc("/notifications/read-all", service.markAllAsRead).Methods("PUT")
	router.HandleFunc("/notifications/settings", service.getSettings).Methods("GET")
	router.HandleFunc("/notifications/settings", service.updateSettings).Methods("PUT")

	// FCM token management
	router.HandleFunc("/notifications/token", service.registerToken).Methods("POST")
	router.HandleFunc("/notifications/token", service.unregisterToken).Methods("DELETE")

	// Test notification
	router.HandleFunc("/notifications/test", service.sendTestNotification).Methods("POST")

	// Health check
	router.HandleFunc("/health", handleHealth).Methods("GET")

	log.Println("Notification service listening on :8085")
	log.Fatal(http.ListenAndServe(":8085", router))
}

func (s *NotificationService) consumeNotifications() {
	for {
		msg, err := s.kafka.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		var event map[string]interface{}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		s.processNotificationEvent(event)
	}
}

func (s *NotificationService) processNotificationEvent(event map[string]interface{}) {
	eventType := event["type"].(string)
	data := event["data"].(map[string]interface{})

	switch eventType {
	case "alert.triggered":
		s.sendAlertNotification(data)
	case "dashboard.shared":
		s.sendShareNotification(data)
	case "price.threshold":
		s.sendPriceNotification(data)
	default:
		log.Printf("Unknown event type: %s", eventType)
	}
}

func (s *NotificationService) sendAlertNotification(data map[string]interface{}) {
	userID := data["user_id"].(string)
	symbol := data["symbol"].(string)
	condition := data["condition"].(string)

	notification := &messaging.Notification{
		Title: "Alert Triggered",
		Body:  fmt.Sprintf("%s alert triggered for %s", condition, symbol),
	}

	s.sendPushNotification(userID, notification, data)
	s.saveNotification(userID, "alert", notification.Title, notification.Body, data)
}

func (s *NotificationService) sendShareNotification(data map[string]interface{}) {
	sharedWith := data["shared_with"].([]interface{})
	dashboardName := data["dashboard_name"].(string)
	sharedBy := data["shared_by"].(string)

	notification := &messaging.Notification{
		Title: "Dashboard Shared",
		Body:  fmt.Sprintf("%s shared '%s' dashboard with you", sharedBy, dashboardName),
	}

	for _, userID := range sharedWith {
		s.sendPushNotification(userID.(string), notification, data)
		s.saveNotification(userID.(string), "share", notification.Title, notification.Body, data)
	}
}

func (s *NotificationService) sendPriceNotification(data map[string]interface{}) {
	userID := data["user_id"].(string)
	symbol := data["symbol"].(string)
	price := data["price"].(float64)
	threshold := data["threshold"].(float64)

	notification := &messaging.Notification{
		Title: "Price Alert",
		Body:  fmt.Sprintf("%s reached $%.2f (threshold: $%.2f)", symbol, price, threshold),
	}

	s.sendPushNotification(userID, notification, data)
	s.saveNotification(userID, "price", notification.Title, notification.Body, data)
}

func (s *NotificationService) sendPushNotification(userID string, notification *messaging.Notification, data map[string]interface{}) {
	// Get user's FCM tokens
	rows, err := s.db.Query(`
        SELECT token FROM fcm_tokens
        WHERE user_id = $1 AND active = true
    `, userID)

	if err != nil {
		log.Printf("Failed to get FCM tokens: %v", err)
		return
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err == nil {
			tokens = append(tokens, token)
		}
	}

	if len(tokens) == 0 {
		return
	}

	// Send multicast message
	message := &messaging.MulticastMessage{
		Notification: notification,
		Data: map[string]string{
			"type": data["type"].(string),
			"data": string(mustMarshal(data)),
		},
		Tokens: tokens,
	}

	response, err := s.fcm.SendMulticast(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send FCM message: %v", err)
		return
	}

	// Handle failed tokens
	if response.FailureCount > 0 {
		for i, result := range response.Responses {
			if !result.Success {
				// Mark token as inactive
				s.db.Exec(`
                    UPDATE fcm_tokens
                    SET active = false
                    WHERE token = $1
                `, tokens[i])
			}
		}
	}
}

func (s *NotificationService) saveNotification(userID, notifType, title, body string, data map[string]interface{}) {
	_, err := s.db.Exec(`
        INSERT INTO notifications (user_id, type, title, body, data)
        VALUES ($1, $2, $3, $4, $5)
    `, userID, notifType, title, body, mustMarshal(data))

	if err != nil {
		log.Printf("Failed to save notification: %v", err)
	}
}

func (s *NotificationService) getNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	rows, err := s.db.Query(`
        SELECT id, type, title, body, data, read, created_at
        FROM notifications
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT 50
    `, userID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		var dataJSON []byte
		err := rows.Scan(&n.ID, &n.Type, &n.Title, &n.Body, &dataJSON, &n.Read, &n.CreatedAt)
		if err != nil {
			continue
		}

		json.Unmarshal(dataJSON, &n.Data)
		n.UserID = userID
		notifications = append(notifications, n)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func (s *NotificationService) markAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notifID := vars["id"]
	userID := r.Header.Get("X-User-ID")

	_, err := s.db.Exec(`
        UPDATE notifications
        SET read = true, read_at = NOW()
        WHERE id = $1 AND user_id = $2
    `, notifID, userID)

	if err != nil {
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification marked as read"})
}

func (s *NotificationService) markAllAsRead(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	_, err := s.db.Exec(`
        UPDATE notifications
        SET read = true, read_at = NOW()
        WHERE user_id = $1 AND read = false
    `, userID)

	if err != nil {
		http.Error(w, "Failed to update notifications", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "All notifications marked as read"})
}

func (s *NotificationService) registerToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var req struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
        INSERT INTO fcm_tokens (user_id, token, platform)
        VALUES ($1, $2, $3)
        ON CONFLICT (token) 
        DO UPDATE SET user_id = $1, active = true, updated_at = NOW()
    `, userID, req.Token, req.Platform)

	if err != nil {
		http.Error(w, "Failed to register token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Token registered successfully"})
}

func (s *NotificationService) unregisterToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
        UPDATE fcm_tokens
        SET active = false
        WHERE user_id = $1 AND token = $2
    `, userID, req.Token)

	if err != nil {
		http.Error(w, "Failed to unregister token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Token unregistered successfully"})
}

func (s *NotificationService) getSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var settings struct {
		EmailEnabled  bool `json:"email_enabled"`
		PushEnabled   bool `json:"push_enabled"`
		AlertsEnabled bool `json:"alerts_enabled"`
		NewsEnabled   bool `json:"news_enabled"`
	}

	err := s.db.QueryRow(`
        SELECT email_notifications, push_notifications, price_alerts, news_alerts
        FROM notification_settings
        WHERE user_id = $1
    `, userID).Scan(
		&settings.EmailEnabled,
		&settings.PushEnabled,
		&settings.AlertsEnabled,
		&settings.NewsEnabled,
	)

	if err == sql.ErrNoRows {
		// Create default settings
		s.db.Exec(`
            INSERT INTO notification_settings (user_id)
            VALUES ($1)
        `, userID)

		settings = struct {
			EmailEnabled  bool `json:"email_enabled"`
			PushEnabled   bool `json:"push_enabled"`
			AlertsEnabled bool `json:"alerts_enabled"`
			NewsEnabled   bool `json:"news_enabled"`
		}{true, true, true, true}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (s *NotificationService) updateSettings(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var settings struct {
		EmailEnabled  bool `json:"email_enabled"`
		PushEnabled   bool `json:"push_enabled"`
		AlertsEnabled bool `json:"alerts_enabled"`
		NewsEnabled   bool `json:"news_enabled"`
	}

	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
        UPDATE notification_settings
        SET email_notifications = $1, push_notifications = $2, price_alerts = $3, news_alerts = $4, updated_at = NOW()
        WHERE user_id = $5
    `, settings.EmailEnabled, settings.PushEnabled, settings.AlertsEnabled, settings.NewsEnabled, userID)

	if err != nil {
		http.Error(w, "Failed to update settings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Settings updated successfully"})
}

func (s *NotificationService) sendTestNotification(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	notification := &messaging.Notification{
		Title: "Test Notification",
		Body:  "This is a test notification from the financial analytics dashboard",
	}

	data := map[string]interface{}{
		"type": "test",
		"time": "now",
	}

	s.sendPushNotification(userID, notification, data)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Test notification sent"})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		log.Panicf("failed to marshal: %v", err)
	}
	return b
}
