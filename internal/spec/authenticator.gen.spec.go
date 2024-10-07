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

// Group defines model for Group.
type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	// Este objeto é composto pela definição de permissões onde a chave é o identificador do recurso e o valor é o número da permissão
	Permissions map[string]interface{} `json:"permissions"`
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

// N400 defines model for 400.
type N400 Error

// N401 defines model for 401.
type N401 Unauthorized

// N500 defines model for 500.
type N500 InternalServerError

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

// NewApplicationJSONBody defines parameters for NewApplication.
type NewApplicationJSONBody NewApplication

// NewGroupJSONBody defines parameters for NewGroup.
type NewGroupJSONBody NewGroup

// NewUserJSONBody defines parameters for NewUser.
type NewUserJSONBody NewUser

// FindUserByEmailJSONBody defines parameters for FindUserByEmail.
type FindUserByEmailJSONBody NewUserStatus

// LoginJSONBody defines parameters for Login.
type LoginJSONBody LoginCredentials

// NewApplicationJSONRequestBody defines body for NewApplication for application/json ContentType.
type NewApplicationJSONRequestBody NewApplicationJSONBody

// Bind implements render.Binder.
func (NewApplicationJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// NewGroupJSONRequestBody defines body for NewGroup for application/json ContentType.
type NewGroupJSONRequestBody NewGroupJSONBody

// Bind implements render.Binder.
func (NewGroupJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// NewUserJSONRequestBody defines body for NewUser for application/json ContentType.
type NewUserJSONRequestBody NewUserJSONBody

// Bind implements render.Binder.
func (NewUserJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// FindUserByEmailJSONRequestBody defines body for FindUserByEmail for application/json ContentType.
type FindUserByEmailJSONRequestBody FindUserByEmailJSONBody

// Bind implements render.Binder.
func (FindUserByEmailJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody LoginJSONBody

// Bind implements render.Binder.
func (LoginJSONRequestBody) Bind(*http.Request) error {
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

// ApplicationsListJSON200Response is a constructor method for a ApplicationsList response.
// A *Response is returned with the configured status code and content type from the spec.
func ApplicationsListJSON200Response(body []Application) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// ApplicationsListJSON401Response is a constructor method for a ApplicationsList response.
// A *Response is returned with the configured status code and content type from the spec.
func ApplicationsListJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// ApplicationsListJSON500Response is a constructor method for a ApplicationsList response.
// A *Response is returned with the configured status code and content type from the spec.
func ApplicationsListJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// NewApplicationJSON201Response is a constructor method for a NewApplication response.
// A *Response is returned with the configured status code and content type from the spec.
func NewApplicationJSON201Response(body BasicCreationResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// NewApplicationJSON400Response is a constructor method for a NewApplication response.
// A *Response is returned with the configured status code and content type from the spec.
func NewApplicationJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// NewApplicationJSON401Response is a constructor method for a NewApplication response.
// A *Response is returned with the configured status code and content type from the spec.
func NewApplicationJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// NewApplicationJSON500Response is a constructor method for a NewApplication response.
// A *Response is returned with the configured status code and content type from the spec.
func NewApplicationJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// FindApplicationByIDJSON200Response is a constructor method for a FindApplicationByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindApplicationByIDJSON200Response(body Application) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// FindApplicationByIDJSON400Response is a constructor method for a FindApplicationByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindApplicationByIDJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// FindApplicationByIDJSON401Response is a constructor method for a FindApplicationByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindApplicationByIDJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// FindApplicationByIDJSON500Response is a constructor method for a FindApplicationByID response.
// A *Response is returned with the configured status code and content type from the spec.
func FindApplicationByIDJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// GroupsListJSON200Response is a constructor method for a GroupsList response.
// A *Response is returned with the configured status code and content type from the spec.
func GroupsListJSON200Response(body []Group) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GroupsListJSON400Response is a constructor method for a GroupsList response.
// A *Response is returned with the configured status code and content type from the spec.
func GroupsListJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GroupsListJSON401Response is a constructor method for a GroupsList response.
// A *Response is returned with the configured status code and content type from the spec.
func GroupsListJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// GroupsListJSON500Response is a constructor method for a GroupsList response.
// A *Response is returned with the configured status code and content type from the spec.
func GroupsListJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// NewGroupJSON201Response is a constructor method for a NewGroup response.
// A *Response is returned with the configured status code and content type from the spec.
func NewGroupJSON201Response(body BasicCreationResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// NewGroupJSON400Response is a constructor method for a NewGroup response.
// A *Response is returned with the configured status code and content type from the spec.
func NewGroupJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// NewGroupJSON401Response is a constructor method for a NewGroup response.
// A *Response is returned with the configured status code and content type from the spec.
func NewGroupJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// NewGroupJSON500Response is a constructor method for a NewGroup response.
// A *Response is returned with the configured status code and content type from the spec.
func NewGroupJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// UsersListJSON200Response is a constructor method for a UsersList response.
// A *Response is returned with the configured status code and content type from the spec.
func UsersListJSON200Response(body []User) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// UsersListJSON400Response is a constructor method for a UsersList response.
// A *Response is returned with the configured status code and content type from the spec.
func UsersListJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// UsersListJSON401Response is a constructor method for a UsersList response.
// A *Response is returned with the configured status code and content type from the spec.
func UsersListJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// UsersListJSON500Response is a constructor method for a UsersList response.
// A *Response is returned with the configured status code and content type from the spec.
func UsersListJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// NewUserJSON201Response is a constructor method for a NewUser response.
// A *Response is returned with the configured status code and content type from the spec.
func NewUserJSON201Response(body BasicResponse) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// NewUserJSON400Response is a constructor method for a NewUser response.
// A *Response is returned with the configured status code and content type from the spec.
func NewUserJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// NewUserJSON401Response is a constructor method for a NewUser response.
// A *Response is returned with the configured status code and content type from the spec.
func NewUserJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// NewUserJSON500Response is a constructor method for a NewUser response.
// A *Response is returned with the configured status code and content type from the spec.
func NewUserJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// FindUserByEmailJSON200Response is a constructor method for a FindUserByEmail response.
// A *Response is returned with the configured status code and content type from the spec.
func FindUserByEmailJSON200Response(body BasicResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// FindUserByEmailJSON400Response is a constructor method for a FindUserByEmail response.
// A *Response is returned with the configured status code and content type from the spec.
func FindUserByEmailJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// FindUserByEmailJSON401Response is a constructor method for a FindUserByEmail response.
// A *Response is returned with the configured status code and content type from the spec.
func FindUserByEmailJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// FindUserByEmailJSON500Response is a constructor method for a FindUserByEmail response.
// A *Response is returned with the configured status code and content type from the spec.
func FindUserByEmailJSON500Response(body InternalServerError) *Response {
	return &Response{
		body:        body,
		Code:        500,
		contentType: "application/json",
	}
}

// LoginJSON200Response is a constructor method for a Login response.
// A *Response is returned with the configured status code and content type from the spec.
func LoginJSON200Response(body LoginResponse) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// LoginJSON400Response is a constructor method for a Login response.
// A *Response is returned with the configured status code and content type from the spec.
func LoginJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// LoginJSON401Response is a constructor method for a Login response.
// A *Response is returned with the configured status code and content type from the spec.
func LoginJSON401Response(body Unauthorized) *Response {
	return &Response{
		body:        body,
		Code:        401,
		contentType: "application/json",
	}
}

// LoginJSON500Response is a constructor method for a Login response.
// A *Response is returned with the configured status code and content type from the spec.
func LoginJSON500Response(body InternalServerError) *Response {
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
	ApplicationsList(w http.ResponseWriter, r *http.Request) *Response
	// Cadastra uma aplicação
	// (POST /applications)
	NewApplication(w http.ResponseWriter, r *http.Request) *Response
	// Todas as informações de uma aplicação
	// (GET /applications/{id})
	FindApplicationByID(w http.ResponseWriter, r *http.Request, id string) *Response
	// Lista os grupos de permissões de uma aplicação
	// (GET /applications/{id}/groups)
	GroupsList(w http.ResponseWriter, r *http.Request, id string) *Response
	// Cadastra um novo grupo de permissões de uma aplicação
	// (POST /applications/{id}/groups)
	NewGroup(w http.ResponseWriter, r *http.Request, id string) *Response
	// Lista os usuários de uma aplicação
	// (GET /applications/{id}/users)
	UsersList(w http.ResponseWriter, r *http.Request, id string) *Response
	// Cadastra um novo usuário
	// (POST /applications/{id}/users)
	NewUser(w http.ResponseWriter, r *http.Request, id string) *Response
	// Atualiza o status de um usuário
	// (PATCH /applications/{id}/users/{byEmail})
	FindUserByEmail(w http.ResponseWriter, r *http.Request, id string, byEmail openapi_types.Email) *Response
	// Autentica usuário
	// (POST /login)
	Login(w http.ResponseWriter, r *http.Request) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// ApplicationsList operation middleware
func (siw *ServerInterfaceWrapper) ApplicationsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.ApplicationsList(w, r)
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

// NewApplication operation middleware
func (siw *ServerInterfaceWrapper) NewApplication(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.NewApplication(w, r)
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

// FindApplicationByID operation middleware
func (siw *ServerInterfaceWrapper) FindApplicationByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.FindApplicationByID(w, r, id)
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

// GroupsList operation middleware
func (siw *ServerInterfaceWrapper) GroupsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GroupsList(w, r, id)
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

// NewGroup operation middleware
func (siw *ServerInterfaceWrapper) NewGroup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.NewGroup(w, r, id)
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

// UsersList operation middleware
func (siw *ServerInterfaceWrapper) UsersList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.UsersList(w, r, id)
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

// NewUser operation middleware
func (siw *ServerInterfaceWrapper) NewUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.NewUser(w, r, id)
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

// FindUserByEmail operation middleware
func (siw *ServerInterfaceWrapper) FindUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "id" -------------
	var id string

	if err := runtime.BindStyledParameter("simple", false, "id", chi.URLParam(r, "id"), &id); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "id"})
		return
	}

	// ------------- Path parameter "byEmail" -------------
	var byEmail openapi_types.Email

	if err := runtime.BindStyledParameter("simple", false, "byEmail", chi.URLParam(r, "byEmail"), &byEmail); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "byEmail"})
		return
	}

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{""})

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.FindUserByEmail(w, r, id, byEmail)
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

// Login operation middleware
func (siw *ServerInterfaceWrapper) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.Login(w, r)
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
		r.Get("/applications", wrapper.ApplicationsList)
		r.Post("/applications", wrapper.NewApplication)
		r.Get("/applications/{id}", wrapper.FindApplicationByID)
		r.Get("/applications/{id}/groups", wrapper.GroupsList)
		r.Post("/applications/{id}/groups", wrapper.NewGroup)
		r.Get("/applications/{id}/users", wrapper.UsersList)
		r.Post("/applications/{id}/users", wrapper.NewUser)
		r.Patch("/applications/{id}/users/{byEmail}", wrapper.FindUserByEmail)
		r.Post("/login", wrapper.Login)
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

	"H4sIAAAAAAAC/+RZzW7bOBB+FYK7RyV2my6wMNBDkqaLLIqi2KTYQ5rDRBzbbCRSS1JOnEAPs+ih2Afo",
	"Za9+sYKkLEsO5Z/E+WuBopElivxm5uM3M9Q1jWWaSYHCaNq7pgp1JoVG9+NVt2v/xFIYFMZeQpYlPAbD",
	"peh81lLYezoeYgr26leFfdqjv3Rmc3b8U905UEoqWhRFRBnqWPHMTkJ7dA8YUfhPjtrQIqKvui82tuZH",
	"AbkZSsWvkIWWbj6P6G8btPdQGFQCkiNUI1St1k+HET+OTAdG5TIuDrszFA4UY9xeQ/JByQyV4TZcfUg0",
	"RjSr3bqmnNn/+1KlYGiP5jlnNKJmnCHtUW0UFwNruIAU7cC5B0VEbWS4QkZ7J9S964aeVnPIs88Yu8Dt",
	"gebxvkIH86+SRmvC7SOyM4jPA1ii1WwJQa5mbYW9ebhzOBZC8DF/lKX/UDLPHpRTEc1QpVxrLoWbrrkf",
	"DrRBYvEZSSb/EbextJEkwwQIwz4XfPJ18kUShqScaPINNZGCIQESD2GE9kVJOENheJ/HwKQiTBKFca60",
	"JEgkGUEilR8nJv+nqCRhUE34RdIbvmrdC02DQk4OScGjRPudHHCxr9B5BhK9JgpoytBiBkT0cmsgt/DS",
	"KNgyMHAzjCDhDIwdVoG2FmAKPGnM6e/cadIMtL6QqsnW6uZtp45SLl7/HqVw+XrnpVf0uv/rPooqK6pV",
	"W6NyH4Jp5DmK5ZTxw5bI5Hu8uH0WCkvBGsGcQ9yag97jxW0U7W74HCN2WlBGdUYcsjbUHzWuqwmb3jOR",
	"f71ole61/fFIW7D0+0p7r3T9kQGTryuHunppUR1YTj0Psnw5hGkGBkWeOlGJDR9Zm7goL08DCbdRzT5K",
	"frlHFi8qKO4UiDm2tMbFrmMLCG7GR3ZGD/8MQaHazc1w9uvt1JI//z6mZR1vZ/JPZ5YNjck8e7noy1Ap",
	"lGHsCpjJV1fjMCC7Hw5JBgqIJDYKWyiYvQ1OYnxh9IlaODbDx2Ck+kS37ZLcJHbNxiMa0REq7Vd7sd3d",
	"7lpvygwFZJz26I67ZXeQGTprOzUtczcG6LokG9CpwPXq3Yp+x7WhUbOrfLlml8UNpkujW09OReViUArG",
	"obbL4rL15Mxz31DX2s/QUpURHTto1i4uHmsH1clDeydN2pycFqcR1XmaghpX0IxkoIn9VwcYUS+LjTpD",
	"01MrtVIHYjGXtj31UZs9ycYba3XnFpkTZKNyLG5QYHNNfrjvDLXab+a3SgwMtFHAwIe+u0rou0+DJvsl",
	"dpKnDaPaOVJEzQ3cueasaN3Fb7lgtajujQ/fOClQkKJBpR1ALlw+N8NpM9TzjVEz+lEtkkvbxhVLAPd6",
	"YV1yJ21ZWVJCJzfOlpk214Pw7Oh0PNUb3jQL78ivzsDW4+3JwpXr0zTxo7FrpczlG5aVc9ZA5ZnUc8cf",
	"z49v3hipw/Ys5l1JqYVZz3v1iVLqXlJwyaMnnHylD/Ys7cpnnXaJkKOpSbeib1gwc+3I2qKXtsX6qeXS",
	"9ZjrVPi5zif/Ki6fs0hWNixmlqfOQl107vuZZNHz5RFUcbkaTqP6Ywni1KoALxcIXuf6bHyQAk9cR5KB",
	"iYfhnsTGc88PfYI8joIYziq8KwDZ0CnuPW+p6jRthY3VfbiN5Z5Z0RzBFX+GXf2uySHhV0Ak8SeRXvCX",
	"bapEDrhzY1j33Remezr8ufFN8YEp0fx6FqDEsTxH4Y76cuOPPp9Yiz4L/hRgKNwa3RdmG/Ci+B4AAP//",
	"8fd+J7YjAAA=",
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
