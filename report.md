# Technical Report: UK Housing Market & Rental Intelligence API

## Executive Summary
The UK Housing Market & Rental Intelligence API is a high-performance system designed to provide deep insights into property transactions and rental trends. This report outlines the architectural decisions, technology stack, and engineering practices employed to build a scalable and maintainable platform.

---

## 1. Technology Stack & Justification

### Backend: Go (Golang) 1.25
*   **Choice:** Go was selected for its exceptional performance, efficient concurrency model (Goroutines), and strong static typing.
*   **Justification:** Given the large-scale dataset (millions of records from UK housing history), Go's low memory footprint and fast execution are critical for real-time analytics and high-throughput API endpoints.

### Frameworks & Libraries
*   **Gin Gonic:** A lightweight, high-performance HTTP web framework. Chosen for its speed and robust middleware ecosystem.
*   **GORM:** An ORM that simplifies database interactions while providing the flexibility to write raw SQL for performance-critical queries (e.g., PostgreSQL estimated counts).
*   **gqlgen:** A schema-first GraphQL library. This ensures a single source of truth for the API contract, enabling type-safe communication between the backend and frontend.
*   **Asynq (Redis-based):** Used for background task processing. Crucial for handling long-running operations like data ingestion and analytics refreshes without blocking the main request-response cycle.

### Database: PostgreSQL
*   **Choice:** PostgreSQL with GORM.
*   **Justification:** Its support for advanced indexing, JSONB, and geospatial extensions (PostGIS) makes it ideal for property data which often involves complex geographic and historical queries.

### Frontend: React 19 & TypeScript
*   **TanStack Suite (Router, Query):** Provides robust state management and type-safe routing.
*   **Shadcn UI & Tailwind CSS:** Accelerates development with highly customizable, accessible components.
*   **Leaflet & Recharts:** Used for interactive geospatial visualization and data analytics.

---

## 2. Architectural Choices

### Layered Architecture
The project follows a strict layered pattern:
**Router → Handler → Service → Repository → Database**
*   **Separated Concerns:** Logic is decoupled; services handle business rules while repositories manage data persistence.
*   **Dependency Injection:** Interfaces are used for all dependencies, facilitating easy mocking and testing.

### Hybrid API Strategy (REST & GraphQL)
The system employs a dual-protocol approach to maximize flexibility and performance:
*   **REST:** Utilized for standard, predictable operations such as authentication, administrative user management, and file uploads. REST provides a simple, standard contract for these well-defined tasks.
*   **GraphQL:** Leveraged for complex property searching and multi-dimensional analytics. GraphQL allows the frontend to request specific fields and nested relationships in a single round-trip, significantly reducing over-fetching and improving performance on data-intensive views.

### Performance Optimizations
*   **Estimated Row Counts:** For large datasets, standard `SELECT COUNT(*)` is slow. The repository utilizes PostgreSQL internal statistics (`pg_class`) to provide instant estimates for large result sets.
*   **Batch Processing:** Data ingestion uses `CreateInBatches` to minimize database roundtrips.
*   **Caching Strategy:** Proactive cache clearing and scheduled analytics refreshes via `Asynq` ensure users see up-to-date data without overloading the database.

---

## 3. Challenges & Lessons Learned

### Challenges
*   **Dataset Scale:** Handling the `dataset.csv` required efficient ingestion scripts and optimized database indexes to keep search latency under 100ms.
*   **GraphQL Complexity:** Balancing the flexibility of GraphQL with the performance constraints of deep nested queries.

### Lessons Learned
*   **Schema-First is Better:** Defining the GraphQL schema first significantly reduced frontend-backend integration friction.
*   **Observability is Non-Negotiable:** Implementing structured logging (slog) and Prometheus metrics early was vital for debugging performance bottlenecks in production-like environments.

---

## 4. Testing Approach
*   **Unit Testing:** Comprehensive coverage of business logic using `testify` and `mock`.
*   **Middleware Testing:** Validation of security and logging layers using `httptest`.
*   **Integration Testing:** Schema-based testing of GraphQL resolvers to ensure API compliance.

---

## 5. Limitations & Future Development
*   **Geospatial Depth:** Currently using basic geographic filters; future versions will integrate PostGIS for radius-based and polygon searches.
*   **Search Optimization:** While PostgreSQL `ILIKE` is sufficient for now, a migration to Meilisearch or Elasticsearch is planned for advanced full-text search capabilities.
*   **Real-time Enhancements:** Expanding WebSocket usage for live transaction alerts and collaborative property analysis.

---

## 6. Generative AI Declaration
This project utilized Generative AI (Gemini) in the following capacities:
*   **Boilerplate Generation:** Creating repetitive code structures for models and repositories to accelerate development.
*   **Architectural Analysis:** Assisting in the design of the layered architecture and dependency injection patterns.
*   **Reference Implementation:** I provided example code from my previous projects (e.g., `example/stock-manage`, `example/GEMS`) to serve as a high-quality reference for patterns and conventions.
*   **Documentation:** Synthesizing complex technical documentation and creating this report.

**Analysis:** The use of AI allowed for a "strategic orchestration" approach, where the developer focuses on high-level architecture and security while the AI handles implementation details based on established patterns. For exact details on AI usage and specific reference patterns, please refer to **`GEMINI.md`**. This approach significantly reduced time-to-market while maintaining high code quality and consistency across the stack. All AI-generated code was rigorously reviewed and tested to ensure it met project standards.
