package api

import (
	postgresql "authenticator/interfaces"
	"authenticator/spec"
	"authenticator/utils"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const applicationsIdentifier = "applications"

// Lista todas as aplicações
// (GET /application)
func (api API) ApplicationsList(w http.ResponseWriter, r *http.Request) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := utils.CheckPermissions(r.Context(), applicationsIdentifier, utils.KeyToRead); err != nil {
		return spec.ApplicationsListJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Tenta listar as aplicações no banco de dados e trata possíveis erros
	rows, err := api.store.ListApplicaions(r.Context())
	if err != nil {
		api.logger.Error("Falha ao tentar listar aplicações", zap.Error(err))
		return spec.ApplicationsListJSON500Response(spec.InternalServerError{Feedback: utils.INTERNAL_SERVER_ERROR})
	}

	// Converte os resultados para a estrutura de resposta
	var applications []spec.Application
	for _, row := range rows {
		applications = append(applications, spec.Application{
			ID:   row.ID.String(),
			Name: row.Name,
		})
	}

	return spec.ApplicationsListJSON200Response(applications)
}

// Lista todas as aplicações
// (GET /applications/{id})
func (api API) FindApplicationByID(w http.ResponseWriter, r *http.Request, id string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := utils.CheckPermissions(r.Context(), applicationsIdentifier, utils.KeyToRead); err != nil {
		return spec.FindApplicationByIDJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID
	applicationId, err := uuid.Parse(id)
	if err != nil {
		return spec.FindApplicationByIDJSON400Response(spec.Error{Feedback: "ID inválido"})
	}

	// Tenta buscar a aplicação no banco de dados e trata possíveis erros
	row, err := api.store.GetApplication(r.Context(), applicationId)
	if err != nil {
		api.logger.Error("Falha ao tentar buscar aplicação", zap.Error(err))
		return spec.FindApplicationByIDJSON500Response(spec.InternalServerError{Feedback: utils.INTERNAL_SERVER_ERROR})
	}

	// Converte os resultados para a estrutura de resposta
	application := spec.Application{
		ID:   row.ID.String(),
		Name: row.Name,
	}

	return spec.FindApplicationByIDJSON200Response(application)
}

// Cadastra uma aplicação
// (POST /applications)
func (api API) NewApplication(w http.ResponseWriter, r *http.Request) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := utils.CheckPermissions(r.Context(), applicationsIdentifier, utils.KeyToWrite); err != nil {
		return spec.ApplicationsListJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Decodifica body e valida dados
	var application spec.NewApplication
	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		return spec.NewApplicationJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}
	if err := api.validator.validate.Struct(application); err != nil {
		return spec.NewApplicationJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	// Verifica se já existe uma aplicação com esse nome no banco de dados
	_, err = api.store.GetApplicationByName(r.Context(), application.Name)
	if err == nil {
		return spec.NewApplicationJSON400Response(spec.Error{Feedback: "Já existe uma aplicação cadastrada com esse nome"})
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar aplicação", zap.Error(err), zap.String("aplicação", application.Name))
		return spec.NewApplicationJSON500Response(spec.InternalServerError{Feedback: utils.INTERNAL_SERVER_ERROR})
	}

	// Cadastra nova aplicação no banco de dados e trata possíveis erros
	id, err := api.store.InsertApplication(r.Context(), application.Name)
	if err != nil {
		api.logger.Error("Falha ao cadastrar nova aplicação", zap.Error(err), zap.String("aplicação", application.Name))
		return spec.NewApplicationJSON500Response(spec.InternalServerError{Feedback: utils.INTERNAL_SERVER_ERROR})
	}

	return spec.NewApplicationJSON201Response(spec.BasicCreationResponse{Feedback: "aplicação cadastrada", ID: id.String()})
}

// Renomeia uma aplicação
// (PATCH /applications/{id})
func (api API) RenameApplication(w http.ResponseWriter, r *http.Request, id string) *spec.Response {
	// Verifica se requisição possui a permissão necessária
	if err := utils.CheckPermissions(r.Context(), applicationsIdentifier, utils.KeyToWrite); err != nil {
		return spec.RenameApplicationJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	// Valida UUID
	applicationUUID, err := uuid.Parse(id)
	if err != nil {
		return spec.RenameApplicationJSON400Response(spec.Error{Feedback: utils.INVALID_APPLICATION_ID + err.Error()})
	}

	// Decodifica body e valida dados
	var updatedApplication spec.UpdatedApplication
	err = json.NewDecoder(r.Body).Decode(&updatedApplication)
	if err != nil {
		return spec.RenameApplicationJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}
	if err := api.validator.validate.Struct(updatedApplication); err != nil {
		return spec.RenameApplicationJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	// Renomeia aplicação no banco de dados e trata possíveis erros
	err = api.store.RenameApplication(r.Context(), postgresql.RenameApplicationParams{
		ID:  	applicationUUID,
		Name: updatedApplication.NewName,
	})
	if err != nil {
		api.logger.Error("Falha ao renomear aplicação", zap.Error(err), zap.String("ID", applicationUUID.String()))
		return spec.RenameApplicationJSON500Response(spec.InternalServerError{Feedback: utils.INTERNAL_SERVER_ERROR})
	}

	return spec.RenameApplicationJSON200Response(spec.BasicResponse{Feedback: "aplicação renomeada"})
}


