package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/spec"
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
// (POST /login)
func (api API) PostLogin(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.Credentials

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return spec.PostLoginJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
	}

	if err := api.validator.validate.Struct(body); err != nil {
		return spec.PostLoginJSON400Response(spec.Error{Feedback: api.validator.Translate(err)})
	}

	user, err := api.store.GetUser(r.Context(), string(body.Email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spec.PostLoginJSON400Response(spec.Error{Feedback: "E-mail ou senha inválidos"})
		}
		api.logger.Error("Falha ao buscar usuário", zap.String("email", string(body.Email)), zap.Error(err))
		return spec.PostLoginJSON400Response(spec.Error{Feedback: "Falha ao tentar fazer login, tente novamente em alguns minutos"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return spec.PostLoginJSON400Response(spec.Error{Feedback: "E-mail ou senha inválidos"})
		}
		api.logger.Error("Falha comparar hash com senha", zap.String("password", string(body.Password)), zap.Error(err))
		return spec.PostLoginJSON400Response(spec.Error{Feedback: "Falha ao tentar fazer login, tente novamente em alguns minutos"})
	}

	secret := "implementar-depois"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "authenticator_ID",
		"sub": body.Email,
		"aud": "authenticator",
		"exp": time.Now().Add(time.Hour * 12).Unix(),
		"iat": time.Now().Unix(),
	})

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		api.logger.Error("Falha ao assinar token", zap.String("token", token.Raw), zap.Error(err))
		return spec.PostLoginJSON400Response(spec.Error{Feedback: "Falha ao tentar fazer login, tente novamente em alguns minutos"})
	}

	return spec.PostLoginJSON200Response(spec.LoginResponse{Token: signedToken})
}
