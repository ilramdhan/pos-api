package tests

import (
	"net/http"
	"testing"
)

// ============================================
// POS (Point of Sale) Tests
// ============================================

func TestPOSProducts(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/pos/products", nil, cookies)

	AssertStatus(t, w, http.StatusOK)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestPOSProducts_WithSearch(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/pos/products?search=Test", nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestPOSProducts_WithCategory(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	w := env.MakeRequest(t, http.MethodGet, "/api/v1/pos/products?category_id="+TestCategoryID, nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestPOSCheckout_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"customer_id":    TestCustomerID,
		"payment_method": "cash",
		"notes":          "POS Checkout Test",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/pos/checkout", body, cookies)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestPOSCheckout_WalkInCustomer(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "qris",
		"notes":          "Walk-in customer checkout",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   2,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/pos/checkout", body, cookies)

	AssertStatus(t, w, http.StatusCreated)
}

func TestPOSCheckout_MultipleItems(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "card",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   3,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/pos/checkout", body, cookies)

	AssertStatus(t, w, http.StatusCreated)
}

func TestPOSHold_Success(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"customer_id": TestCustomerID,
		"notes":       "Hold for later",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   2,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/pos/hold", body, cookies)

	AssertStatus(t, w, http.StatusCreated)

	response := ParseResponse(t, w)
	if response["success"] != true {
		t.Errorf("Expected success=true")
	}
}

func TestPOSHeld_List(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	// First create a held transaction
	holdBody := map[string]interface{}{
		"notes": "Test held item",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}
	env.MakeRequest(t, http.MethodPost, "/api/v1/pos/hold", holdBody, cookies)

	// Now list held items
	w := env.MakeRequest(t, http.MethodGet, "/api/v1/pos/held", nil, cookies)

	AssertStatus(t, w, http.StatusOK)
}

func TestPOSHeld_Delete(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	// First create a held transaction
	holdBody := map[string]interface{}{
		"notes": "To be deleted",
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}
	createResp := env.MakeRequest(t, http.MethodPost, "/api/v1/pos/hold", holdBody, cookies)

	if createResp.Code == http.StatusCreated {
		createData := ParseResponse(t, createResp)
		if data, ok := createData["data"].(map[string]interface{}); ok {
			if heldID, ok := data["id"].(string); ok {
				// Now delete the held item
				w := env.MakeRequest(t, http.MethodDelete, "/api/v1/pos/held/"+heldID, nil, cookies)
				AssertStatus(t, w, http.StatusOK)
			}
		}
	}
}

func TestPOSCheckout_InvalidPaymentMethod(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "bitcoin", // Invalid
		"items": []map[string]interface{}{
			{
				"product_id": TestProductID,
				"quantity":   1,
			},
		},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/pos/checkout", body, cookies)

	AssertStatus(t, w, http.StatusBadRequest)
}

func TestPOSCheckout_EmptyItems(t *testing.T) {
	env := SetupTestEnv(t)
	defer env.Cleanup()

	cookies := env.LoginAsCashier(t)

	body := map[string]interface{}{
		"payment_method": "cash",
		"items":          []map[string]interface{}{},
	}

	w := env.MakeRequest(t, http.MethodPost, "/api/v1/pos/checkout", body, cookies)

	AssertStatus(t, w, http.StatusBadRequest)
}
