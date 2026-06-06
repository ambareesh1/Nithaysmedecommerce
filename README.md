# MedCart SaaS

## Project Overview

MedCart SaaS is a medical product ordering platform built for hospitals, clinics, and diagnostic centers. Users can browse medical products, search and filter them, add products to a cart, place orders, and view their order history. An admin dashboard lets an administrator add, view, and delete products.

The project is a full stack application with a **Next.js (React + TypeScript)** frontend and a **GoLang (Gin)** backend, using **PostgreSQL** for storage and **Redis** for caching. It is fully containerized with **Docker**, has **Kubernetes** manifests, and a **GitHub Actions** CI pipeline.

## Why This Project Was Built

This project was built to demonstrate the practical full stack skills required for a Full Stack Developer role. It maps directly to a job description that asks for GoLang, React, Next.js, REST APIs, SQL, Redis, Docker, Kubernetes, CI/CD, testing, and eCommerce/SaaS understanding. It is intentionally kept beginner-friendly and clean so it is easy to read, run, and explain in an interview.

## Job Description Skill Mapping

| Skill Required | Where It Is Used in MedCart |
| --- | --- |
| GoLang | Backend API server built with Go and Gin |
| ReactJS | Frontend components and pages |
| NextJS | App Router based frontend |
| REST APIs | Auth, products, cart, and orders endpoints |
| SQL database | PostgreSQL tables and queries |
| Redis caching | Product list caching |
| Docker | Dockerfiles for frontend and backend |
| Docker Compose | Runs the whole stack together |
| Kubernetes | Deployment and service YAML files |
| CI/CD | GitHub Actions workflow |
| Testing | Go tests for health and product handler |
| eCommerce/SaaS | Cart, orders, and admin product management |

## Features

- Landing page with project overview
- Product listing with search, category filter, clear filter, and product count
- Product details page
- Cart with quantity increase/decrease, remove, clear, totals, and place order
- Order history with status badges and item subtotals
- Login and register pages with basic validation
- Admin dashboard with product and order stats
- Admin product management with add form, live preview, list table, and delete
- Backend REST APIs for auth, products, cart, and orders
- JWT based authentication for protected routes
- Redis caching for the product list
- Docker, Docker Compose, and Kubernetes setup
- GitHub Actions CI pipeline

## Tech Stack

**Frontend**
- Next.js (App Router)
- React
- TypeScript
- Tailwind CSS
- localStorage for temporary cart and order storage

**Backend**
- GoLang
- Gin framework
- PostgreSQL
- Redis
- JWT authentication
- REST APIs

**DevOps**
- Docker and Docker Compose
- Kubernetes
- GitHub Actions

## Project Architecture

```
Browser
   |
   v
Next.js Frontend  (localStorage for cart/orders demo)
   |
   |  REST API (JSON)
   v
GoLang Gin Backend
   |          |
   v          v
PostgreSQL   Redis (product cache)
```

The frontend talks to the backend over REST APIs. The backend reads and writes data in PostgreSQL and uses Redis to cache the product list so repeated reads are faster. In the current beginner version, the frontend cart and orders use localStorage so the UI works without login, while the backend already exposes the full API for later integration.

## Folder Structure

```
medcart-saas/
  backend/
    cmd/
      main.go
      main_test.go
    internal/
      config/
      database/
      cache/
      middleware/
      modules/
        auth/
        product/
        cart/
        order/
    Dockerfile
    go.mod
    go.sum

  frontend/
    src/
      app/
        page.tsx
        layout.tsx
        globals.css
        products/
          page.tsx
          [id]/page.tsx
        cart/page.tsx
        orders/page.tsx
        login/page.tsx
        register/page.tsx
        admin/
          page.tsx
          products/page.tsx
      components/
        Header.tsx
        ProductCard.tsx
        FormInput.tsx
        AppButton.tsx
      data/products.ts
      lib/
        cart.ts
        orders.ts
        api.ts
    public/
    Dockerfile
    package.json
    next.config.ts
    tsconfig.json
    tailwind.config.ts

  infra/
    docker-compose.yml
    k8s/

  .github/workflows/ci.yml
  README.md
```

