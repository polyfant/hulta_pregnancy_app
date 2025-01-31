package models

import jwt "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
    Email string `json:"email"`
    Role  string `json:"role"`
    jwt.RegisteredClaims  
}
