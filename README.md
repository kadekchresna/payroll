# payroll

This is a backend service for managing payroll, attendance, and compensation for employees. It is designed by implementing **Separation of Concerns** and **Domain-Driven Design (DDD)** principles to promote scalability, testability, and maintainability.

## Features

- Employee attendance tracking
- Overtime and reimbursement management
- Payslip and payroll generation
- JWT-based authentication
- User role seperation
- GORM-powered database access layer
- Unit test support with mocking

## Tech Stack

| Layer            | Tech/Library                                                                                  |
| ---------------- | --------------------------------------------------------------------------------------------- |
| Language         | Go (Golang)                                                                                   |
| Web Framework    | [Echo v4](https://echo.labstack.com)                                                          |
| ORM              | [GORM](https://gorm.io/)                                                                      |
| Testing          | [Testify](https://github.com/stretchr/testify) + [Mockery](https://github.com/vektra/mockery) |
| DB Migration     | Native SQL with custom runner (`migrate.go`)                                                  |
| Auth             | JWT (JSON Web Tokens)                                                                         |
| Logging          | [Zap](https://github.com/uber-go/zap)                                                         |
| Containerization | Docker + Docker Compose                                                                       |

## Folder Structure

```bash
internal/v1/
├── attendance/              # Attendance domain
    ├── model/
    ├── repository/
    ├── usecase/
    └── delivery/
├── audit/                   # Auditing and logging
    └── ...
├── auth/                    # Authentication logic
    └── ...
├── compensation/            # Overtime and reimbursement
    └── ...
├── employee/                # Employee data
    └── ...
└── payroll/                 # Payroll and salary logic
    └── ...

infrastructure/
└── db/                      # DB connection & transaction helpers

helper/
├── env/                     # Env config loader
├── jwt/                     # JWT utilities and middleware
├── logger/                  # Logging and request tracing
└── time/                    # Time parsing and formatting

migration/                   # SQL schema setup
cmd/                         # Main service entry point
```

This project dividing code into layers:

- Model (Domain) – pure business objects and logic
- DTO (Application) – for external data exchange
- Repository – actual DB queries using GORM
- Usecase – orchestrates business logic
- Delivery – handles routing and request parsing

## Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## How to Run the Service locally

1. Makesure **Docker** and **Docker Compose** is already setup on your local

2. Deploy the **payroll service** and PostgreSQL

```
    make up
```

3. Access the service through `:8081` with this [postman collection]()

## How to Run Unit Coverage

1. Run unit test coverage on main functionality

```
make unit-test
```

## User Credential

1. Login as Employee

```
    "username": "user1"
    "password": "secret123"
```

2. Login as Admin

```
    "username": "admin"
    "password": "1234567890"
```

## Note

1. Each endpoint must have `Authorization` header in order to access the enpoint

```
Authorization:Bearer xxxxx.xxxxxxxxxxxxx.xxxxxxx
```

2. In order the employee to able to submit overtime, they have to clocked out first by submit attendance second time at the same specified date
