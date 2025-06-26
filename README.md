# 🍽️ Restaurant Microservices Architecture 

[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-Web%20Framework-green.svg)](https://github.com/gin-gonic/gin)
[![Kafka](https://img.shields.io/badge/Kafka-Event%20Streaming-red.svg)](https://kafka.apache.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Containerization-blue.svg)](https://www.docker.com/)
[![Traefik](https://img.shields.io/badge/Traefik-API%20Gateway-orange.svg)](https://traefik.io/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Orchestration-326CE5.svg)](https://kubernetes.io/)

> 🚀 A modern, scalable restaurant ordering system built with microservices

This project demonstrates a microservices architecture for a restaurant ordering system. The system consists of two primary microservices that communicate asynchronously via Kafka, showcasing best practices for distributed systems design.

## 🏗️ Architecture Overview

<div align="center">
  <img src="docs/architecture-with-gateway.svg" alt="Microservices Architecture Diagram with API Gateway" width="700">
</div>

### Deployment Options

This project can be deployed using either:

- **Docker Compose**: For local development and testing
- **Kubernetes**: For production and scalable deployments (see the [kubernetes](./kubernetes) directory)

<details>
<summary>View ASCII Architecture Diagram</summary>

```
                 ┌─────────────────────┐
                 │                     │
                 │   API Gateway       │
                 │   (Traefik)         │
                 │                     │
                 └───────┬─────────────┘
                         │
          ┌──────────────┴───────────────┐
          │                              │
          ▼                              ▼
┌─────────────────────┐     ┌──────────────────┐
│                     │     │                  │
│  Restaurant Order   │     │  User Feedback   │
│     Service         │◄────┤     Service      │
│                     │     │                  │
└─────────┬───────────┘     └──────────┬───────┘
          │                            │
          │                            │
          ▼                            ▼
┌─────────────────────┐     ┌──────────────────┐
│                     │     │                  │
│  Restaurant DB      │     │   Feedback DB    │
│   (PostgreSQL)      │     │   (PostgreSQL)   │
│                     │     │                  │
└─────────────────────┘     └──────────────────┘
          ▲                            ▲
          │                            │
          └────────────┬───────────────┘
                       │
                       ▼
              ┌─────────────────┐
              │                 │
              │      Kafka      │
              │                 │
              └─────────────────┘
```
</details>

## 🔍 Services Description

### 🔀 API Gateway

| Feature | Description |
|---------|-------------|
| **Features** | Routing, load balancing, API endpoints consolidation |
| **Tech Stack** | Traefik |
| **Dashboard Port** | 8082 |
| **HTTP Port** | 80 |
| **Gateway Paths** | `/api/restaurant/*`, `/api/feedback/*` |

### 🍲 Restaurant Ordering Service

| Feature | Description |
|---------|-------------|
| **Features** | Authentication, food menu, ordering, transactions |
| **Tech Stack** | Go, Gin, PostgreSQL (SQL), Kafka Producer |
| **Direct API Port** | 8080 |
| **Gateway Path** | `/api/restaurant` |
| **Database** | PostgreSQL on port 5434 (mapped to container's 5432) |

### 🌟 User Feedback Service

| Feature | Description |
|---------|-------------|
| **Features** | Feedback collection and analysis for orders |
| **Tech Stack** | Go, Gin, PostgreSQL (GORM), Kafka Consumer |
| **Direct API Port** | 8081 |
| **Gateway Path** | `/api/feedback` |
| **Database** | PostgreSQL on port 5433 (mapped to container's 5432) |

## 🌊 Data Flow

<div align="center">
  <img src="https://mermaid.ink/img/pako:eNptkLFuwzAMRH-F4JRO-QED3bpk6tahuLBI9AXWkqKRLh00Dv7ekpE2cUTNh3ePR97lIQDW4IMNMbjInVn0XmYhoAblzJAN2qZp1eDPM-QYkDIcweoZHLEVqbinNKL2ANdW3YJrYFICjTlb59FzDsp-gR5EM40VT4a6_eaRw1Rtsu4UHptOyo7xj9UKc_wkwaXn5McLe0C2_1NxKMtxzwb27JBbuU9RW_ziVtKFovxg-FcXOPGId6QyI7zE5uN66VCcRaXNdQLV3aGB9TYZ0Fon69d5YvVBqe8fvAZiKA" alt="Data Flow" style="max-width:600px">
</div>

1. 👤 A user places an order through the Restaurant Ordering Service
2. 💳 Once the transaction is completed, the order event is published to Kafka
3. 📡 The User Feedback Service consumes the order event from Kafka
4. 📝 The user can now submit feedback for the order through the Feedback Service
5. 📊 All feedback is stored and analyzed by the Feedback Service

## 📋 Prerequisites

- 🐳 Docker and Docker Compose
- <img src="https://golang.org/favicon.ico" width="16" height="16"> Go 1.22 or later (for local development)
- <img src="https://git-scm.com/favicon.ico" width="16" height="16"> Git

## 🚀 Getting Started

### Setup and Deployment

#### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd go_microsvc
```

#### Step 2: Start All Microservices

```bash
docker-compose up -d
```

#### Step 3: Verify the Deployment

```bash
./verify_deployment.sh
```

#### Step 4: Run Integration Tests

```bash
./test_microservices.sh
```

#### Step 5: Access the API Gateway

Once all services are up and running, you can access them through the API Gateway:

- **Traefik Dashboard**: [http://localhost:8082/dashboard/](http://localhost:8082/dashboard/)
- **Restaurant Service API**: [http://localhost/api/restaurant/food-items](http://localhost/api/restaurant/food-items)
- **Feedback Service API**: [http://localhost/api/feedback/health](http://localhost/api/feedback/health)

The API Gateway provides:
- **Centralized Routing**: All requests go through a single entry point
- **Path-Based Routing**: Routes requests based on URL path prefixes
- **Rate Limiting**: Prevents abuse with built-in rate limiting
- **CORS Headers**: Automatically adds needed CORS headers
- **Dashboard**: Visual interface to monitor traffic and routes

### Kubernetes Deployment

To deploy to Kubernetes instead of Docker Compose:

```bash
# Deploy with default local registry
./deploy_kubernetes.sh

# Or specify a custom registry
./deploy_kubernetes.sh my-registry.io
```

For more details on Kubernetes deployment, see the [kubernetes](./kubernetes) directory.

### 💻 Development Environment

<details>
<summary>Running Individual Services for Development</summary>

```bash
# For Restaurant Ordering Service:
cd restaurant_ordering_service
docker-compose up -d

# For User Feedback Service:
cd user_feedback_service
docker-compose up -d
```
</details>

<style>
.steps-container {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin: 20px 0;
}
.step {
  flex: 1;
  min-width: 300px;
  padding: 15px;
  border: 1px solid #ddd;
  border-radius: 5px;
  background-color: #f9f9f9;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}
.step h4 {
  margin-top: 0;
  border-bottom: 2px solid #4CAF50;
  padding-bottom: 5px;
}
</style>

## 📘 API Documentation

### 🔀 API Gateway Endpoints

| Service | Direct Access | Gateway Access |
|---------|---------------|---------------|
| Restaurant Service | `http://localhost:8080` | `http://localhost/api/restaurant` |
| Feedback Service | `http://localhost:8081` | `http://localhost/api/feedback` |
| Traefik Dashboard | N/A | `http://localhost:8082` |

### 🍽️ Restaurant Ordering Service

#### 🔑 Authentication

**POST /api/restaurant/auth** - Authenticate a user and get a JWT token
**POST /auth** - Direct access endpoint

<details>
<summary>Example Request</summary>

```json
{
  "username": "testuser",
  "password": "password123"
}
```
</details>

#### 🍔 Food Items

**GET /api/restaurant/food-items** - Get list of available food items (via Gateway)
**GET /food-items** - Direct access endpoint

#### 👤 User Profile

**GET /api/restaurant/profile** - Get authenticated user's profile (Requires JWT, via Gateway)
**GET /profile** - Direct access endpoint (Requires JWT)

#### 📋 Orders

**POST /api/restaurant/orders** - Place a new order (Requires JWT, via Gateway)
**POST /orders** - Direct access endpoint (Requires JWT)

<details>
<summary>Example Request</summary>

```json
{
  "items": [
    {"food_item_id": 1, "quantity": 2},
    {"food_item_id": 3, "quantity": 1}
  ]
}
```
</details>

#### 💳 Transactions

**POST /api/restaurant/transactions** - Complete a transaction for an order (Requires JWT, via Gateway)
**POST /transactions** - Direct access endpoint (Requires JWT)

<details>
<summary>Example Request</summary>

```json
{
  "order_id": 1
}
```
</details>
</div>


### 🌟 User Feedback Service

#### 🔑 Authentication

**POST /api/feedback/auth** - Authenticate a user and get a JWT token (via Gateway)
**POST /auth** - Direct access endpoint

<details>
<summary>Example Request</summary>

```json
{
  "username": "testuser",
  "password": "password123"
}
```
</details>

#### 📝 Feedback Management

**POST /api/feedback/feedback** - Submit feedback for an order (Requires JWT, via Gateway)
**POST /feedback** - Direct access endpoint (Requires JWT)

<details>
<summary>Example Request</summary>

```json
{
  "order_id": 1,
  "rating": 5,
  "comment": "Great service!"
}
```
</details>

**GET /api/feedback/feedback** - Get all feedback from the authenticated user (Requires JWT, via Gateway)
**GET /feedback** - Direct access endpoint (Requires JWT)

**PUT /api/feedback/feedback/:id** - Update feedback (Requires JWT, via Gateway)
**PUT /feedback/:id** - Direct access endpoint (Requires JWT)

<details>
<summary>Example Request</summary>

```json
{
  "rating": 4,
  "comment": "Good service but could be better"
}
```
</details>

**DELETE /api/feedback/feedback/:id** - Delete feedback (Requires JWT, via Gateway)
**DELETE /feedback/:id** - Direct access endpoint (Requires JWT)

#### 📊 Analytics

**GET /api/feedback/feedback/stats** - Get feedback statistics (Requires JWT, via Gateway)
**GET /feedback/stats** - Direct access endpoint (Requires JWT)

<style>
.put {
  background-color: #fca130;
  color: white;
}
.delete {
  background-color: #f93e3e;
  color: white;
}
</style>

## 🧪 Testing

### 🔍 Deployment Verification

The `verify_deployment.sh` script checks if all services and connections are working correctly, including the API Gateway.

```bash
./verify_deployment.sh
```

### 🔀 Testing the API Gateway

You can use the provided script to test all API Gateway endpoints:

```bash
./test_api_gateway.sh
```

Or test individual endpoints manually:

```bash
# Check if the API Gateway is accessible
curl http://localhost/

# Get food items through the API Gateway
curl http://localhost/api/restaurant/food-items

# Login through the API Gateway (returns a JWT token)
curl -X POST -H "Content-Type: application/json" \
  -d '{"username":"testuser", "password":"password123"}' \
  http://localhost/api/restaurant/auth
```

### 🔄 Integration Tests

The `test_microservices.sh` script performs end-to-end integration tests across all services.

```bash
./test_microservices.sh
```

## 📁 Project Structure

```
go_microsvc/
├── 📄 docker-compose.yml            # Main container orchestration
├── 📄 test_microservices.sh         # Integration test script
├── 📄 test_api_gateway.sh           # API Gateway test script
├── 📄 verify_deployment.sh          # Deployment verification script
├── 📁 traefik/                      # API Gateway configuration
├── 📁 kubernetes/                   # Kubernetes deployment manifests
│   ├── 📄 traefik.yml               # Main Traefik configuration
│   ├── 📁 dynamic/                  # Dynamic config directory
│   │   └── 📄 conf.yml              # Routes, middlewares, services
│   └── 📄 Dockerfile                # Traefik container definition
├── 📁 restaurant_ordering_service/  # Restaurant ordering service
│   ├── 📁 cmd/                      # Service entry point
│   ├── 📁 internal/                 # Service implementation
│   │   ├── 📁 api/                  # API handlers
│   │   ├── 📁 db/                   # Database operations
│   │   ├── 📁 kafka/                # Kafka producer
│   │   ├── 📁 middleware/           # Service middleware
│   │   └── 📁 models/               # Data models
│   ├── 📄 .env                      # Environment variables
│   ├── 📄 Dockerfile                # Container definition
│   └── 📄 go.mod                    # Go module file
└── 📁 user_feedback_service/        # User feedback service
    ├── 📁 cmd/                      # Service entry point
    ├── 📁 internal/                 # Service implementation
    │   ├── 📁 api/                  # API handlers
    │   ├── 📁 db/                   # GORM database operations
    │   ├── 📁 kafka/                # Kafka consumer
    │   ├── 📁 middleware/           # Service middleware
    │   └── 📁 models/               # Data models
    ├── 📄 .env                      # Environment variables
    ├── 📄 Dockerfile                # Container definition
    └── 📄 go.mod                    # Go module file
```



## 📝 Notes

- 🍔 All food items are initialized with 1000 units of quantity and a price of 10.
- 👤 A default test user is created with username `testuser` and password `password123`.
- 🔐 For simplicity, authentication uses plain text password comparison.
- 🏛️ The project demonstrates key microservices principles:
  - **Service Independence**: Each service has its own codebase and database
  - **Decentralized Data**: Each service manages its own data
  - **Async Communication**: Services communicate via events through Kafka
  - **Tech Diversity**: Different database approaches (SQL vs GORM)
  - **Independent Deployment**: Each service can be deployed separately

## 🔮 Future Improvements

1.  **Service Discovery**: Implement service registry for dynamic service discovery
2. 🛡️ **Circuit Breakers**: Add resilience patterns to handle service outages
3. 🔍 **Distributed Tracing**: Implement tracing to monitor request flows
4. 📊 **Centralized Logging**: Set up centralized logging across services
5. 🔒 **Enhanced API Gateway Security**: Add rate limiting, JWT validation, and other security features to the API Gateway
6. ☁️ **Kubernetes Autoscaling**: Configure Horizontal Pod Autoscaler for automatic scaling in Kubernetes

---

🍽️ Restaurant Microservices Architecture - Developed with ❤️ using Go and Kafka
