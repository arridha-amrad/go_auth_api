# Go Authentication API

A robust authentication API built with Go.

Features include:

✅ Register
✅ Login
✅ Email verification
✅ Resend email verification
✅ Forgot password
✅ Reset password
✅ Get auth info (me)
✅ Logout
✅ Refresh token

## 🔧 Requirements

Make sure the following tools and dependencies are installed:

- Go: >= 1.24.4
- PostgreSQL: Used as the primary database
- Redis: For token caching/session management
- Make: Used to run project commands

  ```sh
  sudo apt install make
  ```

- Air (Live reload for Go):

  ```sh
  go install github.com/air-verse/air@latest
  ```

- mockgen (For generating mock interfaces):

  ```sh
  go install go.uber.org/mock/mockgen@latest
  ```

- migrate (For database migration):

  ```sh
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

#### 🔔 Make sure `$GOPATH/bin` is in your `$PATH`

Usually this is `$HOME/go/bin`. If not already added:

```sh
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.bashrc && source ~/.bashrc
```

## 🧪 Running Tests

Tests rely on PostgreSQL, so you'll need to run the test environment with Docker Compose:

```sh
make compose-test-up
```

## 🧰 Generating Mocks

Mocks are generated using mockgen, and are stored in the `/mocks folder`.
To check if mockgen is installed:

```sh
mockgen -version
```

If not installed:

```sh
go install go.uber.org/mock/mockgen@latest
```

## 📦 Environment Variables

Before running the app, create a .env file in the project root directory with the following variables:

```.env
# Server
PORT="5000"

# PostgreSQL Configuration
DB_URL="postgres://<user>:<password>@localhost:5432/<database>?sslmode=disable"
DB_MAX_OPEN_CONNS=50         # Max number of open DB connections
DB_MAX_IDLE_CONNS=25         # Max number of idle connections
DB_MAX_IDLE_TIME="5m"        # How long a connection can stay idle (e.g. 5m)

# JWT / Application Secret
SECRET_KEY="<your-secret-key>"   # Used for JWT signing

# Google OAuth / Gmail API (used for sending verification emails)
GOOGLE_PROJECT_ID="<your-project-id>"
GOOGLE_CLIENT_ID="<your-client-id>"
GOOGLE_CLIENT_SECRET="<your-client-secret>"
GOOGLE_REFRESH_TOKEN="<your-refresh-token>"

# Application URI (used in email links)
APP_URI="http://localhost:5000"

# Redis Configuration
REDIS_ADDR="localhost:6379"
REDIS_PWD="redis123"             # Password for Redis instance
REDIS_DB=0                       # Redis database number

```

## 🚀 Development

Before execute the dev, make sure you have run db migration

```sh
make db-migrate-up
```

```sh
make run-dev
```
