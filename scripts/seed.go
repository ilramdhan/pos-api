package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Ensure data directory exists
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Failed to create data directory: %v", err)
	}

	// Open database
	dbPath := filepath.Join(dataDir, "pos.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Run migrations
	log.Println("Running migrations...")
	migrationPath := "internal/database/migrations/001_init.sql"
	migrationContent, err := os.ReadFile(migrationPath)
	if err != nil {
		log.Fatalf("Failed to read migration file: %v", err)
	}
	if _, err := db.Exec(string(migrationContent)); err != nil {
		log.Printf("Migration warning (may already exist): %v", err)
	}
	log.Println("Migrations complete")

	ctx := context.Background()
	now := time.Now()

	log.Println("=== Seeding Database ===")

	// ============================================
	// SEED USERS
	// ============================================
	log.Println("\n[1/6] Seeding users...")
	users := []struct {
		id       string
		email    string
		password string
		name     string
		role     string
	}{
		{uuid.New().String(), "admin@pos.local", "Admin123!", "Admin User", "admin"},
		{uuid.New().String(), "manager@pos.local", "Manager123!", "Manager User", "manager"},
		{uuid.New().String(), "cashier@pos.local", "Cashier123!", "Cashier User", "cashier"},
		{uuid.New().String(), "cashier2@pos.local", "Cashier123!", "Cashier Two", "cashier"},
	}

	userIDs := make(map[string]string)
	for _, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		_, err := db.ExecContext(ctx, `
			INSERT OR REPLACE INTO users (id, email, password_hash, name, role, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, 1, ?, ?)
		`, u.id, u.email, string(hash), u.name, u.role, now, now)
		if err != nil {
			log.Printf("  ❌ Failed to seed user %s: %v", u.email, err)
		} else {
			log.Printf("  ✓ User: %s (%s)", u.email, u.role)
			userIDs[u.role] = u.id
		}
	}

	// ============================================
	// SEED CATEGORIES
	// ============================================
	log.Println("\n[2/6] Seeding categories...")
	categories := []struct {
		id          string
		name        string
		description string
		slug        string
	}{
		{uuid.New().String(), "Food", "Main dishes and meals including rice, noodles, and grilled items", "food"},
		{uuid.New().String(), "Beverages", "Hot and cold drinks including coffee, tea, and fresh juices", "beverages"},
		{uuid.New().String(), "Snacks", "Light bites, chips, chocolates, and quick snacks", "snacks"},
		{uuid.New().String(), "Electronics", "Electronic devices, accessories, and gadgets", "electronics"},
		{uuid.New().String(), "Household", "Household items, cleaning supplies, and daily essentials", "household"},
		{uuid.New().String(), "Personal Care", "Personal hygiene and care products", "personal-care"},
	}

	categoryIDs := make(map[string]string)
	for _, c := range categories {
		categoryIDs[c.slug] = c.id
		_, err := db.ExecContext(ctx, `
			INSERT OR REPLACE INTO categories (id, name, description, slug, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, 1, ?, ?)
		`, c.id, c.name, c.description, c.slug, now, now)
		if err != nil {
			log.Printf("  ❌ Failed to seed category %s: %v", c.name, err)
		} else {
			log.Printf("  ✓ Category: %s", c.name)
		}
	}

	// ============================================
	// SEED PRODUCTS
	// ============================================
	log.Println("\n[3/6] Seeding products...")
	products := []struct {
		id           string
		categorySlug string
		sku          string
		name         string
		description  string
		price        float64
		stock        int
	}{
		// Food
		{uuid.New().String(), "food", "FOOD-001", "Nasi Goreng Spesial", "Fried rice with egg, chicken, and vegetables", 28000, 100},
		{uuid.New().String(), "food", "FOOD-002", "Mie Goreng Seafood", "Fried noodles with shrimp and squid", 32000, 80},
		{uuid.New().String(), "food", "FOOD-003", "Ayam Bakar Madu", "Honey grilled chicken with rice and sambal", 38000, 50},
		{uuid.New().String(), "food", "FOOD-004", "Nasi Uduk Komplit", "Coconut rice with fried chicken and sides", 25000, 60},
		{uuid.New().String(), "food", "FOOD-005", "Soto Ayam", "Traditional chicken soup with rice cake", 22000, 70},

		// Beverages
		{uuid.New().String(), "beverages", "BEV-001", "Kopi Hitam", "Strong black coffee", 10000, 200},
		{uuid.New().String(), "beverages", "BEV-002", "Kopi Susu Gula Aren", "Coffee with milk and palm sugar", 18000, 150},
		{uuid.New().String(), "beverages", "BEV-003", "Es Teh Manis", "Sweet iced tea", 8000, 200},
		{uuid.New().String(), "beverages", "BEV-004", "Es Jeruk Segar", "Fresh orange juice with ice", 15000, 100},
		{uuid.New().String(), "beverages", "BEV-005", "Matcha Latte", "Japanese matcha green tea latte", 25000, 80},

		// Snacks
		{uuid.New().String(), "snacks", "SNK-001", "Keripik Singkong", "Crispy cassava chips", 12000, 150},
		{uuid.New().String(), "snacks", "SNK-002", "Cokelat Silverqueen", "Premium chocolate bar 65g", 18000, 100},
		{uuid.New().String(), "snacks", "SNK-003", "Kacang Garuda", "Roasted peanuts 100g", 15000, 120},
		{uuid.New().String(), "snacks", "SNK-004", "Pocky Strawberry", "Strawberry coated biscuit sticks", 14000, 80},

		// Electronics
		{uuid.New().String(), "electronics", "ELEC-001", "USB Cable Type-C 1m", "Fast charging USB-C cable", 35000, 50},
		{uuid.New().String(), "electronics", "ELEC-002", "Power Bank 10000mAh", "Portable charger with dual USB", 180000, 25},
		{uuid.New().String(), "electronics", "ELEC-003", "Earphone Basic", "Wired earphone with mic", 45000, 40},
		{uuid.New().String(), "electronics", "ELEC-004", "Phone Holder", "Universal car phone holder", 55000, 30},

		// Household
		{uuid.New().String(), "household", "HH-001", "Tissue Paseo 250s", "Facial tissue box", 18000, 100},
		{uuid.New().String(), "household", "HH-002", "Sabun Cuci Piring", "Dish soap 500ml", 12000, 80},
		{uuid.New().String(), "household", "HH-003", "Kantong Plastik L", "Large plastic bags 50pcs", 8000, 200},

		// Personal Care
		{uuid.New().String(), "personal-care", "PC-001", "Hand Sanitizer 100ml", "Antibacterial hand sanitizer", 25000, 150},
		{uuid.New().String(), "personal-care", "PC-002", "Masker Medis 50pcs", "Disposable medical masks", 45000, 100},
		{uuid.New().String(), "personal-care", "PC-003", "Sabun Mandi Lifebuoy", "Antibacterial bath soap", 8000, 200},
	}

	productMap := make(map[string]struct {
		id    string
		name  string
		price float64
	})
	for _, p := range products {
		categoryID := categoryIDs[p.categorySlug]
		_, err := db.ExecContext(ctx, `
			INSERT OR REPLACE INTO products (id, category_id, sku, name, description, price, stock, image_url, is_active, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, '', 1, ?, ?)
		`, p.id, categoryID, p.sku, p.name, p.description, p.price, p.stock, now, now)
		if err != nil {
			log.Printf("  ❌ Failed to seed product %s: %v", p.name, err)
		} else {
			log.Printf("  ✓ Product: %s (Rp %,.0f)", p.name, p.price)
			productMap[p.sku] = struct {
				id    string
				name  string
				price float64
			}{p.id, p.name, p.price}
		}
	}

	// ============================================
	// SEED CUSTOMERS
	// ============================================
	log.Println("\n[4/6] Seeding customers...")
	customers := []struct {
		id      string
		name    string
		email   string
		phone   string
		address string
		points  int
	}{
		{uuid.New().String(), "John Doe", "john.doe@email.com", "081234567890", "Jl. Merdeka No. 1, Jakarta", 150},
		{uuid.New().String(), "Jane Smith", "jane.smith@email.com", "081234567891", "Jl. Sudirman No. 25, Jakarta", 200},
		{uuid.New().String(), "Bob Wilson", "bob.wilson@email.com", "081234567892", "Jl. Thamrin No. 10, Jakarta", 75},
		{uuid.New().String(), "Alice Brown", "alice.brown@email.com", "081234567893", "Jl. Gatot Subroto No. 50, Jakarta", 300},
		{uuid.New().String(), "Charlie Davis", "charlie.d@email.com", "081234567894", "Jl. Kuningan No. 15, Jakarta", 50},
		{uuid.New().String(), "Diana Lee", "diana.lee@email.com", "081234567895", "Jl. Rasuna Said No. 8, Jakarta", 100},
		{uuid.New().String(), "Edward Kim", "edward.kim@email.com", "081234567896", "Jl. Casablanca No. 20, Jakarta", 25},
	}

	customerIDs := make([]string, 0)
	for _, c := range customers {
		_, err := db.ExecContext(ctx, `
			INSERT OR REPLACE INTO customers (id, name, email, phone, address, loyalty_points, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`, c.id, c.name, c.email, c.phone, c.address, c.points, now, now)
		if err != nil {
			log.Printf("  ❌ Failed to seed customer %s: %v", c.name, err)
		} else {
			log.Printf("  ✓ Customer: %s (%d points)", c.name, c.points)
			customerIDs = append(customerIDs, c.id)
		}
	}

	// ============================================
	// SEED TRANSACTIONS
	// ============================================
	log.Println("\n[5/6] Seeding transactions...")

	// Create sample transactions for the past 30 days
	cashierID := userIDs["cashier"]
	transactionData := []struct {
		daysAgo       int
		customerIdx   int
		paymentMethod string
		items         []struct {
			sku string
			qty int
		}
	}{
		// Today's transactions
		{0, 0, "cash", []struct {
			sku string
			qty int
		}{{"FOOD-001", 2}, {"BEV-001", 2}}},
		{0, 1, "card", []struct {
			sku string
			qty int
		}{{"FOOD-003", 1}, {"BEV-002", 1}, {"SNK-002", 2}}},
		{0, -1, "ewallet", []struct {
			sku string
			qty int
		}{{"BEV-003", 3}, {"SNK-001", 1}}},

		// Yesterday
		{1, 2, "cash", []struct {
			sku string
			qty int
		}{{"FOOD-002", 1}, {"BEV-004", 1}}},
		{1, 3, "card", []struct {
			sku string
			qty int
		}{{"ELEC-001", 2}, {"PC-001", 1}}},
		{1, 0, "ewallet", []struct {
			sku string
			qty int
		}{{"FOOD-004", 2}, {"BEV-005", 2}}},

		// 2 days ago
		{2, 4, "cash", []struct {
			sku string
			qty int
		}{{"FOOD-005", 3}, {"BEV-001", 3}}},
		{2, 5, "card", []struct {
			sku string
			qty int
		}{{"SNK-003", 2}, {"SNK-004", 2}, {"BEV-002", 1}}},

		// 3 days ago
		{3, 1, "ewallet", []struct {
			sku string
			qty int
		}{{"HH-001", 2}, {"HH-002", 1}, {"PC-003", 3}}},
		{3, 6, "cash", []struct {
			sku string
			qty int
		}{{"FOOD-001", 1}, {"FOOD-002", 1}, {"BEV-003", 2}}},

		// 5 days ago
		{5, 2, "card", []struct {
			sku string
			qty int
		}{{"ELEC-002", 1}, {"ELEC-003", 1}}},
		{5, 3, "ewallet", []struct {
			sku string
			qty int
		}{{"FOOD-003", 2}, {"BEV-004", 2}, {"SNK-001", 1}}},

		// 7 days ago
		{7, 0, "cash", []struct {
			sku string
			qty int
		}{{"FOOD-001", 3}, {"BEV-001", 3}, {"BEV-002", 2}}},
		{7, 4, "card", []struct {
			sku string
			qty int
		}{{"PC-001", 3}, {"PC-002", 1}}},

		// 14 days ago
		{14, 5, "ewallet", []struct {
			sku string
			qty int
		}{{"FOOD-004", 2}, {"FOOD-005", 1}, {"BEV-005", 2}}},
		{14, 1, "cash", []struct {
			sku string
			qty int
		}{{"ELEC-004", 2}, {"HH-003", 5}}},

		// 21 days ago
		{21, 2, "card", []struct {
			sku string
			qty int
		}{{"FOOD-002", 2}, {"FOOD-003", 1}, {"BEV-001", 3}}},
		{21, 6, "ewallet", []struct {
			sku string
			qty int
		}{{"SNK-002", 3}, {"SNK-003", 2}, {"BEV-003", 2}}},

		// 30 days ago
		{30, 0, "cash", []struct {
			sku string
			qty int
		}{{"FOOD-001", 4}, {"BEV-002", 4}}},
		{30, 3, "card", []struct {
			sku string
			qty int
		}{{"ELEC-001", 1}, {"ELEC-002", 1}, {"PC-001", 2}}},
	}

	transactionCount := 0
	for _, td := range transactionData {
		txID := uuid.New().String()
		txDate := now.AddDate(0, 0, -td.daysAgo)
		invoiceNum := fmt.Sprintf("INV-%s-%s", txDate.Format("20060102"), txID[:8])

		var customerID *string
		if td.customerIdx >= 0 && td.customerIdx < len(customerIDs) {
			customerID = &customerIDs[td.customerIdx]
		}

		// Calculate totals
		var subtotal float64
		var items []struct {
			id       string
			prodID   string
			prodName string
			price    float64
			qty      int
			subtotal float64
		}

		for _, item := range td.items {
			prod := productMap[item.sku]
			itemSubtotal := prod.price * float64(item.qty)
			items = append(items, struct {
				id       string
				prodID   string
				prodName string
				price    float64
				qty      int
				subtotal float64
			}{
				uuid.New().String(),
				prod.id,
				prod.name,
				prod.price,
				item.qty,
				itemSubtotal,
			})
			subtotal += itemSubtotal
		}

		taxAmount := subtotal * 0.10
		totalAmount := subtotal + taxAmount

		// Insert transaction
		_, err := db.ExecContext(ctx, `
			INSERT INTO transactions (id, user_id, customer_id, invoice_number, subtotal, tax_amount, 
				discount_amount, total_amount, payment_method, status, notes, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, 0, ?, ?, 'completed', '', ?, ?)
		`, txID, cashierID, customerID, invoiceNum, subtotal, taxAmount, totalAmount, td.paymentMethod, txDate, txDate)
		if err != nil {
			log.Printf("  ❌ Failed to create transaction: %v", err)
			continue
		}

		// Insert transaction items
		for _, item := range items {
			_, err := db.ExecContext(ctx, `
				INSERT INTO transaction_items (id, transaction_id, product_id, product_name, unit_price, quantity, subtotal, created_at)
				VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			`, item.id, txID, item.prodID, item.prodName, item.price, item.qty, item.subtotal, txDate)
			if err != nil {
				log.Printf("  ❌ Failed to create transaction item: %v", err)
			}
		}

		transactionCount++
		log.Printf("  ✓ Transaction: %s (Rp %,.0f) - %s", invoiceNum, totalAmount, txDate.Format("2006-01-02"))
	}

	// ============================================
	// SUMMARY
	// ============================================
	log.Println("\n[6/6] Summary...")
	log.Println("===============================================")
	log.Println("✅ Database seeding completed successfully!")
	log.Println("===============================================")
	log.Println("")
	log.Printf("  📊 Users:        %d", len(users))
	log.Printf("  📊 Categories:   %d", len(categories))
	log.Printf("  📊 Products:     %d", len(products))
	log.Printf("  📊 Customers:    %d", len(customers))
	log.Printf("  📊 Transactions: %d", transactionCount)
	log.Println("")
	log.Println("🔐 Dummy users for testing:")
	log.Println("   ┌─────────────────────────┬─────────────┬──────────┐")
	log.Println("   │ Email                   │ Password    │ Role     │")
	log.Println("   ├─────────────────────────┼─────────────┼──────────┤")
	log.Println("   │ admin@pos.local         │ Admin123!   │ admin    │")
	log.Println("   │ manager@pos.local       │ Manager123! │ manager  │")
	log.Println("   │ cashier@pos.local       │ Cashier123! │ cashier  │")
	log.Println("   │ cashier2@pos.local      │ Cashier123! │ cashier  │")
	log.Println("   └─────────────────────────┴─────────────┴──────────┘")
	log.Println("")
	log.Println("Run the server: go run ./cmd/api/main.go")
}
