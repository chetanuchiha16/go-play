// Package integration contains HTTP-level integration tests that validate
// the API contract. These tests are framework-agnostic — they make real HTTP
// requests and verify the response shape, status codes, and headers.
//
// Prerequisites:
//   - The server must be running on localhost:9999
//   - A PostgreSQL database must be accessible with the credentials in .env
//
// Run with:
//
//	go test -tags=integration ./tests/integration/ -v
package integration

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:9999"

// --- Response structs (framework-agnostic) ---

// APIResponse represents a successful GenericResponse[T]
type APIResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Errors  json.RawMessage `json:"errors"`
}

// ErrorResponse represents Fuego's RFC 7807 error format
type ErrorResponse struct {
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}

type UserData struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt any    `json:"created_at"`
}

type LoginData struct {
	Token string   `json:"token"`
	User  UserData `json:"user"`
}

// --- Helpers ---

func randomEmail() string {
	return fmt.Sprintf("testuser_%d@integration.test", rand.Int())
}

// readBody reads and returns the raw response body bytes.
func readBody(t *testing.T, resp *http.Response) []byte {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
	return body
}

// parseSuccess parses a success response (GenericResponse shape).
func parseSuccess(t *testing.T, body []byte) APIResponse {
	t.Helper()
	var apiResp APIResponse
	err := json.Unmarshal(body, &apiResp)
	require.NoError(t, err, "Failed to parse success response: %s", string(body))
	return apiResp
}

func createUser(t *testing.T, name, email, password string) (*http.Response, []byte) {
	t.Helper()
	form := url.Values{}
	form.Set("name", name)
	form.Set("email", email)
	form.Set("password", password)

	resp, err := http.Post(
		baseURL+"/users/",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	require.NoError(t, err)

	return resp, readBody(t, resp)
}

func login(t *testing.T, email, password string) (*http.Response, []byte) {
	t.Helper()
	form := url.Values{}
	form.Set("email", email)
	form.Set("password", password)

	resp, err := http.Post(
		baseURL+"/auth/login",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	require.NoError(t, err)

	return resp, readBody(t, resp)
}

// --- Tests ---

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL + "/health")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCreateUser_Success(t *testing.T) {
	email := randomEmail()
	resp, body := createUser(t, "Integration Test User", email, "securepassword123")

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := parseSuccess(t, body)
	assert.Equal(t, "success", apiResp.Status)
	assert.Contains(t, apiResp.Message, "created successfully")

	var userData UserData
	err := json.Unmarshal(apiResp.Data, &userData)
	require.NoError(t, err)

	assert.Equal(t, "Integration Test User", userData.Name)
	assert.Equal(t, email, userData.Email)
	assert.NotZero(t, userData.ID)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	email := randomEmail()

	// Create the first user
	resp1, _ := createUser(t, "First User", email, "password123!")
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	// Try to create a second user with the same email — expect RFC 7807 error
	resp2, body := createUser(t, "Second User", email, "password456!")
	assert.Equal(t, http.StatusConflict, resp2.StatusCode)

	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &errResp))
	assert.Equal(t, "Conflict", errResp.Title)
}

func TestCreateUser_InvalidInput(t *testing.T) {
	// Name too short (min=3)
	resp, body := createUser(t, "Ab", randomEmail(), "password123!")
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &errResp))
	assert.Equal(t, 400, errResp.Status)
}

func TestGetUser_Success(t *testing.T) {
	// First create a user
	email := randomEmail()
	_, createBody := createUser(t, "Get Test User", email, "password123!")
	createResp := parseSuccess(t, createBody)
	var createdUser UserData
	err := json.Unmarshal(createResp.Data, &createdUser)
	require.NoError(t, err)

	// Now fetch that user
	resp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, createdUser.ID))
	require.NoError(t, err)

	body := readBody(t, resp)
	apiResp := parseSuccess(t, body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", apiResp.Status)

	var userData UserData
	err = json.Unmarshal(apiResp.Data, &userData)
	require.NoError(t, err)

	assert.Equal(t, createdUser.ID, userData.ID)
	assert.Equal(t, "Get Test User", userData.Name)
	assert.Equal(t, email, userData.Email)
}

