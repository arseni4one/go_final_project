package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
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

func signingKey(password string) []byte {
	return []byte(os.Getenv("TODO_JWT_SECRET") + password)
}

// хэш пароля, чтобы не хранить и не передавать пароль напрямую в токене
func passwordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", sum)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req signInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ошибка десериализации JSON")
		return
	}

	storedPassword := os.Getenv("TODO_PASSWORD")

	if req.Password != storedPassword {
		writeError(w, http.StatusUnauthorized, "Неверный пароль")
		return
	}

	claims := jwt.MapClaims{
		"hash": passwordHash(storedPassword),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret-key")) // лучше вынести в константу/конфиг
	if err != nil {
		writeError(w, http.StatusInternalServerError, "ошибка формирования токена")
		return
	}

	writeJson(w, http.StatusOK, signInResponse{Token: signedToken})
}

// middleware для защищённых маршрутов
func Auth(storedPassword string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(storedPassword) == 0 {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			writeError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		token, err := jwt.ParseWithClaims(
			cookie.Value,
			&jwt.RegisteredClaims{},
			func(t *jwt.Token) (interface{}, error) {
				return []byte("secret-key"), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		)
		if err != nil || !token.Valid {
			writeError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		hashFromToken, ok := claims["hash"].(string)
		if !ok || hashFromToken != passwordHash(storedPassword) {
			writeError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		next(w, r)
	}
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("writeJson: ошибка сериализации: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"внутренняя ошибка сервера"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		log.Printf("writeJson: ошибка записи ответа: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJson(w, status, map[string]string{"error": message})
}
