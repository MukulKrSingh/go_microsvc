# 🍽️ Restaurant Microservices Architecture 

[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org/)
[![Gin](https://img.shields.io/badge/Gin-Web%20Framework-green.svg)](https://github.com/gin-gonic/gin)
[![Kafka](https://img.shields.io/badge/Kafka-Event%20Streaming-red.svg)](https://kafka.apache.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-blue.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Containerization-blue.svg)](https://www.docker.com/)

> 🚀 A modern, scalable restaurant ordering system built with microservices

This project demonstrates a microservices architecture for a restaurant ordering system. The system consists of two primary microservices that communicate asynchronously via Kafka, showcasing best practices for distributed systems design.

## 🏗️ Architecture Overview

<div align="center">
  <img src="docs/architecture-diagram.svg" alt="Microservices Architecture Diagram" width="700">
</div>

<details>
<summary>View ASCII Architecture Diagram</summary>

```
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

### 🍲 Restaurant Ordering Service
<table>
  <tr>
    <td><strong>Features</strong></td>
    <td>Authentication, food menu, ordering, transactions</td>
  </tr>
  <tr>
    <td><strong>Tech Stack</strong></td>
    <td>Go, Gin, PostgreSQL (SQL), Kafka Producer</td>
  </tr>
  <tr>
    <td><strong>API Port</strong></td>
    <td>8080</td>
  </tr>
  <tr>
    <td><strong>Database</strong></td>
    <td>PostgreSQL on port 5434 (mapped to container's 5432)</td>
  </tr>
</table>

### 🌟 User Feedback Service
<table>
  <tr>
    <td><strong>Features</strong></td>
    <td>Feedback collection and analysis for orders</td>
  </tr>
  <tr>
    <td><strong>Tech Stack</strong></td>
    <td>Go, Gin, PostgreSQL (GORM), Kafka Consumer</td>
  </tr>
  <tr>
    <td><strong>API Port</strong></td>
    <td>8081</td>
  </tr>
  <tr>
    <td><strong>Database</strong></td>
    <td>PostgreSQL on port 5433 (mapped to container's 5432)</td>
  </tr>
</table>

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

<div class="steps-container">
<div class="step">

#### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd go_microsvc
```
</div>

<div class="step">

#### Step 2: Start All Microservices

```bash
docker-compose up -d
```
</div>

<div class="step">

#### Step 3: Verify the Deployment

```bash
./verify_deployment.sh
```
</div>

<div class="step">

#### Step 4: Run Integration Tests

```bash
./test_microservices.sh
```
</div>
</div>

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

### 🍽️ Restaurant Ordering Service (port 8080)

<div class="api-container">

#### 🔑 Authentication
<div class="api-endpoint">
  <div class="method post">POST</div>
  <div class="path">/auth</div>
  <div class="description">Authenticate a user and get a JWT token</div>
</div>

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
<div class="api-endpoint">
  <div class="method get">GET</div>
  <div class="path">/food-items</div>
  <div class="description">Get list of available food items</div>
</div>

#### 👤 User Profile
<div class="api-endpoint">
  <div class="method get">GET</div>
  <div class="path">/profile</div>
  <div class="description">Get authenticated user's profile (Requires JWT)</div>
</div>

#### 📋 Orders
<div class="api-endpoint">
  <div class="method post">POST</div>
  <div class="path">/orders</div>
  <div class="description">Place a new order (Requires JWT)</div>
</div>

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
<div class="api-endpoint">
  <div class="method post">POST</div>
  <div class="path">/transactions</div>
  <div class="description">Complete a transaction for an order (Requires JWT)</div>
</div>

<details>
<summary>Example Request</summary>

```json
{
  "order_id": 1
}
```
</details>
</div>

<style>
.api-container {
  margin-bottom: 30px;
}
.api-endpoint {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border-radius: 5px;
  background-color: #f8f9fa;
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}
.method {
  font-weight: bold;
  padding: 5px 10px;
  border-radius: 4px;
  margin-right: 15px;
  min-width: 60px;
  text-align: center;
}
.get {
  background-color: #61affe;
  color: white;
}
.post {
  background-color: #49cc90;
  color: white;
}
.path {
  font-family: monospace;
  font-size: 16px;
  margin-right: 15px;
  color: #3b4151;
}
.description {
  color: #555;
}
details {
  margin-left: 20px;
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f0f0f0;
  border-radius: 5px;
}
details summary {
  cursor: pointer;
  color: #0077cc;
  font-weight: bold;
}
</style>

### 🌟 User Feedback Service (port 8081)

<div class="api-container">

#### 🔑 Authentication
<div class="api-endpoint">
  <div class="method post">POST</div>
  <div class="path">/auth</div>
  <div class="description">Authenticate a user and get a JWT token</div>
</div>

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

<div class="api-endpoint">
  <div class="method post">POST</div>
  <div class="path">/feedback</div>
  <div class="description">Submit feedback for an order (Requires JWT)</div>
</div>

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

<div class="api-endpoint">
  <div class="method get">GET</div>
  <div class="path">/feedback</div>
  <div class="description">Get all feedback from the authenticated user (Requires JWT)</div>
</div>

<div class="api-endpoint">
  <div class="method put">PUT</div>
  <div class="path">/feedback/:id</div>
  <div class="description">Update feedback (Requires JWT)</div>
</div>

<details>
<summary>Example Request</summary>

```json
{
  "rating": 4,
  "comment": "Good service but could be better"
}
```
</details>

<div class="api-endpoint">
  <div class="method delete">DELETE</div>
  <div class="path">/feedback/:id</div>
  <div class="description">Delete feedback (Requires JWT)</div>
</div>

#### 📊 Analytics

<div class="api-endpoint">
  <div class="method get">GET</div>
  <div class="path">/feedback/stats</div>
  <div class="description">Get feedback statistics (Requires JWT)</div>
</div>

</div>

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

<div class="testing-section">
  <div class="test-card">
    <div class="test-icon">🔍</div>
    <div class="test-content">
      <h4>Deployment Verification</h4>
      <p>The <code>verify_deployment.sh</code> script checks if all services and connections are working correctly.</p>
      <div class="test-command">
        <code>./verify_deployment.sh</code>
      </div>
    </div>
  </div>
  
  <div class="test-card">
    <div class="test-icon">🔄</div>
    <div class="test-content">
      <h4>Integration Tests</h4>
      <p>The <code>test_microservices.sh</code> script performs end-to-end integration tests across all services.</p>
      <div class="test-command">
        <code>./test_microservices.sh</code>
      </div>
    </div>
  </div>
</div>

## 📁 Project Structure

<div class="directory-structure">

```
go_microsvc/
├── 📄 docker-compose.yml            # Main container orchestration
├── 📄 test_microservices.sh         # Integration test script
├── 📄 verify_deployment.sh          # Deployment verification script
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

</div>

<style>
.testing-section {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin: 20px 0;
}
.test-card {
  flex: 1;
  min-width: 300px;
  display: flex;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 15px;
  background-color: #f9f9f9;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}
.test-icon {
  font-size: 2em;
  margin-right: 15px;
  display: flex;
  align-items: center;
}
.test-content {
  flex: 1;
}
.test-content h4 {
  margin-top: 0;
  margin-bottom: 10px;
  color: #333;
}
.test-command {
  background-color: #f0f0f0;
  padding: 8px 12px;
  border-radius: 4px;
  margin-top: 10px;
  display: inline-block;
}
.directory-structure {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 15px;
  box-shadow: inset 0 0 5px rgba(0,0,0,0.1);
  overflow-x: auto;
}
</style>

## 📝 Notes

<div class="notes">
  <ul>
    <li>🍔 All food items are initialized with 1000 units of quantity and a price of 10.</li>
    <li>👤 A default test user is created with username <code>testuser</code> and password <code>password123</code>.</li>
    <li>🔐 For simplicity, authentication uses plain text password comparison.</li>
    <li>🏛️ The project demonstrates key microservices principles:
      <div class="principles">
        <div class="principle">Service<br>Independence</div>
        <div class="principle">Decentralized<br>Data</div>
        <div class="principle">Async<br>Communication</div>
        <div class="principle">Tech<br>Diversity</div>
        <div class="principle">Independent<br>Deployment</div>
      </div>
    </li>
  </ul>
</div>

## 🔮 Future Improvements

<div class="future-improvements">
  <div class="improvement">
    <div class="improvement-icon">🔀</div>
    <div class="improvement-content">
      <h4>API Gateway</h4>
      <p>Add an API gateway for routing and cross-cutting concerns</p>
    </div>
  </div>
  
  <div class="improvement">
    <div class="improvement-icon">🔍</div>
    <div class="improvement-content">
      <h4>Service Discovery</h4>
      <p>Implement service registry for dynamic service discovery</p>
    </div>
  </div>
  
  <div class="improvement">
    <div class="improvement-icon">🛡️</div>
    <div class="improvement-content">
      <h4>Circuit Breakers</h4>
      <p>Add resilience patterns to handle service outages</p>
    </div>
  </div>
  
  <div class="improvement">
    <div class="improvement-icon">🔍</div>
    <div class="improvement-content">
      <h4>Distributed Tracing</h4>
      <p>Implement tracing to monitor request flows</p>
    </div>
  </div>
  
  <div class="improvement">
    <div class="improvement-icon">📊</div>
    <div class="improvement-content">
      <h4>Centralized Logging</h4>
      <p>Set up centralized logging across services</p>
    </div>
  </div>
</div>

<style>
.notes {
  background-color: #f8f9fa;
  border-left: 4px solid #4CAF50;
  padding: 15px;
  margin: 20px 0;
  border-radius: 0 5px 5px 0;
}
.notes ul {
  padding-left: 20px;
  margin: 0;
}
.notes li {
  margin-bottom: 10px;
}
.principles {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin: 10px 0 0 20px;
}
.principle {
  background-color: #e7f3ff;
  border: 1px solid #b3d8ff;
  border-radius: 5px;
  padding: 10px;
  text-align: center;
  font-size: 0.9em;
  flex: 1;
  min-width: 100px;
}
.future-improvements {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin: 20px 0;
}
.improvement {
  display: flex;
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 5px rgba(0,0,0,0.1);
}
.improvement-icon {
  font-size: 2em;
  margin-right: 15px;
  display: flex;
  align-items: center;
}
.improvement-content h4 {
  margin-top: 0;
  margin-bottom: 5px;
  color: #333;
}
.improvement-content p {
  margin: 0;
  color: #555;
}
footer {
  margin-top: 50px;
  padding-top: 20px;
  border-top: 1px solid #eee;
  text-align: center;
  color: #777;
}
</style>

<footer>
  <p>🍽️ Restaurant Microservices Architecture - Developed with ❤️ using Go and Kafka</p>
</footer>
