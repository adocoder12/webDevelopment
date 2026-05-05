# Go Web Development Learning

A dedicated repository documenting my journey and progress in learning the Go (Golang) programming language, with a focus on building a real-world web application from scratch.

## 🎯 Learning Objectives

This project tracks my hands-on experience with the Go standard library and common web patterns:

* **HTTP Foundations:** Utilizing the `net/http` package to build robust servers with timeouts and graceful shutdown.
* **Request Handling:** Implementing handlers, muxers, middleware chains, and processing URL parameters.
* **Data Layer:** Building a clean repository pattern with SQLite, transactions, and a `DBTX` interface for testability.
* **Authentication:** Password hashing with `bcrypt`, session-based auth, and session fixation protection.
* **Sessions:** SQLite-backed sessions via `alexedwards/scs` with flash messages and auth middleware.
* **Templating:** A cached, thread-safe template renderer with layout and partial support.
* **Middleware:** Structured logging with response status capture, panic recovery, and `AuthRequired` guards.

## 📁 Project Structure

```
.
├── cmd/
│   └── main.go                  # Entry point — wiring, DB init, server config
├── internal/
│   ├── handler/
│   │   ├── application.go       # Application struct and constructor
│   │   ├── handler.go           # RegisterHandler, LoginHandler, LogoutHandler
│   │   ├── middlewares.go       # Logger, RecoverPanic, AuthRequired
│   │   ├── routes.go            # SetupRoutes — middleware chain + mux
│   │   ├── pages.go             # Home, About
│   │   └── templates.go        # TemplateRenderer, TemplateData, renderTemplate
│   ├── repository/
│   │   ├── user_repository.go   # UserRepository interface, schemas, ErrNotFound
│   │   └── sql_user_repository.go # SQLite implementation
│   └── model/
│       └── user.go              # User and Profile structs
└── templates/
    └── html/
        ├── layouts/             # base.layout.html
        ├── partials/            # Reusable HTML partials
        ├── index.html
        ├── login.html
        └── register.html
```

## 🛠 Setup and Usage

### Prerequisites

* **Go 1.22+** (required for method-based routing e.g. `GET /path`)
* **GCC** (required by `go-sqlite3` for CGO compilation)

### Installation & Execution

1. **Clone the repository:**
    ```bash
    git clone https://github.com/adocoder12/webDevelopment.git
    ```

2. **Navigate to the folder:**
    ```bash
    cd webDevelopment
    ```

3. **Install dependencies:**
    ```bash
    go mod tidy
    ```

4. **Run from the project root:**
    ```bash
    go run ./cmd/main.go
    ```

   > ⚠️ Always run from the project root. The SQLite database (`app.db`) and template paths are relative to where you execute the command.

5. **Open in browser:**
    ```
    http://localhost:8280
    ```

## 📝 Roadmap

- [x] Basic server setup with timeouts
- [x] Route mapping with Go 1.22 method syntax
- [x] Structured logging with response status capture
- [x] Panic recovery middleware
- [x] SQLite database with foreign keys
- [x] Repository pattern with `DBTX` interface
- [x] User registration and login
- [x] Password hashing with bcrypt
- [x] Session management (SQLite-backed)
- [x] Flash messages
- [x] Auth middleware (`AuthRequired`)
- [x] Thread-safe template renderer with caching
- [ ] Protected dashboard route
- [ ] User profile page
- [ ] Avatar upload
- [ ] HTTPS + secure cookies for production

## 👤 Author

**adocoder12**
* GitHub: [@adocoder12](https://github.com/adocoder12)

---
*Built with the goal of mastering idiomatic Go for the web.*