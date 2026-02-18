package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type AuthService struct {
	db           *sql.DB
	firebaseAuth *auth.Client
	jwtSecret    []byte
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Provider string `json:"provider,omitempty"`
	Token    string `json:"token,omitempty"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	User         User   `json:"user"`
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
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

	firebaseAuth, err := app.Auth(ctx)
	if err != nil {
		log.Fatal("Failed to get Firebase Auth client:", err)
	}

	service := &AuthService{
		db:           db,
		firebaseAuth: firebaseAuth,
		jwtSecret:    []byte(os.Getenv("JWT_SECRET")),
	}

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc("/login", service.handleLogin).Methods("POST")
	router.HandleFunc("/register", service.handleRegister).Methods("POST")
	router.HandleFunc("/refresh", service.handleRefreshToken).Methods("POST")
	router.HandleFunc("/verify", service.handleVerifyToken).Methods("POST")
	router.HandleFunc("/logout", service.handleLogout).Methods("POST")
	router.HandleFunc("/health", handleHealth).Methods("GET")

	// Start gRPC server for internal communication
	go service.startGRPCServer()

	log.Println("Auth service listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}

func (s *AuthService) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user User
	var hashedPassword string

	// Check if user exists
	err := s.db.QueryRow(`
        SELECT id, email, provider, password_hash, created_at 
        FROM users WHERE email = $1
    `, req.Email).Scan(&user.ID, &user.Email, &user.Provider, &hashedPassword, &user.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password for email/password login
	if req.Provider == "" || req.Provider == "email" {
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	} else {
		// Verify OAuth token with Firebase
		token, err := s.firebaseAuth.VerifyIDToken(context.Background(), req.Token)
		if err != nil {
			http.Error(w, "Invalid OAuth token", http.StatusUnauthorized)
			return
		}

		if token.Claims["email"] != req.Email {
			http.Error(w, "Email mismatch", http.StatusUnauthorized)
			return
		}
	}

	// Generate JWT tokens
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Update last login
	if _, err := s.db.Exec("UPDATE users SET last_login = NOW() WHERE id = $1", user.ID); err != nil {
		log.Printf("Failed to update last login for user %s: %v", user.ID, err)
	}

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		User:         user,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (s *AuthService) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if user already exists
	var exists bool
	if err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists); err != nil {
		log.Printf("Failed to check user existence: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	// Create user
	var userID string
	err = s.db.QueryRow(`
        INSERT INTO users (email, provider, password_hash) 
        VALUES ($1, $2, $3) 
        RETURNING id
    `, req.Email, "email", string(hashedPassword)).Scan(&userID)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	// Create user preferences
	if _, err := s.db.Exec("INSERT INTO user_preferences (user_id) VALUES ($1)", userID); err != nil {
		log.Printf("Failed to create user preferences for user %s: %v", userID, err)
		// Non-critical, so we don't return error
	}

	// Generate tokens
	user := User{
		ID:        userID,
		Email:     req.Email,
		Provider:  "email",
		CreatedAt: time.Now(),
	}

	accessToken, _ := s.generateAccessToken(user)
	refreshToken, _ := s.generateRefreshToken(user)

	response := AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    3600,
		User:         user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (s *AuthService) generateAccessToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) generateRefreshToken(user User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Verify refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["type"] != "refresh" {
		http.Error(w, "Invalid token type", http.StatusUnauthorized)
		return
	}

	// Get user
	var user User
	err = s.db.QueryRow(`
        SELECT id, email, provider, created_at 
        FROM users WHERE id = $1
    `, claims["user_id"]).Scan(&user.ID, &user.Email, &user.Provider, &user.CreatedAt)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Generate new access token
	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token": accessToken,
		"expires_in":   3600,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (s *AuthService) handleVerifyToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(map[string]bool{"valid": false}); err != nil {
			log.Printf("Failed to encode response: %v", err)
		}
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"valid":   true,
		"user_id": claims["user_id"],
		"email":   claims["email"],
	}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (s *AuthService) handleLogout(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, you might want to blacklist the token
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"}); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

func (s *AuthService) startGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	// Register gRPC service here
	log.Println("Auth gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal("Failed to serve:", err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "healthy"}); err != nil {
		log.Printf("Failed to encode health response: %v", err)
	}
}
