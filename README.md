# Chirpy

Chirpy is a simple microblogging API server written in Go. It supports user registration, authentication, posting and managing "chirps" (short messages), and basic admin/dev endpoints.

## Table of Contents

- [Setup](#setup)
- [Environment Variables](#environment-variables)
- [API Endpoints](#api-endpoints)
  - [Health & Metrics](#health--metrics)
  - [User Endpoints](#user-endpoints)
  - [Authentication & Tokens](#authentication--tokens)
  - [Chirp Endpoints](#chirp-endpoints)
  - [Polka Webhook](#polka-webhook)
  - [Development/Admin](#developmentadmin)
- [Data Models](#data-models)
- [Development](#development)

---

## Setup

1. Clone the repository.
2. Install dependencies:
   ```sh
   go mod download