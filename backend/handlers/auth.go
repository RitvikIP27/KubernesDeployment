package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/trainwithshubham/skillpulse/database"
	"github.com/trainwithshubham/skillpulse/models"
	"golang.org/x/crypto/bcrypt"
)

const googleOAuthStateCookie = "helixacore_google_oauth_state"

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

type googleUserInfo struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Verified bool   `json:"verified_email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
}

func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	var user models.User
	var name, avatarURL, provider, providerID sql.NullString
	err = database.DB.QueryRow(
		"INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id, email, name, avatar_url, oauth_provider, oauth_id, created_at",
		req.Email, string(hashedPassword), req.Name,
	).Scan(&user.ID, &user.Email, &name, &avatarURL, &provider, &providerID, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.Name = nullStringValue(name)
	user.AvatarURL = nullStringValue(avatarURL)
	user.Provider = nullStringValue(provider)
	user.ProviderID = nullStringValue(providerID)

	token, err := createJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{Token: token, User: user})
}

func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	var passwordHash, name, avatarURL, provider, providerID sql.NullString
	err := database.DB.QueryRow(
		"SELECT id, email, password_hash, name, avatar_url, oauth_provider, oauth_id, created_at FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &passwordHash, &name, &avatarURL, &provider, &providerID, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	user.Name = nullStringValue(name)
	user.AvatarURL = nullStringValue(avatarURL)
	user.Provider = nullStringValue(provider)
	user.ProviderID = nullStringValue(providerID)

	if !passwordHash.Valid || passwordHash.String == "" || bcrypt.CompareHashAndPassword([]byte(passwordHash.String), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := createJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{Token: token, User: user})
}

func GetMe(c *gin.Context) {
	userID, ok := getCurrentUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	var name, avatarURL, provider, providerID sql.NullString
	err := database.DB.QueryRow(
		"SELECT id, email, name, avatar_url, oauth_provider, oauth_id, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &name, &avatarURL, &provider, &providerID, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}
	user.Name = nullStringValue(name)
	user.AvatarURL = nullStringValue(avatarURL)
	user.Provider = nullStringValue(provider)
	user.ProviderID = nullStringValue(providerID)

	c.JSON(http.StatusOK, user)
}

func GoogleLogin(c *gin.Context) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google login is not configured"})
		return
	}

	state, err := randomState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start Google login"})
		return
	}

	setOAuthStateCookie(c, state)
	redirectURI := googleRedirectURI(c)
	values := url.Values{}
	values.Set("client_id", clientID)
	values.Set("redirect_uri", redirectURI)
	values.Set("response_type", "code")
	values.Set("scope", "openid email profile")
	values.Set("state", state)
	values.Set("access_type", "online")
	values.Set("prompt", "select_account")

	c.Redirect(http.StatusFound, "https://accounts.google.com/o/oauth2/v2/auth?"+values.Encode())
}

func GoogleCallback(c *gin.Context) {
	if errMessage := c.Query("error"); errMessage != "" {
		redirectWithAuthError(c, errMessage)
		return
	}

	stateCookie, err := c.Cookie(googleOAuthStateCookie)
	if err != nil || stateCookie == "" || stateCookie != c.Query("state") {
		redirectWithAuthError(c, "invalid OAuth state")
		return
	}
	clearOAuthStateCookie(c)

	code := c.Query("code")
	if code == "" {
		redirectWithAuthError(c, "missing OAuth code")
		return
	}

	token, err := exchangeGoogleCode(c, code)
	if err != nil {
		redirectWithAuthError(c, err.Error())
		return
	}

	googleUser, err := fetchGoogleUser(token.AccessToken)
	if err != nil {
		redirectWithAuthError(c, err.Error())
		return
	}
	if googleUser.Email == "" || !googleUser.Verified {
		redirectWithAuthError(c, "Google email is not verified")
		return
	}

	user, err := upsertGoogleUser(googleUser)
	if err != nil {
		redirectWithAuthError(c, "failed to save Google user")
		return
	}

	jwtToken, err := createJWT(user.ID)
	if err != nil {
		redirectWithAuthError(c, "failed to create token")
		return
	}

	values := url.Values{}
	values.Set("token", jwtToken)
	c.Redirect(http.StatusFound, frontendURL(c)+"/?"+values.Encode())
}

func createJWT(userID int) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", nil
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}

func exchangeGoogleCode(c *gin.Context, code string) (googleTokenResponse, error) {
	var token googleTokenResponse
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		return token, errors.New("Google login is not configured")
	}

	form := url.Values{}
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("code", code)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", googleRedirectURI(c))

	req, err := http.NewRequest(http.MethodPost, "https://oauth2.googleapis.com/token", strings.NewReader(form.Encode()))
	if err != nil {
		return token, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return token, err
	}
	if err := json.Unmarshal(body, &token); err != nil {
		return token, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 || token.AccessToken == "" {
		if token.Description != "" {
			return token, errors.New(token.Description)
		}
		if token.Error != "" {
			return token, errors.New(token.Error)
		}
		return token, fmt.Errorf("Google token exchange failed with status %d", res.StatusCode)
	}

	return token, nil
}

func fetchGoogleUser(accessToken string) (googleUserInfo, error) {
	var googleUser googleUserInfo
	req, err := http.NewRequest(http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return googleUser, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return googleUser, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	if err != nil {
		return googleUser, err
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		return googleUser, err
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return googleUser, fmt.Errorf("Google profile request failed with status %d", res.StatusCode)
	}

	return googleUser, nil
}

func upsertGoogleUser(googleUser googleUserInfo) (models.User, error) {
	var user models.User
	var name, avatarURL, provider, providerID sql.NullString
	err := database.DB.QueryRow(
		`INSERT INTO users (email, name, avatar_url, oauth_provider, oauth_id)
		 VALUES ($1, $2, $3, 'google', $4)
		 ON CONFLICT (email) DO UPDATE SET
		     name = COALESCE(NULLIF(EXCLUDED.name, ''), users.name),
		     avatar_url = COALESCE(NULLIF(EXCLUDED.avatar_url, ''), users.avatar_url),
		     oauth_provider = 'google',
		     oauth_id = EXCLUDED.oauth_id
		 RETURNING id, email, name, avatar_url, oauth_provider, oauth_id, created_at`,
		googleUser.Email, googleUser.Name, googleUser.Picture, googleUser.ID,
	).Scan(&user.ID, &user.Email, &name, &avatarURL, &provider, &providerID, &user.CreatedAt)
	if err != nil {
		return user, err
	}
	user.Name = nullStringValue(name)
	user.AvatarURL = nullStringValue(avatarURL)
	user.Provider = nullStringValue(provider)
	user.ProviderID = nullStringValue(providerID)
	return user, nil
}

func googleRedirectURI(c *gin.Context) string {
	if redirectURI := os.Getenv("GOOGLE_REDIRECT_URI"); redirectURI != "" {
		return redirectURI
	}
	return externalBaseURL(c) + "/api/auth/google/callback"
}

func frontendURL(c *gin.Context) string {
	if appURL := os.Getenv("APP_URL"); appURL != "" {
		return strings.TrimRight(appURL, "/")
	}
	return externalBaseURL(c)
}

func externalBaseURL(c *gin.Context) string {
	scheme := c.GetHeader("X-Forwarded-Proto")
	if scheme == "" {
		scheme = "http"
	}
	host := c.Request.Host
	if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}
	return scheme + "://" + host
}

func randomState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func setOAuthStateCookie(c *gin.Context, state string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(googleOAuthStateCookie, state, 600, "/", "", c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https", true)
}

func clearOAuthStateCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(googleOAuthStateCookie, "", -1, "/", "", c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https", true)
}

func redirectWithAuthError(c *gin.Context, message string) {
	values := url.Values{}
	values.Set("auth_error", message)
	c.Redirect(http.StatusFound, frontendURL(c)+"/?"+values.Encode())
}

func nullStringValue(value sql.NullString) string {
	if value.Valid {
		return value.String
	}
	return ""
}
