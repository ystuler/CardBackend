package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"back/config"
	"back/db"
	"back/internal/handler"
	"back/internal/models"
	"back/internal/repository"
	"back/internal/service"
	"back/internal/util"
)

func setupTestEnvironment(t *testing.T) (*httptest.Server, *repository.Repository) {
	// Инициализация конфигурации
	cfg := config.NewConfig()

	// Инициализация базы данных
	database, err := db.NewDatabase(cfg.Database.DSN())
	if err != nil {
		t.Fatalf("could not initialize database connection: %s", err)
	}

	// Инициализация компонентов
	validator := util.NewValidator()
	repos := repository.NewRepository(database.GetDB())
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, validator)
	r := handlers.InitRoutes()

	// Создание тестового сервера
	ts := httptest.NewServer(r)

	return ts, repos
}

func teardownTestEnvironment(ts *httptest.Server) {
	ts.Close()
}

func TestUserSignup(t *testing.T) {
	// Подготовка тестового окружения
	ts, repos := setupTestEnvironment(t)
	defer teardownTestEnvironment(ts)

	hashedPassword, err := util.HashPassword("testpassword")
	if err != nil {
		t.Fatalf("could not hash password: %s", err)
	}

	// Создание пользователя для теста
	testUser := models.User{
		Username:     "testuser",
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	// Подготовка данных для запроса
	signupData := map[string]string{
		"username": "testuser",
		"password": "testpassword",
	}
	jsonData, err := json.Marshal(signupData)
	if err != nil {
		t.Fatalf("could not marshal signup data: %s", err)
	}

	// Выполнение запроса на эндпоинт /auth/signup
	resp, err := http.Post(ts.URL+"/auth/signup", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("could not send signup request: %s", err)
	}
	defer resp.Body.Close()

	// Проверка статуса ответа
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Проверка созданного пользователя в базе данных
	createdUser, err := repos.GetUserByUsername("testuser")
	if err != nil {
		t.Fatalf("could not get user by username: %s", err)
	}

	assert.Equal(t, testUser.Username, createdUser.Username)
	assert.NoError(t, util.CheckPassword("testpassword", createdUser.PasswordHash))
}
