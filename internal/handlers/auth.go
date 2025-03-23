package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/kiricle/file-uploader/internal/models"
)

type AuthService interface {
	SignUp(dto models.SignUpDTO) (int64, error)
	SignIn(dto models.SignInDTO) (string, error)
}

type AuthHandler struct {
	validate    *validator.Validate
	authService AuthService
}

func NewAuthHandler(validate *validator.Validate, authService AuthService) *AuthHandler {
	return &AuthHandler{validate: validate, authService: authService}
}

func (ah *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var dto models.SignUpDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ah.validate.Struct(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, serviceErr := ah.authService.SignUp(dto)
	if serviceErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error signing up: %v", serviceErr)))
		return
	}

	w.Write([]byte("User registered successfully"))
}

func (ah *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	var dto models.SignInDTO
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := ah.validate.Struct(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := ah.authService.SignIn(dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error signing in: %v", err)))
		return
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error marshalling result: %v", err)))
		return
	}

	w.Write(resultJSON)
}
