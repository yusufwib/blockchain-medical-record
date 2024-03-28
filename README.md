# Mekari Backend Golang Technical Test

This document provides instructions on how to set up, run, and use the Mekari Backend Golang application, which utilizes the Echo framework, PostgreSQL, Docker, Go-Wire for dependency injection, Swagger for API documentation, and follows a clean architecture pattern and SOLID principles.


## Getting Started

Ensure you have Docker installed on your machine.
Open a terminal and navigate to the project root directory.
Run the following command to start the application and PostgreSQL container:
```bash
docker-compose up
```
This command reads the docker-compose.yml file and sets up the PostgreSQL and Golang API containers.

The application will be accessible at http://localhost:9009

## API Endpoints

**Base URL:** http://localhost:9009

**Swagger UI:** http://localhost:9009/swagger/

**Employees Endpoints:**

-   **GET /employees:**
    -   Retrieves a list of all employees.
    -   Input: None
    -   Output: JSON array of employee objects
-   **GET /employees/:id:**
    -   Retrieves a single employee by ID.
    -   Input: Path parameter  `id`  (employee ID)
    -   Output: JSON object representing the employee
-   **POST /employees:**
    -   Creates a new employee.
    -   Input: JSON object with employee data
    -   Output: JSON object representing the created employee with its assigned ID
-   **PUT /employees/:id:**
    -   Updates an existing employee by ID.
    -   Input: Path parameter  `id`  (employee ID) and JSON object with updated employee data
    -   Output: JSON object representing the updated employee
-   **DELETE /employees/:id:**
    -   Deletes an employee by ID.
    -   Input: Path parameter  `id`  (employee ID)
    -   Output: Empty response with a 200 status code if successful

## Additional Notes

-   **Environment Variables:**  The API relies on environment variables to configure database connection details. These are set in the  `docker-compose.yml`  file.
-   **Migrations:**  The project includes database migrations to create and manage the schema.
-   **Clean Architecture:**  The project adheres to clean architecture principles, promoting modularity and testability.

## Troubleshooting

-   If you encounter any issues, check the logs for both the Golang API and PostgreSQL containers using  `docker-compose logs`.
-   Ensure that you have the correct environment variables set in your  `docker-compose.yml`  file.
-   Verify that the database schema is up-to-date by running any necessary migrations.

## Project Structure

```plaintext
.
├── Dockerfile.api
├── Makefile
├── README.md
├── config
│   ├── config.go
│   └── env.example
├── datasource
│   ├── migrations
│   │   ├── 000000_employee_down.sql
│   │   └── 000000_employee_up.sql
│   └── pgsql.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── handler
│   ├── employee.go
│   └── response.go
├── infrastructure
│   ├── application.go
│   ├── dependency.go
│   ├── http
│   │   ├── main.go
│   │   ├── middleware.go
│   │   └── router.go
│   └── wire_gen.go
├── main.go
├── models
│   ├── demployee
│   │   ├── employee.go
│   │   ├── request.go
│   │   └── response.go
│   └── server
│       └── server.go
├── repository
│   ├── employee_repository.go
│   └── employeee_repository_test.go
├── service
│   └── employee_service.go
└── utils
    ├── envar
    │   └── envar.go
    ├── logger
    │   ├── logger.go
    │   └── mlog.go
    ├── mvalidator
    │   └── mvalidator.go
    ├── slicer
    │   ├── check.go
    │   ├── check_test.go
    │   ├── chunk.go
    │   └── chunk_test.go
    └── trace_id
        └── trace_id.go
```