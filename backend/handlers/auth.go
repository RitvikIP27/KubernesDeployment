package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/trainwithshubham/skillpulse/database"
	"github.com/trainwithshubham/skillpulse/models"
)

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
	err = database.DB.QueryRow(
		"INSERT INTO users (email, password_hash, name) VALUES ($1, $2, $3) RETURNING id, email, name, avatar_url, oauth_provider, oauth_id, created_at",
		req.Email, string(hashedPassword), req.Name,
	).Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.Provider, &user.ProviderID, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
	var passwordHash string
	err := database.DB.QueryRow(
		"SELECT id, email, password_hash, name, avatar_url, oauth_provider, oauth_id, created_at FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &passwordHash, &user.Name, &user.AvatarURL, &user.Provider, &user.ProviderID, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if passwordHash == "" || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
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
	err := database.DB.QueryRow(
		"SELECT id, email, name, avatar_url, oauth_provider, oauth_id, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Email, &user.Name, &user.AvatarURL, &user.Provider, &user.ProviderID, &user.CreatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	c.JSON(http.StatusOK, user)
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
