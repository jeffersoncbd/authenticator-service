package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/permissions"
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

const usersIdentifier = "users"

// Lista todos os usuários
// (GET /users)
func (api API) GetUsers(w http.ResponseWriter, r *http.Request) *spec.Response {
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToRead); err != nil {
		return spec.GetUsersJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	rows, err := api.store.ListUsers(r.Context())
	if err != nil {
		api.logger.Error("Falha ao listar usuários", zap.Error(err))
		return spec.GetUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	var users []spec.User
	for _, row := range rows {
		status := spec.UserStatus{}
		status.FromValue(row.Status.String)
		users = append(users, spec.User{
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
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToWrite); err != nil {
		return spec.PostUsersJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	var user spec.NewUser

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}

	if err := api.validator.validate.Struct(user); err != nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	_, err = api.store.GetUser(r.Context(), string(user.Email))
	if err == nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Já existe usuário cadastrado com este e-mail"})
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar email", zap.Error(err), zap.String("email", string(user.Email)))
		return spec.PostUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		api.logger.Error("Falha ao gerar hash de senha", zap.Error(err))
		return spec.PostUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}
	err = api.store.InsertUser(r.Context(), postgresql.InsertUserParams{
		Email:    string(user.Email),
		Name:     user.Name,
		Password: string(hash),
	})
	if err != nil {
		api.logger.Error("Falha ao inserir novo usuário", zap.Error(err))
		return spec.PostUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.PostUsersJSON201Response(spec.BasicResponse{Feedback: "usuário registrado"})
}

// Atualiza o status de um usuário
// (PATCH /users/{byEmail})
func (api API) PatchUsersByEmail(w http.ResponseWriter, r *http.Request, byEmail openapi_types.Email) *spec.Response {
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToDelete); err != nil {
		return spec.PostUsersJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

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
		return spec.PatchUsersByEmailJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.PatchUsersByEmailJSON200Response(spec.BasicResponse{Feedback: "status do usuário atualizado"})
}
