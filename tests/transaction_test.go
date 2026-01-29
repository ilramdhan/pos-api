package tests

import (
	"net/http"
	"testing"
)

// ============================================
// Transaction Tests
// ============================================

func TestTransactionList(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/transactions", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestTransactionCreate_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"customer_id":    TestCustomerID,
		"payment_method": "cash",
		"notes":          "Test transaction",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   2,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", body, cookies)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}

	// Check transaction data
	data := response["data"].(map[string]interface{})
	if data["status"] != "completed" {
		t.Errorf("Expected status 'completed', got '%v'", data["status"])
	}
	if data["payment_method"] != "cash" {
		t.Errorf("Expected payment_method 'cash', got '%v'", data["payment_method"])
	}
}

func TestTransactionCreate_WithoutCustomer(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "cash",
		"notes":          "Walk-in customer transaction",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", body, cookies)

	AssertStatus(t, w, http.StatusCreated)
}

func TestTransactionCreate_InvalidPaymentMethod(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "invalid-method",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", body, cookies)

	// Should fail due to invalid payment method
	AssertStatus(t, w, http.StatusBadRequest)
}

func TestTransactionCreate_EmptyItems(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "cash",
		"items":          []map[string]interface{}{},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", body, cookies)

	// Should fail due to empty items
	AssertStatus(t, w, http.StatusBadRequest)
}

func TestTransactionCreate_InvalidProduct(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "cash",
		"items": []map[string]interface{}{
			{
				"product_id": "non-existent-product",
				"quantity":   1,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", body, cookies)

	// Should fail due to invalid product
	if w.Code == http.StatusCreated {
		t.Error("Expected transaction creation to fail with invalid product")
	}
}

func TestTransactionGet(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a transaction
	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "cash",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", body, cookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	transactionID := createData["data"].(map[string]interface{})["id"].(string)

	// Now get it
	w := env.MakeRequest(t, http.MethodGet, "/api/v1/transactions/"+transactionID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	data := response["data"].(map[string]interface{})
	if data["id"] != transactionID {
		t.Errorf("Expected transaction ID '%s', got '%v'", transactionID, data["id"])
	}

	// Transaction should include items
	items := data["items"].([]interface{})
	if len(items) == 0 {
		t.Error("Expected transaction to include items")
	}
}

func TestTransactionGet_NotFound(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/transactions/non-existent-id", nil, cookies)

	AssertStatus(t, w, http.StatusNotFound)
}

func TestTransactionUpdateStatus_AsAdmin(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a transaction as cashier
	cashierCookies := env.LoginAsCashier(t)

	createBody := map[string]interface{}{
		"payment_method": "cash",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", createBody, cashierCookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	transactionID := createData["data"].(map[string]interface{})["id"].(string)

	// Update status as admin
	adminCookies := env.LoginAsAdmin(t)

	updateBody := map[string]interface{}{
		"status": "cancelled",
	}

	w := env.MakeRequest(t, http.MethodPatch, "/api/v1/transactions/"+transactionID+"/status", updateBody, adminCookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestTransactionUpdateStatus_AsCashier_Forbidden(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	// First create a transaction
	cashierCookies := env.LoginAsCashier(t)

	createBody := map[string]interface{}{
		"payment_method": "cash",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", createBody, cashierCookies)
	AssertStatus(t, createResp, http.StatusCreated)

	createData := ParseResponse(t, createResp)
	transactionID := createData["data"].(map[string]interface{})["id"].(string)

	// Try to update status as cashier (should fail)
	updateBody := map[string]interface{}{
		"status": "cancelled",
	}

	w := env.MakeRequest(t, http.MethodPatch, "/api/v1/transactions/"+transactionID+"/status", updateBody, cashierCookies)

	// Cashiers should NOT be able to update transaction status
	AssertStatus(t, w, http.StatusForbidden)
}

func TestTransactionList_WithFilters(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsAdmin(t)

	// Test status filter
	w := env.MakeRequest(t, http.MethodGet, "/api/v1/transactions?status=completed", nil, cookies)
	AssertStatus(t, w, http.StatusOK)

	// Test payment method filter
	w = env.MakeRequest(t, http.MethodGet, "/api/v1/transactions?payment_method=cash", nil, cookies)
	AssertStatus(t, w, http.StatusOK)

	// Test date filter
	w = env.MakeRequest(t, http.MethodGet, "/api/v1/transactions?date_from=2025-01-01&date_to=2025-12-31", nil, cookies)
	AssertStatus(t, w, http.StatusOK)
}

func TestTransactionCreate_VerifyStockDeduction(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	// Get initial stock
	getProductResp := env.MakeRequest(t, http.MethodGet, "/api/v1/products/"+TestProductID, nil, cookies)
	AssertStatus(t, getProductResp, http.StatusOK)
	productData := ParseResponse(t, getProductResp)
	initialStock := int(productData["data"].(map[string]interface{})["stock"].(float64))

	// Create transaction with 5 items
	quantity := 5
	body := map[string]interface{}{
		"payment_method": "cash",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   quantity,
			},
		},
	}

	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/transactions", body, cookies)
	AssertStatus(t, createResp, http.StatusCreated)

	// Check stock was deducted
	getProductResp = env.MakeRequest(t, http.MethodGet, "/api/v1/products/"+TestProductID, nil, cookies)
	AssertStatus(t, getProductResp, http.StatusOK)
	productData = ParseResponse(t, getProductResp)
	newStock := int(productData["data"].(map[string]interface{})["stock"].(float64))

	expectedStock := initialStock - quantity
	if newStock != expectedStock {
		t.Errorf("Expected stock %d after transaction, got %d", expectedStock, newStock)
	}
}