## Frontend Setup

```
cd frontend
npm install
npm run dev
```

Open the app at:

```
http://localhost:3000
```

## Backend Setup

Make sure PostgreSQL and Redis are running locally (or use Docker Compose). Then:

```
cd backend
go mod tidy
go run cmd/main.go
```

Check the health endpoint:

```
http://localhost:8080/health
```

The backend reads configuration from environment variables with sensible defaults:

- `PORT` (default `8080`)
- `DATABASE_URL` (default `postgres://medcart:medcart@localhost:5432/medcart?sslmode=disable`)
- `REDIS_ADDR` (default `localhost:6379`)
- `JWT_SECRET` (default `medcart_secret_key`)

To create the database tables, run the SQL file:

```
backend/internal/database/schema.sql
```

## Docker Setup

The easiest way to run everything (frontend, backend, PostgreSQL, Redis) is Docker Compose:

```
cd infra
docker compose up --build
```

Services and ports:

- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- PostgreSQL: `localhost:5432`
- Redis: `localhost:6379`

The PostgreSQL container automatically loads `schema.sql` on first start.

## Kubernetes Setup

Basic manifests are in `infra/k8s`. After building and loading the images into your cluster, apply them:

```
kubectl apply -f infra/k8s/
```

This creates deployments and services for the frontend, backend, PostgreSQL, and Redis.

## Database Schema

**users**
- id, name, email, password_hash, role, created_at

**products**
- id, name, description, category, price, stock, image_url, created_at

**carts**
- id, user_id, product_id, quantity, created_at

**orders**
- id, user_id, order_number, total_amount, total_quantity, status, created_at

**order_items**
- id, order_id, product_id, quantity, price, subtotal

The full schema is in `backend/internal/database/schema.sql` and includes sample product data.

## API Documentation

**Health**
- `GET /health` - check server status

**Auth**
- `POST /api/auth/register` - register a new user
- `POST /api/auth/login` - login and get a JWT token

**Products**
- `GET /api/products` - list all products (uses Redis cache)
- `GET /api/products/:id` - get one product
- `POST /api/products` - create a product
- `PUT /api/products/:id` - update a product
- `DELETE /api/products/:id` - delete a product

**Cart** (protected, requires JWT)
- `GET /api/cart` - get the current user cart
- `POST /api/cart` - add an item to the cart
- `PUT /api/cart/:productId` - update item quantity
- `DELETE /api/cart/:productId` - remove an item

**Orders** (protected, requires JWT)
- `POST /api/orders` - place an order
- `GET /api/orders/my` - list the current user orders
- `GET /api/orders/:id` - get one order
- `PATCH /api/orders/:id/status` - update order status

Protected routes need an `Authorization: Bearer <token>` header.

## Sample API Payloads

**Register** — `POST /api/auth/register`

```json
{
  "name": "Ambareesh",
  "email": "ambareesh@example.com",
  "password": "password123"
}
```

**Login** — `POST /api/auth/login`

```json
{
  "email": "ambareesh@example.com",
  "password": "password123"
}
```

**Create product** — `POST /api/products`

```json
{
  "name": "Digital Thermometer",
  "category": "Diagnostics",
  "description": "A digital thermometer used to measure body temperature.",
  "price": 499,
  "stock": 25,
  "imageUrl": "https://example.com/thermometer.png"
}
```

**Create order** — `POST /api/orders`

```json
{
  "items": [
    {
      "productId": 1,
      "quantity": 2
    }
  ]
}
```

## Frontend Page Flow

- `/` — Home landing page with links to Products and Admin
- `/products` — Browse, search, and filter products
- `/products/[id]` — View product details and add to cart
- `/cart` — Review cart, change quantities, and place an order
- `/orders` — See order history with status
- `/login` — Login form with basic validation
- `/register` — Register form with basic validation
- `/admin` — Admin dashboard with stats
- `/admin/products` — Add, preview, list, and delete products

