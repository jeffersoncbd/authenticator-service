// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.3.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/discord-gophers/goapi-gen/runtime"
	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for Status.
var (
	UnknownStatus = Status{}

	StatusActive = Status{"active"}

	StatusInactive = Status{"inactive"}
)

// Application defines model for Application.
type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BasicCreationResponse defines model for BasicCreationResponse.
type BasicCreationResponse struct {
	Feedback string `json:"feedback"`
	ID       string `json:"id"`
}

// BasicResponse defines model for BasicResponse.
type BasicResponse struct {
	Feedback string `json:"feedback"`
}

// Error defines model for Error.
type Error struct {
	Feedback string `json:"feedback"`
}

// InternalServerError defines model for InternalServerError.
type InternalServerError struct {
	Feedback string `json:"feedback"`
}

// LoginCredentials defines model for LoginCredentials.
type LoginCredentials struct {
	Application string              `json:"application" validate:"required"`
	Email       openapi_types.Email `json:"email" validate:"required"`
	Password    string              `json:"password" validate:"required,min=8,max=32"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Feedback string `json:"feedback"`
	Token    string `json:"token"`
}

// NewApplication defines model for NewApplication.
type NewApplication struct {
	Name string `json:"name" validate:"required"`
}

// NewGroup defines model for NewGroup.
type NewGroup struct {
	Name string `json:"name" validate:"required,min=3"`
}

// NewUser defines model for NewUser.
type NewUser struct {
	Email    openapi_types.Email `json:"email" validate:"required,email"`
	Name     string              `json:"name" validate:"required,min=3"`
	Password string              `json:"password" validate:"required,min=8,max=32"`
}

// NewUserStatus defines model for NewUserStatus.
type NewUserStatus struct {
	Status Status `json:"status"`
}

// Unauthorized defines model for Unauthorized.
type Unauthorized struct {
	Feedback string `json:"feedback"`
}

// User defines model for User.
type User struct {
	Email  openapi_types.Email `json:"email"`
	Name   string              `json:"name"`
	Status Status              `json:"status"`
}

// Status defines model for Status.
type Status struct {
	value string
}

func (t *Status) ToValue() string {
	return t.value
}
func (t Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}
func (t *Status) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return t.FromValue(value)
}
func (t *Status) FromValue(value string) error {
	switch value {

	case StatusActive.value:
		t.value = value
		return nil

	case StatusInactive.value:
		t.value = value
		return nil

	}
	return fmt.Errorf("unknown enum value: %v", value)
}

// PostApplicationsJSONBody defines parameters for PostApplications.
type PostApplicationsJSONBody NewApplication

// PostApplicationsIDGroupsJSONBody defines parameters for PostApplicationsIDGroups.
type PostApplicationsIDGroupsJSONBody NewGroup

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody LoginCredentials

// PostUsersJSONBody defines parameters for PostUsers.
type PostUsersJSONBody NewUser

// PatchUsersByEmailJSONBody defines parameters for PatchUsersByEmail.
type PatchUsersByEmailJSONBody NewUserStatus

// PostApplicationsJSONRequestBody defines body for PostApplications for application/json ContentType.
type PostApplicationsJSONRequestBody PostApplicationsJSONBody

// Bind implements render.Binder.
func (PostApplicationsJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostApplicationsIDGroupsJSONRequestBody defines body for PostApplicationsIDGroups for application/json ContentType.
type PostApplicationsIDGroupsJSONRequestBody PostApplicationsIDGroupsJSONBody

// Bind implements render.Binder.
func (PostApplicationsIDGroupsJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostLoginJSONRequestBody defines body for PostLogin for application/json ContentType.
type PostLoginJSONRequestBody PostLoginJSONBody

// Bind implements render.Binder.
func (PostLoginJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostUsersJSONRequestBody defines body for PostUsers for application/json ContentType.
type PostUsersJSONRequestBody PostUsersJSONBody

// Bind implements render.Binder.
func (PostUsersJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PatchUsersByEmailJSONRequestBody defines body for PatchUsersByEmail for application/json ContentType.
type PatchUsersByEmailJSONRequestBody PatchUsersByEmailJSONBody

// Bind implements render.Binder.
func (PatchUsersByEmailJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// GetApplicationsJSON200Response is a constructor method for a GetApplications response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsJSON200Response(body []Application) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetApplicationsJSON401Response is a constructor method for a GetApplications response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// GetApplicationsJSON500Response is a constructor method for a GetApplications response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// PostApplicationsJSON201Response is a constructor method for a PostApplications response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsJSON201Response(body BasicCreationResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostApplicationsJSON400Response is a constructor method for a PostApplications response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostApplicationsJSON401Response is a constructor method for a PostApplications response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// PostApplicationsJSON500Response is a constructor method for a PostApplications response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// GetApplicationsIDJSON200Response is a constructor method for a GetApplicationsID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsIDJSON200Response(body Application) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetApplicationsIDJSON400Response is a constructor method for a GetApplicationsID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsIDJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetApplicationsIDJSON401Response is a constructor method for a GetApplicationsID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsIDJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// GetApplicationsIDJSON500Response is a constructor method for a GetApplicationsID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsIDJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// GetApplicationsIDGroupsJSON200Response is a constructor method for a GetApplicationsIDGroups response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsIDGroupsJSON200Response(body BasicCreationResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetApplicationsIDGroupsJSON401Response is a constructor method for a GetApplicationsIDGroups response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsIDGroupsJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// GetApplicationsIDGroupsJSON500Response is a constructor method for a GetApplicationsIDGroups response.
// A *Response is returned with the configured status code and content type from the spec.
func GetApplicationsIDGroupsJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// PostApplicationsIDGroupsJSON201Response is a constructor method for a PostApplicationsIDGroups response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsIDGroupsJSON201Response(body BasicCreationResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostApplicationsIDGroupsJSON400Response is a constructor method for a PostApplicationsIDGroups response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsIDGroupsJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostApplicationsIDGroupsJSON401Response is a constructor method for a PostApplicationsIDGroups response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsIDGroupsJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// PostApplicationsIDGroupsJSON500Response is a constructor method for a PostApplicationsIDGroups response.
// A *Response is returned with the configured status code and content type from the spec.
func PostApplicationsIDGroupsJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// PostLoginJSON200Response is a constructor method for a PostLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLoginJSON200Response(body LoginResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PostLoginJSON400Response is a constructor method for a PostLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLoginJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostLoginJSON401Response is a constructor method for a PostLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLoginJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// PostLoginJSON500Response is a constructor method for a PostLogin response.
// A *Response is returned with the configured status code and content type from the spec.
func PostLoginJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// GetUsersJSON200Response is a constructor method for a GetUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func GetUsersJSON200Response(body []User) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetUsersJSON401Response is a constructor method for a GetUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func GetUsersJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// GetUsersJSON500Response is a constructor method for a GetUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func GetUsersJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// PostUsersJSON201Response is a constructor method for a PostUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUsersJSON201Response(body BasicResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostUsersJSON400Response is a constructor method for a PostUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUsersJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostUsersJSON401Response is a constructor method for a PostUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUsersJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// PostUsersJSON500Response is a constructor method for a PostUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUsersJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// PatchUsersByEmailJSON200Response is a constructor method for a PatchUsersByEmail response.
// A *Response is returned with the configured status code and content type from the spec.
func PatchUsersByEmailJSON200Response(body BasicResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// PatchUsersByEmailJSON400Response is a constructor method for a PatchUsersByEmail response.
// A *Response is returned with the configured status code and content type from the spec.
func PatchUsersByEmailJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PatchUsersByEmailJSON500Response is a constructor method for a PatchUsersByEmail response.
// A *Response is returned with the configured status code and content type from the spec.
func PatchUsersByEmailJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Lista todas as aplicações
	// (GET /applications)
	GetApplications(w http.ResponseWriter, r *http.Request) *Response
	// Cadastra uma aplicação
	// (POST /applications)
	PostApplications(w http.ResponseWriter, r *http.Request) *Response
	// Informações de uma aplicação
	// (GET /applications/{id})
	GetApplicationsID(w http.ResponseWriter, r *http.Request, id string) *Response
	// Lista os grupos de permissões de um aplicativo
	// (GET /applications/{id}/groups)
	GetApplicationsIDGroups(w http.ResponseWriter, r *http.Request, id string) *Response
	// Cadastra um novo grupo de permissões para um aplicativo
	// (POST /applications/{id}/groups)
	PostApplicationsIDGroups(w http.ResponseWriter, r *http.Request, id string) *Response
	// Autentica usuário
	// (POST /login)
	PostLogin(w http.ResponseWriter, r *http.Request) *Response
	// Lista todos os usuários
	// (GET /users)
	GetUsers(w http.ResponseWriter, r *http.Request) *Response
	// Cadastra um novo usuário
	// (POST /users)
	PostUsers(w http.ResponseWriter, r *http.Request) *Response
	// Atualiza o status de um usuário
	// (PATCH /users/{byEmail})
	PatchUsersByEmail(w http.ResponseWriter, r *http.Request, byEmail openapi_types.Email) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// GetApplications operation middleware
func (siw *ServerInterfaceWrapper) GetApplications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetApplications(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostApplications operation middleware
func (siw *ServerInterfaceWrapper) PostApplications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostApplications(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetApplicationsID operation middleware
func (siw *ServerInterfaceWrapper) GetApplicationsID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetApplicationsID(w, r, id)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetApplicationsIDGroups operation middleware
func (siw *ServerInterfaceWrapper) GetApplicationsIDGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetApplicationsIDGroups(w, r, id)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostApplicationsIDGroups operation middleware
func (siw *ServerInterfaceWrapper) PostApplicationsIDGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostApplicationsIDGroups(w, r, id)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostLogin operation middleware
func (siw *ServerInterfaceWrapper) PostLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostLogin(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetUsers operation middleware
func (siw *ServerInterfaceWrapper) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetUsers(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostUsers operation middleware
func (siw *ServerInterfaceWrapper) PostUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostUsers(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PatchUsersByEmail operation middleware
func (siw *ServerInterfaceWrapper) PatchUsersByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "byEmail" -------------
	var byEmail openapi_types.Email

	if err := runtime.BindStyledParameter("simple", false, "byEmail", chi.URLParam(r, "byEmail"), &byEmail); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "byEmail"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PatchUsersByEmail(w, r, byEmail)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:    "/",
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Get("/applications", wrapper.GetApplications)
		r.Post("/applications", wrapper.PostApplications)
		r.Get("/applications/{id}", wrapper.GetApplicationsID)
		r.Get("/applications/{id}/groups", wrapper.GetApplicationsIDGroups)
		r.Post("/applications/{id}/groups", wrapper.PostApplicationsIDGroups)
		r.Post("/login", wrapper.PostLogin)
		r.Get("/users", wrapper.GetUsers)
		r.Post("/users", wrapper.PostUsers)
		r.Patch("/users/{byEmail}", wrapper.PatchUsersByEmail)
	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xZzW7jNhB+FYHtUYm9mxYoBOzBSdMgxWIRNBv0kM1hIo1tbiRSS1JOHEMPU/SwT9An",
	"8IsVJGVZUiRbju04aX1JZIo/w5nvm/lITYjPo5gzZEoSb0KkP8QIzGMvjkPqg6Kc6Z8QBFQ/Q3gheIxC",
	"UZTE60Mo0SVxoWlCaKD/9rmIQBGPJAkNiEvUOEbiEakEZQOSuoRBhLpj5UXqEoHfEiowIN41MWNN15t8",
	"Dn77FX2l5zgGSf0TgcbMP1DGnElc0dw+YnAL/l2NLW67vdSZnM/aaPbmza3YsdCEUyG42M3S50yhYBBe",
	"ohih2KEhH/mAshOBATJFIZQrWgFlhiwGiUseDgb8AB+UgAMFAzPDCEIagNLdcqP1DjACGpbmtC1rTRqD",
	"lPdclAGdNz53ajei7MMvbgQPH47ek7Tq/6KP3HwX+aqNUdkGlxW/Q7YcMrbbEgZ/wvvnJ8j6zLdCMCsW",
	"N6bHT3h/JngSv6h9BhFHDVa6RUScB01WX0lcNSdsmjOuHZ42VqqV/bEjCmZ+b8W9zPWXClSyajqU+aAf",
	"BfaJR37ozMVFJ1MWnWzqqpHZ4Dqb5sYgSyKTVHxFR3pPlGWPNzX64opBooZc0EcMdlNftojiZv3krheI",
	"Cloa46LXQT8RVI0v9YzW/FsEgaKXqOH812+znfz+52c9o+lNvOztfGdDpWKLXsr6XI8PUPqCxjbDklMZ",
	"o0/71Ifp9+k/KJ0AnN7FuRODAIc7OgoHyALdDCbFTL9P/+bOF6LN0RXeB8XFF3Kol6Qq1GuWXhGXjFBI",
	"u9q7w+5hV3uTx8ggpsQjR6ZJM0gNzW47hVxmGgao9D8d0FmC88gZql6xn/a3LW5mzPtuV//zOVPIVEVX",
	"dL5KW11s1Iy6VhgtDW6xNqW5h0EIGFsXl137kUoFTgBy7jntYB8CkEpAwKWe5afuu5UsXWRgiZw1FpXf",
	"u+TnFb20aO06/VljwqybY/s5s45z4BPvugz565v0xiUyiSIQ49yvihvPlp2rQWhSekkjSXKjywSXNTi6",
	"4PIpkL4lKNUxD8Yb805F2VTKiRIJpk8QvDlc1B/oaqLTK1A8xylYmG4OKo3gOIbAyby/p8ZzqHGSxcxJ",
	"olK+buZF6pYTbmdCg7Rt1j3/1SRuAREqFNKYR5lRX2o4u2Hw7NG9jHa34LqldxotBZsZnmqHrFUKWleA",
	"ugiavcwraTEEexK9ERJVoohrkqkz0IdF2Z5TZ7b//4tZrUuUrf4DkcR7/bSOfuIyc6JGeIwiolLOAZ/h",
	"XdFREe4ZkturqdcO5q0IPXs59Fol3pkOeuEQsq9Lb1DcOYyPuOVvlb7m4L6UwLpShXxAzUabuWxujbd0",
	"JHrynaAVY7qbXX8RUz7zO2TavZAoe52xF3KvgjA5I3qzwDiJTKZ/CVrEukRpbp0s2BNpKk+zBrsyHV7i",
	"GslcYLa/P+Iy393+8mhzl0dcOkXPFoBjobJY5MzRshUJYSGyAwWxKB9eZb7ai4f/gnioyZgz4Of5sjO5",
	"HZ9GQENzIxSD8oc1dNDNhg/HtnMrsX+b922h+Df0zXGLmr/wfe+FZcxS2pp3OumN4JHu7C73bbKmpxII",
	"6SM43LHf7LLT8SLypOm/AQAA//90PDXTACYAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
