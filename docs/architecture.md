# System Architecture

## Overview

This project implements a production-style three-tier application architecture deployed on AWS EC2 using Docker Compose and automated through GitHub Actions CI/CD.

The system consists of:

- Presentation Layer (Nginx)
- Application Layer (Go Backend API)
- Data Layer (PostgreSQL/Supabase Database)
- CI/CD Automation Layer (GitHub Actions)
- Container Registry (Docker Hub)
- Deployment Infrastructure (AWS EC2)

---

# High-Level Architecture

```text
                    ┌─────────────────┐
                    │   Developer     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ GitHub Repo     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ GitHub Actions  │
                    │      CI         │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Docker Hub      │
                    │ Container Repo  │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ GitHub Actions  │
                    │      CD         │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ AWS EC2         │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Docker Compose  │
                    └────────┬────────┘
                             │
        ┌────────────────────┼────────────────────┐
        ▼                    ▼                    ▼
 ┌────────────┐      ┌────────────┐      ┌────────────┐
 │   Nginx    │─────▶│ Go Backend │─────▶│ PostgreSQL/│
 │ Port: 80   │      │ Port:8080  │      │ Supabase   │
 └────────────┘      └────────────┘      └────────────┘
```

---

# Application Architecture

The application follows a traditional three-tier architecture.

## Presentation Layer

### Nginx

Responsibilities:

- Serves frontend application
- Acts as reverse proxy
- Receives incoming HTTP requests
- Routes API traffic to backend services

Container:

```yaml
nginx:alpine
```

Port:

```text
80
```

Publicly Accessible:

```text
http://<EC2_PUBLIC_IP>
```

---

## Application Layer

### Go Backend API

Responsibilities:

- Business Logic
- API Endpoints
- Database Communication
- Request Processing

Container:

```yaml
ritvikkantip/k8app:latest
```

Port:

```text
8080
```

The backend is only accessible internally through Docker networking.

---

## Data Layer

### PostgreSQL / Supabase

Responsibilities:

- Data Persistence
- Application State Storage
- User Settings Storage (Target Role & Path settings)

Container:

```yaml
postgres:16-alpine
```

Port:

```text
5433 (Host mapped), 5432 (Internal Docker Network)
```

Persistent Storage:

```yaml
postgres_data
```

The database remains persistent even after container recreation and supports direct connections to Supabase for production.

---

# Docker Networking

Docker Compose automatically provisions an internal network.

Services communicate using Docker DNS service discovery.

Example:

```text
backend
   │
   ▼
db
```

Instead of:

```text
192.168.x.x
```

the backend simply connects to:

```text
db
```

This improves portability and simplifies service discovery.

---

# Deployment Architecture

The deployment process is fully automated using GitHub Actions.

---

## Continuous Integration (CI)

The CI workflow is triggered whenever code is pushed to the repository.

### CI Responsibilities

- Checkout Source Code
- Configure Docker Buildx
- Authenticate with Docker Hub
- Build Backend Docker Image
- Push Docker Image to Docker Hub

Pipeline:

```text
Push
  │
  ▼
Checkout
  │
  ▼
Docker Build
  │
  ▼
Docker Hub Push
```

Generated Image:

```text
ritvikkantip/k8app:latest
```

---

## Continuous Deployment (CD)

The CD workflow executes automatically after a successful CI run.

### CD Responsibilities

- Checkout Repository
- Transfer Deployment Files
- SSH Into EC2
- Inject Runtime Secrets
- Pull Latest Images
- Deploy Updated Containers

Pipeline:

```text
Successful CI
       │
       ▼
Trigger CD
       │
       ▼
SCP Files To EC2
       │
       ▼
SSH Into Server
       │
       ▼
Inject Secrets
       │
       ▼
Docker Compose Pull
       │
       ▼
Docker Compose Up
```

---

# Secret Management

Sensitive values are managed using GitHub Actions Secrets.

Configured Secrets:

| Secret | Purpose |
|----------|----------|
| DOCKERHUB_USERNAME | Docker Hub Username |
| DOCKERHUB_PASS | Docker Hub Access Token |
| EC2_HOST | AWS EC2 Public IP |
| EC2_USER | EC2 SSH User |
| EC2_SSH_KEY | EC2 SSH Key |
| DATABASE_URL | Supabase/PostgreSQL connection string (optional) |
| DB_NAME | Database Name |
| DB_USER | Database User |
| DB_PASSWORD | Database Password |

No sensitive credentials are committed to source control.

---

# Data Persistence Strategy

PostgreSQL data is stored using Docker Volumes.

```yaml
volumes:
  postgres_data:
```

Benefits:

- Data survives container restarts
- Data survives image updates
- Database state remains persistent

---

# Health Checks

The PostgreSQL container implements health checks.

```yaml
healthcheck:
  test: ["CMD-SHELL", "pg_isready -U ${DB_USER} -d ${DB_NAME}"]
```

Purpose:

- Verify database readiness
- Prevent dependent services from starting too early
- Improve deployment reliability

---

# Security Considerations

Current Security Measures:

- GitHub Secrets for credential storage
- SSH Key-based authentication
- Docker network isolation
- Environment-based configuration
- No hardcoded credentials in source code

Potential Improvements:

- AWS Secrets Manager
- HTTPS via Let's Encrypt
- Reverse Proxy Hardening
- Container Image Signing
- Vulnerability Scanning
- IAM-Based Access Controls

---

# Challenges Encountered

During implementation, several real-world deployment issues were identified and resolved:

- Missing environment variables
- Docker image naming mismatches
- Docker Hub authentication failures
- Container startup dependency issues
- PostgreSQL readiness timing failures
- Security Group networking issues

Resolving these challenges provided practical experience with infrastructure troubleshooting and production deployment debugging.

---

# Future Architecture Roadmap

This project serves as the foundation for a larger cloud-native platform.

Planned Improvements:

```text
Current
│
├── Docker Compose
├── GitHub Actions
└── AWS EC2

Future
│
├── Terraform
├── Kubernetes
├── ArgoCD
├── Prometheus
├── Grafana
├── AWS Load Balancer
├── Multi-Environment Deployments
├── GitOps
└── DevSecOps Integration
```

---

# Architecture Summary

The system demonstrates:

- Three-Tier Application Design
- Containerization with Docker
- Container Orchestration with Docker Compose
- CI/CD Automation using GitHub Actions
- AWS Cloud Deployment
- Secret Management
- Persistent Data Storage
- Service Health Monitoring
- Production-Style Deployment Workflows

This architecture provides a scalable foundation for future migration toward Kubernetes, GitOps, Infrastructure as Code, and cloud-native deployment patterns.