func TestGetUser_NotFound(t *testing.T) {
	resp, err := http.Get(baseURL + "/users/999999")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetUser_InvalidID(t *testing.T) {
	resp, err := http.Get(baseURL + "/users/abc")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestListUsers(t *testing.T) {
	// Create a couple of users to ensure there's data
	createUser(t, "List Test User 1", randomEmail(), "password123!")
	createUser(t, "List Test User 2", randomEmail(), "password123!")

	resp, err := http.Get(baseURL + "/users/?limit=5")
	require.NoError(t, err)

	body := readBody(t, resp)
	apiResp := parseSuccess(t, body)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "success", apiResp.Status)

	// Data should be an array
	var users []UserData
	err = json.Unmarshal(apiResp.Data, &users)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(users), 2)
}

func TestLogin_Success(t *testing.T) {
	email := randomEmail()
	password := "password123!"

	// Create user first
	createResp, _ := createUser(t, "Login Test User", email, password)
	assert.Equal(t, http.StatusOK, createResp.StatusCode)

	// Login
	resp, body := login(t, email, password)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	apiResp := parseSuccess(t, body)
	var loginData LoginData
	err := json.Unmarshal(apiResp.Data, &loginData)
	require.NoError(t, err)

	assert.NotEmpty(t, loginData.Token, "JWT token should not be empty")
	assert.Equal(t, email, loginData.User.Email)
	assert.Equal(t, "Login Test User", loginData.User.Name)
}

func TestLogin_WrongPassword(t *testing.T) {
	email := randomEmail()

	// Create user
	createResp, _ := createUser(t, "Wrong Password User", email, "correctpassword")
	assert.Equal(t, http.StatusOK, createResp.StatusCode)

	// Login with wrong password — expect RFC 7807 error
	resp, body := login(t, email, "wrongpassword")
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &errResp))
	assert.Equal(t, "Authentication Failed", errResp.Title)
}

func TestLogin_NonExistentUser(t *testing.T) {
	resp, body := login(t, "nobody@doesnotexist.com", "password")
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var errResp ErrorResponse
	require.NoError(t, json.Unmarshal(body, &errResp))
	assert.Equal(t, 404, errResp.Status)
}

func TestDeleteUser_Unauthorized(t *testing.T) {
	// Create a user
	email := randomEmail()
	_, createBody := createUser(t, "Delete Test User", email, "password123!")
	createResp := parseSuccess(t, createBody)
	var userData UserData
	err := json.Unmarshal(createResp.Data, &userData)
	require.NoError(t, err)

	// Try to delete without auth token
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d", baseURL, userData.ID), nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestDeleteUser_Success(t *testing.T) {
	email := randomEmail()
	password := "password123!"

	// Create a user
	_, createBody := createUser(t, "Delete Me User", email, password)
	createResp := parseSuccess(t, createBody)
	var userData UserData
	err := json.Unmarshal(createResp.Data, &userData)
	require.NoError(t, err)

	// Login to get a token
	_, loginBody := login(t, email, password)
	loginResp := parseSuccess(t, loginBody)
	var loginData LoginData
	err = json.Unmarshal(loginResp.Data, &loginData)
	require.NoError(t, err)

	// Delete with auth token
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d", baseURL, userData.ID), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+loginData.Token)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify user is actually deleted
	getResp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, userData.ID))
	require.NoError(t, err)
	defer getResp.Body.Close()

	assert.Equal(t, http.StatusNotFound, getResp.StatusCode)
}

// TestFullUserLifecycle tests the complete CRUD + auth flow in sequence
func TestFullUserLifecycle(t *testing.T) {
	email := randomEmail()
	password := "lifecycle_pass_123"

	// 1. Create
	createHTTP, createBody := createUser(t, "Lifecycle User", email, password)
	assert.Equal(t, http.StatusOK, createHTTP.StatusCode)

	createResp := parseSuccess(t, createBody)
	var createdUser UserData
	require.NoError(t, json.Unmarshal(createResp.Data, &createdUser))
	assert.NotZero(t, createdUser.ID)

	// 2. Read
	getResp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, createdUser.ID))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResp.StatusCode)
	getResp.Body.Close()

	// 3. Login
	_, loginBody := login(t, email, password)
	loginResp := parseSuccess(t, loginBody)
	var loginData LoginData
	require.NoError(t, json.Unmarshal(loginResp.Data, &loginData))
	assert.NotEmpty(t, loginData.Token)

	// 4. Delete (authenticated)
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/users/%d", baseURL, createdUser.ID), nil)
	req.Header.Set("Authorization", "Bearer "+loginData.Token)
	delResp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, delResp.StatusCode)
	delResp.Body.Close()

	// 5. Verify deleted
	verifyResp, err := http.Get(fmt.Sprintf("%s/users/%d", baseURL, createdUser.ID))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, verifyResp.StatusCode)
	verifyResp.Body.Close()
}
