package middlewares

import (
	postgresql "authenticator/interfaces"
	"authenticator/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type JwtMiddleware struct {
	logger *zap.Logger
	store  *postgresql.Queries
}

func NewJwtMiddleware(logger *zap.Logger, pool *pgxpool.Pool) *JwtMiddleware {
	store := postgresql.New(pool)
	return &JwtMiddleware{logger, store}
}

type middleware func(http.Handler) http.Handler

func (m *JwtMiddleware) Middleware() middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// ignora rotas publicas
			if r.URL.Path == "/login" || strings.HasPrefix(r.URL.Path, "/docs") {
				next.ServeHTTP(w, r)
				return
			}

			token, err := m.extractJwt(r.Header.Get("Authorization"), r.Context())
			if err != nil {
				http.Error(w, fmt.Sprintf("{ \"feedback\": \"%v\" }", err.Error()), http.StatusUnauthorized)
				return
			}

			// recupera os dados do JWT
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				m.logger.Error("Token inválido", zap.Any("claims", claims))
				http.Error(w, "{ \"feedback\": \"Falha ao tentar ler token, tente novamente em alguns minutos\" }", http.StatusUnauthorized)
				return
			}

			// recupera permissões do grupo na aplicação
			permissionsString := fmt.Sprintf("%v", claims["roles"])
			permissions := make(map[string]*int)
			json.Unmarshal([]byte(permissionsString), &permissions)

			// cria contexto com as permissões
			ctx := context.WithValue(r.Context(), utils.ContextKey, permissions)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m *JwtMiddleware) extractJwt(header string, ctx context.Context) (*jwt.Token, error) {
			parts := strings.Split(header, " ")
			if header == "" || len(parts) < 1 {
				return nil, fmt.Errorf("não foi fornecido token de autenticação")
			}
			bearerToken := parts[1]

			// recupera o token de assinatura da aplicação do JWT
			token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
				// valida metodo de assinatura
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
				}

				// recupera claims do token
				claims, _ := token.Claims.(jwt.MapClaims)

				// recupera a aplicação que assinou o token
				application, err := m.store.GetApplication(ctx, uuid.MustParse(fmt.Sprintf("%v", claims["aud"])))
				if err != nil {
					m.logger.Error("Falha ao tentar buscar aplicação por applicationId", zap.Error(err))
					return nil, fmt.Errorf("internal server error")
				}

				// retorna token de assinatura
				return []byte(application.Secret.String()), nil
			})

			return token, err
}