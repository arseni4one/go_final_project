package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type signInRequest struct {
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token,omitempty"`
	Error string `json:"error,omitempty"`
}

// хэш пароля, чтобы не хранить и не передавать пароль напрямую в токене
func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJson(w, http.StatusBadRequest, signInResponse{Error: "ошибка десериализации JSON"})
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")

	if req.Password != storedPassword {
		writeJson(w, http.StatusOK, signInResponse{Error: "Неверный пароль"})
		return
	}

	claims := jwt.MapClaims{
		"hash": passwordHash(storedPassword),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret-key")) // лучше вынести в константу/конфиг
	if err != nil {
		writeJson(w, http.StatusInternalServerError, signInResponse{Error: "ошибка формирования токена"})
		return
	}

	writeJson(w, http.StatusOK, signInResponse{Token: signedToken})
}

// middleware для защищённых маршрутов
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		storedPassword := os.Getenv("TODO_PASSWORD")
		if len(storedPassword) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		hashFromToken, ok := claims["hash"].(string)
		if !ok || hashFromToken != passwordHash(storedPassword) {
			http.Error(w, "Authentification required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJson(w, status, map[string]string{"error": message})
}
