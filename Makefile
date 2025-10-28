APP_NAME=my-app

DB_URL="postgres://admin:Admin20250915@localhost:5434/tpespecialweb?sslmode=disable"

# Significa que si se ejecuta 'make' por defecto es 'make build'
all: build

# 'make run' ejecuta el air, antes levanta la base de datos. Luego air ejecuta 'make build' cuando ve un cambio.
run: setup
	air

generate:
#	echo ">= Generating SQLC code..."
#	sqlc generate 
#	@echo ">= Generating Templ code..."
#	@templ generate

#db/migrate:

# Cada vez que air vea un cambio en un archivo ejecuta 'make generate' y 'make build'
build: generate
	echo ">= Building application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) .

clean:
	@rm -f $(APP_NAME)

# Levantar solo la base de datos
up_db_and_wait:
	@echo "Iniciando servicio 'db'..."
	docker compose up -d db

setup: up_db_and_wait
	@echo "Esperando a que 'db' esté saludable (healthy)..."
	@while [ "$$(docker inspect -f '{{.State.Health.Status}}' tpespecial_db 2>/dev/null)" != "healthy" ]; do \
		echo "Aún no está lista..."; \
		sleep 2; \
	done
	@echo "\033[1;32m[OK] Base de datos lista!\033[0m"

stop-db: 
	docker compose stop db 