package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"seismic-monitor/backend/internal/auth"
	"seismic-monitor/backend/internal/database"
	"seismic-monitor/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository es un mock de ports.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUserLocation(userID string, latitude, longitude, alertRadius, minMagnitude float64) error {
	args := m.Called(userID, latitude, longitude, alertRadius, minMagnitude)
	return args.Error(0)
}

func (m *MockUserRepository) GetAffectedUsers(sismo models.Feature) ([]models.User, error) {
	args := m.Called(sismo)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetUsersNearLocation(lon, lat float64) ([]models.User, error) {
	args := m.Called(lon, lat)
	return args.Get(0).([]models.User), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		jwtService := auth.NewJWTService("secret")
		handler := NewAuthHandler(mockRepo, jwtService)

		mockRepo.On("FindUserByEmail", "test@example.com").Return(nil, nil)
		mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

		r := gin.Default()
		r.POST("/register", handler.Register)

		reqBody := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		assert.Contains(t, w.Body.String(), "usuario registrado con éxito")
		mockRepo.AssertExpectations(t)
	})

	t.Run("UserAlreadyExists", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		jwtService := auth.NewJWTService("secret")
		handler := NewAuthHandler(mockRepo, jwtService)

		mockRepo.On("FindUserByEmail", "test@example.com").Return(&models.User{Email: "test@example.com"}, nil)

		r := gin.Default()
		r.POST("/register", handler.Register)

		reqBody := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.Contains(t, w.Body.String(), "el correo ya está registrado")
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		jwtService := auth.NewJWTService("secret")
		handler := NewAuthHandler(mockRepo, jwtService)

		password := "password123"
		hashedPassword, _ := database.HashPassword(password)
		user := &models.User{
			ID:           "123",
			Email:        "test@example.com",
			PasswordHash: hashedPassword,
		}

		mockRepo.On("FindUserByEmail", "test@example.com").Return(user, nil)

		r := gin.Default()
		r.POST("/login", handler.Login)

		reqBody := models.LoginRequest{
			Email:    "test@example.com",
			Password: password,
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotEmpty(t, response["token"])
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		jwtService := auth.NewJWTService("secret")
		handler := NewAuthHandler(mockRepo, jwtService)

		mockRepo.On("FindUserByEmail", "test@example.com").Return(nil, errors.New("not found"))

		r := gin.Default()
		r.POST("/login", handler.Login)

		reqBody := models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.Contains(t, w.Body.String(), "credenciales inválidas")
	})
}
