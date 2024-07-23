package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/spec"
	"encoding/json"
	"errors"
	"net/http"

	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Lista todos os usuários
// (GET /users)
func (api API) GetUsers(w http.ResponseWriter, r *http.Request) *spec.Response {
	rows, err := api.store.ListUsers(r.Context())
	if err != nil {
		api.logger.Error("Falha ao listar usuários", zap.Error(err))
		return spec.GetUsersJSON500Response(spec.InternalServerError{})
	}

	var users []spec.UserData
	for _, row := range rows {
		status := spec.UserStatus{}
		status.FromValue(row.Status.String)
		users = append(users, spec.UserData{
			Name:   row.Name,
			Email:  openapi_types.Email(row.Email),
			Status: status,
		})
	}

	return spec.GetUsersJSON200Response(users)
}

// Cadastra um novo usuário
// (POST /users)
func (api API) PostUsers(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.User

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	if err := api.validator.validate.Struct(body); err != nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	_, err = api.store.GetUser(r.Context(), string(body.Email))
	if err == nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Já existe usuário cadastrado com este e-mail"})
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar email", zap.Error(err), zap.String("email", string(body.Email)))
		return spec.PostUsersJSON500Response(spec.InternalServerError{})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		api.logger.Error("Falha ao gerar hash de senha", zap.Error(err))
		return spec.PostUsersJSON500Response(spec.InternalServerError{})
	}
	err = api.store.InsertUser(r.Context(), postgresql.InsertUserParams{
		Email:    string(body.Email),
		Name:     body.Name,
		Password: string(hash),
	})
	if err != nil {
		api.logger.Error("Falha ao inserir novo usuário", zap.Error(err))
		return spec.PostUsersJSON500Response(spec.InternalServerError{})
	}

	return nil
}

// Atualiza o status de um usuário
// (PATCH /users/{byEmail})
func (api API) PatchUsersByEmail(w http.ResponseWriter, r *http.Request, byEmail openapi_types.Email) *spec.Response {
	var body spec.PatchUserStatus

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return spec.PatchUsersByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	if err := api.validator.validate.Struct(body); err != nil {
		return spec.PatchUsersByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	if err := api.store.UpdateUserStatus(r.Context(), postgresql.UpdateUserStatusParams{
		Email:  string(byEmail),
		Status: pgtype.Text{String: body.Status.ToValue(), Valid: true},
	}); err != nil {
		api.logger.Error("Falha ao atualizar status do usuário", zap.Error(err))
		return spec.PatchUsersByEmailJSON500Response(spec.InternalServerError{})
	}

	return nil
}
