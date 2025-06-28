# 🔗 Go-URLX – Multi-Tenant URL Shortener with Authentication

A production-ready Golang-based URL Shortener with:
- ✅ User Signup/Login with JWT
- ✅ Multi-Tenant RBAC (Admin/User)
- ✅ Redis Caching for Performance
- ✅ PostgreSQL as Primary DB
- ✅ Fiber Web Framework
- ✅ Rate Limiting, Token Rotation
- ✅ Fully Deployed on Render

## Why This Project?
A production-ready URL shortening service built with **Go** to showcase:
✅ Real-world auth, RBAC, rate-limiting, and Redis caching
✅ Deployable in cloud environments like Render, Docker, AWS


## 🌐 Live Demo
👉 [https://go-shorturl.onrender.com](https://go-shorturl.onrender.com)  
📂 Example Credentials:

| Role  | Email                | Password |
|-------|----------------------|----------|
| Admin | Prajwal@gmail.com    | 1234     |
| User  | Praful@gmail.com     | 1234     |


## 🗃️ PostgreSQL Schema

You can find the schema SQL in [`init.sql`](./pkg/db/init.sql)  
Includes: `tenants`, `users`, `urls` with proper constraints and sample updates.

---

## 📁 Project Structure

/config → Load env vars
/controllers → All handler logic
/models → DB queries (PostgreSQL)
/routes → API routes setup
/middleware → JWT, RBAC, RateLimit
/db → PostgreSQL pool setup
/utils → JWT, Redis, Validator
main.go → App entrypoint

## ⚙️ Features

- 🔐 JWT-based Auth (Access + Refresh tokens)
- 👥 Role-based Access Control (RBAC)
- 🌐 Shorten + Redirect URL support
- 📈 Stats API (click count, timestamps)
- 💥 Rate Limiting (per user)
- 🗃️ Multi-tenant architecture
- 💾 Redis used for caching + performance
- 🐘 PostgreSQL with pgx driver
- 🚀 Deployed on Render (Free Tier)

---

## 🛠️ Tech Stack

- **Language**: Go (Golang)
- **Web Framework**: Fiber
- **Database**: PostgreSQL (pgx)
- **Cache**: Redis (Upstash)
- **Auth**: JWT (access & refresh)
- **Rate Limiting**: Redis + Fiber middleware
- **ORM**: Raw pgx queries
- **Deployment**: Render
- **Other**: Docker (coming soon), pgAdmin

---

## 🧪 API Endpoints

| Method | Endpoint             | Description |
|--------|----------|-------------|
| POST   | `/user/signup` | Register new user |
| POST   | `/user/login` | Login with JWT |
| POST   | `/url/shorten` | Create short URL |
| GET    | `/url/redirect/:shortCode` | Redirect to original URL |
| GET    | `/url/stats` | Get click stats |
| GET    | `/url/AllUrls` | Admin: View all URLs |
| POST   | `/url/deleteShortCode` | Delete a short code |
| POST   | `/tenant/createNewTenant` | Create tenant (admin only) |

---

## 🧑‍💻 Local Setup Instructions

```bash
git clone https://github.com/prajwalzaware/Go_ShortUrl.git
cd Go_ShortUrl

# Setup .env
cp .env.example .env

# Edit your .env file with DB, Redis URLs
# Example:
# DB_URL=postgres://username:password@host:port/dbname
# REDIS_URL=your-upstash-url
# JWT_SECRET=some-secret-key
# BASE_URL=https://yourdomain.com/

# Run the server
go run main.go

