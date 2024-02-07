# EchoSync-Commerce

EchoSync-Commerce is a simple e-commerce platform built in Golang using microservice architecture. It allows users to create their own stores to sell products.

## Services

### UserService

- **Transport:** gRPC
- **Database:** PostgreSQL
- **Authorization:** JWT tokens

### MarketService

- **Transport:** gRPC
- **Database:** PostgreSQL

### Gateway

- **Description:** Provides API Gateway to all services based on gRPC
- **Router:** Fiber
- **Middleware:** JWT

## Endpoints

- **Auth Service Endpoints**
  - `/auth/sign-up`: POST method to sign up
  - `/auth/sign-in`: GET method to sign in
  - `/auth/refresh`: GET method to refresh token

- **Market Service Endpoints**
  - `/market/store`: POST method to create a store

- **Product Service Endpoints**
  - `/product/create`: POST method to create a product


## Request Structures

### CreateStoreRequest

```json
{
  "name": "Example Store",
  "owner_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

### SignUpRequest

```json
{
  "username": "example_user",
  "email": "user@example.com",
  "password": "password123"
}
```

### SignInRequest
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```



