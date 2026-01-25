# Backend API Endpoint Requirements

Dokumen ini berisi daftar lengkap endpoint yang dibutuhkan frontend.

---

## 📌 Status Legend

- ✅ **Sudah ada di Swagger** - Endpoint sudah tersedia
- ⚠️ **Perlu verifikasi** - Endpoint ada tapi perlu dicek/diperbaiki
- ❌ **Belum ada** - Endpoint baru yang perlu dibuat

---

## 🔐 Role & Permission Matrix

| Fitur           | Admin   | Manager        | Cashier        |
| --------------- | ------- | -------------- | -------------- |
| Dashboard       | ✅ Full | ✅ Full        | ✅ Read Only   |
| POS (Transaksi) | ✅ Full | ✅ Full        | ✅ Full        |
| Products        | ✅ CRUD | ✅ CRUD        | 👁️ Read Only   |
| Categories      | ✅ CRUD | ✅ CRUD        | 👁️ Read Only   |
| Customers       | ✅ CRUD | ✅ CRUD        | 👁️ Read Only   |
| Transactions    | ✅ Full | ✅ Full        | 👁️ Read Own    |
| Reports         | ✅ Full | ❌ No Access   | ❌ No Access   |
| Settings        | ✅ Full | ✅ Own Profile | ✅ Own Profile |
| User Management | ✅ Full | ❌ No Access   | ❌ No Access   |

---

## 🧪 Seed Accounts (untuk Testing)

Backend perlu seed 3 akun test:

```
Admin:
  email: admin@gopos.local
  password: Admin123!
  role: admin

Manager:
  email: manager@gopos.local
  password: Manager123!
  role: manager

Cashier:
  email: cashier@gopos.local
  password: Cashier123!
  role: cashier
```

---

# ENDPOINT DETAILS

---

## 1. POS / Point of Sale ❌ NEW

### 1.1 `GET /api/v1/pos/products`

**Description:** Mengambil daftar produk untuk tampilan POS (dengan stok dan harga).

**Query Parameters:**

- `category_id`: uuid (optional)
- `search`: string (optional)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "Americano",
      "sku": "BEV-001",
      "price": 25000,
      "stock": 100,
      "category_id": "uuid",
      "category_name": "Beverages",
      "image_url": "/images/americano.jpg"
    }
  ]
}
```

---

### 1.2 `POST /api/v1/pos/transactions`

**Description:** Membuat transaksi penjualan baru dari POS.

**Request:**

```json
{
  "customer_id": "uuid",
  "payment_method": "cash",
  "items": [
    {
      "product_id": "uuid",
      "quantity": 2,
      "unit_price": 25000,
      "discount": 0
    }
  ],
  "subtotal": 50000,
  "tax_amount": 5000,
  "discount_amount": 0,
  "total_amount": 55000,
  "amount_paid": 60000,
  "change_amount": 5000,
  "notes": "No sugar"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "invoice_number": "INV-20250125-001",
    "status": "completed",
    "total_amount": 55000,
    "created_at": "2025-01-25T10:30:00Z"
  }
}
```

---

### 1.3 `GET /api/v1/pos/hold`

**Description:** Mengambil daftar transaksi yang di-hold/pending.

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "hold_number": "HOLD-001",
      "customer_name": "John Doe",
      "items_count": 3,
      "total_amount": 75000,
      "created_at": "2025-01-25T10:00:00Z"
    }
  ]
}
```

---

### 1.4 `POST /api/v1/pos/hold`

**Description:** Menyimpan transaksi sementara (hold).

**Request:**

```json
{
  "customer_id": "uuid",
  "items": [...],
  "notes": "Customer will return at 3pm"
}
```

---

## 2. Reports - Sales ⚠️ FIX NEEDED

### 2.1 `GET /api/v1/reports/sales/daily` ✅

**Query Parameters:**

- `date_from`: YYYY-MM-DD
- `date_to`: YYYY-MM-DD

**Response:**

```json
{
  "success": true,
  "data": {
    "period": "daily",
    "date_from": "2025-01-01",
    "date_to": "2025-01-25",
    "total_revenue": 15000000,
    "total_transactions": 450,
    "avg_order_value": 33333,
    "chart_data": [
      { "date": "2025-01-01", "revenue": 500000, "transactions": 15 },
      { "date": "2025-01-02", "revenue": 620000, "transactions": 18 }
    ]
  }
}
```

---

### 2.2 `GET /api/v1/reports/sales/weekly` ❌ NEW

**Query Parameters:**

- `date_from`: YYYY-MM-DD
- `date_to`: YYYY-MM-DD

**Response:**

