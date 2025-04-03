package authenticate

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	userclient "github.com/samObot19/shopverse/api-gate-way/user-client"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getGoogleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

var refreshTokens = make(map[string]string)
var refreshTokensMutex sync.Mutex

var invalidatedTokens = struct {
	sync.RWMutex
	tokens map[string]struct{}
}{tokens: make(map[string]struct{})}

type contextKey string

const userKey contextKey = "user"

var users = make(map[string]struct {
	Email          string
	Username       string
	Role           string
	GoogleID       string
	ProfilePicture string
})

func GenerateAccessToken(email, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(15 * time.Minute).Unix(), // 15 minutes expiry
		"type":     "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken(email, username, role string) (string, error) {
	claims := jwt.MapClaims{
		"email":    email,
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(30 * 24 * time.Hour).Unix(), // 30 days expiry
		"type":     "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	refreshTokensMutex.Lock()
	refreshTokens[tokenString] = email
	refreshTokensMutex.Unlock()

	return tokenString, nil
}

func HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	config := getGoogleOauthConfig()
	url := config.AuthCodeURL("randomstate")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	config := getGoogleOauthConfig()

	state := r.URL.Query().Get("state")
	if state != "randomstate" {
		http.Error(w, "Invalid state", http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		GoogleID       string `json:"id"`              
		Email          string `json:"email"`          
		Name           string `json:"name"`           
		ProfilePicture string `json:"picture"`         
	}
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &userInfo)

	userServiceAddress := os.Getenv("USER_SERVICE_ADDRESS")
	if userServiceAddress == "" {
		http.Error(w, "USER_SERVICE_ADDRESS not set in environment", http.StatusInternalServerError)
		return
	}

	conn, err := userclient.ConnectToUserService(userServiceAddress)
	if err != nil {
		http.Error(w, "Failed to connect to user service", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	userClient := userclient.NewUserClient(conn)

	
	existingUser, err := userClient.GetUser(context.Background(), userInfo.Email)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
	}

	if existingUser == nil || existingUser.User == nil {
		_, err := userClient.AddUser(context.Background(), userInfo.Name, userInfo.Email, "default_password", userInfo.GoogleID, userInfo.ProfilePicture)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		log.Printf("New user created: %s", userInfo.Email)
	}

	
	role := "user"
	if userInfo.Email == "admin@example.com" {
		role = "admin"
	}

	accessToken, err := GenerateAccessToken(userInfo.Email, userInfo.Name, role)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := GenerateRefreshToken(userInfo.Email, userInfo.Name, role)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Save user info in memory (for demonstration purposes)
	users[userInfo.Email] = struct {
		Email          string
		Username       string
		Role           string
		GoogleID       string
		ProfilePicture string
	}{
		Email:          userInfo.Email,
		Username:       userInfo.Name,
		Role:           role,
		GoogleID:       userInfo.GoogleID,
		ProfilePicture: userInfo.ProfilePicture,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(request.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	refreshTokensMutex.Lock()
	email, exists := refreshTokens[request.RefreshToken]
	refreshTokensMutex.Unlock()
	if !exists {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}
	username := claims["username"].(string)
	role := claims["role"].(string)

	accessToken, err := GenerateAccessToken(email, username, role)
	if err != nil {
		http.Error(w, "Failed to generate new access token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Authorization header missing or invalid", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	invalidatedTokens.Lock()
	invalidatedTokens.tokens[tokenStr] = struct{}{}
	invalidatedTokens.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}

func IsTokenInvalidated(token string) bool {
	invalidatedTokens.RLock()
	defer invalidatedTokens.RUnlock()
	_, exists := invalidatedTokens.tokens[token]
	return exists
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		if IsTokenInvalidated(tokenStr) {
			http.Error(w, "Token has been invalidated", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["type"] != "access" {
			http.Error(w, "Invalid access token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}