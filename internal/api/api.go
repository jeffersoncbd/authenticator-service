package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/spec"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	pt_br_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func (v Validator) Translate(err error) string {
	var errorMessages []string
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		errorMessages = append(errorMessages, e.Translate(v.translator))
	}
	return strings.Join(errorMessages, ", ")
}

type API struct {
	store     *postgresql.Queries
	logger    *zap.Logger
	validator Validator
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger) API {
	pt_br := pt_BR.New()
	universal_translator := ut.New(pt_br, pt_br)

	translator, err := universal_translator.GetTranslator("pt_BR")
	if !err {
		logger.Fatal("Falha ao carregar tradutores")
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	pt_br_translations.RegisterDefaultTranslations(validate, translator)

	return API{
		store:  postgresql.New(pool),
		logger: logger,
		validator: Validator{
			validate:   validate,
			translator: translator,
		},
	}
}

// Autentica usuário
// (POST /applications/{id}/login)
func (api API) PostLogin(w http.ResponseWriter, r *http.Request) *spec.Response {
	var credentials spec.LoginCredentials

	// decodifica body armazenando dados nas credenciais
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return spec.PostLoginJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	// valida dados de entrada
	if err := api.validator.validate.Struct(credentials); err != nil {
		return spec.PostLoginJSON400Response(spec.Error{Feedback: api.validator.Translate(err)})
	}

	applicationUUID, err := uuid.Parse(credentials.Application)
	if err != nil {
		return spec.PostLoginJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	// busca usuário no banco de dados
	user, err := api.store.GetUser(r.Context(), postgresql.GetUserParams{
		ApplicationID: applicationUUID,
		Email:         string(credentials.Email),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PostLoginJSON400Response(spec.Error{Feedback: "E-mail ou senha inválidos"})
		}
		api.logger.Error("Falha ao buscar usuário", zap.String("email", string(credentials.Email)), zap.Error(err))
		return spec.PostLoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// verifica status
	if user.Status == "inactive" {
		return spec.PostLoginJSON401Response(spec.Unauthorized{Feedback: "Usuário está inativo"})
	}

	// compara senhas
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return spec.PostLoginJSON400Response(spec.Error{Feedback: "E-mail ou senha inválidos"})
		}
		api.logger.Error("Falha comparar hash com senha", zap.String("password", string(credentials.Password)), zap.Error(err))
		return spec.PostLoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// recupera grupo de permissões
	group, err := api.store.GetPermissionsGroup(r.Context(), user.GroupID)
	if err != nil {
		api.logger.Error("Falha ao tentar buscar grupo por user.GroupId", zap.Error(err))
		return spec.PostLoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// mapeia grupos (json)
	groups := make(map[string]interface{})
	json.Unmarshal(group, &groups)

	// verifica se o grupo da aplicação está na lista de grupos do usuário
	if groups[credentials.Application] == nil {
		return spec.PostLoginJSON401Response(spec.Unauthorized{Feedback: "Usuário cadastrado, mas sem acesso à aplicação"})
	}

	// monta UUIDs
	groupId := uuid.MustParse(fmt.Sprintf("%+v", groups[credentials.Application]))
	applicationId := uuid.MustParse(credentials.Application)

	// busca informações da aplicação
	application, err := api.store.GetApplication(r.Context(), applicationId)
	if err != nil {
		api.logger.Error("Falha ao tentar buscar aplicação por applicationId", zap.Error(err))
		return spec.PostLoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// busca permissões do grupo na aplicação
	permissions, err := api.store.GetPermissionsGroup(r.Context(), groupId)
	if err != nil {
		api.logger.Error("Falha ao tentar buscar permissões por groupId", zap.Error(err))
		return spec.PostLoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// cria token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "authenticator",
		"sub":   user.Email,
		"aud":   credentials.Application,
		"exp":   time.Now().Add(time.Hour * 12).Unix(),
		"iat":   time.Now().Unix(),
		"roles": string(permissions[:]),
	})

	// assina token JWT usando a chave secreta da aplicação
	signedToken, err := token.SignedString([]byte(application.Secret.String()))
	if err != nil {
		api.logger.Error("Falha ao assinar token", zap.String("token", token.Raw), zap.Error(err))
		return spec.PostLoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// retorna token JWT para o cliente
	return spec.PostLoginJSON200Response(spec.LoginResponse{Token: signedToken, Feedback: "sessão iniciada"})
}
