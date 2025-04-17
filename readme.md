# E-Commerce Microservices Application

This is a scalable e-commerce application built with **Go** using a microservices architecture. It leverages modern technologies like **gRPC**, **Kafka**, **GraphQL**, **Redis**, **MongoDB**, and **MySQL** to deliver a performant and reliable system.

## Architecture Overview

The application is divided into several microservices, each handling specific functionality. Below is the architecture diagram:

```
                   +------------------+
                   |     CLIENT       |
                   | (Web / Mobile)   |
                   +--------+---------+
                            |
                            v
                +-----------+-------------+
                |       API GATEWAY       |
                | - GraphQL (Queries)     |
                | - REST (Auth)           |
                | - Caching (Redis)       |
                | - Rate Limiting         |
                +-----------+-------------+
                            |
        +-------------------+----------------------------+
        |                   |                            |
        v                   v                            v
+---------------+   +---------------+        +---------------+
|   USER        |   |   PRODUCT     |<------>|    ORDER      |
| Microservice  |   | Microservice  |        | Microservice  |
+------+--------+   +--------+------+        +------+--------+
       |                     |                      |
       |                     |                      |
       +---------------------+----------------------+
                             |
                             v
                   +----------------------+
                   |     APACHE KAFKA     |
                   |  (Message Broker)    |
                   +----------+-----------+
                              |
                              v
                   +----------------------+
                   |  NOTIFICATION SERVICE|
                   | - Sends Email/SMS    |
                   +----------------------+
```

## Key Components

### API Gateway
- **Purpose**: Acts as the entry point for client requests (web/mobile).
- **Features**:
  - **GraphQL**: Handles queries to reduce under-fetching/over-fetching.
  - **REST**: Manages authentication/authorization via **Auth0** (Google Auth with JWT).
  - **Redis**: Caches responses to reduce database load and latency.
  - **Rate Limiting**: Prevents abuse using Redis-based rate limiting.
- **Run**: `go run server.go`

### User Service
- **Purpose**: Manages user profiles and information.
- **Features**:
  - Communicates with other services using **gRPC** for low-latency, synchronous calls.
  - Subscribes to Kafka's `order-created` event to fetch user data and publishes events to the **Notification Service**.
  - Publishes a `user-created` event to Kafka when a new user is created.
- **Run**: `go run cmd/main.go`

### Order Service
- **Purpose**: Handles order-related operations (e.g., create, update status).
- **Features**:
  - Uses **gRPC** to synchronously check product stock with the Product Service.
  - Publishes events to Kafka (e.g., `order-created`).
- **Run**: `go run cmd/main.go`

### Product Service
- **Purpose**: Manages product-related requests (e.g., product details, stock).
- **Features**:
  - Communicates with the Order Service via **gRPC** for stock checks.
- **Run**: `go run cmd/main.go`

### Notification Service
- **Purpose**: Sends email/SMS notifications.
- **Features**:
  - Subscribes to Kafka events (e.g., `user-created`, `order-created`) to trigger notifications.

### Supporting Technologies
- **Apache Kafka**: Fault-tolerant message broker for asynchronous communication between services.
- **Redis**: Used for caching and rate limiting to improve performance.
- **MongoDB/MySQL**: Databases for storing user, product, and order data.

## Prerequisites
1. **Go**: Install Go (version 1.21 or later).
2. **Kafka**: Ensure Kafka is running.
3. **Redis**: Ensure Redis server is running.
4. **MongoDB/MySQL**: Set up databases and configure connection strings.
5. **Auth0**: Configure Google authentication for the API Gateway.

## Setup Instructions
1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Start Kafka and Redis servers:
   ```bash
   # Start Kafka
   kafka-server-start.sh config/server.properties

   # Start Redis
   redis-server
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run each service:
   - API Gateway: `go run server.go`
   - User Service: `cd user-service && go run cmd/main.go`
   - Order Service: `cd order-service && go run cmd/main.go`
   - Product Service: `cd product-service && go run cmd/main.go`

5. Ensure the Notification Service is configured to listen to Kafka events.

## Why This Stack?
- **Go**: Fast, lightweight, and ideal for microservices.
- **gRPC**: Efficient, low-latency communication compared to REST.
- **Kafka**: Scalable and fault-tolerant for event-driven architecture.
- **GraphQL**: Flexible querying to optimize client-server interactions.
- **Redis**: Reduces latency through caching and rate limiting.
- **MongoDB/MySQL**: Reliable storage for structured and unstructured data.

## Contributing
Feel free to submit issues or pull requests to improve the project. Ensure you follow the coding standards and include tests for new features.

