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

type ServerInterface interface {
	Hello(ctx echo.Context, params HelloParams) error
	Registration(ctx echo.Context, params registrationParameters) error
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

		errorBindRequestBody := ctx.Bind(&parameters)
		if errorBindRequestBody != nil {
			errorResult = echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter: %s", errorBindRequestBody))
		}

		errorResult = w.Handler.Registration(ctx, parameters)

	}

	return errorResult
}
