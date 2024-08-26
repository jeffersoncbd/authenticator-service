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
func (api API) PostApplicationsIDGroups(w http.ResponseWriter, r *http.Request, applicationId string) *spec.Response {
	if err := permissions.Check(r.Context(), groupsIdentifier, permissions.ToWrite); err != nil {
		return spec.PostApplicationsIDGroupsJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	var group spec.NewGroup

	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		return spec.PostApplicationsIDGroupsJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}

	if err := api.validator.validate.Struct(group); err != nil {
		return spec.PostApplicationsIDGroupsJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	applicationUuidID, err := uuid.Parse(applicationId)
	if err != nil {
		return spec.PostApplicationsIDGroupsJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	groupID, err := api.store.InsertGroup(r.Context(), postgresql.InsertGroupParams{
		Name:          group.Name,
		ApplicationID: applicationUuidID,
		Permissions:   []byte{},
	})
	if err != nil {
		api.logger.Error("Falha ao tentar inserir grupo", zap.Error(err))
		return spec.PostApplicationsIDGroupsJSON500Response(spec.InternalServerError{Feedback: "Erro ao cadastrar grupo: " + err.Error()})
	}

	return spec.PostApplicationsIDGroupsJSON201Response(spec.BasicCreationResponse{Feedback: "grupo cadastrado", ID: groupID.String()})
}

// Lista os grupos de permissões de um aplicativo
// (GET /applications/{id}/groups)
func (api API) GetApplicationsIDGroups(w http.ResponseWriter, r *http.Request, applicationId string) *spec.Response {
	if err := permissions.Check(r.Context(), groupsIdentifier, permissions.ToWrite); err != nil {
		return spec.GetApplicationsIDGroupsJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	applicationUuidID, err := uuid.Parse(applicationId)
	if err != nil {
		return spec.PostApplicationsIDGroupsJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	rows, err := api.store.ListGrousByApplicationId(r.Context(), applicationUuidID)
	if err != nil {
		api.logger.Error("Falha ao tentar listar grupos da aplicação", zap.Error(err))
		return spec.GetApplicationsIDGroupsJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	groups := []spec.Group{}

	for _, row := range rows {
		permissions := make(map[string]interface{})
		if err := json.Unmarshal(row.Permissions, &permissions); err != nil {
			api.logger.Error("Falha ao tentar converter permissões", zap.Error(err))
			return spec.GetApplicationsIDGroupsJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
		}
		groups = append(groups, spec.Group{
			ID:          row.ID.String(),
			Name:        row.Name,
			Permissions: permissions,
		})
	}

	println(groups)
	return spec.GetApplicationsIDGroupsJSON200Response(groups)
}
