# Example: Gin Server with NG Framework

This example demonstrates how to create a Gin-based HTTP server using the NG framework. The server includes a simple `/hello` endpoint that responds with a JSON message.

## Features

- HTTP server with a `/hello` endpoint.
- Uses the NG framework for application structure and routing.
- Demonstrates the use of the Gin adapter for route handling.

## File Structure

- `main.go`: Contains the main application logic and server setup.
- `adapter/gin.go`: Provides the Gin adapter for the NG framework.

## How to Run

1. Ensure you have Go installed on your system.
2. Navigate to the `example/gin` directory.
3. Run the following command to start the server:

   ```bash
   go run main.go
   ```

4. The server will start on `http://localhost:8080`.
5. Test the `/hello` endpoint by running:

   ```bash
   curl http://localhost:8080/hello
   ```

   You should receive the following response:

   ```json
   { "code": "OK", "data": "Hello, World!" }
   ```

## Dependencies

- [NG Framework](https://github.com/foxie-io/ng): Used for application structure and routing.
- [Gin Web Framework](https://github.com/gin-gonic/gin): Provides the HTTP server and routing capabilities.

## Notes

- The server uses the Gin adapter for NG framework integration.
- The `/api` prefix is added to all routes.