```json
{
  "success": true,
  "data": {
    "period": "weekly",
    "total_revenue": 45000000,
    "total_transactions": 1350,
    "avg_order_value": 33333,
    "chart_data": [
      {
        "week": "2025-W01",
        "week_start": "2025-01-01",
        "revenue": 5000000,
        "transactions": 150
      },
      {
        "week": "2025-W02",
        "week_start": "2025-01-08",
        "revenue": 5500000,
        "transactions": 165
      }
    ]
  }
}
```

---

### 2.3 `GET /api/v1/reports/sales/monthly` ✅

**Response:**

```json
{
  "success": true,
  "data": {
    "period": "monthly",
    "year": 2025,
    "total_revenue": 180000000,
    "total_transactions": 5400,
    "avg_order_value": 33333,
    "chart_data": [
      { "month": "2025-01", "revenue": 15000000, "transactions": 450 },
      { "month": "2025-02", "revenue": 16000000, "transactions": 480 }
    ]
  }
}
```

---

### 2.4 `GET /api/v1/reports/categories/performance` ❌ NEW

**Description:** Performa penjualan per kategori.

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "category_id": "uuid",
      "category_name": "Beverages",
      "total_sold": 1250,
      "total_revenue": 31250000,
      "percentage": 45.5
    },
    {
      "category_id": "uuid",
      "category_name": "Food",
      "total_sold": 890,
      "total_revenue": 22250000,
      "percentage": 32.4
    }
  ]
}
```

---

### 2.5 `GET /api/v1/reports/products/top` ✅

**Note:** Endpoint sudah ada, tapi return 403 untuk cashier. Pastikan permission sesuai (hanya admin).

**Query Parameters:**

- `limit`: number (default: 10)
- `period`: `day` | `week` | `month` (default: month)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "product_id": "uuid",
      "product_name": "Americano",
      "sku": "BEV-001",
      "quantity_sold": 450,
      "total_revenue": 11250000,
      "rank": 1
    }
  ]
}
```

---

## 3. Notifications ❌ NEW

### 3.1 `GET /api/v1/notifications`

**Description:** Mengambil notifikasi untuk user yang login.

**Query Parameters:**

- `limit`: number (default: 10)
- `unread_only`: boolean (default: false)

**Response:**

```json
{
  "success": true,
  "data": {
    "unread_count": 5,
    "notifications": [
      {
        "id": "uuid",
        "type": "low_stock",
        "title": "Low Stock Alert",
        "message": "Americano stock is low (5 remaining)",
        "is_read": false,
        "created_at": "2025-01-25T10:30:00Z",
        "action_url": "/products?id=uuid"
      },
      {
        "id": "uuid",
        "type": "new_order",
        "title": "New Order",
        "message": "Order #2034 received",
        "is_read": true,
        "created_at": "2025-01-25T10:25:00Z",
        "action_url": "/transactions/uuid"
      }
    ]
  }
}
```

**Notification Types:**

- `low_stock` - Stok produk menipis
- `out_of_stock` - Stok produk habis
- `new_order` - Pesanan baru masuk
- `order_completed` - Pesanan selesai
- `system` - Notifikasi sistem

---

### 3.2 `PUT /api/v1/notifications/:id/read`

**Description:** Menandai notifikasi sudah dibaca.

**Response:**

```json
{
  "success": true,
  "message": "Notification marked as read"
}
```

---

### 3.3 `PUT /api/v1/notifications/read-all`

**Description:** Menandai semua notifikasi sudah dibaca.

---

## 4. User Management (Admin Only) ❌ NEW

### 4.1 `GET /api/v1/users`

**Description:** Mengambil daftar semua user (admin only).

**Query Parameters:**

