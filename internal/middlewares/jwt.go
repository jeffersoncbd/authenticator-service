package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JwtMiddleware struct {
	logger *zap.Logger
}

func NewJwtMiddleware(logger *zap.Logger) *JwtMiddleware {
	return &JwtMiddleware{logger}
}

type middleware func(http.Handler) http.Handler

func (m *JwtMiddleware) Middleware() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/login" || strings.HasPrefix(r.URL.Path, "/docs") {
				next.ServeHTTP(w, r)
				return
			}

			bearerToken := strings.Split(r.Header.Get("Authorization"), " ")[1]

			token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
				}
				return []byte("implementar-depois"), nil
			})

			if err != nil {
				http.Error(w, fmt.Sprintf("{ \"feedback\": \"%v\" }", err.Error()), http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				m.logger.Error("Token inválido", zap.Any("claims", claims))
				http.Error(w, "{ \"feedback\": \"Falha ao tentar ler token, tente novamente em alguns minutos\" }", http.StatusUnauthorized)
				return
			}

			// email := fmt.Sprintf("%v", claims["sub"])

			next.ServeHTTP(w, r)
		})
	}
}
