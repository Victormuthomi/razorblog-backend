# 🚀 Razorblog Backend: Go-Powered Community Engine
**High-Performance Content Orchestration for Alcodist Hub**

![Language](https://img.shields.io/badge/Language-Go_1.22+-00ADD8?style=for-the-badge&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean_Internal-blue?style=for-the-badge)
![API](https://img.shields.io/badge/Spec-Swagger_/_OpenAPI-85EA2D?style=for-the-badge&logo=swagger)

## 📌 Executive Summary
The **Razorblog Backend** is the high-performance content engine specifically engineered to power the **Community Logs** within the **Alcodist Hub**. Built with **Go (Golang)**, this system prioritizes execution speed, low memory footprint, and strict data integrity.

Unlike experimental sandboxes, this repository follows a production-grade **Internal Architecture**, ensuring that core business logic is decoupled from the transport layer and protected from external misuse.

---

## 🏗️ Architectural Foundations

### 1. The `internal/` Pattern
The system adheres to a clean, encapsulated project layout:
* **`cmd/`**: The main entry point. Handles server initialization and dependency injection.
* **`internal/`**: The core of the engine. Contains the domain models, services, and repository layers that drive the Community Logs.
* **`api/`**: Defines the RESTful surface area and handles request/response lifecycle.

### 2. Live-Reload Development (`.air.toml`)
Integrated with **Air** to provide a seamless developer experience, enabling instant binary recompilation and hot-reloading on every file change.

### 3. Schema Integrity & Migrations
Includes a dedicated `migrations/` directory for reproducible database schema evolution, ensuring the data layer remains consistent across all environments.

---

## 🛠️ Tech Stack
* **Language:** Go (Golang)
* **API Documentation:** Swagger (OpenAPI 3.0)
* **Hot Reload:** Air (`.air.toml`)
* **Infrastructure:** **Docker** & **Docker Compose** for containerized orchestration.
* **Automation:** GitHub Actions for automated linting and build verification.

---

## 🚀 Deployment & Local Setup
The system is ready for containerized execution:

```bash
# Clone the Engine
git clone [https://github.com/Victormuthomi/razorblog-backend.git](https://github.com/Victormuthomi/razorblog-backend.git)

# Install Go Dependencies
go mod download

# Run in Development Mode (Hot Reload with Air)
air

# Orchestrate with Docker
docker-compose up --build
```

## 📊 API Documentation
The API is fully documented via **Swagger**. Once the server is running locally, the interactive documentation is accessible at:
`http://localhost:8080/swagger/index.html`

---
**© 2026 ALCODIST_LABS_RND.** // *Performance-First Content Engineering.*
