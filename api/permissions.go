package api

import (
	postgresql "authenticator/interfaces"
	"authenticator/spec"
	"authenticator/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Adiciona ou atualiza permissão de um grupo de usuários de uma aplicação
// (POST /applications/{application_id}/groups/{group_id}/permissions)
func (api API) SetPermission(w http.ResponseWriter, r *http.Request, applicationID string, groupID string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := utils.CheckPermissions(r.Context(), groupsIdentifier, utils.KeyToWrite); err != nil {
		return spec.SetPermissionJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID da aplicação
	applicationUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return spec.SetPermissionJSON400Response(spec.Error{Feedback: utils.INVALID_APPLICATION_ID + err.Error()})
	}

	// Valida UUID do grupo
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		return spec.SetPermissionJSON400Response(spec.Error{Feedback: "ID do grupo inválido: " + err.Error()})
	}

	// Decodifica e valida os dados da nova permissão
	var permission spec.Permission
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		return spec.SetPermissionJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}
	if err := api.validator.validate.Struct(permission); err != nil {
		return spec.SetPermissionJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	newPermission := postgresql.AddOrUpdateKeyInGroupParams{
		ApplicationID: applicationUUID,
		ID:            groupUUID,
		Path:          fmt.Sprintf(`{%v}`, permission.Key),
		Replacement:   []byte(fmt.Sprintf(`%v`, permission.Permission)),
	}

	if err := api.store.AddOrUpdateKeyInGroup(r.Context(), newPermission); err != nil {
		api.logger.Error("Falha ao tentar adicionar permissão ao grupo", zap.Error(err))
		return spec.SetPermissionJSON500Response(spec.InternalServerError{Feedback: "Erro ao adicionar permissão ao grupo: " + err.Error()})
	}

	return spec.SetPermissionJSON201Response(spec.BasicResponse{Feedback: "Permissão adicionada ao grupo"})
}

// Exclui uma permissão de um grupo de usuários de uma aplicação
// (DELETE /applications/{application_id}/groups/{group_id}/permissions/{permission_key})
func (api API) DeletePermission(w http.ResponseWriter, r *http.Request, applicationID string, groupID string, permissionKey string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := utils.CheckPermissions(r.Context(), groupsIdentifier, utils.KeyToDelete); err != nil {
		return spec.DeletePermissionJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID da aplicação
	applicationUUID, err := uuid.Parse(applicationID)
	if err != nil {
		return spec.DeletePermissionJSON400Response(spec.Error{Feedback: utils.INVALID_APPLICATION_ID + err.Error()})
	}

	// Valida UUID do grupo
	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		return spec.DeletePermissionJSON400Response(spec.Error{Feedback: "ID do grupo inválido: " + err.Error()})
	}

	args := postgresql.RemoveKeyInGroupParams{
		ApplicationID: applicationUUID,
		ID:            groupUUID,
		Permissions:  []byte(permissionKey),
	}

	if err := api.store.RemoveKeyInGroup(r.Context(), args); err != nil {
		api.logger.Error("Falha ao tentar remover permissão ao grupo", zap.Error(err))
		return spec.SetPermissionJSON500Response(spec.InternalServerError{Feedback: "Erro ao remover permissão ao grupo: " + err.Error()})
	}

	return spec.DeletePermissionJSON200Response(spec.BasicResponse{Feedback: "Permissão removida do grupo"})
}

