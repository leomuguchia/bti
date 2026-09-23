# Zuri Bank (demo)

Small banking transaction backend and UI built for an interview task. Not a production system, a demo showing the core operations working end to end with a real audit trail behind them.

## What it does

- Create an account (basic KYC fields: full name, national ID, phone, address, KRA PIN optional, ID photo, passport photo)
- Log in with national ID + phone number (no password, this is a demo, see notes below)
- Deposit
- Withdraw
- Transfer between two accounts
- Delete a dormant account (zero balance and no activity for 90 days, will not delete an account that still has money)
- Bonus: create a loan account and disburse 10,000 KES into it
- Full transaction ledger per account, append only

## Stack

- Backend: Go, Gin
- DB: SQLite
- Frontend: React, Vite

## Requirements

- Go 1.22 or newer
- Node 18 or newer, npm
- sqlite3 CLI, optional, only useful for poking at the database directly

## Backend setup

From the project root, set these env vars (a `.env` file works if you're loading one, or export them directly):

```
PORT=8080
API_KEY=some-long-random-string
DB_PATH=banking.db
FRONTEND_ORIGIN=http://localhost:5173
```

`API_KEY` has no default and the server refuses to start without it, this is intentional, an empty key would otherwise accept any request.

Install deps and run:

```
go mod tidy
go run .
```

## Frontend setup

```
cd frontend
npm install
```

Create `frontend/.env`:

```
VITE_API_KEY=some-long-random-string
VITE_API_BASE_URL=http://localhost:8080
```

`VITE_API_KEY` has to match `API_KEY` on the backend, this is what the frontend sends on every request in the `Authorization` header.

## Running it

Two ways to run this, depending on what you're doing.

**Dev mode, with hot reload on the frontend:**

Terminal 1, from the root:
```
go run .
```

Terminal 2:
```
cd frontend
npm run dev
```

Frontend runs on `localhost:5173`, backend on `localhost:8080`, CORS is already set up for that split.

**Single server, closer to how it'd actually run:**

```
cd frontend
npm run build
cd ..
go run .
```

This builds the frontend into `frontend/dist`, and the Go server serves those static files directly at `localhost:8080/`, along with the API. One process, one port.

## Notes

Login is national ID + phone number, not a password. There's no password field anywhere in this design. A real bank would gate this with a PIN or OTP, this demo identifies you by the KYC details you gave at account opening instead. Worth knowing before you judge it as a login screen.

ID photo and passport photo are captured in the UI and their filenames get stored, the actual image bytes are never uploaded or persisted anywhere. A real implementation would push those to something like S3 and store a signed URL. Out of scope here.

Known gaps, not hidden: rate limiting exists as a placeholder, not a real implementation. Auth is one shared API key rather than per-user credentials. No TLS locally (that's a deployment concern, not something the app code handles). No idempotency protection on money-moving endpoints yet, meaning a retried request could double-process. That last one is the first thing to build if this went further.

Dormancy threshold (90 days) is hardcoded for the demo. In a real system this would be configurable and probably market/regulator dependent.