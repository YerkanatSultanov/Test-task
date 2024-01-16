migrate-up:
	migrate -path db/migration/ -database "postgresql://postgres:postgres@localhost:5432/test-task-database?sslmode=disable" -verbose up
