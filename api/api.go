package api

import (
	postgresql "authenticator/interfaces"
	"authenticator/spec"
	"encoding/json"
	"errors"
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
	ptBr := pt_BR.New()
	universalTranslator := ut.New(ptBr, ptBr)

	translator, err := universalTranslator.GetTranslator("pt_BR")
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
func (api API) Login(w http.ResponseWriter, r *http.Request) *spec.Response {
	// decodifica body armazenando as credenciais
	var credentials spec.LoginCredentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		return spec.LoginJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	// valida dados de entrada
	if err := api.validator.validate.Struct(credentials); err != nil {
		return spec.LoginJSON400Response(spec.Error{Feedback: api.validator.Translate(err)})
	}

	// valida UUID
	applicationUUID, err := uuid.Parse(credentials.Application)
	if err != nil {
		return spec.LoginJSON400Response(spec.Error{Feedback: "ID da aplicação inválido: " + err.Error()})
	}

	// busca usuário no banco de dados
	user, err := api.store.GetUser(r.Context(), postgresql.GetUserParams{
		ApplicationID: applicationUUID,
		Email:         string(credentials.Email),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.LoginJSON400Response(spec.Error{Feedback: "E-mail ou senha inválidos"})
		}
		api.logger.Error("Falha ao buscar usuário", zap.String("email", string(credentials.Email)), zap.Error(err))
		return spec.LoginJSON500Response(spec.InternalServerError{Feedback: INTERNAL_SERVER_ERROR})
	}

	// verifica status
	if user.Status == "inactive" {
		return spec.LoginJSON401Response(spec.Unauthorized{Feedback: "Usuário está inativo"})
	}

	// compara senhas
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return spec.LoginJSON400Response(spec.Error{Feedback: "E-mail ou senha inválidos"})
		}
		api.logger.Error("Falha comparar hash com senha", zap.String("password", string(credentials.Password)), zap.Error(err))
		return spec.LoginJSON500Response(spec.InternalServerError{Feedback: INTERNAL_SERVER_ERROR})
	}

	// monta UUID
	applicationUUID = uuid.MustParse(credentials.Application)

	// busca informações da aplicação
	application, err := api.store.GetApplication(r.Context(), applicationUUID)
	if err != nil {
		api.logger.Error("Falha ao tentar buscar aplicação por applicationId", zap.Error(err))
		return spec.LoginJSON500Response(spec.InternalServerError{Feedback: INTERNAL_SERVER_ERROR})
	}

	// busca permissões do grupo na aplicação
	permissions, err := api.store.GetPermissionsGroup(r.Context(), user.GroupID)
	if err != nil {
		api.logger.Error("Falha ao tentar buscar permissões por groupId", zap.Error(err))
		return spec.LoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// cria token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   "authenticator",
		"sub":   user.Email,
		"aud":   applicationUUID,
		"exp":   time.Now().Add(time.Hour * 12).Unix(),
		"iat":   time.Now().Unix(),
		"roles": string(permissions[:]),
	})

	// assina token JWT usando a chave secreta da aplicação
	signedToken, err := token.SignedString([]byte(application.Secret.String()))
	if err != nil {
		api.logger.Error("Falha ao assinar token", zap.String("token", token.Raw), zap.Error(err))
		return spec.LoginJSON500Response(spec.InternalServerError{Feedback: "internal server error"})
	}

	// retorna token JWT para o cliente
	return spec.LoginJSON200Response(spec.LoginResponse{Token: signedToken, Feedback: "Sessão iniciada"})
}
