# Datastation Backend

## Overview

Padt.AI Backend powers copy-trading and market intelligence workloads. The codebase now focuses on two core trading rails:

- **CEX service** for managing custodial copy-trade wallets, author subscriptions, and user portfolio metadata.
- **DEX service** for non-custodial integrations where users connect trading accounts across supported decentralized exchanges.

Both services share the same authentication, CRM, and notification stack, while exposing dedicated REST endpoints under `/api/v2`.

## Features

- Market summary processing and retrieval
- Integration with multiple databases (MySQL, TimescaleDB, PlanetScale, Postgres)
- RESTful API endpoints for copy-trade operations across CEX and DEX integrations
- Author subscription workflows with per-wallet risk controls (TP/SL, leverage, holding period)

## Prerequisites

- Go (version X.X or higher)
- MySQL
- TimescaleDB
- PlanetScale (optional)

## Configuration

- All credentials committed in this repository are **mock placeholders**. Provide your own secrets through environment variables or secret managers before running anything in non-local environments.
- Copy `example.yaml` to `config.yaml` (or provide the equivalent env vars) and replace the mock values with your actual configuration.
- Docker assets (`Dockerfile`, `docker-compose.yaml`) default to mock usernames/passwords. Override them at build/run time, for example:
  ```bash
  docker build --build-arg GITHUB_USERNAME=real-user --build-arg GITHUB_TOKEN=real-token -t {reponame} .
  docker compose --env-file .env.local up
  ```
  where `.env.local` holds the real DSNs/passwords.

Ensure the following configuration sections are supplied (either via YAML or env vars):

- `app`: General application settings
- `stock_db`: MySQL database for stock data
- `analytic_db`: MySQL database for analytics
- `app_db`: PlanetScale database connection
- `timescale_db`: TimescaleDB connection details

## CEX Service

The CEX service (`internal/core/service/cex_service.go`) orchestrates centralized exchange copy-trade wallets. Key capabilities include:

- Listing wallet metadata, balances, and author subscriptions for authenticated users.
- Managing wallet lifecycle actions (activate/deactivate, rename, priority updates).
- Updating position size, leverage, and stop-loss settings per wallet.

API surface: `/api/v2/cex/*` (see `api/v2/cex.go`). The handler layer depends on `port.CexService` and the repository implementations in `internal/core/adapter/repo/cex_repo.go`.

### Recommended configuration

- Ensure `crypto_db` points to a Postgres instance with the crypto schema.
- Provide CEX-specific secrets (API keys, trading tokens) through secure env vars or secret managers.

## DEX Service

The DEX service (`internal/core/service/dex_service.go`) provides non-custodial wallet connectivity and risk management for decentralized exchanges.

Highlights:

- Connect user wallets by validating API/private-key pairs and auto-deriving addresses when possible.
- Manage wallet activation, author subscriptions, and trade guardrails (leverage, SL, holding period).
- Share wallet-info structures with other CRM components via `model.WalletInfo`.

API surface: `/api/v2/dex/*` (see `api/v2/dex.go`). Storage is handled through `internal/core/adapter/repo/dex_repo.go` with credential validation hooks under `port.DexRepo`.

### Recommended configuration

- Populate `crypto_db` with the DEX wallet tables (`crypto_copytrade_wallet_dex`, author relations, etc.).
- Supply `crypto_trading_bot.baseURL` and `crypto_trading_bot.token` if you proxy credential validation through an external bot service.
- Review swagger docs (`docs/swagger.yaml`) for request/response shapes, keeping examples aligned with your target exchange.

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/your-username/reponame.git
   ```

2. Navigate to the project directory:
   ```
   cd datastation-backend
   ```

3. Install dependencies:
   ```
   go mod tidy
   ```

4. Set up your `config.yaml` (or env vars) with the appropriate database credentials and application settings.

## Running the Application

5. To start the server, run:

    ```
    go run cmd/api/main.go
    ```

## API Endpoints

6. To regenerate the Swagger documentation, run:
   ```
   swag init -g cmd/api/main.go -o docs
   ```
   http://localhost:8080/swagger/index.html


