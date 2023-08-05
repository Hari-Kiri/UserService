package handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

	"github.com/Hari-Kiri/UserService/repository"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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
		errorResponse.Message = fmt.Sprintf("password parameter value (%s) more than 64 characters", parameters.Password)
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

	// Hash password
	hashedPassword, errorHashingPassword := bcrypt.GenerateFromPassword([]byte(parameters.Password), bcrypt.DefaultCost)
	if errorHashingPassword != nil {
		var errorResponse registrationResponse
		errorResponse.Id = 0
		errorResponse.Message = "failed hashing password"
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
		Password:    string(hashedPassword),
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
	})
	if errorResult != nil {
		fmt.Printf("%s", errorResult)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}

	// Compare password
	byteHashedPassword := []byte(result.Password)
	errorComparePassword := bcrypt.CompareHashAndPassword(byteHashedPassword, []byte(parameters.Password))
	if errorComparePassword != nil {
		fmt.Printf("%s", errorComparePassword)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}

	// Update successful login
	updateSuccessfullLogin, errorUpdateSuccessfullLogin := repo.UpdateSuccessfullLogin(ctx.Request().Context(), repository.UpdateSuccessfullLoginInput{
		GetUserDataOutput: result,
	})
	if errorUpdateSuccessfullLogin != nil {
		fmt.Printf("%s", errorUpdateSuccessfullLogin)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}
	if updateSuccessfullLogin.Id != result.Id {
		fmt.Printf("failed update successfull login user id: %d", result.Id)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}

	// Create jwt token
	parsedRsaPrivateKey, errorParseRsaPrivateKey := jwt.ParseRSAPrivateKeyFromPEM([]byte(rsaPrivateKey))
	if errorParseRsaPrivateKey != nil {
		fmt.Printf("%s", errorParseRsaPrivateKey)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"id":       result.Id,
		"password": result.Password,
	})
	jwtTokenString, errorSigning := jwtToken.SignedString(parsedRsaPrivateKey)
	if errorSigning != nil {
		fmt.Printf("%s", errorSigning)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}

	var response loginResponse
	response.Id = result.Id
	response.Token = jwtTokenString
	return ctx.JSON(http.StatusOK, response)
}

func (s *Server) Profile(ctx echo.Context, parameters profileParameters) error {
	// Parse jwt token
	parsedRsaPublicKey, errorParseRsaPublicKey := jwt.ParseRSAPublicKeyFromPEM([]byte(rsaPublicKey))
	if errorParseRsaPublicKey != nil {
		fmt.Printf("%s", errorParseRsaPublicKey)
		return ctx.JSON(http.StatusForbidden, profileResponse{})
	}
	jwtTokenString := strings.Replace(parameters.Authorization, "Bearer ", "", -1)
	parsedJwtToken, errorParsingToken := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return parsedRsaPublicKey, nil
	})
	if errorParsingToken != nil {
		fmt.Printf("%s", errorParsingToken)
		return ctx.JSON(http.StatusForbidden, profileResponse{})
	}

	// Get jwt claims
	var (
		userId       int64
		userPassword string
	)
	if claims, ok := parsedJwtToken.Claims.(jwt.MapClaims); ok && parsedJwtToken.Valid {
		userId = int64(claims["id"].(float64))
		userPassword = claims["password"].(string)
	}

	// Get profile from database
	var repo = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: os.Getenv("DATABASE_URL"),
	})
	result, errorResult := repo.GetUserProfile(ctx.Request().Context(), repository.GetUserProfileInput{
		Id:       userId,
		Password: userPassword,
	})
	if errorResult != nil {
		fmt.Printf("%s", errorResult)
		return ctx.JSON(http.StatusBadRequest, loginResponse{})
	}

	var response profileResponse
	response.Name = result.Name
	response.PhoneNumber = result.PhoneNumber
	return ctx.JSON(http.StatusOK, response)
}

