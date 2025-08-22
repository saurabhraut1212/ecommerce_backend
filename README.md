# 🛒 E-commerce Backend (Go + Fiber + MongoDB)

Minimal, production-ready starter for an e-commerce API with **Auth**, **Products**, and **Orders**.
- Stack: **Go**, **Fiber**, **MongoDB**, **JWT**, **bcrypt**
- Structure: clean, modular (config, db, models, repo, handlers, middleware, router)

---

## ⚙️ Quick Start

### 1) Prerequisites
- Go 1.24+ (or newer)
- MongoDB atlas
- (Optional) Air for hot reload

### 2) Clone & Setup
```bash
git clone https://github.com/saurabhraut1212/ecommerce_backend.git
cd ecommerce_backend
go mod tidy
```
### 3) Environment
- PORT=8080
- MONGO_URI=atlas_url
- MONGO_DB=ecommerce
- JWT_SECRET=supersecretkey

### 4) Run
```bash
go run ./cmd/server
# or with Air (if installed): air
```

## Folder Structure
- models: pure data types (no DB or HTTP code).
- repo: DB operations (CRUD), easy to mock/test.
- handlers: parse/validate requests, call repos, return responses.
- middleware: cross-cutting concerns (auth).
- router: central route registry.
  
## Authentication Routes
| Method | Endpoint    | Description       |
| ------ | ----------- | ----------------- |
| POST   | `/register` | Register new user |
| POST   | `/login`    | Login & get JWT   |

## Product Routes
| Method | Endpoint        | Description       |
| ------ | --------------- | ----------------- |
| POST   | `/products`     | Create product    |
| GET    | `/products`     | Get all products  |
| GET    | `/products/:id` | Get product by ID |
| PUT    | `/products/:id` | Update product    |
| DELETE | `/products/:id` | Delete product    |

## Order Routes
| Method | Endpoint      | Description      |
| ------ | ------------- | ---------------- |
| POST   | `/orders`     | Create new order |
| GET    | `/orders`     | Get all orders   |
| GET    | `/orders/:id` | Get order by ID  |
| PUT    | `/orders/:id/status` | Update order     |
| DELETE | `/orders/:id` | Delete order     |

## Postman Testing
https://web.postman.co/workspace/388302e8-5eb7-4c3f-821d-5523c39dad56/collection/26119400-da1f5e96-9041-4cf7-986a-26b27b561ce6?action=share&source=copy-link&creator=26119400

Use the token generated after login in Authorization Header for all Product & Order routes
  
