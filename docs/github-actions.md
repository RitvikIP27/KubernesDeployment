# GitHub Actions CI/CD Documentation

## Overview

This project implements a complete CI/CD pipeline using GitHub Actions to automate:

- Source Code Validation
- Docker Image Build
- Docker Hub Publishing
- AWS EC2 Deployment
- Container Orchestration using Docker Compose

The application follows a three-tier architecture:

```text
Frontend (Nginx)
        │
        ▼
Backend (Go API)
        │
        ▼
PostgreSQL/Supabase Database
```

---

# CI Pipeline

Workflow File:

```text
.github/workflows/ci.yml
```

## Purpose

The Continuous Integration (CI) workflow is responsible for:

- Checking out source code
- Setting up Docker Buildx
- Authenticating with Docker Hub
- Building Docker images
- Publishing Docker images to Docker Hub

---

## CI Flow

```text
Developer Push
        │
        ▼
GitHub Actions Runner
        │
        ▼
Checkout Repository
        │
        ▼
Setup Docker Buildx
        │
        ▼
Docker Hub Login
        │
        ▼
Build Docker Image
        │
        ▼
Push Image to Docker Hub
```

---

## Workflow Breakdown

### Step 1: Checkout Source Code

```yaml
uses: actions/checkout@v4
```

This action:

- Clones repository code into the GitHub Actions runner
- Makes the source code available for subsequent workflow steps

Every workflow execution runs on a fresh runner instance.

---

### Step 2: Setup Docker Buildx

```yaml
uses: docker/setup-buildx-action@v3
```

Buildx provides:

- Multi-platform image builds
- Improved build caching
- Advanced Docker build capabilities

---

### Step 3: Authenticate with Docker Hub

```yaml
uses: docker/login-action@v3
```

Credentials are securely retrieved from GitHub Secrets.

```yaml
username: ${{ secrets.DOCKERHUB_USERNAME }}
password: ${{ secrets.DOCKERHUB_PASS }}
```

No credentials are stored in source control.

---

### Step 4: Build and Push Docker Image

```yaml
uses: docker/build-push-action@v6
```

Configuration:

```yaml
context: ./backend
push: true
tags: ${{ secrets.DOCKERHUB_USERNAME }}/k8app:latest
```

This step:

- Builds the backend application image
- Tags the image
- Pushes the image to Docker Hub

Example:

```text
ritvikkantip/k8app:latest
```

---

## CI Workflow

```yaml
name: Ritvik's CI Workflow

on:
  push:

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout Code
        uses: actions/checkout@v4

      - name: Setup Docker
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_PASS }}

      - name: Build and Push Docker Image
        uses: docker/build-push-action@v6
        with:
          context: ./backend
          push: true
          tags: ${{ secrets.DOCKERHUB_USERNAME }}/k8app:latest
```

---

# CD Pipeline

Workflow File:

```text
.github/workflows/cd.yml
```

## Purpose

The Continuous Deployment (CD) workflow is responsible for:

- Triggering after successful CI completion
- Transferring deployment assets to AWS EC2
- Injecting runtime configuration
- Pulling the latest Docker images
- Deploying updated containers

---

## CD Flow

```text
Successful CI Run
        │
        ▼
Trigger CD Workflow
        │
        ▼
Checkout Repository
        │
        ▼
Copy Files to EC2
        │
        ▼
SSH Into EC2
        │
        ▼
Inject Secrets
        │
        ▼
Docker Compose Pull
        │
        ▼
Docker Compose Up
        │
        ▼
Deployment Verification
```

---

## Deployment Process

### Step 1: Trigger After Successful CI

```yaml
workflow_run:
```

The CD pipeline executes only when:

```text
CI Workflow
      ↓
Status = Success
```

---

### Step 2: Transfer Deployment Assets

```yaml
uses: appleboy/scp-action@v0.1.7
```

Files are copied from GitHub Actions runner to EC2.

Target:

```text
/home/ubuntu/KubernesDeployment
```

---

### Step 3: SSH Into EC2

```yaml
uses: appleboy/ssh-action@v1.0.3
```

Authentication is performed using:

```text
EC2_HOST
EC2_USER
EC2_SSH_KEY
```

stored securely in GitHub Secrets.

---

### Step 4: Runtime Secret Injection

Application secrets are injected during deployment.

```bash
cat > .env << EOF
DOCKERHUB_USERNAME=${{ secrets.DOCKERHUB_USERNAME }}
DATABASE_URL=${{ secrets.DATABASE_URL }}
DB_NAME=${{ secrets.DB_NAME }}
DB_USER=${{ secrets.DB_USER }}
DB_PASSWORD=${{ secrets.DB_PASSWORD }}
EOF
```

### Why?

- Secrets are never committed to Git
- Sensitive values remain encrypted in GitHub Secrets
- Deployment remains environment-agnostic

---

### Step 5: Deploy Containers

```bash
docker compose pull
docker compose up -d
```

This:

- Pulls latest images
- Updates containers
- Preserves database volumes
- Applies application updates

---

### Step 6: Deployment Validation

```bash
sleep 30
docker compose ps
```

Used to verify:

```text
backend   Up
db        Up (healthy)
nginx     Up
```

---

# Secret Management

The following secrets are stored in GitHub Actions:

| Secret | Purpose |
|----------|----------|
| DOCKERHUB_USERNAME | Docker Hub Username |
| DOCKERHUB_PASS | Docker Hub Token |
| EC2_HOST | EC2 Public IP |
| EC2_USER | SSH User |
| EC2_SSH_KEY | SSH Private Key |
| DATABASE_URL | Supabase/Postgres Connection String (optional) |
| DB_NAME | Application Database |
| DB_USER | Database User |
| DB_PASSWORD | Database Password |

---

# Lessons Learned

During implementation, several real-world deployment issues were encountered and resolved:

- Docker image naming mismatches
- Missing environment variables
- Invalid Docker image references
- Docker Hub authentication issues
- PostgreSQL health check timing failures
- AWS Security Group configuration issues

Resolving these issues provided hands-on experience with production-style troubleshooting and deployment debugging.

---

# Future Roadmap

This project serves as the foundation for a larger platform that will include:

- Kubernetes Deployment
- Terraform Infrastructure Provisioning
- ArgoCD GitOps
- Prometheus Monitoring
- Grafana Dashboards
- AWS Load Balancer Integration
- Automated Infrastructure Deployment
- Multi-Environment Deployments

---

# Final Architecture

```text
Developer
    │
    ▼
GitHub Repository
    │
    ▼
GitHub Actions CI
    │
    ├── Checkout Code
    ├── Build Docker Image
    └── Push To Docker Hub
    │
    ▼
GitHub Actions CD
    │
    ├── SCP To EC2
    ├── SSH Deployment
    ├── Secret Injection
    └── Docker Compose Deployment
    │
    ▼
AWS EC2
    │
    ├── Nginx
    ├── Go Backend
    └── PostgreSQL/Supabase
```

# Infrastructure Pipeline

A dedicated infrastructure workflow has been introduced.

Purpose:

- Terraform formatting validation
- Terraform configuration validation
- Infrastructure planning

Current Workflow:

Terraform fmt
↓
Terraform validate
↓
Terraform plan

Future versions will include:

- Terraform apply
- Terraform destroy
- Remote state management
- Infrastructure approval gates