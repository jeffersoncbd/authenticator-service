package api

import (
	postgresql "authenticator/interfaces"
	"authenticator/spec"
	"authenticator/utils"
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
	if err := utils.CheckPermissions(r.Context(), groupsIdentifier, utils.KeyToWrite); err != nil {
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
		return spec.NewGroupJSON400Response(spec.Error{Feedback: utils.INVALID_APPLICATION_ID + err.Error()})
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
	if err := utils.CheckPermissions(r.Context(), groupsIdentifier, utils.KeyToWrite); err != nil {
		return spec.GroupsListJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID da aplicação
	applicationUUID, err := uuid.Parse(applicationId)
	if err != nil {
		return spec.GroupsListJSON400Response(spec.Error{Feedback: utils.INVALID_APPLICATION_ID + err.Error()})
	}

	// Tenta buscar os grupos da aplicação no banco de dados e trata possíveis erros
	rows, err := api.store.ListGrousByApplicationId(r.Context(), applicationUUID)
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

// Deleta um grupo de permissões de uma aplicação
// (DELETE /applications/{application_id}/groups/{group_id})
func (api API) DeleteGroup(w http.ResponseWriter, r *http.Request, applicationID string, groupID string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := utils.CheckPermissions(r.Context(), groupsIdentifier, utils.KeyToDelete); err != nil {
		return spec.DeleteGroupJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID da aplicação
	applicationUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return spec.DeleteGroupJSON400Response(spec.Error{Feedback: utils.INVALID_APPLICATION_ID + err.Error()})
	}
	// Valida UUID do grupo
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		return spec.DeleteGroupJSON400Response(spec.Error{Feedback: utils.INVALID_GROUP_ID + err.Error()})
	}

	// remove grupo da aplicação
	err = api.store.DeleteGroup(r.Context(), postgresql.DeleteGroupParams{
		ID: groupUUID,
		ApplicationID: applicationUUID,
	})
	if err != nil {
		api.logger.Error("Falha ao tentar remover grupo da aplicação", zap.Error(err))
		return spec.DeleteGroupJSON500Response(spec.InternalServerError{Feedback: utils.INTERNAL_SERVER_ERROR})
	}

	return spec.DeleteGroupJSON200Response(spec.BasicResponse{Feedback: "grupo removido com sucesso"})
}


