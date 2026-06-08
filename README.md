# Seismic Monitor

**Live Demo / Production:** [http://129.212.161.29](http://129.212.161.29)

**Seismic Monitor** is a full-stack application designed for monitoring seismic events. It provides an interface to visualize recent earthquakes, register users, configure alerts, and view report history. The project uses a **Go** backend and a frontend developed with **Vue.js and Vite**, backed by a **PostgreSQL** database.

This repository is a monorepo containing both the frontend and backend code, as well as end-to-end (E2E) tests.

---

## Prerequisites

To run and collaborate on this project, you need the following tools installed in your local environment:

- **[Go](https://go.dev/dl/)** (v1.24 or higher) - For the backend.
- **[Node.js](https://nodejs.org/)** (v22 or higher) - For the frontend and E2E tests.
- **[PostgreSQL](https://www.postgresql.org/)** - Main database.
- **[Docker](https://www.docker.com/) & Docker Compose** (Optional, but recommended) - To run the database locally.
- **[Make](https://www.gnu.org/software/make/)** - For quick execution of common development commands. _(Note: Windows users may need to install `make` via Winget/Chocolatey or run the commands within the Makefile manually)._

---

## Local Setup and Execution

### 1. Clone the repository

```bash
git clone https://github.com/your-username/seismic-monitor.git
cd seismic-monitor
```

### 2. Configure Environment Variables

The project requires environment variables to connect to the database and other external APIs.
You must configure the corresponding variables (for example, by copying `.env.example` to `.env` or exporting the variables locally) for the backend environment.

Ensure you correctly define the database connection, secrets for JWT token generation, and any necessary API keys (like Gemini or seismic data providers).

_Database configuration example:_

```env
DB_URL=postgres://user:password@localhost:5432/seismic_db?sslmode=disable
```

### 3. Start Local Infrastructure

You can quickly start the necessary services (like PostgreSQL) using Docker Compose:

```bash
docker-compose up -d
```

### 4. Install Dependencies

From the project root, run the following command to install Go development tools (like `air` and `golangci-lint`) and Node.js dependencies:

```bash
make install
```

### 5. Run the Application

Open two terminals to run the concurrent development environment:

**Terminal 1 (Backend):**
Start the Go server with hot reload:

```bash
make run-backend
```

The backend server will start and expose its API routes.

**Terminal 2 (Frontend):**
Start the Vite development server:

```bash
make run-frontend
```

The frontend will be available at `http://localhost:5173`. Requests to `/api` are automatically proxied to the backend via Vite.

---

## Running Tests

The project includes different test suites. You can run them as follows:

- **Unit and Integration Tests (Go & Frontend):**

  ```bash
  make test
  ```

  _(Runs Go and frontend tests)_

- **End-to-End (E2E) Tests with Playwright:**
  To run E2E tests, ensure both the backend and frontend are running first.

  ```bash
  make e2e
  ```

- **Code Linter:**
  To ensure code quality and conventions, run:
  ```bash
  make lint
  ```

---

### Contribution Workflow

1. **Fork** the repository and/or clone it to your environment.
2. **Create a branch** from `main` with a descriptive name for your new feature or fix:
   ```bash
   git checkout -b feat/new-feature
   # or
   git checkout -b fix/bug-fix
   ```
3. Make your changes and use semantic commits (see below).
4. Push your branch to the remote repository:
   ```bash
   git push origin feat/new-feature
   ```
5. Open a **Pull Request (PR)** against the `main` branch, clearly detailing the changes made, the problem it solves, and how to test it.

### Commit Conventions

Please structure your commit messages to explain the changes made.

---

## Technical Documentation

All the technical documentation for the project is organized within the `/docs` directory.

**[View Technical Documentation (`/docs` Directory)](./docs/)**

To get started, we suggest reviewing:

- [Getting Started Guide](./docs/getting-started.md)
- [Monorepo Structure](./docs/monorepo.md)
- [Backend Development Guide (Golang)](./docs/golang.md)
- [Architecture Decision Records (ADR)](./docs/adr/)
