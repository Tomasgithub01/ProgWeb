test:
	hurl --test script.hurl 

testdelete:
	hurl --test delete.hurl

dev: 
	docker compose -f docker-compose.dev.yml up --build

prod:
	docker compose up --build 

down:
	docker compose -f docker-compose.dev.yml down -v

wait-for-db: 
	@echo "Waiting to database to be ready..."
	@until docker exec tpespecial_db pg_isready -U user -d mydatabase_dev; do \
		sleep 1; \
	done
	@echo "Database is ready!"

wait-for-building:
	@echo "Esperando que la API esté disponible..."
	@until curl -s http://localhost:8080/health >/dev/null; do \
		printf '.'; \
		sleep 2; \
	done
	@echo "\nAPI lista!"

testdev:
	docker compose -f docker-compose.dev.yml up -d --build
	make wait-for-db
	make wait-for-building
	hurl --test script.hurl

templ: 
	go run github.com/a-h/templ/cmd/templ@latest generate

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

users:
	curl -X GET http://localhost:8080/users

generate:
	make templ
	make sqlc

health:
	@curl -f http://localhost:8080/health