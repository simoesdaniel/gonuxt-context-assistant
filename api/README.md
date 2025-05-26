
## Description

This project is structured to provide a clean and modular Go application. Below is a brief description of the key components:

- **`go.mod`**: Defines the Go module and its dependencies.
- **`main.go`**: The main entry point of the application.
- **`internal/`**: Contains private packages that are not meant to be imported by other projects.
  - **`tools/`**: A subpackage for utility functions.
    - **`tools.go`**: Implements utility functions such as `GetCurrentDateTime` and `GetWeather`.

## Usage

To get started, clone the repository and navigate to the project directory. Ensure you have Go installed and properly configured on your system.

```bash
git clone <repository-url>
cd gonuxt-context-assistant
go run main.go