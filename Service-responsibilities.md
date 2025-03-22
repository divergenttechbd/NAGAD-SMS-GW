# Nagad SMS Gateway (SMSGW) Service Responsibilities

This document outlines the responsibilities of the six microservices in the Nagad SMS Gateway (SMSGW) project, based on the provided Statement of Work (SoW). The services are designed to meet requirements for scalability, high availability, security, and comprehensive SMS functionality, leveraging a microservices architecture with RabbitMQ for decoupling.

---

## 1. API Gateway Service

**Purpose**: Acts as the centralized entry point for all external requests (e.g., from DFS, Web Frontend, third-party systems), handling routing, security, and cross-cutting concerns to decouple clients from microservices.

**Responsibilities**:
- **API Exposure and Routing** (SoW 3.1.1, 3.10.4):
  - Exposes secure APIs (HTTPS, HTTP, JSON, SMPP) for external systems and Web Frontend.
  - Routes requests to appropriate services (e.g., `/core` to Core Service, `/reporting` to Reporting Service).
- **Security Enforcement** (SoW 3.10.2, 3.10.4):
  - Enforces authentication (OAuth 2.0, API keys), HTTPS/TLS, IP whitelisting, and payload validation to prevent injection attacks.
  - Centralizes rate limiting and request throttling.
- **Request Logging** (SoW 3.11.2):
  - Logs incoming requests (method, path, timestamp) to RabbitMQ for Logging Service consumption.
- **Load Balancing and Redundancy** (SoW 3.13.1):
  - Distributes traffic across microservice instances for scalability and fault tolerance.
- **Health Check**:
  - Provides a `/health` endpoint for monitoring gateway status.
- **Cross-Cutting Concerns**:
  - Handles retries for failed downstream requests (SoW 3.1.8) and basic error handling.

**Notes**: Offloads API-related tasks from Core Service, ensuring a single, secure entry point. Scales independently with multiple instances behind a load balancer.

---

## 2. Core Service

**Purpose**: Manages SMS ingestion, queueing, and orchestration, serving as the backend hub for business logic and internal coordination.

**Responsibilities**:
- **Message Ingestion and Validation** (SoW 3.1, 3.3, 3.6):
  - Receives SMS requests from API Gateway via internal APIs.
  - Validates payloads, assigns unique SMS IDs, and sub-application IDs (3.7.1, 3.7.2).
  - Implements prefix-based MNO selection if MNO is missing (3.6).
- **Priority Queue Management** (SoW 3.5):
  - Manages RabbitMQ priority queues (e.g., OTP: 0, Transactions: 1, Bulk: 2, Marketing: 3) (3.5.1-3.5.4).
  - Enqueues messages with priority levels from DFS or Web Frontend (3.5.2).
- **Store-and-Dispatch Logic** (SoW 3.4):
  - Stores priority SMS (e.g., type 50) and dispatches them based on user triggers (3.4.1-3.4.3).
- **Channel Routing Preparation** (SoW 3.3):
  - Determines MNO and delivery channel based on priority and configuration (3.3.1, 3.5.3).
- **Internal API Exposure**:
  - Provides APIs for Web Frontend (e.g., configurations, submissions) and Consumer Service (e.g., queue metadata).
- **Security Basics** (SoW 3.10):
  - Encrypts sensitive data in transit and applies masking (e.g., OTPs) (3.10.1).

**Notes**: Focuses on business logic, delegating client-facing APIs to API Gateway and delivery to Consumer Service. Uses a shared database with Consumer for consistency.

---

## 3. Logging Service

**Purpose**: Centralizes logging, auditing, and monitoring to ensure traceability, compliance, and system health tracking.

**Responsibilities**:
- **Comprehensive Logging** (SoW 3.10.3, 3.11.2):
  - Logs all SMS requests (source system, module/API, Unique Request ID, Transaction ID, Reference ID) (3.16.2).
  - Records receipt status with timestamps (3.11.2).
- **Audit Trails** (SoW 3.10.3, 3.11.4):
  - Maintains audit logs for configuration changes, user activities, and system events (3.10.9).
- **System Monitoring** (SoW 3.15):
  - Tracks performance (API responses, database, queues) and VM resources (CPU, memory, disk) (3.15.1-3.15.3, 3.15.8).
  - Monitors MNO/channel failures and submission rates (3.15.4-3.15.7).
- **Alert Notifications** (SoW 3.15):
  - Sends alerts for failures (e.g., API downtime, channel issues) to L1 teams.
- **Data Retention and Purging** (SoW 3.10.12, 3.15.9):
  - Manages log retention and purging per Bangladesh Bank regulations.
- **Troubleshooting Support** (SoW 3.11.4):
  - Provides detailed logs for debugging across services.

**Notes**: Consumes logs from API Gateway and other services via RabbitMQ. Uses a dedicated time-series DB (e.g., Elasticsearch) for high write throughput.

---

