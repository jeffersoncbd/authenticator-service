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

// Defines values for InternalServerErrorFeedback.
var (
	UnknownInternalServerErrorFeedback = InternalServerErrorFeedback{}

	InternalServerErrorFeedbackInternalServerError = InternalServerErrorFeedback{"internal server error"}
)

// Defines values for UserStatus.
var (
	UnknownUserStatus = UserStatus{}

	UserStatusActive = UserStatus{"active"}

	UserStatusInactive = UserStatus{"inactive"}
)

// Credentials defines model for Credentials.
type Credentials struct {
	Application string              `json:"application" validate:"required"`
	Email       openapi_types.Email `json:"email" validate:"required"`
	Password    string              `json:"password" validate:"required,min=8,max=32"`
}

// Error defines model for Error.
type Error struct {
	Feedback string `json:"feedback"`
}

// InternalServerError defines model for InternalServerError.
type InternalServerError struct {
	Feedback *InternalServerErrorFeedback `json:"feedback,omitempty"`
}

// LoginResponse defines model for LoginResponse.
type LoginResponse struct {
	Token string `json:"token"`
}

// PatchUserStatus defines model for PatchUserStatus.
type PatchUserStatus struct {
	Status UserStatus `json:"status"`
}

// User defines model for User.
type User struct {
	Email    openapi_types.Email `json:"email" validate:"required,email"`
	Name     string              `json:"name" validate:"required,min=3"`
	Password string              `json:"password" validate:"required,min=8,max=32"`
}

// UserData defines model for UserData.
type UserData struct {
	Email  openapi_types.Email `json:"email"`
	Name   string              `json:"name"`
	Status UserStatus          `json:"status"`
}

// InternalServerErrorFeedback defines model for InternalServerError.Feedback.
type InternalServerErrorFeedback struct {
	value string
}

func (t *InternalServerErrorFeedback) ToValue() string {
	return t.value
}
func (t InternalServerErrorFeedback) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}
func (t *InternalServerErrorFeedback) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return t.FromValue(value)
}
func (t *InternalServerErrorFeedback) FromValue(value string) error {
	switch value {

	case InternalServerErrorFeedbackInternalServerError.value:
		t.value = value
		return nil

	}
	return fmt.Errorf("unknown enum value: %v", value)
}

// UserStatus defines model for UserStatus.
type UserStatus struct {
	value string
}

func (t *UserStatus) ToValue() string {
	return t.value
}
func (t UserStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}
func (t *UserStatus) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return t.FromValue(value)
}
func (t *UserStatus) FromValue(value string) error {
	switch value {

	case UserStatusActive.value:
		t.value = value
		return nil

	case UserStatusInactive.value:
		t.value = value
		return nil

	}
	return fmt.Errorf("unknown enum value: %v", value)
}

// PostLoginJSONBody defines parameters for PostLogin.
type PostLoginJSONBody Credentials

// PostUsersJSONBody defines parameters for PostUsers.
type PostUsersJSONBody User

// PatchUsersByEmailJSONBody defines parameters for PatchUsersByEmail.
type PatchUsersByEmailJSONBody PatchUserStatus

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
func GetUsersJSON200Response(body []UserData) *Response {
	return &Response{
		body:        body,
		Code:        200,
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

// PostUsersJSON400Response is a constructor method for a PostUsers response.
// A *Response is returned with the configured status code and content type from the spec.
func PostUsersJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
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

	"H4sIAAAAAAAC/+xWzW4btxN/FWL+/+PacuIWKBbIwU7dwkUORp2gB0WH0XIkMdklt+SsakXQwxQ95An6",
	"BHqxYkhp9eF1bdUK2kNOWlLD+fzNb2YOhatqZ8lygHwOoZhQhfHztSdNlg2W8YhaGzbOYnnjXU2eDQXI",
	"R1gGyqDeupoD1nVpChRpOY6cr5Ahh6YxGjLgWU2QQ2Bv7BgyuDsZuxO6Y48njOOoYYql0cgi5unXxnjS",
	"sFhkQBWackdnunmW0hpD+M15vaO3vfynqrPK2FffZRXevTp/CQsx1FrN+zs5ytooWquD1qwbfqCCYZHB",
	"lffOH1iKEZEeYvFRvnfj2HeolewyfW2ZvMXylvyUfOvIw7bINpVoNauHKsSXiuLTQdbhzD2jb9zY2J8p",
	"1M4GOjBudh/JPh50EuuK+Aa5mLwL5G8ZuTm0A0L76P+eRpDD/3qbPuutmqy3pX7fsZWCLs/k1YHuHLtt",
	"svRcnLZY0f08H9Qk5/9eF0bvn9h+kvfvkfEL5f7BXC6y58NpL85H0LWB/LqNsWAzFQ3Grj47WzhQ0XjD",
	"s1vxKIU/JPTkLxqebE4/rDPx0y9vxZsoDfnq301mJsx1qpmxIyfvNYXCmzpNFrgKNRVmZApcfl7+SUFp",
	"VBc316pGj8op4aITslquMbLt8vPyD6feg7gjo61Adv49nIpJw6XY3PkLMpiSD8nai9Oz0zNJkavJYm0g",
	"h/N4JbjhSYy2VwpnRWp0geVXoBBZ/lpDDjcucKQ1SPWhwJdOz0SwcJbJ8t787H0IaYimIj8Gge2ZvQd2",
	"9g3Fi8Sn0d2XZ2dHM73L1tH4brXeCtcqTQobThmO5ZCMfnNEP9Jw6rB/iVqtci42vz2iza7x2OHBWkwl",
	"ObUWzCA0VYV+lvCXkqOa0Cx/98YJOCPB9SFQiGAcyJteE8jHMo6pA2o/Er+LAs+suWGqnkQ9kRs3cxy9",
	"x1lXFt6YwKi0C22EQRWoMbBH7cJ/rzYrUoO8v0tn/cFisF26FBg7CW07uq0CppINZNw9yA+bqh2fH+Lu",
	"8CRieHGfbmNzS4hT/GTwa98ego3XK4CrplLWTV1Xd6/B0fZ2bz6cXcnIXsSJIhtpB2TWi2q4TMJxHnms",
	"iCM99OcgEynOKFjvGDBsZXeBkG1l7tjLomTkS4B6f1V/+uD7iu9j4fuCGyzNJ1m70nIpg76p/hbmi8Vf",
	"AQAA///HLcjVfhAAAA==",
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
