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
func (api API) UsersList(w http.ResponseWriter, r *http.Request, id string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToRead); err != nil {
		return spec.UsersListJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID
	applicationUUID, err := uuid.Parse(id)
	if err != nil {
		return spec.UsersListJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	// Tenta listar os usuários da aplicação no banco de dados e trata possíveis erros
	rows, err := api.store.ListUsers(r.Context(), applicationUUID)
	if err != nil {
		api.logger.Error("Falha ao listar usuários", zap.Error(err))
		return spec.UsersListJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// Converte os resultados para a estrutura de resposta
	users := []spec.User{}
	for _, row := range rows {
		status := spec.Status{}
		status.FromValue(row.Status)
		users = append(users, spec.User{
			Name:   row.Name,
			Email:  openapi_types.Email(row.Email),
			Status: status,
			Group: struct {
				ID   string "json:\"id\""
				Name string "json:\"name\""
			}{
				ID:   row.GroupID.String(),
				Name: row.GroupName,
			},
		})
	}

	return spec.UsersListJSON200Response(users)
}

// Cadastra um novo usuário
// (POST /applications/{id}/users)
func (api API) NewUser(w http.ResponseWriter, r *http.Request, id string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToWrite); err != nil {
		return spec.NewUserJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Decodifica body e valida dados
	var user spec.NewUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return spec.NewUserJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}

	// Valida UUID da aplicação
	applicationUUID, err := uuid.Parse(id)
	if err != nil {
		return spec.NewUserJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	// Valida UUID do grupo
	groupUUID, err := uuid.Parse(user.GroupID)
	if err != nil {
		return spec.NewUserJSON400Response(spec.Error{Feedback: "ID do grupo de permissões inválido: " + err.Error()})
	}

	// Valida dados de entrada
	if err := api.validator.validate.Struct(user); err != nil {
		return spec.NewUserJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	// Verifica se já existe um usuário com esse e-mail na aplicação no banco de dados
	_, err = api.store.GetUser(r.Context(), postgresql.GetUserParams{
		ApplicationID: applicationUUID,
		Email:         string(user.Email),
	})
	if err == nil {
		return spec.NewUserJSON400Response(spec.Error{Feedback: "Já existe usuário cadastrado com este e-mail"})
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar email", zap.Error(err), zap.String("email", string(user.Email)))
		return spec.NewUserJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// Gera hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		api.logger.Error("Falha ao gerar hash de senha", zap.Error(err))
		return spec.NewUserJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// Insere novo usuário no banco de dados e trata possíveis erros
	err = api.store.InsertUser(r.Context(), postgresql.InsertUserParams{
		Email:         string(user.Email),
		Name:          user.Name,
		Password:      string(hash),
		ApplicationID: applicationUUID,
		GroupID:       groupUUID,
	})
	if err != nil {
		api.logger.Error("Falha ao inserir novo usuário", zap.Error(err))
		return spec.NewUserJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.NewUserJSON201Response(spec.BasicResponse{Feedback: "usuário registrado"})
}

// Atualiza o status de um usuário
// (PATCH /applications/{id}/users/{byEmail})
func (api API) FindUserByEmail(w http.ResponseWriter, r *http.Request, id string, byEmail openapi_types.Email) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := permissions.Check(r.Context(), usersIdentifier, permissions.ToDelete); err != nil {
		return spec.FindUserByEmailJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID da aplicação
	applicationUUID, err := uuid.Parse(id)
	if err != nil {
		return spec.FindUserByEmailJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	// Decodifica body e valida dados
	var body spec.NewUserStatus
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return spec.FindUserByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	// Valida dados de entrada
	if err := api.validator.validate.Struct(body); err != nil {
		return spec.FindUserByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	// Atualiza status do usuário no banco de dados e trata possíveis erros
	if err := api.store.UpdateUserStatus(r.Context(), postgresql.UpdateUserStatusParams{
		ApplicationID: applicationUUID,
		Email:         string(byEmail),
		Status:        body.Status.ToValue(),
	}); err != nil {
		api.logger.Error("Falha ao atualizar status do usuário", zap.Error(err))
		return spec.FindUserByEmailJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.FindUserByEmailJSON200Response(spec.BasicResponse{Feedback: "status do usuário atualizado"})
}
