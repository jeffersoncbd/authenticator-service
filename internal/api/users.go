package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/permissions"
	"authenticator/internal/spec"
	"encoding/json"
	"errors"
	"net/http"

	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const usersIdentifier = "users"

// Lista todos os usuários
// (GET /applications/{id}/users)
func (api API) GetApplicationsIDUsers(w http.ResponseWriter, r *http.Request, id string) *spec.Response {
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToRead); err != nil {
		return spec.GetApplicationsIDUsersJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	applicationUUID, err := uuid.Parse(id)
	if err != nil {
		return spec.PostApplicationsIDGroupsJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	rows, err := api.store.ListUsers(r.Context(), applicationUUID)
	if err != nil {
		api.logger.Error("Falha ao listar usuários", zap.Error(err))
		return spec.GetApplicationsIDUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	var users []spec.User
	for _, row := range rows {
		status := spec.Status{}
		status.FromValue(row.Status)
		users = append(users, spec.User{
			Name:   row.Name,
			Email:  openapi_types.Email(row.Email),
			Status: status,
		})
	}

	return spec.GetApplicationsIDUsersJSON200Response(users)
}

// Cadastra um novo usuário
// (POST /applications/{id}/users)
func (api API) PostApplicationsIDUsers(w http.ResponseWriter, r *http.Request, id string) *spec.Response {
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToWrite); err != nil {
		return spec.PostApplicationsIDUsersJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	applicationUUID, err := uuid.Parse(id)
	if err != nil {
		return spec.PostApplicationsIDGroupsJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	var user spec.NewUser

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return spec.PostApplicationsIDUsersJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}

	if err := api.validator.validate.Struct(user); err != nil {
		return spec.PostApplicationsIDUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	_, err = api.store.GetUser(r.Context(), postgresql.GetUserParams{
		ApplicationID: applicationUUID,
		Email:         string(user.Email),
	})
	if err == nil {
		return spec.PostApplicationsIDUsersJSON400Response(spec.Error{Feedback: "Já existe usuário cadastrado com este e-mail"})
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar email", zap.Error(err), zap.String("email", string(user.Email)))
		return spec.PostApplicationsIDUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		api.logger.Error("Falha ao gerar hash de senha", zap.Error(err))
		return spec.PostApplicationsIDUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}
	err = api.store.InsertUser(r.Context(), postgresql.InsertUserParams{
		Email:         string(user.Email),
		Name:          user.Name,
		Password:      string(hash),
		ApplicationID: applicationUUID,
	})
	if err != nil {
		api.logger.Error("Falha ao inserir novo usuário", zap.Error(err))
		return spec.PostApplicationsIDUsersJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.PostApplicationsIDUsersJSON201Response(spec.BasicResponse{Feedback: "usuário registrado"})
}

// Atualiza o status de um usuário
// (PATCH /applications/{id}/users/{byEmail})
func (api API) PatchApplicationsIDUsersByEmail(w http.ResponseWriter, r *http.Request, id string, byEmail openapi_types.Email) *spec.Response {
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToDelete); err != nil {
		return spec.PatchApplicationsIDUsersByEmailJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	applicationUUID, err := uuid.Parse(id)
	if err != nil {
		return spec.PostApplicationsIDGroupsJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	var body spec.NewUserStatus

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return spec.PatchApplicationsIDUsersByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	if err := api.validator.validate.Struct(body); err != nil {
		return spec.PatchApplicationsIDUsersByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	if err := api.store.UpdateUserStatus(r.Context(), postgresql.UpdateUserStatusParams{
		ApplicationID: applicationUUID,
		Email:         string(byEmail),
		Status:        body.Status.ToValue(),
	}); err != nil {
		api.logger.Error("Falha ao atualizar status do usuário", zap.Error(err))
		return spec.PatchApplicationsIDUsersByEmailJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.PatchApplicationsIDUsersByEmailJSON200Response(spec.BasicResponse{Feedback: "status do usuário atualizado"})
}
