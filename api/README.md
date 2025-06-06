
## Description

This project is structured to provide a clean and modular Go application. Below is a brief description of the key components:

## Project Structure

```text
gonuxt-context-assistant/
├── cmd/                          <-- For main applications
│   └── api/                      <-- Our HTTP API server
│       └── main.go               <-- Entry point for the API
├── internal/                     <-- Private application code (not importable by external projects)
│   ├── app/                      <-- Core application logic
│   │   └── assistant/            <-- Contains our assistant's core logic (e.g., orchestrator)
│   │       └── assistant.go
│   ├── api/                      <-- Internal API-specific components (e.g., handlers, routes, request/response models)
│   │   └── handler.go
│   │   └── models.go
│   ├── tools/                    <-- Our helper tools (already exists)
│   │   └── tools.go
│   └── config/                   <-- Application configuration
│       └── config.go
├── pkg/                          <-- Public library code (potentially reusable by other projects)
│   └── errors/                   <-- Custom error types (optional, but good for structured errors)
│       └── errors.go
├── go.mod
├── go.sum
```

## Usage

To get started, clone the repository and navigate to the project directory. Ensure you have Go installed and properly configured on your system.

```bash
git clone <repository-url>
cd gonuxt-context-assistant
go run main.go