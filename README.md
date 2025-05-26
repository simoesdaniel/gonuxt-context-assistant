# gonuxt-context-assistant

## Overview

This is a pet project to learn Go and MCP (Message Control Protocol) servers.

- The **`api`** folder will contain a Go project that serves as an API. This API will be publicly exposed and will include a set of tools to communicate with an LLM (Language Model) that integrates an MCP server.
- The **`app`** folder will contain a simple Nuxt project. It will use a server-side rendering approach to generate pages and explore new Nuxt concepts.

## Project Structure

## Getting Started

### Prerequisites

- Go (latest version)
- Node.js and npm (for the Nuxt app)

### Running the API

Navigate to the `api` folder and run the Go server:

```bash
cd api
go run main.go