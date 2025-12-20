# Example FX Framework

This project demonstrates the use of the NG framework with Uber's `fx` and Fiber, mirroring the structure and philosophy of NestJS. It showcases how to build modular, scalable, and maintainable applications in Go.

## Features

- **Modular Design**: Components like `users` and `orders` are self-contained modules.
- **Dependency Injection**: Powered by Uber's `fx` for managing dependencies.
- **Lifecycle Management**: Hooks for application startup and shutdown.
- **Fiber Integration**: High-performance HTTP server integration.
- **DTOs and Models**: Clear separation of data transfer objects and business models.
- **Inspired by NestJS**: Brings the modular and organized structure of NestJS to Go.

## Project Structure

```
example/fx/
├── adapter/          # Adapters for Fiber
├── components/       # Modular components (e.g., users, orders)
│   ├── orders/       # Orders module
│   │   ├── dtos/     # Data Transfer Objects for orders
│   │   ├── order.controller.go
│   │   ├── order.module.go
│   │   ├── order.service.go
│   ├── users/        # Users module
│       ├── dtos/     # Data Transfer Objects for users
│       ├── user.controller.go
│       ├── user.module.go
│       ├── user.service.go
├── models/           # Business models (e.g., Order, User)
├── router/           # Router setup and grouping
├── main.go           # Application entry point
```

## Getting Started

### Prerequisites

- Go 1.18 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/example-fx.git
   ```
2. Navigate to the `example/fx` directory:
   ```bash
   cd example/fx
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Running the Application

To start the application, run:

```bash
go run main.go
```

The server will start on `http://localhost:8080`.

### Adding a New Module

1. Create a new directory under `components/` (e.g., `products`).
2. Add `module.go`, `controller.go`, `service.go`, and `dtos/` as needed.
3. Register the module in `main.go`:
   ```go
   fx.New(
       ...existing code...
       products.Module,
       ...existing code...
   )
   ```

## About NG Framework

NG is a Go framework inspired by NestJS, designed to simplify the development of modular and scalable applications. It provides features like middleware, guards, interceptors, and dynamic adapters for various HTTP frameworks. This example demonstrates how NG can be used with `fx` for dependency injection and Fiber for HTTP server handling.

## License

This project is licensed under the MIT License. See the [LICENSE](../../LICENSE) file for details.
