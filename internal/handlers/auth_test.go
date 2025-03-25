package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/kiricle/file-uploader/internal/handlers"
	mock_handlers "github.com/kiricle/file-uploader/internal/mocks"
	"github.com/kiricle/file-uploader/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_handlers.MockAuthService, signUpDto models.SignUpDTO)

	testCases := []struct {
		name               string
		inputBody          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "Success - Valid Input",
			inputBody: `{"email":"test@example.com","password": "securepassword"}`,
			mockBehavior: func(s *mock_handlers.MockAuthService, signUpDto models.SignUpDTO) {
				s.EXPECT().SignUp(signUpDto).Return(int64(1), nil)
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody:       "User registered successfully",
		},
		{
			name:               "Failure - Invalid Email",
			inputBody:          `{"email": "invalid-email","password": "somepass123"}`,
			mockBehavior:       func(s *mock_handlers.MockAuthService, signUpDto models.SignUpDTO) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Key: 'SignUpDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name:               "Failure - Invalid Password",
			inputBody:          `{"email": "some@gmail.com","password": "short"}`,
			mockBehavior:       func(s *mock_handlers.MockAuthService, signUpDto models.SignUpDTO) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Key: 'SignUpDTO.Password' Error:Field validation for 'Password' failed on the 'min' tag",
		},
		{
			name:      "Failure - AuthService Returns Error",
			inputBody: `{"email": "existed-email@gmail.com", "password": "someSecurePassword123"}`,
			mockBehavior: func(s *mock_handlers.MockAuthService, signUpDto models.SignUpDTO) {
				s.EXPECT().SignUp(signUpDto).Return(int64(0), errors.New("user already exists"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Error signing up: user already exists",
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			var signUpDto models.SignUpDTO
			err := json.Unmarshal([]byte(tc.inputBody), &signUpDto)
			if err != nil {
				t.Fatalf("Failed to unmarshal input body: %v", err)
			}

			authService := mock_handlers.NewMockAuthService(c)
			tc.mockBehavior(authService, signUpDto)

			validate := validator.New()
			handler := handlers.NewAuthHandler(validate, authService)

			r := chi.NewRouter()

			r.Post("/api/v1/sign-up", handler.SignUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sign-up", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_handlers.MockAuthService, signInDto models.SignInDTO)

	testCases := []struct {
		name               string
		inputBody          string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "Success - Valid Credentials",
			inputBody: `{"email":"test@example.com","password":"securepassword"}`,
			mockBehavior: func(s *mock_handlers.MockAuthService, signInDto models.SignInDTO) {
				s.EXPECT().SignIn(signInDto).Return("mocked-jwt-token", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"token":"mocked-jwt-token"}`,
		},
		{
			name:               "Failure - Invalid Email Format",
			inputBody:          `{"email": "invalid-email","password": "somepassword"}`,
			mockBehavior:       func(s *mock_handlers.MockAuthService, signInDto models.SignInDTO) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Key: 'SignInDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag",
		},
		{
			name:               "Failure - Missing Password",
			inputBody:          `{"email": "test@example.com","password": ""}`,
			mockBehavior:       func(s *mock_handlers.MockAuthService, signInDto models.SignInDTO) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       "Key: 'SignInDTO.Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
		{
			name:      "Failure - Wrong Password",
			inputBody: `{"email": "test@example.com", "password": "wrongpassword"}`,
			mockBehavior: func(s *mock_handlers.MockAuthService, signInDto models.SignInDTO) {
				s.EXPECT().SignIn(signInDto).Return("", errors.New("invalid credentials"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBody:       "Error signing in: invalid credentials",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			var signInDto models.SignInDTO
			err := json.Unmarshal([]byte(tc.inputBody), &signInDto)
			if err != nil {
				t.Fatalf("Failed to unmarshal input body: %v", err)
			}

			authService := mock_handlers.NewMockAuthService(c)
			tc.mockBehavior(authService, signInDto)

			validate := validator.New()

			handler := handlers.NewAuthHandler(validate, authService)

			r := chi.NewRouter()
			r.Post("/api/v1/sign-in", handler.SignIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/sign-in", bytes.NewBufferString(tc.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}
}
