# Request Additional Endpoints

Berikut adalah daftar endpoint yang dibutuhkan frontend namun belum tersedia di Swagger API saat ini.

---

## 1. Dashboard Statistics

### `GET /api/v1/reports/dashboard/stats`

**Description:** Mengambil statistik ringkasan untuk halaman dashboard.

**Response:**

```json
{
  "success": true,
  "data": {
    "today_sales": {
      "amount": 4821.5,
      "change_percent": 14.2,
      "comparison": "yesterday"
    },
    "active_orders": {
      "count": 24,
      "avg_prep_time_minutes": 12
    },
    "net_margin": {
      "percent": 32.4,
      "change_percent": 2.1,
      "comparison": "last_week"
    }
  }
}
```

---

## 2. Live/Realtime Sales Data

### `GET /api/v1/reports/sales/realtime`

**Description:** Mengambil data penjualan real-time untuk grafik live sales di dashboard.

**Query Parameters:**

- `interval`: `hourly` | `15min` (default: `hourly`)
- `date`: `YYYY-MM-DD` (default: today)

**Response:**

```json
{
  "success": true,
  "data": {
    "interval": "hourly",
    "data_points": [
      { "time": "08:00", "amount": 250.0 },
      { "time": "09:00", "amount": 480.0 },
      { "time": "10:00", "amount": 620.0 },
      { "time": "11:00", "amount": 890.0 },
      { "time": "12:00", "amount": 1150.0 }
    ]
  }
}
```

---

## 3. Recent Transactions/Activity

### `GET /api/v1/transactions/recent`

**Description:** Mengambil transaksi terbaru untuk widget Recent Activity di dashboard.

**Query Parameters:**

- `limit`: number (default: 5)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "order_number": "2034",
      "table_or_type": "Table 4",
      "total_amount": 42.5,
      "status": "paid",
      "created_at": "2024-01-25T10:30:00Z",
      "time_ago": "Just now"
    }
  ]
}
```

---

## 4. System Health Status

### `GET /api/v1/system/health/detailed`

**Description:** Mengambil status kesehatan sistem yang detail untuk widget System Health.

**Response:**

```json
{
  "success": true,
  "data": {
    "terminals": [
      { "id": "POS-01", "type": "pos", "latency_ms": 35, "status": "online" },
      { "id": "POS-02", "type": "pos", "latency_ms": 42, "status": "online" },
      { "id": "KDS-01", "type": "kds", "latency_ms": 32, "status": "online" }
    ],
    "servers": [
      { "id": "Server A", "uptime_percent": 99.9, "status": "online" },
      { "id": "Server B", "uptime_percent": 99.5, "status": "online" }
    ],
    "integrations": [
      { "name": "Stripe", "status": "ok" },
      { "name": "Twilio", "status": "ok" },
      { "name": "SendGrid", "status": "ok" }
    ]
  }
}
```

---

## 5. Customer Statistics

### `GET /api/v1/customers/stats`

**Description:** Mengambil statistik pelanggan untuk halaman Customer Management.

**Response:**

```json
{
  "success": true,
  "data": {
    "total_customers": 4821,
    "change_percent": 3.2,
    "comparison": "last_month",
    "new_this_month": 156,
    "new_change_percent": 12.8,
    "avg_loyalty_points": 845
  }
}
```

---

## 6. Product/Inventory Statistics

### `GET /api/v1/products/stats`

**Description:** Mengambil statistik produk dan inventori.

**Response:**

```json
{
  "success": true,
  "data": {
    "total_sku": 842,
    "low_stock_count": 14,
    "out_of_stock_count": 3,
    "inventory_value": 42850.0,
    "value_change_percent": 2.4,
    "comparison": "last_week"
  }
}
```

---

## 7. Stock Movement Log

### `GET /api/v1/products/stock-movements`

**Description:** Mengambil log pergerakan stok untuk halaman Product Manager.

**Query Parameters:**

- `page`: number (default: 1)
- `per_page`: number (default: 10)
- `product_id`: uuid (optional, filter by product)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "product_id": "uuid",
      "product_name": "Premium Roast",
      "sku": "SKU-001",
      "type": "restock",
      "quantity_change": 50,
      "new_balance": 124,
      "user": "Admin",
      "created_at": "2024-01-25T10:45:12Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

**Note:** Type values: `sale`, `restock`, `waste`, `adjustment`, `return`

---

## 8. Category Statistics

### `GET /api/v1/categories/stats`

**Description:** Mengambil statistik kategori untuk halaman Category Manager.

**Response:**

```json
{
  "success": true,
  "data": {
    "total_categories": 24,
    "new_this_week": 2,
    "active_categories": 18,
    "coverage_percent": 75,
    "most_popular": {
      "id": "uuid",
      "name": "Coffee Specials",
      "items_sold": 1240
    }
  }
}
```

---

## 9. Category Activity Log

### `GET /api/v1/categories/activity-log`

**Description:** Mengambil log aktivitas kategori.

**Query Parameters:**

- `limit`: number (default: 10)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "event": "Created \"Seasonal Specials\"",
      "user": "J. Doe",
      "details": "New ID: CAT-001",
      "created_at": "2024-01-25T11:05:22Z"
    }
  ]
}
```

---

## 10. User Profile Update

### `PUT /api/v1/auth/me`

**Description:** Update profil user yang sedang login.

**Request:**

```json
{
  "name": "Jane Doe",
  "email": "j.doe@bentopos.system",
  "phone": "+1 (555) 019-2834"
}
```

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "Jane Doe",
    "email": "j.doe@bentopos.system",
    "phone": "+1 (555) 019-2834",
    "role": "admin",
    "is_active": true
  }
}
```

---

## 11. User Activity Log

### `GET /api/v1/auth/me/activity`

**Description:** Mengambil log aktivitas user yang sedang login.

**Query Parameters:**

- `limit`: number (default: 10)

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "action": "User Login",
      "device_ip": "192.168.1.1 (Chrome)",
      "status": "success",
      "created_at": "2024-01-25T10:42:05Z"
    }
  ]
}
```

---

## 12. Auth Refresh Token

### `POST /api/v1/auth/refresh`

**Description:** Refresh access token menggunakan refresh token.

**Request:** (refresh token dari httpOnly cookie)

**Response:**

```json
{
  "success": true,
  "data": {
    "access_token": "new_jwt_token",
    "expires_in": 3600
  }
}
```

---

## 13. Auth Logout

### `POST /api/v1/auth/logout`

**Description:** Logout user dan invalidate refresh token.

**Response:**

```json
{
  "success": true,
  "message": "Successfully logged out"
}
```
