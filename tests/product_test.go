package tests

import (
	"net/http"
	"testing"
)

// ============================================
// Product CRUD Tests
// ============================================

func TestProductList(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/products", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	// Should have at least the seeded product
	data := response["data"].([]interface{})
	if len(data) == 0 {
		t.Error("Expected at least one product")
	}
}

func TestProductList_WithFilters(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	// Test category filter
	w := env.MakeRequest(t, http.MethodGet, "/api/v1/products?category_id="+TestCategoryID, nil, cookies)
	AssertStatus(t, w, http.StatusOK)

	// Test search filter
	w = env.MakeRequest(t, http.MethodGet, "/api/v1/products?search=Test", nil, cookies)
	AssertStatus(t, w, http.StatusOK)

	// Test in_stock filter
	w = env.MakeRequest(t, http.MethodGet, "/api/v1/products?in_stock=true", nil, cookies)
	AssertStatus(t, w, http.StatusOK)
}

func TestProductGet(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/products/"+TestProductID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	data := response["data"].(map[string]interface{})
	if data["name"] != "Test Product" {
		t.Errorf("Expected product name 'Test Product', got '%v'", data["name"])
	}

	// Product should include category info
	category := data["category"]
	if category == nil {
		t.Error("Expected product to include category information")
	}
}

func TestProductGet_NotFound(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/products/non-existent-id", nil, cookies)

	AssertStatus(t, w, http.StatusNotFound)
}

func TestProductCreate_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"category_id": TestCategoryID,
		"sku":         "NEW-001",
		"name":        "New Product",
		"description": "New product description",
		"price":       25000,
		"stock":       50,
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/products", body, cookies)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestProductCreate_AsManager(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsManager(t)

	body := map[string]interface{}{
		"category_id": TestCategoryID,
		"sku":         "MGR-001",
		"name":        "Manager Product",
		"description": "Created by manager",
		"price":       15000,
		"stock":       30,
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/products", body, cookies)

	// Managers should be able to create products
	AssertStatus(t, w, http.StatusCreated)
}

func TestProductCreate_AsCashier_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"category_id": TestCategoryID,
		"sku":         "CSH-001",
		"name":        "Cashier Product",
		"description": "Should fail",
		"price":       10000,
		"stock":       10,
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/products", body, cookies)

	// Cashiers should NOT be able to create products
	AssertStatus(t, w, http.StatusForbidden)
}

func TestProductUpdate_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"category_id": TestCategoryID,
		"sku":         "TEST-001",
		"name":        "Updated Product",
		"description": "Updated description",
		"price":       15000,
		"stock":       80,
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPut, "/api/v1/products/"+TestProductID, body, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestProductUpdateStock_AsCashier(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"operation": "set",
		"quantity":  90,
	}

	w := env.MakeRequest(t, http.MethodPatch, "/api/v1/products/"+TestProductID+"/stock", body, cookies)

	// All authenticated users should be able to update stock
	AssertStatus(t, w, http.StatusOK)
}

func TestProductDelete_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a product to delete
	cookies := env.LoginAsAdmin(t)

	createBody := map[string]interface{}{
		"category_id": TestCategoryID,
		"sku":         "DEL-001",
		"name":        "Product To Delete",
		"description": "Will be deleted",
		"price":       5000,
		"stock":       10,
		"is_active":   true,
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/products", createBody, cookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	productID := createData["data"].(map[string]interface{})["id"].(string)

	// Now delete it
	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/products/"+productID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestProductDelete_AsManager_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a product as admin
	adminCookies := env.LoginAsAdmin(t)

	createBody := map[string]interface{}{
		"category_id": TestCategoryID,
		"sku":         "MGR-DEL-001",
		"name":        "Manager Delete Test",
		"description": "Test delete by manager",
		"price":       5000,
		"stock":       10,
		"is_active":   true,
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/products", createBody, adminCookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	productID := createData["data"].(map[string]interface{})["id"].(string)

	// Try to delete as manager (should fail)
	managerCookies := env.LoginAsManager(t)
	w := env.MakeRequest(t, http.MethodDelete, "/api/v1/products/"+productID, nil, managerCookies)

	// Only admin can delete
	AssertStatus(t, w, http.StatusForbidden)
}

func TestProductCreate_InvalidCategory(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"category_id": "non-existent-category",
		"sku":         "INV-001",
		"name":        "Invalid Category Product",
		"description": "Should fail",
		"price":       10000,
		"stock":       10,
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/products", body, cookies)

	// Should fail due to invalid category
	if w.Code == http.StatusCreated {
		t.Error("Expected product creation to fail with invalid category")
	}
}

func TestProductCreate_DuplicateSKU(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	body := map[string]interface{}{
		"category_id": TestCategoryID,
		"sku":         "TEST-001", // Already exists
		"name":        "Duplicate SKU Product",
		"description": "Should fail",
		"price":       10000,
		"stock":       10,
		"is_active":   true,
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/products", body, cookies)

	// Should fail due to duplicate SKU
	if w.Code == http.StatusCreated {
		t.Error("Expected product creation to fail with duplicate SKU")
	}
}
