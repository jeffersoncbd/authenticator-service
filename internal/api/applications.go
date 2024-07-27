package api

import (
	"authenticator/internal/databases/postgresql"
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
func (api API) GetApplications(w http.ResponseWriter, r *http.Request) *spec.Response {
	if err := permissions.Check(r.Context(), applicationsIdentifier, permissions.ToRead); err != nil {
		return spec.GetApplicationsJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	rows, err := api.store.ListApplicaions(r.Context())
	if err != nil {
		api.logger.Error("Falha ao tentar listar aplicações", zap.Error(err))
		return spec.GetApplicationsJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	var applications []spec.Application
	for _, row := range rows {
		applications = append(applications, spec.Application{
			ID:   row.ID.String(),
			Name: row.Name,
			Keys: row.Keys,
		})
	}

	return spec.GetApplicationsJSON200Response(applications)
}

// Cadastra uma aplicação
// (POST /applications)
func (api API) PostApplications(w http.ResponseWriter, r *http.Request) *spec.Response {
	if err := permissions.Check(r.Context(), applicationsIdentifier, permissions.ToWrite); err != nil {
		return spec.GetApplicationsJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	var application spec.NewApplication

	err := json.NewDecoder(r.Body).Decode(&application)
	if err != nil {
		return spec.PostApplicationsJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}

	if err := api.validator.validate.Struct(application); err != nil {
		return spec.PostApplicationsJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	_, err = api.store.GetApplicationByName(r.Context(), application.Name)
	if err == nil {
		return spec.PostApplicationsJSON400Response(spec.Error{Feedback: "Já existe uma aplicação cadastrada com esse nome"})
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar aplicação", zap.Error(err), zap.String("aplicação", application.Name))
		return spec.PostApplicationsJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	id, err := api.store.InsertApplication(r.Context(), postgresql.InsertApplicationParams{
		Name: application.Name,
		Keys: []string{},
	})
	if err != nil {
		api.logger.Error("Falha ao cadastrar nova aplicação", zap.Error(err), zap.String("aplicação", application.Name))
		return spec.PostApplicationsJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.PostApplicationsJSON201Response(spec.NewApplicationResponse{Feedback: "aplicação cadastrada", ID: id.String()})
}

// Adiciona uma chave de permissão na aplicação
// (POST /applications/{id}/keys)
func (api API) PostApplicationsIDKeys(w http.ResponseWriter, r *http.Request, id string) *spec.Response {
	if err := permissions.Check(r.Context(), applicationsIdentifier, permissions.ToRead); err != nil {
		return spec.PostApplicationsIDKeysJSON401Response(spec.Unauthorized{Feedback: err.Error()})
	}

	applicationId, err := uuid.Parse(id)
	if err != nil {
		return spec.PostApplicationsIDKeysJSON400Response(spec.Error{Feedback: "ID inválido"})
	}

	var body spec.Keys

	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return spec.PostApplicationsIDKeysJSON400Response(spec.Error{Feedback: "Erro de decodificação: " + err.Error()})
	}

	if err := api.validator.validate.Struct(body); err != nil {
		return spec.PostApplicationsJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	err = api.store.InsertKey(r.Context(), postgresql.InsertKeyParams{ID: applicationId, ArrayCat: body.NewKeys})
	if err != nil {
		api.logger.Error("Falha ao adicionar chave à aplicação", zap.Error(err), zap.String("aplicação", id))
		return spec.PostApplicationsIDKeysJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	return spec.PostApplicationsIDKeysJSON200Response(spec.BasicResponse{Feedback: "chave adicionada"})
}
