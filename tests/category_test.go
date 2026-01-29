package tests

import (
	"net/http"
	"testing"
)

// ============================================
// Category CRUD Tests
// ============================================

func TestCategoryList(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/categories", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	// Should have at least the seeded category
	data := response["data"].([]interface{})
	if len(data) == 0 {
		t.Error("Expected at least one category")
	}
}

func TestCategoryGet(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/categories/"+TestCategoryID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	data := response["data"].(map[string]interface{})
	if data["name"] != "Test Category" {
		t.Errorf("Expected category name 'Test Category', got '%v'", data["name"])
	}
}

func TestCategoryGet_NotFound(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/categories/non-existent-id", nil, cookies)

	AssertStatus(t, w, http.StatusNotFound)
}

func TestCategoryCreate_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"name":        "New Category",
		"description": "New category description",
		"slug":        "new-category",
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/categories", body, cookies)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestCategoryCreate_AsManager(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsManager(t)

	body := map[string]interface{}{
		"name":        "Manager Category",
		"description": "Created by manager",
		"slug":        "manager-category",
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/categories", body, cookies)

	// Managers should be able to create categories
	AssertStatus(t, w, http.StatusCreated)
}

func TestCategoryCreate_AsCashier_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"name":        "Cashier Category",
		"description": "Should fail",
		"slug":        "cashier-category",
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/categories", body, cookies)

	// Cashiers should NOT be able to create categories
	AssertStatus(t, w, http.StatusForbidden)
}

func TestCategoryUpdate_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"name":        "Updated Category",
		"description": "Updated description",
		"slug":        "test-category",
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPut, "/api/v1/categories/"+TestCategoryID, body, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestCategoryDelete_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a category to delete
	cookies := env.LoginAsAdmin(t)

	createBody := map[string]interface{}{
		"name":        "Category To Delete",
		"description": "Will be deleted",
		"slug":        "delete-me",
		"is_active":   true,
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/categories", createBody, cookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	categoryID := createData["data"].(map[string]interface{})["id"].(string)

	// Now delete it
	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/categories/"+categoryID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestCategoryDelete_AsManager_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a category as admin
	adminCookies := env.LoginAsAdmin(t)

	createBody := map[string]interface{}{
		"name":        "Manager Delete Test",
		"description": "Test delete by manager",
		"slug":        "manager-delete-test",
		"is_active":   true,
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/categories", createBody, adminCookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	categoryID := createData["data"].(map[string]interface{})["id"].(string)

	// Try to delete as manager (should fail)
	managerCookies := env.LoginAsManager(t)
	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/categories/"+categoryID, nil, managerCookies)

	// Only admin can delete
	AssertStatus(t, w, http.StatusForbidden)
}

func TestCategoryCreate_ValidationError(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	// Missing required field 'name'
	body := map[string]interface{}{
		"description": "Missing name",
		"slug":        "missing-name",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/categories", body, cookies)

	AssertStatus(t, w, http.StatusBadRequest)
}
