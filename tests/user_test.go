package tests

import (
	"net/http"
	"testing"
)

// ============================================
// User Management Tests (Admin Only)
// ============================================

func TestUserList_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/users", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	// Should have the seeded users
	data := response["data"].([]interface{})
	if len(data) < 3 {
		t.Errorf("Expected at least 3 users, got %d", len(data))
	}
}

func TestUserList_AsManager_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsManager(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/users", nil, cookies)

	// Managers should NOT be able to list users
	AssertStatus(t, w, http.StatusForbidden)
}

func TestUserList_AsCashier_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/users", nil, cookies)

	// Cashiers should NOT be able to list users
	AssertStatus(t, w, http.StatusForbidden)
}

func TestUserGet_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/users/"+TestManagerID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	data := response["data"].(map[string]interface{})
	if data["email"] != "manager@test.local" {
		t.Errorf("Expected email 'manager@test.local', got '%v'", data["email"])
	}
}

func TestUserCreate_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"email":    "newstaff@test.local",
		"password": "NewStaff123!",
		"name":     "New Staff Member",
		"role":     "cashier",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/users", body, cookies)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestUserCreate_AsManager_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsManager(t)

	body := map[string]interface{}{
		"email":    "managerstaff@test.local",
		"password": "ManagerStaff123!",
		"name":     "Manager Created Staff",
		"role":     "cashier",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/users", body, cookies)

	// Managers should NOT be able to create users
	AssertStatus(t, w, http.StatusForbidden)
}

func TestUserUpdate_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"email":     "manager@test.local",
		"name":      "Updated Manager Name",
		"role":      "manager",
		"is_active": true,
	}

	w := env.MakeRequest(t, http.MethodPut, "/api/v1/users/"+TestManagerID, body, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestUserDelete_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a user to delete
	cookies := env.LoginAsAdmin(t)

	createBody := map[string]interface{}{
		"email":    "deleteuser@test.local",
		"password": "DeleteUser123!",
		"name":     "User To Delete",
		"role":     "cashier",
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/users", createBody, cookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	userID := createData["data"].(map[string]interface{})["id"].(string)

	// Now delete it
	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/users/"+userID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestUserDelete_AsManager_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsManager(t)

	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/users/"+TestCashierID, nil, cookies)

	// Managers should NOT be able to delete users
	AssertStatus(t, w, http.StatusForbidden)
}

func TestUserCreate_DuplicateEmail(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"email":    "admin@test.local", // Already exists
		"password": "DuplicateAdmin123!",
		"name":     "Duplicate Admin",
		"role":     "admin",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/users", body, cookies)

	// Should fail due to duplicate email
	if w.Code == http.StatusCreated {
		t.Error("Expected user creation to fail with duplicate email")
	}
}

func TestUserCreate_ValidationError(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	// Invalid role
	body := map[string]interface{}{
		"email":    "invalidrole@test.local",
		"password": "InvalidRole123!",
		"name":     "Invalid Role User",
		"role":     "superadmin", // Invalid role
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/users", body, cookies)

	// Should fail due to invalid role
	AssertStatus(t, w, http.StatusBadRequest)
}

func TestUserList_WithRoleFilter(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	// Filter by cashier role
	w := env.MakeRequest(t, http.MethodGet, "/api/v1/users?role=cashier", nil, cookies)
	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	data := response["data"].([]interface{})

	// All returned users should be cashiers
	for _, user := range data {
		userData := user.(map[string]interface{})
		if userData["role"] != "cashier" {
			t.Errorf("Expected role 'cashier', got '%v'", userData["role"])
		}
	}
}