## 4. Reporting Service

**Purpose**: Handles reporting, analytics, and data aggregation, serving data to Web Frontend for display and external systems for integration.

**Responsibilities**:
- **MIS Reporting** (SoW 3.9):
  - Generates detailed and summary reports (MNO-wise, user-wise, channel-wise, monthly/weekly) (3.9.1).
  - Offers customizable report filters.
- **Analytical Dashboards** (SoW 3.16.6, 3.16.8):
  - Provides data for trends (e.g., delayed OTPs, non-DFS users).
- **Real-Time Data Feeds** (SoW 3.16.9):
  - Produces near real-time feeds for external systems (e.g., every 5 minutes or 10MB).
- **Health Check Reports** (SoW 4.2):
  - Generates monthly health check reports for application and database.
- **Data Supply for GUI** (SoW 3.8.2, 3.8.3):
  - Supplies real-time message status and failure data to Web Frontend via APIs.

**Notes**: Uses a dedicated analytical DB (e.g., ClickHouse) synced from Core/Consumer via CDC or ETL. Serves Web Frontend and external systems.

---

## 5. Consumer Service

**Purpose**: Consumes queued messages and manages SMS delivery to MNOs, including channel logic and retries.

**Responsibilities**:
- **SMS Delivery** (SoW 3.1, 3.3):
  - Sends messages to MNOs via HTTP, HTTPS, JSON, SMPP (3.1.1, 3.3.1).
  - Handles bulk SMS (personalized/predefined) and scheduled SMS (3.1.6, 3.1.7).
  - Supports push, OTP, and Unicode SMS (3.1.5, 3.1.3).
- **Channel Management** (SoW 3.3):
  - Manages multiple channels per MNO (e.g., HTTP and SMPP) with failover (3.3.1, 3.3.3).
  - Segregates TPS by SMS type/channel (3.3.2).
- **Retry Mechanism** (SoW 3.1.8):
  - Executes retries for failed deliveries.
- **MNP Checking** (SoW 3.1.10):
  - Supports MNP checking (future requirement).
- **Push-Pull SMS** (SoW 3.2):
  - Processes incoming requests and sends template responses (3.2.1, 3.2.2).
- **Queue Consumption** (SoW 3.5):
  - Dequeues messages from RabbitMQ and routes to MNOs (3.5.4).

**Notes**: Scales independently for delivery load. Shares a database with Core for message lifecycle consistency.

---

## 6. Web Frontend Service

**Purpose**: Provides the user interface for managing the SMS Gateway, including real-time status, configurations, and reporting.

**Responsibilities**:
- **Web GUI Interface** (SoW 3.8):
  - Implements a web portal with role-based access (Admin, Super Admin, User) (3.8.1).
  - Supports maker-checker workflows for configurations (3.8.1.b).
  - Displays real-time message status and failure notifications (3.8.2, 3.8.3).
  - Offers SMS submission and delivery status checks for individual numbers (3.8.4).
- **Configuration Management** (SoW 3.3.2, 3.8):
  - Manages MNO channels, TPS segregation, and priority settings via GUI.
  - Submits changes to Core Service via API Gateway.
- **Reporting Integration** (SoW 3.8, 3.9):
  - Displays MIS reports and dashboards by querying Reporting Service via API Gateway.
- **User Interaction** (SoW 3.1.6, 3.1.7):
  - Allows bulk SMS uploads (XLS, CSV) and scheduling, forwarding to Core via API Gateway.
- **Security** (SoW 3.10.2):
  - Enforces RBAC, MFA, and IP whitelisting for GUI access.

**Notes**: Stateless UI (e.g., Angular/React), uses API Gateway for all backend interactions. Logs user actions to Logging Service.

---

## Architecture Overview

- **API Gateway Service**: External entry point, routes requests, enforces security.
- **Core Service**: Ingestion and queueing (RabbitMQ producer).
- **Consumer Service**: Delivery to MNOs (RabbitMQ consumer).
- **Logging Service**: Centralized logging and monitoring.
- **Reporting Service**: Analytics and reporting.
- **Web Frontend Service**: User interface.

**Data Flow**:
1. DFS/Web Frontend → API Gateway → Core → RabbitMQ.
2. Consumer ← RabbitMQ → MNOs.
3. Web Frontend ← API Gateway ← Core/Reporting.
4. All services → Logging (via RabbitMQ).

**Database Strategy**:
- Core + Consumer: Shared relational DB (e.g., PostgreSQL).
- Logging: Dedicated time-series DB (e.g., Elasticsearch).
- Reporting: Dedicated analytical DB (e.g., ClickHouse).
- API Gateway/Web Frontend: No DB (cache optional, e.g., Redis).

---

## Conclusion

These six services comprehensively cover the SoW requirements, ensuring scalability, security, and functionality for the SMSGW. The API Gateway enhances the microservices architecture by centralizing client access, while each service maintains a clear, focused responsibility.