# 🏦 Go Dynatrace Banking Demo

Go Dynatrace Banking Demo is a sample banking application built with **Go (Gin)** and **PostgreSQL** to demonstrate end-to-end observability using **Dynatrace** and **OpenTelemetry**.

The application simulates common banking operations such as user management, account inquiry, money transfer, and payment. It is designed as a hands-on project for showcasing distributed tracing, structured logging, and business transaction monitoring.

---

# Architecture

                +----------------+
                | Client         |
                | Postman / API  |
                +--------+-------+
                         |
                         ▼
                +----------------+
                | Kong Gateway   |
                +--------+-------+
                         |
                         ▼
                +----------------+
                | Banking API    |
                | (Go + Gin)     |
                +--------+-------+
                         |
          +--------------+--------------+
          |              |              |
          ▼              ▼              ▼
      User Service  Transfer Service  Payment Service
                         |
                         ▼
                  Repository Layer
                         |
                         ▼
                    PostgreSQL

               OpenTelemetry + Dynatrace
```

---

# Request Flow

Client
   │
   ▼
Kong Gateway
   │
   ▼
Banking API
   │
   ▼
Business Service
   │
   ▼
Repository
   │
   ▼
PostgreSQL
   │
   ▼
Response

          │
          ▼
OpenTelemetry Trace
          │
          ▼
Dynatrace
```

---

# Components

| Component | Description |
|----------|-------------|
| Banking API | REST API built with Go and Gin |
| Kong Gateway | API Gateway for routing incoming requests |
| PostgreSQL | Stores customer, account, and transaction data |
| OpenTelemetry | Generates traces and business attributes |
| Zap Logger | Structured application logging |
| Docker Compose | Local deployment environment |
| Dynatrace | Observability platform for monitoring services, traces, logs, and infrastructure |

---

# Business Features

- User Management
- Login
- Account Inquiry
- Money Transfer
- Payment
- Transaction History
- Health & Readiness Endpoints

---

# Run Locally

```bash
docker compose up --build
```

Application

```
http://localhost:8080
```

Kong Gateway

```
http://localhost:8000
```