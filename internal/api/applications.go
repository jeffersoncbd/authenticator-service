package api

import (
	"authenticator/internal/permissions"
	"authenticator/internal/spec"
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
	if err := permissions.Check(r.Context(), applicationsIdentifier, permissions.ToRead); err != nil {
		return spec.ApplicationsListJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	rows, err := api.store.ListApplicaions(r.Context())
	if err != nil {
		api.logger.Error("Falha ao tentar listar aplicações", zap.Error(err))
		return spec.ApplicationsListJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

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
	if err := permissions.Check(r.Context(), applicationsIdentifier, permissions.ToRead); err != nil {
		return spec.FindApplicationByIDJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	applicationId, err := uuid.Parse(id)
	if err != nil {
		return spec.FindApplicationByIDJSON400Response(spec.Error{Feedback: "ID inválido"})
	}
	row, err := api.store.GetApplication(r.Context(), applicationId)
	if err != nil {
		api.logger.Error("Falha ao tentar buscar aplicação", zap.Error(err))
		return spec.FindApplicationByIDJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	application := spec.Application{
		ID:   row.ID.String(),
		Name: row.Name,
	}

	return spec.FindApplicationByIDJSON200Response(application)
}

// Cadastra uma aplicação
// (POST /applications)
func (api API) NewApplication(w http.ResponseWriter, r *http.Request) *spec.Response {
	if err := permissions.Check(r.Context(), applicationsIdentifier, permissions.ToWrite); err != nil {
		return spec.ApplicationsListJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	var application spec.NewApplication

	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		return spec.NewApplicationJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}

	if err := api.validator.validate.Struct(application); err != nil {
		return spec.NewApplicationJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	_, err = api.store.GetApplicationByName(r.Context(), application.Name)
	if err == nil {
		return spec.NewApplicationJSON400Response(spec.Error{Feedback: "Já existe uma aplicação cadastrada com esse nome"})
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar aplicação", zap.Error(err), zap.String("aplicação", application.Name))
		return spec.NewApplicationJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	id, err := api.store.InsertApplication(r.Context(), application.Name)
	if err != nil {
		api.logger.Error("Falha ao cadastrar nova aplicação", zap.Error(err), zap.String("aplicação", application.Name))
		return spec.NewApplicationJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.NewApplicationJSON201Response(spec.BasicCreationResponse{Feedback: "aplicação cadastrada", ID: id.String()})
}
