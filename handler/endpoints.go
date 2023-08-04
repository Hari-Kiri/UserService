package handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"unicode"

	"github.com/Hari-Kiri/UserService/repository"
	jwt "github.com/golang-jwt/jwt/v5"
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
	// Check phone number length
	phoneNumberLength := 0
	for i := 0; i < len(parameters.PhoneNumber[3:]); i++ {
		phoneNumberLength++
	}
	if phoneNumberLength <= 9 {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = fmt.Sprintf("phoneNumber parameter value (%s) less than 10 characters", parameters.PhoneNumber)
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if phoneNumberLength >= 14 {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = fmt.Sprintf("phoneNumber parameter value (%s) more than 13 characters", parameters.PhoneNumber)
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	// Check name length
	nameLength := 0
	for i := 0; i < len(parameters.Name); i++ {
		nameLength++
	}
	if nameLength <= 2 {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = fmt.Sprintf("name parameter value (%s) less than 3 characters", parameters.Name)
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if nameLength >= 61 {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = fmt.Sprintf("name parameter value (%s) more than 60 characters", parameters.Name)
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	// Check password length
	passwordLength := 0
	for i := 0; i < len(parameters.Password); i++ {
		passwordLength++
	}
	if passwordLength <= 5 {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = fmt.Sprintf("password parameter value (%s) less than 6 characters", parameters.Password)
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if passwordLength >= 65 {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = fmt.Sprintf("password parameter value (%s) more than 65 characters", parameters.Password)
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	// Check password have any capital character, number and special (non alpha-numeric character)
	passwordRune := []rune(parameters.Password)
	hasAnyCapitalCharacter := false
	hasAnyNumber := false
	alphanumeric := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	hasAnyNonAlphanumeric := false
	for i := 0; i < len(passwordRune); i++ {
		if unicode.IsUpper(passwordRune[i]) {
			hasAnyCapitalCharacter = true
		}
		if unicode.IsNumber(passwordRune[i]) {
			hasAnyNumber = true
		}
		if !alphanumeric.MatchString(string(passwordRune[i])) {
			hasAnyNonAlphanumeric = true
		}
	}
	if !hasAnyCapitalCharacter {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = "password not contain any capital character"
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if !hasAnyNumber {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = "password not contain any number"
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}
	if !hasAnyNonAlphanumeric {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = "password not contain any non alpha-numeric character"
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	// Insert data to database
	var repo = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: os.Getenv("DATABASE_URL"),
	})
	result, errorResult := repo.InsertUserData(ctx.Request().Context(), repository.InsertUserDataInput{
		Name:        parameters.Name,
		PhoneNumber: parameters.PhoneNumber,
		Password:    parameters.Password,
	})
	if errorResult != nil {
		fmt.Printf("%s", errorResult)
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = "failed create new user, error: can't create new user data"
		errorResponse.Response = false
		return ctx.JSON(http.StatusBadRequest, errorResponse)
	}

	var response registrationResponse
	response.Id = result.Id
	response.Message = "success create new user"
	response.Response = true
	return ctx.JSON(http.StatusOK, response)
}

func (s *Server) Login(ctx echo.Context, parameters loginParameters) error {
	// Get user data from database
	var repo = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: os.Getenv("DATABASE_URL"),
	})
	result, errorResult := repo.GetUserData(ctx.Request().Context(), repository.GetUserDataInput{
		PhoneNumber: parameters.PhoneNumber,
		Password:    parameters.Password,
	})
	if errorResult != nil {
		fmt.Printf("%s", errorResult)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}

	// Create jwt token
	rsaPrivateKey, errorParseRsaPrivateKey := jwt.ParseRSAPrivateKeyFromPEM([]byte(rsaPrivateKey))
	if errorParseRsaPrivateKey != nil {
		fmt.Printf("%s", errorParseRsaPrivateKey)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"phoneNumber": parameters.PhoneNumber,
		"password":    parameters.Password,
	})
	jwtTokenString, errorSigning := jwtToken.SignedString(rsaPrivateKey)
	if errorSigning != nil {
		fmt.Printf("%s", errorSigning)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}

	var response loginResponse
	response.Id = result.Id
	response.Token = jwtTokenString
	return ctx.JSON(http.StatusOK, response)
}
