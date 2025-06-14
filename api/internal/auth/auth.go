// internal/auth/auth.go
package auth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/htstinson/stinsondataapi/api/internal/model"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID     string                            `json:"user_id"`
	Username   string                            `json:"username"`
	Roles      string                            `json:"roles"`
	IP_Address string                            `json:"ip_address"`
	Subscribed []model.User_Subscriber_Role_View `json:"subscribed"`
	jwt.RegisteredClaims
}

type Config struct {
	SecretKey     string
	TokenDuration time.Duration
}

type JWTAuth struct {
	Config Config
}

func New(config Config) *JWTAuth {
	return &JWTAuth{Config: config}
}

func (a *JWTAuth) GenerateToken(user model.User, roles model.Roles, subscribed []model.User_Subscriber_Role_View) (string, error) {
	now := time.Now()

	fmt.Println("GenerateToken")

	fmt.Println(user.IP_address)
	fmt.Println(user.ID)

	claims := Claims{
		UserID:     user.ID,
		Username:   user.Username,
		Roles:      roles.Names,
		IP_Address: user.IP_address,
		Subscribed: subscribed,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(a.Config.TokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.Config.SecretKey))
}

func (a *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.Config.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (a *JWTAuth) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("handler err", err.Error())
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := a.ValidateToken(tokenParts[1])
		if err != nil {
			if err == ErrExpiredToken {
				http.Error(w, "Token has expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
