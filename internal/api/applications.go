package api

import (
	"authenticator/internal/permissions"
	"authenticator/internal/spec"
	"net/http"

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
