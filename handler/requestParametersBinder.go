package handler

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// HelloResponse defines model for HelloResponse.
type HelloResponse struct {
	Message string `json:"message"`
}

// HelloParams defines parameters for Hello.
type HelloParams struct {
	Id int `form:"id" json:"id"`
}

type registrationResponse struct {
	Id       int64  `json:"id"`
	Message  string `json:"message"`
	Response bool   `json:"response"`
}
type registrationParameters struct {
	PhoneNumber string `json:"phoneNumber"`
	Name        string `json:"name"`
	Password    string `json:"password"`
}

type loginResponse struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}
type loginParameters struct {
	PhoneNumber string `form:"phoneNumber"`
	Password    string `form:"password"`
}

type profileResponse struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
}
type profileParameters struct {
	Authorization string `header:"Authorization"`
}

type updateProfileParameters struct {
	Authorization string  `header:"Authorization"`
	PhoneNumber   *string `json:"phoneNumber"`
	Name          *string `json:"name"`
}
type updateProfileResponse struct {
	registrationResponse
}

type ServerInterface interface {
	Hello(ctx echo.Context, params HelloParams) error
	Registration(ctx echo.Context, params registrationParameters) error
	Login(ctx echo.Context, params loginParameters) error
	Profile(ctx echo.Context, params profileParameters) error
	UpdateProfile(ctx echo.Context, params updateProfileParameters) error
}
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

func (w *ServerInterfaceWrapper) requestParametersBinder(ctx echo.Context) error {
	var errorResult error

	switch ctx.Path()[1:] {

	case "hello":
		// Parameter object where we will unmarshal all parameters from the context
		var params HelloParams
		// ------------- Required query parameter "id" -------------

		errorBindQueryParameter := runtime.BindQueryParameter("form", true, true, "id", ctx.QueryParams(), &params.Id)
		if errorBindQueryParameter != nil {
			errorResult = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", errorBindQueryParameter))
		}

		// Invoke the callback with all the unmarshaled arguments
		errorResult = w.Handler.Hello(ctx, params)
	case "registration":
		var parameters registrationParameters

		// Bind json (application/json) data from request body
		errorBindRequestBody := ctx.Bind(&parameters)
		if errorBindRequestBody != nil {
			errorResult = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter: %s", errorBindRequestBody))
		}

		errorResult = w.Handler.Registration(ctx, parameters)
	case "login":
		var parameters loginParameters

		// Bind form (application/x-www-form-urlencoded) data from request body
		errorBindRequestBody := ctx.Bind(&parameters)
		if errorBindRequestBody != nil {
			errorResult = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter: %s", errorBindRequestBody))
		}

		errorResult = w.Handler.Login(ctx, parameters)
	case "profile":
		var parameters profileParameters

		// Bind request headers data
		binder := &echo.DefaultBinder{}
		errorBindRequestHeader := binder.BindHeaders(ctx, &parameters)
		if errorBindRequestHeader != nil {
			errorResult = echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Invalid format for parameter: %s", errorBindRequestHeader))
		}

		errorResult = w.Handler.Profile(ctx, parameters)
	case "update-profile":
		var parameters updateProfileParameters

		// Bind request headers data
		binder := &echo.DefaultBinder{}
		errorBindRequestHeader := binder.BindHeaders(ctx, &parameters)
		if errorBindRequestHeader != nil {
			errorResult = echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("Invalid format for parameter: %s", errorBindRequestHeader))
		}

		// Bind json (application/json) data from request body
		errorBindRequestBody := ctx.Bind(&parameters)
		if errorBindRequestBody != nil {
			errorResult = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter: %s", errorBindRequestBody))
		}

		errorResult = w.Handler.UpdateProfile(ctx, parameters)
	}

	return errorResult
}
