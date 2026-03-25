package tests

import (
	"back/internal/models"
	"back/internal/util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errorResponse struct {
	Error string `json:"error"`
}

func makeRequest(t *testing.T, method, url string, body []byte, token string) *http.Response {
	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decodeErrorBody(t *testing.T, resp *http.Response) errorResponse {
	t.Helper()
	defer resp.Body.Close()

	var e errorResponse
	err := json.NewDecoder(resp.Body).Decode(&e)
	require.NoError(t, err)
	return e
}

func createTestUserAndToken(t *testing.T, repos interface {
	CreateUser(user *models.User) (*models.User, error)
}) (int, string) {
	t.Helper()

	username := fmt.Sprintf("user_%d", time.Now().UnixNano())
	hashed, err := util.HashPassword("testpassword")
	require.NoError(t, err)

	user, err := repos.CreateUser(&models.User{Username: username, PasswordHash: hashed})
	require.NoError(t, err)

	token, err := util.GenerateJWT(user)
	require.NoError(t, err)

	return user.ID, token
}

func createCollectionForUser(t *testing.T, repos interface {
	CreateCollection(collection *models.Collection) (*models.Collection, error)
}, userID int) int {
	t.Helper()

	name := fmt.Sprintf("col_%d", time.Now().UnixNano())
	desc := "test"
	collection, err := repos.CreateCollection(&models.Collection{
		Name:        name,
		Description: &desc,
		UserID:      userID,
	})
	require.NoError(t, err)

	return collection.ID
}

func TestError400InvalidJSON(t *testing.T) {
	ts, _ := setupTestEnvironment(t)
	defer teardownTestEnvironment(ts)

	resp := makeRequest(t, http.MethodPost, ts.URL+"/auth/signup", []byte("{"), "")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	errResp := decodeErrorBody(t, resp)
	assert.Equal(t, "invalid JSON format", errResp.Error)
}

func TestError422Validation(t *testing.T) {
	ts, _ := setupTestEnvironment(t)
	defer teardownTestEnvironment(ts)

	resp := makeRequest(t, http.MethodPost, ts.URL+"/auth/signup", []byte(`{"password":"x"}`), "")
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	errResp := decodeErrorBody(t, resp)
	assert.Contains(t, errResp.Error, "invalid input")
}

func TestError401UnauthorizedWithoutToken(t *testing.T) {
	ts, _ := setupTestEnvironment(t)
	defer teardownTestEnvironment(ts)

	resp := makeRequest(t, http.MethodGet, ts.URL+"/collections/", nil, "")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	errResp := decodeErrorBody(t, resp)
	assert.Equal(t, "unauthorized", errResp.Error)
}

func TestError403ForbiddenForeignCollection(t *testing.T) {
	ts, repos := setupTestEnvironment(t)
	defer teardownTestEnvironment(ts)

	ownerID, _ := createTestUserAndToken(t, repos)
	_, attackerToken := createTestUserAndToken(t, repos)
	collectionID := createCollectionForUser(t, repos, ownerID)

	resp := makeRequest(t, http.MethodGet, fmt.Sprintf("%s/collections/%d/", ts.URL, collectionID), nil, attackerToken)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	errResp := decodeErrorBody(t, resp)
	assert.Equal(t, "forbidden", errResp.Error)
}

func TestError404NotFound(t *testing.T) {
	ts, repos := setupTestEnvironment(t)
	defer teardownTestEnvironment(ts)

	_, token := createTestUserAndToken(t, repos)

	resp := makeRequest(t, http.MethodGet, ts.URL+"/collections/999999999/", nil, token)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	errResp := decodeErrorBody(t, resp)
	assert.Equal(t, "collection not found", errResp.Error)
}

func TestError405MethodNotAllowed(t *testing.T) {
	ts, _ := setupTestEnvironment(t)
	defer teardownTestEnvironment(ts)

	resp := makeRequest(t, http.MethodGet, ts.URL+"/auth/login", nil, "")
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	errResp := decodeErrorBody(t, resp)
	assert.Equal(t, "method not allowed", errResp.Error)
}
