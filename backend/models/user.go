package models

import "time"

type User struct {
    ID           int       `json:"id"`
    Email        string    `json:"email"`
    Name         string    `json:"name,omitempty"`
    AvatarURL    string    `json:"avatar_url,omitempty"`
    Provider     string    `json:"provider,omitempty"`
    ProviderID   string    `json:"provider_id,omitempty"`
    CreatedAt    time.Time `json:"created_at"`
}

type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}
