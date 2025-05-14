# 🛍️ Product Search API (Go + Bleve)

A high-performance full-text search API written in **Go**, capable of searching over **1 million in-memory products** using **Bleve** for indexing. Built with **Chi Router**, **Dockerized** for deployment, and designed with **graceful shutdown**, this project simulates a real-world e-commerce search microservice.


📄 [Download DOCUMENTATION { Decision Log (DOCX) }](https://github.com/ak2k4/product-search/DecisionLog_Atharv-Kaushik.docx)


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
```
Visit in browser:
```bash
http://localhost:8080/search?q=Books
```
🐳 Docker Setup
📦 Build Docker Image
```bash
docker build -t product-search
```
▶️ Run Container
```bash
docker run -p 8080:8080 product-search
```
Visit in browser:
```bash
http://localhost:8080/search?q=Electronics
```

🔍 API Usage
```bash
GET /search?q=<term>
```
Searches over product Name and Category.

Returns up to 50 matching product names in JSON.

✅ Example:
```bash
http://localhost:8080/search?q=Clothing
```
📤 Output:



![image](https://github.com/user-attachments/assets/d8640441-fb88-4584-86c1-0d4462f77cd8)



.



✅Examples with proof of Search Queries:


```bash
http://localhost:8080/search?q=Category:Electronics&from=100&size=20
```
📤 Output:



![image](https://github.com/user-attachments/assets/84350b67-0c4f-40ac-89e9-f890fd244bd5)


```bash
http://localhost:8080/search?q=ID:99999
```
📤 Output:


![image](https://github.com/user-attachments/assets/daf971b9-0651-4216-b86b-c52d8c29e2b4)




Note: Query is case-insensitive and supports partial matches.




🧠 Architecture Overview
```graphql

.
├── main.go             # App entry point
├── product.go          # Data generation (1M products)
├── search.go           # Bleve index & search logic
├── router.go           # API routes via Chi
├── Dockerfile          # Multi-stage container build
├── products.bleve/     # Bleve index folder
└── README.md           # Project documentation

```

⚙️ How It Works
Data Layer: 1M Product structs with ID, Name, Category.

Index Layer: Indexed using Bleve and stored in products.bleve.

API Layer: Chi exposes a /search endpoint.

Shutdown Layer: Handles SIGINT/SIGTERM and shuts down gracefully using context.

Container Layer: Docker-ready for isolated, reproducible builds.




📊 Performance Metrics
Metric	Value
Indexing Time	~15 seconds (on first run)
Query Response	<100ms average
Max Results	50 per search query
Memory Usage	~150MB (1M products in memory)
Disk Index Size	~250MB (products.bleve)


.


.


🔐 Graceful Shutdown
Listens for OS signals (CTRL+C / SIGTERM)

Waits up to 5 seconds to finish in-progress requests

Shuts down the HTTP server using context.WithTimeout

.

.


📄 License
This project is intended for educational and assessment use.
Feel free to fork and build upon it!

.


.


👨‍💻 Author
GitHub: @ak2k4 (ATHARV KAUSHIK)


Internship Submission — Backend Engineering & Search API Task
