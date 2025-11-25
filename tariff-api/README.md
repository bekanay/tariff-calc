# Tariff API Endpoints

HTTP API served by `cmd/main.go` on port `8081`. Auth flows use Keycloak; all station routes require the configured `REQUIRED_ROLE`. The `Authorization` header is expected to contain the raw access token (no `Bearer` prefix).

## Auth

### POST `/login`
- Body:
```json
{
  "username": "alice",
  "password": "secret"
}
```
- Success `200`:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "id_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_in": 1800
}
```
- Errors: `400` invalid body, `401` invalid credentials, `403` missing required role, `500` on server issues.

### POST `/refresh-token`
- Body:
```json
{ "refresh_token": "<refresh token>" }
```
- Success `200`: same shape as `/login`.
- Errors: `400` invalid body, `401` if refresh fails.

### POST `/logout`
- Body:
```json
{ "refresh_token": "<refresh token>" }
```
- Success `200`:
```json
{ "message": "Logged out successfully" }
```
- Errors: `400` invalid body, `500` on server issues.

### GET `/roles/{role}`
- Headers: `Authorization: <access token>`
- Success `200`:
```json
{ "role": "admin", "has_role": true }
```
- Errors: `400` missing role param, `401` missing token, `500` on token decode.

## Stations (protected; requires `Authorization` header with access token)

### GET `/stations`
- Query params:
  - `page` (default `1`), `page_size` (default `20`, max `200`)
  - `sort` in `[id,-id,stan_kod,-stan_kod,stan_name,-stan_name]` (default `id`)
  - `name` substring match (case-insensitive)
- Success `200`:
```json
{
  "stations": [
    {
      "id": 1,
      "stan_kod": "1234",
      "stan_name": "Central",
      "stan_priznak": 1,
      "paragraph": "A1"
    }
  ],
  "metadata": {
    "current_page": 1,
    "page_size": 20,
    "first_page": 1,
    "last_page": 1,
    "total_records": 1
  }
}
```
- Errors: `400` invalid filters, `401/403` missing or insufficient role, `500` server issues.

### GET `/stations/{kod}`
- Success `200`:
```json
{
  "id": 1,
  "stan_kod": "1234",
  "stan_name": "Central",
  "stan_priznak": 1,
  "paragraph": "A1"
}
```
- Errors: `400` invalid code, `401/403` auth/role failure, `404` not found.

### POST `/stations`
- Body:
```json
{
  "stan_kod": "1234",
  "stan_name": "Central",
  "stan_priznak": 1,
  "paragraph": "A1"
}
```
- Success `201`: returns created station with `id`.
- Errors: `400` invalid data, `401/403` auth/role failure, `409` duplicate code, `500` server issues.

### PUT `/stations/{kod}`
- Body (replaces data for the station found by `{kod}`):
```json
{
  "stan_kod": "5678",
  "stan_name": "New Name",
  "stan_priznak": 2,
  "paragraph": "B2"
}
```
- Success `200`: updated station.
- Errors: `400` invalid data, `401/403` auth/role failure, `404` not found, `409` duplicate `stan_kod`.

### DELETE `/stations/{kod}`
- Success `204` (no body).
- Errors: `400` invalid code, `401/403` auth/role failure, `404` not found, `500` server issues.