func (s *Server) UpdateProfile(ctx echo.Context, parameters updateProfileParameters) error {
	// Parse jwt token
	parsedRsaPublicKey, errorParseRsaPublicKey := jwt.ParseRSAPublicKeyFromPEM([]byte(rsaPublicKey))
	if errorParseRsaPublicKey != nil {
		fmt.Printf("%s", errorParseRsaPublicKey)
		return ctx.JSON(http.StatusForbidden, updateProfileResponse{})
	}
	jwtTokenString := strings.Replace(parameters.Authorization, "Bearer ", "", -1)
	parsedJwtToken, errorParsingToken := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return parsedRsaPublicKey, nil
	})
	if errorParsingToken != nil {
		fmt.Printf("%s", errorParsingToken)
		return ctx.JSON(http.StatusForbidden, updateProfileResponse{})
	}

	// Get jwt claims
	var (
		userId       int64
		userPassword string
	)
	if claims, ok := parsedJwtToken.Claims.(jwt.MapClaims); ok && parsedJwtToken.Valid {
		userId = int64(claims["id"].(float64))
		userPassword = claims["password"].(string)
	}

	// Get profile from database
	var repo = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: os.Getenv("DATABASE_URL"),
	})
	profile, errorGetProfile := repo.GetUserProfile(ctx.Request().Context(), repository.GetUserProfileInput{
		Id:       userId,
		Password: userPassword,
	})
	if errorGetProfile != nil {
		fmt.Printf("%s", errorGetProfile)
		return ctx.JSON(http.StatusBadRequest, updateProfileResponse{})
	}

	// Conflict of the phone number
	if parameters.PhoneNumber != nil && *parameters.PhoneNumber == profile.PhoneNumber {
		fmt.Printf("CONFLICT: phone number can't updated with same string. parameter=%s | database=%s.", *parameters.PhoneNumber, profile.PhoneNumber)
		return ctx.JSON(http.StatusConflict, updateProfileResponse{})
	}

	var (
		response                                updateProfileResponse
		updatePhoneNumber                       repository.UpdatePhoneNumberOutput
		updateName                              repository.UpdateNameOutput
		errorUpdatePhoneNumber, errorUpdateName error
	)
	// Update phone number
	if parameters.PhoneNumber != nil && *parameters.PhoneNumber != profile.PhoneNumber {
		updatePhoneNumber, errorUpdatePhoneNumber = repo.UpdatePhoneNumber(ctx.Request().Context(), repository.UpdatePhoneNumberInput{
			Id:          userId,
			Password:    userPassword,
			PhoneNumber: *parameters.PhoneNumber,
		})
	}
	if errorUpdatePhoneNumber != nil {
		fmt.Printf("error update phone number: %s", errorUpdatePhoneNumber)
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if updatePhoneNumber.Id != 0 && updatePhoneNumber.Id != userId {
		fmt.Printf("failed update phone number to \"%s\"", *parameters.PhoneNumber)
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if updatePhoneNumber.Id != 0 && updatePhoneNumber.Id == userId {
		response.Message = fmt.Sprintf("Success update phone number to %s.", *parameters.PhoneNumber)
	}

	// Update name
	if parameters.Name != nil {
		updateName, errorUpdateName = repo.UpdateName(ctx.Request().Context(), repository.UpdateNameInput{
			Id:       userId,
			Password: userPassword,
			Name:     *parameters.Name,
		})
	}
	if errorUpdateName != nil {
		fmt.Printf("error update user name: %s", errorUpdateName)
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if updateName.Id != 0 && updateName.Id != userId {
		fmt.Printf("failed update user name to \"%s\"", *parameters.Name)
		return ctx.JSON(http.StatusInternalServerError, response)
	}
	if updateName.Id != 0 && updateName.Id == userId && response.Message != "" {
		response.Message = fmt.Sprintf("%s Success update user name to %s.", response.Message, *parameters.Name)
	}
	if updateName.Id != 0 && updateName.Id == userId && response.Message == "" {
		response.Message = fmt.Sprintf("Success update user name to %s.", *parameters.Name)
	}

	response.Id = userId
	response.Response = true
	return ctx.JSON(http.StatusOK, response)
}
