package api

import (
	"authenticator/internal/databases/postgresql"
	"authenticator/internal/spec"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/go-playground/locales/pt_BR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	pt_br_translations "github.com/go-playground/validator/v10/translations/pt_BR"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Validator struct {
	validate	*validator.Validate
	translator	ut.Translator
}
func (v Validator) Translate(err error) string {
	var errorMessages []string
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		errorMessages = append(errorMessages, e.Translate(v.translator))
	}
	return strings.Join(errorMessages, ", ")
}

type API struct{
	store		*postgresql.Queries
	logger		*zap.Logger
	validator	Validator
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
        store: postgresql.New(pool),
        logger: logger,
		validator: Validator{
			validate: validate,
			translator: translator,
		},
    }
}

// Autentica usuário
// (POST /login)
func (api API) PostLogin(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.Credentials

    err := json.NewDecoder(r.Body).Decode(&body)
    if err!= nil {
        return spec.PostLoginJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
    }

    if err := api.validator.validate.Struct(body); err!= nil {
        return spec.PostLoginJSON400Response(spec.Error{Feedback: api.validator.Translate(err)})
    }

    user, err := api.store.GetUser(r.Context(), string(body.Email))
    if err!= nil {
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
	if err!= nil {
        api.logger.Error("Falha ao assinar token", zap.String("token", token.Raw), zap.Error(err))
        return spec.PostLoginJSON400Response(spec.Error{Feedback: "Falha ao tentar fazer login, tente novamente em alguns minutos"})
    }

	return spec.PostLoginJSON200Response(spec.LoginResponse{Token:  signedToken})
}

// Lista todos os usuários
// (GET /users)
func (api API) GetUsers(w http.ResponseWriter, r *http.Request) *spec.Response {
	rows, err := api.store.ListUsers(r.Context())
    if err!= nil {
        api.logger.Error("Falha ao listar usuários", zap.Error(err))
        return spec.GetUsersJSON500Response(spec.Error{Feedback: "Falha ao listar usuários, tente novamente em alguns minutos..."})
    }

    var users []spec.UserData
    for _, row := range rows {
		status := spec.UserStatus{}
		status.FromValue(row.Status.String)
        users = append(users, spec.UserData{
            Name:  row.Name,
            Email: openapi_types.Email(row.Email),
			Status: status,
        })
    }

    return spec.GetUsersJSON200Response(users)
}

// Cadastra um novo usuário
// (POST /users)
func (api API) PostUsers(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.User

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
        return spec.PostUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
    }

	if err := api.validator.validate.Struct(body); err != nil {
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
	}

	_, err = api.store.GetUser(r.Context(), string(body.Email))
	if err == nil {
        return spec.PostUsersJSON400Response(spec.Error{Feedback: "Já existe usuário cadastrado com este e-mail"})
    }
	if !errors.Is(err, pgx.ErrNoRows) {
		api.logger.Error("Falha ao consultar email", zap.Error(err), zap.String("email", string(body.Email)))
		return spec.PostUsersJSON400Response(spec.Error{Feedback: "Falha no cadastro, tente novamente em alguns minutos"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err!= nil {
        api.logger.Error("Falha ao gerar hash de senha", zap.Error(err))
        return spec.PostUsersJSON400Response(spec.Error{Feedback: "Falha no cadastro, tente novamente em alguns minutos"})
    }
	err = api.store.InsertUser(r.Context(), postgresql.InsertUserParams{
		Email:    string(body.Email),
		Name:     body.Name,
		Password: string(hash),
	})
	if err!= nil {
        api.logger.Error("Falha ao inserir novo usuário", zap.Error(err))
        return spec.PostUsersJSON400Response(spec.Error{Feedback: "Falha no cadastro, tente novamente em alguns minutos"})
    }

	return nil
}

// Atualiza o status de um usuário
// (PATCH /users/{byEmail})
func (api API) PatchUsersByEmail(w http.ResponseWriter, r *http.Request, byEmail openapi_types.Email) *spec.Response {
	var body spec.PatchUserStatus

    err := json.NewDecoder(r.Body).Decode(&body)
    if err!= nil {
        return spec.PatchUsersByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + err.Error()})
    }

    if err := api.validator.validate.Struct(body); err!= nil {
        return spec.PatchUsersByEmailJSON400Response(spec.Error{Feedback: "Dados inválidos: " + api.validator.Translate(err)})
    }

    if err := api.store.UpdateUserStatus(r.Context(), postgresql.UpdateUserStatusParams{
        Email:  string(byEmail),
        Status: pgtype.Text{String: body.Status.ToValue(), Valid: true},
    }); err != nil {
        api.logger.Error("Falha ao atualizar status do usuário", zap.Error(err))
	}

	return nil
}
