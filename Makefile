
# Migrate
# name (migration name) ex: make migrate-create name=create_table_user
migrate-create:
	migrate create -ext sql -dir ./migration $(name)

migrate-up:
	go run migration/migrate.go up

migrate-down:	
	go run migration/migrate.go down

unit-test:
	go test -coverprofile=coverage.out \
		./internal/v1/attendance/delivery/... \
		./internal/v1/attendance/usecase \
		./internal/v1/attendance/repository \
		./internal/v1/audit/repository \
		./internal/v1/auth/delivery/... \
		./internal/v1/auth/usecase \
		./internal/v1/auth/repository \
		./internal/v1/compensation/delivery/... \
		./internal/v1/compensation/usecase \
		./internal/v1/compensation/repository \
		./internal/v1/employee/usecase \
		./internal/v1/employee/repository \
		./internal/v1/payroll/delivery/... \
		./internal/v1/payroll/usecase \
		./internal/v1/payslip/repository  && go tool cover -func=coverage.out


tidy:
	go mod tidy

up:
	docker-compose up