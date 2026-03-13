# UK Housing Market & Rental Intelligence API

A comprehensive platform for analyzing and accessing UK housing market and rental data. This project features a high-performance Go backend with GraphQL and REST interfaces, a modern React frontend with interactive mapping, and a robust observability stack.

## 🏗 Project Overview

The system is designed to process and analyze large datasets of UK property transactions and rental listings. It provides insights into affordability, growth hotspots, and market distributions through interactive dashboards and geographic visualizations.

### Mission Statement
To provide a scalable, high-performance intelligence platform that transforms raw UK property data into actionable market insights for developers, analysts, and stakeholders.

### Data Source
The primary dataset used in this project is the **UK Property Price Data (1995-2023)**, sourced from Kaggle.
- **Source:** [Kaggle UK Property Price Data](https://www.kaggle.com/datasets/willianoliveiragibin/uk-property-price-data-1995-2023-04?resource=download)
- **Ingestion:** Data is ingested from `dataset.csv` into a PostgreSQL database, with asynchronous processing managed by a Redis-backed task queue.

### Core Value Propositions
- **High-Performance Processing:** Optimized for large-scale housing datasets using Go and PostgreSQL.
- **Geographic Visualizations:** Interactive mapping for regional trend analysis.
- **Hybrid API Architecture:** Combines the flexibility of GraphQL with the standardization of REST.
- **Layered Design:** Strictly followed architectural patterns (Router → Handler → Service → Repository → Database) for maximum maintainability and testability.

## 🛠 Tech Stack

### Backend
- **Language:** Go 1.25.0
- **Web Framework:** [Gin](https://gin-gonic.com/)
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL
- **GraphQL:** [gqlgen](https://gqlgen.com/)
- **Task Queue:** [Asynq](https://github.com/hibiken/asynq) (Redis-backed)
- **Storage:** [RustFS](https://github.com/rustfs/rustfs) (S3-compatible storage)
- **Documentation:** [Swagger/swag](https://github.com/swaggo/swag)

### Frontend
- **Framework:** [React 19](https://react.dev/)
- **Build Tool:** [Vite](https://vitejs.dev/)
- **State Management:** [Zustand](https://github.com/pmndrs/zustand) & [TanStack Query](https://tanstack.com/query/latest)
- **Routing:** [TanStack Router](https://tanstack.com/router/latest)
- **Styling:** [Tailwind CSS v4](https://tailwindcss.com/) & [shadcn/ui](https://ui.shadcn.com/)
- **Visualization:** [Leaflet](https://leafletjs.com/) (Maps) & [Recharts](https://recharts.org/) (Charts)
- **API Client:** generated via [Kubb](https://kubb.dev/)

### Infrastructure & Observability
- **Deployment:** [Pulumi](https://www.pulumi.com/) (IaC)
- **Containerization:** Docker & Docker Compose
- **Metrics:** Prometheus & Grafana
- **Logging:** Loki & Promtail

---

## 🚀 Getting Started (Development)

### Prerequisites
- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Go 1.25+](https://go.dev/dl/)
- [Node.js](https://nodejs.org/) & [pnpm](https://pnpm.io/)
- [Air](https://github.com/cosmtrek/air) (optional, for backend hot-reloading)

### 1. Start Infrastructure
Launch the database, redis, storage, and monitoring stack:
```bash
docker-compose -f infra/docker-compose.yml up -d
```

### 2. Backend Setup
```bash
cd backend
go mod download
# Start the server (or use 'air' for hot-reload)
go run cmd/main.go
```
The API will be available at `http://localhost:8080`.

### 3. Frontend Setup
```bash
cd frontend
pnpm install
pnpm dev
```
The UI will be available at `http://localhost:5173`.

### Default Credentials
- **App Login:** `admin` / `admin`
- **Grafana:** `admin` / `admin`
- **Postgres:** `postgres` / `password`
- **RustFS (S3):** `rustfsadmin` / `rustfsadmin`

---

## 🌐 API & Tools

- **REST API Docs:** `http://localhost:8080/swagger/index.html`
- **GraphQL Playground:** `http://localhost:8080/playground`
- **Asynq Monitoring:** `http://localhost:8090` (Queue status)
- **Grafana Dashboards:** `http://localhost:3000` (Metrics & Logs)
- **Prometheus:** `http://localhost:9090`

---

## 🚢 Production Deployment

The project is designed to be deployed using Pulumi for Infrastructure as Code and GitHub Actions for CI/CD.

### Infrastructure as Code (IaC)
Located in `infra/pulumi/`, the setup manages:
- Kubernetes resources (GKE/EKS/Self-hosted)
- Managed PostgreSQL and Redis instances
- Cloudflare configuration
- Observability stack deployment

### CI/CD Workflow
The `.github/workflows/deploy.yml` handles:
1. **Linting & Testing:** Go and React test suites.
2. **Build:** Docker image creation and pushing to registry.
3. **Deploy:** Automated Pulumi up to update the production environment.

### Production Environment Variables
In production, ensure the following are configured via Pulumi or secrets:
- `DATABASE_URL`: Production Postgres connection string.
- `REDIS_URL`: Production Redis instance.
- `BUCKET_ENDPOINT`: Production S3/Object Storage endpoint.
- `SECRET_KEY`: Long, random string for JWT signing.
- `IS_PROD`: Set to `true` to enable production optimizations.

---

## 📂 Project Structure

- `backend/`: Go source code, GraphQL schemas, and migrations.
- `frontend/`: React application, UI components, and API clients.
- `infra/`: Docker Compose for local dev and Pulumi for production.
- `scripts/`: Utility scripts for data migration and API generation.
- `docs/`: Additional documentation and architecture diagrams.
