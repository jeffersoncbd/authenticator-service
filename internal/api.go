package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/spec"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type store interface {
	GetUser(context.Context, string) (postgresql.User, error)
	InsertUser(context.Context, postgresql.InsertUserParams) error
}

type API struct{
	store		store
	logger		*zap.Logger
	validator	*validator.Validate
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger) API {
	return API{
        store: postgresql.New(pool),
        logger: logger,
		validator: validator.New(validator.WithRequiredStructEnabled()),
    }
}

// Cadastra um novo usuário
// (POST /users)
func (api API) PostUsers(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.User

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
        return spec.PostUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
    }

	if err := api.validator.Struct(body); err != nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	_, err = api.store.GetUser(r.Context(), body.Email)
	if err == nil {
        return spec.PostUsersJSON400Response(spec.Error{Feedback: "Usuário já cadastrado"})
    }
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar email", zap.Error(err), zap.String("email", body.Email))
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Falha no cadastro, tente novamente em alguns minutos"})
	}

	err = api.store.InsertUser(r.Context(), postgresql.InsertUserParams{
		Email:    body.Email,
		Name:     body.Name,
		Password: body.Password,
	})
	if err!= nil {
        api.logger.Error("Falha ao inserir novo usuário", zap.Error(err))
        return spec.PostUsersJSON400Response(spec.Error{Feedback: "Falha no cadastro, tente novamente em alguns minutos"})
    }

	return spec.PostUsersJSON201Response(nil)
}
