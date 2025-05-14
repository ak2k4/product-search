# 🛍️ Product Search API (Go + Bleve)

A high-performance full-text search API written in **Go**, capable of searching over **1 million in-memory products** using **Bleve** for indexing. Built with **Chi Router**, **Dockerized** for deployment, and designed with **graceful shutdown**, this project simulates a real-world e-commerce search microservice.

---

## 📌 Features

- ⚡ Handles 1,000,000 in-memory product entries
- 🔍 Fast full-text search on product `Name` and `Category`
- 📦 Clean, RESTful `GET /search?q=term` endpoint
- 💾 Disk-persisted Bleve index (`products.bleve`)
- 🐳 Docker-ready for consistent builds and deployment
- 🧵 Graceful shutdown via `context.Context` and signal handling
- 🧪 Lightweight and minimal dependencies

---

## 🚀 Getting Started

### ✅ Requirements

- [Go 1.19+](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/) (optional for containerized usage)

---

### 🔧 Local Setup (Without Docker)

```bash
git clone https://github.com/ak2k4/product-search.git
cd product-search
go mod tidy
go run main.go
---

### 🔧 Setup (With Docker)
```bash
docker build -t product-search .
docker run -p 8080:8080 product-search
