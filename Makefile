
# Migrate
# name (migration name) ex: make migrate-create name=create_table_user
migrate-create:
	migrate create -ext sql -dir ./migration $(name)

migrate-up:
	go run migration/migrate.go $(db) up

migrate-down:	
	go run migration/migrate.go $(db) down

