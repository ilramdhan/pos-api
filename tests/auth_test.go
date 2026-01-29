package tests

import (
	"net/http"
	"testing"
)

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	w := env.MakeRequest(t, http.MethodGet, "/health", nil, nil)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%v'", response["status"])
	}
}

// ============================================
// Authentication Tests
// ============================================

func TestAuthLogin_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	body := map[string]string{
		"email":    "admin@test.local",
		"password": "Admin123!",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/login", body, nil)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true, got %v", response["success"])
	}

	// Check cookies are set
	cookies := w.Result().Cookies()
	hasAccessToken := false
	hasRefreshToken := false
	for _, c := range cookies {
		if c.Name == "access_token" {
			hasAccessToken = true
		}
		if c.Name == "refresh_token" {
			hasRefreshToken = true
		}
	}
	if !hasAccessToken {
		t.Error("Expected access_token cookie to be set")
	}
	if !hasRefreshToken {
		t.Error("Expected refresh_token cookie to be set")
	}
}

func TestAuthLogin_InvalidCredentials(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	body := map[string]string{
		"email":    "admin@test.local",
		"password": "wrong-password",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/login", body, nil)

	AssertStatus(t, w, http.StatusUnauthorized)
}

func TestAuthLogin_InvalidEmail(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	body := map[string]string{
		"email":    "notexist@test.local",
		"password": "Admin123!",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/login", body, nil)

	AssertStatus(t, w, http.StatusUnauthorized)
}

func TestAuthLogin_MissingFields(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// Missing password
	body := map[string]string{
		"email": "admin@test.local",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/login", body, nil)

	AssertStatus(t, w, http.StatusBadRequest)
}

func TestAuthMe_Authenticated(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/auth/me", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	data := response["data"].(map[string]interface{})
	if data["email"] != "admin@test.local" {
		t.Errorf("Expected email 'admin@test.local', got '%v'", data["email"])
	}
	if data["role"] != "admin" {
		t.Errorf("Expected role 'admin', got '%v'", data["role"])
	}
}

func TestAuthMe_Unauthenticated(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/auth/me", nil, nil)

	AssertStatus(t, w, http.StatusUnauthorized)
}

func TestAuthRegister_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	body := map[string]string{
		"email":    "newuser@test.local",
		"password": "NewUser123!",
		"name":     "New User",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/register", body, nil)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestAuthRegister_DuplicateEmail(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	body := map[string]string{
		"email":    "admin@test.local", // Already exists
		"password": "NewUser123!",
		"name":     "Duplicate User",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/register", body, nil)

	// Should fail due to duplicate email
	if w.Code == http.StatusCreated {
		t.Error("Expected registration to fail with duplicate email")
	}
}

func TestAuthLogout(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/logout", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	// Check that cookies are cleared
	respCookies := w.Result().Cookies()
	for _, c := range respCookies {
		if c.Name == "access_token" && c.MaxAge > 0 {
			t.Error("Expected access_token cookie to be cleared")
		}
		if c.Name == "refresh_token" && c.MaxAge > 0 {
			t.Error("Expected refresh_token cookie to be cleared")
		}
	}
}

func TestAuthRefresh(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/auth/refresh", nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}
