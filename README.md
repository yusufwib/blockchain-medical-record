# Blockchain Medical Record Services

## Getting Started

Ensure you have Docker installed on your machine.
Open a terminal and navigate to the project root directory.
Run the following command to start the application and PostgreSQL container:
```bash
docker-compose up
```
This command reads the docker-compose.yml file and sets up the PostgreSQL and Golang API containers.

The application will be accessible at http://localhost:9009

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
├── Dockerfile.api
├── Makefile
├── README.md
├── config
│   ├── config.go
│   └── env.example
├── datasource
│   ├── migrations
│   │   ├── 000000_medical_records_down.sql
│   │   ├── 000000_medical_records_dummy.sql
│   │   └── 000000_medical_records_up.sql
│   └── pgsql.go
├── docker-compose.yml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
├── handler
│   ├── appointment.go
│   ├── health.go
│   ├── response.go
│   └── user.go
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
│   ├── dappointment
│   │   ├── appointment.go
│   │   ├── request.go
│   │   └── response.go
│   ├── ddoctor
│   │   └── doctor.go
│   ├── dhealthservice
│   │   └── healthservice.go
│   ├── dmedicalrecordaccess
│   │   └── medicalrecordaccess.go
│   ├── dpatient
│   │   └── patient.go
│   ├── duser
│   │   ├── request.go
│   │   ├── response.go
│   │   └── user.go
│   └── server
│       └── server.go
├── repository
│   ├── appointment_repository.go
│   ├── health_repository.go
│   └── user_repository.go
├── service
│   ├── appointment_service.go
│   ├── health_service.go
│   └── user_service.go
└── utils
    ├── envar
    │   └── envar.go
    ├── logger
    │   ├── logger.go
    │   └── mlog.go
    ├── mvalidator
    │   └── mvalidator.go
    ├── randstr
    │   └── randstr.go
    ├── slicer
    │   ├── check.go
    │   ├── check_test.go
    │   ├── chunk.go
    │   └── chunk_test.go
    └── trace_id
        └── trace_id.go
```