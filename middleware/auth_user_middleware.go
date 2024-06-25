package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/ichtrojan/muzz/database"
	"github.com/ichtrojan/muzz/helpers"
	"github.com/ichtrojan/muzz/models"
	"net/http"
	"strings"
)

type ctxKey string

const (
	userKey ctxKey = "user"
)

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", -1)

		unauthorized := helpers.PrepareMessage("unauthorized")

		if token == "" {
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusUnauthorized)

			_ = json.NewEncoder(w).Encode(unauthorized)

			return
		}

		validation, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("testing"), nil
		})

		var foundUser models.User

		if claims, ok := validation.Claims.(jwt.MapClaims); ok && validation.Valid {
			_ = database.Connection.Where("id = ?", claims["user_id"]).First(&foundUser)

			if foundUser.Empty() {
				w.Header().Set("Content-Type", "application/json")

				w.WriteHeader(http.StatusUnauthorized)

				_ = json.NewEncoder(w).Encode(unauthorized)

				return
			}

			r = r.WithContext(context.WithValue(r.Context(), userKey, foundUser))

			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusUnauthorized)

			_ = json.NewEncoder(w).Encode(unauthorized)

			return
		}
	})
}

func GetUser(ctx context.Context) models.User {
	return ctx.Value(userKey).(models.User)
}