- `role`: `admin` | `manager` | `cashier` (optional filter)
- `page`: number
- `per_page`: number

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "name": "John Doe",
      "email": "john@gopos.local",
      "phone": "+62812345678",
      "role": "manager",
      "is_active": true,
      "last_login_at": "2025-01-25T10:00:00Z",
      "created_at": "2025-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 15,
    "total_pages": 2
  }
}
```

---

### 4.2 `POST /api/v1/users`

**Description:** Membuat user baru (admin only).

**Request:**

```json
{
  "name": "New Employee",
  "email": "employee@gopos.local",
  "password": "Secure123!",
  "phone": "+62812345678",
  "role": "cashier"
}
```

---

### 4.3 `PUT /api/v1/users/:id`

**Description:** Update user (admin only).

**Request:**

```json
{
  "name": "Updated Name",
  "phone": "+62812345678",
  "role": "manager",
  "is_active": true
}
```

---

### 4.4 `DELETE /api/v1/users/:id`

**Description:** Delete/deactivate user (admin only).

---

### 4.5 `PUT /api/v1/users/:id/reset-password`

**Description:** Reset password user (admin only).

**Request:**

```json
{
  "new_password": "NewSecure123!"
}
```

---

## 5. Profile Update ⚠️ FIX NEEDED

### 5.1 `PUT /api/v1/auth/me`

**Issue:** Saat ini endpoint return success tapi data tidak tersimpan.

**Request:**

```json
{
  "name": "Updated Name",
  "phone": "+62812345678"
}
```

**Note:** Email TIDAK boleh diubah (frontend sudah disable field email).

---

## 6. Transaction Statistics ❌ NEW

### 6.1 `GET /api/v1/transactions/stats`

**Description:** Statistik untuk halaman transaksi.

**Response:**

```json
{
  "success": true,
  "data": {
    "total_transactions": 4821,
    "total_revenue": 150000000,
    "avg_order_value": 31125,
    "today_transactions": 45,
    "today_revenue": 1400000
  }
}
```

---

## 7. Endpoints yang Sudah Ada ✅

Berikut endpoint yang sudah tersedia dan berfungsi:

| Endpoint                          | Method         | Status                    |
| --------------------------------- | -------------- | ------------------------- |
| `/api/v1/auth/login`              | POST           | ✅ OK                     |
| `/api/v1/auth/register`           | POST           | ✅ OK                     |
| `/api/v1/auth/logout`             | POST           | ✅ OK                     |
| `/api/v1/auth/refresh`            | POST           | ✅ OK                     |
| `/api/v1/auth/me`                 | GET            | ✅ OK                     |
| `/api/v1/auth/me/activity`        | GET            | ✅ OK                     |
| `/api/v1/categories`              | GET/POST       | ✅ OK (check permissions) |
| `/api/v1/categories/:id`          | GET/PUT/DELETE | ✅ OK (check permissions) |
| `/api/v1/categories/stats`        | GET            | ✅ OK                     |
| `/api/v1/products`                | GET/POST       | ✅ OK (check permissions) |
| `/api/v1/products/:id`            | GET/PUT/DELETE | ✅ OK (check permissions) |
| `/api/v1/products/stats`          | GET            | ✅ OK                     |
| `/api/v1/customers`               | GET/POST       | ✅ OK (check permissions) |
| `/api/v1/customers/:id`           | GET/PUT/DELETE | ✅ OK (check permissions) |
| `/api/v1/customers/stats`         | GET            | ✅ OK                     |
| `/api/v1/transactions`            | GET/POST       | ✅ OK                     |
| `/api/v1/transactions/:id`        | GET            | ✅ OK                     |
| `/api/v1/transactions/recent`     | GET            | ✅ OK                     |
| `/api/v1/reports/dashboard/stats` | GET            | ✅ OK                     |
| `/api/v1/reports/sales/realtime`  | GET            | ✅ OK                     |
| `/api/v1/system/health/detailed`  | GET            | ✅ OK                     |

---

## Summary: New Endpoints Needed

| #   | Endpoint                                     | Priority  |
| --- | -------------------------------------------- | --------- |
| 1   | `GET /api/v1/pos/products`                   | 🔴 High   |
| 2   | `POST /api/v1/pos/transactions`              | 🔴 High   |
| 3   | `GET /api/v1/pos/hold`                       | 🟡 Medium |
| 4   | `POST /api/v1/pos/hold`                      | 🟡 Medium |
| 5   | `GET /api/v1/reports/sales/weekly`           | 🔴 High   |
| 6   | `GET /api/v1/reports/categories/performance` | 🔴 High   |
| 7   | `GET /api/v1/notifications`                  | 🟡 Medium |
| 8   | `PUT /api/v1/notifications/:id/read`         | 🟡 Medium |
| 9   | `PUT /api/v1/notifications/read-all`         | 🟡 Medium |
| 10  | `GET /api/v1/users`                          | 🔴 High   |
| 11  | `POST /api/v1/users`                         | 🔴 High   |
| 12  | `PUT /api/v1/users/:id`                      | 🔴 High   |
| 13  | `DELETE /api/v1/users/:id`                   | 🔴 High   |
| 14  | `PUT /api/v1/users/:id/reset-password`       | 🟡 Medium |
| 15  | `GET /api/v1/transactions/stats`             | 🟡 Medium |

---

## Fixes Needed

| #   | Endpoint                           | Issue                                                 |
| --- | ---------------------------------- | ----------------------------------------------------- |
| 1   | `PUT /api/v1/auth/me`              | Data tidak tersimpan setelah update                   |
| 2   | All CRUD endpoints                 | Permission 403 untuk role yang seharusnya punya akses |
| 3   | `GET /api/v1/reports/products/top` | 403 untuk admin                                       |
