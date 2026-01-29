package tests

import (
	"net/http"
	"testing"
)

// ============================================
// Customer CRUD Tests
// ============================================

func TestCustomerList(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/customers", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	// Should have at least the seeded customer
	data := response["data"].([]interface{})
	if len(data) == 0 {
		t.Error("Expected at least one customer")
	}
}

func TestCustomerList_WithSearch(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	// Test search filter
	w := env.MakeRequest(t, http.MethodGet, "/api/v1/customers?search=Test", nil, cookies)
	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	data := response["data"].([]interface{})
	if len(data) == 0 {
		t.Error("Expected to find customer with 'Test' in name")
	}
}

func TestCustomerGet(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/customers/"+TestCustomerID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	data := response["data"].(map[string]interface{})
	if data["name"] != "Test Customer" {
		t.Errorf("Expected customer name 'Test Customer', got '%v'", data["name"])
	}
}

func TestCustomerGet_NotFound(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/customers/non-existent-id", nil, cookies)

	AssertStatus(t, w, http.StatusNotFound)
}

func TestCustomerCreate_AsCashier(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// Cashiers should be able to create customers
	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"name":    "New Customer",
		"email":   "newcustomer@test.local",
		"phone":   "081999888777",
		"address": "New Customer Address",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/customers", body, cookies)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestCustomerUpdate_AsCashier(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"name":    "Updated Customer",
		"email":   "customer@test.local",
		"phone":   "081234567890",
		"address": "Updated Address",
	}

	w := env.MakeRequest(t, http.MethodPut, "/api/v1/customers/"+TestCustomerID, body, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestCustomerDelete_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a customer to delete
	cookies := env.LoginAsAdmin(t)

	createBody := map[string]interface{}{
		"name":    "Customer To Delete",
		"email":   "delete@test.local",
		"phone":   "081111222333",
		"address": "Delete Address",
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/customers", createBody, cookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	customerID := createData["data"].(map[string]interface{})["id"].(string)

	// Now delete it
	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/customers/"+customerID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestCustomerDelete_AsCashier_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a customer as admin
	adminCookies := env.LoginAsAdmin(t)

	createBody := map[string]interface{}{
		"name":    "Cashier Delete Test",
		"email":   "cashier-delete@test.local",
		"phone":   "081444555666",
		"address": "Test Address",
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/customers", createBody, adminCookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	customerID := createData["data"].(map[string]interface{})["id"].(string)

	// Try to delete as cashier (should fail)
	cashierCookies := env.LoginAsCashier(t)
	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/customers/"+customerID, nil, cashierCookies)

	// Cashiers should NOT be able to delete customers
	AssertStatus(t, w, http.StatusForbidden)
}

func TestCustomerCreate_ValidationError(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	// Missing required field 'name'
	body := map[string]interface{}{
		"email":   "noname@test.local",
		"phone":   "081777888999",
		"address": "No Name Address",
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/customers", body, cookies)

	AssertStatus(t, w, http.StatusBadRequest)
}
