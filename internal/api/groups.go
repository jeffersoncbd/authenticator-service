package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/permissions"
	"authenticator/internal/spec"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const groupsIdentifier = "groups"

// Cadastra um novo grupo de permissões
// (POST /applications/{id}/groups)
func (api API) NewGroup(w http.ResponseWriter, r *http.Request, applicationId string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := permissions.Check(r.Context(), groupsIdentifier, permissions.ToWrite); err != nil {
		return spec.NewGroupJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Decodifica e valida os dados do novo grupo
	var group spec.NewGroup
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		return spec.NewGroupJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}
	if err := api.validator.validate.Struct(group); err != nil {
		return spec.NewGroupJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	// Valida UUID da aplicação
	applicationUuidID, err := uuid.Parse(applicationId)
	if err != nil {
		return spec.NewGroupJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	// Insere o novo grupo no banco de dados e trata possíveis erros
	groupID, err := api.store.InsertGroup(r.Context(), postgresql.InsertGroupParams{
		Name:          group.Name,
		ApplicationID: applicationUuidID,
		Permissions:   []byte("{}"),
	})
	if err != nil {
		api.logger.Error("Falha ao tentar inserir grupo", zap.Error(err))
		return spec.NewGroupJSON500Response(spec.InternalServerError{Feedback: "Erro ao cadastrar grupo: " + err.Error()})
	}

	return spec.NewGroupJSON201Response(spec.BasicCreationResponse{Feedback: "grupo cadastrado", ID: groupID.String()})
}

// Lista os grupos de permissões de um aplicativo
// (GET /applications/{id}/groups)
func (api API) GroupsList(w http.ResponseWriter, r *http.Request, applicationId string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := permissions.Check(r.Context(), groupsIdentifier, permissions.ToWrite); err != nil {
		return spec.GroupsListJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID da aplicação
	applicationUuidID, err := uuid.Parse(applicationId)
	if err != nil {
		return spec.GroupsListJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	// Tenta buscar os grupos da aplicação no banco de dados e trata possíveis erros
	rows, err := api.store.ListGrousByApplicationId(r.Context(), applicationUuidID)
	if err != nil {
		api.logger.Error("Falha ao tentar listar grupos da aplicação", zap.Error(err))
		return spec.GroupsListJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// Converte os resultados para a estrutura de resposta e retorna-os
	groups := []spec.Group{}
	for _, row := range rows {
		permissions := make(map[string]interface{})
		if err := json.Unmarshal(row.Permissions, &permissions); err != nil {
			api.logger.Error("Falha ao tentar converter permissões", zap.Error(err))
			return spec.GroupsListJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
		}
		groups = append(groups, spec.Group{
			ID:          row.ID.String(),
			Name:        row.Name,
			Permissions: permissions,
		})
	}

	return spec.GroupsListJSON200Response(groups)
}
