package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params HelloParams) error {

	var resp HelloResponse
	resp.Message = fmt.Sprintf("Hello User %d", params.Id)
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Registration(ctx echo.Context, parameters registrationParameters) error {
	phoneNumberLength := 0
	for i := 0; i < len(parameters.PhoneNumber[3:]); i++ {
		phoneNumberLength++
	}
	if phoneNumberLength <= 9 {
		var errorResponse registrationErrorResponse
		errorResponse.Message = fmt.Sprintf("phoneNumber parameter value (%s) less than 10 characters", parameters.PhoneNumber)
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if phoneNumberLength >= 14 {
		var errorResponse registrationErrorResponse
		errorResponse.Message = fmt.Sprintf("phoneNumber parameter value (%s) more than 13 characters", parameters.PhoneNumber)
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	nameLength := 0
	for i := 0; i < len(parameters.Name); i++ {
		nameLength++
	}
	if nameLength <= 2 {
		var errorResponse registrationErrorResponse
		errorResponse.Message = fmt.Sprintf("name parameter value (%s) less than 3 characters", parameters.Name)
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if nameLength >= 61 {
		var errorResponse registrationErrorResponse
		errorResponse.Message = fmt.Sprintf("name parameter value (%s) more than 60 characters", parameters.Name)
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	passwordLength := 0
	for i := 0; i < len(parameters.Password); i++ {
		passwordLength++
	}
	if passwordLength <= 5 {
		var errorResponse registrationErrorResponse
		errorResponse.Message = fmt.Sprintf("password parameter value (%s) less than 6 characters", parameters.Password)
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if passwordLength >= 65 {
		var errorResponse registrationErrorResponse
		errorResponse.Message = fmt.Sprintf("password parameter value (%s) more than 65 characters", parameters.Password)
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	var response registrationResponse
	response.Message = fmt.Sprintf("Halo User %s", parameters.Name)
	return ctx.JSON(http.StatusOK, response)
}