## Cart Flow

1. User clicks "Add to Cart" on a product card or details page.
2. The item is saved to localStorage under the key `medcart_items`.
3. The cart page reads items from localStorage.
4. User can increase or decrease quantity, remove an item, or clear the cart.
5. Totals (quantity and amount) are calculated from the items.

## Place Order Flow

1. On the cart page the user clicks "Place Order".
2. `createOrderFromCart` builds an order object with items, subtotals, and totals.
3. The order is saved to localStorage under the key `medcart_orders`.
4. The cart is cleared.
5. The user is redirected to `/orders` to see the new order.

In the backend version, `POST /api/orders` does the same logic against PostgreSQL: it reads product prices, calculates totals, stores the order and order items, and clears the user cart.

## Redis Caching Explanation

The product list endpoint (`GET /api/products`) uses Redis to speed up reads:

1. When a request comes in, the backend first checks Redis for the key `products_list`.
2. If the cache has the products, it returns them directly (marked as `"source": "cache"`).
3. If the cache is empty, it fetches products from PostgreSQL, stores them in Redis with a 5 minute expiry, and returns them (marked as `"source": "database"`).
4. Whenever a product is created, updated, or deleted, the cache is cleared so the next read rebuilds it with fresh data.

This is a simple but realistic caching pattern used in real eCommerce systems to reduce database load.

## CI/CD Explanation

The GitHub Actions workflow in `.github/workflows/ci.yml` runs on every push and pull request. It has two jobs:

- **backend**: checks out the code, sets up Go, runs `go mod tidy`, and runs `go test ./...`.
- **frontend**: checks out the code, sets up Node, installs dependencies, and runs `npm run build`.

The pipeline only validates that the code builds and tests pass. It does not deploy anywhere.

## Testing Explanation

**Backend**
- `backend/cmd/main_test.go` tests the `/health` endpoint returns 200 with the correct JSON.
- `backend/internal/modules/product/handler_test.go` tests product JSON encoding and decoding.

Run backend tests with:

```
cd backend
go test ./...
```

**Frontend**
Frontend testing is optional in this version. To add it you can install React Testing Library and Jest (or Vitest), then write tests for components like `ProductCard` and utility functions like `cart.ts` and `orders.ts`.

## Interview Explanation

When explaining this project in an interview, you can describe it like this:

- **What it is**: MedCart SaaS is a medical eCommerce platform where clinics and hospitals order medical products. It has a customer side (browse, cart, orders) and an admin side (manage products).
- **Frontend**: Built with Next.js App Router and TypeScript. I used reusable components (Header, ProductCard, FormInput, AppButton) and Tailwind CSS for styling. Cart and orders use localStorage so the UI is fully working without a backend, which keeps it easy to demo.
- **Backend**: Built with Go and the Gin framework. It exposes clean REST APIs grouped into modules (auth, product, cart, order). I used JWT for authentication and middleware to protect cart and order routes.
- **Database**: PostgreSQL with five tables (users, products, carts, orders, order_items). The schema models a real ordering system with order items and subtotals.
- **Caching**: Redis caches the product list. Reads check Redis first, fall back to PostgreSQL, then refill the cache, and writes clear the cache. This shows I understand cache invalidation.
- **DevOps**: I containerized both apps with Docker, wired them together with Docker Compose, wrote Kubernetes manifests, and set up a GitHub Actions CI pipeline that builds the frontend and runs Go tests.
- **Why it matters**: It maps one to one with the job description and shows I can work across the full stack, from UI to API to database to deployment.

## Future Improvements

- Connect the frontend cart and orders to the backend APIs instead of localStorage
- Add real admin authentication and role based access on the frontend
- Add product images and image upload
- Add pagination and sorting to the product list
- Add frontend tests with React Testing Library
- Add database migrations instead of a single schema file
- Add payment integration
- Add proper Kubernetes ingress and secrets management
